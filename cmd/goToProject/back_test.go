package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestAddPlace(t *testing.T) {
    // Создайте тестовый сервер и клиент для отправки запросов
    ts := httptest.NewServer(http.HandlerFunc(yourHandlerFunction))
    defer ts.Close()

    // Создайте JSON-данные для нового места
    newPlace := Place{
        Name:        "Тестовое место",
        Description: "Описание тестового места",
        Rating:      4.0,
        PriceGroup:  "дёшево",
    }

    // Преобразуйте данные в JSON
    jsonData, err := json.Marshal(newPlace)
    if err != nil {
        t.Fatal(err)
    }

    // Отправьте POST-запрос на сервер
    resp, err := http.Post(ts.URL+"/places", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    // Проверьте, что статус ответа равен 201 (Created)
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("Ожидается статус 201, получено: %d", resp.StatusCode)
    }
}