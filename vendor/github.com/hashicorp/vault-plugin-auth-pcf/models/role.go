package models

import (
	"time"

	"github.com/hashicorp/go-sockaddr"
)

// RoleEntry is a role as it's reflected in Vault's storage system.
type RoleEntry struct {
	BoundAppIDs       []string                      `json:"bound_application_ids"`
	BoundSpaceIDs     []string                      `json:"bound_space_ids"`
	BoundOrgIDs       []string                      `json:"bound_organization_ids"`
	BoundInstanceIDs  []string                      `json:"bound_instance_ids"`
	BoundCIDRs        []*sockaddr.SockAddrMarshaler `json:"bound_cidrs"`
	Policies          []string                      `json:"policies"`
	DisableIPMatching bool                          `json:"disable_ip_matching"`
	TTL               time.Duration                 `json:"ttl"`
	MaxTTL            time.Duration                 `json:"max_ttl"`
	Period            time.Duration                 `json:"period"`
}
