package langs

import (
	"cucinassistant/database"
)

// Available is a collection of all the available languages
var Available map[string]*Lang = map[string]*Lang{
	english.Tag: english,
	italian.Tag: italian,
}

// Default is the default language
var Default *Lang = english

// Lang is an available language
type Lang struct {
	// Tag is the language tag (es. "en")
	Tag string

	// Name is the language name, in that language
	Name string

	// Strings contains the translations for the Strings
	Strings map[String]string
}

// Get returns the language with that tag.
// If the language is not found, the default one is returned.
func Get(tag string) *Lang {
	l, ok := Available[tag]
	if ok {
		return l
	} else {
		return Default
	}
}

// ParseError returns a String from a database.Error
func ParseError(err error) String {
	return String(err.(database.Error))
}
