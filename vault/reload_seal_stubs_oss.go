// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"io"

	"github.com/hashicorp/go-hclog"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func (c *Core) reloadSealsEnt(secureRandomReader io.Reader, sealAccess Seal, logger hclog.Logger, shouldRewrap bool) error {
	return nil
}
