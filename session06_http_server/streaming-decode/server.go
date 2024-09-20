package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type logLine struct {
	UserIP string `json:"user_ip"`
	Event  string `json:"event"`
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var e *json.UnmarshalTypeError
	for {
		var l logLine
		err := dec.Decode(&l)
		if err == io.EOF {
			break
		}
		if errors.As(err, &e) {
			log.Println(err)
			continue
		}
		if err != nil {
			if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
				http.Error(w, fmt.Sprintf("key %s is not matched: %s", unmarshalErr.Field, unmarshalErr.Error()), http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		fmt.Println(l.UserIP, l.Event)
	}
	fmt.Fprintf(w, "OK")

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/decode", decodeHandler)
	http.ListenAndServe(":8080", mux)
}
