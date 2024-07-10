package database

import (
	"fmt"
	"github.com/alexedwards/argon2id"
	"log"
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
		return fmt.Errorf("Nome utente non valido: lunghezza minima 5 caratteri")
	} else if len(u.Password) < 8 {
		return fmt.Errorf("Password non valida: lunghezza minima 8 caratteri")
	}

	// Hashes the password
	hash, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal("ERR: " + err.Error())
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
			return fmt.Errorf("Nome utente non disponibile")
		} else if strings.HasSuffix(err.Error(), "for key 'email'") {
			return fmt.Errorf("Email non disponibile")
		} else {
			log.Print("ERR: " + err.Error())
			return fmt.Errorf("Errore sconosciuto")
		}
	}

	// Updates the uid
	DB.QueryRow(`SELECT uid FROM users WHERE username = ?;`, u.Username).Scan(&u.UID)
	return nil
}
