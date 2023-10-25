package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"project/internal/middleware"
	"project/internal/router"

	"project/users/handler"
	"project/users/repo"
	"project/users/usecase"
	"project/utils/api"

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

var ErrCantParseConfig = errors.New("cant parce database config")
var ErrCantConnectToDB = errors.New("cant connect to db")

func getPosgres() (*sql.DB, error) {
	dbConfig := DBconfig{
		"GoTo", "GoTo", "qwerty", "127.0.0.1", 5432, "disable",
	}
	connectionConfig := ConnectionConfig{10}

	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=%s",
		dbConfig.user, dbConfig.dbname, dbConfig.password,
		dbConfig.host, dbConfig.port, dbConfig.sslmode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, ErrCantConnectToDB
	}

	err = db.Ping()
	if err != nil {
		return nil, ErrCantParseConfig
	}

	db.SetMaxOpenConns(connectionConfig.maxConnectionCount)
	return db, nil
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

	db, err := getPosgres()
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepo := repo.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, authConfig)
	userHandler := handler.NewUserHandler(userUsecase)

	placeRepo := prepo.NewPlaceRepository(db)

	placeUseCase := pusecase.NewPlaceUseCase(placeRepo)
	placeHandler := phandler.NewPlaceHandler(placeUseCase)

	r := mux.NewRouter()

	apiPath := "/api/v1"
	r.HandleFunc(apiPath+api.Auth, userHandler.CheckAuth).Methods("GET").Name(api.Auth)
	r.HandleFunc(apiPath+api.Login, userHandler.Login).Methods("POST").Name(api.Login)
	r.HandleFunc(apiPath+api.Signup, userHandler.Signup).Methods("POST").Name(api.Signup)
	r.HandleFunc(apiPath+api.Logout, userHandler.Logout).Methods("DELETE").Name(api.Logout)
	r.HandleFunc(apiPath+api.User, userHandler.GetUserInfo).Methods("GET").Name(api.User)
	r.HandleFunc(apiPath+api.Users_by_id, userHandler.PatchUser).Methods("Patch").Name(api.Users_by_id)

	h := router.AddCors(r, []string{"http://localhost:8080/"})

	r.HandleFunc(apiPath+api.Places, placeHandler.CreatePlace).Methods("POST").Name(api.Places)
	r.HandleFunc(apiPath+api.Places, placeHandler.GetPlaces).Methods("GET").Name(api.Places)

	r.Use(middleware.Auth(userUsecase))
	r.Use(middleware.AccessLog)
	r.Use(middleware.Panic)

	fmt.Println("Server is running on :8080")
	err = http.ListenAndServe(":8080", h)
	if err != nil {
		fmt.Println(err)
	}
}
