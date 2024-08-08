package database

import (
	"crypto/rand"
	"fmt"
	"github.com/alexedwards/argon2id"
	"log/slog"
	"net/mail"
)

// checkUsername ensures an username is valid.
// To be valid, it must be at least 5-characters long
// and should not be already used by someone else.
func checkUsername(username string) error {
	if len(username) < 5 {
		return ERR_USER_NAME_TOO_SHORT
	}

	var found bool
	DB.QueryRow(`SELECT 1 FROM users WHERE username = ?;`, username).Scan(&found)
	if found {
		return ERR_USER_NAME_UNAVAIL
	}

	return nil
}

// checkEmail ensures an email is valid.
// To be valid, it must be an actual email address
// and should not be already used by someone else.
func checkEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return ERR_USER_MAIL_INVALID
	}

	var found int
	DB.QueryRow(`SELECT 1 FROM users WHERE email = ?;`, email).Scan(&found)
	if found > 0 {
		return ERR_USER_MAIL_UNAVAIL
	}

	return nil
}

// checkPassword ensures a new password is valid.
// To be valid, it must be at least 8-characters long.
func checkPassword(password string) error {
	if len(password) < 8 {
		return ERR_USER_PASS_TOO_SHORT
	} else {
		return nil
	}
}

// createHash returns the hash of a string
func createHash(plain string) (hash string, err error) {
	hash, err = argon2id.CreateHash(plain, argon2id.DefaultParams)
	if err != nil {
		err = ERR_UNKNOWN
	}

	return
}

// compareHash compare a plain text with its hash
func compareHash(plain string, hash string) (match bool, err error) {
	match, err = argon2id.ComparePasswordAndHash(plain, hash)
	if err != nil {
		slog.Error("while hashing string", "err", err)
		err = ERR_UNKNOWN
	}

	return
}

// generateToken generates a new token. It returns both the plaintext and the hash.
func generateToken() (plain string, hash string, err error) {
	// Generates the token
	buffer := make([]byte, 18)
	rand.Read(buffer)
	plain = fmt.Sprintf("%x", buffer)

	// Hashes it
	hash, err = argon2id.CreateHash(plain, argon2id.DefaultParams)
	if err != nil {
		err = ERR_UNKNOWN
	}

	return
}

// ensureFetched checks if an user is builded
// by hand or fetched form the database. In the
// first case, it replaces it with a fetched one.
func (u *User) ensureFetched() error {
	if u == nil {
		return ERR_USER_UNKNOWN
	}

	if !u.fetched {
		if uFetched, err := GetUser("UID", u.UID); err == nil {
			*u = *uFetched
			return nil
		} else {
			return err
		}
	}

	return nil
}
