package database

import (
	"testing"
)

// testCase represents a single input case for the target function.
// It contains a message and some data
type testCase[D any] struct {
	Message string
	Data    D
}

// testSuite contains a target function and some cases
// with which the target function will be executed
type testSuite[D any] struct {
	Target func(*testing.T, string, D)
	Cases  []testCase[D]
}

// Run executes all the cases
func (ts testSuite[D]) Run(t *testing.T) {
	for _, tc := range ts.Cases {
		ts.Target(t, tc.Message, tc.Data)
	}
}
