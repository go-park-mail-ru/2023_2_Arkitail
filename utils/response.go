package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/microcosm-cc/bluemonday"
)

type UserClaim struct {
	Id uint
	jwt.RegisteredClaims
}

var ErrTokenInvalid = errors.New("token is invalid")

func CreateErrorResponse(errorMsg string) []byte {
	response := ErrorResponse{Error: errorMsg}
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return nil
	}
	return responseJson
}

func WriteResponse(w http.ResponseWriter, status int, body []byte) {
	if body == nil {
		w.WriteHeader(status)
		return
	}
	sanitizer := bluemonday.UGCPolicy()
	body = sanitizer.SanitizeBytes(body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

type ErrorResponse struct {
	Error string `json:"error"`
}
