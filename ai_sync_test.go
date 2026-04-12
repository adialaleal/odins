package main

import (
	"os/exec"
	"testing"
)

func TestAIPacksAreSynced(t *testing.T) {
	t.Parallel()

	cmd := exec.Command("go", "run", "./tools/ai-pack-gen", "--check")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run ./tools/ai-pack-gen --check failed: %v\n%s", err, string(out))
	}
}
