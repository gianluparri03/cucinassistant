package database

import (
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
			return ERR_USER_NAME_UNAVAIL
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

// GetUser returns the user with the given UID. Password
// and Token fields are not fetched.
func GetUser(uid int) (u User) {
	DB.QueryRow(`SELECT username, email FROM users WHERE uid = ?;`, uid).Scan(&u.Username, &u.Email)
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
