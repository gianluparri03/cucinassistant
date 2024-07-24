package database

import (
	"testing"
)

func TestUserSignup(t *testing.T) {
	type testCase struct {
		user        *User
		expected    string
		usersNumber int
		onFail      string
	}

	testCases := []testCase{
		{
			user:        &User{Username: "u"},
			expected:    "Nome utente non valido: lunghezza minima 5 caratteri",
			usersNumber: 0,
			onFail:      "User.SignUp: username length not checked",
		},
		{
			user:        &User{Username: "username", Email: "email", Password: "p"},
			expected:    "Email non valida",
			usersNumber: 0,
			onFail:      "User.SignUp: email not checked",
		},
		{
			user:        &User{Username: "username", Email: "email@email.com", Password: "p"},
			expected:    "Password non valida: lunghezza minima 8 caratteri",
			usersNumber: 0,
			onFail:      "User.SignUp: password length not checked",
		},
		{
			user:        &User{Username: "username", Email: "email@email.com", Password: "password"},
			expected:    "",
			usersNumber: 1,
			onFail:      "User.SignUp: valid user could not sign up",
		},
	}

	for _, tc := range testCases {
		err := tc.user.SignUp()
		if tc.expected == "" && err != nil {
			t.Errorf(tc.onFail + ": expected <nil>, got <" + err.Error() + ">")
		} else if tc.expected != "" && err == nil {
			t.Errorf(tc.onFail + ": expected <" + tc.expected + ">, got <nil>")
		} else if tc.expected != "" && tc.expected != err.Error() {
			t.Errorf(tc.onFail + ": expected <" + tc.expected + ">, got <" + err.Error() + ">")
		} else if GetUsersNumber() != tc.usersNumber {
			t.Errorf(tc.onFail + ": unexpected number of users")
		}
	}
}
