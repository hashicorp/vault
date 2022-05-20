package ssh

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathIssue(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issue/" + framework.GenericNameWithAtRegex("role"),

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathIssue,
			},
		},
	}
}

func (b *backend) pathIssue(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get role? or should it be passed? / Role is always required here
	roleName := data.Get("role").(string)
	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", roleName)), nil
	}
	// If "KeyType" is not "ca" return?
	if role.KeyType != "ca" {
		return logical.ErrorResponse("role key type \"any\" not allowed for issuing certificates, only signing"), nil
	}

	// We are expecting a "key_type" and "key_bits"
	keyType := data.Get("keyType").(string)
	if keyType == "" {
		return logical.ErrorResponse("missing key_type"), nil
	}

	// Can "keyBits" be 0?
	keyBits := data.Get("keyBits").(int)
	if keyBits == 0 {
		return logical.ErrorResponse("missing key_bits"), nil
	}

	//allowed_user_key_lengths | Also return list of allowed key_types?
	keyLengths, present := role.AllowedUserKeyTypesLengths[keyType]
	if !present {
		return logical.ErrorResponse("key_type not in allowed_user_key_lengths"), nil
	}

	// Also return list of allowed key bits?
	present = false
	for _, kb := range keyLengths {
		if keyBits == kb {
			present = true
			break
		}
	}
	if !present {
		return logical.ErrorResponse("key_bits not in list of valid key lengths"), nil
	}

	// Create key pair
	publicKey, privateKey, err := generateSSHKeyPair(b.Backend.GetRandomReader(), keyType, keyBits)
	if err != nil {
		return nil, err
	}

	if publicKey == "" || privateKey == "" {
		return nil, fmt.Errorf("failed to generate or parse the keys")
	}

	// Let's see if this is working
	respData := map[string]interface{}{}

	respData["private_key"] = privateKey
	respData["public_key"] = publicKey

	// Create response
	resp := &logical.Response{
		Data: respData,
	}

	return resp, nil
}
