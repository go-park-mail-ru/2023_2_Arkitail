package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func AccessLog(logger *logrus.Entry) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// start := time.Now()
			next.ServeHTTP(w, r)
			logger.WithFields(logrus.Fields{
				// "method": r.Method,
				// "remote_addr": r.RemoteAddr,
				// "work_time":   time.Since(start),
				// "status_code": w.Header().Get("Status"),
			}).Info(r.URL.Path)
		})
	}
}
