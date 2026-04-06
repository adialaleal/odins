package components

import "github.com/adialaleal/odins/internal/tui/styles"

// StatusDot returns a colored dot based on the active state.
func StatusDot(active bool) string {
	if active {
		return styles.DotActive
	}
	return styles.DotInactive
}
