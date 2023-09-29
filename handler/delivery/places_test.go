package delivery

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreatePlace(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(CreatePlace))
    defer ts.Close()

    newPlace := Place{
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
}