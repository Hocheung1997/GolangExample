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

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client

}

func createHTTPGetRequestWithTrace(ctx context.Context, url string) (*http.Request, error) {
	// 1. make a request with URL
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	// define tracy hooks and generate trace context
	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
	}
	// generate a trace context (sub-conext of parent-context)
	ctxTrace := httptrace.WithClientTrace(req.Context(), trace)
	// change request's context
	req = req.WithContext(ctxTrace)
	return req, err
}

func main() {
	d := 5 * time.Second
	ctx := context.Background()
	client := createHTTPClientWithTimeout(d)

	req, err := createHTTPGetRequestWithTrace(ctx, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for {
		resp, err := client.Do(req)
		if err != nil {
			panic("disconnect")
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		time.Sleep(1 * time.Second)
		fmt.Println("---------")
	}
}
