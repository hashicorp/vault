// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func (s *Session) MustRegisterPlugin(pluginName, binaryPath, pluginType string) {
	s.t.Helper()

	f, err := os.Open(binaryPath)
	require.NoError(s.t, err)
	defer func() { _ = f.Close() }()

	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	require.NoError(s.t, err)

	shaSum := hex.EncodeToString(hasher.Sum(nil))

	payload := map[string]any{
		"sha256":  shaSum,
		"command": filepath.Base(binaryPath),
		"type":    pluginType,
	}

	s.MustWrite(filepath.Join("sys/plugins/catalog", pluginType, pluginName), payload)
}

func (s *Session) MustEnablePlugin(path, pluginName, pluginType string) {
	s.t.Helper()

	switch pluginType {
	case "auth":
		s.MustEnableAuth(path, &api.EnableAuthOptions{Type: pluginName})
	case "secret":
		s.MustEnableSecretsEngine(path, &api.MountInput{Type: pluginName})
	default:
		s.t.Fatalf("unknown plugin type: %s", pluginType)
	}
}

func (s *Session) AssertPluginRegistered(pluginName string) {
	s.t.Helper()

	secret := s.MustRead(filepath.Join("sys/plugins/catalog", pluginName))
	require.NotNil(s.t, secret)
}

func (s *Session) AssertPluginConfigured(path string) {
	s.t.Helper()

	configPath := filepath.Join(path, "config")
	s.MustRead(configPath)
}
