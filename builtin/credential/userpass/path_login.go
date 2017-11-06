package userpass

import (
	"crypto/subtle"
	"fmt"
	"strings"

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

func (b *backend) pathLoginAliasLookahead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

func (b *backend) pathLogin(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("username").(string))

	password := d.Get("password").(string)
	if password == "" {
		return nil, fmt.Errorf("missing password")
	}

	// Get the user and validate auth
	user, err := b.user(req.Storage, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return logical.ErrorResponse("invalid username or password"), nil
	}

	// Check for a password match. Check for a hash collision for Vault 0.2+,
	// but handle the older legacy passwords with a constant time comparison.
	passwordBytes := []byte(password)
	if user.PasswordHash != nil {
		if err := bcrypt.CompareHashAndPassword(user.PasswordHash, passwordBytes); err != nil {
			return logical.ErrorResponse("invalid username or password"), nil
		}
	} else {
		if subtle.ConstantTimeCompare([]byte(user.Password), passwordBytes) != 1 {
			return logical.ErrorResponse("invalid username or password"), nil
		}
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
				Renewable: true,
			},
			Alias: &logical.Alias{
				Name: username,
			},
		},
	}, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the user
	user, err := b.user(req.Storage, req.Auth.Metadata["username"])
	if err != nil {
		return nil, err
	}
	if user == nil {
		// User no longer exists, do not renew
		return nil, nil
	}

	if !policyutil.EquivalentPolicies(user.Policies, req.Auth.Policies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	return framework.LeaseExtend(user.TTL, user.MaxTTL, b.System())(req, d)
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password.
`
