package apihelper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		panic(fmt.Errorf("failed to marshal json: %v", err))
	}
}

type HTTPError struct {
	Type   string            `json:"type"`
	Errors map[string]string `json:"errors"`
}

func ValidationErrResp(w http.ResponseWriter, payload interface{}) {
	// TODO: also handle internal error in validation
	BadRequestErrResp(w, "validation_error", payload)
}

func BadRequestErrResp(w http.ResponseWriter, errType string, payload interface{}) {
	JSON(w, http.StatusBadRequest, map[string]interface{}{
		"type":   errType,
		"errors": payload,
	})
}

func UnauthorizedAccessResp(w http.ResponseWriter, errType string, payload interface{}) {
	JSON(w, http.StatusUnauthorized, map[string]interface{}{
		"type":   errType,
		"errors": payload,
	})
}

type InternalServerError struct {
	Error string `json:"error"`
}

func InternalServerErrResp(w http.ResponseWriter, err error) {
	log.Println("internal server err:", err)
	JSON(w, http.StatusInternalServerError, InternalServerError{Error: "something's wrong on our side :("})
}

func RedirectResp(w http.ResponseWriter, to string) {
	w.Header().Add("Location", to)
	w.WriteHeader(http.StatusFound)
}
