package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"project/internal/middleware"
	"project/internal/router"

	"project/users/handler"
	"project/users/repo"
	"project/users/usecase"

	phandler "project/places/handler"
	prepo "project/places/repo"
	pusecase "project/places/usecase"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"
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

	logger := new(logrus.Logger)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

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

	r.Use(middleware.AccessLog(logger))
	r.Use(middleware.Panic)

	fmt.Println("Server is running on :8080")
	err = http.ListenAndServe(":8080", h)
	if err != nil {
		fmt.Println(err)
	}
}
