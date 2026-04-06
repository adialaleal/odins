package screens

import (
	"time"

	"github.com/adialaleal/odins/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SplashDoneMsg is sent when the splash timer expires.
type SplashDoneMsg struct{}

// SplashModel renders a full-screen splash on startup.
type SplashModel struct {
	width  int
	height int
	tick   int // drives the dot animation
}

var splashLogoLines = []string{
	`   ██████╗ ██████╗ ██╗███╗   ██╗███████╗`,
	`  ██╔═══██╗██╔══██╗██║████╗  ██║██╔════╝`,
	`  ██║   ██║██║  ██║██║██╔██╗ ██║███████╗`,
	`  ██║   ██║██║  ██║██║██║╚██╗██║╚════██║`,
	`  ╚██████╔╝██████╔╝██║██║ ╚████║███████║`,
	`   ╚═════╝ ╚═════╝ ╚═╝╚═╝  ╚═══╝╚══════╝`,
}

var splashGradient = []lipgloss.Color{
	"#E9D5FF",
	"#C4B5FD",
	"#A78BFA",
	"#8B5CF6",
	"#7C3AED",
	"#6D28D9",
}

// splashTickMsg drives the dot animation.
type splashTickMsg struct{}

// NewSplash creates the splash screen.
func NewSplash(width, height int) SplashModel {
	return SplashModel{width: width, height: height}
}

// Init schedules the dismiss timer and the first animation tick.
func (m SplashModel) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(1800*time.Millisecond, func(t time.Time) tea.Msg { return SplashDoneMsg{} }),
		tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg { return splashTickMsg{} }),
	)
}

// Update handles size changes and the animation tick.
func (m SplashModel) Update(msg tea.Msg) (SplashModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case splashTickMsg:
		m.tick++
		return m, tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg { return splashTickMsg{} })
	}
	return m, nil
}

// View renders the full-screen splash.
func (m SplashModel) View() string {
	// Build gradient logo
	var logoLines []string
	for i, line := range splashLogoLines {
		color := splashGradient[i]
		rendered := lipgloss.NewStyle().
			Foreground(color).
			Bold(true).
			Render(line)
		logoLines = append(logoLines, rendered)
	}

	logo := lipgloss.JoinVertical(lipgloss.Center, logoLines...)

	runeDecor := lipgloss.NewStyle().
		Foreground(styles.ColorBorder).
		Render("ᚦ ᚢ ᚱ ᛋ ᛏ ᚨ ᛉ ᚾ")

	tagline := lipgloss.NewStyle().
		Foreground(styles.ColorAccent).
		Italic(true).
		Render("The All-Father of Local DNS")

	// Animated dots
	dots := [4]string{"   ", ".  ", ".. ", "..."}
	dot := dots[m.tick%4]
	loading := lipgloss.NewStyle().
		Foreground(styles.ColorMuted).
		Render("loading" + dot)

	block := lipgloss.JoinVertical(lipgloss.Center,
		logo,
		"",
		runeDecor,
		tagline,
		"",
		loading,
	)

	// Center the whole block on screen
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(styles.ColorBg).
		Align(lipgloss.Center, lipgloss.Center).
		Render(block)
}
