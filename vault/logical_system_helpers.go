// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	invalidateMFAConfig                      = func(context.Context, *SystemBackend, string) {}
	invalidateLoginMFAConfig                 = func(context.Context, *SystemBackend, string) {}
	invalidateLoginMFALoginEnforcementConfig = func(context.Context, *SystemBackend, string) {}

	sysInvalidate = func(b *SystemBackend) func(context.Context, string) {
		return nil
	}

	getSystemSchemas = func() []func() *memdb.TableSchema { return nil }

	getEGPListResponseKeyInfo = func(*SystemBackend, *namespace.Namespace) map[string]interface{} { return nil }
	addSentinelPolicyData     = func(map[string]interface{}, *Policy) {}
	inputSentinelPolicyData   = func(*framework.FieldData, *Policy) *logical.Response { return nil }

	controlGroupUnwrap = func(context.Context, *SystemBackend, string, bool) (string, error) {
		return "", errors.New("control groups unavailable")
	}

	pathInternalUINamespacesRead = func(b *SystemBackend) framework.OperationFunc {
		return func(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
			// Short-circuit here if there's no client token provided
			if req.ClientToken == "" {
				return nil, fmt.Errorf("client token empty")
			}

			// Load the ACL policies so we can check for access and filter namespaces
			_, te, entity, _, err := b.Core.fetchACLTokenEntryAndEntity(ctx, req)
			if err != nil {
				return nil, err
			}
			if entity != nil && entity.Disabled {
				b.logger.Warn("permission denied as the entity on the token is disabled")
				return nil, logical.ErrPermissionDenied
			}
			if te != nil && te.EntityID != "" && entity == nil {
				b.logger.Warn("permission denied as the entity on the token is invalid")
				return nil, logical.ErrPermissionDenied
			}

			return logical.ListResponse([]string{""}), nil
		}
	}

	entPaths = func(b *SystemBackend) []*framework.Path {
		buildEnterpriseOnlyPaths := func(paths map[string][]logical.Operation) []*framework.Path {
			var results []*framework.Path
			for pattern, operations := range paths {
				operationsMap := map[logical.Operation]framework.OperationHandler{}

				for _, operation := range operations {
					operationsMap[operation] = &framework.PathOperation{
						Callback: func(context.Context, *logical.Request, *framework.FieldData) (*logical.Response, error) {
							return logical.ErrorResponse("enterprise-only feature"), logical.ErrUnsupportedPath
						},
					}
				}

				results = append(results, &framework.Path{
					Pattern:    pattern,
					Operations: operationsMap,
				})
			}

			return results
		}

		var paths []*framework.Path

		// license paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"license/status$": {logical.ReadOperation},
		})...)

		// group-policy-application paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"config/group-policy-application$": {logical.ReadOperation, logical.UpdateOperation},
		})...)

		// namespaces paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"namespaces/?$": {logical.ListOperation},
			"namespaces/api-lock/lock" + framework.OptionalParamRegex("path"):   {logical.UpdateOperation},
			"namespaces/api-lock/unlock" + framework.OptionalParamRegex("path"): {logical.UpdateOperation},
			"namespaces/(?P<path>.+?)": {logical.DeleteOperation, logical.PatchOperation, logical.ReadOperation, logical.UpdateOperation},
		})...)

		// replication paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"replication/performance/primary/enable":                                               {logical.UpdateOperation},
			"replication/dr/primary/enable":                                                        {logical.UpdateOperation},
			"replication/performance/primary/demote":                                               {logical.UpdateOperation},
			"replication/dr/primary/demote":                                                        {logical.UpdateOperation},
			"replication/performance/primary/disable":                                              {logical.UpdateOperation},
			"replication/dr/primary/disable":                                                       {logical.UpdateOperation},
			"replication/performance/primary/secondary-token":                                      {logical.UpdateOperation},
			"replication/dr/primary/secondary-token":                                               {logical.UpdateOperation},
			"replication/performance/primary/revoke-secondary":                                     {logical.UpdateOperation},
			"replication/dr/primary/revoke-secondary":                                              {logical.UpdateOperation},
			"replication/performance/secondary/generate-public-key":                                {logical.UpdateOperation},
			"replication/dr/secondary/generate-public-key":                                         {logical.UpdateOperation},
			"replication/performance/secondary/enable":                                             {logical.UpdateOperation},
			"replication/dr/secondary/enable":                                                      {logical.UpdateOperation},
			"replication/performance/secondary/promote":                                            {logical.UpdateOperation},
			"replication/dr/secondary/promote":                                                     {logical.UpdateOperation},
			"replication/performance/secondary/disable":                                            {logical.UpdateOperation},
			"replication/dr/secondary/disable":                                                     {logical.UpdateOperation},
			"replication/dr/secondary/operation-token/delete":                                      {logical.UpdateOperation},
			"replication/performance/secondary/update-primary":                                     {logical.UpdateOperation},
			"replication/dr/secondary/update-primary":                                              {logical.UpdateOperation},
			"replication/dr/secondary/license/status":                                              {logical.ReadOperation},
			"replication/dr/secondary/config/reload/(?P<subsystem>.+)":                             {logical.UpdateOperation},
			"replication/recover":                                                                  {logical.UpdateOperation},
			"replication/dr/secondary/recover":                                                     {logical.UpdateOperation},
			"replication/dr/secondary/reindex":                                                     {logical.UpdateOperation},
			"replication/reindex":                                                                  {logical.UpdateOperation},
			"replication/status":                                                                   {logical.ReadOperation},
			"replication/dr/status":                                                                {logical.ReadOperation},
			"replication/performance/status":                                                       {logical.ReadOperation},
			"replication/primary/enable":                                                           {logical.UpdateOperation},
			"replication/primary/demote":                                                           {logical.UpdateOperation},
			"replication/primary/disable":                                                          {logical.UpdateOperation},
			"replication/primary/secondary-token":                                                  {logical.UpdateOperation},
			"replication/performance/primary/paths-filter/" + framework.GenericNameRegex("id"):     {logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
			"replication/performance/primary/dynamic-filter/" + framework.GenericNameRegex("id"):   {logical.ReadOperation},
			"replication/primary/revoke-secondary":                                                 {logical.UpdateOperation},
			"replication/secondary/enable":                                                         {logical.UpdateOperation},
			"replication/secondary/promote":                                                        {logical.UpdateOperation},
			"replication/secondary/disable":                                                        {logical.UpdateOperation},
			"replication/secondary/update-primary":                                                 {logical.UpdateOperation},
			"replication/performance/secondary/dynamic-filter/" + framework.GenericNameRegex("id"): {logical.ReadOperation},
		})...)

		// seal paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"sealwrap/rewrap": {logical.ReadOperation, logical.UpdateOperation},
		})...)

		// mfa paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"mfa/method/?": {logical.ListOperation},
			"mfa/method/totp/" + framework.GenericNameRegex("name") + "/generate$":       {logical.ReadOperation},
			"mfa/method/totp/" + framework.GenericNameRegex("name") + "/admin-generate$": {logical.UpdateOperation},
			"mfa/method/totp/" + framework.GenericNameRegex("name") + "/admin-destroy$":  {logical.UpdateOperation},
			"mfa/method/totp/" + framework.GenericNameRegex("name"):                      {logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
			"mfa/method/okta/" + framework.GenericNameRegex("name"):                      {logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
			"mfa/method/duo/" + framework.GenericNameRegex("name"):                       {logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
			"mfa/method/pingid/" + framework.GenericNameRegex("name"):                    {logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
		})...)

		// control-group paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"control-group/authorize": {logical.UpdateOperation},
			"control-group/request":   {logical.UpdateOperation},
			"config/control-group":    {logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
		})...)

		// sentinel paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"policies/rgp/?$":           {logical.ListOperation},
			"policies/rgp/(?P<name>.+)": {logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
			"policies/egp/?$":           {logical.ListOperation},
			"policies/egp/(?P<name>.+)": {logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
		})...)

		// plugins reload status paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"plugins/reload/backend/status$": {logical.ReadOperation},
		})...)

		// quotas paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"quotas/lease-count/?$": {logical.ListOperation},
			"quotas/lease-count/" + framework.GenericNameRegex("name"): {logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
		})...)

		// raft auto-snapshot paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"storage/raft/snapshot-auto/config/":                                      {logical.ListOperation},
			"storage/raft/snapshot-auto/config/" + framework.GenericNameRegex("name"): {logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
		})...)

		paths = append(paths, buildEnterpriseOnlyPaths(map[string][]logical.Operation{
			"managed-keys/" + framework.GenericNameRegex("type") + "/?":                                                    {logical.ListOperation},
			"managed-keys/" + framework.GenericNameRegex("type") + "/" + framework.GenericNameRegex("name"):                {logical.CreateOperation, logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation},
			"managed-keys/" + framework.GenericNameRegex("type") + "/" + framework.GenericNameRegex("name") + "/test/sign": {logical.CreateOperation, logical.UpdateOperation},
		})...)

		return paths
	}
	handleGlobalPluginReload = func(context.Context, *Core, string, string, []string) error {
		return nil
	}
	handleSetupPluginReload = func(*Core) error {
		return nil
	}
	handleLicenseReload = func(b *SystemBackend) framework.OperationFunc {
		return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
			return nil, nil
		}
	}
	checkRaw = func(b *SystemBackend, path string) error { return nil }
)

