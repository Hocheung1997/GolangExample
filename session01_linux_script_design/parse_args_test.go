package main

import (
	"errors"
	"testing"
)

type testConfig struct {
	args        []string
	expectedErr error
	config
}

func TestParseArgs(t *testing.T) {
	testCases := []testConfig{
		{
			args:        []string{"-h"},
			expectedErr: nil,
			config:      config{numTimes: 0, printUsage: true},
		},
		{
			args:        []string{"10"},
			expectedErr: nil,
			config:      config{numTimes: 10, printUsage: false},
		},
		{
			args:        []string{"abc"},
			expectedErr: errors.New("strconv.Atoi: parsing \"abc\": invalid syntax"),
			config:      config{numTimes: 0, printUsage: false},
		},
		{
			args:        []string{"1", "foo"},
			expectedErr: errors.New("invaild number of arguments"),
			config:      config{numTimes: 0, printUsage: false},
		},
	}

	for _, tc := range testCases {
		c, err := parseAgrs(tc.args)
		if tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
			t.Fatalf("expected error to be: %v, got: %v\n", tc.expectedErr, err)
		}
		if tc.expectedErr == nil && err != nil {
			t.Errorf("expected nil error, got: %v\n", err)
		}
		if c.printUsage != tc.config.printUsage {
			t.Errorf("Expected printUsage to be: %v, got: %v\n", tc.config.printUsage, c.printUsage)
		}
		if c.numTimes != tc.config.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v", tc.config.numTimes, c.numTimes)
		}
	}
}
