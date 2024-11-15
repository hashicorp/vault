// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"time"

	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
)

// RoleEntry is a role as it's reflected in Vault's storage system.
type RoleEntry struct {
	tokenutil.TokenParams

	BoundAppIDs       []string `json:"bound_application_ids"`
	BoundSpaceIDs     []string `json:"bound_space_ids"`
	BoundOrgIDs       []string `json:"bound_organization_ids"`
	BoundInstanceIDs  []string `json:"bound_instance_ids"`
	DisableIPMatching bool     `json:"disable_ip_matching"`

	// Deprecated by TokenParams
	TTL        time.Duration                 `json:"ttl"`
	MaxTTL     time.Duration                 `json:"max_ttl"`
	Period     time.Duration                 `json:"period"`
	Policies   []string                      `json:"policies"`
	BoundCIDRs []*sockaddr.SockAddrMarshaler `json:"bound_cidrs"`
}
