package database

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/alexedwards/argon2id"
	"log/slog"
	"net/mail"
	"strings"
)

// User represents a registered User
type User struct {
	// UID is the User ID
	UID int

	Username string
	Password string
	Email    string

	// Token is an optional string, that
	// can be generated to delete an account
	// or to recover its password
	Token string
}

// SignUp tries to sign up an user. The required fields are Username,
// Password and Email. If the registration is successfull, the UID will be set.
// Password will be overwritten with its hash.
func (u *User) SignUp() error {
	// Ensures the username and the password are big enough
	if len(u.Username) < 5 {
		return ERR_USER_NAME_TOO_SHORT
	} else if _, err := mail.ParseAddress(u.Email); err != nil {
		return ERR_USER_MAIL_INVALID
	} else if len(u.Password) < 8 {
		return ERR_USER_PASS_TOO_SHORT
	}

	// Hashes the password
	hash, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	if err != nil {
		slog.Error("while hashing password:", "err", err)
		return ERR_UNKNOWN
	} else {
		u.Password = hash
	}

	// Tries to save it in the database
	_, err = DB.Exec(
		`INSERT INTO users (uid, username, email, password)
         SELECT IFNULL(MAX(uid), 0) + 1, ?, ?, ? from users;`,
		u.Username,
		u.Email,
		u.Password,
	)

	// Handles errors
	if err != nil {
		if strings.HasSuffix(err.Error(), "for key 'username'") {
			return ERR_USER_NAME_UNAVAIL
		} else if strings.HasSuffix(err.Error(), "for key 'email'") {
			return ERR_USER_MAIL_UNAVAIL
		} else {
			slog.Warn("while signup:", "err", err)
			return ERR_UNKNOWN
		}
	}

	// Updates the uid
	err = DB.QueryRow(`SELECT uid FROM users WHERE username = ?;`, u.Username).Scan(&u.UID)
	if err != nil {
		slog.Error("while retrieving uid on signup:", "err", err)
		return ERR_UNKNOWN
	}

	return nil
}

// SignIn tries to sign in an user. The required fields are Username and
// Password. If the login is successfull, the UID will be set.
// Password will be overwritten with an empty string.
func (u *User) SignIn() error {
	var uid int
	var hash string

	// Fetches the hash
	err := DB.QueryRow(`SELECT uid, password FROM users WHERE username = ?;`, u.Username).Scan(&uid, &hash)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			return ERR_USER_WRONG_CREDENTIALS
		} else {
			slog.Error("while retrieving data on signin:", "err", err)
			return ERR_UNKNOWN
		}
	}

	// Ensures the user has been found
	if uid == 0 {
		return ERR_USER_WRONG_CREDENTIALS
	}

	// Compare the passwords
	match, err := argon2id.ComparePasswordAndHash(u.Password, hash)
	if err != nil {
		slog.Error("while comparing hashes on signin:", "err", err)
		return ERR_UNKNOWN
	} else if !match {
		return ERR_USER_WRONG_CREDENTIALS
	}

	u.UID = uid
	u.Password = ""
	return nil
}

// GenerateToken generates a new token for the user, then returns it
// The Token field will be overwritten with its hash. The plain text
// is returned.
func (u *User) GenerateToken() (token string, err error) {
	// Generates the token
	buffer := make([]byte, 18)
	rand.Read(buffer)
	token = fmt.Sprintf("%x", buffer)

	// Hashes it
	var hash string
	hash, err = argon2id.CreateHash(token, argon2id.DefaultParams)
	if err != nil {
		slog.Error("while hashing token:", "err", err)
		return "", ERR_UNKNOWN
	}

	// Saves it in the database
	var res sql.Result
	res, err = DB.Exec(`UPDATE users SET token = ? WHERE uid = ?;`, hash, u.UID)
	if err != nil {
		slog.Error("while saving token:", "err", err)
		return "", ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra != 1 {
		return "", ERR_UNKNOWN
	}

	u.Token = hash
	return
}

// ResetPassword tries to reset the password of the user whose
// is picked from the struct. The required fields are Email and Token.
// UID will be set, and Token cleared, if the operation does not generate errors.
func (u *User) ResetPassword(password string) (err error) {
	var hashedToken *string
	var uid int

	// Fetches the user
	err = DB.QueryRow(`SELECT uid, token FROM users WHERE email = ?;`, u.Email).Scan(&uid, &hashedToken)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			return ERR_UNKNOWN
		} else {
			slog.Error("while retrieving data on reset:", "err", err)
			return ERR_UNKNOWN
		}
	} else if hashedToken == nil {
		return ERR_UNKNOWN
	}

	// Compares the tokens
	match, err := argon2id.ComparePasswordAndHash(u.Token, *hashedToken)
	if err != nil {
		slog.Error("while comparing hashes on password_reset:", "err", err)
		return ERR_UNKNOWN
	} else if !match {
		return ERR_UNKNOWN
	}

	// Checks password length
	if len(password) < 8 {
		return ERR_USER_PASS_TOO_SHORT
	}

	// Hashes it
	var hashedPassword string
	hashedPassword, err = argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		slog.Error("while hashing password:", "err", err)
		return ERR_UNKNOWN
	}

	// Saves everything
	_, err = DB.Exec(`UPDATE users SET password=?, token=NULL WHERE uid = ?;`, hashedPassword, uid)
	if err != nil {
		slog.Error("while resetting password:", "err", err)
	}

	u.UID = uid
	u.Token = ""
	return nil
}

// GetUser returns the user with the given UID. Password
// and Token fields are not fetched.
func GetUser(uid int) (u User, err error) {
	err = DB.QueryRow(`SELECT uid, username, email FROM users WHERE uid = ?;`, uid).Scan(&u.UID, &u.Username, &u.Email)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "no rows in result set") {
			slog.Error("while retrieving data:", "err", err)
		}

		return User{}, ERR_UNKNOWN
	}

	return
}

// GetUserFromEmail returns the user with the given email.
// Password and Token fields are not fetched.
func GetUserFromEmail(email string) (u User, err error) {
	err = DB.QueryRow(`SELECT uid, username, email FROM users WHERE email = ?;`, email).Scan(&u.UID, &u.Username, &u.Email)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "no rows in result set") {
			slog.Error("while retrieving data:", "err", err)
		}

		return User{}, ERR_UNKNOWN
	}

	return
}

// GetUsersNumber returns the number of users currently registered
func GetUsersNumber() (n int) {
	err := DB.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&n)
	if err != nil {
		slog.Error("while selecting users number:", "err", err)
	}

	return
}
