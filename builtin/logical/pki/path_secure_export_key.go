// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/observe"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathSecureExportCAKey registers WRITE /pki/keys/:ca-key-uuid/export.
func pathSecureExportCAKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("ca_key_uuid") + "/export$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationVerb:   "export",
			OperationSuffix: "ca-key",
		},

		Fields: map[string]*framework.FieldSchema{
			"ca_key_uuid": {
				Type:        framework.TypeString,
				Description: `UUID of the CA key to export.`,
			},
			"public_key": {
				Type:        framework.TypeString,
				Description: `PEM-encoded public key used to wrap the CA private key. Supported types: RSA (any size), EC P-256/P-384/P-521, ML-KEM-768, ML-KEM-1024.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathSecureExportCAKeyHandler,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"ca_key_uuid": {
								Type:        framework.TypeString,
								Description: `Echo of the CA key UUID that was exported.`,
								Required:    true,
							},
							"wrapped_key": {
								Type:        framework.TypeString,
								Description: `Base64-encoded JSON blob containing the wrapped CA private key. The wrapping algorithm depends on the supplied public_key type.`,
								Required:    true,
							},
							"export_key_hmac": {
								Type:        framework.TypeString,
								Description: `HMAC of the public key used to wrap the CA private key. Pass this alongside wrapped_key to the import endpoint to identify which export key to decrypt with.`,
								Required:    true,
							},
							"exported_at": {
								Type:        framework.TypeString,
								Description: `RFC3339 timestamp of when the export occurred.`,
								Required:    true,
							},
						},
					}},
				},
				// This is a root-only, highly-privileged write — never serve from standby or secondary.
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathSecureExportCAKeyHelpSyn,
		HelpDescription: pathSecureExportCAKeyHelpDesc,
	}
}

// pathSecureExportCAKeyHandler handles WRITE /pki/keys/:ca-key-uuid/export.
func (b *backend) pathSecureExportCAKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	caKeyUUID := data.Get("ca_key_uuid").(string)
	if caKeyUUID == "" {
		return logical.ErrorResponse("ca_key_uuid is required"), nil
	}

	pubKeyPEM := data.Get("public_key").(string)
	if pubKeyPEM == "" {
		return logical.ErrorResponse("public_key is required"), nil
	}

	// Fetch and validate the CA key exists before doing any crypto.
	keyEntry, err := issuing.FetchKeyById(ctx, req.Storage, issuing.KeyID(caKeyUUID))
	if err != nil {
		if _, ok := err.(errutil.UserError); ok {
			return logical.ErrorResponse(err.Error()), nil
		}
		return nil, err
	}

	if keyEntry.IsManagedPrivateKey() {
		return logical.ErrorResponse("secure export is not supported for managed (KMS) keys"), nil
	}

	blobJSON, err := wrapCAPrivateKey(keyEntry.PrivateKey, pubKeyPEM)
	if err != nil {
		if _, ok := err.(errutil.UserError); ok {
			return logical.ErrorResponse(err.Error()), nil
		}
		return nil, err
	}

	// Compute the HMAC over the public key DER bytes so the import side can
	// identify which export key was used. PEM validity is already guaranteed by
	// wrapCAPrivateKey above.
	pubKeyHMAC := computeExportKeyHMACFromPEM(pubKeyPEM)
	now := time.Now().UTC()

	// Persist an export audit record keyed on the SPKI fingerprint so the
	// record survives a delete-and-reimport of this key entry.
	spkiFingerprint, err := caPublicKeyFingerprint(keyEntry.PrivateKey)
	if err != nil {
		return nil, err
	}

	record := CAKeyExportRecord{
		ExportKeyHMAC: pubKeyHMAC,
		ExportedAt:    now,
	}
	if err := writeCAKeyExportRecord(ctx, req.Storage, spkiFingerprint, record); err != nil {
		return nil, err
	}

	b.pkiObserver.RecordPKIObservation(ctx, req, observe.ObservationTypePKICAKeyExport,
		observe.NewAdditionalPKIMetadata("ca_key_uuid", caKeyUUID),
		observe.NewAdditionalPKIMetadata("export_key_hmac", pubKeyHMAC),
	)

	b.Logger().Warn("CA private key material was securely exported via BYOK; if this is unexpected, your key has been compromised.",
		"ca_key_uuid", caKeyUUID,
		"ca_key_name", keyEntry.Name,
		"ca_key_fingerprint", spkiFingerprint,
		"export_key_hmac", pubKeyHMAC,
	)

	return &logical.Response{
		Data: map[string]interface{}{
			"ca_key_uuid":     caKeyUUID,
			"wrapped_key":     base64.StdEncoding.EncodeToString(blobJSON),
			"export_key_hmac": pubKeyHMAC,
			"exported_at":     now.Format(time.RFC3339),
		},
	}, nil
}

const (
	pathSecureExportCAKeyHelpSyn  = `Securely wrap a CA private key with a caller-supplied public key and return the encrypted blob.`
	pathSecureExportCAKeyHelpDesc = `
WRITE: Fetches the private key for the given CA key UUID, encrypts it with the
provided public_key, and returns the wrapped blob. The raw private key is never
present in the response or in any Vault log.

The wrapped_key field is a base64-encoded JSON envelope. The wrapping algorithm
is selected automatically based on the public_key type:
  - RSA:              RSA-OAEP-SHA256 + AES-256-GCM
  - EC P-256/384/521: ECDH-ES + HKDF-SHA256 + AES-256-GCM
  - ML-KEM-768/1024:  KEM encapsulate + HKDF-SHA256 + AES-256-GCM

Vault does not retain the wrapped blob; the caller is responsible for storing it
until it is imported on the destination Vault instance.
`
)
