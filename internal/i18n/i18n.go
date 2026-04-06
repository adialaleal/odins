// Package i18n provides a simple translation system for ODINS.
// Supported languages: pt (Portuguese), en (English), es (Spanish).
// Language priority: config file > LANG/LC_ALL env vars > pt (default).
package i18n

import "fmt"

// Lang represents a supported language.
type Lang string

const (
	PT Lang = "pt"
	EN Lang = "en"
	ES Lang = "es"
)

var current = PT

var catalogs = map[Lang]map[string]string{
	PT: ptStrings,
	EN: enStrings,
	ES: esStrings,
}

// SetLang sets the active language for all T/Tf calls.
func SetLang(l Lang) {
	current = l
}

// Current returns the active language.
func Current() Lang {
	return current
}

// T returns the translation for key in the active language.
// Falls back to PT, then returns the key itself if not found.
func T(key string) string {
	if m, ok := catalogs[current]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	if v, ok := catalogs[PT][key]; ok {
		return v
	}
	return key
}

// Tf returns the translation for key formatted with args (fmt.Sprintf).
func Tf(key string, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}
