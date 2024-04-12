// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "context"

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func (c *Core) StartCensusReports(ctx context.Context) {}
