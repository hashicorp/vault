// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package limits

import (
	"context"
	"errors"
)

var ErrCapacity = errors.New("Vault server temporarily overloaded")

type RequestLimiter struct{}

func (l *RequestLimiter) Acquire(_ctx context.Context) (*RequestListener, bool) {
	return &RequestListener{}, true
}

func (l *RequestLimiter) EstimatedLimit() int { return 0 }
