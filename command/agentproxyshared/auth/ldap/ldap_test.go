// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"bytes"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
)

func TestIngressPass(t *testing.T) {
	const (
		dir       = "dir"
		file      = "file"
		empty     = "empty"
		missing   = "missing"
		symlinked = "symlinked"
	)

	rootDir, err := os.MkdirTemp("", "vault-agent-ldap-auth-test")
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
			"[ERROR] password file is not a regular file or symlink",
		},
		{
			"password file path is symlink",
			symlinked,
			"",
		},
		{
			"password file path is missing (implies nothing for ingressPass to do)",
			missing,
			"",
		},
		{
			"password file path is empty file",
			empty,
			"[WARN]  empty password file read",
		},
	} {
		testDir := setupTestDir()
		logBuffer := bytes.Buffer{}
		ldapAuth := &ldapMethod{
			logger: hclog.New(&hclog.LoggerOptions{
				Output: &logBuffer,
			}),
			latestPass:       new(atomic.Value),
			passwordFilePath: path.Join(testDir, tc.path),
		}

		ldapAuth.ingressPass()

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
		rootDir, err := os.MkdirTemp("", "vault-agent-ldap-auth-test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %s", err)
		}
		defer os.RemoveAll(rootDir)
		passPath := path.Join(rootDir, "pass")
		err = os.WriteFile(passPath, []byte("test"), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		config := &auth.AuthConfig{
			Config: map[string]interface{}{
				"password_file_path": passPath,
				"username":           "testuser",
			},
			Logger: hclog.Default(),
		}
		if tc.configValue != "" {
			config.Config["remove_password_after_reading"] = tc.configValue
		}

		ldapAuth, err := NewLdapAuthMethod(config)
		if err != nil {
			t.Fatal(err)
		}

		ldapAuth.(*ldapMethod).ingressPass()

		if _, err := os.Lstat(passPath); tc.shouldDelete {
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
		configValue               string
		shouldDelete              bool
		removePassFollowsSymlinks bool
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
		"default + removePassFollowsSymlinks": {
			"",
			true,
			true,
		},
		"explicit true + removePassFollowsSymlinks": {
			"true",
			true,
			true,
		},
		"false + removePassFollowsSymlinks": {
			"false",
			false,
			true,
		},
	} {
		rootDir, err := os.MkdirTemp("", "vault-agent-ldap-auth-test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %s", err)
		}
		defer os.RemoveAll(rootDir)
		passPath := path.Join(rootDir, "pass")
		err = os.WriteFile(passPath, []byte("test"), 0o644)
		if err != nil {
			t.Fatal(err)
		}

		symlink, err := os.CreateTemp("", "auth.ldap.symlink.test.")
		if err != nil {
			t.Fatal(err)
		}
		symlinkName := symlink.Name()
		symlink.Close()
		os.Remove(symlinkName)
		os.Symlink(passPath, symlinkName)

		config := &auth.AuthConfig{
			Config: map[string]interface{}{
				"password_file_path": symlinkName,
				"username":           "testuser",
			},
			Logger: hclog.Default(),
		}
		if tc.configValue != "" {
			config.Config["remove_password_after_reading"] = tc.configValue
		}
		config.Config["remove_password_follows_symlinks"] = tc.removePassFollowsSymlinks

		ldapAuth, err := NewLdapAuthMethod(config)
		if err != nil {
			t.Fatal(err)
		}

		ldapAuth.(*ldapMethod).ingressPass()

		pathToCheck := symlinkName
		if tc.removePassFollowsSymlinks {
			pathToCheck = passPath
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
