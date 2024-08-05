package database

import (
	"database/sql"
	"log/slog"
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
	// can be generated to delete an user
	// or to reset its password
	Token string
}

// SignUp tries to sign up an user. The required fields are Username,
// Used fields: Username, Password, Email
// Updated fields: UID, Password (with the hash)
func (u *User) SignUp() error {
	// Checks if data is valid
	if err := checkUsername(u.Username); err != nil {
		return err
	} else if err := checkEmail(u.Email); err != nil {
		return err
	} else if err := checkPassword(u.Password); err != nil {
		return err
	}

	// Generates the UID
	var uid int
	DB.QueryRow(`SELECT IFNULL(MAX(uid), 0) + 1 FROM users;`).Scan(&uid)

	// Hashes the password
	hash, err := createHash(u.Password)
	if err != nil {
		return err
	}

	// Tries to save it in the database
	if _, err = DB.Exec(`INSERT INTO users (uid, username, email, password) VALUES (?, ?, ?, ?);`, uid, u.Username, u.Email, hash); err != nil {
		slog.Error("while signup:", "err", err)
		return ERR_UNKNOWN
	}

	u.UID = uid
	u.Password = hash
	return nil
}

// SignIn tries to sign in an user.
// Used fields: Username, Password
// Updated fields: UID, Password (cleared)
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

	// Compare the passwords
	if match, err := compareHash(u.Password, hash); err != nil {
		return err
	} else if !match {
		return ERR_USER_WRONG_CREDENTIALS
	}

	// Updates the struct
	u.UID = uid
	u.Password = ""
	return nil
}

// GenerateToken generates a new token for the user, then returns it.
// The Token field will be overwritten with its hash. The plain text
// is returned.
// Used fields: UID
// Updated fields: Token (with the hash)
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
// Used fields: Email, Token, Password (the new one)
// Updated fields: UID, Token (cleared), Password (with the hash)
func (u *User) ResetPassword() (err error) {
	// Checks if password is valid
	if err = checkPassword(u.Password); err != nil {
		return err
	}

	// Fetches the user
	var hashedToken *string
	var uid int
	err = DB.QueryRow(`SELECT uid, token FROM users WHERE email = ?;`, u.Email).Scan(&uid, &hashedToken)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			return ERR_USER_UNKNOWN
		} else {
			slog.Error("while retrieving data on reset:", "err", err)
			return ERR_UNKNOWN
		}
	} else if hashedToken == nil {
		return ERR_USER_WRONG_TOKEN
	}

	// Compares the tokens
	if match, err := compareHash(u.Token, *hashedToken); err != nil {
		return err
	} else if !match {
		return ERR_USER_WRONG_TOKEN
	}

	// Hashes the password
	hashedPassword, err := createHash(u.Password)
	if err != nil {
		return err
	}

	// Saves the new password (and resets the token)
	_, err = DB.Exec(`UPDATE users SET password=?, token=NULL WHERE uid = ?;`, hashedPassword, uid)
	if err != nil {
		slog.Error("while resetting password:", "err", err)
		return ERR_UNKNOWN
	}

	// Updates the struct
	u.UID = uid
	u.Token = ""
	u.Password = hashedPassword
	return nil
}

// Delete deletes the user and all of its content
// Used fields: UID, Token
func (u *User) Delete() (err error) {
	// Fetches the user
	var hashedToken *string
	var found int
	err = DB.QueryRow(`SELECT 1, token FROM users WHERE uid = ?;`, u.UID).Scan(&found, &hashedToken)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no rows in result set") {
			return ERR_USER_UNKNOWN
		} else {
			slog.Error("while retrieving data on reset:", "err", err)
			return ERR_UNKNOWN
		}
	} else if hashedToken == nil {
		return ERR_USER_WRONG_TOKEN
	}

	// Compares the tokens
	if match, err := compareHash(u.Token, *hashedToken); err != nil {
		return err
	} else if !match {
		return ERR_USER_WRONG_TOKEN
	}

	// Deletes the user
	_, err = DB.Exec(`DELETE FROM users WHERE uid = ?;`, u.UID)
	if err != nil {
		slog.Error("while deleting user:", "err", err)
		return ERR_UNKNOWN
	}

	return nil
}

// GetUser returns the user with the given UID.
// Fetched fields: UID, Username, Email
func GetUser(uid int) (u User, err error) {
	err = DB.QueryRow(`SELECT uid, username, email FROM users WHERE uid = ?;`, uid).Scan(&u.UID, &u.Username, &u.Email)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "no rows in result set") {
			slog.Error("while retrieving user data (from uid):", "err", err)
			err = ERR_USER_UNKNOWN
		} else {
			err = ERR_UNKNOWN
		}
	}

	return
}

// GetUserFromEmail returns the user with the given email.
// Fetched fields: UID, Username, Email
func GetUserFromEmail(email string) (u User, err error) {
	err = DB.QueryRow(`SELECT uid, username, email FROM users WHERE email = ?;`, email).Scan(&u.UID, &u.Username, &u.Email)
	if err != nil {
		if !strings.HasSuffix(err.Error(), "no rows in result set") {
			slog.Error("while retrieving user data (from email):", "err", err)
			err = ERR_USER_UNKNOWN
		} else {
			err = ERR_UNKNOWN
		}
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
