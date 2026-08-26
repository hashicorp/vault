// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"context"
	"time"

	consumptionv1 "github.com/hashicorp/calistoga-control-plane/sdks/secure-products/scp/api/vault/consumption/v1"
)

func (c *Core) buildMetricsProto(ctx context.Context, currentMonth time.Time) (*consumptionv1.ConsumptionMetrics, error) {
	return nil, nil
}
