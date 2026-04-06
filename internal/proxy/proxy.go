package proxy

import "github.com/adialaleal/odins/internal/state"

// Backend is the interface all reverse proxy implementations must satisfy.
type Backend interface {
	Name() string
	IsInstalled() bool
	IsRunning() bool
	Install() error
	Start() error
	Stop() error
	Restart() error
	AddRoute(r state.Route) error
	RemoveRoute(subdomain string) error
	Reload() error
	LogPath() string
}
