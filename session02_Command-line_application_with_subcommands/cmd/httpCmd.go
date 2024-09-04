package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type httpConfig struct {
	url  string
	verb string
}

func fetchRemoteResource(verb string, url string) ([]byte, int, error) {
	if verb == "GET" {
		r, err := http.Get(url)
		if err != nil {
			return nil, 0, err
		}
		defer r.Body.Close()
		responseBody, err := io.ReadAll(r.Body)
		return responseBody, r.StatusCode, err
	}
	return nil, 0, errors.New("inputed verb not be supported yet")
}

func HandleHttp(w io.Writer, args []string) error {
	var v string
	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&v, "verb", "GET", "HTTP method")

	fs.Usage = func() {
		var usageString = `
http: A HTTP client.

http: <options> server`
		fmt.Fprint(w, usageString)

		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	c := httpConfig{verb: v}
	c.url = fs.Arg(0)
	_, returnCode, err := fetchRemoteResource(c.verb, c.url)
	if err != nil {
		return err
	}
	if returnCode == 200 {
		fmt.Fprint(w, returnCode)
	} else {
		return errors.New("http error code")
	}

	return nil
}
