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
		Pattern: "key-export/" + framework.GenericNameRegex("name") + framework.OptionalParamRegex("version") + "(/hmac)?",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
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
	name := d.Get("name").(string)

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

	// Return HMAC keys, if requesting signing keys
	if strings.HasSuffix(req.Path, "/hmac") {
		retKeys := map[string]string{}
		for k, v := range p.Keys {
			retKeys[strconv.Itoa(k)] = base64.StdEncoding.EncodeToString(v.HMACKey)
		}
		resp.Data["keys"] = retKeys
		return resp, nil
	}

	switch p.Type {
	case keysutil.KeyType_AES256_GCM96:
		retKeys := map[string]string{}
		for k, v := range p.Keys {
			retKeys[strconv.Itoa(k)] = base64.StdEncoding.EncodeToString(v.AESKey)
		}
		resp.Data["keys"] = retKeys

	case keysutil.KeyType_ECDSA_P256:
		retKeys := map[string]string{}
		for k, v := range p.Keys {
			privKey := &ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: elliptic.P256(),
					X:     v.EC_X,
					Y:     v.EC_Y,
				},
				D: v.EC_D,
			}
			ecder, err := x509.MarshalECPrivateKey(privKey)
			if err != nil {
				return nil, err
			}

			block := pem.Block{
				Type:  "EC PRIVATE KEY",
				Bytes: ecder,
			}
			retKeys[strconv.Itoa(k)] = strings.TrimSpace(string(pem.EncodeToMemory(&block)))
		}
		resp.Data["keys"] = retKeys

	default:
		return nil, fmt.Errorf("unknown key type %v", p.Type)
	}

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

const pathExportHelpSyn = `Export named encryption key`

const pathExportHelpDesc = `
This path is used to export the named keys that are configured as
exportable.
`
