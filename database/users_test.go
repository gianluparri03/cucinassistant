package database

import (
	"testing"
)

func TestSignup(t *testing.T) {
	user := generateTestingUser()
	baseUN := GetUsersNumber()

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			user, err := SignUp(tc.User.Username, tc.User.Email, tc.User.Password)
			if err == nil {
				tc.Data["UID"] = user.UID
			}
			return err
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			// Ensures the users number is correct
			expectedUN := tc.Data["UsersNumber"].(int)
			gotUN := GetUsersNumber()
			if expectedUN != gotUN {
				t.Errorf("%s, wrong users number: expected %d got %d", tc.Description, expectedUN, gotUN)
			}

			// Ensures the saved password hash is correct
			if tc.Expected == nil {
				user, _ := GetUser("username", tc.User.Username)
				if match, _ := compareHash(tc.Data["Password"].(string), user.Password); !match {
					t.Errorf("%s, password hash does not match", tc.Description)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "username length not checked",
				User:        &User{Username: "u"},
				Expected:    ERR_USER_NAME_TOO_SHORT,
				Data:        map[string]any{"UsersNumber": baseUN},
			},
			{
				Description: "email not checked",
				User:        &User{Username: user.Username, Email: "email", Password: "p"},
				Expected:    ERR_USER_MAIL_INVALID,
				Data:        map[string]any{"UsersNumber": baseUN},
			},
			{
				Description: "password length not checked",
				User:        &User{Username: user.Username, Email: user.Email, Password: "p"},
				Expected:    ERR_USER_PASS_TOO_SHORT,
				Data:        map[string]any{"UsersNumber": baseUN},
			},
			{
				Description: "could not sign up",
				User:        &User{Username: user.Username, Email: user.Email, Password: user.Password},
				Expected:    nil,
				Data:        map[string]any{"UsersNumber": baseUN + 1, "Password": user.Password},
			},
			{
				Description: "signed up with duplicated username",
				User:        &User{Username: user.Username, Email: user.Email + "+", Password: user.Password},
				Expected:    ERR_USER_NAME_UNAVAIL,
				Data:        map[string]any{"UsersNumber": baseUN + 1},
			},
			{
				Description: "signed up with duplicated email",
				User:        &User{Username: user.Username + "+", Email: user.Email, Password: user.Password},
				Expected:    ERR_USER_MAIL_UNAVAIL,
				Data:        map[string]any{"UsersNumber": baseUN + 1},
			},
			{
				Description: "could not sign up with duplicated password",
				User:        &User{Username: user.Username + "+", Email: user.Email + "+", Password: user.Password},
				Expected:    nil,
				Data:        map[string]any{"UsersNumber": baseUN + 2, "Password": user.Password},
			},
		},
	}.Run(t)
}

