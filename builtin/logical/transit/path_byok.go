// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathBYOKExportKeys() *framework.Path {
	return &framework.Path{
		Pattern: "byok-export/" + framework.GenericNameRegex("destination") + "/" + framework.GenericNameRegex("source") + framework.OptionalParamRegex("version"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "byok",
			OperationSuffix: "key|key-version",
		},

		Fields: map[string]*framework.FieldSchema{
			"destination": {
				Type:        framework.TypeString,
				Description: "Destination key to export to; usually the public wrapping key of another Transit instance.",
			},
			"source": {
				Type:        framework.TypeString,
				Description: "Source key to export; could be any present key within Transit.",
			},
			"version": {
				Type:        framework.TypeString,
				Description: "Optional version of the key to export, else all key versions are exported.",
			},
			"hash": {
				Type:        framework.TypeString,
				Description: "Hash function to use for inner OAEP encryption. Defaults to SHA256.",
				Default:     "SHA256",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathPolicyBYOKExportRead,
		},

		HelpSynopsis:    pathBYOKExportHelpSyn,
		HelpDescription: pathBYOKExportHelpDesc,
	}
}

func (b *backend) pathPolicyBYOKExportRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	dst := d.Get("destination").(string)
	src := d.Get("source").(string)
	version := d.Get("version").(string)
	hash := d.Get("hash").(string)

	dstP, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    dst,
	}, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if dstP == nil {
		return nil, fmt.Errorf("no such destination key to export to")
	}
	if !b.System().CachingDisabled() {
		dstP.Lock(false)
	}
	defer dstP.Unlock()

	srcP, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    src,
	}, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if srcP == nil {
		return nil, fmt.Errorf("no such source key for export")
	}
	if !b.System().CachingDisabled() {
		srcP.Lock(false)
	}
	defer srcP.Unlock()

	if !srcP.Exportable {
		return logical.ErrorResponse("key is not exportable"), nil
	}

	if srcP.Type.IsEnterpriseOnly() && !constants.IsEnterprise {
		return logical.ErrorResponse(fmt.Sprintf(ErrKeyTypeEntOnly, srcP.Type)), logical.ErrInvalidRequest
	}

	retKeys := map[string]string{}
	var exportVersion *int
	switch version {
	case "":
		for k, v := range srcP.Keys {
			exportKey, err := getBYOKExportKey(dstP, srcP, &v, hash)
			if err != nil {
				return nil, err
			}
			retKeys[k] = exportKey
		}

	default:
		var versionValue int
		if version == "latest" {
			versionValue = srcP.LatestVersion
		} else {
			version = strings.TrimPrefix(version, "v")
			versionValue, err = strconv.Atoi(version)
			if err != nil {
				return logical.ErrorResponse("invalid key version"), logical.ErrInvalidRequest
			}
		}

		if versionValue < srcP.MinDecryptionVersion {
			return logical.ErrorResponse("version for export is below minimum decryption version"), logical.ErrInvalidRequest
		}
		key, ok := srcP.Keys[strconv.Itoa(versionValue)]
		if !ok {
			return logical.ErrorResponse("version does not exist or cannot be found"), logical.ErrInvalidRequest
		}

		exportKey, err := getBYOKExportKey(dstP, srcP, &key, hash)
		if err != nil {
			return nil, err
		}

		retKeys[strconv.Itoa(versionValue)] = exportKey
		exportVersion = &versionValue
	}

	metadata := b.keyPolicyObservationMetadata(srcP)
	if exportVersion != nil {
		metadata["export_version"] = *exportVersion
	}
	metadata["destination_key"] = dstP.Name
	b.TryRecordObservationWithRequest(ctx, req, ObservationTypeTransitKeyExportBYOK, metadata)

	resp := &logical.Response{
		Data: map[string]interface{}{
			"name": srcP.Name,
			"type": srcP.Type.String(),
			"keys": retKeys,
		},
	}

	return resp, nil
}

func getBYOKExportKey(dstP *keysutil.Policy, srcP *keysutil.Policy, key *keysutil.KeyEntry, hash string) (string, error) {
	if dstP == nil || srcP == nil {
		return "", errors.New("nil policy provided")
	}

	var targetKey interface{}
	switch srcP.Type {
	case keysutil.KeyType_AES128_GCM96, keysutil.KeyType_AES256_GCM96, keysutil.KeyType_ChaCha20_Poly1305, keysutil.KeyType_HMAC, keysutil.KeyType_AES128_CMAC, keysutil.KeyType_AES256_CMAC, keysutil.KeyType_AES192_CMAC, keysutil.KeyType_AES128_CBC, keysutil.KeyType_AES256_CBC:
		targetKey = key.Key
	case keysutil.KeyType_RSA2048, keysutil.KeyType_RSA3072, keysutil.KeyType_RSA4096:
		targetKey = key.RSAKey
	case keysutil.KeyType_ECDSA_P256, keysutil.KeyType_ECDSA_P384, keysutil.KeyType_ECDSA_P521:
		var curve elliptic.Curve
		switch srcP.Type {
		case keysutil.KeyType_ECDSA_P384:
			curve = elliptic.P384()
		case keysutil.KeyType_ECDSA_P521:
			curve = elliptic.P521()
		default:
			curve = elliptic.P256()
		}
		pubKey := ecdsa.PublicKey{
			Curve: curve,
			X:     key.EC_X,
			Y:     key.EC_Y,
		}
		targetKey = &ecdsa.PrivateKey{
			PublicKey: pubKey,
			D:         key.EC_D,
		}
	case keysutil.KeyType_ED25519:
		targetKey = ed25519.PrivateKey(key.Key)
	default:
		return "", fmt.Errorf("unable to export to unknown key type: %v", srcP.Type)
	}

	hasher, err := parseHashFn(hash)
	if err != nil {
		return "", err
	}

	return dstP.WrapKey(0, targetKey, srcP.Type, hasher)
}

const pathBYOKExportHelpSyn = `Securely export named encryption or signing key`

const pathBYOKExportHelpDesc = `
This path is used to export the named keys that are configured as
exportable.

Unlike the regular /export/:name[/:version] paths, this path uses
the same encryption specification /import, allowing secure migration
of keys between clusters to enable workloads to communicate between
them.

Presently this only works for RSA destination keys.
`
