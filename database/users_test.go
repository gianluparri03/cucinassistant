package database

import (
	"testing"
)

func TestUserSignup(t *testing.T) {
	type testCase struct {
		user        *User
		expected    error
		usersNumber int
		onFail      string
	}

	testCases := []testCase{
		{
			user:        &User{Username: "u"},
			expected:    ERR_USER_NAME_TOO_SHORT,
			usersNumber: 0,
			onFail:      "username length not checked",
		},
		{
			user:        &User{Username: "username", Email: "email", Password: "p"},
			expected:    ERR_USER_MAIL_INVALID,
			usersNumber: 0,
			onFail:      "email not checked",
		},
		{
			user:        &User{Username: "username", Email: "email@email.com", Password: "p"},
			expected:    ERR_USER_PASS_TOO_SHORT,
			usersNumber: 0,
			onFail:      "password length not checked",
		},
		{
			user:        &User{Username: "username", Email: "email@email.com", Password: "password"},
			expected:    nil,
			usersNumber: 1,
			onFail:      "could not sign up",
		},
	}

	for _, tc := range testCases {
		err := tc.user.SignUp()
		if tc.expected != err {
			t.Errorf("%s: error: expected <%v>, got <%v>", tc.onFail, tc.expected, err)
		} else if GetUsersNumber() != tc.usersNumber {
			t.Errorf("%s: wrong number of users: expected %d, got %d", tc.onFail, tc.usersNumber, GetUsersNumber())
		}
	}

	// TODO delete testing user
}

func TestUserSignIn(t *testing.T) {
	u := User{Username: "username2", Email: "email2@email.com", Password: "password"}
	if err := u.SignUp(); err != nil {
		t.Errorf("Cannot create testing user.")
		return
	}

	type testCase struct {
		user     *User
		expected error
		onFail   string
		uid      int
	}

	testCases := []testCase{
		{
			user:     &User{Username: "user", Password: ""},
			expected: ERR_USER_WRONG_CREDENTIALS,
			onFail:   "signed in unknown user",
			uid:      0,
		},
		{
			user:     &User{Username: "username2", Password: "password2"},
			expected: ERR_USER_WRONG_CREDENTIALS,
			onFail:   "signed in with wrong password",
			uid:      0,
		},
		{
			user:     &User{Username: "username2", Password: "password"},
			expected: nil,
			onFail:   "could not sign in",
			uid:      u.UID,
		},
	}

	for _, tc := range testCases {
		err := tc.user.SignIn()
		if tc.expected != err {
			t.Errorf("%s: error: expected <%v>, got <%v>", tc.onFail, tc.expected, err)
		} else if tc.uid != tc.user.UID {
			t.Errorf("%s: wrong uid: expected %d, got %d", tc.onFail, tc.uid, tc.user.UID)
		}
	}

	// TODO delete testing user
}
