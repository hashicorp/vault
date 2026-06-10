// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
)

// MountMetrics holds all metrics collected from a single mount table traversal.
// It distinguishes between replicated and local mounts.
type MountMetrics struct {
	// ReplicatedRoleCounts contains role counts from replicated (non-local) mounts
	ReplicatedRoleCounts *RoleCounts
	// ReplicatedManagedKeys contains managed key counts from replicated mounts
	ReplicatedManagedKeys *ManagedKeyCounts
	// ReplicatedKvMounts contains KV counts from replicated mounts
	ReplicatedKvCounts int
	// LocalRoleCounts contains role counts from local mounts
	LocalRoleCounts *RoleCounts
	// LocalManagedKeys contains managed key counts from local mounts
	LocalManagedKeys *ManagedKeyCounts
	// LocalKvMounts contains KV counts from local mounts
	LocalKvCounts int
}

// CountMetricsFromMounts performs a single iteration through mount entries
// and collects all required metrics for both replicated and local mounts.
func (c *Core) CountMetricsFromMounts(officialOnly bool) (*MountMetrics, error) {
	if c.Sealed() {
		return nil, fmt.Errorf("core is sealed, cannot access mount table")
	}

	c.mountsLock.RLock()
	defer c.mountsLock.RUnlock()

	metrics := &MountMetrics{
		ReplicatedRoleCounts:  &RoleCounts{},
		ReplicatedManagedKeys: &ManagedKeyCounts{},
		ReplicatedKvCounts:    0,
		LocalRoleCounts:       &RoleCounts{},
		LocalManagedKeys:      &ManagedKeyCounts{},
		LocalKvCounts:         0,
	}

	if c.mounts == nil {
		return metrics, nil
	}

	ctx := namespace.RootContext(c.activeContext)

	// Iterate through all mount entries
	for _, entry := range c.mounts.Entries {
		if officialOnly && !c.isOfficialPlugin(ctx, entry) {
			continue
		}

		// Determine which metrics buckets to update based on mount locality
		var targetRoleCounts *RoleCounts
		var targetManagedKeys *ManagedKeyCounts
		targetKvLocal := false

		if entry.Local {
			targetRoleCounts = metrics.LocalRoleCounts
			targetManagedKeys = metrics.LocalManagedKeys
			targetKvLocal = true
		} else {
			targetRoleCounts = metrics.ReplicatedRoleCounts
			targetManagedKeys = metrics.ReplicatedManagedKeys
		}

		// Collect role counts and managed key counts based on plugin type
		c.collectMountMetrics(ctx, entry, targetRoleCounts, targetManagedKeys)

		// If this is a KV mount, gather its secret count and add to the total
		if kvMount := getKVMountMetadata(entry); kvMount != nil {
			c.walkKvMountSecrets(ctx, kvMount)
			if targetKvLocal {
				metrics.LocalKvCounts += kvMount.NumSecrets
			} else {
				metrics.ReplicatedKvCounts += kvMount.NumSecrets
			}
		}
	}

	return metrics, nil
}

