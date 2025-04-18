package langs

import (
	"cucinassistant/database"
)

var italian *Lang = &Lang{
	Tag:  "it",
	Name: "Italiano",

	Strings: map[String]string{
		STR_DELETE_CONFIRM_EMAIL:                        "per eliminare definitivamente il tuo account,",
		STR_GOODBYE_EMAIL:                               "il tuo account è stato eliminato definitivamente.",
		STR_PASSWORD_CHANGED_EMAIL:                      "la tua password è stata cambiata di recente.",
		STR_RESET_PASSWORD_EMAIL:                        "per resettare la tua password,",
		STR_WELCOME_EMAIL:                               "Benvenuto/a su CucinAssistant!",
		String(database.ERR_ARTICLE_DUPLICATED):         "Esiste già un articolo con stesso nome e scadenza",
		String(database.ERR_ARTICLE_EXPIRATION_INVALID): "Scadenza non valida",
		String(database.ERR_ARTICLE_NOT_FOUND):          "Articolo non trovata",
		String(database.ERR_ARTICLE_QUANTITY_INVALID):   "Quantità non valida",
		String(database.ERR_ENTRY_DUPLICATED):           "Elemento già in lista",
		String(database.ERR_ENTRY_NOT_FOUND):            "Elemento non trovato",
		String(database.ERR_MENU_NOT_FOUND):             "Menù non trovato",
		String(database.ERR_RECIPE_DUPLICATED):          "Esiste già una ricetta con questo nome",
		String(database.ERR_RECIPE_NOT_FOUND):           "Ricetta non trovata",
		String(database.ERR_SECTION_DUPLICATED):         "Esiste già una sezione con lo stesso nome",
		String(database.ERR_SECTION_NOT_FOUND):          "Sezione non trovata",
		String(database.ERR_UNKNOWN):                    "Errore sconosciuto",
		String(database.ERR_USER_MAIL_INVALID):          "Email non valida",
		String(database.ERR_USER_MAIL_UNAVAIL):          "Email non disponibile",
		String(database.ERR_USER_NAME_TOO_SHORT):        "Nome utente non valido: lunghezza minima 5 caratteri",
		String(database.ERR_USER_NAME_UNAVAIL):          "Nome utente non disponibile",
		String(database.ERR_USER_PASS_TOO_SHORT):        "Password non valida: lunghezza minima 8 caratteri",
		String(database.ERR_USER_UNKNOWN):               "Utente sconosciuto",
		String(database.ERR_USER_WRONG_CREDENTIALS):     "Credenziali non valide",
		String(database.ERR_USER_WRONG_TOKEN):           "Qualcosa è andato storto",
		STR_ARTICLES_ADDED:                              "Articoli aggiunti correttamente.",
		STR_EMAIL_CHANGED:                               "Email cambiata con successo",
		STR_EMAIL_SENT:                                  "Ti abbiamo inviato un'email: controlla la tua casella di posta",
		STR_LANG_CHANGED:                                "Lingua cambiata correttamente",
		STR_MENU_DELETED:                                "Menù elimnato con successo",
		STR_ORDER_CHANGED:                               "L'ordine degli articoli è cambiato",
		STR_PAGE_NOT_FOUND:                              "Pagina non trovata",
		STR_PASSWORD_CHANGED:                            "Password cambiata con successo",
		STR_RECIPE_DELETED:                              "Ricetta eliminata correttamente",
		STR_SECTION_DELETED:                             "Sezione eliminata con successo",
		STR_SHOPPINGLIST_EMPTIED:                        "Lista svuotata con successo",
		STR_UNKNOWN_LANG:                                "Lingua sconosciuta",
		STR_UNKNOWN_REQUEST:                             "Richiesta sconosciuta",
		STR_UNMATCHING_PASSWORDS:                        "Le due password non corrispondono",
		STR_USERNAME_CHANGED:                            "Nome cambiato con successo",
		STR_USER_CREATED:                                "Account creato con successo",
		STR_USER_DELETED:                                "Account eliminato con successo",
		STR_STORAGE:                                     "Dispensa",
		STR_SUNDAY:                                      "Domenica",
		STR_THURSDAY:                                    "Giovedì",
		STR_TUESDAY:                                     "Martedì",
		STR_TUTORIAL:                                    "Guida",
		STR_USERNAME:                                    "Nome utente",
		STR_WEDNESDAY:                                   "Mercoledì",
		STR_WELCOMEBACK:                                 "Bentornato/a," + placeholder + "!",
		STR_ADD:                                         "Aggiungi",
		STR_ADD_ARTICLES:                                "Aggiungi articoli",
		STR_APPEND_ENTRIES:                              "Aggiungi elementi",
		STR_CANCEL:                                      "Annulla",
		STR_CHANGE_EMAIL:                                "Cambio email",
		STR_CHANGE_PASSWORD:                             "Cambio password",
		STR_CHANGE_USERNAME:                             "Cambio nome utente",
		STR_CLONE:                                       "Clona",
		STR_CONFIRM:                                     "Conferma",
		STR_CURRENT_SEARCH:                              "Ricerca corrente",
		STR_DELETE:                                      "Elimina",
		STR_DELETE_MENU:                                 "Elimina menù",
		STR_DELETE_MENU_TEXT:                            "Sei sicuro di voler cancellare questo menù?",
		STR_DELETE_RECIPE:                               "Elimina ricetta",
		STR_DELETE_RECIPE_TEXT:                          "Sei sicuro di voler eliminare questa ricetta?",
		STR_DELETE_SECTION_TEXT:                         "Sei sicuro? Tutti gli articoli in questa sezione verranno eliminati.",
		STR_DELETE_SELECTED:                             "Elimina selezionati",
		STR_DELETE_USER:                                 "Elimina account",
		STR_DELETE_USER_TEXT1:                           "Sei sicuro di voler eliminare il tuo account?",
		STR_DELETE_USER_TEXT2:                           "Sei DAVVERO sicuro di voler eliminare il tuo account? Questa azione è irreversibile.",
		STR_DIRECTIONS:                                  "Procedimento",
		STR_EDIT:                                        "Modifica",
		STR_EDIT_ARTICLE:                                "Modifica articolo",
		STR_EDIT_ENTRY:                                  "Modifica elemento",
		STR_EDIT_MENU:                                   "Modifica menù",
		STR_EDIT_RECIPE:                                 "Modifica ricetta",
		STR_EDIT_SECTION:                                "Modifica sezione",
		STR_EMAIL:                                       "Email",
		STR_EXPIRATION:                                  "Scadenza",
		STR_FORGOT_PASSWORD:                             "Password dimenticata",
		STR_FRIDAY:                                      "Venerdì",
		STR_GOOD_MORNING:                                "Buongiorno " + placeholder + ",",
		STR_INFO:                                        "Maggiori informazioni",
		STR_INGREDIENTS:                                 "Ingredienti",
		STR_LANGUAGE:                                    "Lingua",
		STR_LOGOUT:                                      "Esci",
		STR_MENUS:                                       "Menù",
		STR_MONDAY:                                      "Lunedì",
		STR_NAME:                                        "Nome",
		STR_NEW_EMAIL:                                   "Nuova email",
		STR_NEW_MENU:                                    "Nuovo menù",
		STR_NEW_PASSWORD:                                "Nuova password",
		STR_NEW_RECIPE:                                  "Nuova ricetta",
		STR_NEW_SECTION:                                 "Nuova sezione",
		STR_NEW_USERNAME:                                "Nuovo nome utente",
		STR_NOREPLY:                                     "Questa email è stata generata automaticamente. Si prega di non rispondere.",
		STR_NOTES:                                       "Note",
		STR_OK:                                          "Va bene",
		STR_OLD_PASSWORD:                                "Vecchia password",
		STR_PASSWORD:                                    "Password",
		STR_QUANTITY:                                    "Quantità",
		STR_RECIPES:                                     "Ricette",
		STR_RECIPES_EMPTY:                               "Nessuna ricetta trovata.",
		STR_REGARDS:                                     "Saluti",
		STR_REPEAT_PASSWORD:                             "Ripeti password",
		STR_RESET_PASSWORD:                              "Reset password",
		STR_SATURDAY:                                    "Sabato",
		STR_SAVE:                                        "Salva",
		STR_SEARCH_ARTICLES:                             "Ricerca articoli",
		STR_SEARCH_EMPTY:                                "Nessun articolo trovato",
		STR_SECTION:                                     "Sezione",
		STR_SECTIONS:                                    "Sezioni della dispensa",
		STR_SECTION_EMPTY:                               "Questa sezione è vuota.",
		STR_SETTINGS:                                    "Impostazioni",
		STR_SET_EMAIL_LANG:                              "Imposta lingua delle email",
		STR_SHOPPINGLIST:                                "Lista della spesa",
		STR_SHOPPINGLIST_EMPTY:                          "La lista è vuota.",
		STR_SIGNIN:                                      "Accedi",
		STR_SIGNUP:                                      "Registrati",
		STR_SIGNUP_DONE:                                 "Registrazione avvenuta",
		STR_STARS:                                       "Stelle",
		STR_STATS:                                       "Statistiche",
		STR_GOODBYE:                                     "Arrivederci",

		STR_STATS_ARTICLES: placeholder + " articoli",
		STR_STATS_SECTIONS: placeholder + " sezioni",
		STR_STATS_ENTRIES:  placeholder + " elementi",
		STR_STATS_MENUS:    placeholder + " menù",
		STR_STATS_RECIPES:  placeholder + " ricette",
		STR_STATS_USERS:    placeholder + " utenti",
		STR_CLICK_HERE:     "clicca qui",
	},
}
