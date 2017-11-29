package approle

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		Fields: map[string]*framework.FieldSchema{
			"role_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Unique identifier of the Role. Required to be supplied when the 'bind_secret_id' constraint is set.",
			},
			"secret_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "SecretID belong to the App role",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLoginUpdate,
			logical.AliasLookaheadOperation: b.pathLoginUpdateAliasLookahead,
		},
		HelpSynopsis:    pathLoginHelpSys,
		HelpDescription: pathLoginHelpDesc,
	}
}

func (b *backend) pathLoginUpdateAliasLookahead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleID := strings.TrimSpace(data.Get("role_id").(string))
	if roleID == "" {
		return nil, fmt.Errorf("missing role_id")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: roleID,
			},
		},
	}, nil
}

// Returns the Auth object indicating the authentication and authorization information
// if the credentials provided are validated by the backend.
func (b *backend) pathLoginUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, roleName, metadata, _, err := b.validateCredentials(req, data)
	if err != nil || role == nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to validate credentials: %v", err)), nil
	}

	// Always include the role name, for later filtering
	metadata["role_name"] = roleName

	auth := &logical.Auth{
		NumUses: role.TokenNumUses,
		Period:  role.Period,
		InternalData: map[string]interface{}{
			"role_name": roleName,
		},
		Metadata: metadata,
		Policies: role.Policies,
		LeaseOptions: logical.LeaseOptions{
			Renewable: true,
		},
		Alias: &logical.Alias{
			Name: role.RoleID,
		},
	}

	// If 'Period' is set, use the value of 'Period' as the TTL.
	// Otherwise, set the normal TokenTTL.
	if role.Period > time.Duration(0) {
		auth.TTL = role.Period
	} else {
		auth.TTL = role.TokenTTL
	}

	return &logical.Response{
		Auth: auth,
	}, nil
}

// Invoked when the token issued by this backend is attempting a renewal.
func (b *backend) pathLoginRenew(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := req.Auth.InternalData["role_name"].(string)
	if roleName == "" {
		return nil, fmt.Errorf("failed to fetch role_name during renewal")
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	// Ensure that the Role still exists.
	role, err := b.roleEntry(req.Storage, strings.ToLower(roleName))
	if err != nil {
		return nil, fmt.Errorf("failed to validate role %s during renewal:%s", roleName, err)
	}
	if role == nil {
		return nil, fmt.Errorf("role %s does not exist during renewal", roleName)
	}

	// If 'Period' is set on the Role, the token should never expire.
	// Replenish the TTL with 'Period's value.
	if role.Period > time.Duration(0) {
		// If 'Period' was updated after the token was issued,
		// token will bear the updated 'Period' value as its TTL.
		req.Auth.TTL = role.Period
		return &logical.Response{Auth: req.Auth}, nil
	} else {
		return framework.LeaseExtend(role.TokenTTL, role.TokenMaxTTL, b.System())(req, data)
	}
}

const pathLoginHelpSys = "Issue a token based on the credentials supplied"

const pathLoginHelpDesc = `
While the credential 'role_id' is required at all times,
other credentials required depends on the properties App role
to which the 'role_id' belongs to. The 'bind_secret_id'
constraint (enabled by default) on the App role requires the
'secret_id' credential to be presented.

'role_id' is fetched using the 'role/<role_name>/role_id'
endpoint and 'secret_id' is fetched using the 'role/<role_name>/secret_id'
endpoint.`
