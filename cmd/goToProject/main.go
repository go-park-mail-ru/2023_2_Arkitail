package main

import (
	"fmt"
	"net/http"

	"project/handler/delivery"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/places", delivery.CreatePlace).Methods("POST")
	r.HandleFunc("/api/v1/places", delivery.GetPlaces).Methods("GET")

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}