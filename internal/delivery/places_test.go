package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePlace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(CreatePlace))
	defer ts.Close()

	newPlace := Place{
		ID:          "1",
		Name:        "Новое место",
		Description: "Описание нового места",
		Rating:      4.0,
		Cost:        "дешево",
		ImageURL:    "https://example.com/image.jpg",
	}

	jsonData, err := json.Marshal(newPlace)
	if err != nil {
		t.Fatalf("Ошибка кодирования JSON: %v", err)
	}

	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Ошибка отправки POST-запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusCreated, resp.StatusCode)
	}
}

func TestGetPlaces(t *testing.T) {
    // Initialize some sample data in the 'places' map
    places = map[string]Place{
        "1": {
            ID:          "1",
            Name:        "Новое место",
            Description: "Описание нового места",
            Rating:      4.0,
            Cost:        "дешево",
            ImageURL:    "https://example.com/image.jpg",
        },
    }

    ts := httptest.NewServer(http.HandlerFunc(GetPlaces))
    defer ts.Close()

    resp, err := http.Get(ts.URL)
    if err != nil {
        t.Fatalf("Ошибка отправки GET-запроса: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Ожидался статус код %d, но получен %d", http.StatusOK, resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Fatalf("Ошибка чтения тела ответа: %v", err)
    }

    expectedResponse := `[{"id":"1","name":"Новое место","description":"Описание нового места","rating":4,"cost":"дешево","imageUrl":"https://example.com/image.jpg"}]
`

    if string(body) != expectedResponse {
        t.Errorf("Ожидался ответ:\n%s\n\nПолучен ответ:\n%s", expectedResponse, string(body))
    }
}



