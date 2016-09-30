package awsiam

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRole(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			"bound_iam_principal": {
				Type:        framework.TypeString,
				Description: "ARN of the IAM principal to bind to this role.",
			},
			"ttl": {
				Type:    framework.TypeDurationSecond,
				Default: 0,
				Description: `Duration in seconds after which the issued token should expire. Defaults
to 0, in whcih case the value will fall back to the system/mount defaults.`,
			},
			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Default:     0,
				Description: "The maximum allowed lifetime of tokens issued using this role",
			},
			"policies": {
				Type:        framework.TypeString,
				Default:     "default",
				Description: "Policies to be set on tokens issued using this role.",
			},
		},

		ExistenceCheck: b.pathRoleExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathRoleCreateUpdate,
			logical.UpdateOperation: b.pathRoleCreateUpdate,
			logical.ReadOperation:   b.pathRoleRead,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleSyn,
		HelpDescription: pathRoleDesc,
	}
}

func pathListRole(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles?/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

type awsRoleEntry struct {
	TTL               time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	MaxTTL            time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	BoundIamPrincipal string        `json:"bound_iam_principal" structs:"bound_iam_principal" mapstructure:"bound_iam_principal"`
	Policies          []string      `json:"policies" structs:"policies" mapstructure:"policies"`
}

func (b *backend) lockedAWSRole(s logical.Storage, role string) (*awsRoleEntry, error) {
	b.roleMutex.RLock()
	defer b.roleMutex.RUnlock()

	return b.nonLockedAWSRole(s, role)
}

func (b *backend) nonLockedAWSRole(s logical.Storage, role string) (*awsRoleEntry, error) {
	entry, err := s.Get("role/" + strings.ToLower(role))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result awsRoleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *backend) pathRoleExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.lockedAWSRole(req.Storage, strings.ToLower(data.Get("role").(string)))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) pathRoleCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	roleName := strings.ToLower(data.Get("role").(string))
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	roleEntry, err := b.nonLockedAWSRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		roleEntry = &awsRoleEntry{}
	}

	if boundIamPrincipal, ok := data.GetOk("bound_iam_principal"); ok {
		roleEntry.BoundIamPrincipal = boundIamPrincipal.(string)
	}

	// We do this here to permit updating an existing role
	if roleEntry.BoundIamPrincipal == "" {
		return logical.ErrorResponse("Must bind to an IAM Principal"), nil
	}

	policiesStr, ok := data.GetOk("policies")
	if ok {
		roleEntry.Policies = policyutil.ParsePolicies(policiesStr.(string))
	} else if req.Operation == logical.CreateOperation {
		roleEntry.Policies = []string{"default"}
	}

	var resp logical.Response

	ttlRaw, ok := data.GetOk("ttl")
	if ok {
		ttl := time.Duration(ttlRaw.(int)) * time.Second
		defaultLeaseTTL := b.System().DefaultLeaseTTL()
		if ttl > defaultLeaseTTL {
			resp.AddWarning(fmt.Sprintf("Given ttl of %d seconds greater than current mount/system default of %d seconds; ttl will be capped at login time", ttl/time.Second, defaultLeaseTTL/time.Second))
		}
		roleEntry.TTL = ttl
	} else if req.Operation == logical.CreateOperation {
		roleEntry.TTL = time.Duration(data.Get("ttl").(int)) * time.Second
	}

	maxTTLInt, ok := data.GetOk("max_ttl")
	if ok {
		maxTTL := time.Duration(maxTTLInt.(int)) * time.Second
		systemMaxTTL := b.System().MaxLeaseTTL()
		if maxTTL > systemMaxTTL {
			resp.AddWarning(fmt.Sprintf("Given max_ttl of %d seconds greater than current mount/system default of %d seconds; max_ttl will be capped at login time", maxTTL/time.Second, systemMaxTTL/time.Second))
		}

		if maxTTL < time.Duration(0) {
			return logical.ErrorResponse("max_ttl cannot be negative"), nil
		}

		roleEntry.MaxTTL = maxTTL
	} else if req.Operation == logical.CreateOperation {
		roleEntry.MaxTTL = time.Duration(data.Get("max_ttl").(int)) * time.Second
	}

	entry, err := logical.StorageEntryJSON("role/"+roleName, roleEntry)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	if len(resp.Warnings()) == 0 {
		return nil, nil
	}

	return &resp, nil
}

func (b *backend) pathRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	roleEntry, err := b.lockedAWSRole(req.Storage, strings.ToLower(data.Get("role").(string)))
	if err != nil {
		return nil, err
	}

	if roleEntry == nil {
		return nil, nil
	}

	respData := structs.New(roleEntry).Map()
	respData["ttl"] = roleEntry.TTL / time.Second
	respData["max_ttl"] = roleEntry.MaxTTL / time.Second

	return &logical.Response{
		Data: respData,
	}, nil
}

func (b *backend) pathRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	roleName := data.Get("role").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	return nil, req.Storage.Delete("role/" + strings.ToLower(roleName))
}

func (b *backend) pathRoleList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	b.roleMutex.RLock()
	defer b.roleMutex.RUnlock()

	roles, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

const pathRoleSyn = `
Create a role and associate policies to it.
`

const pathRoleDesc = `
The login endpoint takes in the role name which the client
would like to access. After the client is authenticated, the
authorization for the client is determined by the policies that
are associated for the role.

Also, a 'max_ttl' can be configured in this endpoint that determines
the maximum duration for which a login can be renewed. Note that the
'max_ttl' has an upper limit of the 'max_ttl' value on the backend's
mount.
`

const pathListRolesHelpSyn = `
List all roles that are registered with Vault.
`

const pathListRolesHelpDesc = `
Roles will be listed by their respective role names.
`
