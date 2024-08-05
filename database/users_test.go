package database

import (
	"github.com/alexedwards/argon2id"
	"testing"
)

func TestUserSignup(t *testing.T) {
	user := generateTestingUser()

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.SignUp()
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			// Ensures the user number is correct
			if un := GetUsersNumber(); un != tc.Data["UsersNumber"].(int) {
				t.Errorf("%s, wrong users number: expected %d, got %d", tc.Description, tc.Data["UsersNumber"], un)
			}

			// Ensures the saved password hash is correct
			// (if the user has been registered)
			if tc.Expected == nil {
				var hash string
				DB.QueryRow(`SELECT password FROM users WHERE uid = ?;`, tc.User.UID).Scan(&hash)
				if match, err := argon2id.ComparePasswordAndHash(tc.Data["Password"].(string), hash); !match {
					t.Errorf("%s, password hash does not match: error: %v", tc.Description, err)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "username length not checked",
				User:        &User{Username: "u"},
				Expected:    ERR_USER_NAME_TOO_SHORT,
				Data:        map[string]any{"UsersNumber": 0},
			},
			{
				Description: "email not checked",
				User:        &User{Username: user.Username, Email: "email", Password: "p"},
				Expected:    ERR_USER_MAIL_INVALID,
				Data:        map[string]any{"UsersNumber": 0},
			},
			{
				Description: "password length not checked",
				User:        &User{Username: user.Username, Email: user.Email, Password: "p"},
				Expected:    ERR_USER_PASS_TOO_SHORT,
				Data:        map[string]any{"UsersNumber": 0},
			},
			{
				Description: "could not sign up",
				User:        &User{Username: user.Username, Email: user.Email, Password: user.Password},
				Expected:    nil,
				Data:        map[string]any{"UsersNumber": 1, "Password": user.Password},
			},
			{
				Description: "signed up with duplicated username",
				User:        &User{Username: user.Username, Email: user.Email + "+", Password: user.Password},
				Expected:    ERR_USER_NAME_UNAVAIL,
				Data:        map[string]any{"UsersNumber": 1},
			},
			{
				Description: "signed up with duplicated email",
				User:        &User{Username: user.Username + "+", Email: user.Email, Password: user.Password},
				Expected:    ERR_USER_MAIL_UNAVAIL,
				Data:        map[string]any{"UsersNumber": 1},
			},
			{
				Description: "could not sign up with duplicated password",
				User:        &User{Username: user.Username + "+", Email: user.Email + "+", Password: user.Password},
				Expected:    nil,
				Data:        map[string]any{"UsersNumber": 2, "Password": user.Password},
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
				User:        &User{Username: "", Password: ""},
				Expected:    ERR_USER_WRONG_CREDENTIALS,
				Data:        map[string]any{"UID": 0},
			},
			{
				Description: "signed in with wrong password",
				User:        &User{Username: user.Username, Password: "pa$$word"},
				Expected:    ERR_USER_WRONG_CREDENTIALS,
				Data:        map[string]any{"UID": 0},
			},
			{
				Description: "could not sign in",
				User:        &User{Username: user.Username, Password: user.Password},
				Expected:    nil,
				Data:        map[string]any{"UID": user.UID},
			},
		},
	}.Run(t)
}

