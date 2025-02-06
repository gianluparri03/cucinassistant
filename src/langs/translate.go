package langs

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"html/template"
	"io"
	"log/slog"
	"path"
)

// Translate returns the string with that id in the given language
func Translate(lang string, id string, data any) string {
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

// ExecuteTemplates parses the given templates and applies to them
// the translations (function t) in the given language, then
// writes them
func ExecuteTemplates(w io.Writer, lang string, templates []string, data any) {
	// Prepares the FuncMap
	funcs := template.FuncMap{
		"t": func(id string, data any) template.HTML {
			return template.HTML(Translate(lang, id, data))
		},
	}

	// Loads the templates
	tmpl, err := template.New(path.Base(templates[0])).Funcs(funcs).ParseFiles(templates...)
	if err != nil {
		slog.Error("while fetching templates:", "err", err, "templates", templates)
	}

	// Executes it
	if err = tmpl.Execute(w, data); err != nil {
		slog.Error("while executing templates:", "err", err, "templates", templates)
	}
}
