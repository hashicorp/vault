package oidc

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"time"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `login`,
		Fields: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "OIDC Identity Token to be used for login.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLogin,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLogin(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	token := d.Get("token").(string)

	config, err := b.Config(req.Storage)
	if err != nil {
		return logical.ErrorResponse("OIDC backend not configured"), nil
	}

	tokenUsername, tokenGroups, tokenStruct, err := b.validateAndExtractClaims(config, token)
	if err != nil {
		b.Logger().Info("auth/oidc: IdToken verification failed", "error", err)
		return logical.ErrorResponse("OIDC backend bad token"), nil
	}
	if b.Logger().IsDebug() {
		b.Logger().Debug("auth/oidc: IdToken verified", "tokenUsername", tokenUsername, "tokenGroups", tokenGroups)
	}

	loginResponse := &logical.Response{
		Data: map[string]interface{}{},
	}
	if len(tokenGroups) == 0 {
		errString := fmt.Sprintf(
			"no identity token groups found; only policies from user-defined groups available")
		loginResponse.AddWarning(errString)
	}

	var allGroups []string
	// Import the custom added tokenGroups from okta backend
	user, err := b.User(req.Storage, tokenUsername)
	if err == nil && user != nil && user.Groups != nil {
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/oidc: adding local user groups", "userGroups", user.Groups)
		}
		allGroups = append(allGroups, user.Groups...)
	}
	// Merge local and Okta tokenGroups
	allGroups = append(allGroups, tokenGroups...)

	// Retrieve policies
	var policies []string
	for _, groupName := range allGroups {
		group, err := b.Group(req.Storage, groupName)
		if err == nil && group != nil && group.Policies != nil {
			policies = append(policies, group.Policies...)
		}
	}
	// Merge local Policies into Okta Policies
	if user != nil && user.Policies != nil {
		policies = append(policies, user.Policies...)
	}

	if len(policies) == 0 {
		errStr := "user is not a member of any authorized policy"
		if len(loginResponse.Warnings()) > 0 {
			errStr = fmt.Sprintf("%s; additionally, %s", errStr, loginResponse.Warnings()[0])
		}

		loginResponse.Data["error"] = errStr
		return loginResponse, nil
	}
	// TODO(mwitkow): Copy pasted from okta, don't think it is correct
	if err != nil {
		return nil, err
	}

	sort.Strings(policies)
	loginResponse.Auth = &logical.Auth{
		Policies: policies,
		Metadata: map[string]string{
			"username":        tokenUsername,
			"token_nonce":     tokenStruct.Nonce,
			"token_issuer":    tokenStruct.Issuer,
			"token_issued_at": fmt.Sprintf("%v", tokenStruct.IssuedAt),
			"policies":        strings.Join(policies, ","),
		},
		DisplayName: tokenUsername,
		LeaseOptions: logical.LeaseOptions{
			TTL:       tokenStruct.Expiry.Sub(time.Now()),
			Renewable: false,
		},
	}
	b.Logger().Debug("Returning login response", "ttl_set", loginResponse.Auth.LeaseOptions.TTL)
	return loginResponse, nil
}

const pathLoginSyn = `
Log in with an OIDC identity token.
`

const pathLoginDesc = `
This endpoint authenticates a user with an OIDC identity token.
`
