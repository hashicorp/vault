package alicloud

import (
	"time"

	"github.com/hashicorp/go-sockaddr"
)

type roleEntry struct {
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
	return map[string]interface{}{
		"arn":         r.ARN.String(),
		"policies":    r.Policies,
		"ttl":         r.TTL / time.Second,
		"max_ttl":     r.MaxTTL / time.Second,
		"period":      r.Period / time.Second,
		"bound_cidrs": cidrs,
	}
}
