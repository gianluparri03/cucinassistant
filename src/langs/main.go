package langs

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log/slog"
	"os"
)

// Available contains all the supported languages
// in the form tag: Name (the name is in that language's language)
var Available map[string]string = map[string]string{
	"en": "English",
	"it": "Italiano",
}

// Default is the default language
var Default string = "en"

// localizers is a collection of localizers
var localizers map[string]*i18n.Localizer

// Load initializes the undle and loads the translations
func Load() {
	// Creates the bundle
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Initializes the localizers map
	localizers = make(map[string]*i18n.Localizer)

	for lang, _ := range Available {
		// Loads the files and the localizers
		if _, err := bundle.LoadMessageFile("langs/active." + lang + ".toml"); err != nil {
			slog.Error("cannot load language:", "lang", lang, "err", err)
			os.Exit(1)
		}

		// Creates the localizers
		localizers[lang] = i18n.NewLocalizer(bundle, lang)
	}
}
