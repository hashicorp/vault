// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/vault/quotas"
)

// quotasPaths returns paths that enable quota management
func (b *SystemBackend) quotasPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "quotas/config$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "rate-limit-quotas",
			},

			Fields: map[string]*framework.FieldSchema{
				"rate_limit_exempt_paths": {
					Type:        framework.TypeStringSlice,
					Description: "Specifies the list of exempt paths from all rate limit quotas. If empty no paths will be exempt.",
				},
				"enable_rate_limit_audit_logging": {
					Type:        framework.TypeBool,
					Description: "If set, starts audit logging of requests that get rejected due to rate limit quota rule violations.",
				},
				"enable_rate_limit_response_headers": {
					Type:        framework.TypeBool,
					Description: "If set, additional rate limit quota HTTP headers will be added to responses.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleQuotasConfigUpdate(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "configure",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleQuotasConfigRead(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationSuffix: "configuration",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"enable_rate_limit_audit_logging": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"enable_rate_limit_response_headers": {
									Type:     framework.TypeBool,
									Required: true,
								},
								"rate_limit_exempt_paths": {
									Type:     framework.TypeStringSlice,
									Required: true,
								},
							},
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(quotasHelp["quotas-config"][0]),
			HelpDescription: strings.TrimSpace(quotasHelp["quotas-config"][1]),
		},
		{
			Pattern: "quotas/rate-limit/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "rate-limit-quotas",
				OperationVerb:   "list",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleRateLimitQuotasList(),
				},
			},
			HelpSynopsis:    strings.TrimSpace(quotasHelp["rate-limit-list"][0]),
			HelpDescription: strings.TrimSpace(quotasHelp["rate-limit-list"][1]),
		},
		{
			Pattern: "quotas/rate-limit/" + framework.GenericNameRegex("name"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "rate-limit-quotas",
			},

			Fields: map[string]*framework.FieldSchema{
				"type": {
					Type:        framework.TypeString,
					Description: "Type of the quota rule.",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the quota rule.",
				},
				"path": {
					Type: framework.TypeString,
					Description: `Path of the mount or namespace to apply the quota. A blank path configures a
global quota. For example namespace1/ adds a quota to a full namespace,
namespace1/auth/userpass adds a quota to userpass in namespace1.`,
				},
				"role": {
					Type: framework.TypeString,
					Description: `Login role to apply this quota to. Note that when set, path must be configured
to a valid auth method with a concept of roles.`,
				},
				"inheritable": {
					Type:        framework.TypeBool,
					Description: `Whether all child namespaces can inherit this namespace quota.`,
				},
				"rate": {
					Type: framework.TypeFloat,
					Description: `The maximum number of requests in a given interval to be allowed by the quota rule.
The 'rate' must be positive.`,
				},
				"interval": {
					Type:        framework.TypeDurationSecond,
					Description: "The duration to enforce rate limiting for (default '1s').",
				},
				"block_interval": {
					Type: framework.TypeDurationSecond,
					Description: `If set, when a client reaches a rate limit threshold, the client will be prohibited
from any further requests until after the 'block_interval' has elapsed.`,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRateLimitQuotasUpdate(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "write",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: http.StatusText(http.StatusNoContent),
						}},
					},
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRateLimitQuotasRead(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "read",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"type": {
									Type:     framework.TypeString,
									Required: true,
								},
								"name": {
									Type:     framework.TypeString,
									Required: true,
								},
								"path": {
									Type:     framework.TypeString,
									Required: true,
								},
								"role": {
									Type:     framework.TypeString,
									Required: true,
								},
								"rate": {
									Type:     framework.TypeFloat,
									Required: true,
								},
								"interval": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"block_interval": {
									Type:     framework.TypeInt,
									Required: true,
								},
								"inheritable": {
									Type:     framework.TypeBool,
									Required: true,
								},
							},
						}},
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleRateLimitQuotasDelete(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "delete",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
				},
			},
			HelpSynopsis:    strings.TrimSpace(quotasHelp["rate-limit"][0]),
			HelpDescription: strings.TrimSpace(quotasHelp["rate-limit"][1]),
		},
	}
}

