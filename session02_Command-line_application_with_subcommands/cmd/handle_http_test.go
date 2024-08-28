package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestHandleHttp(t *testing.T) {
	usageMessage := `Usage: mync [http|grpc] -h
	
                http: A HTTP client.

                http: <options> serverOptions:  -verb string
        HTTP method (default "GET")

                grpc: A gRPC client.

                grpc: <options> server

Options:
  -body string
        body of request
  -method string
        Method to call
`
	testConfigs := []struct {
		args   []string
		output string
		err    error
	}{
		{
			args: []string{},
			err:  ErrNoServiceSpecified,
		},
		{
			args:   []string{"-h"},
			err:    errors.New("flag: help requested"),
			output: usageMessage,
		},
		{
			args:   []string{"http://localhost"},
			err:    nil,
			output: "Executing http command\n",
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testConfigs {
		err := HandleHttp(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("expected error %v, got %v", tc.err, tc.err.Error())
		}
		if len(tc.output) != 0 {
			gotOutput := byteBuf.String()
			if tc.output != gotOutput {
				t.Errorf("expected output to be: %#v, Got: %#v", tc.output, gotOutput)
			}
		}
	}
	byteBuf.Reset()
}
