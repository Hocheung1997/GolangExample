package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

type cmdArgs struct {
	url          string
	numRequests  int
	maxIdleConns int
}

func createHTTPClientWithTimeout(maxIdleConn int, d time.Duration) *http.Client {

	transport := http.Transport{
		MaxIdleConns:    maxIdleConn,
		IdleConnTimeout: 30 * time.Second,
	}
	client := http.Client{Timeout: d, Transport: &transport}
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

func (c *cmdArgs) createFlagSet(w io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.StringVar(&c.url, "url", "", "input URL")
	fs.IntVar(&c.numRequests, "numRequests", 1, "input the number of requests you need")
	fs.IntVar(&c.maxIdleConns, "maxIdleConns", 1, "input the max number of idle connection")
	fs.Usage = func() {
		var usageString = `
http: A HTTP client.
`
		fmt.Fprint(w, usageString)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	return fs
}

func main() {
	var c cmdArgs
	fs := c.createFlagSet(os.Stdout)
	err := fs.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("error in parsing argumment")
		os.Exit(1)
	}

	d := 5 * time.Second
	ctx := context.Background()
	client := createHTTPClientWithTimeout(c.maxIdleConns, d)

	req, err := createHTTPGetRequestWithTrace(ctx, c.url)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < c.numRequests; i++ {
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
