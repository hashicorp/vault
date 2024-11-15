// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloud

import (
	"time"

	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
)

type roleEntry struct {
	tokenutil.TokenParams

	ARN        *arn                          `json:"arn"`
	Policies   []string                      `json:"policies"`
	TTL        time.Duration                 `json:"ttl"`
	MaxTTL     time.Duration                 `json:"max_ttl"`
	Period     time.Duration                 `json:"period"`
	BoundCIDRs []*sockaddr.SockAddrMarshaler `json:"bound_cidrs"`
}

func (r *roleEntry) ToResponseData() map[string]interface{} {
	cidrs := make([]string, len(r.BoundCIDRs))
	for i, cidr := range r.BoundCIDRs {
		cidrs[i] = cidr.String()
	}
	d := map[string]interface{}{
		"arn": r.ARN.String(),
	}
	r.PopulateTokenData(d)

	if len(r.Policies) > 0 {
		d["policies"] = d["token_policies"]
	}
	if len(r.BoundCIDRs) > 0 {
		d["bound_cidrs"] = d["token_bound_cidrs"]
	}
	if r.TTL > 0 {
		d["ttl"] = int64(r.TTL.Seconds())
	}
	if r.MaxTTL > 0 {
		d["max_ttl"] = int64(r.MaxTTL.Seconds())
	}
	if r.Period > 0 {
		d["period"] = int64(r.Period.Seconds())
	}

	return d
}
