package jwt

import (
	"encoding/json"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/hashicorp/vault/helper/uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathIssue(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issue/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"role": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The desired role with configuration for this request",
			},
			"issuer": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The issuer of the token",
			},
			"subject": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The subject of the token",
			},
			"audience": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The audience of the token",
			},
			"expiration": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "This will define the expiration in NumericDate value",
			},
			"not_before": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Defines the time before which the JWT MUST NOT be accepted for processing",
			},
			"issued_at": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "The time the JWT was issued",
			},
			"jti": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Unique identifier for the JWT",
			},
			"claims": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "JSON Object of Claims for the JWT Token",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathIssueWrite,
		},
	}
}

func (b *backend) pathIssueWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	// Get the role
	role, err := b.getRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Unknown role: %s", roleName)), nil
	}

	claims := map[string]interface{}{
		"initial": "ok",
	}

	if role.Issuer != "" {
		claims["iss"] = role.Issuer
	}
	if role.Subject != "" {
		claims["sub"] = role.Subject
	}
	if role.Audience != "" {
		claims["aud"] = role.Audience
	}

	if data.Get("not_before") == 0 {
		claims["nbf"] = int(time.Now().Unix())
	}
	if data.Get("issued_at") == 0 {
		claims["iat"] = int(time.Now().Unix())
	}
	if data.Get("jti") == "" {
		claims["jti"] = uuid.GenerateUUID()
	}

	if data.Get("issuer") != "" {
		claims["iss"] = data.Get("issuer").(string)
	}
	if data.Get("subject") != "" {
		claims["sub"] = data.Get("subject").(string)
	}
	if data.Get("audience") != "" {
		claims["aud"] = data.Get("audience").(string)
	}
	if data.Get("expiration").(int) > 0 {
		claims["exp"] = data.Get("expiration").(int)
	}
	if data.Get("not_before").(int) > 0 {
		claims["nbf"] = data.Get("not_before").(int)
	}
	if data.Get("issued_at").(int) > 0 {
		claims["iat"] = data.Get("issued_at").(int)
	}
	if data.Get("jti") != "" {
		claims["jti"] = data.Get("jti").(string)
	}

	if data.Get("claims").(string) != "" {
		// Parse JSON using unmarshal
		var uc map[string]interface{}
		err := json.Unmarshal([]byte(data.Get("claims").(string)), &uc)
		if err != nil {
			return nil, err
		}

		for k, v := range uc {
			claims[k] = v
		}
	}

	delete(claims, "initial")

	token := jwt.New(jwt.GetSigningMethod(role.Algorithm))
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(role.Key))
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"jti":   claims["jti"].(string),
			"token": tokenString,
		},
	}

	return resp, nil
}
