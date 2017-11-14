package transit

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	exportTypeEncryptionKey = "encryption-key"
	exportTypeSigningKey    = "signing-key"
	exportTypeHMACKey       = "hmac-key"
)

func (b *backend) pathExportKeys() *framework.Path {
	return &framework.Path{
		Pattern: "export/" + framework.GenericNameRegex("type") + "/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("version"),
		Fields: map[string]*framework.FieldSchema{
			"type": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Type of key to export (encryption-key, signing-key, hmac-key)",
			},
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
			"version": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Version of the key",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathPolicyExportRead,
		},

		HelpSynopsis:    pathExportHelpSyn,
		HelpDescription: pathExportHelpDesc,
	}
}

func (b *backend) pathPolicyExportRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	exportType := d.Get("type").(string)
	name := d.Get("name").(string)
	version := d.Get("version").(string)

	switch exportType {
	case exportTypeEncryptionKey:
	case exportTypeSigningKey:
	case exportTypeHMACKey:
	default:
		return logical.ErrorResponse(fmt.Sprintf("invalid export type: %s", exportType)), logical.ErrInvalidRequest
	}

	p, lock, err := b.lm.GetPolicyShared(req.Storage, name)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	if !p.Exportable {
		return logical.ErrorResponse("key is not exportable"), nil
	}

	switch exportType {
	case exportTypeEncryptionKey:
		if !p.Type.EncryptionSupported() {
			return logical.ErrorResponse("encryption not supported for the key"), logical.ErrInvalidRequest
		}
	case exportTypeSigningKey:
		if !p.Type.SigningSupported() {
			return logical.ErrorResponse("signing not supported for the key"), logical.ErrInvalidRequest
		}
	}

	retKeys := map[string]string{}
	switch version {
	case "":
		for k, v := range p.Keys {
			exportKey, err := getExportKey(p, &v, exportType)
			if err != nil {
				return nil, err
			}
			retKeys[strconv.Itoa(k)] = exportKey
		}

	default:
		var versionValue int
		if version == "latest" {
			versionValue = p.LatestVersion
		} else {
			version = strings.TrimPrefix(version, "v")
			versionValue, err = strconv.Atoi(version)
			if err != nil {
				return logical.ErrorResponse("invalid key version"), logical.ErrInvalidRequest
			}
		}

		if versionValue < p.MinDecryptionVersion {
			return logical.ErrorResponse("version for export is below minimun decryption version"), logical.ErrInvalidRequest
		}
		key, ok := p.Keys[versionValue]
		if !ok {
			return logical.ErrorResponse("version does not exist or cannot be found"), logical.ErrInvalidRequest
		}

		exportKey, err := getExportKey(p, &key, exportType)
		if err != nil {
			return nil, err
		}

		retKeys[strconv.Itoa(versionValue)] = exportKey
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"name": p.Name,
			"type": p.Type.String(),
			"keys": retKeys,
		},
	}

	return resp, nil
}

func getExportKey(policy *keysutil.Policy, key *keysutil.KeyEntry, exportType string) (string, error) {
	if policy == nil {
		return "", errors.New("nil policy provided")
	}

	switch exportType {
	case exportTypeHMACKey:
		return strings.TrimSpace(base64.StdEncoding.EncodeToString(key.HMACKey)), nil

	case exportTypeEncryptionKey:
		switch policy.Type {
		case keysutil.KeyType_AES256_GCM96:
			return strings.TrimSpace(base64.StdEncoding.EncodeToString(key.Key)), nil

		case keysutil.KeyType_RSA2048, keysutil.KeyType_RSA4096:
			return encodeRSAPrivateKey(key.RSAKey), nil
		}

	case exportTypeSigningKey:
		switch policy.Type {
		case keysutil.KeyType_ECDSA_P256:
			ecKey, err := keyEntryToECPrivateKey(key, elliptic.P256())
			if err != nil {
				return "", err
			}
			return ecKey, nil

		case keysutil.KeyType_ED25519:
			return strings.TrimSpace(base64.StdEncoding.EncodeToString(key.Key)), nil

		case keysutil.KeyType_RSA2048, keysutil.KeyType_RSA4096:
			return encodeRSAPrivateKey(key.RSAKey), nil
		}
	}

	return "", fmt.Errorf("unknown key type %v", policy.Type)
}

func encodeRSAPrivateKey(key *rsa.PrivateKey) string {
	// When encoding PKCS1, the PEM header should be `RSA PRIVATE KEY`. When Go
	// has PKCS8 encoding support, we may want to change this.
	derBytes := x509.MarshalPKCS1PrivateKey(key)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derBytes,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
	return string(pemBytes)
}

func keyEntryToECPrivateKey(k *keysutil.KeyEntry, curve elliptic.Curve) (string, error) {
	if k == nil {
		return "", errors.New("nil KeyEntry provided")
	}

	privKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     k.EC_X,
			Y:     k.EC_Y,
		},
		D: k.EC_D,
	}
	ecder, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return "", err
	}
	if ecder == nil {
		return "", errors.New("No data returned when marshalling to private key")
	}

	block := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: ecder,
	}
	return strings.TrimSpace(string(pem.EncodeToMemory(&block))), nil
}

const pathExportHelpSyn = `Export named encryption or signing key`

const pathExportHelpDesc = `
This path is used to export the named keys that are configured as
exportable.
`
