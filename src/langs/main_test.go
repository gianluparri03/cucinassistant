package langs

import (
	"testing"

	"cucinassistant/database"
)

// TestAll runs the test on every language
func TestAll(t *testing.T) {
	for _, l := range Available {
		l.test(t)
	}
}

// test checks wether the language contains all the translations
func (l *Lang) test(t *testing.T) {
	ctx := l.Ctx()

	// Checks the database errors
	for i := 1; i < database.ErrorsNumber; i++ {
		if e := database.Error(i); Translate(ctx, ParseError(e)) == "" {
			t.Errorf("Missing <%s> in language <%s>", e.String(), l.Tag)
		}
	}

	// Checks all the Strings
	for i := str_begin + 1; i < str_end; i++ {
		if s := String(i); Translate(ctx, s) == "" {
			t.Errorf("Missing <%s> in language <%s>", s.String(), l.Tag)
		}
	}
}
