package database

import (
	"fmt"
	"testing"
)

var testingUsersN int = 0

// generateTestingUser returns a testing user which has not been registered yet
func generateTestingUser() *User {
	testingUsersN++

	return &User{
		Username: fmt.Sprintf("username%d", testingUsersN),
		Email:    fmt.Sprintf("email%d@email.com", testingUsersN),
		Password: fmt.Sprintf("password%d", testingUsersN),
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

var unknownUser *User = &User{}
