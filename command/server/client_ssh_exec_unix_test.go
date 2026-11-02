// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build unix

package command_test

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

// TestClientCommands_SSHCAExecution verifies that vault-client ssh in CA mode
// runs the SSH executable with the certificate Vault signed and the user's
// arguments, returns the executable's exit status, and removes the certificate
// afterward. The executable is a small script that records its arguments and
// copies the certificate, so no SSH server is needed.
func TestClientCommands_SSHCAExecution(t *testing.T) {
	t.Parallel()

	for _, exitCode := range []int{0, 23} {
		t.Run(strconv.Itoa(exitCode), func(t *testing.T) {
			t.Parallel()

			cluster := minimal.NewTestSoloCluster(t, nil)
			client := cluster.Cores[0].Client
			require.NoError(t, client.Sys().MountWithContext(t.Context(), "custom-ssh", &api.MountInput{
				Type: "ssh",
			}), "enable the SSH secrets engine at a custom mount")
			ca, err := client.Logical().WriteWithContext(t.Context(), "custom-ssh/config/ca", map[string]any{
				"key_type": "ed25519",
			})
			require.NoError(t, err, "generate an SSH CA")
			caKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(ca.Data["public_key"].(string)))
			require.NoError(t, err, "parse the SSH CA public key")
			_, err = client.Logical().WriteWithContext(t.Context(), "custom-ssh/roles/test-role", map[string]any{
				"key_type": "ca", "allow_user_certificates": true, "allowed_users": "alice",
				"allowed_extensions": "permit-X11-forwarding,permit-agent-forwarding,permit-port-forwarding,permit-pty,permit-user-rc",
			})
			require.NoError(t, err, "create a CA role")

			rawUserKey, _, err := ed25519.GenerateKey(rand.Reader)
			require.NoError(t, err, "generate a test-local SSH key")
			userKey, err := ssh.NewPublicKey(rawUserKey)
			require.NoError(t, err, "convert the test-local key to SSH format")
			dir := t.TempDir()
			publicKeyPath := filepath.Join(dir, "id.pub")
			privateKeyPath := filepath.Join(dir, "id")
			knownHostsPath := filepath.Join(dir, "known_hosts")
			argvPath := filepath.Join(dir, "argv")
			certificateCopyPath := filepath.Join(dir, "certificate")
			executablePath := filepath.Join(dir, "ssh-fixture")
			require.NoError(t, os.WriteFile(publicKeyPath, ssh.MarshalAuthorizedKey(userKey), 0o600), "write a test-local public key")
			require.NoError(t, os.WriteFile(privateKeyPath, []byte("test-private-key"), 0o600), "write an inert test-local private key")
			require.NoError(t, os.WriteFile(knownHostsPath, nil, 0o600), "avoid the user's known_hosts file")
			require.NoError(t, os.WriteFile(executablePath, []byte(`#!/bin/sh
set -eu
umask 077
printf '%s\n' "$@" > "${10}"
cat "$4" > "${11}"
exit "${12}"
`), 0o700), "write a finite SSH executable fixture")

			sshArgs := []string{
				"alice@127.0.0.1", "-p", "2222", argvPath, certificateCopyPath,
				strconv.Itoa(exitCode), "remote argument; $(not-a-shell)",
			}
			args := append([]string{
				"ssh", "-mode=ca", "-role=test-role", "-mount-point=custom-ssh",
				"-public-key-path=" + publicKeyPath, "-private-key-path=" + privateKeyPath,
				"-host-key-mount-point=", "-strict-host-key-checking=yes",
				"-user-known-hosts-file=" + knownHostsPath, "-ssh-executable=" + executablePath,
			}, sshArgs...)
			var stdout, stderr bytes.Buffer
			code := clientcmd.RunCustom(args, &base.RunOptions{
				BinaryName: "vault-client", Client: client, HCPTokenHelper: clientAuthNoHCPTokenHelper{},
				Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stderr,
			})
			require.Equal(t, exitCode, code, "propagate the SSH executable status: %s", stderr.String())
			if exitCode != 0 {
				require.Contains(t, stderr.String(), "failed to run ssh command", "retain SSH-specific execution errors")
			} else {
				require.Empty(t, stderr.String(), "successful execution must not emit an error")
			}
			require.NotContains(t, stdout.String()+stderr.String(), client.Token(), "Vault authorization must not be echoed")

			argv, err := os.ReadFile(argvPath)
			require.NoError(t, err, "the SSH fixture must actually execute")
			gotArgs := strings.Split(strings.TrimSuffix(string(argv), "\n"), "\n")
			require.Len(t, gotArgs, 6+len(sshArgs), "preserve the entire SSH argument vector")
			certificatePath := gotArgs[3]
			require.Equal(t, append([]string{
				"-i", privateKeyPath, "-i", certificatePath,
				"-o StrictHostKeyChecking=yes", "-o UserKnownHostsFile=" + knownHostsPath,
			}, sshArgs...), gotArgs, "forward key paths, host checking, destination, port, and remote arguments intact")

			certificate, err := os.ReadFile(certificateCopyPath)
			require.NoError(t, err, "the signed certificate must be readable during execution")
			signed, _, _, _, err := ssh.ParseAuthorizedKey(certificate)
			require.NoError(t, err, "pass the signed certificate to the SSH executable")
			cert, ok := signed.(*ssh.Certificate)
			require.True(t, ok, "the signed key must be an SSH certificate")
			require.Equal(t, userKey.Marshal(), cert.Key.Marshal(), "sign the key from -public-key-path")
			require.Equal(t, caKey.Marshal(), cert.SignatureKey.Marshal(), "sign with the CA at -mount-point")
			require.Equal(t, uint32(ssh.UserCert), cert.CertType, "request a user certificate")
			require.Equal(t, []string{"alice"}, cert.ValidPrincipals, "sign for the user in the SSH destination")
			require.Equal(t, map[string]string{
				"permit-X11-forwarding": "", "permit-agent-forwarding": "", "permit-port-forwarding": "",
				"permit-pty": "", "permit-user-rc": "",
			}, cert.Permissions.Extensions, "request the default user certificate extensions")

			_, err = os.Stat(certificatePath)
			require.ErrorIs(t, err, os.ErrNotExist, "remove the temporary certificate after either execution outcome")
		})
	}
}
