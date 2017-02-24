package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login" + framework.OptionalParamRegex("urlusername"),
		Fields: map[string]*framework.FieldSchema{
			"urlusername": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username to be used for login. (URL parameter)",
			},

			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username to be used for login. (POST request body)",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password for this user.",
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
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	if username == "" {
		username = d.Get("urlusername").(string)
		if username == "" {
			return logical.ErrorResponse("username cannot be empty"), nil
		}
	}

	if password == "" {
		return logical.ErrorResponse("password cannot be empty"), nil
	}

	policies, resp, err := b.HttpLogin(req, username, password)
	// Handle an internal error
	if err != nil {
		return nil, err
	}
	if resp != nil {
		// Handle a logical error
		if resp.IsError() {
			return resp, nil
		}
	}

	resp.Auth = &logical.Auth{
		Policies: policies,
		Metadata: map[string]string{
			"username": username,
			"policies": strings.Join(policies, ","),
		},
		InternalData: map[string]interface{}{
			"password": password,
		},
		DisplayName: username,
		LeaseOptions: logical.LeaseOptions{
			Renewable: true,
		},
	}
	return resp, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var err error

	username := req.Auth.Metadata["username"]
	password := req.Auth.InternalData["password"].(string)

	var resp *logical.Response
	var loginPolicies []string

	loginPolicies, resp, err = b.HttpLogin(req, username, password)
	if err != nil || (resp != nil && resp.IsError()) {
		return resp, err
	}

	if !policyutil.EquivalentPolicies(loginPolicies, req.Auth.Policies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	return framework.LeaseExtend(0, 0, b.System())(req, d)
}

func (b *backend) HttpLogin(req *logical.Request, username string, password string) ([]string, *logical.Response, error) {

	cfg, err := b.Config(req)
	if err != nil {
		return nil, nil, err
	}
	if cfg == nil || cfg.Url == "" {
		return nil, logical.ErrorResponse("HTTP backend not configured"), nil
	}

    client := &http.Client{ Timeout: 10 * time.Second }

    /* Authenticate */
    authreq, err := http.NewRequest("GET", cfg.Url, nil)
    authreq.SetBasicAuth(username, password)
    resp, err := client.Do(authreq)
    if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
    }
    if resp.StatusCode != 200 {
        return nil, logical.ErrorResponse("access denied by the authentication server"), nil
    }

	var policies []string
	// Retrieve user entry from storage
	user, err := b.user(req.Storage, username)
	if user == nil {
		// No user found, check if unregistered users are allowed (unregistered_user_policies not empty)
		if len(policyutil.SanitizePolicies(cfg.UnregisteredUserPolicies, false)) == 0 {
			return nil, logical.ErrorResponse("authentication succeeded but user has no associated policies"), nil
		}
		policies = policyutil.SanitizePolicies(cfg.UnregisteredUserPolicies, true)
	} else {
		policies = policyutil.SanitizePolicies(user.Policies, true)
	}

	return policies, &logical.Response{}, nil
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password. Please be sure to
read the note on escaping from the path-help for the 'config' endpoint.
`
