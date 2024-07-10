package database

import (
	"errors"
	"github.com/alexedwards/argon2id"
	"log/slog"
	"strings"
)

type User struct {
	UID int

	Username string
	Password string
	Email    string
	Token    string
}

func (u *User) SignUp() error {
	// Ensures the username and the password are big enough
	if len(u.Username) < 5 {
		return errors.New("Nome utente non valido: lunghezza minima 5 caratteri")
	} else if len(u.Password) < 8 {
		return errors.New("Password non valida: lunghezza minima 8 caratteri")
	}

	// Hashes the password
	hash, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	if err != nil {
		slog.Error("while hashing password:", "err", err)
		return errors.New("Errore sconosciuto")
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
			return errors.New("Nome utente non disponibile")
		} else if strings.HasSuffix(err.Error(), "for key 'email'") {
			return errors.New("Email non disponibile")
		} else {
			slog.Warn("while signup:", "err", err)
			return errors.New("Errore sconosciuto")
		}
	}

	// Updates the uid
	err = DB.QueryRow(`SELECT uid FROM users WHERE username = ?;`, u.Username).Scan(&u.UID)
	if err != nil {
		slog.Error("while retrieving uid on signup:", "err", err)
		return errors.New("Errore sconosciuto")
	}

	// Logs
	slog.Info("User signed up succesfully", "email", u.Email)

	return nil
}
