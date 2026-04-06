package components

import (
	"strings"

	"github.com/adialaleal/odins/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// KeyHint represents a single key binding hint.
type KeyHint struct {
	Key  string
	Desc string
}

// Footer renders a help bar with key hints.
func Footer(width int, hints []KeyHint) string {
	var parts []string
	for _, h := range hints {
		parts = append(parts, styles.HelpKey(h.Key, h.Desc))
	}

	content := strings.Join(parts, "  ")

	return lipgloss.NewStyle().
		Background(styles.ColorSurface).
		Foreground(styles.ColorMuted).
		Width(width).
		Padding(0, 1).
		Render(content)
}
