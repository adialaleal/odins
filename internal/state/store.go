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
	Domain          string    `json:"domain,omitempty"` // parent domain (e.g. "tatoh")
	DockerContainer string    `json:"docker_container"`
	HTTPS           bool      `json:"https"`
	CreatedAt       time.Time `json:"created_at"`
}

// Domain represents a local domain workspace (e.g. "tatoh" → tatoh.odins).
type Domain struct {
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Store is the persistent route and domain registry.
type Store struct {
	Routes  []Route  `json:"routes"`
	Domains []Domain `json:"domains"`
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

// ByDomain returns all routes attached to a domain name.
func (s *Store) ByDomain(domain string) []Route {
	var out []Route
	for _, r := range s.Routes {
		if r.Domain == domain {
			out = append(out, r)
		}
	}
	return out
}

// AddDomain adds or updates a domain in the store.
func (s *Store) AddDomain(d Domain) {
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now()
	}
	for i, existing := range s.Domains {
		if existing.Name == d.Name {
			s.Domains[i] = d
			return
		}
	}
	s.Domains = append(s.Domains, d)
}

// GetDomain returns a domain by name.
func (s *Store) GetDomain(name string) (Domain, bool) {
	for _, d := range s.Domains {
		if d.Name == name {
			return d, true
		}
	}
	return Domain{}, false
}

// RemoveDomain deletes a domain and all its routes.
func (s *Store) RemoveDomain(name string) bool {
	for i, d := range s.Domains {
		if d.Name == name {
			s.Domains = append(s.Domains[:i], s.Domains[i+1:]...)
			// Also remove all routes that belong to this domain
			var kept []Route
			for _, r := range s.Routes {
				if r.Domain != name {
					kept = append(kept, r)
				}
			}
			s.Routes = kept
			return true
		}
	}
	return false
}
