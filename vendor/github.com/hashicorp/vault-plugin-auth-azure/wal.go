// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azureauth

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

const (
	walRotateRootCreds = "rotateRootCreds"
)

func (b *azureAuthBackend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
	switch kind {
	case walRotateRootCreds:
		return b.rollbackRootWAL(ctx, req, data)
	default:
		return fmt.Errorf("unknown rollback type %q", kind)
	}
}

type walRotateRoot struct{}

func (b *azureAuthBackend) rollbackRootWAL(ctx context.Context, req *logical.Request, data interface{}) error {
	b.Logger().Debug("rolling back config")
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return err
	}

	config.NewClientSecret = ""
	config.NewClientSecretCreated = time.Time{}
	config.NewClientSecretExpirationDate = time.Time{}
	config.NewClientSecretKeyID = ""

	err = b.saveConfig(ctx, config, req.Storage)
	if err != nil {
		return err
	}

	b.updatePassword = false

	return nil
}
