package auth

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/internal/coreapi/utils/apihelper"
	"github.com/awanku/awanku/pkg/core"
	"github.com/go-chi/chi"
)

const oauthRedirectTo = "/oauth/callback"
const oauthAuthorizationCodeLength = 20
const oauthTokenLength = 20

// @Id api.v1.auth.provider.connect
// @Summary Auth provider connect
// @Tags Auth
// @Param provider path string true "Auth provider" Enums(github, google)
// @Param queryParam query getProviderConnectParam true "Query param"
// @Router /v1/auth/{provider}/connect [get]
// @Produce json
// @Success 301
// @Header 301 {string} location "provider login url"
// @Failure 400 {object} apihelper.HTTPError
// @Failure 401 {object} apihelper.HTTPError
// @Failure 500 {object} apihelper.InternalServerError
func HandleOauthProviderConnect(w http.ResponseWriter, r *http.Request) {
	environment := appctx.Environment(r.Context())

	param := &getProviderConnectParam{
		Provider:   chi.URLParam(r, "provider"),
		RedirectTo: r.URL.Query().Get("redirect_to"),
	}
	if err := param.Validate(r.Context()); err != nil {
		apihelper.ValidationErrResp(w, err)
		return
	}

	state, err := param.encodeState()
	if err != nil {
		apihelper.InternalServerErrResp(w, err)
		return
	}

	authHandler := oauth2Provider(param.Provider, environment)
	apihelper.RedirectResp(w, authHandler.LoginURL(state))
}

// @Id api.v1.auth.provider.callback
// @Summary Auth provider callback
// @Tags Auth
// @Param provider path string true "Auth provider" Enums(github, google)
// @Param queryParam query getProviderCallbackParam true "Query param"
// @Router /v1/auth/{provider}/callback [get]
// @Produce json
// @Success 301
// @Header 301 {string} location "return to url"
// @Failure 400 {object} apihelper.HTTPError
// @Failure 401 {object} apihelper.HTTPError
// @Failure 500 {object} apihelper.InternalServerError
func HandleOauthProviderCallback(w http.ResponseWriter, r *http.Request) {
	environment := appctx.Environment(r.Context())

	param := &getProviderCallbackParam{
		Provider: chi.URLParam(r, "provider"),
		Code:     r.URL.Query().Get("code"),
		State:    r.URL.Query().Get("state"),
	}
	if err := param.Validate(); err != nil {
		apihelper.ValidationErrResp(w, err)
		return
	}

	state, err := param.decodeState()
	if err != nil {
		apihelper.InternalServerErrResp(w, err)
		return
	}
	redirectTo := state["redirect_to"]

	authHandler := oauth2Provider(param.Provider, environment)
	userData, err := authHandler.ExchangeCode(param.Code)
	if err != nil {
		apihelper.ValidationErrResp(w, map[string]string{
			"code": "invalid",
		})
		return
	}

	user := &core.User{
		Name:  userData.Name,
		Email: userData.Email,
	}
	user.SetOauth2Identifier(userData.Provider, &userData.Identifier)

	authorizationCode, err := core.BuildOauthAuthorizationCode(oauthAuthorizationCodeLength)
	if err != nil {
		apihelper.InternalServerErrResp(w, err)
		return
	}

	if err := registerUser(r.Context(), user, authorizationCode); err != nil {
		apihelper.InternalServerErrResp(w, err)
		return
	}

	parsedRedirectTo, err := url.Parse(redirectTo)
	if err != nil {
		apihelper.BadRequestErrResp(w, "bad_request", map[string]string{
			"state": "invalid",
		})
		return
	}
	query := parsedRedirectTo.Query()
	query.Set("code", authorizationCode)
	parsedRedirectTo.RawQuery = query.Encode()

	apihelper.RedirectResp(w, parsedRedirectTo.String())
}

// @Id api.v1.auth.exchangeToken
// @Summary Exchange authorization code for authentication token
// @Tags Auth
// @Accept json
// @Param param body postTokenParam true "Request body"
// @Router /v1/auth/token [post]
// @Produce json
// @Success 200 {object} oauth2.Token
// @Failure 400 {object} apihelper.HTTPError
// @Failure 401 {object} apihelper.HTTPError
// @Failure 500 {object} apihelper.InternalServerError
func HandleExchangeOauthToken(oauthTokenSecretKey []byte) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			apihelper.BadRequestErrResp(w, "invalid_request", map[string]string{
				"request_body": "malformed format",
			})
			return
		}

		param := postTokenParam{
			oauthTokenSecretKey: oauthTokenSecretKey,
		}
		err := json.NewDecoder(r.Body).Decode(&param)
		if err != nil {
			apihelper.BadRequestErrResp(w, "invalid_request", map[string]string{
				"request_body": "malformed format",
			})
			return
		}
		if err := param.Validate(r.Context()); err != nil {
			apihelper.ValidationErrResp(w, err)
			return
		}

		token, err := core.BuildOauthToken(oauthTokenSecretKey, oauthTokenLength)
		if err != nil {
			apihelper.InternalServerErrResp(w, err)
			return
		}

		token.RequesterIP = r.Header.Get("X-Real-Ip")
		if token.RequesterIP == "" {
			parts := strings.Split(r.Header.Get("X-Forwarded-For"), " ")
			if len(parts) > 0 {
				token.RequesterIP = parts[len(parts)-1]
			}
		}
		if token.RequesterIP == "" {
			token.RequesterIP = "127.0.0.1"
		}

		token.RequesterUserAgent = r.Header.Get("User-Agent")

		switch param.GrantType {
		case "refresh_token":
			// if grant type is refresh_token, also delete old token
			if err := deleteOauthToken(r.Context(), param.retrievedOauthToken.ID); err != nil {
				apihelper.InternalServerErrResp(w, err)
				return
			}
			token.UserID = param.retrievedOauthToken.UserID
		case "authorization_code":
			token.UserID = param.retrievedCode.UserID
		}

		if err := saveOauthToken(r.Context(), token); err != nil {
			apihelper.InternalServerErrResp(w, err)
			return
		}
		apihelper.JSON(w, http.StatusOK, token.Token())
	}
	return http.HandlerFunc(handler)
}
