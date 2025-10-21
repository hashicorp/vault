// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package limits

// LimiterRegistry holds the map of RequestLimiters mapped to keys.
type LimiterRegistry struct{}
