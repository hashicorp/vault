package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/helper/wrapping"
	"github.com/hashicorp/vault/logical"
)

type dynamicSystemView struct {
	core       *Core
	mountEntry *MountEntry
}

func (d dynamicSystemView) DefaultLeaseTTL() time.Duration {
	def, _ := d.fetchTTLs()
	return def
}

func (d dynamicSystemView) MaxLeaseTTL() time.Duration {
	_, max := d.fetchTTLs()
	return max
}

func (d dynamicSystemView) SudoPrivilege(ctx context.Context, path string, token string) bool {
	// Resolve the token policy
	te, err := d.core.tokenStore.Lookup(ctx, token)
	if err != nil {
		d.core.logger.Error("core: failed to lookup token", "error", err)
		return false
	}

	// Ensure the token is valid
	if te == nil {
		d.core.logger.Error("entry not found for given token")
		return false
	}

	// Construct the corresponding ACL object
	acl, err := d.core.policyStore.ACL(ctx, te.Policies...)
	if err != nil {
		d.core.logger.Error("failed to retrieve ACL for token's policies", "token_policies", te.Policies, "error", err)
		return false
	}

	// The operation type isn't important here as this is run from a path the
	// user has already been given access to; we only care about whether they
	// have sudo
	req := new(logical.Request)
	req.Operation = logical.ReadOperation
	req.Path = path
	authResults := acl.AllowOperation(req)
	return authResults.RootPrivs
}

// TTLsByPath returns the default and max TTLs corresponding to a particular
// mount point, or the system default
func (d dynamicSystemView) fetchTTLs() (def, max time.Duration) {
	def = d.core.defaultLeaseTTL
	max = d.core.maxLeaseTTL

	if d.mountEntry.Config.DefaultLeaseTTL != 0 {
		def = d.mountEntry.Config.DefaultLeaseTTL
	}
	if d.mountEntry.Config.MaxLeaseTTL != 0 {
		max = d.mountEntry.Config.MaxLeaseTTL
	}

	return
}

// Tainted indicates that the mount is in the process of being removed
func (d dynamicSystemView) Tainted() bool {
	return d.mountEntry.Tainted
}

// CachingDisabled indicates whether to use caching behavior
func (d dynamicSystemView) CachingDisabled() bool {
	return d.core.cachingDisabled || (d.mountEntry != nil && d.mountEntry.Config.ForceNoCache)
}

func (d dynamicSystemView) LocalMount() bool {
	return d.mountEntry != nil && d.mountEntry.Local
}

// Checks if this is a primary Vault instance. Caller should hold the stateLock
// in read mode.
func (d dynamicSystemView) ReplicationState() consts.ReplicationState {
	return d.core.ReplicationState()
}

// ResponseWrapData wraps the given data in a cubbyhole and returns the
// token used to unwrap.
func (d dynamicSystemView) ResponseWrapData(ctx context.Context, data map[string]interface{}, ttl time.Duration, jwt bool) (*wrapping.ResponseWrapInfo, error) {
	req := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "sys/wrapping/wrap",
	}

	resp := &logical.Response{
		WrapInfo: &wrapping.ResponseWrapInfo{
			TTL: ttl,
		},
		Data: data,
	}

	if jwt {
		resp.WrapInfo.Format = "jwt"
	}

	_, err := d.core.wrapInCubbyhole(ctx, req, resp, nil)
	if err != nil {
		return nil, err
	}

	return resp.WrapInfo, nil
}

// LookupPlugin looks for a plugin with the given name in the plugin catalog. It
// returns a PluginRunner or an error if no plugin was found.
func (d dynamicSystemView) LookupPlugin(ctx context.Context, name string) (*pluginutil.PluginRunner, error) {
	if d.core == nil {
		return nil, fmt.Errorf("system view core is nil")
	}
	if d.core.pluginCatalog == nil {
		return nil, fmt.Errorf("system view core plugin catalog is nil")
	}
	r, err := d.core.pluginCatalog.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("{{err}}: %s", name), ErrPluginNotFound)
	}

	return r, nil
}

// MlockEnabled returns the configuration setting for enabling mlock on plugins.
func (d dynamicSystemView) MlockEnabled() bool {
	return d.core.enableMlock
}

func (d dynamicSystemView) CalculateTTL(increment, period, backendMaxTTL, explicitMaxTTL time.Duration, startTime time.Time) (ttl time.Duration, warnings []string, errors error) {
	now := time.Now()

	// Start off with the sys default value, and update according to period/TTL
	// from resp.Auth
	ttl = d.DefaultLeaseTTL()

	// Use the mount's configured max unless the backend specifies
	// something more restrictive (perhaps from a role configuration
	// parameter)
	maxTTL := d.MaxLeaseTTL()
	if backendMaxTTL > 0 && backendMaxTTL < maxTTL {
		maxTTL = backendMaxTTL
	}

	// Should never happen, but guard anyways
	if maxTTL < 0 {
		return 0, nil, fmt.Errorf("max TTL is negative")
	}

	switch {
	case period > 0:
		// Cap the period value to the sys max_ttl value. The auth backend should
		// have checked for it on its login path, but we check here again for
		// sanity.
		if period > maxTTL {
			period = maxTTL
			warnings = append(warnings,
				fmt.Sprintf("TTL of %q exceeded the effective max_ttl of %q; Period value is capped accordingly", period, maxTTL))
		}
		ttl = period
	case increment > 0:
		// We cannot go past this time
		maxValidTime := startTime.Add(maxTTL)

		// If we are past the max TTL, we shouldn't be in this function...but
		// fast path out if we are
		if maxValidTime.Before(now) {
			return 0, nil, fmt.Errorf("past the max TTL, cannot renew")
		}

		// We are proposing a time of the current time plus the increment
		proposedExpiration := now.Add(increment)

		// If the proposed expiration is after the maximum TTL of the lease,
		// cap the increment to whatever is left
		if maxValidTime.Before(proposedExpiration) {
			warnings = append(warnings,
				fmt.Sprintf("TTL of %q exceeded the effective max_ttl of %q; TTL value is capped accordingly", period, maxTTL))
			increment = maxValidTime.Sub(now)
		}
		ttl = increment
	}

	// Run some bounding checks if the explicit max TTL is set; we do not check
	// period as it's defined to escape the max TTL
	if explicitMaxTTL > 0 {
		// Limit the lease duration, except for periodic tokens -- in that case the explicit max limits the period, which itself can escape normal max
		if period == 0 && explicitMaxTTL > maxTTL {
			warnings = append(warnings,
				fmt.Sprintf("Explicit max TTL of %q is greater than system/mount allowed value; value is being capped to %q", explicitMaxTTL, maxTTL))
			explicitMaxTTL = maxTTL
		}

		// We cannot go past this time
		maxValidTime := startTime.Add(explicitMaxTTL)

		// If we are past the max TTL, we shouldn't be in this function...but
		// fast path out if we are
		if maxValidTime.Before(now) {
			return 0, nil, fmt.Errorf("past the explicit max TTL, cannot renew")
		}

		// We are proposing a time of the current time plus the increment
		proposedExpiration := now.Add(ttl)

		// If the proposed expiration is after the maximum TTL of the lease,
		// cap the increment to whatever is left
		if maxValidTime.Before(proposedExpiration) {
			ttl = maxValidTime.Sub(now)
			warnings = append(warnings,
				fmt.Sprintf("TTL of %q exceeded the explicit max_ttl of %q; TTL value is capped accordingly", ttl, explicitMaxTTL))
		}
	}

	return ttl, warnings, nil
}
