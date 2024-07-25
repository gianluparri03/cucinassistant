package test

import (
	"testing"

	"cucinassistant/database"
)

func TestUserSignup(t *testing.T) {
	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.SignUp()
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if un := database.GetUsersNumber(); un != tc.Data["UsersNumber"].(int) {
				t.Errorf("%s, wrong users number: expected %d, got %d", tc.Description, tc.Data["UsersNumber"], un)
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "username length not checked",
				User:        &database.User{Username: "u"},
				Expected:    database.ERR_USER_NAME_TOO_SHORT,
				Data:        map[string]any{"UsersNumber": 0},
			},
			{
				Description: "email not checked",
				User:        &database.User{Username: "username", Email: "email", Password: "p"},
				Expected:    database.ERR_USER_MAIL_INVALID,
				Data:        map[string]any{"UsersNumber": 0},
			},
			{
				Description: "password length not checked",
				User:        &database.User{Username: "username", Email: "email@email.com", Password: "p"},
				Expected:    database.ERR_USER_PASS_TOO_SHORT,
				Data:        map[string]any{"UsersNumber": 0},
			},
			{
				Description: "could not sign up",
				User:        &database.User{Username: "username", Email: "email@email.com", Password: "password"},
				Expected:    nil,
				Data:        map[string]any{"UsersNumber": 1},
			},
		},
	}.Run(t)
}

func TestUserSignIn(t *testing.T) {
	user := GetTestingUser(t)

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.SignIn()
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.User.UID != tc.Data["UID"].(int) {
				t.Errorf("%s, wrong uid: expected %d, got %d", tc.Description, tc.Data["UID"], tc.User.UID)
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "signed in unknown user",
				User:        &database.User{Username: "", Password: ""},
				Expected:    database.ERR_USER_WRONG_CREDENTIALS,
				Data:        map[string]any{"UID": 0},
			},
			{
				Description: "signed in with wrong password",
				User:        &database.User{Username: user.Username, Password: "pa$$word"},
				Expected:    database.ERR_USER_WRONG_CREDENTIALS,
				Data:        map[string]any{"UID": 0},
			},
			{
				Description: "could not sign in",
				User:        &database.User{Username: user.Username, Password: user.Password},
				Expected:    nil,
				Data:        map[string]any{"UID": user.UID},
			},
		},
	}.Run(t)
}

func TestGetUserData(t *testing.T) {
	user := GetTestingUser(t)
	user.Password = ""

	TestSuite[database.User]{
		Target: func(tc *TestCase[database.User]) database.User {
			return database.GetUserData(tc.User.UID)
		},

		PostCheck: func(t *testing.T, tc *TestCase[database.User]) {},

		Cases: []TestCase[database.User]{
			{
				Description: "got some random data",
				User:        &database.User{UID: 0},
				Expected:    database.User{},
				Data:        nil,
			},
			{
				Description: "could not get user data",
				User:        &database.User{UID: user.UID},
				Expected:    user,
				Data:        nil,
			},
		},
	}.Run(t)

}