func TestSignIn(t *testing.T) {
	user, password := GetTestingUser(t)

	TestSuite[error]{
		Target: func(tc *TestCase[error]) error {
			user, err := SignIn(tc.User.Username, tc.User.Password)
			if err == nil {
				tc.Data["gotUID"] = user.UID
			}
			return err
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				expected := tc.Data["expectedUID"].(int)
				got := tc.Data["gotUID"].(int)
				if expected != got {
					t.Errorf("%s, wrong uid: expected %d, got %d", tc.Description, expected, got)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "signed in unknown user",
				User:        &User{Username: "", Password: ""},
				Expected:    ERR_USER_WRONG_CREDENTIALS,
				Data:        map[string]any{"expectedUID": 0},
			},
			{
				Description: "signed in with wrong password",
				User:        &User{Username: user.Username, Password: password + "+"},
				Expected:    ERR_USER_WRONG_CREDENTIALS,
				Data:        map[string]any{"expectedUID": 0},
			},
			{
				Description: "could not sign in",
				User:        &User{Username: user.Username, Password: password},
				Expected:    nil,
				Data:        map[string]any{"expectedUID": user.UID},
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
				user, _ := GetUser("UID", tc.User.UID)
				if user.Username != tc.Data["NewUsername"].(string) {
					t.Errorf("%s, new username isn't saved", tc.Description)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "changed username of unknown user",
				User:        &User{},
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
				user, _ := GetUser("UID", tc.User.UID)
				if user.Email != tc.Data["NewEmail"].(string) {
					t.Errorf("%s, new email isn't saved", tc.Description)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "changed email of unknown user",
				User:        &User{},
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
			return tc.User.ChangePassword(tc.Data["OldPassword"].(string), tc.Data["NewPassword"].(string))
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				// Makes sure the new password is saved
				user, _ := GetUser("UID", tc.User.UID)
				if match, _ := compareHash(tc.Data["NewPassword"].(string), user.Password); !match {
					t.Errorf("%s, new password doesn't match", tc.Description)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "changed password of unknown user",
				User:        &User{},
				Expected:    ERR_USER_UNKNOWN,
				Data:        map[string]any{"OldPassword": "", "NewPassword": "newPassword"},
			},
			{
				Description: "changed password with wrong old one",
				User:        user,
				Expected:    ERR_USER_WRONG_CREDENTIALS,
				Data:        map[string]any{"OldPassword": password + "+", "NewPassword": "newPassword"},
			},
			{
				Description: "changed password with an invalid one",
				User:        user,
				Expected:    ERR_USER_PASS_TOO_SHORT,
				Data:        map[string]any{"OldPassword": password, "NewPassword": "p"},
			},
			{
				Description: "could not change password",
				User:        user,
				Expected:    nil,
				Data:        map[string]any{"OldPassword": password, "NewPassword": "newPassword"},
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
				user, _ := GetUser("UID", tc.User.UID)
				if match, _ := compareHash(tc.Data["Token"].(string), user.Token); !match {
					t.Errorf("%s, saved token does not match returned one", tc.Description)
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
			return tc.User.ResetPassword(tc.Data["Token"].(string), tc.Data["NewPassword"].(string))
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			if tc.Expected == nil {
				// Makes sure that the token is dropped, and the new password is saved
				user, _ := GetUser("email", tc.User.Email)
				if user.Token != "" {
					t.Errorf("%s, token wasn't dropped as expected", tc.Description)
				} else if match, _ := compareHash(tc.Data["NewPassword"].(string), user.Password); !match {
					t.Errorf("%s, new password not saved", tc.Description)
				}
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "reset password of unknown user",
				User:        &User{Email: ""},
				Expected:    ERR_USER_UNKNOWN,
				Data:        map[string]any{"Token": "", "NewPassword": "newPassword"},
			},
			{
				Description: "reset password without the token",
				User:        &User{Email: userWithoutToken.Email},
				Expected:    ERR_USER_WRONG_TOKEN,
				Data:        map[string]any{"Token": "", "NewPassword": "newPassword"},
			},
			{
				Description: "reset password with wrong token",
				User:        &User{Email: userWithToken.Email},
				Expected:    ERR_USER_WRONG_TOKEN,
				Data:        map[string]any{"Token": token + "+", "NewPassword": "newPassword"},
			},
			{
				Description: "reset password with an invalid one",
				User:        &User{Email: userWithToken.Email},
				Expected:    ERR_USER_PASS_TOO_SHORT,
				Data:        map[string]any{"Token": token, "NewPassword": "p"},
			},
			{
				Description: "could not reset password",
				User:        &User{Email: userWithToken.Email},
				Expected:    nil,
				Data:        map[string]any{"Token": token, "NewPassword": "newPassword"},
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
			return tc.User.Delete(tc.Data["Token"].(string))
		},

		PostCheck: func(t *testing.T, tc *TestCase[error]) {
			user, _ := GetUser("UID", tc.User.UID)
			if (tc.Expected == nil) && (user != nil) {
				t.Errorf("%s, user wasn't deleted", tc.Description)
			} else if (tc.User.UID > 0) && (tc.Expected != nil) && (user == nil) {
				t.Errorf("%s, user wasn deleted anyway", tc.Description)
			}
		},

		Cases: []TestCase[error]{
			{
				Description: "deleted unknown user",
				User:        &User{},
				Expected:    ERR_USER_UNKNOWN,
				Data:        map[string]any{"Token": ""},
			},
			{
				Description: "deleted user without the token",
				User:        &User{UID: userWithoutToken.UID},
				Expected:    ERR_USER_WRONG_TOKEN,
				Data:        map[string]any{"Token": ""},
			},
			{
				Description: "deleted user with wrong token",
				User:        &User{UID: userWithToken.UID},
				Expected:    ERR_USER_WRONG_TOKEN,
				Data:        map[string]any{"Token": token + "+"},
			},
			{
				Description: "could not delete user",
				User:        &User{UID: userWithToken.UID},
				Expected:    nil,
				Data:        map[string]any{"Token": token},
			},
		},
	}.Run(t)
}

func TestGetUser(t *testing.T) {
	user, _ := GetTestingUser(t)

	TestSuite[Pair[*User, error]]{
		Target: func(tc *TestCase[Pair[*User, error]]) Pair[*User, error] {
			u, e := GetUser(tc.Data["Field"].(string), tc.Data["Value"])
			return Pair[*User, error]{u, e}
		},

		Cases: []TestCase[Pair[*User, error]]{
			{
				Description: "got data of unknown user",
				Expected:    Pair[*User, error]{nil, ERR_USER_UNKNOWN},
				Data:        map[string]any{"Field": "UID", "Value": 0},
			},
			{
				Description: "wrong user data (got from uid)",
				Expected:    Pair[*User, error]{user, nil},
				Data:        map[string]any{"Field": "UID", "Value": user.UID},
			},
			{
				Description: "wrong user data (got from email)",
				Expected:    Pair[*User, error]{user, nil},
				Data:        map[string]any{"Field": "email", "Value": user.Email},
			},
			{
				Description: "wrong user data (got from username)",
				Expected:    Pair[*User, error]{user, nil},
				Data:        map[string]any{"Field": "username", "Value": user.Username},
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
