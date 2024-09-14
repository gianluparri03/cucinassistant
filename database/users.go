package database

import (
	"database/sql"
	"log/slog"
	"strings"
)

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

// SignUp tries to sign up an user.
func SignUp(username string, email string, password string) (user *User, err error) {
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
	_, err = DB.Exec(`INSERT INTO users (username, email, password) VALUES (?, ?, ?);`, username, email, hash)
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
func SignIn(username string, password string) (user *User, err error) {
	// Fetches the hash
	var user_ *User
	if user_, err = GetUser("username", username); err != nil {
		err = ERR_USER_WRONG_CREDENTIALS
		return
	}

	// Compare the passwords
	var match bool
	if match, err = compareHash(password, user_.Password); err != nil {
		return
	} else if !match {
		err = ERR_USER_WRONG_CREDENTIALS
		return
	}

	user = user_
	return
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
	_, err = DB.Exec(`UPDATE users SET username = ? WHERE uid = ?;`, newUsername, u.UID)
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
	_, err = DB.Exec(`UPDATE users SET email = ? WHERE uid = ?;`, newEmail, u.UID)
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
	if match, err := compareHash(oldPassword, u.Password); err != nil {
		return err
	} else if !match {
		return ERR_USER_WRONG_CREDENTIALS
	}

	// Hashes the new one
	hashedPassword, err := createHash(newPassword)
	if err != nil {
		return err
	}

	// Saves the new one
	_, err = DB.Exec(`UPDATE users SET password = ? WHERE uid = ?;`, hashedPassword, u.UID)
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
	res, err = DB.Exec(`UPDATE users SET token = ? WHERE uid = ?;`, hash, u.UID)
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
	u, err = GetUser("email", u.Email)
	if err != nil {
		return
	} else if u.Token == "" {
		err = ERR_USER_WRONG_TOKEN
		return
	}

	// Compares the tokens
	var match bool
	if match, err = compareHash(token, u.Token); err != nil {
		return
	} else if !match {
		err = ERR_USER_WRONG_TOKEN
		return
	}

	// Hashes the password
	var hashedPassword string
	hashedPassword, err = createHash(newPassword)
	if err != nil {
		return
	}

	// Saves the new password (and resets the token)
	_, err = DB.Exec(`UPDATE users SET password=?, token=NULL WHERE uid=?;`, hashedPassword, u.UID)
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
func (u *User) Delete(token string) (err error) {
	// Ensures the token is present
	u, err = GetUser("UID", u.UID)
	if err = u.ensureFetched(); err != nil {
		return
	} else if u.Token == "" {
		err = ERR_USER_WRONG_TOKEN
		return
	}

	// Compares the tokens
	var match bool
	if match, err = compareHash(token, u.Token); err != nil {
		return
	} else if !match {
		err = ERR_USER_WRONG_TOKEN
		return
	}

	// Deletes the user
	_, err = DB.Exec(`DELETE FROM users WHERE uid = ?;`, u.UID)
	if err != nil {
		slog.Error("while deleting user:", "err", err)
		err = ERR_UNKNOWN
		return
	}

	return
}

// GetUser returns the user with the given field.
func GetUser(field string, value any) (user *User, err error) {
	// Makes sure the field is valid
	if field != "UID" && field != "username" && field != "email" {
		slog.Error("invalid user identifier", "field", field)
		err = ERR_UNKNOWN
		return
	}

	// Prepares the user
	user = &User{fetched: true}
	var token *string

	// Queries the data
	err = DB.QueryRow(`SELECT uid, username, email, password, token FROM users WHERE `+field+` = ?;`, value).
		Scan(&user.UID, &user.Username, &user.Email, &user.Password, &token)
	if err != nil {
		user = nil

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

// GetUsersNumber returns the number of users currently registered
func GetUsersNumber() (n int) {
	err := DB.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&n)
	if err != nil {
		slog.Error("while selecting users number:", "err", err)
	}

	return
}
