package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	positiveTests := map[string]struct {
		setCookie bool
	}{
		"user has cookie": {
			setCookie: true,
		},
		"user doesnt have cookie": {
			setCookie: false,
		},
	}
	for name, test := range positiveTests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			reader := bytes.NewReader([]byte{})
			req, err := http.NewRequest("Delete", "/auth", reader)
			assert.NoError(t, err)
			if test.setCookie == true {
				expire := time.Now().Add(-7 * 24 * time.Hour)
				cookie := http.Cookie{
					Name:    "session_id",
					Value:   "value",
					Expires: expire,
				}
				req.AddCookie(&cookie)
			}
			rr := httptest.NewRecorder()
			handler := NewAuthHandler("fdsjhfsidfsd")
			logoutHandler := http.HandlerFunc(handler.Logout)
			logoutHandler.ServeHTTP(rr, req)
			if len(rr.Result().Cookies()) > 0 {
				assert.Greater(t, time.Now(), rr.Result().Cookies()[0].Expires, "Cookie was not deleted")
			} else {
				assert.Equal(t, 0, len(rr.Result().Cookies()), "Cookie was not deleted")
			}
		})
	}
}
