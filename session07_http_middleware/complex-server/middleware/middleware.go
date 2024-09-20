package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Hocheung1997/complex-server/config"
)

func loggingMiddleware(h http.Handler, c config.AppConfig) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			h.ServeHTTP(w, r)
			requestDuration := time.Since(t1).Seconds()
			c.Logger.Printf(
				"protocol=%s path=%s method=%s duration=%f", r.Proto, r.URL.Path, r.Method, requestDuration)
		})
}

func panicMiddleware(h http.Handler, c config.AppConfig) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rVaule := recover(); rVaule != nil {
					c.Logger.Println("panic detected", rVaule)
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "Unexpected server error occurred")
				}
			}()
			h.ServeHTTP(w, r)
		})
}
