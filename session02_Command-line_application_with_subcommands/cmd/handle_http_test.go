package cmd

import (
	"bytes"
	"errors"
	"strconv"
	"testing"
)

func TestHandleHttp(t *testing.T) {
	usageMessage := `
http: A HTTP client.

http: <options> server

Options: 
  -verb string
    	HTTP method (default "GET")
`
	testConfigs := []struct {
		args       []string
		output     string
		err        error
		statusCode int
	}{
		{
			args: []string{},
			err:  ErrNoServerSpecified,
		},
		{
			args:   []string{"-h"},
			err:    errors.New("flag: help requested"),
			output: usageMessage,
		},
		{
			args:       []string{"https://www.example.com/"},
			err:        nil,
			statusCode: 200,
		},
	}
	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleHttp(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error %v, got %v", tc.err, err)
		}

		if tc.statusCode != 0 {
			returnedStatusCode := byteBuf.String()
			if strconv.Itoa(tc.statusCode) != returnedStatusCode {
				t.Errorf("Expected output to be: %d, Got: %s", tc.statusCode, returnedStatusCode)
			}
		}
		byteBuf.Reset()
	}
}
