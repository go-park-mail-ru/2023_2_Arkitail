package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateAuthBody(login string, password string) (res []byte) {
	authRequest := &AuthUser{login, password}
	res, _ = json.Marshal(authRequest)
	return
}
func TestLoginPositives(t *testing.T) {
	positiveTests := map[string]struct {
		login         string
		password      string
		addedLogin    string
		addedPassword string
		expectedBody  string
		expectedCode  int
	}{
		"correct credentials": {
			login:         "abc",
			password:      "1234",
			addedLogin:    "abc",
			addedPassword: "1234",
			expectedBody:  `{"error":""}`,
			expectedCode:  http.StatusOK,
		}}
	for name, test := range positiveTests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			login := test.login
			password := test.password
			jsonBody := CreateAuthBody(login, password)
			reader := bytes.NewReader(jsonBody)
			req, err := http.NewRequest("POST", "/login", reader)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := NewAuthHandler("fdsjhfsidfsd")
			assert.NoError(t, handler.storage.AddUser(test.addedLogin, test.addedPassword))
			loginHandler := http.HandlerFunc(handler.Login)

			loginHandler.ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, test.expectedCode, "Http response has wrong status")
			assert.Equal(t, rr.Body.String(), test.expectedBody, "Json response is wrong.")
		})
	}
}

func TestLoginNegatives(t *testing.T) {
	positiveTests := map[string]struct {
		login         string
		password      string
		addedLogin    string
		addedPassword string
		expectedBody  string
		expectedCode  int
	}{
		"user doesnt exist": {
			login:         "abc",
			password:      "1234",
			addedLogin:    "trq",
			addedPassword: "fjssd",
			expectedBody:  `{"error":"user with this name doesnt exist"}`,
			expectedCode:  http.StatusUnauthorized,
		},
		"wrong password": {
			login:         "abc",
			password:      "1234",
			addedLogin:    "abc",
			addedPassword: "123456789",
			expectedBody:  `{"error":"wrong password"}`,
			expectedCode:  http.StatusUnauthorized,
		},
	}
	for name, test := range positiveTests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			login := test.login
			password := test.password

			jsonBody := CreateAuthBody(login, password)
			reader := bytes.NewReader(jsonBody)
			req, err := http.NewRequest("POST", "/login", reader)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := NewAuthHandler("fdsjhfsidfsd")
			assert.NoError(t, handler.storage.AddUser(test.addedLogin, test.addedPassword))
			loginHandler := http.HandlerFunc(handler.Login)

			loginHandler.ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, test.expectedCode, "Http response has wrong status")
			assert.Equal(t, rr.Body.String(), test.expectedBody, "Json response is wrong.")
		})
	}
}
