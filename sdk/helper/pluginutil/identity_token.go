// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginutil

import "time"

type IdentityTokenRequest struct {
	Key      string
	Audience string
	TTL      time.Duration
}

type IdentityTokenResponse struct {
	Token string
	TTL   time.Duration
}
