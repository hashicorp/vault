// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package connutil

import (
	"context"
	"database/sql"
	"errors"
)

func (c *SQLConnectionProducer) StaticConnection(_ context.Context, _, _ string) (*sql.DB, error) {
	return nil, errors.New("self-managed static roles only available in Vault Enterprise")
}
