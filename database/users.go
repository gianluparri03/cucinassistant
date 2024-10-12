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

// compareHash compare a plain text with its hash. If they
// do not match, it returns NoMatch.
func compareHash(plain string, hash string, noMatch error) error {
	match, err := argon2id.ComparePasswordAndHash(plain, hash)
	if err != nil {
		slog.Error("while hashing string", "err", err)
		return ERR_UNKNOWN
	} else if !match {
		return noMatch
	} else {
		return nil
	}
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

// User represents a registered User
type User struct {
	// UID is the User ID.
	// It is the only required fields when using
	// receivers.
	UID int

	Username string
	Password string
	Email    string

	// Token is an optional string, that
	// can be generated to delete an user
	// or to reset its password
	Token string

	// fetched is true if the user if fetched from
	// the database, and not builded by hand
	fetched bool
}

// GetUser returns the user with the given field.
func GetUser(field string, value any) (user User, err error) {
	// Makes sure the field is valid
	if field != "UID" && field != "username" && field != "email" {
		slog.Error("invalid user identifier", "field", field)
		err = ERR_UNKNOWN
		return
	}

	// Prepares the user
	user.fetched = true
	var token *string

	// Queries the data
	err = db.QueryRow(`SELECT uid, username, email, password, token FROM ca_users WHERE `+field+`=$1;`, value).
		Scan(&user.UID, &user.Username, &user.Email, &user.Password, &token)
	if err != nil {
		// Checks the error
		if !strings.HasSuffix(err.Error(), "no rows in result set") {
			slog.Error("while retrieving user data:", "err", err)
			err = ERR_UNKNOWN
		} else {
			err = ERR_USER_UNKNOWN
		}

		return
	}

	// Dereferences token (can be null)
	if token == nil {
		user.Token = ""
	} else {
		user.Token = *token
	}

	return
}

// SignUp tries to sign up an user.
func SignUp(username string, email string, password string) (user User, err error) {
	// Checks if data is valid
	if err = checkUsername(username); err != nil {
		return
	} else if err = checkEmail(email); err != nil {
		return
	} else if err = checkPassword(password); err != nil {
		return
	}

	// Hashes the password
	hash, err := createHash(password)
	if err != nil {
		return
	}

	// Tries to save it in the database
	_, err = db.Exec(`INSERT INTO ca_users (username, email, password) VALUES ($1, $2, $3);`, username, email, hash)
	if err != nil {
		slog.Error("while signup:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Retrieves the UID
	user, err = GetUser("email", email)
	return
}

// SignIn tries to sign in an user.
func SignIn(username string, password string) (user User, err error) {
	// Fetches the hash
	var fetched User
	if fetched, err = GetUser("username", username); err != nil {
		err = ERR_USER_WRONG_CREDENTIALS
		return
	}

	// Compare the passwords
	err = compareHash(password, fetched.Password, ERR_USER_WRONG_CREDENTIALS)
	if err != nil {
		return
	}

	user = fetched
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
			*u = uFetched
			return nil
		} else {
			return err
		}
	}

	return nil
}

// ChangeUsername changes the user's username with a new one.
func (u *User) ChangeUsername(newUsername string) (err error) {
	// Ensures all data is present
	if err = u.ensureFetched(); err != nil {
		return
	}

	// Ensures the username is actually new, and that it's valid
	if u.Username == newUsername {
		return
	} else if err = checkUsername(newUsername); err != nil {
		return
	}

	// Saves the new one
	_, err = db.Exec(`UPDATE ca_users SET username=$2 WHERE uid=$1;`, u.UID, newUsername)
	if err != nil {
		slog.Error("while changing user username:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Updates struct
	u.Username = newUsername
	return
}

// ChangeEmail changes the user's email with a new one.
func (u *User) ChangeEmail(newEmail string) (err error) {
	// Ensures all data is present
	if err = u.ensureFetched(); err != nil {
		return
	}

	// Ensures the email is actually new, and that it's valid
	if u.Email == newEmail {
		return nil
	} else if err = checkEmail(newEmail); err != nil {
		return err
	}

	// Saves the new one
	_, err = db.Exec(`UPDATE ca_users SET email=$2 WHERE uid=$1;`, u.UID, newEmail)
	if err != nil {
		slog.Error("while changing user email:", "err", err)
		return ERR_UNKNOWN
	}

	// Updates struct
	u.Email = newEmail
	return nil
}

// ChangePassword changes the user's password with a new one.
func (u *User) ChangePassword(oldPassword string, newPassword string) (err error) {
	// Ensures all data is present
	if err = u.ensureFetched(); err != nil {
		return
	}

	// Checks if the new one is valid
	if err = checkPassword(newPassword); err != nil {
		return err
	}

	// Compares the old passwords
	err = compareHash(oldPassword, u.Password, ERR_USER_WRONG_CREDENTIALS)
	if err != nil {
		return
	}

	// Hashes the new one
	hashedPassword, err := createHash(newPassword)
	if err != nil {
		return err
	}

	// Saves the new one
	_, err = db.Exec(`UPDATE ca_users SET password=$2 WHERE uid=$1;`, u.UID, hashedPassword)
	if err != nil {
		slog.Error("while changing user password:", "err", err)
		return ERR_UNKNOWN
	}

	// Updates struct
	u.Password = hashedPassword
	return nil
}

// GenerateToken generates a new token for the user, then returns it.
// The Token field will be overwritten with its hash. The plain text
// is returned.
func (u *User) GenerateToken() (token string, err error) {
	// Generates the token
	token, hash, err := generateToken()

	// Saves it in the database
	var res sql.Result
	res, err = db.Exec(`UPDATE ca_users SET token=$2 WHERE uid=$1;`, u.UID, hash)
	if err != nil {
		slog.Error("while saving token:", "err", err)
		return "", ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra != 1 {
		return "", ERR_USER_UNKNOWN
	}

	// Updates the struct
	u.Token = hash
	return
}

// ResetPassword tries to reset the password of the user whose
// is picked from the struct.
// Note: the required field of the struct is Email, not UID.
func (u *User) ResetPassword(token string, newPassword string) (err error) {
	// Checks if password is valid
	if err = checkPassword(newPassword); err != nil {
		return
	}

	// Fetches the user
	*u, err = GetUser("email", u.Email)
	if err != nil {
		return
	} else if u.Token == "" {
		err = ERR_USER_WRONG_TOKEN
		return
	}

	// Compares the tokens
	if err = compareHash(token, u.Token, ERR_USER_WRONG_TOKEN); err != nil {
		return
	}

	// Hashes the password
	var hashedPassword string
	hashedPassword, err = createHash(newPassword)
	if err != nil {
		return
	}

	// Saves the new password (and resets the token)
	_, err = db.Exec(`UPDATE ca_users SET password=$2, token=NULL WHERE uid=$1;`, u.UID, hashedPassword)
	if err != nil {
		slog.Error("while resetting password:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	// Updates the struct
	u.Token = ""
	u.Password = hashedPassword
	return
}

// Delete deletes the user and all of its content
func (u User) Delete(token string) (err error) {
	// Ensures all data is present and the token
	// has been generated
	if err = u.ensureFetched(); err != nil {
		return
	} else if u.Token == "" {
		err = ERR_USER_WRONG_TOKEN
		return
	}

	// Compares the tokens
	if err = compareHash(token, u.Token, ERR_USER_WRONG_TOKEN); err != nil {
		return
	}

	// Deletes the user
	_, err = db.Exec(`DELETE FROM ca_users WHERE uid=$1;`, u.UID)
	if err != nil {
		slog.Error("while deleting user:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	return
}

// GetUsersNumber returns the number of users currently registered
func GetUsersNumber() (n int) {
	err := db.QueryRow(`SELECT COUNT(*) FROM ca_users;`).Scan(&n)
	if err != nil {
		slog.Error("while selecting users number:", "err", err)
	}

	return
}

// GetUsersEmails returns the emails of all the users
func GetUsersEmails() (emails []string) {
	rows, err := db.Query(`SELECT email FROM ca_users;`)
	if err != nil {
		slog.Error("while selecting emails:", "err", err)
	}

	defer rows.Close()
	for rows.Next() {
		var email string
		rows.Scan(&email)
		emails = append(emails, email)
	}

	return
}
