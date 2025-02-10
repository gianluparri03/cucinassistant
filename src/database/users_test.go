package database

import (
	"fmt"
	"reflect"
	"testing"
)

var testingUsersN int = 0

// generateTestingUser returns a testing user which has not been registered yet
func generateTestingUser() User {
	testingUsersN++

	return User{
		Username: fmt.Sprintf("username%d", testingUsersN),
		Email:    fmt.Sprintf("email%d@email.com", testingUsersN),
		Password: fmt.Sprintf("password%d", testingUsersN),
	}
}

// getTestingUser returns an user to be used for testing purposes
func getTestingUser(t *testing.T) (user User, password string) {
	user = generateTestingUser()
	password = user.Password

	// And tries to sign it up
	var err error
	if user, err = SignUp(user.Username, user.Email, user.Password); err != nil {
		t.Fatalf("Cannot create testing user: %s", err.Error())
	}

	return
}

var unknownUser User = User{}

func TestSignup(t *testing.T) {
	user := generateTestingUser()

	type data struct {
		Username string
		Email    string
		Password string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			preUN := GetStats().UsersNumber

			_, err := SignUp(d.Username, d.Email, d.Password)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			postUN := GetStats().UsersNumber

			// UsersNumber should be incremented only if there are no errors
			if (d.ExpectedErr != nil) != (preUN == postUN) {
				t.Errorf("%s, wrong users number", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"username length not checked",
				data{Username: "u", Email: user.Email, Password: user.Password, ExpectedErr: ERR_USER_NAME_TOO_SHORT},
			},
			{
				"email not checked",
				data{Username: user.Username, Email: "e", Password: user.Password, ExpectedErr: ERR_USER_MAIL_INVALID},
			},
			{
				"password not checked",
				data{Username: user.Username, Email: user.Email, Password: "p", ExpectedErr: ERR_USER_PASS_TOO_SHORT},
			},
			{
				"",
				data{Username: user.Username, Email: user.Email, Password: user.Password},
			},
			{
				"signed up with duplicated username",
				data{Username: user.Username, Email: user.Email + "e", Password: user.Password, ExpectedErr: ERR_USER_NAME_UNAVAIL},
			},
			{
				"signed up with duplicated email",
				data{Username: user.Username + "u", Email: user.Email, Password: user.Password, ExpectedErr: ERR_USER_MAIL_UNAVAIL},
			},
			{
				"",
				data{Username: user.Username + "u", Email: user.Email + "e", Password: user.Password},
			},
		},
	}.Run(t)
}

func TestSignIn(t *testing.T) {
	user, password := getTestingUser(t)

	type data struct {
		Username string
		Password string

		ExpectedErr error
		ExpectedUID int
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			user, err := SignIn(d.Username, d.Password)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if (d.ExpectedErr == nil) && (user.UID != d.ExpectedUID) {
				t.Errorf("%s, wrong uid", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"signed in unknown user",
				data{Username: user.Username + "u", Password: password, ExpectedErr: ERR_USER_WRONG_CREDENTIALS},
			},
			{
				"signed in with wrong password",
				data{Username: user.Username, Password: password + "p", ExpectedErr: ERR_USER_WRONG_CREDENTIALS},
			},
			{
				"",
				data{Username: user.Username, Password: password, ExpectedUID: user.UID},
			},
		},
	}.Run(t)
}

func TestChangeUsername(t *testing.T) {
	user, _ := getTestingUser(t)
	otherUser, _ := getTestingUser(t)

	type data struct {
		User        User
		NewUsername string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.ChangeUsername(d.NewUsername)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				user, _ := GetUser("UID", user.UID)
				if user.Username != d.NewUsername {
					t.Errorf("%s, changes not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"changed username of unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"username not checked",
				data{User: user, NewUsername: "u", ExpectedErr: ERR_USER_NAME_TOO_SHORT},
			},
			{
				"changed username with an unavailable one",
				data{User: user, NewUsername: otherUser.Username, ExpectedErr: ERR_USER_NAME_UNAVAIL},
			},
			{
				"(same)",
				data{User: user, NewUsername: user.Username},
			},
			{
				"(different)",
				data{User: user, NewUsername: user.Username + "u"},
			},
		},
	}.Run(t)
}

func TestChangeEmail(t *testing.T) {
	user, _ := getTestingUser(t)
	otherUser, _ := getTestingUser(t)

	type data struct {
		User     User
		NewEmail string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.ChangeEmail(d.NewEmail)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				user, _ := GetUser("UID", user.UID)
				if user.Email != d.NewEmail {
					t.Errorf("%s, changes not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"changed email of unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"email not checked",
				data{User: user, NewEmail: "e", ExpectedErr: ERR_USER_MAIL_INVALID},
			},
			{
				"changed email with an unavailable one",
				data{User: user, NewEmail: otherUser.Email, ExpectedErr: ERR_USER_MAIL_UNAVAIL},
			},
			{
				"(same)",
				data{User: user, NewEmail: user.Email},
			},
			{
				"(different)",
				data{User: user, NewEmail: user.Email + "e"},
			},
		},
	}.Run(t)
}

