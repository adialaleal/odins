// Package tui provides the Bubble Tea TUI application for ODINS.
package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/detect"
	"github.com/adialaleal/odins/internal/proxy/apache"
	"github.com/adialaleal/odins/internal/proxy/caddy"
	"github.com/adialaleal/odins/internal/proxy/nginx"
	"github.com/adialaleal/odins/internal/state"
	"github.com/adialaleal/odins/internal/tui/components"
	"github.com/adialaleal/odins/internal/tui/screens"
	"github.com/adialaleal/odins/internal/tui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Screen identifies the active TUI screen.
type Screen int

const (
	ScreenSplash Screen = iota
	ScreenDashboard
	ScreenAddRoute
	ScreenSettings
	ScreenLogs
)

// AppModel is the root Bubble Tea model.
type AppModel struct {
	screen    Screen
	width     int
	height    int

	// transition animation
	transOffset int
	transActive bool

	// sub-models
	splash    screens.SplashModel
	dashboard screens.DashboardModel
	addRoute  screens.AddRouteModel
	settings  screens.SettingsModel
	logs      screens.LogsModel

	// app state
	cfg    config.GlobalConfig
	store  *state.Store
	status string
}

// TransitionTickMsg drives the slide-up animation.
type TransitionTickMsg struct{}

// Run starts the TUI application.
func Run() error {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	store, err := state.Load()
	if err != nil {
		return fmt.Errorf("load state: %w", err)
	}

	logPath := resolveLogPath(cfg)

	// Auto-detect project if .odins is absent in the current directory.
	// Set a welcome status message that prompts the user to press [u].
	cwd, _ := os.Getwd()
	autoStatus := ""
	if !config.ExistsProject(cwd) {
		d := detect.Project(cwd)
		if d.Runtime != "unknown" {
			autoStatus = fmt.Sprintf(
				"Projeto '%s' (%s/%s, porta %d) detectado — pressione [u] para ativar as rotas",
				d.Name, d.Runtime, d.Framework, d.Port,
			)
		} else {
			autoStatus = "Nenhum .odins encontrado — pressione [a] para adicionar uma rota manualmente"
		}
	}

	m := AppModel{
		screen: ScreenSplash,
		cfg:    cfg,
		store:  store,
		status: autoStatus,
	}

	// Initialize sub-models with placeholder size (updated on WindowSizeMsg)
	m.splash = screens.NewSplash(80, 24)
	m.dashboard = screens.NewDashboard(store.Routes, 80, 24)
	m.settings = screens.NewSettings(cfg, 80, 24)
	m.logs = screens.NewLogs(logPath, 80, 24)

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err = p.Run()
	return err
}

// Init starts all sub-model inits.
func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.splash.Init(),
		m.dashboard.Init(),
		m.logs.Init(),
	)
}

