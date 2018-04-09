package openpgp

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/keybase/go-crypto/openpgp"
	"github.com/keybase/go-crypto/openpgp/armor"
	"io"
	"strings"
)

func pathDecrypt(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "decrypt/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "The key to use",
			},
			"ciphertext": {
				Type:        framework.TypeString,
				Description: "The ciphertext to decrypt",
			},
			"format": {
				Type:        framework.TypeString,
				Default:     "base64",
				Description: `Encoding format the ciphertext uses. Can be "base64" or "ascii-armor". Defaults to "base64".`,
			},
			"signer_key": {
				Type:        framework.TypeString,
				Description: "The ASCII-armored PGP key of the signer of the ciphertext. If present, the signature must be valid.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathDecryptWrite,
		},

		HelpSynopsis:    pathDecryptHelpSyn,
		HelpDescription: pathDecryptHelpDesc,
	}
}

func (b *backend) pathDecryptWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	format := data.Get("format").(string)
	switch format {
	case "base64":
	case "ascii-armor":
	default:
		return logical.ErrorResponse(fmt.Sprintf("unsupported encoding format %s; must be \"base64\" or \"ascii-armor\"", format)), nil
	}

	keyEntry, err := b.key(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if keyEntry == nil {
		return logical.ErrorResponse("key not found"), logical.ErrInvalidRequest
	}

	r := bytes.NewReader(keyEntry.SerializedKey)
	keyring, err := openpgp.ReadKeyRing(r)
	if err != nil {
		return nil, err
	}

	signerKey := data.Get("signer_key").(string)
	if signerKey != "" {
		el, err := openpgp.ReadArmoredKeyRing(strings.NewReader(signerKey))
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		keyring = append(keyring, el[0])
	}

	ciphertextEncoded := strings.NewReader(data.Get("ciphertext").(string))
	var ciphertextDecoder io.Reader
	switch format {
	case "base64":
		ciphertextDecoder = base64.NewDecoder(base64.StdEncoding, ciphertextEncoded)
	case "ascii-armor":
		block, err := armor.Decode(ciphertextEncoded)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		ciphertextDecoder = block.Body
	}

	md, err := openpgp.ReadMessage(ciphertextDecoder, keyring, nil, nil)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	var plaintext bytes.Buffer
	w := base64.NewEncoder(base64.StdEncoding, &plaintext)
	if _, err = io.Copy(w, md.UnverifiedBody); err != nil {
		return nil, err
	}
	if err = w.Close(); err != nil {
		return nil, err
	}

	if signerKey != "" && (!md.IsSigned || md.SignedBy == nil || md.SignatureError != nil) {
		return logical.ErrorResponse("Signature is invalid or not present"), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"plaintext": plaintext.String(),
		},
	}, nil
}

const pathDecryptHelpSyn = "Decrypt a ciphertext value using a named PGP key"

const pathDecryptHelpDesc = `
This path uses the named PGP key from the request path to decrypt a user
provided ciphertext. The plaintext is returned base64 encoded.
`
