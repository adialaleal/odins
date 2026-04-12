package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/adialaleal/odins/internal/config"
)

// IndexedProject is a project entry in the global index.
type IndexedProject struct {
	Path      string             `json:"path"`
	Name      string             `json:"name"`
	Domain    string             `json:"domain,omitempty"`
	Runtime   string             `json:"runtime"`
	Framework string             `json:"framework"`
	Port      int                `json:"port"`
	Routes    []config.RouteConfig `json:"routes,omitempty"`
}

// ProjectIndex is the global project registry.
type ProjectIndex struct {
	Projects  []IndexedProject `json:"projects"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func indexPath() string {
	return filepath.Join(config.ConfigDir(), "projects.json")
}

// LoadIndex reads the global project index from disk.
func LoadIndex() (ProjectIndex, error) {
	data, err := os.ReadFile(indexPath())
	if err != nil {
		if os.IsNotExist(err) {
			return ProjectIndex{}, nil
		}
		return ProjectIndex{}, err
	}

	var idx ProjectIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return ProjectIndex{}, err
	}
	return idx, nil
}

// SaveIndex writes the global project index to disk.
func SaveIndex(idx ProjectIndex) error {
	idx.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(indexPath())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(indexPath(), data, 0o644)
}

// UpdateIndex merges scan results into the existing index.
func UpdateIndex(result ScanResult) error {
	idx, _ := LoadIndex() // ignore error, start fresh if corrupt

	// Build lookup by path for existing entries.
	existing := make(map[string]int)
	for i, p := range idx.Projects {
		existing[p.Path] = i
	}

	for _, sp := range result.Projects {
		entry := IndexedProject{
			Path:      sp.Path,
			Name:      sp.Name,
			Domain:    sp.Domain,
			Runtime:   sp.Runtime,
			Framework: sp.Framework,
			Port:      sp.Port,
			Routes:    sp.Routes,
		}

		if i, ok := existing[sp.Path]; ok {
			idx.Projects[i] = entry
		} else {
			idx.Projects = append(idx.Projects, entry)
		}
	}

	return SaveIndex(idx)
}
