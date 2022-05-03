package transit

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"path"
	"strconv"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
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
	p, err := getWrappingKey(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if p == nil {
		p, err = generateWrappingKey(ctx, req.Storage, b.GetRandomReader())
		if err != nil {
			return nil, err
		}
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

	publicKeyString := string(pemBytes)

	resp := &logical.Response{
		Data: map[string]interface{}{
			"public_key": publicKeyString,
		},
	}

	return resp, nil
}

func getWrappingKey(ctx context.Context, storage logical.Storage) (*keysutil.Policy, error) {
	p, err := keysutil.LoadPolicy(ctx, storage, path.Join("import", "policy", WrappingKeyName))
	if err != nil {
		return nil, err
	}

	return p, nil
}

func generateWrappingKey(ctx context.Context, storage logical.Storage, rand io.Reader) (*keysutil.Policy, error) {
	p := &keysutil.Policy{
		Name:                 WrappingKeyName,
		Type:                 keysutil.KeyType_RSA4096,
		Derived:              false,
		Exportable:           false,
		AllowPlaintextBackup: false,
		AutoRotatePeriod:     0,
		StoragePrefix:        "import/",
	}

	err := p.Rotate(ctx, storage, rand)
	if err != nil {
		return nil, err
	}

	return p, nil
}

const (
	pathWrappingKeyHelpSyn  = "Returns the public key to use for wrapping imported keys"
	pathWrappingKeyHelpDesc = "This path is used to retrieve the RSA-4096 wrapping key" +
		"for wrapping keys that are being imported into transit."
)
