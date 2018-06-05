package consts

import "errors"

var (
	// ErrSealed is returned if an operation is performed on a sealed barrier.
	// No operation is expected to succeed before unsealing
	ErrSealed = errors.New("Vault is sealed")

	// ErrStandby is returned if an operation is performed on a standby Vault.
	// No operation is expected to succeed until active.
	ErrStandby = errors.New("Vault is in standby mode")

	// ErrVerIncompatible is returned if an operation is performed from a server version
	// that is less than the data version stored in coreDataVersionPath
	ErrVerIncompatible = errors.New("Vault data version is incompatible with the server version")

	// Used when .. is used in a path
	ErrPathContainsParentReferences = errors.New("path cannot contain parent references")
)
