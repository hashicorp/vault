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
	Expiry      string `json:"expiry"`
	Disable     bool   `json:"disable"`
	OcspDisable bool   `json:"ocsp_disable"`
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
			"expiry":       config.Expiry,
			"disable":      config.Disable,
			"ocsp_disable": config.OcspDisable,
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

	var oldDisable bool
	if disableRaw, ok := d.GetOk("disable"); ok {
		oldDisable = config.Disable
		config.Disable = disableRaw.(bool)
	}

	if ocspDisableRaw, ok := d.GetOk("ocsp_disable"); ok {
		config.OcspDisable = ocspDisableRaw.(bool)
	}

	err = sc.writeRevocationConfig(config)
	if err != nil {
		return nil, err
	}

	if oldDisable != config.Disable {
		// It wasn't disabled but now it is, rotate
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
