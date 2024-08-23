package main

import (
	"errors"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	testCases := []struct {
		c           config
		expectedErr error
	}{
		{
			c:           config{},
			expectedErr: errors.New("invaild number of output times"),
		},
		{
			c:           config{numTimes: -1},
			expectedErr: errors.New("invaild number of output times"),
		},
		{
			c:           config{numTimes: 10},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		err := vaildateArgs(tc.c)
		if tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
			t.Errorf("expected error to be: %v, got: %v", tc.expectedErr, err)
		}
		if tc.expectedErr == nil && err != nil {
			t.Errorf("expected error to be: nil, got: %v", err)
		}
	}
}
