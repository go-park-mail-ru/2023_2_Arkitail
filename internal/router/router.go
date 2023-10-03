package router

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func AddCors(router *mux.Router, originNames []string) http.Handler {
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPost})
	handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins(originNames)
	handler := handlers.CORS(credentials, methods, origins)(router)
	return handler
}
