package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#7C3AED") // violet
	ColorAccent    = lipgloss.Color("#A78BFA") // light violet
	ColorSuccess   = lipgloss.Color("#10B981") // emerald
	ColorWarning   = lipgloss.Color("#F59E0B") // amber
	ColorDanger    = lipgloss.Color("#EF4444") // red
	ColorMuted     = lipgloss.Color("#6B7280") // gray
	ColorBg        = lipgloss.Color("#0F0F0F") // near black
	ColorSurface   = lipgloss.Color("#1A1A2E") // dark blue-black
	ColorBorder    = lipgloss.Color("#3730A3") // indigo border
	ColorText      = lipgloss.Color("#E5E7EB") // light gray text
	ColorTextDim   = lipgloss.Color("#9CA3AF") // dimmed text

	// Dot indicators
	DotActive   = lipgloss.NewStyle().Foreground(ColorSuccess).Render("●")
	DotInactive = lipgloss.NewStyle().Foreground(ColorMuted).Render("○")
	DotWarning  = lipgloss.NewStyle().Foreground(ColorWarning).Render("◎")

	// Base styles
	Base = lipgloss.NewStyle().
		Background(ColorBg).
		Foreground(ColorText)

	// App border
	AppBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Background(ColorBg)

	// Header
	Header = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Padding(0, 1)

	HeaderDim = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Padding(0, 1)

	// Title bar
	TitleBar = lipgloss.NewStyle().
		Background(ColorSurface).
		Foreground(ColorAccent).
		Bold(true).
		Padding(0, 2).
		Width(80)

	// Footer / help bar
	Footer = lipgloss.NewStyle().
		Foreground(ColorMuted).
		Padding(0, 1)

	FooterKey = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	// Table
	TableHeader = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ColorBorder).
		Padding(0, 1)

	TableRow = lipgloss.NewStyle().
		Foreground(ColorText).
		Padding(0, 1)

	TableRowSelected = lipgloss.NewStyle().
		Background(ColorSurface).
		Foreground(ColorAccent).
		Bold(true).
		Padding(0, 1)

	// Form
	InputLabel = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Width(14)

	InputActive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(0, 1)

	InputInactive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorMuted).
		Padding(0, 1)

	// Status messages
	StatusSuccess = lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true)

	StatusError = lipgloss.NewStyle().
		Foreground(ColorDanger).
		Bold(true)

	StatusInfo = lipgloss.NewStyle().
		Foreground(ColorAccent)

	// Modal / confirm
	Modal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorWarning).
		Background(ColorSurface).
		Padding(1, 2)

	// Log lines
	Log2xx = lipgloss.NewStyle().Foreground(ColorSuccess)
	Log4xx = lipgloss.NewStyle().Foreground(ColorWarning)
	Log5xx = lipgloss.NewStyle().Foreground(ColorDanger)
	LogDim  = lipgloss.NewStyle().Foreground(ColorMuted)

	// Badge
	BadgeNode   = lipgloss.NewStyle().Foreground(lipgloss.Color("#68D391")).Bold(true)
	BadgeGo     = lipgloss.NewStyle().Foreground(lipgloss.Color("#63B3ED")).Bold(true)
	BadgePython = lipgloss.NewStyle().Foreground(lipgloss.Color("#F6E05E")).Bold(true)
	BadgeDocker = lipgloss.NewStyle().Foreground(lipgloss.Color("#76E4F7")).Bold(true)
)

// RuntimeBadge returns a styled badge for a runtime string.
func RuntimeBadge(runtime string) string {
	switch runtime {
	case "node":
		return BadgeNode.Render("⬡ node")
	case "go":
		return BadgeGo.Render("◉ go")
	case "python":
		return BadgePython.Render("⬡ python")
	case "docker":
		return BadgeDocker.Render("⬡ docker")
	default:
		return LogDim.Render("? " + runtime)
	}
}

// HelpKey renders a key hint for the footer.
func HelpKey(key, desc string) string {
	return FooterKey.Render("["+key+"]") + " " + Footer.Render(desc)
}
