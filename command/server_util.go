// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"os"
	"path/filepath"
	"runtime/pprof"
	"testing"

	"github.com/hashicorp/cli"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	physInmem "github.com/hashicorp/vault/sdk/physical/inmem"
)

func TestServerCommand(tb testing.TB) (*cli.MockUi, *ServerCommand) {
	tb.Helper()
	return testServerCommand(tb)
}

func (c *ServerCommand) StartedCh() chan struct{} {
	return c.startedCh
}

func (c *ServerCommand) ReloadedCh() chan struct{} {
	return c.reloadedCh
}

func (c *ServerCommand) LicenseReloadedCh() chan error {
	return c.licenseReloadedCh
}

func testServerCommand(tb testing.TB) (*cli.MockUi, *ServerCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &ServerCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
		ShutdownCh: MakeShutdownCh(),
		SighupCh:   MakeSighupCh(),
		SigUSR2Ch:  MakeSigUSR2Ch(),
		PhysicalBackends: map[string]physical.Factory{
			"inmem":               physInmem.NewInmem,
			"inmem_ha":            physInmem.NewInmemHA,
			"inmem_transactional": physInmem.NewTransactionalInmem,
		},

		// These prevent us from random sleep guessing...
		startedCh:         make(chan struct{}, 5),
		reloadedCh:        make(chan struct{}, 5),
		licenseReloadedCh: make(chan error, 1),
	}
}

// writeGoroutineDump writes a goroutine pprof profile to a file named
// "goroutine" inside a resolved directory. The directory is determined by the
// VAULT_STACKTRACE_FILE_PATH env var when set and valid; otherwise a new
// timestamped subdirectory is created inside os.TempDir(). Returns the
// absolute path of the written file on success, or an empty string on failure.
func writeGoroutineDump(logger hclog.Logger) string {
	dir := os.Getenv("VAULT_STACKTRACE_FILE_PATH")
	if dir != "" {
		if _, err := os.Stat(dir); err != nil {
			logger.Warn("writeGoroutineDump: could not access VAULT_STACKTRACE_FILE_PATH", "path", dir, "error", err)
			return ""
		}
	} else {
		var err error
		dir, err = os.MkdirTemp(os.TempDir(), "vault-goroutine-dump-")
		if err != nil {
			logger.Warn("writeGoroutineDump: could not create temporary directory", "error", err)
			return ""
		}
	}

	f, err := os.Create(filepath.Join(dir, "goroutine"))
	if err != nil {
		logger.Warn("writeGoroutineDump: could not create goroutine dump file", "error", err)
		return ""
	}
	defer f.Close()

	if err := pprof.Lookup("goroutine").WriteTo(f, 2); err != nil {
		logger.Warn("writeGoroutineDump: could not write goroutine profile", "error", err)
		return ""
	}

	return f.Name()
}
