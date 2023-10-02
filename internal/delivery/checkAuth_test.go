package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckAuthNegative(t *testing.T) {
	positiveTests := map[string]struct {
		cookie       string
		expectedBody string
		expectedCode int
	}{
		"user has wrong cookie": {
			cookie:       "abdfsdc.ssfsfddad.assdfsdfdj",
			expectedCode: http.StatusUnauthorized,
		},
		"user cookie is empty": {
			cookie:       "",
			expectedCode: http.StatusUnauthorized,
		},
	}
	for name, test := range positiveTests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			reader := bytes.NewReader([]byte{})
			req, err := http.NewRequest("GET", "/auth", reader)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := NewAuthHandler("fdsjhfsidfsd")

			expire := time.Now().Add(7 * 24 * time.Hour)
			cookie := http.Cookie{
				Name:    "session_id",
				Value:   test.cookie,
				Expires: expire,
			}
			req.AddCookie(&cookie)

			checkAuthHandler := http.HandlerFunc(handler.CheckAuth)
			checkAuthHandler.ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, test.expectedCode, "Http response has wrong status")
		})
	}
}

func TestCheckAuthPositives(t *testing.T) {
	positiveTests := map[string]struct {
		expectedBody string
		expectedCode int
	}{
		"user has cookie": {
			expectedBody: `{"error":""}`,
			expectedCode: http.StatusOK,
		},
	}
	for name, test := range positiveTests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			login := "hello"
			reader := bytes.NewReader([]byte{})
			req, err := http.NewRequest("GET", "/auth", reader)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := NewAuthHandler("fdsjhfsidfsd")
			cookie, err := handler.createSessionCookie(login)
			assert.NoError(t, err)
			req.AddCookie(cookie)

			checkAuthHandler := http.HandlerFunc(handler.CheckAuth)
			checkAuthHandler.ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, test.expectedCode, "Http response has wrong status")
			assert.Equal(t, rr.Body.String(), test.expectedBody, "Json response is wrong.")
		})
	}
}
