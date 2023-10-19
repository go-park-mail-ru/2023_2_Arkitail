package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"project/internal/router"

	"project/users/handler"
	"project/users/repo"
	"project/users/usecase"

	phandler "project/places/handler"
	prepo "project/places/repo"
	pusecase "project/places/usecase"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
)

type DBconfig struct {
	user     string
	dbname   string
	password string
	host     string
	port     int
	sslmode  string
}

type ConnectionConfig struct {
	maxConnectionCount int
}

func getPosgres() *sql.DB {
	dbConfig := DBconfig{
		"GoTo", "GoTo", "qwerty", "127.0.0.1", 5432, "disable",
	}
	connectionConfig := ConnectionConfig{10}

	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=%s",
		dbConfig.user, dbConfig.dbname, dbConfig.password,
		dbConfig.host, dbConfig.port, dbConfig.sslmode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("cant parce database config")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(connectionConfig.maxConnectionCount)
	return db
}

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

	db := getPosgres()

	userRepo := repo.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, authConfig)
	userHandler := handler.NewUserHandler(userUsecase)

	placeRepo := prepo.NewPlaceRepository(db)
	placeUseCase := pusecase.NewPlaceUseCase(placeRepo)
	placeHandler := phandler.NewPlaceHandler(placeUseCase)

	r := mux.NewRouter()

	apiPath := "/api/v1"
	r.HandleFunc(apiPath+"/auth", userHandler.CheckAuth).Methods("GET")
	r.HandleFunc(apiPath+"/login", userHandler.Login).Methods("POST")
	r.HandleFunc(apiPath+"/signup", userHandler.Signup).Methods("POST")
	r.HandleFunc(apiPath+"/logout", userHandler.Logout).Methods("DELETE")
	r.HandleFunc(apiPath+"/user", userHandler.GetUserInfo).Methods("GET")

	h := router.AddCors(r, []string{"http://localhost:8080/"})

	r.HandleFunc(apiPath+"/places", placeHandler.CreatePlace).Methods("POST")
	r.HandleFunc(apiPath+"/places", placeHandler.GetPlaces).Methods("GET")

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", h)
	if err != nil {
		fmt.Println(err)
	}
}
