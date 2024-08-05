package database

import (
	"errors"
)

var (
	ERR_UNKNOWN = errors.New("Errore sconosciuto")

	ERR_USER_UNKNOWN           = errors.New("Utente sconosciuto")
	ERR_USER_NAME_TOO_SHORT    = errors.New("Nome utente non valido: lunghezza minima 5 caratteri")
	ERR_USER_NAME_UNAVAIL      = errors.New("Nome utente non disponibile")
	ERR_USER_MAIL_INVALID      = errors.New("Email non valida")
	ERR_USER_MAIL_UNAVAIL      = errors.New("Email non disponibile")
	ERR_USER_PASS_TOO_SHORT    = errors.New("Password non valida: lunghezza minima 8 caratteri")
	ERR_USER_WRONG_CREDENTIALS = errors.New("Credenziali non valide")
	ERR_USER_WRONG_TOKEN       = errors.New("Qualcosa è andato storto. Riprova.")
)
