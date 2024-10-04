package database

import (
	"testing"
)

// TestCase represents a single input case for the target function.
// It contains a message and some data
type TestCase[D any] struct {
	Message string
	Data    D
}

// TestSuite contains a target function and some cases
// with which the target function will be executed
type TestSuite[D any] struct {
	Target func(*testing.T, string, D)
	Cases  []TestCase[D]
}

// Run executes all the cases
func (ts TestSuite[D]) Run(t *testing.T) {
	for _, tc := range ts.Cases {
		ts.Target(t, tc.Message, tc.Data)
	}
}
