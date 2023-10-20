package middleware

import (
	"context"
	"fmt"
	"net/http"
)

var (
	apiPath  = "api/v1"
	authUrls = map[string]string{
		apiPath + "/user":    http.MethodGet,
		apiPath + "/logout":  http.MethodDelete,
		apiPath + "/places":  http.MethodPost,
		apiPath + "/auth":    http.MethodGet,
		apiPath + "/users":   http.MethodPatch,
		apiPath + "/reviews": http.MethodDelete,
		apiPath + "/review":  http.MethodPost,
		apiPath + "/trip":    http.MethodPost,
	}
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := sm.Check(r)
		_, canbeWithouthSess := noSessUrls[r.URL.Path]
		if err != nil && !canbeWithouthSess {
			fmt.Println("no auth")
			http.Redirect(w, r, "/", 302)
			return
		}
		ctx := context.WithValue(r.Context(), models.SessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
