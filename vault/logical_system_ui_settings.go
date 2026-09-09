// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *SystemBackend) handleConfigUISettings(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	uiTelemetryEnabled := false

	if config := b.Core.GetCoreConfigInternal(); config != nil && config.UISettings != nil {
		uiTelemetryEnabled = config.UISettings.UITelemetry
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"ui_telemetry_enabled": uiTelemetryEnabled,
		},
	}, nil
}
