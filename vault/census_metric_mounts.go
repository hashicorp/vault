// Copyright IBM Corp. 2016, 2026
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
	// ReplicatedSecretEngineResourceCounts contains secret engine resource counts from replicated mounts
	ReplicatedSecretEngineResourceCounts *SecretEngineResourceCounts
	// ReplicatedKvMounts contains KV counts from replicated mounts
	ReplicatedKvCounts int
	// LocalRoleCounts contains role counts from local mounts
	LocalRoleCounts *RoleCounts
	// LocalManagedKeys contains managed key counts from local mounts
	LocalManagedKeys *ManagedKeyCounts
	// LocalSecretEngineResourceCounts contains secret engine resource counts from local mounts
	LocalSecretEngineResourceCounts *SecretEngineResourceCounts
	// LocalKvMounts contains KV counts from local mounts
	LocalKvCounts int

	// Attribution fields
	//
	// ReplicatedKvAttribution contains attribution data for replicated mounts with KV secrets
	ReplicatedKvAttribution MountAttributionMap
	// ReplicatedRoleAttribution contains a mapping from role type to attribution data for replicated mounts with those roles
	ReplicatedRoleAttribution map[string]MountAttributionMap
	// ReplicatedManagedKeyAttribution contains a mapping from managed key type to attribution data for replicated mounts with those keys
	ReplicatedManagedKeyAttribution map[string]MountAttributionMap
	// LocalKvAttribution contains attribution data for local mounts with KV secrets
	LocalKvAttribution MountAttributionMap
	// LocalRoleAttribution contains a mapping from role type to attribution data for local mounts with those roles
	LocalRoleAttribution map[string]MountAttributionMap
	// LocalManagedKeyAttribution contains a mapping from managed key type to attribution data for local mounts with those keys
	LocalManagedKeyAttribution map[string]MountAttributionMap
}

// MountAttributionMap is a map from mount_accessor to mount metadata for attribution
type MountAttributionMap map[string]logical.MountAttribution

// CountMetricsSecretMounts performs a single iteration through mount entries
// and collects all required metrics for both replicated and local mounts.
func (c *Core) CountMetricsSecretMounts(officialOnly bool, collectAttribution bool) (*MountMetrics, error) {
	if c.Sealed() {
		return nil, fmt.Errorf("core is sealed, cannot access mount table")
	}

	c.mountsLock.RLock()
	defer c.mountsLock.RUnlock()

	metrics := &MountMetrics{
		ReplicatedRoleCounts:                 &RoleCounts{},
		ReplicatedManagedKeys:                &ManagedKeyCounts{},
		ReplicatedSecretEngineResourceCounts: &SecretEngineResourceCounts{},
		ReplicatedKvCounts:                   0,
		LocalRoleCounts:                      &RoleCounts{},
		LocalManagedKeys:                     &ManagedKeyCounts{},
		LocalSecretEngineResourceCounts:      &SecretEngineResourceCounts{},
		LocalKvCounts:                        0,
		ReplicatedKvAttribution:              make(MountAttributionMap),
		ReplicatedRoleAttribution:            make(map[string]MountAttributionMap),
		ReplicatedManagedKeyAttribution:      make(map[string]MountAttributionMap),
		LocalKvAttribution:                   make(MountAttributionMap),
		LocalRoleAttribution:                 make(map[string]MountAttributionMap),
		LocalManagedKeyAttribution:           make(map[string]MountAttributionMap),
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
		var targetSecretEngineResourceCounts *SecretEngineResourceCounts
		var targetRoleAttribution map[string]MountAttributionMap
		var targetManagedKeyAttribution map[string]MountAttributionMap
		var targetKvAttribution MountAttributionMap
		targetKvLocal := false

		if entry.Local {
			targetRoleCounts = metrics.LocalRoleCounts
			targetManagedKeys = metrics.LocalManagedKeys
			targetSecretEngineResourceCounts = metrics.LocalSecretEngineResourceCounts
			targetRoleAttribution = metrics.LocalRoleAttribution
			targetManagedKeyAttribution = metrics.LocalManagedKeyAttribution
			targetKvAttribution = metrics.LocalKvAttribution
			targetKvLocal = true
		} else {
			targetRoleCounts = metrics.ReplicatedRoleCounts
			targetManagedKeys = metrics.ReplicatedManagedKeys
			targetSecretEngineResourceCounts = metrics.ReplicatedSecretEngineResourceCounts
			targetRoleAttribution = metrics.ReplicatedRoleAttribution
			targetManagedKeyAttribution = metrics.ReplicatedManagedKeyAttribution
			targetKvAttribution = metrics.ReplicatedKvAttribution
		}

		// Collect metrics counts and attributions based on plugin type
		c.collectMetricsForSecretMount(ctx, entry, targetRoleCounts, targetManagedKeys, targetSecretEngineResourceCounts,
			targetRoleAttribution, targetManagedKeyAttribution, collectAttribution)

		// If this is a KV mount, gather its secret count and mount attribution
		if kvMount := getKVMountMetadata(entry); kvMount != nil {
			c.walkKvMountSecrets(ctx, kvMount)
			if targetKvLocal {
				metrics.LocalKvCounts += kvMount.NumSecrets
			} else {
				metrics.ReplicatedKvCounts += kvMount.NumSecrets
			}

			// Collect KV attribution if requested and mount has secrets
			if collectAttribution && kvMount.NumSecrets > 0 {
				attribution := logical.MountAttribution{
					Count:            kvMount.NumSecrets,
					MountAccessor:    kvMount.MountAccessor,
					MountPath:        kvMount.MountPoint,
					MountType:        entry.Type,
					NamespaceID:      kvMount.Namespace.ID,
					NamespacePath:    kvMount.Namespace.Path,
					BackendAwareUUID: entry.BackendAwareUUID,
				}
				// Use mount accessor as the map key
				targetKvAttribution[kvMount.MountAccessor] = attribution
			}
		}
	}

	return metrics, nil
}

