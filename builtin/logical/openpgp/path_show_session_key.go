package openpgp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/keybase/go-crypto/openpgp"
	"github.com/keybase/go-crypto/openpgp/armor"
	"github.com/keybase/go-crypto/openpgp/packet"
	"io"
	"strings"
)

func pathShowSessionKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "show-session-key/" + framework.GenericNameRegex("name"),
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
			logical.UpdateOperation: b.pathShowSessionKeyWrite,
		},

		HelpSynopsis:    pathDecryptSessionKeyHelpSyn,
		HelpDescription: pathDecryptSessionKeyHelpDesc,
	}
}

func (b *backend) pathShowSessionKeyWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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

	var p packet.Packet
	var sessionKey string
	for {
		p, err = packet.Read(ciphertextDecoder)
		if err == io.EOF {
			return logical.ErrorResponse("Unable to decrypt session key"), nil
		}
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		switch p := p.(type) {
		case *packet.EncryptedKey:
			encryptedKey := packet.EncryptedKey(*p)
			keys := keyring.KeysById(encryptedKey.KeyId, nil)
			for _, key := range keys {
				encryptedKey.Decrypt(key.PrivateKey, nil)

				if encryptedKey.Key != nil && len(encryptedKey.Key) > 0 {
					sessionKey = fmt.Sprintf("%d:%s", encryptedKey.CipherFunc, strings.ToUpper(hex.EncodeToString(encryptedKey.Key)))
					return &logical.Response{
						Data: map[string]interface{}{
							"session_key": sessionKey,
						},
					}, nil
				}
			}
		}
	}
}

const pathDecryptSessionKeyHelpSyn = "Decrypt a session key of a message using a named PGP key"

const pathDecryptSessionKeyHelpDesc = `
This path uses the named PGP key from the request path to decrypt the session key of a message. 
`
