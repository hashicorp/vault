package userpass

import (
	"context"
	"crypto/subtle"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/cidrutil"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/crypto/bcrypt"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login/" + framework.GenericNameRegex("username"),
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username of the user.",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password for this user.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.pathLoginAliasLookahead,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLoginAliasLookahead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("username").(string))
	if username == "" {
		return nil, fmt.Errorf("missing username")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: username,
			},
		},
	}, nil
}

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("username").(string))

	password := d.Get("password").(string)
	if password == "" {
		return nil, fmt.Errorf("missing password")
	}

	// Get the user and validate auth
	user, userError := b.user(ctx, req.Storage, username)

	var userPassword []byte
	var legacyPassword bool
	// If there was an error or it's nil, we fake a password for the bcrypt
	// check so as not to have a timing leak. Specifics of the underlying
	// storage still leaks a bit but generally much more in the noise compared
	// to bcrypt.
	if user != nil && userError == nil {
		if user.PasswordHash == nil {
			userPassword = []byte(user.Password)
			legacyPassword = true
		} else {
			userPassword = user.PasswordHash
		}
	} else {
		// This is still acceptable as bcrypt will still make sure it takes
		// a long time, it's just nicer to be random if possible
		userPassword = []byte("dummy")
	}

	// Check for a password match. Check for a hash collision for Vault 0.2+,
	// but handle the older legacy passwords with a constant time comparison.
	passwordBytes := []byte(password)
	if !legacyPassword {
		if err := bcrypt.CompareHashAndPassword(userPassword, passwordBytes); err != nil {
			return logical.ErrorResponse("invalid username or password"), nil
		}
	} else {
		if subtle.ConstantTimeCompare(userPassword, passwordBytes) != 1 {
			return logical.ErrorResponse("invalid username or password"), nil
		}
	}

	if userError != nil {
		return nil, userError
	}
	if user == nil {
		return logical.ErrorResponse("invalid username or password"), nil
	}

	// Check for a CIDR match.
	if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, user.BoundCIDRs) {
		return logical.ErrorResponse("login request originated from invalid CIDR"), nil
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Policies: user.Policies,
			Metadata: map[string]string{
				"username": username,
			},
			DisplayName: username,
			LeaseOptions: logical.LeaseOptions{
				TTL:       user.TTL,
				MaxTTL:    user.MaxTTL,
				Renewable: true,
			},
			Alias: &logical.Alias{
				Name: username,
			},
			BoundCIDRs: user.BoundCIDRs,
		},
	}, nil
}

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the user
	user, err := b.user(ctx, req.Storage, req.Auth.Metadata["username"])
	if err != nil {
		return nil, err
	}
	if user == nil {
		// User no longer exists, do not renew
		return nil, nil
	}

	if !policyutil.EquivalentPolicies(user.Policies, req.Auth.TokenPolicies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.TTL = user.TTL
	resp.Auth.MaxTTL = user.MaxTTL
	return resp, nil
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password.
`