// collectMetricsForSecretMount collects role counts, managed key counts, and secret engine resource counts
// for a specific mount entry based on its plugin type. It updates the provided target metrics in place.
func (c *Core) collectMetricsForSecretMount(ctx context.Context, entry *MountEntry, targetRoleCounts *RoleCounts, targetManagedKeys *ManagedKeyCounts, targetSecretEngineResourceCounts *SecretEngineResourceCounts, targetRoleAttribution, targetManagedKeyAttribution map[string]MountAttributionMap, collectAttribution bool) {
	// apiListWithKeyInfo lists keys for a mount's apiPath and also returns the key_info
	// map, used when the caller needs per-key metadata (e.g. SSH key_type for CA vs OTP).
	apiListWithKeyInfo := func(entry *MountEntry, apiPath string) ([]string, map[string]interface{}) {
		listRequest := &logical.Request{
			Operation: logical.ListOperation,
			Path:      entry.namespace.Path + entry.Path + apiPath,
		}

		resp, err := c.router.Route(ctx, listRequest)
		if err != nil || resp == nil || resp.Data == nil {
			return nil, nil
		}

		var keyInfo map[string]interface{}
		if rawKeyInfo, ok := resp.Data["key_info"]; ok {
			if ki, ok := rawKeyInfo.(map[string]interface{}); ok {
				keyInfo = ki
			}
		}

		rawKeys, ok := resp.Data["keys"]
		if !ok || rawKeys == nil {
			return nil, keyInfo
		}

		// Type switch handles both the 'official' behavior and the 'generic' behavior
		switch kt := rawKeys.(type) {
		case []string:
			// Existing plugins likely hit this path
			return kt, keyInfo
		case []interface{}:
			// External/RPC plugins likely hit this path
			keys := make([]string, 0, len(kt))
			for _, k := range kt {
				if s, ok := k.(string); ok {
					keys = append(keys, s)
				}
			}
			return keys, keyInfo
		default:
			// If it's something totally weird, we still fail safely
			return nil, keyInfo
		}
	}

	// apiList returns just the listed keys, discarding key_info.
	apiList := func(entry *MountEntry, apiPath string) []string {
		keys, _ := apiListWithKeyInfo(entry, apiPath)
		return keys
	}

	// apiReadCount calls the given GET endpoint on a mount and returns the integer value
	// of the "count" key in the response, or 0 on any failure.
	apiReadCount := func(entry *MountEntry, apiPath string) int {
		readRequest := &logical.Request{
			Operation: logical.ReadOperation,
			Path:      entry.namespace.Path + entry.Path + apiPath,
		}
		resp, err := c.router.Route(ctx, readRequest)
		if err != nil || resp == nil || resp.Data == nil {
			return 0
		}
		rawCount, ok := resp.Data["count"]
		if !ok || rawCount == nil {
			return 0
		}
		switch v := rawCount.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case uint:
			return int(v)
		case uint64:
			return int(v)
		case float64:
			return int(v)
		default:
			return 0
		}
	}

	// Helper to add role attribution by type
	addRoleAttribution := func(roleType string, count int) {
		if collectAttribution && count > 0 {
			attribution := logical.MountAttribution{
				Count:            count,
				MountAccessor:    entry.Accessor,
				MountPath:        entry.Path,
				MountType:        entry.Type,
				NamespaceID:      entry.NamespaceID,
				NamespacePath:    entry.namespace.Path,
				BackendAwareUUID: entry.BackendAwareUUID,
			}
			// Initialize the inner map if it doesn't exist
			if targetRoleAttribution[roleType] == nil {
				targetRoleAttribution[roleType] = make(MountAttributionMap)
			}
			// Use mount accessor as the key
			targetRoleAttribution[roleType][entry.Accessor] = attribution
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
		addRoleAttribution(billing.AWSDynamicRoles, len(dynamicRoles))
		staticRoles := apiList(entry, "static-roles")
		targetRoleCounts.AWSStaticRoles += len(staticRoles)
		addRoleAttribution(billing.AWSStaticRoles, len(staticRoles))

	case pluginconsts.SecretEngineAzure:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.AzureDynamicRoles += len(dynamicRoles)
		addRoleAttribution(billing.AzureDynamicRoles, len(dynamicRoles))
		staticRoles := apiList(entry, "static-roles")
		targetRoleCounts.AzureStaticRoles += len(staticRoles)
		addRoleAttribution(billing.AzureStaticRoles, len(staticRoles))

	case pluginconsts.SecretEngineDatabase:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.DatabaseDynamicRoles += len(dynamicRoles)
		addRoleAttribution(billing.DatabaseDynamicRoles, len(dynamicRoles))
		staticRoles := apiList(entry, "static-roles")
		targetRoleCounts.DatabaseStaticRoles += len(staticRoles)
		addRoleAttribution(billing.DatabaseStaticRoles, len(staticRoles))

	case pluginconsts.SecretEngineGCP:
		rolesets := apiList(entry, "rolesets")
		targetRoleCounts.GCPRolesets += len(rolesets)
		addRoleAttribution(billing.GCPRolesets, len(rolesets))
		staticAccounts := apiList(entry, "static-accounts")
		targetRoleCounts.GCPStaticAccounts += len(staticAccounts)
		addRoleAttribution(billing.GCPStaticAccounts, len(staticAccounts))
		impersonatedAccounts := apiList(entry, "impersonated-accounts")
		targetRoleCounts.GCPImpersonatedAccounts += len(impersonatedAccounts)
		addRoleAttribution(billing.GCPImpersonatedAccounts, len(impersonatedAccounts))

	case pluginconsts.SecretEngineLDAP:
		dynamicRoles := apiReadCount(entry, "role-count")
		targetRoleCounts.LDAPDynamicRoles += dynamicRoles
		addRoleAttribution(billing.LDAPDynamicRoles, dynamicRoles)
		staticRoles := apiReadCount(entry, "static-role-count")
		targetRoleCounts.LDAPStaticRoles += staticRoles
		addRoleAttribution(billing.LDAPStaticRoles, staticRoles)
		librarySets := apiReadCount(entry, "library-count")
		targetRoleCounts.LDAPLibrarySets += librarySets
		addRoleAttribution(billing.LDAPLibrarySets, librarySets)

	case pluginconsts.SecretEngineOpenLDAP:
		dynamicRoles := apiReadCount(entry, "role-count")
		targetRoleCounts.OpenLDAPDynamicRoles += dynamicRoles
		addRoleAttribution(billing.OpenLDAPDynamicRoles, dynamicRoles)
		staticRoles := apiReadCount(entry, "static-role-count")
		targetRoleCounts.OpenLDAPStaticRoles += staticRoles
		addRoleAttribution(billing.OpenLDAPStaticRoles, staticRoles)
		librarySets := apiReadCount(entry, "library-count")
		targetRoleCounts.OpenLDAPLibrarySets += librarySets
		addRoleAttribution(billing.OpenLDAPLibrarySets, librarySets)

	case pluginconsts.SecretEngineAlicloud:
		dynamicRoles := apiList(entry, "role")
		targetRoleCounts.AlicloudDynamicRoles += len(dynamicRoles)
		addRoleAttribution(billing.AlicloudDynamicRoles, len(dynamicRoles))

	case pluginconsts.SecretEngineRabbitMQ:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.RabbitMQDynamicRoles += len(dynamicRoles)
		addRoleAttribution(billing.RabbitMQDynamicRoles, len(dynamicRoles))

	case pluginconsts.SecretEngineConsul:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.ConsulDynamicRoles += len(dynamicRoles)
		addRoleAttribution(billing.ConsulDynamicRoles, len(dynamicRoles))

	case pluginconsts.SecretEngineNomad:
		dynamicRoles := apiList(entry, "role")
		targetRoleCounts.NomadDynamicRoles += len(dynamicRoles)
		addRoleAttribution(billing.NomadDynamicRoles, len(dynamicRoles))

	case pluginconsts.SecretEngineKubernetes:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.KubernetesDynamicRoles += len(dynamicRoles)
		addRoleAttribution(billing.KubernetesDynamicRoles, len(dynamicRoles))

	case pluginconsts.SecretEngineMongoDBAtlas:
		dynamicRoles := apiList(entry, "roles")
		targetRoleCounts.MongoDBAtlasDynamicRoles += len(dynamicRoles)
		addRoleAttribution(billing.MongoDBAtlasDynamicRoles, len(dynamicRoles))

	case pluginconsts.SecretEngineTerraform:
		dynamicRoles := apiList(entry, "role")
		targetRoleCounts.TerraformCloudDynamicRoles += len(dynamicRoles)
		addRoleAttribution(billing.TerraformCloudDynamicRoles, len(dynamicRoles))

	// Transit secret engine not billed by resource counts so skip attribution collection here
	case pluginconsts.SecretEngineTransit:
		transitKeys := apiList(entry, "keys")
		targetManagedKeys.TransitKeys += len(transitKeys)

	case pluginconsts.SecretEngineTOTP:
		keyCountPerEntry := apiList(entry, "keys")
		targetManagedKeys.TotpKeys += len(keyCountPerEntry)
		if collectAttribution && len(keyCountPerEntry) > 0 {
			attribution := logical.MountAttribution{
				Count:            len(keyCountPerEntry),
				MountAccessor:    entry.Accessor,
				MountPath:        entry.Path,
				MountType:        entry.Type,
				NamespaceID:      entry.NamespaceID,
				NamespacePath:    entry.namespace.Path,
				BackendAwareUUID: entry.BackendAwareUUID,
			}
			// Initialize the inner map if it doesn't exist
			if targetManagedKeyAttribution[billing.TotpKeys] == nil {
				targetManagedKeyAttribution[billing.TotpKeys] = make(MountAttributionMap)
			}
			// Use mount accessor as the key
			targetManagedKeyAttribution[billing.TotpKeys][entry.Accessor] = attribution
		}

	// Transform secret engine not billed by resource counts so skip attribution collection here
	case pluginconsts.SecretEngineTransform:
		transformRoles := apiList(entry, "role")
		targetRoleCounts.TransformRoles += len(transformRoles)

		// Collect Transform secret engine resource counts
		transformations := apiList(entry, "transformation")
		targetSecretEngineResourceCounts.TransformTransformations += len(transformations)

		templates := apiList(entry, "template")
		targetSecretEngineResourceCounts.TransformTemplates += len(templates)

		alphabets := apiList(entry, "alphabet")
		targetSecretEngineResourceCounts.TransformAlphabets += len(alphabets)

		stores := apiList(entry, "stores")
		targetSecretEngineResourceCounts.TransformStores += len(stores)

	// KMIP secret engine not billed by resource counts so skip attribution collection here
	case pluginconsts.SecretEngineKMIP:
		// Collect KMIP secret engine resource counts
		scopes := apiList(entry, "scope")
		targetSecretEngineResourceCounts.KmipScopes += len(scopes)

		// Count roles across all scopes
		for _, scope := range scopes {
			if scope == "" {
				continue
			}
			scopeName := strings.TrimSuffix(scope, "/")
			roles := apiList(entry, "scope/"+scopeName+"/role")
			targetSecretEngineResourceCounts.KmipScopeRoles += len(roles)
		}

		// Count CAs
		cas := apiList(entry, "ca")
		targetSecretEngineResourceCounts.KmipCas += len(cas)

	case pluginconsts.SecretEngineKeymgmt:
		keyCountPerEntry := apiList(entry, "key")
		targetManagedKeys.KmseKeys += len(keyCountPerEntry)
		if collectAttribution && len(keyCountPerEntry) > 0 {
			attribution := logical.MountAttribution{
				Count:            len(keyCountPerEntry),
				MountAccessor:    entry.Accessor,
				MountPath:        entry.Path,
				MountType:        entry.Type,
				NamespaceID:      entry.NamespaceID,
				NamespacePath:    entry.namespace.Path,
				BackendAwareUUID: entry.BackendAwareUUID,
			}
			// Initialize the inner map if it doesn't exist
			if targetManagedKeyAttribution[billing.KmseKeys] == nil {
				targetManagedKeyAttribution[billing.KmseKeys] = make(MountAttributionMap)
			}
			// Use mount accessor as the key
			targetManagedKeyAttribution[billing.KmseKeys][entry.Accessor] = attribution
		}

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
		addRoleAttribution(billing.OSLocalAccountRoles, accountCount)

	// SPIFFE secret engine not billed by resource counts so skip attribution collection here
	case pluginconsts.SecretEngineSpiffe:
		spiffeRoles := apiList(entry, "role")
		targetRoleCounts.SpiffeRoles += len(spiffeRoles)

	// SSH secret engine not billed by resource counts so skip attribution collection here
	case pluginconsts.SecretEngineSsh:
		sshRoles, sshRoleInfo := apiListWithKeyInfo(entry, "roles")
		for _, role := range sshRoles {
			info, ok := sshRoleInfo[role].(map[string]interface{})
			if !ok {
				continue
			}
			switch info["key_type"] {
			case "ca":
				targetRoleCounts.SSHCARoles++
			case "otp":
				targetRoleCounts.SSHOTPRoles++
			}
		}
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
	m, err := c.CountMetricsSecretMounts(false, false)
	if err != nil {
		return nil
	}
	return combineRoleCounts(m.LocalRoleCounts, m.ReplicatedRoleCounts)
}

// GetRoleAndManagedKeyCounts returns the local or replicated role and managed key counts and attribution depending on the mount type
// For use in tests only
func (c *Core) GetRoleAndManagedKeyCountsAndAttribution(mountType string) (*RoleCounts, *ManagedKeyCounts, map[string]MountAttributionMap, map[string]MountAttributionMap) {
	m, err := c.CountMetricsSecretMounts(true, true)
	if err != nil {
		return nil, nil, nil, nil
	}
	if mountType == billing.LocalPrefix {
		return m.LocalRoleCounts, m.LocalManagedKeys, m.LocalRoleAttribution, m.LocalManagedKeyAttribution
	}
	return m.ReplicatedRoleCounts, m.ReplicatedManagedKeys, m.ReplicatedRoleAttribution, m.ReplicatedManagedKeyAttribution
}

// GetKvCounts returns the local or replicated KV counts and attribution depending on the mount type
// For use in tests only
func (c *Core) GetKvCountsAndAttribution(mountType string) (int, MountAttributionMap) {
	m, err := c.CountMetricsSecretMounts(true, true)
	if err != nil {
		return 0, nil
	}
	if mountType == billing.LocalPrefix {
		return m.LocalKvCounts, m.LocalKvAttribution
	}
	return m.ReplicatedKvCounts, m.ReplicatedKvAttribution
}
