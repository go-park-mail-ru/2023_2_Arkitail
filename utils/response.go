package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func CreateErrorResponse(errorMsg string) []byte {
	response := ErrorResponse{Error: errorMsg}
	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		return nil
	}
	return responseJson
}

func WriteResponse(w http.ResponseWriter, status int, body []byte) {
	if body == nil {
		w.WriteHeader(status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

type ErrorResponse struct {
	Error string `json:"error"`
}
