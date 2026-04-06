package components

import (
	"github.com/adialaleal/odins/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

const logo = `  ____  ____  ___ _   _ ____
 / __ \|  _ \|_ _| \ | / ___|
| |  | | | | || ||  \| \___ \
| |__| | |_| || || |\  |___) |
 \____/|____/|___|_| \_|____/ `

// Header renders the ODINS app header with logo and subtitle.
func Header(width int, subtitle string) string {
	logoStyle := lipgloss.NewStyle().
		Foreground(styles.ColorPrimary).
		Bold(true)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorMuted).
		Italic(true)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		logoStyle.Render(logo),
		subtitleStyle.Render("  "+subtitle),
		"",
	)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Left).
		Padding(0, 2).
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
