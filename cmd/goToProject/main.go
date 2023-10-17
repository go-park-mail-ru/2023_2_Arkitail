package main

import (
	"flag"
	"fmt"
	"net/http"

	"project/internal/router"
	hand "project/places/handler"

	"project/users/handler"
	"project/users/repo"
	"project/users/usecase"

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
	authConfig := usecase.AuthConfig{
		Secret: []byte(secret),
	}

	userRepo := repo.NewUserRepository()

	userUsecase := usecase.NewUserUsecase(userRepo, authConfig)

	userHandler := handler.NewUserHandler(userUsecase)


	r := mux.NewRouter()

	apiPath := "/api/v1"
	// Регистрируйте маршруты и обработчики, используя userHandler
	r.HandleFunc(apiPath+"/auth", userHandler.CheckAuth).Methods("GET")
	r.HandleFunc(apiPath+"/login", userHandler.Login).Methods("POST")
	r.HandleFunc(apiPath+"/signup", userHandler.Signup).Methods("POST")
	r.HandleFunc(apiPath+"/logout", userHandler.Logout).Methods("DELETE")
	r.HandleFunc(apiPath+"/user", userHandler.GetUserInfo).Methods("GET")

	handler := router.AddCors(r, []string{"http://localhost:8080/"})

	r.HandleFunc(apiPath+"/places", hand.CreatePlace).Methods("POST")
	r.HandleFunc(apiPath+"/places", hand.GetPlaces).Methods("GET")

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8088", handler)
	if err != nil {
		fmt.Println(err)
	}
}
