// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const WrappingKeyName = "wrapping-key"

func (b *backend) pathWrappingKey() *framework.Path {
	return &framework.Path{
		Pattern: "wrapping_key",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationSuffix: "wrapping-key",
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathWrappingKeyRead,
		},
		HelpSynopsis:    pathWrappingKeyHelpSyn,
		HelpDescription: pathWrappingKeyHelpDesc,
	}
}

func (b *backend) pathWrappingKeyRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	p, err := b.getWrappingKey(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	wrappingKey := p.Keys[strconv.Itoa(p.LatestVersion)]

	derBytes, err := x509.MarshalPKIXPublicKey(wrappingKey.RSAKey.Public())
	if err != nil {
		return nil, fmt.Errorf("error marshaling RSA public key: %w", err)
	}
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
	if pemBytes == nil || len(pemBytes) == 0 {
		return nil, fmt.Errorf("failed to PEM-encode RSA public key")
	}

	publicKeyString := string(pemBytes)

	resp := &logical.Response{
		Data: map[string]interface{}{
			"public_key": publicKeyString,
		},
	}

	return resp, nil
}

func (b *backend) getWrappingKey(ctx context.Context, storage logical.Storage) (*keysutil.Policy, error) {
	polReq := keysutil.PolicyRequest{
		Upsert:               true,
		Storage:              storage,
		Name:                 fmt.Sprintf("import/%s", WrappingKeyName),
		KeyType:              keysutil.KeyType_RSA4096,
		Derived:              false,
		Convergent:           false,
		Exportable:           false,
		AllowPlaintextBackup: false,
		AutoRotatePeriod:     0,
	}
	p, _, err := b.GetPolicy(ctx, polReq, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("error retrieving wrapping key: returned policy was nil")
	}
	if b.System().CachingDisabled() {
		p.Unlock()
	}

	return p, nil
}

const (
	pathWrappingKeyHelpSyn  = "Returns the public key to use for wrapping imported keys"
	pathWrappingKeyHelpDesc = "This path is used to retrieve the RSA-4096 wrapping key " +
		"for wrapping keys that are being imported into transit."
)
