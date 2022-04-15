package transit

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"strconv"
)

const WrappingKeyName = "wrapping-key"

func (b *backend) pathWrappingKey() *framework.Path {
	return &framework.Path{
		Pattern: "wrapping_key",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathWrappingKeyRead,
		},
		HelpSynopsis:    pathWrappingKeyHelpSyn,
		HelpDescription: pathWrappingKeyHelpDesc,
	}
}

func (b *backend) pathWrappingKeyRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	p, err := b.lm.GetWrappingKey(ctx, req.Storage, WrappingKeyName, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("error generating wrapping key: returned policy was nil")
	}
	if b.System().CachingDisabled() {
		p.Unlock()
	}

	rsaPublicKey := p.Keys[strconv.Itoa(p.LatestVersion)]

	derBytes, err := x509.MarshalPKIXPublicKey(rsaPublicKey.RSAKey.Public())
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

	resp := &logical.Response{
		Data: map[string]interface{}{
			"public_key": string(pemBytes),
		},
	}

	return resp, nil
}

const pathWrappingKeyHelpSyn = "Returns the public key to use for wrapping imported keys"
const pathWrappingKeyHelpDesc = "This path is used to retrieve the RSA-4096 wrapping key" +
	"for wrapping keys that are being imported into transit."
