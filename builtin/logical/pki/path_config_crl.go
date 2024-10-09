// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/builtin/logical/pki/pki_backend"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var configCRLFields = map[string]*framework.FieldSchema{
	"expiry": {
		Type: framework.TypeString,
		Description: `The amount of time the generated CRL should be
valid; defaults to 72 hours`,
		Default: "72h",
	},
	"disable": {
		Type:        framework.TypeBool,
		Description: `If set to true, disables generating the CRL entirely.`,
	},
	"ocsp_disable": {
		Type:        framework.TypeBool,
		Description: `If set to true, ocsp unauthorized responses will be returned.`,
	},
	"ocsp_expiry": {
		Type: framework.TypeString,
		Description: `The amount of time an OCSP response will be valid (controls 
the NextUpdate field); defaults to 12 hours`,
		Default: "1h",
	},
	"auto_rebuild": {
		Type:        framework.TypeBool,
		Description: `If set to true, enables automatic rebuilding of the CRL`,
	},
	"auto_rebuild_grace_period": {
		Type:        framework.TypeString,
		Description: `The time before the CRL expires to automatically rebuild it, when enabled. Must be shorter than the CRL expiry. Defaults to 12h.`,
		Default:     "12h",
	},
	"enable_delta": {
		Type:        framework.TypeBool,
		Description: `Whether to enable delta CRLs between authoritative CRL rebuilds`,
	},
	"delta_rebuild_interval": {
		Type:        framework.TypeString,
		Description: `The time between delta CRL rebuilds if a new revocation has occurred. Must be shorter than the CRL expiry. Defaults to 15m.`,
		Default:     "15m",
	},
	"cross_cluster_revocation": {
		Type: framework.TypeBool,
		Description: `Whether to enable a global, cross-cluster revocation queue.
Must be used with auto_rebuild=true.`,
	},
	"unified_crl": {
		Type: framework.TypeBool,
		Description: `If set to true enables global replication of revocation entries,
also enabling unified versions of OCSP and CRLs if their respective features are enabled.
disable for CRLs and ocsp_disable for OCSP.`,
		Default: "false",
	},
	"unified_crl_on_existing_paths": {
		Type: framework.TypeBool,
		Description: `If set to true, 
existing CRL and OCSP paths will return the unified CRL instead of a response based on cluster-local data`,
		Default: "false",
	},
	"max_crl_entries": {
		Type:        framework.TypeInt,
		Description: `The maximum number of entries the CRL can contain.  This is meant as a guard against accidental runaway revocations overloading Vault storage.  If this limit is exceeded writing the CRL will fail.  If set to -1 this limit is disabled.`,
		Default:     pki_backend.DefaultCrlConfig.MaxCRLEntries,
	},
}

func pathConfigCRL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/crl",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
		},

		Fields: configCRLFields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "crl-configuration",
				},
				Callback: b.pathCRLRead,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields:      configCRLFields,
					}},
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathCRLWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "crl",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields:      configCRLFields,
					}},
				},
				// Read more about why these flags are set in backend.go.
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathConfigCRLHelpSyn,
		HelpDescription: pathConfigCRLHelpDesc,
	}
}

func (b *backend) pathCRLRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)

	config, err := b.CrlBuilder().getConfigWithForcedUpdate(sc)
	if err != nil {
		return nil, fmt.Errorf("failed fetching CRL config: %w", err)
	}

	return genResponseFromCrlConfig(config), nil
}

