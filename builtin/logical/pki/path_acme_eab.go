// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var decodedTokenPrefix = mustBase64Decode("vault-eab-0-")

func mustBase64Decode(s string) []byte {
	bytes, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		panic(fmt.Sprintf("Token prefix value: %s failed decoding: %v", s, err))
	}

	// Should be dividable by 3 otherwise our prefix will not be properly honored.
	if len(bytes)%3 != 0 {
		panic(fmt.Sprintf("Token prefix value: %s is not dividable by 3, will not prefix properly", s))
	}
	return bytes
}

/*
 * This file unlike the other path_acme_xxx.go are VAULT APIs to manage the
 * ACME External Account Bindings, this isn't providing any APIs that an ACME
 * client would use.
 */
func pathAcmeEabList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "eab/?$",
		Fields:  map[string]*framework.FieldSchema{},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathAcmeListEab,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationPrefix: operationPrefixPKI,
					OperationVerb:   "list-eab-keys",
					Description:     "List all eab key identifiers yet to be used.",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"keys": {
								Type:        framework.TypeStringSlice,
								Description: `A list of unused eab keys`,
								Required:    true,
							},
							"key_info": {
								Type:        framework.TypeMap,
								Description: `EAB details keyed by the eab key id`,
								Required:    false,
							},
						},
					}},
				},
			},
		},

		HelpSynopsis:    "list external account bindings to be used for ACME",
		HelpDescription: `list identifiers that have been generated but yet to be used.`,
	}
}

func pathAcmeNewEab(b *backend, baseUrl string) *framework.Path {
	return patternAcmeNewEab(b, baseUrl+"/new-eab")
}

func patternAcmeNewEab(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)

	opSuffix := getAcmeOperationSuffix(pattern)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathAcmeCreateEab,
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationPrefix: operationPrefixPKI,
					OperationVerb:   "generate-eab-key",
					OperationSuffix: opSuffix,
					Description:     "Generate an ACME EAB token for a directory",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"id": {
								Type:        framework.TypeString,
								Description: `The EAB key identifier`,
								Required:    true,
							},
							"key_type": {
								Type:        framework.TypeString,
								Description: `The EAB key type`,
								Required:    true,
							},
							"key": {
								Type:        framework.TypeString,
								Description: `The EAB hmac key`,
								Required:    true,
							},
							"acme_directory": {
								Type:        framework.TypeString,
								Description: `The ACME directory to which the key belongs`,
								Required:    true,
							},
							"created_on": {
								Type:        framework.TypeTime,
								Description: `An RFC3339 formatted date time when the EAB token was created`,
								Required:    true,
							},
						},
					}},
				},
			},
		},

		HelpSynopsis:    "Generate external account bindings to be used for ACME",
		HelpDescription: `Generate single use id/key pairs to be used for ACME EAB.`,
	}
}

func pathAcmeEabDelete(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "eab/" + uuidNameRegex("key_id"),

		Fields: map[string]*framework.FieldSchema{
			"key_id": {
				Type:        framework.TypeString,
				Description: "EAB key identifier",
				Required:    true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.DeleteOperation: &framework.PathOperation{
				Callback:                    b.pathAcmeDeleteEab,
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationPrefix: operationPrefixPKI,
					OperationVerb:   "delete-eab-key",
					Description:     "Delete an unused EAB token",
				},
			},
		},

		HelpSynopsis: "Delete an external account binding id prior to its use within an ACME account",
		HelpDescription: `Allows an operator to delete an external account binding,
before its bound to a new ACME account. If the identifier provided does not exist or 
was already consumed by an ACME account a successful response is returned along with 
a warning that it did not exist.`,
	}
}

type eabType struct {
	KeyID         string    `json:"key-id"`
	KeyType       string    `json:"key-type"`
	PrivateBytes  []byte    `json:"private-bytes"`
	AcmeDirectory string    `json:"acme-directory"`
	CreatedOn     time.Time `json:"created-on"`
}

func (b *backend) pathAcmeListEab(ctx context.Context, r *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, r.Storage)

	acmeState := b.GetAcmeState()
	eabIds, err := acmeState.ListEabIds(sc)
	if err != nil {
		return nil, err
	}

	var warnings []string
	var keyIds []string
	keyInfos := map[string]interface{}{}

	for _, eabKey := range eabIds {
		eab, err := acmeState.LoadEab(sc, eabKey)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("failed loading eab entry %s: %v", eabKey, err))
			continue
		}

		keyIds = append(keyIds, eab.KeyID)
		keyInfos[eab.KeyID] = map[string]interface{}{
			"key_type":       eab.KeyType,
			"acme_directory": path.Join(eab.AcmeDirectory, "directory"),
			"created_on":     eab.CreatedOn.Format(time.RFC3339),
		}
	}

	resp := logical.ListResponseWithInfo(keyIds, keyInfos)
	for _, warning := range warnings {
		resp.AddWarning(warning)
	}
	return resp, nil
}

func (b *backend) pathAcmeCreateEab(ctx context.Context, r *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	kid := genUuid()
	size := 32
	bytes, err := uuid.GenerateRandomBytesWithReader(size, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed generating eab key: %w", err)
	}

	acmeDirectory, err := getAcmeDirectory(r)
	if err != nil {
		return nil, err
	}

	eab := &eabType{
		KeyID:         kid,
		KeyType:       "hs",
		PrivateBytes:  append(decodedTokenPrefix, bytes...), // we do this to avoid generating tokens that start with -
		AcmeDirectory: acmeDirectory,
		CreatedOn:     time.Now(),
	}

	sc := b.makeStorageContext(ctx, r.Storage)
	err = b.GetAcmeState().SaveEab(sc, eab)
	if err != nil {
		return nil, fmt.Errorf("failed saving generated eab: %w", err)
	}

	encodedKey := base64.RawURLEncoding.EncodeToString(eab.PrivateBytes)

	return &logical.Response{
		Data: map[string]interface{}{
			"id":             eab.KeyID,
			"key_type":       eab.KeyType,
			"key":            encodedKey,
			"acme_directory": path.Join(eab.AcmeDirectory, "directory"),
			"created_on":     eab.CreatedOn.Format(time.RFC3339),
		},
	}, nil
}

func (b *backend) pathAcmeDeleteEab(ctx context.Context, r *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, r.Storage)
	keyId := d.Get("key_id").(string)

	_, err := uuid.ParseUUID(keyId)
	if err != nil {
		return nil, fmt.Errorf("badly formatted key_id field")
	}

	deleted, err := b.GetAcmeState().DeleteEab(sc, keyId)
	if err != nil {
		return nil, fmt.Errorf("failed deleting key id: %w", err)
	}

	resp := &logical.Response{}
	if !deleted {
		resp.AddWarning("No key id found with id: " + keyId)
	}
	return resp, nil
}

// getAcmeOperationSuffix used mainly to compute the OpenAPI spec suffix value to distinguish
// different versions of ACME Vault APIs based on directory paths
func getAcmeOperationSuffix(pattern string) string {
	hasRole := strings.Contains(pattern, framework.GenericNameRegex("role"))
	hasIssuer := strings.Contains(pattern, framework.GenericNameRegex(issuerRefParam))

	switch {
	case hasRole && hasIssuer:
		return "for-issuer-and-role"
	case hasRole:
		return "for-role"
	case hasIssuer:
		return "for-issuer"
	default:
		return ""
	}
}
