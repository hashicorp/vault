package pki

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathRevoke(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `revoke`,
		Fields: map[string]*framework.FieldSchema{
			"serial_number": {
				Type: framework.TypeString,
				Description: `Certificate serial number, in colon- or
hyphen-separated octal`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.metricsWrap("revoke", noRole, b.pathRevokeWrite),
				// This should never be forwarded. See backend.go for more information.
				// If this needs to write, the entire request will be forwarded to the
				// active node of the current performance cluster, but we don't want to
				// forward invalid revoke requests there.
			},
		},

		HelpSynopsis:    pathRevokeHelpSyn,
		HelpDescription: pathRevokeHelpDesc,
	}
}

func pathRotateCRL(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `crl/rotate`,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathRotateCRLRead,
				// See backend.go; we will read a lot of data prior to calling write,
				// so this request should be forwarded when it is first seen, not
				// when it is ready to write.
				ForwardPerformanceStandby: true,
			},
		},

		HelpSynopsis:    pathRotateCRLHelpSyn,
		HelpDescription: pathRotateCRLHelpDesc,
	}
}

func (b *backend) pathRevokeWrite(ctx context.Context, req *logical.Request, data *framework.FieldData, _ *roleEntry) (*logical.Response, error) {
	serial := data.Get("serial_number").(string)
	if len(serial) == 0 {
		return logical.ErrorResponse("The serial number must be provided"), nil
	}

	if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) {
		return nil, logical.ErrReadOnly
	}

	// We store and identify by lowercase colon-separated hex, but other
	// utilities use dashes and/or uppercase, so normalize
	serial = strings.Replace(strings.ToLower(serial), "-", ":", -1)

	b.revokeStorageLock.Lock()
	defer b.revokeStorageLock.Unlock()

	return revokeCert(ctx, b, req, serial, false)
}

func (b *backend) pathRotateCRLRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	b.revokeStorageLock.RLock()
	defer b.revokeStorageLock.RUnlock()

	crlErr := b.crlBuilder.rebuild(ctx, b, req, false)
	if crlErr != nil {
		switch crlErr.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
		default:
			return nil, fmt.Errorf("error encountered during CRL building: %w", crlErr)
		}
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"success": true,
		},
	}, nil
}

const pathRevokeHelpSyn = `
Revoke a certificate by serial number.
`

const pathRevokeHelpDesc = `
This allows certificates to be revoked using its serial number. A root token is required.
`

const pathRotateCRLHelpSyn = `
Force a rebuild of the CRL.
`

const pathRotateCRLHelpDesc = `
Force a rebuild of the CRL. This can be used to remove expired certificates from it if no certificates have been revoked. A root token is required.
`
