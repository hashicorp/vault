// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
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
func pathAcmeEabCreateList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "acme/eab",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
		},

		Fields: map[string]*framework.FieldSchema{},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "acme-configuration",
				},
				Callback: b.pathAcmeListEab,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathAcmeCreateEab,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "acme",
				},
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
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
				Callback: b.pathAcmeDeleteEab,
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

type eabType struct {
	KeyID     string    `json:"-"`
	KeyType   string    `json:"key-type"`
	KeyBits   string    `json:"key-bits"`
	MacKey    []byte    `json:"mac-key"`
	CreatedOn time.Time `json:"created-on"`
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
			"key_bits":   eab.KeyBits,
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
	macKey, err := generateEabKey(b.GetRandomReader())
	if err != nil {
		return nil, fmt.Errorf("failed generating eab key: %w", err)
	}

	eab := &eabType{
		KeyID:     kid,
		KeyType:   "ec",
		KeyBits:   "256",
		MacKey:    macKey,
		CreatedOn: time.Now(),
	}

	sc := b.makeStorageContext(ctx, r.Storage)
	err = b.acmeState.SaveEab(sc, eab)
	if err != nil {
		return nil, fmt.Errorf("failed saving generated eab: %w", err)
	}

	encodedKey := base64.RawURLEncoding.EncodeToString(macKey)

	return &logical.Response{
		Data: map[string]interface{}{
			"id":          eab.KeyID,
			"key_type":    eab.KeyType,
			"key_bits":    eab.KeyBits,
			"private_key": encodedKey,
			"created_on":  eab.CreatedOn.Format(time.RFC3339),
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

func generateEabKey(random io.Reader) ([]byte, error) {
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), random)
	if err != nil {
		return nil, err
	}

	keyBytes, err := x509.MarshalECPrivateKey(ecKey)
	if err != nil {
		return nil, err
	}

	return keyBytes, nil
}
