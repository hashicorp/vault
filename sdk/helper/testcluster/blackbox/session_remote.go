// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// RemoteHost represents a remote host configuration
type RemoteHost struct {
	PublicIP  string `json:"public_ip"`
	PrivateIP string `json:"private_ip"`
}

// AssertRemoteCLIVersion verifies the Vault CLI version on a remote host via SSH
// This method SSHs to the remote host and runs the vault version command
func (s *Session) AssertRemoteCLIVersion(host RemoteHost, vaultInstallDir, version, sha, buildDate, edition string) {
	s.t.Helper()

	// Build the vault version command
	vaultBinary := fmt.Sprintf("%s/vault", vaultInstallDir)
	remoteCmd := fmt.Sprintf("%s version", vaultBinary)

	// Execute SSH command
	cmd := exec.Command("ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "LogLevel=ERROR",
		host.PublicIP,
		remoteCmd,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		s.t.Fatalf("Failed to execute vault version on remote host %s: %v\nStderr: %s", host.PublicIP, err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())

	// Build expected version string
	expectedVersion := fmt.Sprintf("Vault v%s ('%s'), built %s", version, sha, buildDate)

	switch edition {
	case "ce", "ent":
		// No additional suffix
	case "ent.hsm", "ent.fips1403", "ent.hsm.fips1403":
		expectedVersion += " (cgo)"
	default:
		s.t.Fatalf("unknown Vault edition: %s", edition)
	}

	// Also check version without SHA (some builds may not include it)
	expectedVersionNoSHA := strings.Replace(expectedVersion, fmt.Sprintf("('%s') ", sha), "", 1)
	expectedVersionNoSHA = strings.TrimSpace(strings.Replace(expectedVersionNoSHA, "  ", " ", -1))

	if output != expectedVersion && output != expectedVersionNoSHA {
		s.t.Fatalf("CLI version mismatch on host %s.\nExpected: %s\nor: %s\nGot: %s",
			host.PublicIP, expectedVersion, expectedVersionNoSHA, output)
	}

	s.t.Logf("CLI version verification succeeded on host %s: %s", host.PublicIP, output)
}
