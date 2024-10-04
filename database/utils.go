package database

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alexedwards/argon2id"
	"log/slog"
	"net/mail"
	"strings"
)

const (
	// menuDefaultName is the name given to new menus
	menuDefaultName = "Nuovo Menù"

	// mealSeparator is used to separate meals when packed
	mealSeparator = ";"

	// duplicatedMenuSuffix is added at the end of the name when
	// duplicating a menu
	duplicatedMenuSuffix = " (copia)"
)

// checkUsername ensures an username is valid.
// To be valid, it must be at least 5-characters long
// and should not be already used by someone else.
func checkUsername(username string) error {
	if len(username) < 5 {
		return ERR_USER_NAME_TOO_SHORT
	}

	var found bool
	db.QueryRow(`SELECT 1 FROM ca_users WHERE username=$1;`, username).Scan(&found)
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
	db.QueryRow(`SELECT 1 FROM ca_users WHERE email=$1;`, email).Scan(&found)
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

// CompareHash compare a plain text with its hash
func CompareHash(plain string, hash string) (match bool, err error) {
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
// It is used for testing purposes.
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

// handleNoRowsError is an utility function that does the following.
// If err is sql.ErrNoRows, checks if it happened because the user (with given
// UID) does not exist (and in this case it returns ERR_USER_UNKNOWN), or just because
// there are no rows (and in this case it returns ifExists).
// If err is not sql.ErrNoRows, it prints a log (with whileDoing) and returns
// ERR_UNKNOWN
func handleNoRowsError(err error, UID int, ifExist error, whileDoing string) error {
	if errors.Is(err, sql.ErrNoRows) {
		if _, err = GetUser("UID", UID); err == nil {
			return ifExist
		} else {
			return err
		}
	} else {
		slog.Error("while "+whileDoing+":", "err", err)
		return ERR_UNKNOWN
	}
}

// packMeals packs the 14 meals in a string
func packMeals(meals [14]string) string {
	var sb strings.Builder

	for index, meal := range meals {
		sb.WriteString(meal)

		if index < 13 {
			sb.WriteString(mealSeparator)
		}
	}

	return sb.String()
}

// unpackMeals unpacks a string in an array of meals
func unpackMeals(meals string) (out [14]string) {
	for index, meal := range strings.Split(meals, mealSeparator) {
		if index == 14 {
			break
		}

		out[index] = meal
	}

	return
}
