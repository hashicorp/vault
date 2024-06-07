// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki_backend

import "github.com/hashicorp/vault/builtin/logical/pki/revocation"

type CrlBuilderType interface {
	RebuildIfForced(sc StorageContext) ([]string, error)
	Rebuild(sc StorageContext, forceNew bool) ([]string, error)
	RebuildDeltaCRLsHoldingLock(sc StorageContext, forceNew bool) ([]string, error)
	GetPresentLocalDeltaWALForClearing(sc StorageContext) ([]string, error)
	GetPresentUnifiedDeltaWALForClearing(sc StorageContext) ([]string, error)
	GetConfigWithUpdate(sc StorageContext) (*revocation.CrlConfig, error)
	ClearLocalDeltaWAL(sc StorageContext, walSerials []string) error
	ClearUnifiedDeltaWAL(sc StorageContext, walSerials []string) error
}
