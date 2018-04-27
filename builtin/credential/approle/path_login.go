package approle

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/cidrutil"
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

func (b *backend) pathLoginUpdateAliasLookahead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
func (b *backend) pathLoginUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	metadata := make(map[string]string)

	// RoleID must be supplied during every login
	roleID := strings.TrimSpace(data.Get("role_id").(string))
	if roleID == "" {
		return logical.ErrorResponse("missing role_id"), nil
	}

	// Look for the storage entry that maps the roleID to role
	roleIDIndex, err := b.roleIDEntry(ctx, req.Storage, roleID)
	if err != nil {
		return nil, err
	}
	if roleIDIndex == nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid role_id %q", roleID)), nil
	}

	lock := b.roleLock(roleIDIndex.Name)
	lock.RLock()
	defer lock.RUnlock()

	role, err := b.roleEntry(ctx, req.Storage, roleIDIndex.Name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid role_id %q", roleID)), nil
	}

	var secretID string
	if role.BindSecretID {
		// If 'bind_secret_id' was set on role, look for the field 'secret_id'
		// to be specified and validate it.
		secretID = strings.TrimSpace(data.Get("secret_id").(string))
		if secretID == "" {
			return logical.ErrorResponse("missing secret_id"), nil
		}

		if role.LowerCaseRoleName {
			roleIDIndex.Name = strings.ToLower(roleIDIndex.Name)
		}

		// Check if the SecretID supplied is valid. If use limit was specified
		// on the SecretID, it will be decremented in this call.
		var valid bool
		valid, metadata, err = b.validateBindSecretID(ctx, req, roleIDIndex.Name, secretID, role)
		if err != nil {
			return nil, err
		}
		if !valid {
			return logical.ErrorResponse(fmt.Sprintf("invalid secret_id %q", secretID)), nil
		}
	}

	if len(role.BoundCIDRList) != 0 {
		// If 'bound_cidr_list' was set, verify the CIDR restrictions
		if req.Connection == nil || req.Connection.RemoteAddr == "" {
			return nil, fmt.Errorf("failed to get connection information")
		}

		belongs, err := cidrutil.IPBelongsToCIDRBlocksSlice(req.Connection.RemoteAddr, role.BoundCIDRList)
		if err != nil {
			return nil, errwrap.Wrapf("failed to verify the CIDR restrictions set on the role: {{err}}", err)
		}
		if !belongs {
			return logical.ErrorResponse(fmt.Sprintf("source address %q unauthorized by CIDR restrictions on the role", req.Connection.RemoteAddr)), nil
		}
	}

	// Always include the role name, for later filtering
	metadata["role_name"] = roleIDIndex.Name

	auth := &logical.Auth{
		NumUses: role.TokenNumUses,
		Period:  role.Period,
		InternalData: map[string]interface{}{
			"role_name": roleIDIndex.Name,
		},
		Metadata: metadata,
		Policies: role.Policies,
		LeaseOptions: logical.LeaseOptions{
			Renewable: true,
			TTL:       role.TokenTTL,
			MaxTTL:    role.TokenMaxTTL,
		},
		Alias: &logical.Alias{
			Name: role.RoleID,
		},
	}

	return &logical.Response{
		Auth: auth,
	}, nil
}

// Invoked when the token issued by this backend is attempting a renewal.
func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := req.Auth.InternalData["role_name"].(string)
	if roleName == "" {
		return nil, fmt.Errorf("failed to fetch role_name during renewal")
	}

	lock := b.roleLock(roleName)
	lock.RLock()
	defer lock.RUnlock()

	// Ensure that the Role still exists.
	role, err := b.roleEntry(ctx, req.Storage, strings.ToLower(roleName))
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to validate role %q during renewal: {{err}}", roleName), err)
	}
	if role == nil {
		return nil, fmt.Errorf("role %q does not exist during renewal", roleName)
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.TTL = role.TokenTTL
	resp.Auth.MaxTTL = role.TokenMaxTTL
	resp.Auth.Period = role.Period
	return resp, nil
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
