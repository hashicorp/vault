// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

/*
 * This file unlike the other path_acme_xxx.go are VAULT APIs to manage the
 * ACME External Account Bindings, this isn't providing any APIs that an ACME
 * client would use.
 */
func pathAcmeEabList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "acme/eab/?$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
		},

		Fields: map[string]*framework.FieldSchema{},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "list-eab-key",
					OperationSuffix: "acme",
				},
				Callback: b.pathAcmeListEab,
			},
		},

		HelpSynopsis:    "list external account bindings to be used for ACME",
		HelpDescription: `list identifiers that have been generated but yet to be used.`,
	}
}

func pathAcmeEabCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "acme/new-eab",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
		},

		Fields: map[string]*framework.FieldSchema{},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathAcmeCreateEab,
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "generate-eab-key",
					OperationSuffix: "acme",
				},
			},
		},

		HelpSynopsis: "Generate or list external account bindings to be used for ACME",
		HelpDescription: `Generate single use id/key pairs to be used for ACME EAB or list 
identifiers that have been generated but yet to be used.`,
	}
}

func pathAcmeEabDelete(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "acme/eab/" + uuidNameRegex("key_id"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
		},

		Fields: map[string]*framework.FieldSchema{
			"key_id": {
				Type:        framework.TypeString,
				Description: "EAB key identifier",
				Required:    true,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.DeleteOperation: &framework.PathOperation{
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "acme-configuration",
				},
				Callback:                    b.pathAcmeDeleteEab,
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
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
	KeyID        string    `json:"-"`
	KeyType      string    `json:"key-type"`
	KeySize      int       `json:"key-size"`
	PrivateBytes []byte    `json:"private-bytes"`
	CreatedOn    time.Time `json:"created-on"`
}

func (b *backend) pathAcmeListEab(ctx context.Context, r *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, r.Storage)

	eabIds, err := b.acmeState.ListEabIds(sc)
	if err != nil {
		return nil, err
	}

	var warnings []string
	var keyIds []string
	keyInfos := map[string]interface{}{}

	for _, eabKey := range eabIds {
		eab, err := b.acmeState.LoadEab(sc, eabKey)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("failed loading eab entry %s: %v", eabKey, err))
			continue
		}

		keyIds = append(keyIds, eab.KeyID)
		keyInfos[eab.KeyID] = map[string]interface{}{
			"key_type":   eab.KeyType,
			"key_size":   eab.KeySize,
			"created_on": eab.CreatedOn.Format(time.RFC3339),
		}
	}

	resp := logical.ListResponseWithInfo(keyIds, keyInfos)
	for _, warning := range warnings {
		resp.AddWarning(warning)
	}
	return resp, nil
}

func (b *backend) pathAcmeCreateEab(ctx context.Context, r *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	kid := genUuid()
	size := 32
	bytes, err := uuid.GenerateRandomBytesWithReader(size, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed generating eab key: %w", err)
	}

	eab := &eabType{
		KeyID:        kid,
		KeyType:      "HS",
		KeySize:      size,
		PrivateBytes: bytes,
		CreatedOn:    time.Now(),
	}

	sc := b.makeStorageContext(ctx, r.Storage)
	err = b.acmeState.SaveEab(sc, eab)
	if err != nil {
		return nil, fmt.Errorf("failed saving generated eab: %w", err)
	}

	encodedKey := base64.RawURLEncoding.EncodeToString(eab.PrivateBytes)

	return &logical.Response{
		Data: map[string]interface{}{
			"id":         eab.KeyID,
			"key_type":   eab.KeyType,
			"key_size":   eab.KeySize,
			"key":        encodedKey,
			"created_on": eab.CreatedOn.Format(time.RFC3339),
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

	deleted, err := b.acmeState.DeleteEab(sc, keyId)
	if err != nil {
		return nil, fmt.Errorf("failed deleting key id: %w", err)
	}

	resp := &logical.Response{}
	if !deleted {
		resp.AddWarning("No key id found with id: " + keyId)
	}
	return resp, nil
}
