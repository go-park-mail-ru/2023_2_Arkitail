package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	auth "project/internal/delivery"
	"project/internal/middleware"
	"project/internal/router"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	var secret string
	flag.StringVar(&secret, "secret", "", "secret for jwt encoding")
	flag.Parse()
	if secret == "" {
		flag.Usage()
		return
	}

	logger := new(logrus.Logger)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	r := mux.NewRouter()
	authHandler := auth.NewAuthHandler(secret)
	apiPath := "/api/v1"

	r.HandleFunc(apiPath+"/auth", authHandler.CheckAuth).Methods("GET")
	r.HandleFunc(apiPath+"/login", authHandler.Login).Methods("POST")
	r.HandleFunc(apiPath+"/signup", authHandler.Signup).Methods("POST")
	r.HandleFunc(apiPath+"/logout", authHandler.Logout).Methods("Delete")
	r.HandleFunc(apiPath+"/user", authHandler.GetUserInfo).Methods("Get")

	r.HandleFunc(apiPath+"/places", auth.CreatePlace).Methods("POST")
	r.HandleFunc(apiPath+"/places", auth.GetPlaces).Methods("GET")

	r.Use(middleware.AccessLog(logger))
	handler := router.AddCors(r, []string{"http://localhost:8080/"})

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println(err)
	}
}
