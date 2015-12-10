package awsKms

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"key_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Key ID",
			},

			"ciphertext": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Encrypted token",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathLogin,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

type token struct {
	NotBefore time.Time `json:"notBefore"`
	NotAfter  time.Time `json:"notAfter"`
}

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cipherText := data.Get("ciphertext").(string)

	// Ensure cipherText is provided
	if cipherText == "" {
		return logical.ErrorResponse("missing 'ciphertext'"), nil
	}

	decodedCiphertext, err := hex.DecodeString(cipherText)
	if err != nil {
		return logical.ErrorResponse("Could not decode the token"), nil
	}

	// Decrypt ciphertext
	svc, err := clientKMS(req.Storage)
	if err != nil {
		return logical.ErrorResponse("Could not connect to AWS"), err
	}
	params := &kms.DecryptInput{
		/*EncryptionContext: map[string]*string{
			"Key": aws.String(keyId),
		},*/
		CiphertextBlob: decodedCiphertext,
	}
	resp, err := svc.Decrypt(params)

	if err != nil {
		return logical.ErrorResponse("Could not decrypt the token"), err
	}

	keyId := (*resp.KeyId)[len(*resp.KeyId)-36:]

	var tok token

	err = json.Unmarshal(resp.Plaintext, &tok)

	if err != nil {
		return logical.ErrorResponse("Cannot parse token"), nil
	}
	
	var now = time.Now().UTC()
	if now.Before(tok.NotBefore) || now.After(tok.NotAfter) {
		return logical.ErrorResponse("Token expired/not valid yet"), nil
	}

	keyRaw, err := b.MapKey.Get(req.Storage, keyId)

	if keyRaw == nil {
		return logical.ErrorResponse("invalid key ID: " + keyId), nil
	}

	if err != nil {
		return nil, err
	}

	policies, err := b.MapKey.Policies(req.Storage, keyId)

	if err != nil {
		return nil, err
	}

	var displayName string
	if raw, ok := keyRaw["display_name"]; ok {
		displayName = raw.(string)
	}
	metadata := map[string]string{
		"key-id":    keyId,
		"cleartext": string(resp.Plaintext[:]),
	}

	return &logical.Response{
		Auth: &logical.Auth{
			DisplayName: displayName,
			Policies:    policies,
			Metadata:    metadata,
		},
	}, nil
}

const pathLoginSyn = `
Log in with an Key ID and encrypted token
`

const pathLoginDesc = `
This endpoint authenticates using an AWS KMS key ID and encrypted token
`
