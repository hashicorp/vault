// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package spiffe

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/client/auth"
	"github.com/stretchr/testify/require"
)

// generateTestCertAndKey writes a self-signed ECDSA client certificate and
// private key as PEM files in dir and returns their paths.
func generateTestCertAndKey(t *testing.T, dir string) (string, string) {
	t.Helper()

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"SPIFFE Test"},
		},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	require.NoError(t, err)

	certPath := filepath.Join(dir, "cert.pem")
	certOut, err := os.Create(certPath)
	require.NoError(t, err)
	defer certOut.Close()
	require.NoError(t, pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))

	keyBytes, err := x509.MarshalECPrivateKey(priv)
	require.NoError(t, err)

	keyPath := filepath.Join(dir, "key.pem")
	keyOut, err := os.Create(keyPath)
	require.NoError(t, err)
	defer keyOut.Close()
	require.NoError(t, pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}))

	return certPath, keyPath
}

// TestNewSpiffeAuthMethod_ConfigValidation verifies that NewSpiffeAuthMethod
// rejects nil, incomplete, or invalid configs with a descriptive error, and
// accepts valid jwt and cert configs.
func TestNewSpiffeAuthMethod_ConfigValidation(t *testing.T) {
	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "svid.token")
	require.NoError(t, os.WriteFile(tokenFile, []byte("test-jwt-token"), 0o600))

	certPath, keyPath := generateTestCertAndKey(t, tempDir)
	validTLSCfg := &api.TLSConfig{
		ClientCert: certPath,
		ClientKey:  keyPath,
	}

	tests := []struct {
		name      string
		config    *auth.AuthConfig
		tlsCfg    *api.TLSConfig
		wantErr   bool
		errString string
	}{
		{
			name:      "nil config",
			config:    nil,
			wantErr:   true,
			errString: "empty config",
		},
		{
			name: "missing config map",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
			},
			wantErr:   true,
			errString: "missing config",
		},
		{
			name: "missing type",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config:    map[string]interface{}{},
			},
			wantErr:   true,
			errString: "missing required 'type' config value",
		},
		{
			name: "invalid type value",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type": "unsupported",
				},
			},
			wantErr:   true,
			errString: "invalid auth type \"unsupported\"",
		},
		{
			name: "jwt type missing svid_token_file",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type": "jwt",
				},
			},
			wantErr:   true,
			errString: "missing required 'svid_token_file' config value",
		},
		{
			name: "jwt type empty svid_token_file",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type":            "jwt",
					"svid_token_file": "",
				},
			},
			wantErr:   true,
			errString: "'svid_token_file' config value cannot be empty",
		},
		{
			name: "jwt valid config with role",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type":            "jwt",
					"role":            "dev-role",
					"svid_token_file": tokenFile,
				},
			},
			wantErr: false,
		},
		{
			name: "jwt valid config without role",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type":            "jwt",
					"svid_token_file": tokenFile,
				},
			},
			wantErr: false,
		},
		{
			name: "cert type no tls anywhere",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type": "cert",
				},
			},
			tlsCfg:    nil,
			wantErr:   true,
			errString: "missing required client certificate and key in tls configuration",
		},
		{
			name: "cert type only client_cert in stanza",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type":        "cert",
					"client_cert": certPath,
				},
			},
			tlsCfg:    nil,
			wantErr:   true,
			errString: "missing required client certificate and key in tls configuration",
		},
		{
			name: "cert type non-string client_cert",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type":        "cert",
					"client_cert": 1,
					"client_key":  keyPath,
				},
			},
			wantErr:   true,
			errString: "could not convert 'client_cert' config value to string",
		},
		{
			name: "cert type non-string client_key",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type":        "cert",
					"client_cert": certPath,
					"client_key":  1,
				},
			},
			wantErr:   true,
			errString: "could not convert 'client_key' config value to string",
		},
		{
			name: "cert valid config via stanza",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type":        "cert",
					"client_cert": certPath,
					"client_key":  keyPath,
				},
			},
			tlsCfg:  nil,
			wantErr: false,
		},
		{
			name: "cert valid config",
			config: &auth.AuthConfig{
				Logger:    hclog.NewNullLogger(),
				MountPath: "auth/spiffe",
				Config: map[string]interface{}{
					"type": "cert",
					"role": "cert-role",
				},
			},
			tlsCfg:  validTLSCfg,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, err := NewSpiffeAuthMethod(tt.config, tt.tlsCfg)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errString != "" {
					require.True(t, strings.Contains(err.Error(), tt.errString), "error %q did not contain %q", err.Error(), tt.errString)
				}
				return
			}
			require.NoError(t, err)
			require.NotNil(t, method)
			defer method.Shutdown()
		})
	}
}

