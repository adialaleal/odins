package service

import (
	"github.com/adialaleal/odins/internal/config"
	"github.com/adialaleal/odins/internal/detect"
)

// DetectResult is the machine-readable inspection output for a project directory.
type DetectResult struct {
	Directory           string                 `json:"directory"`
	ProjectConfigPath   string                 `json:"project_config_path"`
	ProjectConfigExists bool                   `json:"project_config_exists"`
	Detected            detect.DetectedProject `json:"detected"`
	RecommendedConfig   config.ProjectConfig   `json:"recommended_config"`
	ExistingConfig      *config.ProjectConfig  `json:"existing_config,omitempty"`
}

// Detect inspects a project directory without changing disk state.
func (m *Manager) Detect(dir string) (DetectResult, []string, error) {
	normalizedDir, err := normalizeDir(dir)
	if err != nil {
		return DetectResult{}, nil, runtimeFailure(err, "não foi possível resolver o diretório do projeto")
	}

	detected := detect.Project(normalizedDir)
	result := DetectResult{
		Directory:           normalizedDir,
		ProjectConfigPath:   projectConfigPath(normalizedDir),
		ProjectConfigExists: config.ExistsProject(normalizedDir),
		Detected:            detected,
		RecommendedConfig:   recommendedConfig(detected.Name, detected.Runtime, detected.Framework, detected.Port),
	}

	if result.ProjectConfigExists {
		cfg, err := config.LoadProject(result.ProjectConfigPath)
		if err != nil {
			return DetectResult{}, nil, configurationError(err, "não foi possível ler o arquivo .odins em %s", normalizedDir)
		}
		result.ExistingConfig = &cfg
	}

	var warnings []string
	if detected.Runtime == "unknown" {
		warnings = append(warnings, "ODINS não conseguiu detectar automaticamente runtime/framework neste diretório.")
	}
	if detected.HasDocker {
		warnings = append(warnings, "Dockerfile encontrado; revise se as rotas devem apontar para portas locais expostas pelo container.")
	}
	if detected.HasCompose {
		warnings = append(warnings, "Arquivo de compose encontrado; confirme se o nome do container exposto deve entrar no campo docker_container.")
	}

	return result, warnings, nil
}
