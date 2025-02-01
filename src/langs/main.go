package langs

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log/slog"
	"os"
)

// Supported is the list of supported languages
var Supported []language.Tag = []language.Tag{language.English, language.Italian}

// localizers is a collection of localizers
var localizers map[string]*i18n.Localizer

// Load initializes the undle and loads the translations
func Load() {
	// Creates the bundle
	bundle := i18n.NewBundle(Supported[0])
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Initializes the localizers map
	localizers = make(map[string]*i18n.Localizer)

	for _, l := range Supported {
		lang := l.String()

		// Loads the files and the localizers
		if _, err := bundle.LoadMessageFile("langs/active." + lang + ".toml"); err != nil {
			slog.Error("cannot load language:", "lang", lang, "err", err)
			os.Exit(1)
		}

		// Creates the localizers
		localizers[lang] = i18n.NewLocalizer(bundle, lang)
	}
}

// Translate returns the string with that id
// in the given language
func Translate(lang string, id string, data any) string {
	// Gets the required localizer (or the default one)
	l, found := localizers[lang]
	if !found {
		l = localizers[Supported[0].String()]
	}

	// Gets the string
	str, err := l.Localize(&i18n.LocalizeConfig{MessageID: id, TemplateData: data})
	if err != nil {
		slog.Error("while translating:", "id", id, "lang", lang, "err", err)
	}

	return str
}
