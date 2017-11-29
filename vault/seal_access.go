package vault

// SealAccess is a wrapper around Seal that exposes accessor methods
// through Core.SealAccess() while restricting the ability to modify
// Core.seal itself.
type SealAccess struct {
	seal Seal
}

func NewSealAccess(seal Seal) *SealAccess {
	return &SealAccess{seal: seal}
}

func (s *SealAccess) StoredKeysSupported() bool {
	return s.seal.StoredKeysSupported()
}

func (s *SealAccess) BarrierConfig() (*SealConfig, error) {
	return s.seal.BarrierConfig()
}

func (s *SealAccess) RecoveryKeySupported() bool {
	return s.seal.RecoveryKeySupported()
}

func (s *SealAccess) RecoveryConfig() (*SealConfig, error) {
	return s.seal.RecoveryConfig()
}

func (s *SealAccess) VerifyRecoveryKey(key []byte) error {
	return s.seal.VerifyRecoveryKey(key)
}

func (s *SealAccess) ClearCaches() {
	s.seal.SetBarrierConfig(nil)
	if s.RecoveryKeySupported() {
		s.seal.SetRecoveryConfig(nil)
	}
}
