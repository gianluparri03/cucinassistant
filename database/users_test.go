package database

import (
	"github.com/alexedwards/argon2id"
	"testing"
)

func TestSignup(t *testing.T) {
	user := generateTestingUser()

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.SignUp()
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			// Ensures the user has been registered only if there
			// were no errors
			var found int
			DB.QueryRow(`SELECT 1 FROM users WHERE username = ? AND email = ?;`, tc.User.Username, tc.User.Email).Scan(&found)
			if tc.Expected == nil && found == 0 {
				t.Errorf("%s, not registered", tc.Description)
			} else if tc.Expected != nil && found > 0 {
				t.Errorf("%s, registered anyway", tc.Description)
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

func TestSignIn(t *testing.T) {
	user, password := GetTestingUser(t)

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
				User:        &User{Username: user.Username, Password: password + "+"},
				Expected:    ERR_USER_WRONG_CREDENTIALS,
				Data:        map[string]any{"UID": 0},
			},
			{
				Description: "could not sign in",
				User:        &User{Username: user.Username, Password: password},
				Expected:    nil,
				Data:        map[string]any{"UID": user.UID},
			},
		},
	}.Run(t)
}

func TestChangeUsername(t *testing.T) {
	user, _ := GetTestingUser(t)

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.ChangeUsername(tc.Data["NewUsername"].(string))
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				// Makes sure the new username is saved
				var username string
				DB.QueryRow(`SELECT username FROM users WHERE uid = ?;`, tc.User.UID).Scan(&username)
				if username != tc.Data["NewUsername"].(string) {
					t.Errorf("%s, new username isn't saved", tc.Description)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "changed username of unknown user",
				User:        &User{UID: 0},
				Expected:    ERR_USER_UNKNOWN,
				Data:        map[string]any{"NewUsername": "newUsername"},
			},
			{
				Description: "changed username with an invalid one",
				User:        &User{UID: user.UID},
				Expected:    ERR_USER_NAME_TOO_SHORT,
				Data:        map[string]any{"NewUsername": "u"},
			},
			{
				Description: "could not change username",
				User:        &User{UID: user.UID},
				Expected:    nil,
				Data:        map[string]any{"NewUsername": "newUsername"},
			},
			{
				Description: "could not keep username",
				User:        &User{UID: user.UID},
				Expected:    nil,
				Data:        map[string]any{"NewUsername": "newUsername"},
			},
		},
	}.Run(t)
}

func TestChangeEmail(t *testing.T) {
	user, _ := GetTestingUser(t)

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.ChangeEmail(tc.Data["NewEmail"].(string))
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				// Makes sure the new email is saved
				var email string
				DB.QueryRow(`SELECT email FROM users WHERE uid = ?;`, tc.User.UID).Scan(&email)
				if email != tc.Data["NewEmail"].(string) {
					t.Errorf("%s, new email isn't saved", tc.Description)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "changed email of unknown user",
				User:        &User{UID: 0},
				Expected:    ERR_USER_UNKNOWN,
				Data:        map[string]any{"NewEmail": "newemail@email.com"},
			},
			{
				Description: "changed email with an invalid one",
				User:        &User{UID: user.UID},
				Expected:    ERR_USER_MAIL_INVALID,
				Data:        map[string]any{"NewEmail": "e"},
			},
			{
				Description: "could not change email",
				User:        &User{UID: user.UID},
				Expected:    nil,
				Data:        map[string]any{"NewEmail": "newemail@email.com"},
			},
			{
				Description: "could not keep email",
				User:        &User{UID: user.UID},
				Expected:    nil,
				Data:        map[string]any{"NewEmail": "newemail@email.com"},
			},
		},
	}.Run(t)
}

func TestChangePassword(t *testing.T) {
	user, password := GetTestingUser(t)

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.ChangePassword(tc.Data["NewPassword"].(string))
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				// Makes sure the new password is saved
				var pHash string
				DB.QueryRow(`SELECT password FROM users WHERE uid = ?;`, tc.User.UID).Scan(&pHash)
				if match, err := argon2id.ComparePasswordAndHash(tc.Data["NewPassword"].(string), pHash); !match {
					t.Errorf("%s, new password doesn't match (err: %v)", tc.Description, err)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "changed password of unknown user",
				User:        &User{UID: 0, Password: ""},
				Expected:    ERR_USER_UNKNOWN,
				Data:        map[string]any{"NewPassword": "newPassword"},
			},
			{
				Description: "changed password with wrong old one",
				User:        &User{UID: user.UID, Password: password + "+"},
				Expected:    ERR_USER_WRONG_CREDENTIALS,
				Data:        map[string]any{"NewPassword": "newPassword"},
			},
			{
				Description: "changed password with an invalid one",
				User:        &User{UID: user.UID, Password: password + "+"},
				Expected:    ERR_USER_PASS_TOO_SHORT,
				Data:        map[string]any{"NewPassword": "p"},
			},
			{
				Description: "could not change password",
				User:        &User{UID: user.UID, Password: password},
				Expected:    nil,
				Data:        map[string]any{"NewPassword": "newPassword"},
			},
		},
	}.Run(t)
}

