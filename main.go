package main

import (
	"os"
	"strings"

	"github.com/adialaleal/odins/cmd"
	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/i18n"
)

// Injected by GoReleaser via ldflags (-X main.version=...).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	initLang()
	cmd.SetVersion(version, commit, date)
	cmd.Execute()
}

// initLang sets the active language.
// Priority: config file > LANG/LC_ALL/LC_MESSAGES env vars > pt (default).
func initLang() {
	cfg, err := config.LoadGlobal()
	if err == nil && cfg.Language != "" {
		switch i18n.Lang(cfg.Language) {
		case i18n.EN:
			i18n.SetLang(i18n.EN)
			return
		case i18n.ES:
			i18n.SetLang(i18n.ES)
			return
		case i18n.PT:
			i18n.SetLang(i18n.PT)
			return
		}
	}

	for _, env := range []string{"LANG", "LC_ALL", "LC_MESSAGES"} {
		v := os.Getenv(env)
		if strings.HasPrefix(v, "pt") {
			i18n.SetLang(i18n.PT)
			return
		}
		if strings.HasPrefix(v, "es") {
			i18n.SetLang(i18n.ES)
			return
		}
		if strings.HasPrefix(v, "en") {
			i18n.SetLang(i18n.EN)
			return
		}
	}

	// Default: Portuguese
	i18n.SetLang(i18n.PT)
}
