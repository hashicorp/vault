// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"io"

	"github.com/hashicorp/go-hclog"
)

func (c *Core) reloadSealsEnt(secureRandomReader io.Reader, sealAccess Seal, logger hclog.Logger, shouldRewrap bool) error {
	return nil
}
