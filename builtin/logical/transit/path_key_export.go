package transit

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathExportKeys() *framework.Path {
	return &framework.Path{
		Pattern: "export/" + framework.GenericNameRegex("export_type") + "/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("version"),
		Fields: map[string]*framework.FieldSchema{
			"export_type": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Type of key to export (encryption, signing)",
			},
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
			"version": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathPolicyExport,
		},

		HelpSynopsis:    pathExportHelpSyn,
		HelpDescription: pathExportHelpDesc,
	}
}

func (b *backend) pathPolicyExport(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	exportType := d.Get("export_type").(string)
	name := d.Get("name").(string)
	version := d.Get("version").(string)

	if exportType != "encryptor" && exportType != "signer" {
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

	resp := &logical.Response{
		Data: map[string]interface{}{
			"name": p.Name,
		},
	}

	if version != "" {
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
		resp.Data["version"] = versionValue

		key, ok := p.Keys[versionValue]
		if !ok {
			return logical.ErrorResponse("version does not exist or is no longer valid"), logical.ErrInvalidRequest
		}

		switch exportType {
		case "signer":
			resp.Data["key"] = strings.TrimSpace(base64.StdEncoding.EncodeToString(key.HMACKey))
		case "encryptor":
			switch p.Type {
			case keysutil.KeyType_AES256_GCM96:
				resp.Data["key"] = strings.TrimSpace(base64.StdEncoding.EncodeToString(key.AESKey))
			case keysutil.KeyType_ECDSA_P256:
				ecKey, err := keyEntryToECPrivateKey(key)
				if err != nil {
					return nil, err
				}
				resp.Data["key"] = ecKey
			default:
				return nil, fmt.Errorf("unknown key type %v", p.Type)
			}
		}

		return resp, nil
	}

	retKeys := map[string]string{}
	switch exportType {
	case "signer":
		for k, v := range p.Keys {
			retKeys[strconv.Itoa(k)] = base64.StdEncoding.EncodeToString(v.HMACKey)
		}
	case "encryptor":
		switch p.Type {
		case keysutil.KeyType_AES256_GCM96:
			for k, v := range p.Keys {
				retKeys[strconv.Itoa(k)] = base64.StdEncoding.EncodeToString(v.AESKey)
			}
		case keysutil.KeyType_ECDSA_P256:
			for k, v := range p.Keys {
				ecKey, err := keyEntryToECPrivateKey(v)
				if err != nil {
					return nil, err
				}
				retKeys[strconv.Itoa(k)] = ecKey
			}
		default:
			return nil, fmt.Errorf("unknown key type %v", p.Type)
		}
	}
	resp.Data["keys"] = retKeys

	return resp, nil
}

func keyEntryToECPrivateKey(k keysutil.KeyEntry) (string, error) {
	privKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     k.EC_X,
			Y:     k.EC_Y,
		},
		D: k.EC_D,
	}
	ecder, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return "", err
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
