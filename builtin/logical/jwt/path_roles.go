package jwt

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the Role",
			},
			"algorithm": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "RS256",
				Description: "Algorithm for JWT Signing",
			},
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Private Key (RSA or EC) or String for HMAC Algorithm",
			},
			"default_issuer": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Default Issuer for the Role for the JWT Tokens",
			},
			"default_subject": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Default Subject for the Role for the JWT Token",
			},
			"default_audience": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Default Audience for the Role for the JWT Token",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.WriteOperation:  b.pathRoleCreate,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRolesHelpSyn,
		HelpDescription: pathRolesHelpDesc,
	}
}

func (b *backend) getRole(s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get("role/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("role/" + data.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, err := b.getRole(req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	var r = structs.New(role).Map()

	delete(r, "key")

	resp := &logical.Response{
		Data: r,
	}

	return resp, nil
}

func (b *backend) pathRoleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	key := data.Get("key").(string)
	alg := data.Get("algorithm").(string)

	signingMethod := jwt.GetSigningMethod(data.Get("algorithm").(string))
	if signingMethod == nil {
		return nil, fmt.Errorf("Invalid Signing Algorithm")
	}

	if key == "" {
		return nil, fmt.Errorf("Key is Required")
	}

	if strings.HasPrefix(alg, "RS") {
		// need RSA Private Key
		if strings.Contains(key, "RSA PRIVATE KEY") == false {
			return nil, fmt.Errorf("Key is not a PEM formatted RSA Private Key")
		}

		block, _ := pem.Decode([]byte(key))
		_, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse the private key: %s", err)
		}
	} else if strings.HasPrefix(alg, "HS") {
		// need a string
		if key == "" {
			return nil, fmt.Errorf("Key must not be blank")
		}
	} else if strings.HasPrefix(alg, "EC") {
		// need EC Private Key
		if strings.Contains(key, "EC PRIVATE KEY") == false {
			return nil, fmt.Errorf("Key is not a PEM formatted EC Private Key")
		}

		block, _ := pem.Decode([]byte(key))
		_, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse the private key: %s", err)
		}
	}

	entry := &roleEntry{
		Algorithm: alg,
		Key:       key,
		Issuer:    data.Get("default_issuer").(string),
		Subject:   data.Get("default_subject").(string),
		Audience:  data.Get("default_audience").(string),
	}

	// Store it
	jsonEntry, err := logical.StorageEntryJSON("role/"+name, entry)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(jsonEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

type roleEntry struct {
	Algorithm string `json:"algorithm" structs:"algorithm" mapstructure:"algorithm"`
	Key       string `json:"key" structs:"key" mapstructure:"key"`
	Issuer    string `json:"iss" structs:"iss" mapstructure:"iss"`
	Subject   string `json:"sub" structs:"sub" mapstructure:"sub"`
	Audience  string `json:"aud" structs:"aud" mapstructure:"aud"`
}

const pathRolesHelpSyn = `
Read and write basic configuration for generating signed JWT Tokens.
`

const pathRolesHelpDesc = `
This path allows you to read and write roles that are used to
create JWT tokens. These roles have a few settings that dictated
what signing algorithm is used for the JWT token. For example,
if the backend is mounted at "jwt" and you create a role at
"jwt/roles/auth" then a user can request a JWT token at "jwt/issue/auth".
`