// tuneMount is used to set config on a mount point
func (b *SystemBackend) tuneMountTTLs(ctx context.Context, path string, me *MountEntry, newDefault, newMax time.Duration) error {
	zero := time.Duration(0)

	switch {
	case newDefault == zero && newMax == zero:
		// No checks needed

	case newDefault == zero && newMax != zero:
		// No default/max conflict, no checks needed

	case newDefault != zero && newMax == zero:
		// No default/max conflict, no checks needed

	case newDefault != zero && newMax != zero:
		if newMax < newDefault {
			return fmt.Errorf("backend max lease TTL of %d would be less than backend default lease TTL of %d", int(newMax.Seconds()), int(newDefault.Seconds()))
		}
	}

	origMax := me.Config.MaxLeaseTTL
	origDefault := me.Config.DefaultLeaseTTL

	me.Config.MaxLeaseTTL = newMax
	me.Config.DefaultLeaseTTL = newDefault

	// Update the mount table
	var err error
	switch {
	case strings.HasPrefix(path, credentialRoutePrefix):
		err = b.Core.persistAuth(ctx, b.Core.auth, &me.Local)
	default:
		err = b.Core.persistMounts(ctx, b.Core.mounts, &me.Local)
	}
	if err != nil {
		me.Config.MaxLeaseTTL = origMax
		me.Config.DefaultLeaseTTL = origDefault
		return fmt.Errorf("failed to update mount table, rolling back TTL changes")
	}
	if b.Core.logger.IsInfo() {
		b.Core.logger.Info("mount tuning of leases successful", "path", path)
	}

	return nil
}
