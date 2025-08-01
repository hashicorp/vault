// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/random"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/license"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/rotation"
	"github.com/hashicorp/vault/vault/plugincatalog"
	"github.com/hashicorp/vault/version"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ctxKeyForwardedRequestMountAccessor struct{}

func (c ctxKeyForwardedRequestMountAccessor) String() string {
	return "forwarded-req-mount-accessor"
}

type dynamicSystemView struct {
	core        *Core
	mountEntry  *MountEntry
	perfStandby bool
}

func (d dynamicSystemView) DefaultLeaseTTL() time.Duration {
	def, _ := d.fetchTTLs()
	return def
}

func (d dynamicSystemView) MaxLeaseTTL() time.Duration {
	_, max := d.fetchTTLs()
	return max
}

// TTLsByPath returns the default and max TTLs corresponding to a particular
// mount point, or the system default
func (d dynamicSystemView) fetchTTLs() (def, max time.Duration) {
	def = d.core.defaultLeaseTTL
	max = d.core.maxLeaseTTL

	if d.mountEntry != nil {
		if d.mountEntry.Config.DefaultLeaseTTL != 0 {
			def = d.mountEntry.Config.DefaultLeaseTTL
		}
		if d.mountEntry.Config.MaxLeaseTTL != 0 {
			max = d.mountEntry.Config.MaxLeaseTTL
		}
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
	state := d.core.ReplicationState()
	if d.perfStandby {
		state |= consts.ReplicationPerformanceStandby
	}
	return state
}

func (d dynamicSystemView) HasFeature(feature license.Features) bool {
	return d.core.HasFeature(feature)
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

func (d dynamicSystemView) NewPluginClient(ctx context.Context, config pluginutil.PluginClientConfig) (pluginutil.PluginClient, error) {
	if d.core == nil {
		return nil, fmt.Errorf("system view core is nil")
	}
	if d.core.pluginCatalog == nil {
		return nil, fmt.Errorf("system view core plugin catalog is nil")
	}

	c, err := d.core.pluginCatalog.NewPluginClient(ctx, config)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// LookupPlugin looks for a plugin with the given name in the plugin catalog. It
// returns a PluginRunner or an error if no plugin was found.
func (d dynamicSystemView) LookupPlugin(ctx context.Context, name string, pluginType consts.PluginType) (*pluginutil.PluginRunner, error) {
	return d.LookupPluginVersion(ctx, name, pluginType, "")
}

// LookupPluginVersion looks for a plugin with the given name and version in the plugin catalog. It
// returns a PluginRunner or an error if no plugin was found.
func (d dynamicSystemView) LookupPluginVersion(ctx context.Context, name string, pluginType consts.PluginType, version string) (*pluginutil.PluginRunner, error) {
	if d.core == nil {
		return nil, fmt.Errorf("system view core is nil")
	}
	if d.core.pluginCatalog == nil {
		return nil, fmt.Errorf("system view core plugin catalog is nil")
	}
	r, err := d.core.pluginCatalog.Get(ctx, name, pluginType, version)
	if err != nil {
		return nil, err
	}
	if r == nil {
		errContext := name
		if version != "" {
			errContext += fmt.Sprintf(", version=%s", version)
		}
		return nil, fmt.Errorf("%w: %s", plugincatalog.ErrPluginNotFound, errContext)
	}

	return r, nil
}

// ListVersionedPlugins returns information about all plugins of a certain
// typein the catalog, including any versioning information stored for them.
func (d dynamicSystemView) ListVersionedPlugins(ctx context.Context, pluginType consts.PluginType) ([]pluginutil.VersionedPlugin, error) {
	if d.core == nil {
		return nil, fmt.Errorf("system view core is nil")
	}
	if d.core.pluginCatalog == nil {
		return nil, fmt.Errorf("system view core plugin catalog is nil")
	}
	return d.core.pluginCatalog.ListVersionedPlugins(ctx, pluginType)
}

// MlockEnabled returns the configuration setting for enabling mlock on plugins.
func (d dynamicSystemView) MlockEnabled() bool {
	return d.core.enableMlock
}

func (d dynamicSystemView) EntityInfo(entityID string) (*logical.Entity, error) {
	// Requests from token created from the token backend will not have entity information.
	// Return missing entity instead of error when requesting from MemDB.
	if entityID == "" {
		return nil, nil
	}

	if d.core == nil {
		return nil, fmt.Errorf("system view core is nil")
	}
	if d.core.identityStore == nil {
		return nil, fmt.Errorf("system view identity store is nil")
	}

	// Retrieve the entity from MemDB
	entity, err := d.core.identityStore.MemDBEntityByID(entityID, false)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	// Return a subset of the data
	ret := &logical.Entity{
		ID:       entity.ID,
		Name:     entity.Name,
		Disabled: entity.Disabled,
	}

	if entity.Metadata != nil {
		ret.Metadata = make(map[string]string, len(entity.Metadata))
		for k, v := range entity.Metadata {
			ret.Metadata[k] = v
		}
	}

	aliases := make([]*logical.Alias, 0, len(entity.Aliases))
	for _, a := range entity.Aliases {

		// Don't return aliases from other namespaces
		if a.NamespaceID != d.mountEntry.NamespaceID {
			continue
		}

		alias := identity.ToSDKAlias(a)

		// MountType is not stored with the entity and must be looked up
		if mount := d.core.router.ValidateMountByAccessor(a.MountAccessor); mount != nil {
			alias.MountType = mount.MountType
		}

		aliases = append(aliases, alias)
	}
	ret.Aliases = aliases

	return ret, nil
}

func (d dynamicSystemView) GroupsForEntity(entityID string) ([]*logical.Group, error) {
	// Requests from token created from the token backend will not have entity information.
	// Return missing entity instead of error when requesting from MemDB.
	if entityID == "" {
		return nil, nil
	}

	if d.core == nil {
		return nil, fmt.Errorf("system view core is nil")
	}
	if d.core.identityStore == nil {
		return nil, fmt.Errorf("system view identity store is nil")
	}

	groups, inheritedGroups, err := d.core.identityStore.groupsByEntityID(entityID)
	if err != nil {
		return nil, err
	}

	groups = append(groups, inheritedGroups...)

	logicalGroups := make([]*logical.Group, 0, len(groups))
	for _, g := range groups {
		// Don't return groups from other namespaces
		if g.NamespaceID != d.mountEntry.NamespaceID {
			continue
		}

		logicalGroups = append(logicalGroups, identity.ToSDKGroup(g))
	}

	return logicalGroups, nil
}

func (d dynamicSystemView) PluginEnv(_ context.Context) (*logical.PluginEnvironment, error) {
	v := version.GetVersion()

	buildDate, err := version.GetVaultBuildDate()
	if err != nil {
		return nil, err
	}

	return &logical.PluginEnvironment{
		VaultVersion:           v.Version,
		VaultVersionPrerelease: v.VersionPrerelease,
		VaultVersionMetadata:   v.VersionMetadata,
		VaultBuildDate:         timestamppb.New(buildDate),
	}, nil
}

func (d dynamicSystemView) VaultVersion(_ context.Context) (string, error) {
	return version.GetVersion().Version, nil
}

func (d dynamicSystemView) GeneratePasswordFromPolicy(ctx context.Context, policyName string) (password string, err error) {
	if policyName == "" {
		return "", fmt.Errorf("missing password policy name")
	}

	// Ensure there's a timeout on the context of some sort
	if _, hasTimeout := ctx.Deadline(); !hasTimeout {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
	}

	ctx = namespace.ContextWithNamespace(ctx, d.mountEntry.Namespace())

	policyCfg, err := d.retrievePasswordPolicy(ctx, policyName)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve password policy: %w", err)
	}

	if policyCfg == nil {
		return "", fmt.Errorf("no password policy found")
	}

	passPolicy, err := random.ParsePolicy(policyCfg.HCLPolicy)
	if err != nil {
		return "", fmt.Errorf("stored password policy is invalid: %w", err)
	}

	return passPolicy.Generate(ctx, nil)
}

func (d dynamicSystemView) ClusterID(ctx context.Context) (string, error) {
	clusterInfo, err := d.core.Cluster(ctx)
	if err != nil || clusterInfo.ID == "" {
		return "", fmt.Errorf("unable to retrieve cluster info or empty ID: %w", err)
	}

	return clusterInfo.ID, nil
}

func (d dynamicSystemView) GenerateIdentityToken(ctx context.Context, req *pluginutil.IdentityTokenRequest) (*pluginutil.IdentityTokenResponse, error) {
	mountEntry := d.mountEntry
	if mountEntry == nil {
		return nil, fmt.Errorf("no mount entry")
	}
	nsCtx := namespace.ContextWithNamespace(ctx, mountEntry.Namespace())
	storage := d.core.router.MatchingStorageByAPIPath(nsCtx, mountPathIdentity)
	if storage == nil {
		return nil, fmt.Errorf("failed to find storage entry for identity mount")
	}

	token, ttl, err := d.core.IdentityStore().generatePluginIdentityToken(nsCtx, storage, d.mountEntry, req.Audience, req.TTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate plugin identity token: %w", err)
	}

	return &pluginutil.IdentityTokenResponse{
		Token: pluginutil.IdentityToken(token),
		TTL:   ttl,
	}, nil
}

func (d dynamicSystemView) RegisterRotationJob(ctx context.Context, req *rotation.RotationJobConfigureRequest) (string, error) {
	mountEntry := d.mountEntry
	if mountEntry == nil {
		return "", fmt.Errorf("no mount entry")
	}
	nsCtx := namespace.ContextWithNamespace(ctx, mountEntry.Namespace())

	job, err := rotation.ConfigureRotationJob(req)
	if err != nil {
		return "", fmt.Errorf("error configuring rotation job: %s", err)
	}

	id, err := d.core.RegisterRotationJob(nsCtx, job)
	if err != nil {
		return "", fmt.Errorf("error registering rotation job: %s", err)
	}

	job.RotationID = id
	return id, nil
}

func (d dynamicSystemView) DeregisterRotationJob(ctx context.Context, req *rotation.RotationJobDeregisterRequest) (err error) {
	mountEntry := d.mountEntry
	if mountEntry == nil {
		return fmt.Errorf("no mount entry")
	}
	nsCtx := namespace.ContextWithNamespace(ctx, mountEntry.Namespace())

	return d.core.DeregisterRotationJob(nsCtx, req)
}
