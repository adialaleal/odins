package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
)

// Route represents a single active proxy route.
type Route struct {
	ID              string    `json:"id"`
	Subdomain       string    `json:"subdomain"`
	Port            int       `json:"port"`
	Project         string    `json:"project"`
	Runtime         string    `json:"runtime"`
	DockerContainer string    `json:"docker_container"`
	HTTPS           bool      `json:"https"`
	CreatedAt       time.Time `json:"created_at"`
}

// Store is the persistent route registry.
type Store struct {
	Routes []Route `json:"routes"`
}

func dataDir() string {
	return filepath.Join(xdg.DataHome, "odins")
}

func storePath() string {
	return filepath.Join(dataDir(), "routes.json")
}

// Load reads the route store from disk.
func Load() (*Store, error) {
	path := storePath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Store{Routes: []Route{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Save persists the store to disk.
func (s *Store) Save() error {
	if err := os.MkdirAll(dataDir(), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(storePath(), data, 0644)
}

// Add adds a route to the store (or updates if same subdomain exists).
func (s *Store) Add(r Route) {
	r.ID = "odins-" + r.Subdomain
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	for i, existing := range s.Routes {
		if existing.Subdomain == r.Subdomain {
			s.Routes[i] = r
			return
		}
	}
	s.Routes = append(s.Routes, r)
}

// Remove deletes a route by subdomain.
func (s *Store) Remove(subdomain string) bool {
	for i, r := range s.Routes {
		if r.Subdomain == subdomain {
			s.Routes = append(s.Routes[:i], s.Routes[i+1:]...)
			return true
		}
	}
	return false
}

// Get finds a route by subdomain.
func (s *Store) Get(subdomain string) (Route, bool) {
	for _, r := range s.Routes {
		if r.Subdomain == subdomain {
			return r, true
		}
	}
	return Route{}, false
}

// ByProject returns all routes belonging to a project.
func (s *Store) ByProject(project string) []Route {
	var out []Route
	for _, r := range s.Routes {
		if r.Project == project {
			out = append(out, r)
		}
	}
	return out
}
