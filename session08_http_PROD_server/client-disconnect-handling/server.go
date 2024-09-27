package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	log.Println("Ping: Got a request")
	time.Sleep(10 * time.Second)
	fmt.Fprint(w, "pong")
}

func doSomeWork(data []byte) {
	time.Sleep(15 * time.Second)
}

func createHTTPGetRequestWithTrace(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	trace := &httptrace.ClientTrace{
		DNSStart:    func(di httptrace.DNSStartInfo) { log.Printf("DNS Start Info: %+v\n", di) },
		DNSDone:     func(di httptrace.DNSDoneInfo) { log.Printf("DNS Done Info: %+v\n", di) },
		GotConn:     func(connInfo httptrace.GotConnInfo) { log.Printf("Got Conn: %+v\n", connInfo) },
		PutIdleConn: func(err error) { log.Printf("Put Idle Conn Error: %+v\n", err) },
	}
	ctxTrace := httptrace.WithClientTrace(req.Context(), trace)
	req = req.WithContext(ctxTrace)
	return req, err
}

func handleUserAPI(w http.ResponseWriter, r *http.Request) {
	done := make(chan bool)
	log.Println("I started prcoessing the request")

	req, err := createHTTPGetRequestWithTrace(r.Context(), "http://localhost:8080/ping")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	log.Println("Outgoing HTTP request")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	log.Println("Processing the response i got")

	go func() {
		doSomeWork(data)
		done <- true
	}()

	select {
	case <-done:
		log.Println("doSomeWork done:Continuing request processing")
	case <-r.Context().Done():
		log.Printf("Aborting request processing: %v\n", r.Context().Err())
		return
	}

	fmt.Fprint(w, string(data))
	log.Println("I finished processing the request")
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	timeoutDuration := 30 * time.Second

	userHandler := http.HandlerFunc(handleUserAPI)
	hTimeout := http.TimeoutHandler(
		userHandler,
		timeoutDuration,
		"I ran out of time")
	mux := http.NewServeMux()
	mux.Handle("/api/users/", hTimeout)
	mux.HandleFunc("/ping", handlePing)

	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
