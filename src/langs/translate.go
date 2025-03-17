package langs

import (
	"context"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"log/slog"
)

const contextKey string = "lang"

// Lang returns a context in which is saved the lang
func Lang(lang string) context.Context {
	return context.WithValue(context.Background(), contextKey, lang)
}

// Translate returns the string with that id in the language specified in the
// context
func Translate(ctx context.Context, id string, data any) string {
	lang, ok := ctx.Value(contextKey).(string)
	if !ok {
		lang = Default
	}

	// Gets the required localizer (or the default one)
	l, found := localizers[lang]
	if !found {
		l = localizers[Default]
	}

	// Gets the string
	str, err := l.Localize(&i18n.LocalizeConfig{MessageID: id, TemplateData: data})
	if err != nil {
		slog.Error("while translating:", "id", id, "lang", lang, "err", err)
	}

	return str
}
