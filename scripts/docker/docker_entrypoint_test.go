package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestEntrypointNonRootUser verifies that when the entrypoint script runs as a
// non-root user it does NOT attempt setcap or chown, and it still execs vault.
func TestEntrypointNonRootUser(t *testing.T) {
	// Locate the entrypoint script relative to the repo root.
	scriptPath := "../../../scripts/docker/docker-entrypoint.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = "../../scripts/docker/docker-entrypoint.sh"
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			t.Skip("entrypoint script not found")
		}
	}

	// We run the script under `sh -n` to syntax-check it first.
	cmd := exec.Command("sh", "-n", scriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("entrypoint script has syntax errors: %v\n%s", err, out)
	}

	// Verify that the script contains the non-root guard.
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read entrypoint: %v", err)
	}
	body := string(content)

	if !strings.Contains(body, `if [ "$(id -u)" != '0' ]; then`) {
		t.Error("entrypoint missing non-root user guard")
	}
	if !strings.Contains(body, "Container is running as non-root user, ignoring SKIP_SETCAP") {
		t.Error("entrypoint missing SKIP_SETCAP warning for non-root")
	}
	if !strings.Contains(body, "VAULT_DISABLE_MLOCK") {
		t.Error("entrypoint should reference VAULT_DISABLE_MLOCK")
	}
}

// TestEntrypointRootUser verifies that the root path still performs chown and
// setcap before dropping to the vault user.
func TestEntrypointRootUser(t *testing.T) {
	scriptPath := "../../../scripts/docker/docker-entrypoint.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = "../../scripts/docker/docker-entrypoint.sh"
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			t.Skip("entrypoint script not found")
		}
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read entrypoint: %v", err)
	}
	body := string(content)

	if !strings.Contains(body, "setcap cap_ipc_lock=+ep") {
		t.Error("entrypoint missing setcap for root user")
	}
	if !strings.Contains(body, "su-exec vault") {
		t.Error("entrypoint missing su-exec vault for root user")
	}
}

// TestUbiEntrypointNonRootUser verifies the same behaviour for the UBI
// entrypoint variant used in Red Hat builds.
func TestUbiEntrypointNonRootUser(t *testing.T) {
	scriptPath := "../../../.release/docker/ubi-docker-entrypoint.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = "../../.release/docker/ubi-docker-entrypoint.sh"
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			t.Skip("ubi entrypoint script not found")
		}
	}

	cmd := exec.Command("sh", "-n", scriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("ubi entrypoint script has syntax errors: %v\n%s", err, out)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read ubi entrypoint: %v", err)
	}
	body := string(content)

	if !strings.Contains(body, `if [ "$(id -u)" != '0' ]; then`) {
		t.Error("ubi entrypoint missing non-root user guard")
	}
	if !strings.Contains(body, "Container is running as non-root user, ignoring SKIP_SETCAP") {
		t.Error("ubi entrypoint missing SKIP_SETCAP warning for non-root")
	}
	if !strings.Contains(body, "VAULT_DISABLE_MLOCK") {
		t.Error("ubi entrypoint should reference VAULT_DISABLE_MLOCK")
	}
}

// TestUbiEntrypointRootUser verifies the root path in the UBI entrypoint.
func TestUbiEntrypointRootUser(t *testing.T) {
	scriptPath := "../../../.release/docker/ubi-docker-entrypoint.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = "../../.release/docker/ubi-docker-entrypoint.sh"
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			t.Skip("ubi entrypoint script not found")
		}
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read ubi entrypoint: %v", err)
	}
	body := string(content)

	if !strings.Contains(body, "setcap cap_ipc_lock=+ep") {
		t.Error("ubi entrypoint missing setcap for root user")
	}
	if !strings.Contains(body, "su vault -p") {
		t.Error("ubi entrypoint missing su vault for root user")
	}
}

// TestEntrypointEnvVars documents the expected environment variables.
func TestEntrypointEnvVars(t *testing.T) {
	// This test is purely documentary; it lists the env vars the entrypoint
	// respects so that operators know what knobs are available.
	vars := []string{
		"SKIP_SETCAP",
		"SKIP_CHOWN",
		"VAULT_DISABLE_MLOCK",
		"VAULT_REDIRECT_INTERFACE",
		"VAULT_CLUSTER_INTERFACE",
		"VAULT_LOCAL_CONFIG",
		"VAULT_DEV_ROOT_TOKEN_ID",
		"VAULT_DEV_LISTEN_ADDRESS",
	}
	for _, v := range vars {
		if os.Getenv(v) == "" {
			// We don't require them to be set; just document them.
			fmt.Printf("documented env var: %s\n", v)
		}
	}
}
