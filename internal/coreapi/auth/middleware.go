package auth

import (
	"context"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/internal/coreapi/utils/apihelper"
	"github.com/awanku/awanku/pkg/core"
)

func OauthTokenValidatorMiddleware(oauthSecretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerParts := strings.Split(r.Header.Get("authorization"), " ")
			if len(headerParts) != 2 {
				apihelper.BadRequestErrResp(w, "invalid_request", map[string]string{
					"authorization_header": "malformed format",
				})
				return
			}

			tokenType := headerParts[0]
			if strings.ToLower(tokenType) != "bearer" {
				apihelper.BadRequestErrResp(w, "invalid_request", map[string]string{
					"authentication_type": "not supported",
				})
				return
			}

			accessTokenRaw := headerParts[1]
			accessTokenParts := strings.Split(accessTokenRaw, ":")
			if len(accessTokenParts) != 2 {
				apihelper.BadRequestErrResp(w, "invalid_request", map[string]string{
					"authentication_credentials": "malformed format",
				})
				return
			}

			accessTokenID := accessTokenParts[0]
			decodedAccessToken, err := base64.URLEncoding.DecodeString(accessTokenParts[1])
			if err != nil {
				apihelper.BadRequestErrResp(w, "invalid_request", map[string]string{
					"authentication_credentials": "malformed format",
				})
				return
			}

			parsedTokenID, _ := strconv.ParseInt(accessTokenID, 10, 64)
			if parsedTokenID <= 0 {
				apihelper.BadRequestErrResp(w, "invalid_request", map[string]string{
					"authentication_credentials": "malformed format",
				})
				return
			}

			storedToken, err := getOauthTokenByID(r.Context(), parsedTokenID)
			if err != nil {
				apihelper.InternalServerErrResp(w, err)
				return
			}
			if err == errOauthTokenExpired {
				apihelper.BadRequestErrResp(w, "invalid_request", map[string]string{
					"token": "expired",
				})
				return
			}

			valid, err := core.ValidateHMAC(oauthSecretKey, []byte(decodedAccessToken), storedToken.AccessTokenHash)
			if err != nil {
				apihelper.InternalServerErrResp(w, err)
				return
			}

			if !valid {
				apihelper.UnauthorizedAccessResp(w, "access_denied", map[string]string{
					"access_token": "invalid",
				})
				return
			}

			user, err := getUserByID(r.Context(), storedToken.UserID)
			if err != nil {
				apihelper.InternalServerErrResp(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), appctx.KeyAuthenticatedUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