func (b *SystemBackend) handleQuotasConfigUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		config, err := quotas.LoadConfig(ctx, b.Core.systemBarrierView)
		if err != nil {
			return nil, err
		}

		config.EnableRateLimitAuditLogging = d.Get("enable_rate_limit_audit_logging").(bool)
		config.EnableRateLimitResponseHeaders = d.Get("enable_rate_limit_response_headers").(bool)
		config.RateLimitExemptPaths = d.Get("rate_limit_exempt_paths").([]string)

		entry, err := logical.StorageEntryJSON(quotas.ConfigPath, config)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		entry, err = logical.StorageEntryJSON(quotas.DefaultRateLimitExemptPathsToggle, true)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		b.Core.quotaManager.SetEnableRateLimitAuditLogging(config.EnableRateLimitAuditLogging)
		b.Core.quotaManager.SetEnableRateLimitResponseHeaders(config.EnableRateLimitResponseHeaders)
		b.Core.quotaManager.SetRateLimitExemptPaths(config.RateLimitExemptPaths)

		return nil, nil
	}
}

func (b *SystemBackend) handleQuotasConfigRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		config := b.Core.quotaManager.Config()
		return &logical.Response{
			Data: map[string]interface{}{
				"enable_rate_limit_audit_logging":    config.EnableRateLimitAuditLogging,
				"enable_rate_limit_response_headers": config.EnableRateLimitResponseHeaders,
				"rate_limit_exempt_paths":            config.RateLimitExemptPaths,
			},
		}, nil
	}
}

func (b *SystemBackend) handleRateLimitQuotasList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		names, err := b.Core.quotaManager.QuotaNames(quotas.TypeRateLimit)
		if err != nil {
			return nil, err
		}

		return logical.ListResponse(names), nil
	}
}

