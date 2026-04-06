package components

import (
	"github.com/adialaleal/odins/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmModal renders a confirmation dialog box.
func ConfirmModal(width int, message string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.StatusError.Render("⚠  "+message),
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			styles.FooterKey.Render("[y]")+" confirmar  ",
			styles.FooterKey.Render("[n]")+" cancelar",
		),
	)

	modal := styles.Modal.
		Width(40).
		Align(lipgloss.Center).
		Render(content)

	return lipgloss.Place(width, 10, lipgloss.Center, lipgloss.Center, modal)
}
