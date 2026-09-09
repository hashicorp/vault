// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEntrypointUserSetup(t *testing.T) {
	// The entrypoints use substring expansion supported by their container
	// shells and bash, but not by every host's /bin/sh (for example, dash).
	shell, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash is required to execute the entrypoint scripts")
	}

	for _, script := range []struct {
		name, path, switchUser string
	}{
		{"development", "docker-entrypoint.sh", "su-exec"},
		{"release", "../../.release/docker/docker-entrypoint.sh", "su-exec"},
		{"ubi", "../../.release/docker/ubi-docker-entrypoint.sh", "su"},
	} {
		t.Run(script.name, func(t *testing.T) {
			scriptPath, err := filepath.Abs(script.path)
			if err != nil {
				t.Fatal(err)
			}
			if out, err := exec.Command(shell, "-n", scriptPath).CombinedOutput(); err != nil {
				t.Fatalf("entrypoint syntax: %v\n%s", err, out)
			}
			for _, tc := range []struct {
				name                        string
				root, skipChown, skipSetcap bool
				versionFails, disableMlock  bool
				exitCode                    int
			}{
				{name: "nonroot"},
				{name: "nonroot_skip_chown", skipChown: true},
				{name: "nonroot_skip_setcap", skipSetcap: true},
				{name: "nonroot_skip_both", skipChown: true, skipSetcap: true},
				{name: "nonroot_mlock_disabled", disableMlock: true},
				{name: "nonroot_exit_status", exitCode: 23},
				{name: "root", root: true},
				{name: "root_skip_chown", root: true, skipChown: true},
				{name: "root_skip_setcap", root: true, skipSetcap: true},
				{name: "root_skip_both", root: true, skipChown: true, skipSetcap: true},
				{name: "root_capability_fallback", root: true, versionFails: true},
				{name: "root_exit_status", root: true, exitCode: 23},
			} {
				t.Run(tc.name, func(t *testing.T) {
					binDir := t.TempDir()
					logPath := filepath.Join(binDir, "calls")
					for _, command := range []string{"id", "stat", "chown", "setcap", "readlink", "which", "su-exec", "su", "vault"} {
						if err := os.WriteFile(filepath.Join(binDir, command), []byte(entrypointCommandStub), 0o755); err != nil {
							t.Fatal(err)
						}
					}
					cmd := exec.Command(shell, scriptPath, "vault", "server", "-config=/tmp/config with spaces")
					// Start from a clean environment so host Vault settings cannot
					// write config files or change the behavior under test.
					cmd.Env = []string{
						"PATH=" + binDir + string(os.PathListSeparator) + os.Getenv("PATH"),
						"ENTRYPOINT_SHELL=" + shell,
						"COMMAND_LOG=" + logPath,
						fmt.Sprintf("MOCK_EXIT_CODE=%d", tc.exitCode),
					}
					if tc.root {
						cmd.Env = append(cmd.Env, "MOCK_UID=0")
					} else {
						cmd.Env = append(cmd.Env, "MOCK_UID=1000")
					}
					if tc.skipChown {
						cmd.Env = append(cmd.Env, "SKIP_CHOWN=true")
					}
					if tc.skipSetcap {
						cmd.Env = append(cmd.Env, "SKIP_SETCAP=true")
					}
					if tc.versionFails {
						cmd.Env = append(cmd.Env, "MOCK_VERSION_EXIT=1")
					}
					if tc.disableMlock {
						cmd.Env = append(cmd.Env, "VAULT_DISABLE_MLOCK=true")
					}
					out, err := cmd.CombinedOutput()
					if cmd.ProcessState == nil || cmd.ProcessState.ExitCode() != tc.exitCode {
						t.Fatalf("entrypoint exit status: got %v, want %d\n%s", err, tc.exitCode, out)
					}
					calls, err := os.ReadFile(logPath)
					if err != nil {
						t.Fatal(err)
					}
					log := string(calls)
					assertCount := func(call string, want int) {
						t.Helper()
						if got := strings.Count(log, call); got != want {
							t.Errorf("got %d calls containing %q, want %d\n%s", got, call, want, log)
						}
					}
					assertCount("vault <server> <-config=/tmp/config with spaces>", 1)
					for _, dir := range []string{"config", "logs", "file"} {
						want := 0
						if tc.root && !tc.skipChown {
							want = 1
						}
						assertCount("chown <-R> <vault:vault> </vault/"+dir+">", want)
					}
					setcapCount, fallbackCount, switchCount := 0, 0, 0
					if tc.root {
						switchCount = 1
						if !tc.skipSetcap {
							setcapCount = 1
							if tc.versionFails {
								fallbackCount = 1
							}
						}
					}
					assertCount("setcap <cap_ipc_lock=+ep>", setcapCount)
					assertCount("setcap <cap_ipc_lock=-ep>", fallbackCount)
					assertCount(script.switchUser+" <vault>", switchCount)
					if !tc.root {
						for name, enabled := range map[string]bool{"SKIP_CHOWN": tc.skipChown, "SKIP_SETCAP": tc.skipSetcap} {
							warning := "Container is running as non-root user, ignoring " + name
							if strings.Contains(string(out), warning) != enabled {
								t.Errorf("unexpected %s warning: %s", name, out)
							}
						}
					}
				})
			}
		})
	}
}

// All privileged commands are replaced with stubs. Tests never modify the
// host's ownership, capabilities, or user identity.
const entrypointCommandStub = `#!/bin/sh
command=${0##*/}
if [ "$command" = vault ] && [ "$1" = --help ]; then
    exit 1
fi
{
    printf '%s' "$command"
    printf ' <%s>' "$@"
    printf '\n'
} >> "$COMMAND_LOG"
case "$command" in
    id)
        if [ "$2" = vault ]; then printf '1000\n'; else printf '%s\n' "$MOCK_UID"; fi
        ;;
    stat) printf '0\n' ;;
    which) command -v "$1" ;;
    readlink) printf '%s\n' "$2" ;;
    chown|setcap)
        [ "$MOCK_UID" = 0 ] || exit 99
        ;;
    su-exec)
        shift
        export MOCK_UID=1000
        exec "$@"
        ;;
    su)
        shift 2
        script=$1
        shift 2
        export MOCK_UID=1000
        exec "$ENTRYPOINT_SHELL" "$script" "$@"
        ;;
    vault)
        if [ "$1" = -version ]; then exit "${MOCK_VERSION_EXIT:-0}"; fi
        exit "$MOCK_EXIT_CODE"
        ;;
esac
`