func (b *backend) pathCRLWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)
	config, err := b.CrlBuilder().getConfigWithForcedUpdate(sc)
	if err != nil {
		return nil, err
	}

	if expiryRaw, ok := d.GetOk("expiry"); ok {
		expiry := expiryRaw.(string)
		_, err := parseutil.ParseDurationSecond(expiry)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("given expiry could not be decoded: %s", err)), nil
		}
		config.Expiry = expiry
	}

	oldDisable := config.Disable
	if disableRaw, ok := d.GetOk("disable"); ok {
		config.Disable = disableRaw.(bool)
	}

	if ocspDisableRaw, ok := d.GetOk("ocsp_disable"); ok {
		config.OcspDisable = ocspDisableRaw.(bool)
	}

	if expiryRaw, ok := d.GetOk("ocsp_expiry"); ok {
		expiry := expiryRaw.(string)
		duration, err := parseutil.ParseDurationSecond(expiry)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("given ocsp_expiry could not be decoded: %s", err)), nil
		}
		if duration < 0 {
			return logical.ErrorResponse(fmt.Sprintf("ocsp_expiry must be greater than or equal to 0 got: %s", duration)), nil
		}
		config.OcspExpiry = expiry
	}

	oldAutoRebuild := config.AutoRebuild
	if autoRebuildRaw, ok := d.GetOk("auto_rebuild"); ok {
		config.AutoRebuild = autoRebuildRaw.(bool)
	}

	if autoRebuildGracePeriodRaw, ok := d.GetOk("auto_rebuild_grace_period"); ok {
		autoRebuildGracePeriod := autoRebuildGracePeriodRaw.(string)
		if _, err := parseutil.ParseDurationSecond(autoRebuildGracePeriod); err != nil {
			return logical.ErrorResponse(fmt.Sprintf("given auto_rebuild_grace_period could not be decoded: %s", err)), nil
		}
		config.AutoRebuildGracePeriod = autoRebuildGracePeriod
	}

	oldEnableDelta := config.EnableDelta
	if enableDeltaRaw, ok := d.GetOk("enable_delta"); ok {
		config.EnableDelta = enableDeltaRaw.(bool)
	}

	if deltaRebuildIntervalRaw, ok := d.GetOk("delta_rebuild_interval"); ok {
		deltaRebuildInterval := deltaRebuildIntervalRaw.(string)
		if _, err := parseutil.ParseDurationSecond(deltaRebuildInterval); err != nil {
			return logical.ErrorResponse(fmt.Sprintf("given delta_rebuild_interval could not be decoded: %s", err)), nil
		}
		config.DeltaRebuildInterval = deltaRebuildInterval
	}

	if useGlobalQueue, ok := d.GetOk("cross_cluster_revocation"); ok {
		config.UseGlobalQueue = useGlobalQueue.(bool)
	}

	oldUnifiedCRL := config.UnifiedCRL
	if unifiedCrlRaw, ok := d.GetOk("unified_crl"); ok {
		config.UnifiedCRL = unifiedCrlRaw.(bool)
	}

	if unifiedCrlOnExistingPathsRaw, ok := d.GetOk("unified_crl_on_existing_paths"); ok {
		config.UnifiedCRLOnExistingPaths = unifiedCrlOnExistingPathsRaw.(bool)
	}

	if maxCRLEntriesRaw, ok := d.GetOk("max_crl_entries"); ok {
		v := maxCRLEntriesRaw.(int)
		if v == -1 || v > 0 {
			config.MaxCRLEntries = v
		}
	}

	if config.UnifiedCRLOnExistingPaths && !config.UnifiedCRL {
		return logical.ErrorResponse("unified_crl_on_existing_paths cannot be enabled if unified_crl is disabled"), nil
	}

	expiry, _ := parseutil.ParseDurationSecond(config.Expiry)
	if config.AutoRebuild {
		gracePeriod, _ := parseutil.ParseDurationSecond(config.AutoRebuildGracePeriod)
		if gracePeriod >= expiry {
			return logical.ErrorResponse(fmt.Sprintf("CRL auto-rebuilding grace period (%v) must be strictly shorter than CRL expiry (%v) value when auto-rebuilding of CRLs is enabled", config.AutoRebuildGracePeriod, config.Expiry)), nil
		}
	}

	if config.EnableDelta {
		deltaRebuildInterval, _ := parseutil.ParseDurationSecond(config.DeltaRebuildInterval)
		if deltaRebuildInterval >= expiry {
			return logical.ErrorResponse(fmt.Sprintf("CRL delta rebuild window (%v) must be strictly shorter than CRL expiry (%v) value when delta CRLs are enabled", config.DeltaRebuildInterval, config.Expiry)), nil
		}
	}

	if !config.AutoRebuild {
		if config.EnableDelta {
			return logical.ErrorResponse("Delta CRLs cannot be enabled when auto rebuilding is disabled as the complete CRL is always regenerated!"), nil
		}

		if config.UseGlobalQueue {
			return logical.ErrorResponse("Global, cross-cluster revocation queue cannot be enabled when auto rebuilding is disabled as the local cluster may not have the certificate entry!"), nil
		}
	}

	if !constants.IsEnterprise && config.UseGlobalQueue {
		return logical.ErrorResponse("Global, cross-cluster revocation queue (cross_cluster_revocation) can only be enabled on Vault Enterprise."), nil
	}

	if !constants.IsEnterprise && config.UnifiedCRL {
		return logical.ErrorResponse("unified_crl can only be enabled on Vault Enterprise"), nil
	}

	isLocalMount := b.System().LocalMount()
	if isLocalMount && config.UseGlobalQueue {
		return logical.ErrorResponse("Global, cross-cluster revocation queue (cross_cluster_revocation) cannot be enabled on local mounts."),
			nil
	}

	if isLocalMount && config.UnifiedCRL {
		return logical.ErrorResponse("unified_crl cannot be enabled on local mounts."), nil
	}

	if !config.AutoRebuild && config.UnifiedCRL {
		return logical.ErrorResponse("unified_crl=true requires auto_rebuild=true, as unified CRLs cannot be rebuilt on every revocation."), nil
	}

	if _, err := b.CrlBuilder().writeConfig(sc, config); err != nil {
		return nil, fmt.Errorf("failed persisting CRL config: %w", err)
	}

	resp := genResponseFromCrlConfig(config)

	// Note this only affects/happens on the main cluster node, if you need to
	// notify something based on a configuration change on all server types
	// have a look at CrlBuilder::reloadConfigIfRequired
	if oldDisable != config.Disable || (oldAutoRebuild && !config.AutoRebuild) || (oldEnableDelta != config.EnableDelta) || (oldUnifiedCRL != config.UnifiedCRL) {
		// It wasn't disabled but now it is (or equivalently, we were set to
		// auto-rebuild and we aren't now or equivalently, we changed our
		// mind about delta CRLs and need a new complete one or equivalently,
		// we changed our mind about unified CRLs), rotate the CRLs.
		warnings, crlErr := b.CrlBuilder().Rebuild(sc, true)
		if crlErr != nil {
			switch crlErr.(type) {
			case errutil.UserError:
				return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
			default:
				return nil, fmt.Errorf("error encountered during CRL building: %w", crlErr)
			}
		}
		for index, warning := range warnings {
			resp.AddWarning(fmt.Sprintf("Warning %d during CRL rebuild: %v", index+1, warning))
		}
	}

	return resp, nil
}

