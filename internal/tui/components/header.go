package components

import (
	"strings"

	"github.com/adialaleal/odins/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// logoLines are rendered with an individual gradient color per line.
var logoLines = []string{
	`   ██████╗ ██████╗ ██╗███╗   ██╗███████╗`,
	`  ██╔═══██╗██╔══██╗██║████╗  ██║██╔════╝`,
	`  ██║   ██║██║  ██║██║██╔██╗ ██║███████╗`,
	`  ██║   ██║██║  ██║██║██║╚██╗██║╚════██║`,
	`  ╚██████╔╝██████╔╝██║██║ ╚████║███████║`,
	`   ╚═════╝ ╚═════╝ ╚═╝╚═╝  ╚═══╝╚══════╝`,
}

// Gradient: light violet at top → deep violet at bottom
var logoGradient = []lipgloss.Color{
	"#E9D5FF", // violet-200
	"#C4B5FD", // violet-300
	"#A78BFA", // violet-400
	"#8B5CF6", // violet-500
	"#7C3AED", // violet-600
	"#6D28D9", // violet-700
}

// Header renders the ODINS app header with gradient logo and subtitle.
func Header(width int, subtitle string) string {
	// Top margin — keeps logo from sticking to terminal edge
	topMargin := "\n"

	var logoRendered []string
	for i, line := range logoLines {
		color := logoGradient[i]
		rendered := lipgloss.NewStyle().
			Foreground(color).
			Bold(true).
			Render(line)
		logoRendered = append(logoRendered, rendered)
	}

	decorLine := lipgloss.NewStyle().
		Foreground(styles.ColorBorder).
		Render("  ᚦ ᚢ ᚱ ᛋ ᛏ ᚨ ᛉ ᚾ")

	subtitleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Italic(true)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		topMargin,
		strings.Join(logoRendered, "\n"),
		decorLine,
		subtitleStyle.Render("    "+subtitle),
		"",
	)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Left).
		Padding(0, 1).
		Render(content)
}

// TitleBar renders a compact title bar for screens.
func TitleBar(width int, title, badge string) string {
	left := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Bold(true).
		Render("  " + title)

	right := lipgloss.NewStyle().
		Foreground(styles.ColorMuted).
		Render(badge + "  ")

	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	spacer := lipgloss.NewStyle().Width(gap).Render("")

	bar := lipgloss.JoinHorizontal(lipgloss.Top, left, spacer, right)

	return lipgloss.NewStyle().
		Background(styles.ColorSurface).
		Width(width).
		Render(bar)
}
