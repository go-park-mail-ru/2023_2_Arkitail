package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"project/users/model"
	"project/users/usecase"

	"github.com/gorilla/mux"
)

var (
	AuthNames = map[string]struct{}{
		"CreatePlace": struct{}{},
		"Auth":        struct{}{},
		"User":        struct{}{},
	}
	apiPath = "api/v1"
)

var errTokenInvalid = errors.New("token is invalid")

// TODO: выделить урлы, для которых не должно быть сессии
func Auth(ucase usecase.UserUseCase) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := AuthNames[mux.CurrentRoute(r).GetName()]; !ok {
				next.ServeHTTP(w, r)
				return
			}
			cookie, err := r.Cookie("session_id")
			if err != nil {
				writeResponse(w, http.StatusUnauthorized, createErrorResponse(errTokenInvalid.Error()))
				return
			}

			token := cookie.Value
			user, err := ucase.ValidateToken(token)
			if err != nil {
				writeResponse(w, http.StatusUnauthorized, createErrorResponse(errTokenInvalid.Error()))
				return
			}

			ctx := context.WithValue(r.Context(), "userClaim", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func createErrorResponse(errorMsg string) []byte {
	response := model.ErrorResponse{Error: errorMsg}
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return nil
	}
	return responseJson
}

func createUserResponse(user *model.User) ([]byte, error) {
	responseJson, err := json.Marshal(user)
	return responseJson, err
}

func writeResponse(w http.ResponseWriter, status int, body []byte) {
	if body == nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}
