package logical

import "time"

// SystemView exposes system configuration information in a safe way
// for logical backends to consume
type SystemView interface {
	// DefaultLeaseTTL returns the default lease TTL set in Vault configuration
	DefaultLeaseTTL() time.Duration

	// MaxLeaseTTL returns the max lease TTL set in Vault configuration; backend
	// authors should take care not to issue credentials that last longer than
	// this value, as Vault will revoke them
	MaxLeaseTTL() time.Duration

	// SudoPrivilege returns true if given path has sudo privileges
	// for the given client token
	SudoPrivilege(path string, token string) bool

	// Returns true if the mount is tainted. A mount is tainted if it is in the
	// process of being unmounted. This should only be used in special
	// circumstances; a primary use-case is as a guard in revocation functions.
	// If revocation of a backend's leases fails it can keep the unmounting
	// process from being successful. If the reason for this failure is not
	// relevant when the mount is tainted (for instance, saving a CRL to disk
	// when the stored CRL will be removed during the unmounting process
	// anyways), we can ignore the errors to allow unmounting to complete.
	Tainted() bool

	// Returns true if caching is disabled. If true, no caches should be used,
	// despite known slowdowns.
	CachingDisabled() bool
}

type StaticSystemView struct {
	DefaultLeaseTTLVal time.Duration
	MaxLeaseTTLVal     time.Duration
	SudoPrivilegeVal   bool
	TaintedVal         bool
	CachingDisabledVal bool
}

func (d StaticSystemView) DefaultLeaseTTL() time.Duration {
	return d.DefaultLeaseTTLVal
}

func (d StaticSystemView) MaxLeaseTTL() time.Duration {
	return d.MaxLeaseTTLVal
}

func (d StaticSystemView) SudoPrivilege(path string, token string) bool {
	return d.SudoPrivilegeVal
}

func (d StaticSystemView) Tainted() bool {
	return d.TaintedVal
}

func (d StaticSystemView) CachingDisabled() bool {
	return d.CachingDisabledVal
}