func TestGenerateToken(t *testing.T) {
	user, _ := GetTestingUser(t)

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

func TestResetPassword(t *testing.T) {
	userWithoutToken, _ := GetTestingUser(t)
	userWithToken, _ := GetTestingUser(t)
	token, _ := userWithToken.GenerateToken()

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			return tc.User.ResetPassword(tc.Data["NewPassword"].(string))
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				// Makes sure that the token is dropped, and the new password is saved
				var tHash *string
				var pHash string
				DB.QueryRow(`SELECT token, password FROM users WHERE uid = ?;`, tc.User.UID).Scan(&tHash, &pHash)

				if tHash != nil {
					t.Errorf("%s, token wasn't dropped as expected", tc.Description)
				} else if match, err := argon2id.ComparePasswordAndHash(tc.Data["NewPassword"].(string), pHash); !match {
					t.Errorf("%s, new password doesn't match (err: %v)", tc.Description, err)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "reset password of unknown user",
				User:        &User{Email: "", Token: ""},
				Expected:    ERR_USER_UNKNOWN,
				Data:        map[string]any{"NewPassword": "newPassword"},
			},
			{
				Description: "reset password without the token",
				User:        &User{Email: userWithoutToken.Email, Token: ""},
				Expected:    ERR_USER_WRONG_TOKEN,
				Data:        map[string]any{"NewPassword": "newPassword"},
			},
			{
				Description: "reset password with wrong token",
				User:        &User{Email: userWithToken.Email, Token: token + "+"},
				Expected:    ERR_USER_WRONG_TOKEN,
				Data:        map[string]any{"NewPassword": "newPassword"},
			},
			{
				Description: "reset password with an invalid one",
				User:        &User{Email: userWithToken.Email, Token: token},
				Expected:    ERR_USER_PASS_TOO_SHORT,
				Data:        map[string]any{"NewPassword": "p"},
			},
			{
				Description: "could not reset password",
				User:        &User{Email: userWithToken.Email, Token: token},
				Expected:    nil,
				Data:        map[string]any{"NewPassword": "newPassword"},
			},
		},
	}.Run(t)
}

func TestDeleteUser(t *testing.T) {
	userWithoutToken, _ := GetTestingUser(t)
	userWithToken, _ := GetTestingUser(t)
	token, _ := userWithToken.GenerateToken()

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			// TODO: create some content, to check the foreign keys cascade
			tc.User.AppendEntries("...")
			tc.User.NewMenu()
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
	user, _ := GetTestingUser(t)

	TestSuite[Pair[User, error]]{
		Target: func(tc *TestCase[Pair[User, error]]) Pair[User, error] {
			u, e := GetUser(tc.User.UID)
			return Pair[User, error]{u, e}
		},

		Cases: []TestCase[Pair[User, error]]{
			{
				Description: "got data of unknown user",
				User:        &User{UID: 0},
				Expected:    Pair[User, error]{User{}, ERR_USER_UNKNOWN},
			},
			{
				Description: "wrong user data",
				User:        &User{UID: user.UID},
				Expected:    Pair[User, error]{user, nil},
			},
		},
	}.Run(t)
}

func TestGetUserFromEmail(t *testing.T) {
	user, _ := GetTestingUser(t)

	TestSuite[Pair[User, error]]{
		Target: func(tc *TestCase[Pair[User, error]]) Pair[User, error] {
			u, e := GetUserFromEmail(tc.User.Email)
			return Pair[User, error]{u, e}
		},

		Cases: []TestCase[Pair[User, error]]{
			{
				Description: "got some random user's data (from email)",
				User:        &User{Email: "email@email.net"},
				Expected:    Pair[User, error]{User{}, ERR_USER_UNKNOWN},
			},
			{
				Description: "wrong user data (from email)",
				User:        &User{Email: user.Email},
				Expected:    Pair[User, error]{user, nil},
			},
		},
	}.Run(t)
}

func TestGetUsersNumber(t *testing.T) {
	TestSuite[int]{
		Target: func(tc *TestCase[int]) int {
			return GetUsersNumber()
		},

		Cases: []TestCase[int]{
			{
				Description: "wrong users number",
				Expected:    userN,
			},
		},
	}.Run(t)
}