// Update handles all messages for the root model.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		var splashCmd tea.Cmd
		m.splash, splashCmd = m.splash.Update(msg)
		cmds = append(cmds, splashCmd)
		m.dashboard.SetSize(m.width, m.height)
		m.logs.SetSize(m.width, m.height)
		return m, tea.Batch(cmds...)

	case screens.SplashDoneMsg:
		m.screen = ScreenDashboard
		return m, nil

	case tea.KeyMsg:
		// Dismiss splash immediately on any key press
		if m.screen == ScreenSplash {
			m.screen = ScreenDashboard
			return m, nil
		}
		// Global key handlers
		switch msg.String() {
		case "ctrl+c", "q":
			if m.screen == ScreenDashboard {
				return m, tea.Quit
			}
		case "u":
			if m.screen == ScreenDashboard {
				m.status = "Ativando rotas..."
				return m, odinsUpCmd(m.cfg, m.store)
			}
		case "a":
			if m.screen == ScreenDashboard {
				return m.navigateTo(ScreenAddRoute)
			}
		case "s":
			if m.screen == ScreenDashboard {
				return m.navigateTo(ScreenSettings)
			}
		case "l":
			if m.screen == ScreenDashboard {
				return m.navigateTo(ScreenLogs)
			}
		case "esc":
			if m.screen != ScreenDashboard {
				return m.navigateTo(ScreenDashboard)
			}
		}

	case TransitionTickMsg:
		if m.transActive {
			m.transOffset -= 4
			if m.transOffset <= 0 {
				m.transOffset = 0
				m.transActive = false
			} else {
				cmds = append(cmds, transitionTick())
			}
		}

	case screens.DeleteRouteMsg:
		return m.handleDeleteRoute(msg.Subdomain)

	case screens.AddRouteSubmitMsg:
		return m.handleAddRoute(msg.Route)

	case screens.AddRouteCancelMsg:
		return m.navigateTo(ScreenDashboard)

	case UpDoneMsg:
		if msg.Err != nil {
			m.status = "Erro: " + msg.Err.Error()
		} else if msg.Applied == 0 {
			m.status = "Nenhuma rota aplicada"
		} else {
			m.status = fmt.Sprintf("✓ %d rota(s) ativada(s)!", msg.Applied)
		}
		m.dashboard.SetRoutes(m.store.Routes)
		return m, nil

	case screens.SettingsSavedMsg:
		m.cfg = msg.Config
		if err := config.SaveGlobal(m.cfg); err != nil {
			m.status = "Erro ao salvar config: " + err.Error()
		} else {
			m.status = "Configurações salvas!"
		}
		return m.navigateTo(ScreenDashboard)
	}

	// Delegate to active screen
	switch m.screen {
	case ScreenSplash:
		var cmd tea.Cmd
		m.splash, cmd = m.splash.Update(msg)
		cmds = append(cmds, cmd)
	case ScreenDashboard:
		var cmd tea.Cmd
		m.dashboard, cmd = m.dashboard.Update(msg)
		cmds = append(cmds, cmd)
	case ScreenAddRoute:
		var cmd tea.Cmd
		m.addRoute, cmd = m.addRoute.Update(msg)
		cmds = append(cmds, cmd)
	case ScreenSettings:
		var cmd tea.Cmd
		m.settings, cmd = m.settings.Update(msg)
		cmds = append(cmds, cmd)
	case ScreenLogs:
		var cmd tea.Cmd
		m.logs, cmd = m.logs.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the current screen with optional slide-up transition.
func (m AppModel) View() string {
	// Splash renders full-screen; skip header/status/footer wrapping.
	if m.screen == ScreenSplash {
		return m.splash.View()
	}

	var screenView string

	switch m.screen {
	case ScreenDashboard:
		screenView = m.dashboard.View()
	case ScreenAddRoute:
		screenView = m.addRoute.View()
	case ScreenSettings:
		screenView = m.settings.View()
	case ScreenLogs:
		screenView = m.logs.View()
	}

	// Apply transition offset (slide-up)
	if m.transActive && m.transOffset > 0 {
		placeholder := lipgloss.NewStyle().
			Height(m.transOffset).
			Width(m.width).
			Background(styles.ColorBg).
			Render("")
		screenView = lipgloss.JoinVertical(lipgloss.Left, placeholder, screenView)
	}

	// Status message bar — uses accent color for info/detect messages, green for confirmations
	if m.status != "" {
		prefix := "  ✦ "
		fg := styles.ColorAccent
		// Switch to green checkmark for success-style messages
		if strings.HasPrefix(m.status, "✓") || strings.HasSuffix(m.status, "!") {
			prefix = "  ✓ "
			fg = styles.ColorSuccess
		}
		statusBar := lipgloss.NewStyle().
			Background(styles.ColorSurface).
			Foreground(fg).
			Width(m.width).
			Padding(0, 1).
			Render(prefix + m.status)
		screenView = lipgloss.JoinVertical(lipgloss.Left,
			components.Header(m.width, "The All-Father of Local DNS"),
			screenView,
			statusBar,
		)
	} else {
		screenView = lipgloss.JoinVertical(lipgloss.Left,
			components.Header(m.width, "The All-Father of Local DNS"),
			screenView,
		)
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(styles.ColorBg).
		Render(screenView)
}

func (m AppModel) navigateTo(s Screen) (AppModel, tea.Cmd) {
	m.transOffset = 8
	m.transActive = true

	switch s {
	case ScreenAddRoute:
		cwd, _ := os.Getwd()
		d := detect.Project(cwd)
		m.addRoute = screens.NewAddRoute(m.width, m.height, &d, m.cfg.TLD)
	case ScreenSettings:
		m.settings = screens.NewSettings(m.cfg, m.width, m.height)
	}

	m.screen = s
	m.status = ""
	return m, transitionTick()
}

func (m AppModel) handleAddRoute(r state.Route) (AppModel, tea.Cmd) {
	if r.Runtime == "" {
		cwd, _ := os.Getwd()
		d := detect.Project(cwd)
		r.Runtime = d.Runtime
	}

	m.store.Add(r)
	if err := m.store.Save(); err != nil {
		m.status = "Erro ao salvar: " + err.Error()
	} else {
		if err := addToProxy(m.cfg, r); err != nil {
			m.status = "Erro no proxy: " + err.Error()
		} else {
			m.status = fmt.Sprintf("✓ %s → :%d adicionado!", r.Subdomain, r.Port)
		}
	}

	m.dashboard.SetRoutes(m.store.Routes)
	return m.navigateTo(ScreenDashboard)
}

func (m AppModel) handleDeleteRoute(subdomain string) (AppModel, tea.Cmd) {
	m.store.Remove(subdomain)
	if err := m.store.Save(); err != nil {
		m.status = "Erro ao salvar: " + err.Error()
	} else {
		if err := removeFromProxy(m.cfg, subdomain); err != nil {
			m.status = "Proxy error: " + err.Error()
		} else {
			m.status = subdomain + " removido."
		}
	}
	m.dashboard.SetRoutes(m.store.Routes)
	return m, nil
}

// UpDoneMsg carries the result of running odins up from the TUI.
type UpDoneMsg struct {
	Applied int
	Err     error
}

// odinsUpCmd runs the equivalent of `odins up` as a background tea.Cmd.
func odinsUpCmd(cfg config.GlobalConfig, store *state.Store) tea.Cmd {
	return func() tea.Msg {
		cwd, err := os.Getwd()
		if err != nil {
			return UpDoneMsg{Err: err}
		}

		var routes []config.RouteConfig
		var projName, projRuntime string

		projectCfgPath := filepath.Join(cwd, config.ProjectConfigFile)

		if config.ExistsProject(cwd) {
			projCfg, err := config.LoadProject(projectCfgPath)
			if err != nil {
				return UpDoneMsg{Err: fmt.Errorf("ler .odins: %w", err)}
			}
			routes = projCfg.Routes
			projName = projCfg.Project.Name
			projRuntime = projCfg.Project.Runtime
		} else {
			d := detect.Project(cwd)
			if d.Runtime == "unknown" {
				return UpDoneMsg{Err: fmt.Errorf("projeto não detectado em %s", cwd)}
			}
			projName = d.Name
			projRuntime = d.Runtime
			routes = []config.RouteConfig{{Subdomain: d.Name, Port: d.Port, HTTPS: true}}

			projCfg := config.ProjectConfig{
				Project: config.ProjectInfo{Name: d.Name, Runtime: d.Runtime, Framework: d.Framework},
				Routes:  routes,
			}
			_ = config.SaveProject(projectCfgPath, projCfg)
		}

		applied := 0
		for _, rc := range routes {
			fqdn := tuiUpFQDN(rc.Subdomain, projName, cfg.TLD)
			r := state.Route{
				ID:              "odins-" + fqdn,
				Subdomain:       fqdn,
				Port:            rc.Port,
				Project:         projName,
				Runtime:         projRuntime,
				DockerContainer: rc.DockerContainer,
				HTTPS:           rc.HTTPS,
				CreatedAt:       time.Now(),
			}
			if err := addToProxy(cfg, r); err != nil {
				continue
			}
			store.Add(r)
			applied++
		}
		_ = store.Save()
		return UpDoneMsg{Applied: applied}
	}
}

func tuiUpFQDN(subdomain, project, tld string) string {
	for _, c := range subdomain {
		if c == '.' {
			return subdomain + "." + tld
		}
	}
	return subdomain + "." + project + "." + tld
}

func transitionTick() tea.Cmd {
	return func() tea.Msg { return TransitionTickMsg{} }
}

func resolveLogPath(cfg config.GlobalConfig) string {
	switch cfg.ProxyBackend {
	case config.BackendNginx:
		return nginx.New().LogPath()
	case config.BackendApache:
		return apache.New().LogPath()
	default:
		return caddy.New().LogPath()
	}
}

func addToProxy(cfg config.GlobalConfig, r state.Route) error {
	switch cfg.ProxyBackend {
	case config.BackendNginx:
		return nginx.New().AddRoute(r)
	case config.BackendApache:
		return apache.New().AddRoute(r)
	default:
		return caddy.New().AddRoute(r)
	}
}

func removeFromProxy(cfg config.GlobalConfig, subdomain string) error {
	switch cfg.ProxyBackend {
	case config.BackendNginx:
		return nginx.New().RemoveRoute(subdomain)
	case config.BackendApache:
		return apache.New().RemoveRoute(subdomain)
	default:
		return caddy.New().RemoveRoute(subdomain)
	}
}
