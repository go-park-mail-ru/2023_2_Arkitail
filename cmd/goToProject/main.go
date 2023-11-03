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

	userHandler "project/users/handler"
	userRepo "project/users/repo"
	userUsecase "project/users/usecase"

	reviewHandler "project/reviews/handler"
	reviewRepo "project/reviews/repo"
	reviewUsecase "project/reviews/usecase"

	"project/utils/api"

	placeHandler "project/places/handler"
	placeRepo "project/places/repo"
	placeUsecase "project/places/usecase"

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
	authConfig := userUsecase.AuthConfig{
		Secret: []byte(secret),
	}

	db, err := getPosgres()
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepo := userRepo.NewUserRepository(db)
	userUsecase := userUsecase.NewUserUsecase(userRepo, authConfig)
	userHandler := userHandler.NewUserHandler(userUsecase)

	placeRepo := placeRepo.NewPlaceRepository(db)
	placeUseCase := placeUsecase.NewPlaceUseCase(placeRepo)
	placeHandler := placeHandler.NewPlaceHandler(placeUseCase)

	reviewRepo := reviewRepo.NewReviewRepository(db)
	reviewUseCase := reviewUsecase.NewReviewUsecase(reviewRepo)
	reviewHandler := reviewHandler.NewReviewHandler(reviewUseCase)

	r := mux.NewRouter()

	apiPath := "/api/v1"
	r.HandleFunc(apiPath+api.Auth, userHandler.CheckAuth).Methods("GET").Name(api.Auth)
	r.HandleFunc(apiPath+api.Login, userHandler.Login).Methods("POST").Name(api.Login)
	r.HandleFunc(apiPath+api.Signup, userHandler.Signup).Methods("POST").Name(api.Signup)
	r.HandleFunc(apiPath+api.Logout, userHandler.Logout).Methods("DELETE").Name(api.Logout)
	r.HandleFunc(apiPath+api.User, userHandler.GetUserInfo).Methods("GET").Name(api.User)
	r.HandleFunc(apiPath+api.User, userHandler.PatchUser).Methods("Patch").Name(api.User)
	r.HandleFunc(apiPath+api.UserById, userHandler.GetCleanUser).Methods("GET").Name(api.UserById)
	r.HandleFunc(apiPath+api.UserAvatar, userHandler.UploadAvatar).Methods("Post").Name(api.UserAvatar)

	r.HandleFunc(apiPath+api.ReviewById, reviewHandler.DeleteReview).Methods("Delete").Name(api.ReviewById)
	r.HandleFunc(apiPath+api.Review, reviewHandler.AddReview).Methods("POST").Name(api.Review)
	r.HandleFunc(apiPath+api.PlaceReviews, reviewHandler.GetPlaceReviews).Methods("GET").Name(api.PlaceReviews)

	h := router.AddCors(r, []string{"http://localhost:8080/"})

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
