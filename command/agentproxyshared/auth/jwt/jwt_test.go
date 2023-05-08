// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jwt

import (
	"bytes"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/auth"
)

func TestIngressToken(t *testing.T) {
	const (
		dir       = "dir"
		file      = "file"
		empty     = "empty"
		missing   = "missing"
		symlinked = "symlinked"
	)

	rootDir, err := os.MkdirTemp("", "vault-agent-jwt-auth-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(rootDir)

	setupTestDir := func() string {
		testDir, err := os.MkdirTemp(rootDir, "")
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(path.Join(testDir, file), []byte("test"), 0o644)
		if err != nil {
			t.Fatal(err)
		}
		_, err = os.Create(path.Join(testDir, empty))
		if err != nil {
			t.Fatal(err)
		}
		err = os.Mkdir(path.Join(testDir, dir), 0o755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.Symlink(path.Join(testDir, file), path.Join(testDir, symlinked))
		if err != nil {
			t.Fatal(err)
		}

		return testDir
	}

	for _, tc := range []struct {
		name      string
		path      string
		errString string
	}{
		{
			"happy path",
			file,
			"",
		},
		{
			"path is directory",
			dir,
			"[ERROR] jwt file is not a regular file or symlink",
		},
		{
			"path is symlink",
			symlinked,
			"",
		},
		{
			"path is missing (implies nothing for ingressToken to do)",
			missing,
			"",
		},
		{
			"path is empty file",
			empty,
			"[WARN]  empty jwt file read",
		},
	} {
		testDir := setupTestDir()
		logBuffer := bytes.Buffer{}
		jwtAuth := &jwtMethod{
			logger: hclog.New(&hclog.LoggerOptions{
				Output: &logBuffer,
			}),
			latestToken: new(atomic.Value),
			path:        path.Join(testDir, tc.path),
		}

		jwtAuth.ingressToken()

		if tc.errString != "" {
			if !strings.Contains(logBuffer.String(), tc.errString) {
				t.Fatal("logs did no contain expected error", tc.errString, logBuffer.String())
			}
		} else {
			if strings.Contains(logBuffer.String(), "[ERROR]") || strings.Contains(logBuffer.String(), "[WARN]") {
				t.Fatal("logs contained unexpected error", logBuffer.String())
			}
		}
	}
}

func TestDeleteAfterReading(t *testing.T) {
	for _, tc := range map[string]struct {
		configValue  string
		shouldDelete bool
	}{
		"default": {
			"",
			true,
		},
		"explicit true": {
			"true",
			true,
		},
		"false": {
			"false",
			false,
		},
	} {
		rootDir, err := os.MkdirTemp("", "vault-agent-jwt-auth-test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %s", err)
		}
		defer os.RemoveAll(rootDir)
		tokenPath := path.Join(rootDir, "token")
		err = os.WriteFile(tokenPath, []byte("test"), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		config := &auth.AuthConfig{
			Config: map[string]interface{}{
				"path": tokenPath,
				"role": "unusedrole",
			},
			Logger: hclog.Default(),
		}
		if tc.configValue != "" {
			config.Config["remove_jwt_after_reading"] = tc.configValue
		}

		jwtAuth, err := NewJWTAuthMethod(config)
		if err != nil {
			t.Fatal(err)
		}

		jwtAuth.(*jwtMethod).ingressToken()

		if _, err := os.Lstat(tokenPath); tc.shouldDelete {
			if err == nil || !os.IsNotExist(err) {
				t.Fatal(err)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestDeleteAfterReadingSymlink(t *testing.T) {
	for _, tc := range map[string]struct {
		configValue              string
		shouldDelete             bool
		removeJWTFollowsSymlinks bool
	}{
		"default": {
			"",
			true,
			false,
		},
		"explicit true": {
			"true",
			true,
			false,
		},
		"false": {
			"false",
			false,
			false,
		},
		"default + removeJWTFollowsSymlinks": {
			"",
			true,
			true,
		},
		"explicit true + removeJWTFollowsSymlinks": {
			"true",
			true,
			true,
		},
		"false + removeJWTFollowsSymlinks": {
			"false",
			false,
			true,
		},
	} {
		rootDir, err := os.MkdirTemp("", "vault-agent-jwt-auth-test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %s", err)
		}
		defer os.RemoveAll(rootDir)
		tokenPath := path.Join(rootDir, "token")
		err = os.WriteFile(tokenPath, []byte("test"), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		symlink, err := os.CreateTemp("", "auth.jwt.symlink.test.")
		if err != nil {
			t.Fatal(err)
		}
		symlinkName := symlink.Name()
		symlink.Close()
		os.Remove(symlinkName)
		os.Symlink(tokenPath, symlinkName)

		config := &auth.AuthConfig{
			Config: map[string]interface{}{
				"path": symlinkName,
				"role": "unusedrole",
			},
			Logger: hclog.Default(),
		}
		if tc.configValue != "" {
			config.Config["remove_jwt_after_reading"] = tc.configValue
		}
		config.Config["remove_jwt_follows_symlinks"] = tc.removeJWTFollowsSymlinks

		jwtAuth, err := NewJWTAuthMethod(config)
		if err != nil {
			t.Fatal(err)
		}

		jwtAuth.(*jwtMethod).ingressToken()

		pathToCheck := symlinkName
		if tc.removeJWTFollowsSymlinks {
			pathToCheck = tokenPath
		}
		if _, err := os.Lstat(pathToCheck); tc.shouldDelete {
			if err == nil || !os.IsNotExist(err) {
				t.Fatal(err)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}
