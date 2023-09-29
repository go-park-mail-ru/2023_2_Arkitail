package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreatePlace(t *testing.T) {
	// Создаем фейковый сервер
	r := mux.NewRouter()
	r.HandleFunc("/places", createPlace).Methods("POST")

	// Создаем JSON-данные для нового места
	newPlace := Place{
		Name:        "Новое место",
		Description: "Описание нового места",
		Rating:      4.0,
		Cost:        "дешево",
		ImageURL:    "https://example.com/image.jpg",
	}

	// Кодируем JSON-данные
	jsonData, err := json.Marshal(newPlace)
	if err != nil {
		t.Fatalf("Ошибка кодирования JSON: %v", err)
	}

	// Создаем POST-запрос с данными
	req, err := http.NewRequest("POST", "/places", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}

	// Создаем запись запроса (response recorder) для записи ответа
	rr := httptest.NewRecorder()

	// Запускаем запрос
	r.ServeHTTP(rr, req)

	// Проверяем статус код ответа
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusCreated, status)
	}

	// Проверяем, что место было добавлено
	if len(places) != 1 {
		t.Errorf("Ожидалось, что место будет добавлено, но не добавлено")
	}
}
