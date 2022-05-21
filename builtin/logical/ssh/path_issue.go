package ssh

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"log"
)

func pathIssue(b *backend) *framework.Path {
	log.Println("HELLO WORLD")
	return &framework.Path{
		Pattern: "issue/" + framework.GenericNameWithAtRegex("role"),

		//Callbacks: map[logical.Operation]framework.OperationFunc{
		//	logical.UpdateOperation: b.pathIssue,
		//},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathIssue,
			},
		},
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: `The desired role with configuration for this request.`,
			},
			"key_type": {
				Type:        framework.TypeString,
				Description: "TBD",
			},
			"key_bits": {
				Type:        framework.TypeInt,
				Description: "TBD",
			},
		},
		HelpSynopsis:    "TBD - HelpSynopsis",
		HelpDescription: "TBD - HelpDescription",
	}
}

func (b *backend) pathIssue(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	log.Println("HELLO WORLD in pathIssue")
	// Get role or should it be passed?
	roleName := data.Get("role").(string)
	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", roleName)), nil
	}

	// If "KeyType" is not "ca" return?
	if role.KeyType != "ca" {
		return logical.ErrorResponse("role key type \"any\" not allowed for issuing certificates, only signing"), nil
	}

	// We are expecting a "key_type" and "key_bits" | REQUIRED?
	keyType := data.Get("key_type").(string)
	//if keyType == "" {
	//	return logical.ErrorResponse("missing key_type"), nil
	//}

	keyBits := data.Get("key_bits").(int)
	// If statement ?
	//if keyBits == 0 {
	//	return logical.ErrorResponse("missing key_bits"), nil
	//}

	//allowed_user_key_lengths | What if the users provides an invalid key type?
	keyTypeLengths, keyPresent := role.AllowedUserKeyTypesLengths[keyType]
	if keyPresent {
		var bitsPresent bool
		for _, kb := range keyTypeLengths {
			if keyBits == kb {
				bitsPresent = true
				break
			}
		}
		if !bitsPresent {
			// Also return list of allowed key bits?
			return logical.ErrorResponse("key_bits not in list of allowed value  for key_type provided"), nil
		}
	}

	// Create key pair | key_type and key_bits checked here?
	publicKey, privateKey, err := generateSSHKeyPair(b.Backend.GetRandomReader(), keyType, keyBits)
	if err != nil {
		return nil, err
	}

	if publicKey == "" || privateKey == "" {
		return nil, fmt.Errorf("failed to generate or parse the keys")
	}

	// Sign key
	// Raw or Schema?
	data.Raw["public_key"] = publicKey
	data.Raw["private_key"] = privateKey
	log.Printf("Public Key:\n%s\n", publicKey)
	log.Printf("Private Key:\n%s\n", privateKey)
	//return b.pathSignCertificate(ctx, req, data, role)

	return nil, nil

	//	// Everything after this is creating a response
	/*
		respData := map[string]interface{}{}

		respData["private_key"] = privateKey
		respData["public_key"] = publicKey

		// Create response
		resp := &logical.Response{
			Data: respData,
		}
	*/
}
