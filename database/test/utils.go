package test

import (
	"fmt"
	"testing"

	"cucinassistant/database"
)

// TestCase represents a single input case for a function
type TestCase[R comparable] struct {
	// Description will be print if the case fails
	Description string

	// User is the user who will run the test
	User *database.User

	// Expected is the expected output of the function
	Expected R

	// Data contains some optional data used by the PostCheck
	Data map[string]any
}

// TestSuite groups all the TestCases for a function
type TestSuite[R comparable] struct {
	// Target is the function on which the tests are made
	Target func(tc *TestCase[R]) R

	// Cases is the list of all the cases to be run
	Cases []TestCase[R]

	// PostCheck is a function which is run after the test,
	// only if it is successfull
	PostCheck func(t *testing.T, tc *TestCase[R])
}

// Run executes all the cases
func (ts TestSuite[R]) Run(t *testing.T) {
	for _, tc := range ts.Cases {
		if got := ts.Target(&tc); got != tc.Expected {
			t.Errorf("%s: error: expected <%v>, got <%v>", tc.Description, tc.Expected, got)
		} else {
			ts.PostCheck(t, &tc)
		}
	}
}

var userN int = 0

// GetTestingUser returns an user to be used for testing purposes
func GetTestingUser(t *testing.T) (user database.User) {
	userN++

	// Generate an user with incremental data
	password := fmt.Sprintf("password%d", userN)
	user = database.User{
		Username: fmt.Sprintf("username%d", userN),
		Email:    fmt.Sprintf("email%d@email.com", userN),
		Password: password,
	}

	// And tries to sign it up
	if err := user.SignUp(); err != nil {
		t.Fatalf("Cannot create testing user: %s", err.Error())
	} else {
		// Restore the plain text password
		user.Password = password
	}

	return
}