// TestSpiffeAuthMethod_JWTAuthenticate verifies that jwt-type Authenticate
// reads the SVID token from disk, trims trailing whitespace, and sends it as a
// Bearer Authorization header to the mount's login path along with the role.
func TestSpiffeAuthMethod_JWTAuthenticate(t *testing.T) {
	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "svid.token")
	tokenContent := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.dummy.sig"
	require.NoError(t, os.WriteFile(tokenFile, []byte(tokenContent+"\n"), 0o600))

	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/spiffe",
		Config: map[string]interface{}{
			"type":            "jwt",
			"role":            "myrole",
			"svid_token_file": tokenFile,
		},
	}

	method, err := NewSpiffeAuthMethod(conf, nil)
	require.NoError(t, err)
	defer method.Shutdown()

	path, header, data, err := method.Authenticate(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, "auth/spiffe/login", path)
	require.Equal(t, "Bearer "+tokenContent, header.Get("Authorization"))
	require.Equal(t, "myrole", data["role"])
}

// TestSpiffeAuthMethod_CertAuthClientAndAuthenticate verifies that cert-type
// auth builds an auth client from the configured TLS cert and key, and that
// Authenticate sends no header and omits the role when none is set.
func TestSpiffeAuthMethod_CertAuthClientAndAuthenticate(t *testing.T) {
	tempDir := t.TempDir()
	certPath, keyPath := generateTestCertAndKey(t, tempDir)

	tlsCfg := &api.TLSConfig{
		ClientCert: certPath,
		ClientKey:  keyPath,
	}

	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/spiffe",
		Config: map[string]interface{}{
			"type": "cert",
		},
	}

	method, err := NewSpiffeAuthMethod(conf, tlsCfg)
	require.NoError(t, err)
	defer method.Shutdown()

	client, err := api.NewClient(api.DefaultConfig())
	require.NoError(t, err)

	authWithClient, ok := method.(auth.AuthMethodWithClient)
	require.True(t, ok)

	authClient, err := authWithClient.AuthClient(client)
	require.NoError(t, err)
	require.NotNil(t, authClient)

	path, header, data, err := method.Authenticate(context.Background(), authClient)
	require.NoError(t, err)
	require.Equal(t, "auth/spiffe/login", path)
	require.Nil(t, header)
	require.Equal(t, "cert", data["type"])
	require.NotContains(t, data, "role")
}

// TestNewSpiffeAuthMethod_StanzaOverridesVaultTLS verifies that client_cert and
// client_key in the method stanza take precedence over the values inherited
// from the vault stanza, that other inherited TLS fields are preserved, and
// that the caller's shared *api.TLSConfig is not mutated.
func TestNewSpiffeAuthMethod_StanzaOverridesVaultTLS(t *testing.T) {
	tempDir := t.TempDir()
	stanzaCert, stanzaKey := generateTestCertAndKey(t, tempDir)

	vaultTLS := &api.TLSConfig{
		CACert:        "/vault/ca.pem",
		ClientCert:    "/vault/cert.pem",
		ClientKey:     "/vault/key.pem",
		TLSServerName: "vault.example.com",
	}

	conf := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "auth/spiffe",
		Config: map[string]interface{}{
			"type":        "cert",
			"client_cert": stanzaCert,
			"client_key":  stanzaKey,
		},
	}

	method, err := NewSpiffeAuthMethod(conf, vaultTLS)
	require.NoError(t, err)
	defer method.Shutdown()

	got := method.(*spiffeMethod).tlsCfg
	require.Equal(t, stanzaCert, got.ClientCert)
	require.Equal(t, stanzaKey, got.ClientKey)
	require.Equal(t, "/vault/ca.pem", got.CACert)
	require.Equal(t, "vault.example.com", got.TLSServerName)

	require.Equal(t, "/vault/cert.pem", vaultTLS.ClientCert, "shared vault TLS config must not be mutated")
	require.Equal(t, "/vault/key.pem", vaultTLS.ClientKey, "shared vault TLS config must not be mutated")
}

