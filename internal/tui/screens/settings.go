package screens

import (
	"fmt"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/tui/components"
	"github.com/adialaleal/odins/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SettingsSavedMsg is emitted when the user saves settings.
type SettingsSavedMsg struct {
	Config config.GlobalConfig
}

// SettingsModel is the settings screen.
type SettingsModel struct {
	cfg        config.GlobalConfig
	tldIndex   int
	backendIdx int
	focusRow   int // 0 = TLD row, 1 = backend row
	width      int
	height     int
	saved      bool
}

var backends = []string{"caddy", "nginx", "apache"}

// NewSettings creates the settings screen.
func NewSettings(cfg config.GlobalConfig, width, height int) SettingsModel {
	tldIdx := 0
	for i, t := range config.SupportedTLDs {
		if t.TLD == cfg.TLD {
			tldIdx = i
			break
		}
	}
	backendIdx := 0
	for i, b := range backends {
		if string(cfg.ProxyBackend) == b {
			backendIdx = i
			break
		}
	}
	return SettingsModel{
		cfg:        cfg,
		tldIndex:   tldIdx,
		backendIdx: backendIdx,
		width:      width,
		height:     height,
	}
}

func (m SettingsModel) Init() tea.Cmd { return nil }

func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return AddRouteCancelMsg{} }

		case "tab", "down":
			m.focusRow = (m.focusRow + 1) % 2

		case "shift+tab", "up":
			m.focusRow = (m.focusRow - 1 + 2) % 2

		case "left", "h":
			if m.focusRow == 0 {
				m.tldIndex = (m.tldIndex - 1 + len(config.SupportedTLDs)) % len(config.SupportedTLDs)
			} else {
				m.backendIdx = (m.backendIdx - 1 + len(backends)) % len(backends)
			}

		case "right", "l":
			if m.focusRow == 0 {
				m.tldIndex = (m.tldIndex + 1) % len(config.SupportedTLDs)
			} else {
				m.backendIdx = (m.backendIdx + 1) % len(backends)
			}

		case "enter", "s":
			m.cfg.TLD = config.SupportedTLDs[m.tldIndex].TLD
			m.cfg.ProxyBackend = config.ProxyBackend(backends[m.backendIdx])
			m.saved = true
			return m, func() tea.Msg { return SettingsSavedMsg{Config: m.cfg} }
		}
	}
	return m, nil
}

func (m SettingsModel) View() string {
	titleBar := components.TitleBar(m.width, "Configurações", "")

	// TLD selector
	tldLabel := styles.InputLabel.Render("TLD:")
	tldItems := make([]string, len(config.SupportedTLDs))
	for i, t := range config.SupportedTLDs {
		item := "." + t.TLD
		if i == m.tldIndex {
			item = styles.StatusSuccess.Render("▶ " + item)
		} else {
			item = styles.LogDim.Render("  " + item)
		}
		if t.Warning != "" {
			item += styles.StatusError.Render(" ⚠")
		}
		tldItems[i] = item
	}
	tldRow := lipgloss.JoinHorizontal(lipgloss.Top, tldLabel, "  ")
	for _, item := range tldItems {
		tldRow = lipgloss.JoinHorizontal(lipgloss.Top, tldRow, item+"  ")
	}

	var tldBorder lipgloss.Style
	if m.focusRow == 0 {
		tldBorder = styles.InputActive
	} else {
		tldBorder = styles.InputInactive
	}
	tldSection := tldBorder.Width(m.width - 8).Render(tldRow)

	// Proxy backend selector
	backendLabel := styles.InputLabel.Render("Proxy:")
	backendItems := make([]string, len(backends))
	for i, b := range backends {
		item := b
		if i == m.backendIdx {
			item = styles.StatusSuccess.Render("▶ " + item)
		} else {
			item = styles.LogDim.Render("  " + item)
		}
		backendItems[i] = item
	}
	backendRow := lipgloss.JoinHorizontal(lipgloss.Top, backendLabel, "  ")
	for _, item := range backendItems {
		backendRow = lipgloss.JoinHorizontal(lipgloss.Top, backendRow, item+"  ")
	}

	var backendBorder lipgloss.Style
	if m.focusRow == 1 {
		backendBorder = styles.InputActive
	} else {
		backendBorder = styles.InputInactive
	}
	backendSection := backendBorder.Width(m.width - 8).Render(backendRow)

	// Warning if .local selected
	warnLine := ""
	if config.SupportedTLDs[m.tldIndex].Warning != "" {
		warnLine = styles.StatusError.Render(
			"  ⚠  .local conflita com mDNS/Bonjour no macOS — use com cuidado",
		)
	}

	currentTLD := "." + config.SupportedTLDs[m.tldIndex].TLD
	currentBackend := backends[m.backendIdx]
	info := styles.StatusInfo.Render(
		fmt.Sprintf("  Domínios: *%s → 127.0.0.1  |  Proxy: %s", currentTLD, currentBackend),
	)

	hints := []components.KeyHint{
		{Key: "←/→", Desc: "selecionar"},
		{Key: "Tab", Desc: "navegar"},
		{Key: "Enter", Desc: "salvar"},
		{Key: "Esc", Desc: "voltar"},
	}
	footer := components.Footer(m.width, hints)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleBar,
		"",
		lipgloss.NewStyle().Padding(0, 2).Render(tldSection),
		"",
		lipgloss.NewStyle().Padding(0, 2).Render(backendSection),
		"",
		info,
	)
	if warnLine != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, warnLine)
	}

	return lipgloss.JoinVertical(lipgloss.Left, content, "", footer)
}
