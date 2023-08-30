// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

type enterprisePathStub struct {
	parameters []string
	operations []logical.Operation
}

var (
	invalidateMFAConfig                      = func(context.Context, *SystemBackend, string) {}
	invalidateLoginMFAConfig                 = func(context.Context, *SystemBackend, string) {}
	invalidateLoginMFALoginEnforcementConfig = func(context.Context, *SystemBackend, string) {}

	sysInvalidate = func(b *SystemBackend) func(context.Context, string) {
		return nil
	}

	sysInitialize = func(b *SystemBackend) func(context.Context, *logical.InitializationRequest) error {
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

	entPaths     = entStubPaths
	entStubPaths = func(b *SystemBackend) []*framework.Path {
		buildEnterpriseOnlyPaths := func(paths map[string]enterprisePathStub) []*framework.Path {
			var results []*framework.Path
			for pattern, pathSpec := range paths {
				path := &framework.Path{
					Pattern:    pattern,
					Operations: make(map[logical.Operation]framework.OperationHandler),
					Fields:     make(map[string]*framework.FieldSchema),
					DisplayAttrs: &framework.DisplayAttributes{
						// Since we lack full information for Fields, and all information for Responses, the generated
						// OpenAPI won't be good for much other than identifying the endpoint exists at all. Thus, it
						// is useful to make it clear that this is only a stub. Code generation will use this to ignore
						// these operations.
						OperationPrefix: "enterprise-stub",
					},
				}

				for _, parameter := range pathSpec.parameters {
					path.Fields[parameter] = &framework.FieldSchema{
						Type:     framework.TypeString,
						Required: true,
					}
				}

				for _, operation := range pathSpec.operations {
					path.Operations[operation] = &framework.PathOperation{
						Callback: func(context.Context, *logical.Request, *framework.FieldData) (*logical.Response, error) {
							return logical.ErrorResponse("enterprise-only feature"), logical.ErrUnsupportedPath
						},
					}

					// There is a correctness check that verifies there is an ExistenceFunc for all paths that have
					// a CreateOperation, so we must define a stub one to pass that check if needed.
					if operation == logical.CreateOperation {
						path.ExistenceCheck = func(context.Context, *logical.Request, *framework.FieldData) (bool, error) {
							return false, nil
						}
					}
				}

				results = append(results, path)
			}

			return results
		}

		var paths []*framework.Path

		// license paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"license/status$": {operations: []logical.Operation{logical.ReadOperation}},
		})...)

		// group-policy-application paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"config/group-policy-application$": {operations: []logical.Operation{logical.ReadOperation, logical.UpdateOperation}},
		})...)

		// namespaces paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"namespaces/?$": {operations: []logical.Operation{logical.ListOperation}},
			"namespaces/api-lock/lock" + framework.OptionalParamRegex("path"):   {parameters: []string{"path"}, operations: []logical.Operation{logical.UpdateOperation}},
			"namespaces/api-lock/unlock" + framework.OptionalParamRegex("path"): {parameters: []string{"path"}, operations: []logical.Operation{logical.UpdateOperation}},
			"namespaces/(?P<path>.+?)": {parameters: []string{"path"}, operations: []logical.Operation{logical.DeleteOperation, logical.PatchOperation, logical.ReadOperation, logical.UpdateOperation}},
		})...)

		// replication paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"replication/performance/primary/enable":                                               {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/primary/enable":                                                        {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/primary/demote":                                               {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/primary/demote":                                                        {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/primary/disable":                                              {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/primary/disable":                                                       {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/primary/secondary-token":                                      {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/primary/secondary-token":                                               {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/primary/revoke-secondary":                                     {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/primary/revoke-secondary":                                              {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/secondary/generate-public-key":                                {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/secondary/generate-public-key":                                         {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/secondary/enable":                                             {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/secondary/enable":                                                      {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/secondary/promote":                                            {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/secondary/promote":                                                     {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/secondary/disable":                                            {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/secondary/disable":                                                     {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/secondary/operation-token/delete":                                      {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/secondary/update-primary":                                     {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/secondary/update-primary":                                              {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/secondary/license/status":                                              {operations: []logical.Operation{logical.ReadOperation}},
			"replication/dr/secondary/config/reload/(?P<subsystem>.+)":                             {parameters: []string{"subsystem"}, operations: []logical.Operation{logical.UpdateOperation}},
			"replication/recover":                                                                  {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/secondary/recover":                                                     {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/secondary/reindex":                                                     {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/reindex":                                                                  {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/dr/status":                                                                {operations: []logical.Operation{logical.ReadOperation}},
			"replication/performance/status":                                                       {operations: []logical.Operation{logical.ReadOperation}},
			"replication/primary/enable":                                                           {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/primary/demote":                                                           {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/primary/disable":                                                          {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/primary/secondary-token":                                                  {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/primary/paths-filter/" + framework.GenericNameRegex("id"):     {parameters: []string{"id"}, operations: []logical.Operation{logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
			"replication/performance/primary/dynamic-filter/" + framework.GenericNameRegex("id"):   {parameters: []string{"id"}, operations: []logical.Operation{logical.ReadOperation}},
			"replication/primary/revoke-secondary":                                                 {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/secondary/enable":                                                         {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/secondary/promote":                                                        {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/secondary/disable":                                                        {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/secondary/update-primary":                                                 {operations: []logical.Operation{logical.UpdateOperation}},
			"replication/performance/secondary/dynamic-filter/" + framework.GenericNameRegex("id"): {parameters: []string{"id"}, operations: []logical.Operation{logical.ReadOperation}},
		})...)
		// This path, though an enterprise path, has always been handled in OSS.
		paths = append(paths, &framework.Path{
			Pattern: "replication/status",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
						resp := &logical.Response{
							Data: map[string]interface{}{
								"mode": "disabled",
							},
						}
						return resp, nil
					},
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "read",
						OperationSuffix: "replication-status",
					},
				},
			},
		})

		// seal paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"sealwrap/rewrap": {operations: []logical.Operation{logical.ReadOperation, logical.UpdateOperation}},
		})...)

		// mfa paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"mfa/method/?": {operations: []logical.Operation{logical.ListOperation}},
			"mfa/method/totp/" + framework.GenericNameRegex("name") + "/generate$":       {parameters: []string{"name"}, operations: []logical.Operation{logical.ReadOperation}},
			"mfa/method/totp/" + framework.GenericNameRegex("name") + "/admin-generate$": {parameters: []string{"name"}, operations: []logical.Operation{logical.UpdateOperation}},
			"mfa/method/totp/" + framework.GenericNameRegex("name") + "/admin-destroy$":  {parameters: []string{"name"}, operations: []logical.Operation{logical.UpdateOperation}},
			"mfa/method/totp/" + framework.GenericNameRegex("name"):                      {parameters: []string{"name"}, operations: []logical.Operation{logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
			"mfa/method/okta/" + framework.GenericNameRegex("name"):                      {parameters: []string{"name"}, operations: []logical.Operation{logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
			"mfa/method/duo/" + framework.GenericNameRegex("name"):                       {parameters: []string{"name"}, operations: []logical.Operation{logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
			"mfa/method/pingid/" + framework.GenericNameRegex("name"):                    {parameters: []string{"name"}, operations: []logical.Operation{logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
		})...)

		// control-group paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"control-group/authorize": {operations: []logical.Operation{logical.UpdateOperation}},
			"control-group/request":   {operations: []logical.Operation{logical.UpdateOperation}},
			"config/control-group":    {operations: []logical.Operation{logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
		})...)

		// sentinel paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"policies/rgp/?$":           {operations: []logical.Operation{logical.ListOperation}},
			"policies/rgp/(?P<name>.+)": {parameters: []string{"name"}, operations: []logical.Operation{logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
			"policies/egp/?$":           {operations: []logical.Operation{logical.ListOperation}},
			"policies/egp/(?P<name>.+)": {parameters: []string{"name"}, operations: []logical.Operation{logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
		})...)

		// plugins reload status paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"plugins/reload/backend/status$": {operations: []logical.Operation{logical.ReadOperation}},
		})...)

		// quotas paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"quotas/lease-count/?$": {operations: []logical.Operation{logical.ListOperation}},
			"quotas/lease-count/" + framework.GenericNameRegex("name"): {parameters: []string{"name"}, operations: []logical.Operation{logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
		})...)

		// raft auto-snapshot paths
		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"storage/raft/snapshot-auto/config/":                                      {operations: []logical.Operation{logical.ListOperation}},
			"storage/raft/snapshot-auto/config/" + framework.GenericNameRegex("name"): {parameters: []string{"name"}, operations: []logical.Operation{logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
			"storage/raft/snapshot-auto/status/" + framework.GenericNameRegex("name"): {parameters: []string{"name"}, operations: []logical.Operation{logical.ReadOperation}},
		})...)

		paths = append(paths, buildEnterpriseOnlyPaths(map[string]enterprisePathStub{
			"managed-keys/" + framework.GenericNameRegex("type") + "/?":                                                    {parameters: []string{"type"}, operations: []logical.Operation{logical.ListOperation}},
			"managed-keys/" + framework.GenericNameRegex("type") + "/" + framework.GenericNameRegex("name"):                {parameters: []string{"type", "name"}, operations: []logical.Operation{logical.CreateOperation, logical.DeleteOperation, logical.ReadOperation, logical.UpdateOperation}},
			"managed-keys/" + framework.GenericNameRegex("type") + "/" + framework.GenericNameRegex("name") + "/test/sign": {parameters: []string{"type", "name"}, operations: []logical.Operation{logical.CreateOperation, logical.UpdateOperation}},
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
