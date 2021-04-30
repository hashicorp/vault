package jwt

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestIngressToken(t *testing.T) {
	const (
		dir       = "dir"
		file      = "file"
		empty     = "empty"
		missing   = "missing"
		symlinked = "symlinked"
	)

	rootDir, err := ioutil.TempDir("", "vault-agent-jwt-auth-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
	}
	defer os.RemoveAll(rootDir)

	setupTestDir := func() string {
		testDir, err := ioutil.TempDir(rootDir, "")
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(path.Join(testDir, file), []byte("test"), 0644)
		if err != nil {
			t.Fatal(err)
		}
		_, err = os.Create(path.Join(testDir, empty))
		if err != nil {
			t.Fatal(err)
		}
		err = os.Mkdir(path.Join(testDir, dir), 0755)
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
				t.Fatal("logs contained unexepected error", logBuffer.String())
			}
		}
	}
}
