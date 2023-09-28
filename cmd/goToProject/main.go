package main

import (
	"flag"
	"fmt"
	"net/http"

	auth "project/internal/delivery"

	"github.com/gorilla/mux"
)

func main() {
	var secret string
	flag.StringVar(&secret, "secret", "", "secret for jwt encoding")
	flag.Parse()
	if secret == "" {
		flag.Usage()
		return
	}

	r := mux.NewRouter()
	authHandler := auth.NewAuthHandler(secret)
	apiPath := "/api/v1"
	r.HandleFunc(apiPath+"/auth", authHandler.CheckAuth).Methods("GET")
	r.HandleFunc(apiPath+"/login", authHandler.Login).Methods("POST")
	r.HandleFunc(apiPath+"/signup", authHandler.Signup).Methods("POST")
	r.HandleFunc(apiPath+"/logout", authHandler.Logout).Methods("Delete")

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
