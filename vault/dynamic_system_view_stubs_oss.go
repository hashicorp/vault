// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// DownloadExtractVerifyPlugin returns an error as this is an enterprise only feature
func (d dynamicSystemView) DownloadExtractVerifyPlugin(_ context.Context, _ *pluginutil.PluginRunner) error {
	return fmt.Errorf("enterprise only feature")
}

func (d dynamicSystemView) TPMByID(ctx context.Context, id string) (*logical.TPM, error) {
	return nil, fmt.Errorf("enterprise only feature")
}

func (d dynamicSystemView) GroupsForTPM(ctx context.Context, id string) ([]*logical.TPMGroup, error) {
	return nil, fmt.Errorf("enterprise only feature")
}
