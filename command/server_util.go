// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	"github.com/hashicorp/cli"
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
		licenseReloadedCh: make(chan error),
	}
}
