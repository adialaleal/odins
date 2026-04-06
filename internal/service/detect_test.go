package service

import (
	"path/filepath"
	"testing"
)

func TestDetectFixtures(t *testing.T) {
	t.Parallel()

	manager := New(DefaultRuntime())
	root := filepath.Join("..", "..", "testdata", "fixtures")

	tests := []struct {
		name      string
		fixture   string
		runtime   string
		framework string
		port      int
		docker    bool
		compose   bool
	}{
		{name: "node vite", fixture: "node-vite", runtime: "node", framework: "vite", port: 4173},
		{name: "go gin", fixture: "go-gin", runtime: "go", framework: "gin", port: 8088},
		{name: "python fastapi", fixture: "python-fastapi", runtime: "python", framework: "fastapi", port: 9000},
		{name: "docker node", fixture: "docker-node", runtime: "node", framework: "express", port: 3000, docker: true, compose: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dir := filepath.Join(root, tt.fixture)
			result, _, err := manager.Detect(dir)
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			if result.Detected.Runtime != tt.runtime {
				t.Fatalf("runtime = %q, want %q", result.Detected.Runtime, tt.runtime)
			}
			if result.Detected.Framework != tt.framework {
				t.Fatalf("framework = %q, want %q", result.Detected.Framework, tt.framework)
			}
			if result.Detected.Port != tt.port {
				t.Fatalf("port = %d, want %d", result.Detected.Port, tt.port)
			}
			if result.Detected.HasDocker != tt.docker {
				t.Fatalf("HasDocker = %t, want %t", result.Detected.HasDocker, tt.docker)
			}
			if result.Detected.HasCompose != tt.compose {
				t.Fatalf("HasCompose = %t, want %t", result.Detected.HasCompose, tt.compose)
			}
			if got := result.RecommendedConfig.Project.Name; got != tt.fixture {
				t.Fatalf("recommended project name = %q, want %q", got, tt.fixture)
			}
			if len(result.RecommendedConfig.Routes) != 1 {
				t.Fatalf("recommended routes len = %d, want 1", len(result.RecommendedConfig.Routes))
			}
		})
	}
}
