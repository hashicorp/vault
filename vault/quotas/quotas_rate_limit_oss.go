// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package quotas

import "context"

type entRateLimitRequest struct{}

type entRateLimitQuota struct{}

// getGroupKey returns the identifier to the request for rate limiting purposes. On CE we only support IP-based grouping.
func (rlq *RateLimitQuota) getGroupKey(req *Request) (key string, isSecondaryGroup bool, err error) {
	if rlq.GroupBy != "" && rlq.GroupBy != GroupByIp {
		return "", false, ErrGroupByNotSupported
	}
	return req.ClientAddress, false, nil
}

func (rlq *RateLimitQuota) take(ctx context.Context, key string, isSecondaryGroup bool) (tokens, remaining, reset uint64, allow bool, err error) {
	return rlq.store.Take(ctx, key)
}

func (rlq *RateLimitQuota) initializeEnt() error {
	return nil
}

func (rlq *RateLimitQuota) closeEnt(ctx context.Context) error {
	return nil
}
