package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {
	testCases := []struct {
		c           config
		input       string
		output      string
		expectedErr error
	}{
		{
			c:      config{printUsage: true},
			output: usageString,
		},
		{
			c:           config{numTimes: 5},
			input:       "",
			output:      "please input your name, press enter to process:\n",
			expectedErr: errors.New("blank value, this feild should input your name"),
		},
		{
			c:      config{numTimes: 5},
			input:  "Hocheung",
			output: "please input your name, press enter to process:\n" + strings.Repeat("Nice to meet you, Hocheung\n", 5),
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testCases {
		rd := strings.NewReader(tc.input)
		err := runCmd(rd, byteBuf, tc.c)
		if err != nil && tc.expectedErr == nil {
			t.Fatalf("expected nil error, got: %v\n", err)
		}
		if tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
			t.Errorf("expected error: %v, got: %v\n", tc.expectedErr, err)
		}
		gotMsg := byteBuf.String()
		if gotMsg != tc.output {
			t.Errorf("expected standard output message to be: %v, Got: %v\n", tc.output, gotMsg)
		}
		byteBuf.Reset()
	}

}