func maxCRLEntriesOrDefault(size int) int {
	if size == 0 {
		return pki_backend.DefaultCrlConfig.MaxCRLEntries
	}
	return size
}

func genResponseFromCrlConfig(config *pki_backend.CrlConfig) *logical.Response {
	return &logical.Response{
		Data: map[string]interface{}{
			"expiry":                        config.Expiry,
			"disable":                       config.Disable,
			"ocsp_disable":                  config.OcspDisable,
			"ocsp_expiry":                   config.OcspExpiry,
			"auto_rebuild":                  config.AutoRebuild,
			"auto_rebuild_grace_period":     config.AutoRebuildGracePeriod,
			"enable_delta":                  config.EnableDelta,
			"delta_rebuild_interval":        config.DeltaRebuildInterval,
			"cross_cluster_revocation":      config.UseGlobalQueue,
			"unified_crl":                   config.UnifiedCRL,
			"unified_crl_on_existing_paths": config.UnifiedCRLOnExistingPaths,
			"max_crl_entries":               maxCRLEntriesOrDefault(config.MaxCRLEntries),
		},
	}
}

const pathConfigCRLHelpSyn = `
Configure the CRL expiration.
`

const pathConfigCRLHelpDesc = `
This endpoint allows configuration of the CRL lifetime.
`
