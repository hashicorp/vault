// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package limits

import (
	"context"
)

type RequestLimiter struct{}

// Acquire is a no-op on CE
func (l *RequestLimiter) Acquire(_ctx context.Context) (*RequestListener, bool) {
	return &RequestListener{}, true
}

// EstimatedLimit is effectively 0, since we're not limiting requests on CE.
func (l *RequestLimiter) EstimatedLimit() int { return 0 }
