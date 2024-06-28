// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/helper/pluginidentityutil"
	"github.com/hashicorp/vault/sdk/logical"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func (i *IdentityStore) generatePluginIdentityToken(_ context.Context, _ logical.Storage, _ *MountEntry, _ string, _ time.Duration) (string, time.Duration, error) {
	return "", 0, pluginidentityutil.ErrPluginWorkloadIdentityUnsupported
}

func validChildIssuer(child string) bool {
	return child == baseIdentityTokenIssuer
}
