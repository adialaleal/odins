package screens

import (
	"fmt"
	"time"

	"github.com/adialaleal/odins/internal/docker"
	"github.com/adialaleal/odins/internal/state"
	"github.com/adialaleal/odins/internal/tui/components"
	"github.com/adialaleal/odins/internal/tui/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatusTickMsg is sent on each polling tick.
type StatusTickMsg struct{}

// RouteStatus maps a subdomain to its online status.
type RouteStatus map[string]bool

// DashboardModel is the main screen showing all active routes.
type DashboardModel struct {
	table    table.Model
	routes   []state.Route
	statuses RouteStatus
	width    int
	height   int
	confirm  bool // confirm delete?
	selected string
}

// NewDashboard creates the initial dashboard model.
func NewDashboard(routes []state.Route, width, height int) DashboardModel {
	cols := []table.Column{
		{Title: "STATUS", Width: 7},
		{Title: "SUBDOMAIN", Width: 30},
		{Title: "PORT", Width: 6},
		{Title: "PROTO", Width: 6},
		{Title: "RUNTIME", Width: 14},
		{Title: "PROJECT", Width: 16},
	}

	m := DashboardModel{
		routes:   routes,
		statuses: make(RouteStatus),
		width:    width,
		height:   height,
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(m.buildRows()),
		table.WithFocused(true),
		table.WithHeight(height-10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.ColorBorder).
		BorderBottom(true).
		Foreground(styles.ColorAccent).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(styles.ColorAccent).
		Background(styles.ColorSurface).
		Bold(true)
	t.SetStyles(s)

	m.table = t
	return m
}

// Init starts the status polling ticker.
func (m DashboardModel) Init() tea.Cmd {
	return tickStatus()
}

// Update handles messages for the dashboard.
func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case StatusTickMsg:
		newStatuses := make(RouteStatus)
		for _, r := range m.routes {
			newStatuses[r.Subdomain] = docker.CheckSubdomain(r.Port)
		}
		m.statuses = newStatuses
		m.table.SetRows(m.buildRows())
		cmds = append(cmds, tickStatus())

	case tea.KeyMsg:
		if m.confirm {
			switch msg.String() {
			case "y", "Y":
				m.confirm = false
				// Signal deletion — handled by parent
				return m, deleteRoute(m.selected)
			case "n", "N", "esc":
				m.confirm = false
			}
			return m, nil
		}

		switch msg.String() {
		case "d", "delete", "backspace":
			if len(m.routes) > 0 {
				row := m.table.SelectedRow()
				if len(row) >= 2 {
					m.selected = row[1]
					m.confirm = true
				}
			}
		}
	}

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	return m, tea.Batch(cmds...)
}

// View renders the dashboard.
func (m DashboardModel) View() string {
	titleBar := components.TitleBar(m.width,
		"Dashboard — Rotas Ativas",
		fmt.Sprintf("%d rotas", len(m.routes)),
	)

	tableView := lipgloss.NewStyle().
		Padding(0, 2).
		Render(m.table.View())

	hints := []components.KeyHint{
		{Key: "a", Desc: "adicionar"},
		{Key: "u", Desc: "odins up"},
		{Key: "d", Desc: "remover"},
		{Key: "s", Desc: "settings"},
		{Key: "l", Desc: "logs"},
		{Key: "q", Desc: "sair"},
	}
	footer := components.Footer(m.width, hints)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleBar,
		tableView,
	)

	if m.confirm {
		overlay := components.ConfirmModal(m.width, "Remover "+m.selected+"?")
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", overlay)
	}

	return lipgloss.JoinVertical(lipgloss.Left, content, footer)
}

// SetRoutes updates the route list and refreshes the table.
func (m *DashboardModel) SetRoutes(routes []state.Route) {
	m.routes = routes
	m.table.SetRows(m.buildRows())
}

// SetSize updates the terminal dimensions.
func (m *DashboardModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.table.SetHeight(h - 10)
}

func (m DashboardModel) buildRows() []table.Row {
	var rows []table.Row
	for _, r := range m.routes {
		proto := "HTTP"
		if r.HTTPS {
			proto = "HTTPS"
		}
		status := components.StatusDot(m.statuses[r.Subdomain])
		runtime := r.Runtime
		if r.DockerContainer != "" {
			runtime = "docker"
		}
		rows = append(rows, table.Row{
			status,
			r.Subdomain,
			fmt.Sprintf("%d", r.Port),
			proto,
			styles.RuntimeBadge(runtime),
			r.Project,
		})
	}
	return rows
}

// DeleteRouteMsg signals that a route should be deleted.
type DeleteRouteMsg struct{ Subdomain string }

func deleteRoute(subdomain string) tea.Cmd {
	return func() tea.Msg {
		return DeleteRouteMsg{Subdomain: subdomain}
	}
}

func tickStatus() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return StatusTickMsg{}
	})
}