// TestWaitForCreds_JWT verifies that WaitForCreds returns nil as soon as the
// SVID token file contains a non-empty token, and blocks until then. It covers
// three cases: file already present on entry, file written after entry, and
// context cancellation before the file appears.
func TestWaitForCreds_JWT(t *testing.T) {
	t.Parallel()

	t.Run("returns immediately when file already present", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		tokenFile := filepath.Join(dir, "svid.jwt")
		require.NoError(t, os.WriteFile(tokenFile, []byte("my-jwt-token"), 0o600))

		method, err := NewSpiffeAuthMethod(&auth.AuthConfig{
			Logger:    hclog.NewNullLogger(),
			MountPath: "auth/spiffe",
			Config: map[string]interface{}{
				"type":            "jwt",
				"svid_token_file": tokenFile,
			},
		}, nil)
		require.NoError(t, err)
		defer method.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = method.(auth.AuthMethodReady).WaitForCreds(ctx)
		require.NoError(t, err)
	})

	t.Run("blocks until file is written", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		tokenFile := filepath.Join(dir, "svid.jwt")
		// File does not exist yet.

		method, err := NewSpiffeAuthMethod(&auth.AuthConfig{
			Logger:    hclog.NewNullLogger(),
			MountPath: "auth/spiffe",
			Config: map[string]interface{}{
				"type":            "jwt",
				"svid_token_file": tokenFile,
			},
		}, nil)
		require.NoError(t, err)
		defer method.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Write the file after a short delay.
		go func() {
			time.Sleep(300 * time.Millisecond)
			_ = os.WriteFile(tokenFile, []byte("delayed-jwt-token"), 0o600)
		}()

		err = method.(auth.AuthMethodReady).WaitForCreds(ctx)
		require.NoError(t, err)
	})

	t.Run("returns error when context is cancelled", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		tokenFile := filepath.Join(dir, "svid.jwt")
		// File will never be written.

		method, err := NewSpiffeAuthMethod(&auth.AuthConfig{
			Logger:    hclog.NewNullLogger(),
			MountPath: "auth/spiffe",
			Config: map[string]interface{}{
				"type":            "jwt",
				"svid_token_file": tokenFile,
			},
		}, nil)
		require.NoError(t, err)
		defer method.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		err = method.(auth.AuthMethodReady).WaitForCreds(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})
}

// TestWaitForCreds_Cert verifies that WaitForCreds returns nil as soon as the
// client cert and key files form a valid TLS key pair, and blocks until then.
// It mirrors the JWT cases: file already present, file written after entry,
// and context cancellation.
func TestWaitForCreds_Cert(t *testing.T) {
	t.Parallel()

	t.Run("returns immediately when cert/key already present", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		certPath, keyPath := generateTestCertAndKey(t, dir)

		method, err := NewSpiffeAuthMethod(&auth.AuthConfig{
			Logger:    hclog.NewNullLogger(),
			MountPath: "auth/spiffe",
			Config: map[string]interface{}{
				"type": "cert",
			},
		}, &api.TLSConfig{
			ClientCert: certPath,
			ClientKey:  keyPath,
		})
		require.NoError(t, err)
		defer method.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = method.(auth.AuthMethodReady).WaitForCreds(ctx)
		require.NoError(t, err)
	})

	t.Run("blocks until cert/key are written", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		certPath := filepath.Join(dir, "cert.pem")
		keyPath := filepath.Join(dir, "key.pem")
		// Files do not exist yet.

		method, err := NewSpiffeAuthMethod(&auth.AuthConfig{
			Logger:    hclog.NewNullLogger(),
			MountPath: "auth/spiffe",
			Config: map[string]interface{}{
				"type": "cert",
			},
		}, &api.TLSConfig{
			ClientCert: certPath,
			ClientKey:  keyPath,
		})
		require.NoError(t, err)
		defer method.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Write the cert/key after a short delay.
		go func() {
			time.Sleep(300 * time.Millisecond)
			// Reuse generateTestCertAndKey but write to the fixed paths.
			tmpDir := t.TempDir()
			srcCert, srcKey := generateTestCertAndKey(t, tmpDir)
			certBytes, _ := os.ReadFile(srcCert)
			keyBytes, _ := os.ReadFile(srcKey)
			_ = os.WriteFile(certPath, certBytes, 0o600)
			_ = os.WriteFile(keyPath, keyBytes, 0o600)
		}()

		err = method.(auth.AuthMethodReady).WaitForCreds(ctx)
		require.NoError(t, err)
	})

	t.Run("returns error when context is cancelled", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		certPath := filepath.Join(dir, "cert.pem")
		keyPath := filepath.Join(dir, "key.pem")
		// Files will never be written.

		method, err := NewSpiffeAuthMethod(&auth.AuthConfig{
			Logger:    hclog.NewNullLogger(),
			MountPath: "auth/spiffe",
			Config: map[string]interface{}{
				"type": "cert",
			},
		}, &api.TLSConfig{
			ClientCert: certPath,
			ClientKey:  keyPath,
		})
		require.NoError(t, err)
		defer method.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancel()

		err = method.(auth.AuthMethodReady).WaitForCreds(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
