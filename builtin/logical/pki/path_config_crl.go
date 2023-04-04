// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const latestCrlConfigVersion = 1

// CRLConfig holds basic CRL configuration information
type crlConfig struct {
	Version                   int    `json:"version"`
	Expiry                    string `json:"expiry"`
	Disable                   bool   `json:"disable"`
	OcspDisable               bool   `json:"ocsp_disable"`
	AutoRebuild               bool   `json:"auto_rebuild"`
	AutoRebuildGracePeriod    string `json:"auto_rebuild_grace_period"`
	OcspExpiry                string `json:"ocsp_expiry"`
	EnableDelta               bool   `json:"enable_delta"`
	DeltaRebuildInterval      string `json:"delta_rebuild_interval"`
	UseGlobalQueue            bool   `json:"cross_cluster_revocation"`
	UnifiedCRL                bool   `json:"unified_crl"`
	UnifiedCRLOnExistingPaths bool   `json:"unified_crl_on_existing_paths"`
}

// Implicit default values for the config if it does not exist.
var defaultCrlConfig = crlConfig{
	Version:                   latestCrlConfigVersion,
	Expiry:                    "72h",
	Disable:                   false,
	OcspDisable:               false,
	OcspExpiry:                "12h",
	AutoRebuild:               false,
	AutoRebuildGracePeriod:    "12h",
	EnableDelta:               false,
	DeltaRebuildInterval:      "15m",
	UseGlobalQueue:            false,
	UnifiedCRL:                false,
	UnifiedCRLOnExistingPaths: false,
}

func pathConfigCRL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/crl",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
		},

		Fields: map[string]*framework.FieldSchema{
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
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "crl-configuration",
				},
				Callback: b.pathCRLRead,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"expiry": {
								Type: framework.TypeString,
								Description: `The amount of time the generated CRL should be
valid; defaults to 72 hours`,
								Required: true,
							},
							"disable": {
								Type:        framework.TypeBool,
								Description: `If set to true, disables generating the CRL entirely.`,
								Required:    true,
							},
							"ocsp_disable": {
								Type:        framework.TypeBool,
								Description: `If set to true, ocsp unauthorized responses will be returned.`,
								Required:    true,
							},
							"ocsp_expiry": {
								Type: framework.TypeString,
								Description: `The amount of time an OCSP response will be valid (controls 
the NextUpdate field); defaults to 12 hours`,
								Required: true,
							},
							"auto_rebuild": {
								Type:        framework.TypeBool,
								Description: `If set to true, enables automatic rebuilding of the CRL`,
								Required:    true,
							},
							"auto_rebuild_grace_period": {
								Type:        framework.TypeString,
								Description: `The time before the CRL expires to automatically rebuild it, when enabled. Must be shorter than the CRL expiry. Defaults to 12h.`,
								Required:    true,
							},
							"enable_delta": {
								Type:        framework.TypeBool,
								Description: `Whether to enable delta CRLs between authoritative CRL rebuilds`,
								Required:    true,
							},
							"delta_rebuild_interval": {
								Type:        framework.TypeString,
								Description: `The time between delta CRL rebuilds if a new revocation has occurred. Must be shorter than the CRL expiry. Defaults to 15m.`,
								Required:    true,
							},
							"cross_cluster_revocation": {
								Type: framework.TypeBool,
								Description: `Whether to enable a global, cross-cluster revocation queue.
Must be used with auto_rebuild=true.`,
								Required: true,
							},
							"unified_crl": {
								Type: framework.TypeBool,
								Description: `If set to true enables global replication of revocation entries,
also enabling unified versions of OCSP and CRLs if their respective features are enabled.
disable for CRLs and ocsp_disable for OCSP.`,
								Required: true,
							},
							"unified_crl_on_existing_paths": {
								Type: framework.TypeBool,
								Description: `If set to true, 
existing CRL and OCSP paths will return the unified CRL instead of a response based on cluster-local data`,
								Required: true,
							},
						},
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
						Fields: map[string]*framework.FieldSchema{
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
								Required: false,
							},
							"unified_crl": {
								Type: framework.TypeBool,
								Description: `If set to true enables global replication of revocation entries,
also enabling unified versions of OCSP and CRLs if their respective features are enabled.
disable for CRLs and ocsp_disable for OCSP.`,
								Required: false,
							},
							"unified_crl_on_existing_paths": {
								Type: framework.TypeBool,
								Description: `If set to true, 
existing CRL and OCSP paths will return the unified CRL instead of a response based on cluster-local data`,
								Required: false,
							},
						},
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
	config, err := sc.getRevocationConfig()
	if err != nil {
		return nil, err
	}

	return genResponseFromCrlConfig(config), nil
}

func (b *backend) pathCRLWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)
	config, err := sc.getRevocationConfig()
	if err != nil {
		return nil, err
	}

	if expiryRaw, ok := d.GetOk("expiry"); ok {
		expiry := expiryRaw.(string)
		_, err := time.ParseDuration(expiry)
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
		duration, err := time.ParseDuration(expiry)
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
		if _, err := time.ParseDuration(autoRebuildGracePeriod); err != nil {
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
		if _, err := time.ParseDuration(deltaRebuildInterval); err != nil {
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

	if config.UnifiedCRLOnExistingPaths && !config.UnifiedCRL {
		return logical.ErrorResponse("unified_crl_on_existing_paths cannot be enabled if unified_crl is disabled"), nil
	}

	expiry, _ := time.ParseDuration(config.Expiry)
	if config.AutoRebuild {
		gracePeriod, _ := time.ParseDuration(config.AutoRebuildGracePeriod)
		if gracePeriod >= expiry {
			return logical.ErrorResponse(fmt.Sprintf("CRL auto-rebuilding grace period (%v) must be strictly shorter than CRL expiry (%v) value when auto-rebuilding of CRLs is enabled", config.AutoRebuildGracePeriod, config.Expiry)), nil
		}
	}

	if config.EnableDelta {
		deltaRebuildInterval, _ := time.ParseDuration(config.DeltaRebuildInterval)
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

	entry, err := logical.StorageEntryJSON("config/crl", config)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	b.crlBuilder.markConfigDirty()
	b.crlBuilder.reloadConfigIfRequired(sc)

	// Note this only affects/happens on the main cluster node, if you need to
	// notify something based on a configuration change on all server types
	// have a look at crlBuilder::reloadConfigIfRequired
	if oldDisable != config.Disable || (oldAutoRebuild && !config.AutoRebuild) || (oldEnableDelta != config.EnableDelta) || (oldUnifiedCRL != config.UnifiedCRL) {
		// It wasn't disabled but now it is (or equivalently, we were set to
		// auto-rebuild and we aren't now or equivalently, we changed our
		// mind about delta CRLs and need a new complete one or equivalently,
		// we changed our mind about unified CRLs), rotate the CRLs.
		crlErr := b.crlBuilder.rebuild(sc, true)
		if crlErr != nil {
			switch crlErr.(type) {
			case errutil.UserError:
				return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
			default:
				return nil, fmt.Errorf("error encountered during CRL building: %w", crlErr)
			}
		}
	}

	return genResponseFromCrlConfig(config), nil
}

func genResponseFromCrlConfig(config *crlConfig) *logical.Response {
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
		},
	}
}

const pathConfigCRLHelpSyn = `
Configure the CRL expiration.
`

const pathConfigCRLHelpDesc = `
This endpoint allows configuration of the CRL lifetime.
`
