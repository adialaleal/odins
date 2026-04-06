package screens

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/adialaleal/odins/internal/detect"
	"github.com/adialaleal/odins/internal/state"
	"github.com/adialaleal/odins/internal/tui/components"
	"github.com/adialaleal/odins/internal/tui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	fieldSubdomain = iota
	fieldPort
	fieldProject
	fieldDocker
	fieldCount
)

// AddRouteSubmitMsg is emitted when the user confirms the form.
type AddRouteSubmitMsg struct {
	Route state.Route
}

// AddRouteCancelMsg is emitted when the user cancels.
type AddRouteCancelMsg struct{}

// AddRouteModel is the form screen for adding a new route.
type AddRouteModel struct {
	inputs     []textinput.Model
	focusIndex int
	width      int
	height     int
	err        string
	detected   *detect.DetectedProject
}

// NewAddRoute creates the add-route form, pre-filling from detection if available.
func NewAddRoute(width, height int, detected *detect.DetectedProject, tld string) AddRouteModel {
	inputs := make([]textinput.Model, fieldCount)

	for i := range inputs {
		t := textinput.New()
		t.CharLimit = 80
		inputs[i] = t
	}

	inputs[fieldSubdomain].Placeholder = "api.rankly." + tld
	inputs[fieldSubdomain].Focus()
	inputs[fieldPort].Placeholder = "3000"
	inputs[fieldProject].Placeholder = "rankly"
	inputs[fieldDocker].Placeholder = "container-name (opcional)"

	m := AddRouteModel{
		inputs:   inputs,
		width:    width,
		height:   height,
		detected: detected,
	}

	// Pre-fill from detection
	if detected != nil {
		if detected.Framework != "" {
			inputs[fieldSubdomain].SetValue(detected.Name + "." + tld)
		}
		inputs[fieldPort].SetValue(fmt.Sprintf("%d", detected.Port))
		inputs[fieldProject].SetValue(detected.Name)
	}

	return m
}

// Init starts cursor blink.
func (m AddRouteModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles keyboard navigation and submission.
func (m AddRouteModel) Update(msg tea.Msg) (AddRouteModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return AddRouteCancelMsg{} }

		case "enter":
			if m.focusIndex == fieldCount-1 || msg.String() == "enter" {
				// Validate and submit
				route, err := m.validate()
				if err != nil {
					m.err = err.Error()
					return m, nil
				}
				m.err = ""
				return m, func() tea.Msg { return AddRouteSubmitMsg{Route: route} }
			}
			m.nextField()

		case "tab", "down":
			m.nextField()

		case "shift+tab", "up":
			m.prevField()
		}
	}

	// Update focused input
	var cmds []tea.Cmd
	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

// View renders the add-route form.
func (m AddRouteModel) View() string {
	titleBar := components.TitleBar(m.width, "Adicionar Rota", "")

	labels := []string{"Subdomínio", "Porta", "Projeto", "Docker"}
	fields := make([]string, fieldCount)

	for i, inp := range m.inputs {
		borderStyle := styles.InputInactive
		if i == m.focusIndex {
			borderStyle = styles.InputActive
		}
		fields[i] = lipgloss.JoinHorizontal(
			lipgloss.Top,
			styles.InputLabel.Render(labels[i]+":"),
			borderStyle.Width(40).Render(inp.View()),
		)
	}

	form := lipgloss.JoinVertical(lipgloss.Left, fields...)

	var detectionInfo string
	if m.detected != nil && m.detected.Runtime != "unknown" {
		detectionInfo = styles.StatusInfo.Render(
			fmt.Sprintf("  ✦ Detectado: %s/%s (porta %d)",
				m.detected.Runtime, m.detected.Framework, m.detected.Port),
		)
	}

	errMsg := ""
	if m.err != "" {
		errMsg = styles.StatusError.Render("  ✗ " + m.err)
	}

	hints := []components.KeyHint{
		{Key: "Tab", Desc: "próximo campo"},
		{Key: "Enter", Desc: "confirmar"},
		{Key: "Esc", Desc: "cancelar"},
	}
	footer := components.Footer(m.width, hints)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleBar,
		"",
		lipgloss.NewStyle().Padding(0, 2).Render(form),
		"",
	)
	if detectionInfo != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, detectionInfo, "")
	}
	if errMsg != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, errMsg, "")
	}

	return lipgloss.JoinVertical(lipgloss.Left, content, footer)
}

func (m *AddRouteModel) nextField() {
	m.inputs[m.focusIndex].Blur()
	m.focusIndex = (m.focusIndex + 1) % fieldCount
	m.inputs[m.focusIndex].Focus()
}

func (m *AddRouteModel) prevField() {
	m.inputs[m.focusIndex].Blur()
	m.focusIndex = (m.focusIndex - 1 + fieldCount) % fieldCount
	m.inputs[m.focusIndex].Focus()
}

func (m AddRouteModel) validate() (state.Route, error) {
	subdomain := strings.TrimSpace(m.inputs[fieldSubdomain].Value())
	portStr := strings.TrimSpace(m.inputs[fieldPort].Value())
	project := strings.TrimSpace(m.inputs[fieldProject].Value())
	docker := strings.TrimSpace(m.inputs[fieldDocker].Value())

	if subdomain == "" {
		return state.Route{}, fmt.Errorf("subdomínio é obrigatório")
	}
	if portStr == "" {
		return state.Route{}, fmt.Errorf("porta é obrigatória")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return state.Route{}, fmt.Errorf("porta inválida")
	}

	return state.Route{
		Subdomain:       subdomain,
		Port:            port,
		Project:         project,
		DockerContainer: docker,
		HTTPS:           true,
	}, nil
}