func TestUserGenerateToken(t *testing.T) {
	user := GetTestingUser(t)

	TestSuite[error]{
		Target: func(tc *TestCase[error]) (got error) {
			tc.Data["Token"], got = tc.User.GenerateToken()
			return got
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			// Ensures the returned token matches the saved one
			if tc.Expected == nil {
				var hash string
				DB.QueryRow(`SELECT token FROM users WHERE uid = ?;`, tc.User.UID).Scan(&hash)
				if match, err := argon2id.ComparePasswordAndHash(tc.Data["Token"].(string), hash); !match {
					t.Errorf("%s, saved token does not match returned one: error: %v", tc.Description, err)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "generated token for unknown user",
				User:        &User{UID: -1},
				Expected:    ERR_USER_UNKNOWN,
				Data:        make(map[string]any),
			},
			{
				Description: "could not generate token",
				User:        &User{UID: user.UID},
				Expected:    nil,
				Data:        make(map[string]any),
			},
		},
	}.Run(t)
}

func TestUserResetPassword(t *testing.T) {
	userWithoutToken := GetTestingUser(t)
	userWithToken := GetTestingUser(t)
	token, _ := userWithToken.GenerateToken()

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.ResetPassword()
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				// Makes sure that the token is dropped, and the new password is saved
				var tHash *string
				var pHash string
				DB.QueryRow(`SELECT token, password FROM users WHERE uid = ?;`, tc.User.UID).Scan(&tHash, &pHash)

				if tHash != nil {
					t.Errorf("%s, token wasn't dropped as expected", tc.Description)
				} else if match, err := argon2id.ComparePasswordAndHash("newPassword", pHash); !match {
					t.Errorf("%s, new password doesn't match %v %s", tc.Description, err, pHash)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "reset password of unknown user",
				User:        &User{Email: "", Token: "", Password: "newPassword"},
				Expected:    ERR_USER_UNKNOWN,
				Data:        nil,
			},
			{
				Description: "reset password without the token",
				User:        &User{Email: userWithoutToken.Email, Token: "", Password: "newPassword"},
				Expected:    ERR_USER_WRONG_TOKEN,
				Data:        nil,
			},
			{
				Description: "reset password with wrong token",
				User:        &User{Email: userWithToken.Email, Token: token + "+", Password: "newPassword"},
				Expected:    ERR_USER_WRONG_TOKEN,
				Data:        nil,
			},
			{
				Description: "reset password with an invalid one",
				User:        &User{Email: userWithToken.Email, Token: token, Password: "p"},
				Expected:    ERR_USER_PASS_TOO_SHORT,
				Data:        nil,
			},
			{
				Description: "could not reset password",
				User:        &User{Email: userWithToken.Email, Token: token, Password: "newPassword"},
				Expected:    nil,
				Data:        nil,
			},
		},
	}.Run(t)
}

func TestUserDelete(t *testing.T) {
	userWithoutToken := GetTestingUser(t)
	userWithToken := GetTestingUser(t)
	token, _ := userWithToken.GenerateToken()

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			// TODO: create some content, to check the foreign keys cascade
			return tc.User.Delete()
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				// Makes sure that the token is dropped, and the new password is saved
				var found int
				DB.QueryRow(`SELECT 1 FROM users WHERE uid = ?;`, tc.User.UID).Scan(&found)

				if found > 0 {
					t.Errorf("%s, user wasn't deleted", tc.Description)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "deleted unknown user",
				User:        &User{UID: 0, Token: ""},
				Expected:    ERR_USER_UNKNOWN,
				Data:        nil,
			},
			{
				Description: "deleted user without the token",
				User:        &User{UID: userWithoutToken.UID, Token: ""},
				Expected:    ERR_USER_WRONG_TOKEN,
				Data:        nil,
			},
			{
				Description: "deleted user with wrong token",
				User:        &User{UID: userWithToken.UID, Token: token + "+"},
				Expected:    ERR_USER_WRONG_TOKEN,
				Data:        nil,
			},
			{
				Description: "could not delete user",
				User:        &User{UID: userWithToken.UID, Token: token},
				Expected:    nil,
				Data:        nil,
			},
		},
	}.Run(t)
}

func TestGetUser(t *testing.T) {
	user := GetTestingUser(t)
	user.Password = ""

	TestSuite[User]{
		Target: func(tc *TestCase[User]) User {
			got, _ := GetUser(tc.User.UID)
			return got
		},

		PostCheck: func(t *testing.T, tc *TestCase[User]) {},

		Cases: []TestCase[User]{
			{
				Description: "got some random data",
				User:        &User{UID: 0},
				Expected:    User{},
				Data:        nil,
			},
			{
				Description: "could not get user data",
				User:        &User{UID: user.UID},
				Expected:    user,
				Data:        nil,
			},
		},
	}.Run(t)
}

func TestGetUserFromEmail(t *testing.T) {
	user := GetTestingUser(t)
	user.Password = ""

	TestSuite[User]{
		Target: func(tc *TestCase[User]) User {
			got, _ := GetUserFromEmail(tc.User.Email)
			return got
		},

		PostCheck: func(t *testing.T, tc *TestCase[User]) {},

		Cases: []TestCase[User]{
			{
				Description: "got some random data",
				User:        &User{Email: "email@email.net"},
				Expected:    User{},
				Data:        nil,
			},
			{
				Description: "could not get user data",
				User:        &User{Email: user.Email},
				Expected:    user,
				Data:        nil,
			},
		},
	}.Run(t)
}
