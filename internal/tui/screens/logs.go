package screens

import (
	"os"
	"strings"
	"time"

	"github.com/adialaleal/odins/internal/i18n"
	"github.com/adialaleal/odins/internal/tui/components"
	"github.com/adialaleal/odins/internal/tui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LogLineMsg carries a new log line from the tail goroutine.
type LogLineMsg struct{ Line string }

// LogsModel tails the proxy access log.
type LogsModel struct {
	viewport viewport.Model
	logPath  string
	lines    []string
	offset   int64
	width    int
	height   int
}

// NewLogs creates the logs screen.
func NewLogs(logPath string, width, height int) LogsModel {
	vp := viewport.New(width-4, height-8)
	vp.Style = lipgloss.NewStyle().
		Background(styles.ColorBg).
		Foreground(styles.ColorText)

	return LogsModel{
		viewport: vp,
		logPath:  logPath,
		width:    width,
		height:   height,
	}
}

// Init starts the log tail goroutine.
func (m LogsModel) Init() tea.Cmd {
	return m.tailCmd()
}

// Update handles new log lines and viewport scroll.
func (m LogsModel) Update(msg tea.Msg) (LogsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case LogLineMsg:
		if msg.Line == "" {
			return m, m.tailCmd()
		}
		colored := colorLogLine(msg.Line)
		m.lines = append(m.lines, colored)
		if len(m.lines) > 500 {
			m.lines = m.lines[len(m.lines)-500:]
		}
		m.viewport.SetContent(strings.Join(m.lines, "\n"))
		m.viewport.GotoBottom()
		return m, m.tailCmd()

	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, func() tea.Msg { return AddRouteCancelMsg{} }
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the logs screen.
func (m LogsModel) View() string {
	titleBar := components.TitleBar(m.width, i18n.T("logs.title"), m.logPath)

	hints := []components.KeyHint{
		{Key: "↑/↓", Desc: i18n.T("hint.scroll")},
		{Key: "Esc", Desc: i18n.T("hint.back")},
	}
	footer := components.Footer(m.width, hints)

	return lipgloss.JoinVertical(lipgloss.Left,
		titleBar,
		lipgloss.NewStyle().Padding(0, 2).Render(m.viewport.View()),
		footer,
	)
}

// SetSize updates viewport dimensions.
func (m *LogsModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w - 4
	m.viewport.Height = h - 8
}

// tailCmd reads new bytes from the log file and emits a LogLineMsg.
func (m LogsModel) tailCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(500 * time.Millisecond)

		f, err := os.Open(m.logPath)
		if err != nil {
			return LogLineMsg{}
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return LogLineMsg{}
		}

		if fi.Size() <= m.offset {
			return LogLineMsg{}
		}

		f.Seek(m.offset, 0)
		buf := make([]byte, fi.Size()-m.offset)
		n, _ := f.Read(buf)
		if n == 0 {
			return LogLineMsg{}
		}

		// Return the last non-empty line
		lines := strings.Split(strings.TrimSpace(string(buf[:n])), "\n")
		for i := len(lines) - 1; i >= 0; i-- {
			if lines[i] != "" {
				return LogLineMsg{Line: lines[i]}
			}
		}
		return LogLineMsg{}
	}
}

func colorLogLine(line string) string {
	switch {
	case strings.Contains(line, " 2"):
		return styles.Log2xx.Render(line)
	case strings.Contains(line, " 4"):
		return styles.Log4xx.Render(line)
	case strings.Contains(line, " 5"):
		return styles.Log5xx.Render(line)
	default:
		return styles.LogDim.Render(line)
	}
}
