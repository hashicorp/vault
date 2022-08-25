package pki

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// CRLConfig holds basic CRL configuration information
type crlConfig struct {
	Expiry                 string `json:"expiry"`
	Disable                bool   `json:"disable"`
	OcspDisable            bool   `json:"ocsp_disable"`
	AutoRebuild            bool   `json:"auto_rebuild"`
	AutoRebuildGracePeriod string `json:"auto_rebuild_grace_period"`
	OcspExpiry             string `json:"ocsp_expiry"`
}

// Implicit default values for the config if it does not exist.
var defaultCrlConfig = crlConfig{
	Expiry:                 "72h",
	Disable:                false,
	OcspDisable:            false,
	OcspExpiry:             "12h",
	AutoRebuild:            false,
	AutoRebuildGracePeriod: "12h",
}

func pathConfigCRL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/crl",
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
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathCRLRead,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathCRLWrite,
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

	return &logical.Response{
		Data: map[string]interface{}{
			"expiry":                    config.Expiry,
			"disable":                   config.Disable,
			"ocsp_disable":              config.OcspDisable,
			"ocsp_expiry":               config.OcspExpiry,
			"auto_rebuild":              config.AutoRebuild,
			"auto_rebuild_grace_period": config.AutoRebuildGracePeriod,
		},
	}, nil
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

	if config.AutoRebuild {
		expiry, _ := time.ParseDuration(config.Expiry)
		gracePeriod, _ := time.ParseDuration(config.AutoRebuildGracePeriod)

		if gracePeriod >= expiry {
			return logical.ErrorResponse(fmt.Sprintf("CRL auto-rebuilding grace period (%v) must be strictly shorter than CRL expiry (%v) value when auto-rebuilding of CRLs is enabled", config.AutoRebuildGracePeriod, config.Expiry)), nil
		}
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

	if oldDisable != config.Disable || (oldAutoRebuild && !config.AutoRebuild) {
		// It wasn't disabled but now it is (or equivalently, we were set to
		// auto-rebuild and we aren't now), so rotate the CRL.
		crlErr := b.crlBuilder.rebuild(ctx, b, req, true)
		if crlErr != nil {
			switch crlErr.(type) {
			case errutil.UserError:
				return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
			default:
				return nil, fmt.Errorf("error encountered during CRL building: %w", crlErr)
			}
		}
	}

	return nil, nil
}

const pathConfigCRLHelpSyn = `
Configure the CRL expiration.
`

const pathConfigCRLHelpDesc = `
This endpoint allows configuration of the CRL lifetime.
`
