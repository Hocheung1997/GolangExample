package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type RequestLog struct {
	Method     string `json:"method"`
	Path       string `json:"path"`
	RemoteAddr string `json:"remote_addr"`
	Duration   string `json:"duration"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(startTime).String()
		logEntry := RequestLog{
			Method:     r.Method,
			Path:       r.URL.Path,
			RemoteAddr: r.RemoteAddr,
			Duration:   duration,
		}
		logJSON, err := json.Marshal(logEntry)
		if err != nil {
			log.Println("Error marshaling log entry:", err)
			return
		}
		log.Println(string(logJSON))
	})
}

func apiHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "ok")
}

func setupHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthCheckHandler)
	mux.HandleFunc("/api", apiHandler)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}
	mux := http.NewServeMux()
	setupHandlers(mux)
	log.Fatal(http.ListenAndServe(listenAddr, loggingMiddleware(mux)))
}
