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

	ERR_MENU_NOT_FOUND = errors.New("Menù non trovato")

	ERR_SECTION_DUPLICATED         = errors.New("Esiste già una sezione con lo stesso nome")
	ERR_SECTION_NOT_FOUND          = errors.New("Sezione non trovata")
	ERR_ARTICLE_NOT_FOUND          = errors.New("Articolo non trovata")
	ERR_ARTICLE_QUANTITY_INVALID   = errors.New("Quantità non valida")
	ERR_ARTICLE_EXPIRATION_INVALID = errors.New("Scadenza non valida")
	ERR_ARTICLE_DUPLICATED         = errors.New("Esiste già un articolo con stesso nome e scadenza")

	ERR_ENTRY_NOT_FOUND  = errors.New("Elemento non trovato")
	ERR_ENTRY_DUPLICATED = errors.New("Elemento già in lista")
)
