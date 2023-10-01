package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func AddCors(router *mux.Router, originNames []string) {
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST"})
	handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins(originNames)
	handlers.CORS(credentials, methods, origins)(router)
}