// collectMountMetrics collects role counts and managed key counts for a specific mount entry
// based on its plugin type. It updates the provided target metrics in place.
func (c *Core) collectMountMetrics(ctx context.Context, entry *MountEntry, targetRoleCounts *RoleCounts, targetManagedKeys *ManagedKeyCounts) {
	apiList := func(entry *MountEntry, apiPath string) []string {
		listRequest := &logical.Request{
			Operation: logical.ListOperation,
			Path:      entry.namespace.Path + entry.Path + apiPath,
		}

		resp, err := c.router.Route(ctx, listRequest)
		if err != nil || resp == nil || resp.Data == nil {
			return nil
		}

		rawKeys, ok := resp.Data["keys"]
		if !ok || rawKeys == nil {
			return nil
		}

		// Type switch handles both the 'official' behavior and the 'generic' behavior
		switch kt := rawKeys.(type) {
		case []string:
			// Existing plugins likely hit this path
			return kt
		case []interface{}:
			// External/RPC plugins likely hit this path
			keys := make([]string, 0, len(kt))
			for _, k := range kt {
				if s, ok := k.(string); ok {
					keys = append(keys, s)
				}
			}
			return keys
		default:
			// If it's something totally weird, we still fail safely
			return nil
		}
	}

	pluginName := getAdjustedPluginType(entry)
	if pluginName == "" {
		return
	}

	switch pluginName {
	case pluginconsts.SecretEngineAWS:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.AWSDynamicRoles += len(dynamicRoles)
		staticRoles := apiList(entry, "static-roles")
		targetRoleCounts.AWSStaticRoles += len(staticRoles)

	case pluginconsts.SecretEngineAzure:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.AzureDynamicRoles += len(dynamicRoles)
		staticRoles := apiList(entry, "static-roles")
		targetRoleCounts.AzureStaticRoles += len(staticRoles)

	case pluginconsts.SecretEngineDatabase:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.DatabaseDynamicRoles += len(dynamicRoles)
		staticRoles := apiList(entry, "static-roles")
		targetRoleCounts.DatabaseStaticRoles += len(staticRoles)

	case pluginconsts.SecretEngineGCP:
		rolesets := apiList(entry, "rolesets")
		targetRoleCounts.GCPRolesets += len(rolesets)
		staticAccounts := apiList(entry, "static-accounts")
		targetRoleCounts.GCPStaticAccounts += len(staticAccounts)
		impersonatedAccounts := apiList(entry, "impersonated-accounts")
		targetRoleCounts.GCPImpersonatedAccounts += len(impersonatedAccounts)

	case pluginconsts.SecretEngineLDAP:
		dynamicRoles := apiList(entry, "role")
		targetRoleCounts.LDAPDynamicRoles += len(dynamicRoles)
		staticRoles := apiList(entry, "static-role")
		targetRoleCounts.LDAPStaticRoles += len(staticRoles)

	case pluginconsts.SecretEngineOpenLDAP:
		dynamicRoles := apiList(entry, "role")
		targetRoleCounts.OpenLDAPDynamicRoles += len(dynamicRoles)
		staticRoles := apiList(entry, "static-role")
		targetRoleCounts.OpenLDAPStaticRoles += len(staticRoles)

	case pluginconsts.SecretEngineAlicloud:
		dynamicRoles := apiList(entry, "role")
		targetRoleCounts.AlicloudDynamicRoles += len(dynamicRoles)

	case pluginconsts.SecretEngineRabbitMQ:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.RabbitMQDynamicRoles += len(dynamicRoles)

	case pluginconsts.SecretEngineConsul:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.ConsulDynamicRoles += len(dynamicRoles)

	case pluginconsts.SecretEngineNomad:
		dynamicRoles := apiList(entry, "role")
		targetRoleCounts.NomadDynamicRoles += len(dynamicRoles)

	case pluginconsts.SecretEngineKubernetes:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.KubernetesDynamicRoles += len(dynamicRoles)

	case pluginconsts.SecretEngineMongoDBAtlas:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.MongoDBAtlasDynamicRoles += len(dynamicRoles)

	case pluginconsts.SecretEngineTerraform:
		dynamicRoles := apiList(entry, "role")
		targetRoleCounts.TerraformCloudDynamicRoles += len(dynamicRoles)

	case pluginconsts.SecretEngineTOTP:
		keyCountPerEntry := apiList(entry, "keys")
		targetManagedKeys.TotpKeys += len(keyCountPerEntry)

	case pluginconsts.SecretEngineKeymgmt:
		keyCountPerEntry := apiList(entry, "key")
		targetManagedKeys.KmseKeys += len(keyCountPerEntry)

	case pluginconsts.SecretEngineOS:
		// OS plugin stores all accounts within each host entry
		// List all hosts, then list accounts for each host
		hosts := apiList(entry, "hosts/")
		accountCount := 0
		for _, host := range hosts {
			if host == "" {
				continue
			}
			hostName := strings.TrimSuffix(host, "/")
			accounts := apiList(entry, "hosts/"+hostName+"/accounts/")
			accountCount += len(accounts)
		}
		targetRoleCounts.OSLocalAccountRoles += accountCount
	}
}

// getKVMountMetadata creates and returns KV mount metadata if the entry is a KV mount.
// Returns nil if the entry is not a KV mount. Includes all KV versions.
func getKVMountMetadata(entry *MountEntry) *kvMount {
	// Check if this is a KV mount (type kv or generic)
	if entry.Type != pluginconsts.SecretEngineKV && entry.Type != pluginconsts.SecretEngineGeneric {
		return nil
	}

	// Determine KV version from mount options
	version, ok := entry.Options["version"]
	if !ok || version == "" {
		version = "1"
	}

	// Create KV mount metadata entry
	return &kvMount{
		Namespace:            entry.namespace,
		MountPoint:           entry.Path,
		MountAccessor:        entry.Accessor,
		Version:              version,
		NumSecrets:           0,
		Local:                entry.Local,
		RunningPluginVersion: entry.RunningVersion,
	}
}

// isOfficialPlugin is a helper to determine if a mount entry is an official or builtin plugin
func (c *Core) isOfficialPlugin(ctx context.Context, entry *MountEntry) bool {
	pluginName := getAdjustedPluginType(entry)
	if pluginName == "" {
		return false
	}

	pluginVersion := entry.RunningVersion

	runner, err := c.pluginCatalog.Get(ctx, pluginName, consts.PluginTypeSecrets, pluginVersion)
	if err != nil {
		return false
	}

	return isOfficialOrBuiltin(runner)
}

// GetRoleCounts returns the combined local and replicated role counts across all mounts
// For use in tests only
func (c *Core) GetRoleCounts() *RoleCounts {
	m, err := c.CountMetricsFromMounts(false)
	if err != nil {
		return nil
	}
	return combineRoleCounts(m.LocalRoleCounts, m.ReplicatedRoleCounts)
}

// GetRoleAndManagedKeyCounts returns the local or replicated role and managed key counts depending on the mount type
// For use in tests only
func (c *Core) GetRoleAndManagedKeyCounts(mountType string) (*RoleCounts, *ManagedKeyCounts) {
	m, err := c.CountMetricsFromMounts(true)
	if err != nil {
		return nil, nil
	}
	if mountType == billing.LocalPrefix {
		return m.LocalRoleCounts, m.LocalManagedKeys
	}
	return m.ReplicatedRoleCounts, m.ReplicatedManagedKeys
}

// GetKvCounts returns the local or replicated KV counts depending on the mount type
// For use in tests only
func (c *Core) GetKvCounts(mountType string) int {
	m, err := c.CountMetricsFromMounts(true)
	if err != nil {
		return 0
	}
	if mountType == billing.LocalPrefix {
		return m.LocalKvCounts
	}
	return m.ReplicatedKvCounts
}
