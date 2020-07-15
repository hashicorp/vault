package vault

import (
	"context"
	"strings"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/quotas"
)

// quotasPaths returns paths that enable quota management
func (b *SystemBackend) quotasPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "quotas/config$",
			Fields: map[string]*framework.FieldSchema{
				"enable_rate_limit_audit_logging": {
					Type:        framework.TypeBool,
					Description: "If set, starts audit logging of requests that get rejected due to rate limit quota rule violations.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleQuotasConfigUpdate(),
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleQuotasConfigRead(),
				},
			},
			HelpSynopsis:    strings.TrimSpace(quotasHelp["quotas-config"][0]),
			HelpDescription: strings.TrimSpace(quotasHelp["quotas-config"][1]),
		},
		{
			Pattern: "quotas/rate-limit/?$",
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
				"rate": {
					Type: framework.TypeFloat,
					Description: `The rate at which allowed requests are refilled per second by the quota rule.
Internally, a token-bucket algorithm is used which has a size of 'burst', initially full. The quota
limits requests to 'rate' per-second, with a maximum burst size of 'burst'. Each request takes a single
token from this bucket. The 'rate' must be positive.`,
				},
				"burst": {
					Type: framework.TypeInt,
					Description: `The maximum number of requests at any given second to be allowed by the quota
rule. There is a one-to-one mapping between requests and tokens in the rate limit quota. A client
may perform up to 'burst' requests at once, at which they they may invoke additional requests at
'rate' per-second.`,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleRateLimitQuotasUpdate(),
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleRateLimitQuotasRead(),
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleRateLimitQuotasDelete(),
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

		entry, err := logical.StorageEntryJSON(quotas.ConfigPath, config)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		b.Core.quotaManager.SetEnableRateLimitAuditLogging(config.EnableRateLimitAuditLogging)
		return nil, nil
	}
}

func (b *SystemBackend) handleQuotasConfigRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		config := b.Core.quotaManager.Config()
		return &logical.Response{
			Data: map[string]interface{}{
				"enable_rate_limit_audit_logging": config.EnableRateLimitAuditLogging,
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

		burst := d.Get("burst").(int)
		if burst < int(rate) {
			return logical.ErrorResponse("'burst' must be greater than or equal to 'rate' as an integer value"), nil
		}

		mountPath := sanitizePath(d.Get("path").(string))
		ns := b.Core.namespaceByPath(mountPath)
		if ns.ID != namespace.RootNamespaceID {
			mountPath = strings.TrimPrefix(mountPath, ns.Path)
		}

		if mountPath != "" {
			match := b.Core.router.MatchingMount(namespace.ContextWithNamespace(ctx, ns), mountPath)
			if match == "" {
				return logical.ErrorResponse("invalid mount path %q", mountPath), nil
			}
		}

		// If a quota already exists, fetch and update it.
		quota, err := b.Core.quotaManager.QuotaByName(qType, name)
		if err != nil {
			return nil, err
		}

		switch {
		case quota == nil:
			// Disallow creation of new quota that has properties similar to an
			// existing quota.
			quotaByFactors, err := b.Core.quotaManager.QuotaByFactors(ctx, qType, ns.Path, mountPath)
			if err != nil {
				return nil, err
			}
			if quotaByFactors != nil && quotaByFactors.QuotaName() != name {
				return logical.ErrorResponse("quota rule with similar properties exists under the name %q", quotaByFactors.QuotaName()), nil
			}

			quota = quotas.NewRateLimitQuota(name, ns.Path, mountPath, rate, burst)
		default:
			rlq := quota.(*quotas.RateLimitQuota)
			rlq.NamespacePath = ns.Path
			rlq.MountPath = mountPath
			rlq.Rate = rate
			rlq.Burst = burst
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
			"type":  qType,
			"name":  rlq.Name,
			"path":  nsPath + rlq.MountPath,
			"rate":  rlq.Rate,
			"burst": rlq.Burst,
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
		`A rate limit quota will enforce rate limiting using a token bucket algorithm. A
rate limit quota can be created at the root level or defined on a namespace or
mount by specifying a 'path'. The rate limiter is applied to each unique client
IP address. A client may invoke 'burst' requests at any given second, after
which they may invoke additional requests at 'rate' per-second.`,
	},
	"rate-limit-list": {
		"Lists the names of all the rate limit quotas.",
		"This list contains quota definitions from all the namespaces.",
	},
}
