// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package revocation

import (
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"time"
)

type RevocationInfo struct {
	CertificateBytes  []byte           `json:"certificate_bytes"`
	RevocationTime    int64            `json:"revocation_time"`
	RevocationTimeUTC time.Time        `json:"revocation_time_utc"`
	CertificateIssuer issuing.IssuerID `json:"issuer_id"`
}
