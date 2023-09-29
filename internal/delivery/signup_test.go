package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignupPositives(t *testing.T) {
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
			addedLogin:    "fdsos",
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
			loginHandler := http.HandlerFunc(handler.Signup)

			loginHandler.ServeHTTP(rr, req)
			assert.Equal(t, rr.Body.String(), test.expectedBody, "Json response is wrong.")
		})
	}
}

func TestSignupNegatives(t *testing.T) {
	positiveTests := map[string]struct {
		login         string
		password      string
		addedLogin    string
		addedPassword string
		expectedBody  string
		expectedCode  int
	}{
		"Username taken": {
			login:         "abc",
			password:      "1234",
			addedLogin:    "abc",
			addedPassword: "1234",
			expectedBody:  `{"error":"user with this name already exists"}`,
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
			loginHandler := http.HandlerFunc(handler.Signup)

			loginHandler.ServeHTTP(rr, req)
			assert.Equal(t, rr.Body.String(), test.expectedBody, "Json response is wrong.")
		})
	}
}
