// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command_test

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

// TestClientCommands_SSHNoExecCA verifies that vault-client ssh in CA mode has
// Vault sign the user's public key with the given mount and role, and that
// -no-exec prints the signed certificate instead of running ssh.
func TestClientCommands_SSHNoExecCA(t *testing.T) {
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
	require.NoError(t, os.WriteFile(publicKeyPath, ssh.MarshalAuthorizedKey(userKey), 0o600), "write a test-local public key")

	var stdout, stderr bytes.Buffer
	code := clientcmd.RunCustom([]string{
		"ssh", "-mode=ca", "-role=test-role", "-mount-point=custom-ssh",
		"-no-exec", "-field=signed_key", "-public-key-path=" + publicKeyPath,
		"-ssh-executable=" + filepath.Join(dir, "does-not-exist"), "alice@127.0.0.1",
	}, &base.RunOptions{
		BinaryName: "vault-client", Client: client, HCPTokenHelper: clientAuthNoHCPTokenHelper{},
		Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stderr,
	})
	require.Equal(t, 0, code, "no-exec must succeed without an SSH executable: %s", stderr.String())
	require.Empty(t, stderr.String(), "explicit mode and role must not warn")
	require.NotContains(t, stdout.String(), client.Token(), "Vault authorization must not be echoed")

	signed, _, _, _, err := ssh.ParseAuthorizedKey(stdout.Bytes())
	require.NoError(t, err, "-field=signed_key must print the signed certificate")
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
}

// TestClientCommands_SSHNoExecOTP verifies that vault-client ssh in OTP mode
// has Vault issue an OTP for the user and IP in the SSH destination, with the
// given mount and role, and that -no-exec prints the OTP instead of running
// ssh. The role's default user differs from the destination's user, so the OTP
// is only issued for alice if the command sends her name.
func TestClientCommands_SSHNoExecOTP(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	require.NoError(t, client.Sys().MountWithContext(t.Context(), "custom-ssh", &api.MountInput{
		Type: "ssh",
	}), "enable the SSH secrets engine at a custom mount")
	_, err := client.Logical().WriteWithContext(t.Context(), "custom-ssh/roles/test-role", map[string]any{
		"key_type": "otp", "default_user": "test-default", "allowed_users": "alice",
		"cidr_list": "127.0.0.1/32", "port": 2222,
	})
	require.NoError(t, err, "create an OTP role")

	var stdout, stderr bytes.Buffer
	code := clientcmd.RunCustom([]string{
		"ssh", "-mode=otp", "-role=test-role", "-mount-point=custom-ssh", "-no-exec", "-field=key",
		"-ssh-executable=" + filepath.Join(t.TempDir(), "does-not-exist"), "alice@127.0.0.1",
	}, &base.RunOptions{
		BinaryName: "vault-client", Client: client, HCPTokenHelper: clientAuthNoHCPTokenHelper{},
		Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stderr,
	})
	require.Equal(t, 0, code, "no-exec must succeed without an SSH executable: %s", stderr.String())
	require.Empty(t, stderr.String(), "explicit mode and role must not warn")
	require.NotContains(t, stdout.String(), client.Token(), "Vault authorization must not be echoed")

	otp := strings.TrimSpace(stdout.String())
	require.NotEmpty(t, otp, "-field=key must print the OTP")
	verify, err := client.Logical().WriteWithContext(t.Context(), "custom-ssh/verify", map[string]any{
		"otp": otp,
	})
	require.NoError(t, err, "Vault must accept the printed OTP")
	require.NotNil(t, verify, "Vault must recognize the printed OTP")
	require.Equal(t, map[string]any{
		"username": "alice", "ip": "127.0.0.1", "role_name": "test-role",
	}, verify.Data, "issue the OTP for the SSH destination with the given role")
}