func (b *SystemBackend) handleRateLimitQuotasUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		name := d.Get("name").(string)

		qType := quotas.TypeRateLimit.String()
		rate := d.Get("rate").(float64)
		if rate <= 0 {
			return logical.ErrorResponse("'rate' is invalid"), nil
		}

		interval := time.Second * time.Duration(d.Get("interval").(int))
		if interval == 0 {
			interval = time.Second
		}

		blockInterval := time.Second * time.Duration(d.Get("block_interval").(int))
		if blockInterval < 0 {
			return logical.ErrorResponse("'block' is invalid"), nil
		}

		mountPath := sanitizePath(d.Get("path").(string))
		ns := b.Core.namespaceByPath(mountPath)
		if ns.ID != namespace.RootNamespaceID {
			mountPath = strings.TrimPrefix(mountPath, ns.Path)
		}

		var pathSuffix string
		if mountPath != "" {
			me := b.Core.router.MatchingMountEntry(namespace.ContextWithNamespace(ctx, ns), mountPath)
			if me == nil {
				return logical.ErrorResponse("invalid mount path %q", mountPath), nil
			}

			mountAPIPath := me.APIPathNoNamespace()
			pathSuffix = strings.TrimSuffix(strings.TrimPrefix(mountPath, mountAPIPath), "/")
			mountPath = mountAPIPath
		}

		role := d.Get("role").(string)
		// If this is a quota with a role, ensure the backend supports role resolution
		if role != "" {
			if pathSuffix != "" {
				return logical.ErrorResponse("Quotas cannot contain both a path suffix and a role. If a role is provided, path must be a valid auth mount with a concept of roles"), nil
			}
			authBackend := b.Core.router.MatchingBackend(namespace.ContextWithNamespace(ctx, ns), mountPath)
			if authBackend == nil || authBackend.Type() != logical.TypeCredential {
				return logical.ErrorResponse("Mount path %q is not a valid auth method and therefore unsuitable for use with role-based quotas", mountPath), nil
			}
			// We will always error as we aren't supplying real data, but we're looking for "unsupported operation" in particular
			_, err := authBackend.HandleRequest(ctx, &logical.Request{
				Path:      "login",
				Operation: logical.ResolveRoleOperation,
			})
			if err != nil && (err == logical.ErrUnsupportedOperation || err == logical.ErrUnsupportedPath) {
				return logical.ErrorResponse("Mount path %q does not support use with role-based quotas", mountPath), nil
			}
		}

		var inheritable bool
		// All global quotas should be inherited by default
		if ns.Path == "" {
			inheritable = true
		}

		if inheritableRaw, ok := d.GetOk("inheritable"); ok {
			inheritable = inheritableRaw.(bool)
			if inheritable {
				if pathSuffix != "" || role != "" || mountPath != "" {
					return logical.ErrorResponse("only namespace quotas can be configured as inheritable"), nil
				}
			} else if ns.Path == "" {
				// User should not try to configure a global quota that cannot be inherited
				return logical.ErrorResponse("all global quotas must be inheritable"), nil
			}
		}

		// User should not try to configure a global quota to be uninheritable
		if ns.Path == "" && !inheritable {
			return logical.ErrorResponse("all global quotas must be inheritable"), nil
		}

		// Disallow creation of new quota that has properties similar to an
		// existing quota.
		quotaByFactors, err := b.Core.quotaManager.QuotaByFactors(ctx, qType, ns.Path, mountPath, pathSuffix, role)
		if err != nil {
			return nil, err
		}
		if quotaByFactors != nil && quotaByFactors.QuotaName() != name {
			return logical.ErrorResponse("quota rule with similar properties exists under the name %q", quotaByFactors.QuotaName()), nil
		}

		// If a quota already exists, fetch and update it.
		quota, err := b.Core.quotaManager.QuotaByName(qType, name)
		if err != nil {
			return nil, err
		}

		switch {
		case quota == nil:
			quota = quotas.NewRateLimitQuota(name, ns.Path, mountPath, pathSuffix, role, inheritable, interval, blockInterval, rate)
		default:
			// Re-inserting the already indexed object in memdb might cause problems.
			// So, clone the object. See https://github.com/hashicorp/go-memdb/issues/76.
			clonedQuota := quota.Clone()
			rlq := clonedQuota.(*quotas.RateLimitQuota)
			rlq.NamespacePath = ns.Path
			rlq.MountPath = mountPath
			rlq.PathSuffix = pathSuffix
			rlq.Rate = rate
			rlq.Inheritable = inheritable
			rlq.Interval = interval
			rlq.BlockInterval = blockInterval
			quota = rlq
		}

		entry, err := logical.StorageEntryJSON(quotas.QuotaStoragePath(qType, name), quota)
		if err != nil {
			return nil, err
		}

		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		if err := b.Core.quotaManager.SetQuota(ctx, qType, quota, false); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (b *SystemBackend) handleRateLimitQuotasRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		name := d.Get("name").(string)
		qType := quotas.TypeRateLimit.String()

		quota, err := b.Core.quotaManager.QuotaByName(qType, name)
		if err != nil {
			return nil, err
		}
		if quota == nil {
			return nil, nil
		}

		rlq := quota.(*quotas.RateLimitQuota)

		nsPath := rlq.NamespacePath
		if rlq.NamespacePath == "root" {
			nsPath = ""
		}

		data := map[string]interface{}{
			"type":           qType,
			"name":           rlq.Name,
			"path":           nsPath + rlq.MountPath + rlq.PathSuffix,
			"role":           rlq.Role,
			"rate":           rlq.Rate,
			"inheritable":    rlq.Inheritable,
			"interval":       int(rlq.Interval.Seconds()),
			"block_interval": int(rlq.BlockInterval.Seconds()),
		}

		return &logical.Response{
			Data: data,
		}, nil
	}
}

func (b *SystemBackend) handleRateLimitQuotasDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		name := d.Get("name").(string)
		qType := quotas.TypeRateLimit.String()

		if err := req.Storage.Delete(ctx, quotas.QuotaStoragePath(qType, name)); err != nil {
			return nil, err
		}

		if err := b.Core.quotaManager.DeleteQuota(ctx, qType, name); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

var quotasHelp = map[string][2]string{
	"quotas-config": {
		"Create, update and read the quota configuration.",
		"",
	},
	"rate-limit": {
		`Get, create or update rate limit resource quota for an optional namespace or
mount.`,
		`A rate limit quota will enforce API rate limiting in a specified interval. A
rate limit quota can be created at the root level or defined on a namespace or
mount by specifying a 'path'. The rate limiter is applied to each unique client
IP address.`,
	},
	"rate-limit-list": {
		"Lists the names of all the rate limit quotas.",
		"This list contains quota definitions from all the namespaces.",
	},
}
