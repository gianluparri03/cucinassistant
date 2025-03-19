package langs

import (
	"log/slog"
	"os"

	"cucinassistant/database"
)

// CheckAll checks all the languages
func CheckAll() {
	missing := 0

	for _, l := range Available {
		missing += l.check()
	}

	if missing > 0 {
		os.Exit(1)
	}
}

// check makes sure that a language has all the translations.
// The returned value is the number of missing ones.
func (l *Lang) check() int {
	var missing []string
	ctx := l.Ctx()

	// Checks the database errors
	for i := 1; i < database.ErrorsNumber; i++ {
		if Translate(ctx, ParseError(database.Error(i))) == "" {
			missing = append(missing, database.Error(i).String())
		}
	}

	// Checks all the Strings
	for i := str_begin + 1; i < str_end; i++ {
		if Translate(ctx, String(i)) == "" {
			missing = append(missing, String(i).String())
		}
	}

	// Prints the missing ones, if any
	if len(missing) > 0 {
		slog.Warn("Missing translations", "lang", l.Tag, "missing", missing)
	}

	return len(missing)
}
