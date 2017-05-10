package logical

import (
	"errors"
	"time"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/helper/wrapping"
)

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

	// ReplicationState indicates the state of cluster replication
	ReplicationState() consts.ReplicationState

	// ResponseWrapData wraps the given data in a cubbyhole and returns the
	// token used to unwrap.
	ResponseWrapData(data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error)

	// LookupPlugin looks into the plugin catalog for a plugin with the given
	// name. Returns a PluginRunner or an error if a plugin can not be found.
	LookupPlugin(string) (*pluginutil.PluginRunner, error)

	// MlockEnabled returns the configuration setting for enabling mlock on
	// plugins.
	MlockEnabled() bool
}

type StaticSystemView struct {
	DefaultLeaseTTLVal  time.Duration
	MaxLeaseTTLVal      time.Duration
	SudoPrivilegeVal    bool
	TaintedVal          bool
	CachingDisabledVal  bool
	Primary             bool
	EnableMlock         bool
	ReplicationStateVal consts.ReplicationState
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

func (d StaticSystemView) ReplicationState() consts.ReplicationState {
	return d.ReplicationStateVal
}

func (d StaticSystemView) ResponseWrapData(data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error) {
	return nil, errors.New("ResponseWrapData is not implemented in StaticSystemView")
}

func (d StaticSystemView) LookupPlugin(name string) (*pluginutil.PluginRunner, error) {
	return nil, errors.New("LookupPlugin is not implemented in StaticSystemView")
}

func (d StaticSystemView) MlockEnabled() bool {
	return d.EnableMlock
}
