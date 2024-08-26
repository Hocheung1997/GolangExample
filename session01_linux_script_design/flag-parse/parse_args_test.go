package main

import (
	"bytes"
	"errors"
	"testing"
)

type testConfig struct {
	args        []string
	expectedErr error
	numTimes    int
}

func TestParseArgs(t *testing.T) {
	testCases := []testConfig{
		{
			args:        []string{"-h"},
			expectedErr: nil,
			numTimes:    0,
		},
		{
			args:        []string{"-n", "10"},
			expectedErr: nil,
			numTimes:    10,
		},
		{
			args:        []string{"-n", "abc"},
			expectedErr: errors.New("invalid value \"abc\" for flag -n: parse error"),
			numTimes:    0,
		},
		{
			args:        []string{"-n", "1", "foo"},
			expectedErr: errors.New("positional arguments specified"),
			numTimes:    1,
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range testCases {
		c, err := parseAgrs(byteBuf, tc.args)
		if tc.expectedErr == nil && err != nil {
			t.Errorf("expected nil error, got: %v\n", err)
		}
		if tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
			t.Errorf("expected error to be: %v, got: %v\n", tc.expectedErr, err)
		}
		if c.numTimes != tc.numTimes {
			t.Errorf("expected numTimes to be: %v, got: %v", tc.numTimes, c.numTimes)
		}
		byteBuf.Reset()
	}
}
