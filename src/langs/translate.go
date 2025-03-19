package langs

import (
	"context"
	"strings"
)

const contextKey string = "lang"
const placeholder string = "%%"

// GetCtx returns a context in which is saved the language tag
func (l *Lang) Ctx() context.Context {
	return context.WithValue(context.Background(), contextKey, l)
}

// Translate translates a String
func Translate(ctx context.Context, s String) string {
	lang, ok := ctx.Value(contextKey).(*Lang)
	if !ok {
		lang = Default
	}

	return lang.Strings[s]
}

// TranslateArg translates a String, then replaces the placeholder
// with the given argument
func TranslateArg(ctx context.Context, s String, a string) string {
	return strings.ReplaceAll(Translate(ctx, s), placeholder, a)
}
