// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"

	"github.com/hashicorp/vault/vault/seal"
)

// SealAccess is a wrapper around Seal that exposes accessor methods
// through Core.SealAccess() while restricting the ability to modify
// Core.seal itself.
type SealAccess struct {
	seal Seal
}

func NewSealAccess(seal Seal) *SealAccess {
	return &SealAccess{seal: seal}
}

func (s *SealAccess) StoredKeysSupported() seal.StoredKeysSupport {
	return s.seal.StoredKeysSupported()
}

func (s *SealAccess) BarrierSealConfigType() SealConfigType {
	return s.seal.BarrierSealConfigType()
}

func (s *SealAccess) BarrierConfig(ctx context.Context) (*SealConfig, error) {
	return s.seal.BarrierConfig(ctx)
}

func (s *SealAccess) RecoveryKeySupported() bool {
	return s.seal.RecoveryKeySupported()
}

func (s *SealAccess) RecoveryConfig(ctx context.Context) (*SealConfig, error) {
	return s.seal.RecoveryConfig(ctx)
}

func (s *SealAccess) VerifyRecoveryKey(ctx context.Context, key []byte) error {
	return s.seal.VerifyRecoveryKey(ctx, key)
}

func (s *SealAccess) ClearCaches(ctx context.Context) {
	s.seal.ClearBarrierConfig(ctx)
	if s.RecoveryKeySupported() {
		s.seal.ClearRecoveryConfig(ctx)
	}
}

func (s *SealAccess) GetAccess() seal.Access {
	return s.seal.GetAccess()
}
