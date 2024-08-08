package database

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"cucinassistant/config"
)

// // Pair is used to pack the results of functions
// // that returns two values
type Pair[A any, B any] struct {
	First  A
	Second B
}

// TestMain sets up the testing environment
func TestMain(m *testing.M) {
	// Loads the configuration
	config.Read(os.Args[len(os.Args)-1])

	// Connects to the database
	// and creates the missing tables
	Connect()
	Bootstrap("schema.sql")

	// Runs the actual tests
	m.Run()
}

// TestCase represents a single input case for a function
type TestCase[R any] struct {
	// Description will be print if the case fails
	Description string

	// User is the user who will run the test
	User *User

	// Expected is the expected output of the function
	Expected R

	// Data contains some optional data used by the PostCheck
	Data map[string]any
}

// TestSuite groups all the TestCases for a function
type TestSuite[R any] struct {
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
		if got := ts.Target(&tc); !reflect.DeepEqual(got, tc.Expected) {
			t.Errorf("%s: error: expected <%v>, got <%v>", tc.Description, tc.Expected, got)
		} else if ts.PostCheck != nil {
			ts.PostCheck(t, &tc)
		}
	}
}

var userN int = 0

// generateTestingUser returns a testing user which has not been registered yet
func generateTestingUser() *User {
	userN++

	return &User{
		Username: fmt.Sprintf("username%d", userN),
		Email:    fmt.Sprintf("email%d@email.com", userN),
		Password: fmt.Sprintf("password%d", userN),
	}
}

// GetTestingUser returns an user to be used for testing purposes
func GetTestingUser(t *testing.T) (user *User, password string) {
	user = generateTestingUser()
	password = user.Password

	// And tries to sign it up
	var err error
	if user, err = SignUp(user.Username, user.Email, user.Password); err != nil {
		t.Fatalf("Cannot create testing user: %s", err.Error())
	}

	return
}
