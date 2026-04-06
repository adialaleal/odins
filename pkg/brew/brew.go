package brew

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsInstalled returns true if Homebrew is available on PATH.
func IsInstalled() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

// Install installs a Homebrew formula if not already present.
func Install(formula string) error {
	if IsFormulaInstalled(formula) {
		return nil
	}
	fmt.Printf("  → brew install %s\n", formula)
	cmd := exec.Command("brew", "install", formula)
	cmd.Stdout = nil
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("brew install %s: %w\n%s", formula, err, string(out))
	}
	return nil
}

// IsFormulaInstalled returns true if the formula is already installed.
func IsFormulaInstalled(formula string) bool {
	cmd := exec.Command("brew", "list", "--formula", formula)
	return cmd.Run() == nil
}

// ServiceStart starts a Homebrew service.
func ServiceStart(service string) error {
	return runService("start", service)
}

// ServiceStop stops a Homebrew service.
func ServiceStop(service string) error {
	return runService("stop", service)
}

// ServiceRestart restarts a Homebrew service.
func ServiceRestart(service string) error {
	return runService("restart", service)
}

// ServiceReload sends a reload signal to a Homebrew service.
func ServiceReload(service string) error {
	return runService("reload", service)
}

// ServiceRunning returns true if the service is currently running.
func ServiceRunning(service string) bool {
	out, err := exec.Command("brew", "services", "list").CombinedOutput()
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == service {
			return fields[1] == "started"
		}
	}
	return false
}

func runService(action, service string) error {
	out, err := exec.Command("brew", "services", action, service).CombinedOutput()
	if err != nil {
		return fmt.Errorf("brew services %s %s: %w\n%s", action, service, err, string(out))
	}
	return nil
}