func TestChangePassword(t *testing.T) {
	user, password := getTestingUser(t)

	type data struct {
		User        *User
		OldPassword string
		NewPassword string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.ChangePassword(d.OldPassword, d.NewPassword)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				user, _ := GetUser("UID", d.User.UID)
				if err := compareHash(d.NewPassword, user.Password, ERR_UNKNOWN); err != nil {
					t.Errorf("%s, changes not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"changed password of unknown user",
				data{User: &unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"changed password with an invalid one",
				data{User: &user, NewPassword: "p", ExpectedErr: ERR_USER_PASS_TOO_SHORT},
			},
			{
				"changed password with wrong old one",
				data{User: &user, OldPassword: "p", NewPassword: password + "p", ExpectedErr: ERR_USER_WRONG_CREDENTIALS},
			},
			{
				"",
				data{User: &user, OldPassword: password, NewPassword: password + "p"},
			},
			{
				"",
				data{User: &user, OldPassword: password + "p", NewPassword: password + "p"},
			},
		},
	}.Run(t)
}

func TestGenerateToken(t *testing.T) {
	user, _ := getTestingUser(t)

	type data struct {
		User        User
		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			token, err := d.User.GenerateToken()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				user, _ := GetUser("UID", d.User.UID)
				if err := compareHash(token, user.Token, ERR_UNKNOWN); err != nil {
					t.Errorf("%s, saved token does not match returned one", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"generated token of unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{User: user},
			},
		},
	}.Run(t)
}

func TestResetPassword(t *testing.T) {
	user, password := getTestingUser(t)
	token, _ := user.GenerateToken()
	otherUser, _ := getTestingUser(t)

	type data struct {
		User        User
		Token       string
		NewPassword string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.ResetPassword(d.Token, d.NewPassword)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				user, _ := GetUser("email", d.User.Email)
				if user.Token != "" {
					t.Errorf("%s, token wasn't dropped", msg)
				} else if err := compareHash(d.NewPassword, user.Password, ERR_UNKNOWN); err != nil {
					t.Errorf("%s, new password not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"reset password of unknown user",
				data{User: unknownUser, NewPassword: password, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"reset password of unknown user",
				data{User: otherUser, NewPassword: password, ExpectedErr: ERR_USER_WRONG_TOKEN},
			},
			{
				"reset password with wrong token",
				data{User: user, Token: token + "t", NewPassword: password, ExpectedErr: ERR_USER_WRONG_TOKEN},
			},
			{
				"reset password with an invalid one",
				data{User: user, Token: token, NewPassword: "p", ExpectedErr: ERR_USER_PASS_TOO_SHORT},
			},
			{
				"",
				data{User: user, Token: token, NewPassword: password},
			},
		},
	}.Run(t)
}

func TestDeleteUser(t *testing.T) {
	user, _ := getTestingUser(t)
	token, _ := user.GenerateToken()

	user.ShoppingList().Append("e")
	user.Menus().New("m")
	section, _ := user.Storage().NewSection("s")
	user.Storage().AddArticles(StringArticle{Name: "article", Section: section.SID})
	testingArticlesN++
	user.Recipes().New("r")

	otherUser, _ := getTestingUser(t)

	type data struct {
		User  User
		Token string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.Delete(d.Token)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			_, err = GetUser("UID", d.User.UID)
			if (d.ExpectedErr == nil) && (err == nil) {
				t.Errorf("%s, user wasn't deleted", msg)
			} else if (d.User.UID > 0) && (d.ExpectedErr != nil) && (err != nil) {
				t.Errorf("%s, user was deleted anyway", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"deleted unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"deleted user without the token",
				data{User: otherUser, ExpectedErr: ERR_USER_WRONG_TOKEN},
			},
			{
				"deleted user with wrong token",
				data{User: user, Token: token + "t", ExpectedErr: ERR_USER_WRONG_TOKEN},
			},
			{
				"",
				data{User: user, Token: token},
			},
		},
	}.Run(t)
}

func TestGetUser(t *testing.T) {
	user, _ := getTestingUser(t)

	type data struct {
		Field string
		Value any

		ExpectedErr  error
		ExpectedUser User
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			user, err := GetUser(d.Field, d.Value)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(user, d.ExpectedUser) {
				t.Errorf("%s: expected user <%v>, got <%v>", msg, d.ExpectedUser, user)
			}
		},

		Cases: []testCase[data]{
			{
				"got data of unknown user",
				data{Field: "UID", Value: 0, ExpectedUser: User{fetched: true}, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"(UID)",
				data{Field: "UID", Value: user.UID, ExpectedUser: user},
			},
			{
				"(username)",
				data{Field: "username", Value: user.Username, ExpectedUser: user},
			},
			{
				"(email)",
				data{Field: "email", Value: user.Email, ExpectedUser: user},
			},
		},
	}.Run(t)
}

func TestSetEmailLang(t *testing.T) {
	user, _ := getTestingUser(t)

	type data struct {
		User    User
		NewLang string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.SetEmailLang(d.NewLang)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				user, _ := GetUser("UID", user.UID)
				if user.EmailLang != d.NewLang {
					t.Errorf("%s, changes not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"changed email_lang of unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{User: user, NewLang: "it"},
			},
		},
	}.Run(t)
}
