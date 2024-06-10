// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package revocation

import (
	"context"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/parsing"
	"github.com/hashicorp/vault/sdk/logical"
)

type UnifiedRevocationEntry struct {
	SerialNumber      string           `json:"-"`
	CertExpiration    time.Time        `json:"certificate_expiration_utc"`
	RevocationTimeUTC time.Time        `json:"revocation_time_utc"`
	CertificateIssuer issuing.IssuerID `json:"issuer_id"`
}

const (
	UnifiedRevocationReadPathPrefix  = "unified-revocation/"
	UnifiedRevocationWritePathPrefix = UnifiedRevocationReadPathPrefix + "{{clusterId}}/"
)

func WriteUnifiedRevocationEntry(ctx context.Context, storage logical.Storage, ure *UnifiedRevocationEntry) error {
	json, err := logical.StorageEntryJSON(UnifiedRevocationWritePathPrefix+parsing.NormalizeSerialForStorage(ure.SerialNumber), ure)
	if err != nil {
		return err
	}

	return storage.Put(ctx, json)
}
