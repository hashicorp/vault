// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathListRoles() *framework.Path {
	return &framework.Path{
		Pattern: "role/?$",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAliCloud,
			OperationVerb:   "list",
			OperationSuffix: "roles",
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.operationRolesList,
		},
		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

func (b *backend) pathRole() *framework.Path {
	return &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("name"),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAliCloud,
			OperationSuffix: "role",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeLowerCaseString,
				Description: "The name of the role.",
			},
			"role_arn": {
				Type: framework.TypeString,
				Description: `ARN of the role to be assumed. If provided, inline_policies and 
remote_policies should be blank. At creation time, this role must have configured trusted actors,
and the access key and secret that will be used to assume the role (in /config) must qualify
as a trusted actor.`,
			},
			"inline_policies": {
				Type:        framework.TypeString,
				Description: "JSON of policies to be dynamically applied to users of this role.",
			},
			"remote_policies": {
				Type: framework.TypeStringSlice,
				Description: `The name and type of each remote policy to be applied. 
Example: "name:AliyunRDSReadOnlyAccess,type:System".`,
			},
			"ttl": {
				Type: framework.TypeDurationSecond,
				Description: `Duration in seconds after which the issued token should expire. Defaults
to 0, in which case the value will fallback to the system/mount defaults.`,
			},
			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "The maximum allowed lifetime of tokens issued using this role.",
			},
		},
		ExistenceCheck: b.operationRoleExistenceCheck,
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.operationRoleCreateUpdate,
			logical.UpdateOperation: b.operationRoleCreateUpdate,
			logical.ReadOperation:   b.operationRoleRead,
			logical.DeleteOperation: b.operationRoleDelete,
		},
		HelpSynopsis:    pathRolesHelpSyn,
		HelpDescription: pathRolesHelpDesc,
	}
}

func (b *backend) operationRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := readRole(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) operationRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return nil, errors.New("name is required")
	}

	role, err := readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil && req.Operation == logical.UpdateOperation {
		return nil, fmt.Errorf("no role found to update for %s", roleName)
	} else if role == nil {
		role = &roleEntry{}
	}

	if raw, ok := data.GetOk("role_arn"); ok {
		role.RoleARN = raw.(string)
	}
	if raw, ok := data.GetOk("inline_policies"); ok {
		policyDocsStr := raw.(string)

		var policyDocs []map[string]interface{}
		if err := json.Unmarshal([]byte(policyDocsStr), &policyDocs); err != nil {
			return nil, err
		}

		// If any inline policies were set before, we need to clear them and consider
		// these the new ones.
		role.InlinePolicies = make([]*inlinePolicy, len(policyDocs))

		for i, policyDoc := range policyDocs {
			uid, err := uuid.GenerateUUID()
			if err != nil {
				return nil, err
			}
			uid = strings.Replace(uid, "-", "", -1)
			role.InlinePolicies[i] = &inlinePolicy{
				UUID:           uid,
				PolicyDocument: policyDoc,
			}
		}
	}
	if raw, ok := data.GetOk("remote_policies"); ok {
		strPolicies := raw.([]string)

		// If any remote policies were set before, we need to clear them and consider
		// these the new ones.
		role.RemotePolicies = make([]*remotePolicy, len(strPolicies))

		for i, strPolicy := range strPolicies {
			policy := &remotePolicy{}
			kvPairs := strings.Split(strPolicy, ",")
			for _, kvPair := range kvPairs {
				kvFields := strings.Split(kvPair, ":")
				if len(kvFields) != 2 {
					return nil, fmt.Errorf("unable to recognize pair in %s", kvPair)
				}
				switch kvFields[0] {
				case "name":
					policy.Name = kvFields[1]
				case "type":
					policy.Type = kvFields[1]
				default:
					return nil, fmt.Errorf("invalid key: %s", kvFields[0])
				}
			}
			if policy.Name == "" {
				return nil, fmt.Errorf("policy name is required in %s", strPolicy)
			}
			if policy.Type == "" {
				return nil, fmt.Errorf("policy type is required in %s", strPolicy)
			}
			role.RemotePolicies[i] = policy
		}
	}
	if raw, ok := data.GetOk("ttl"); ok {
		role.TTL = time.Duration(raw.(int)) * time.Second
	}
	if raw, ok := data.GetOk("max_ttl"); ok {
		role.MaxTTL = time.Duration(raw.(int)) * time.Second
	}

	// Now that the role is built, validate it.
	if role.MaxTTL > 0 && role.TTL > role.MaxTTL {
		return nil, errors.New("ttl exceeds max_ttl")
	}
	if role.Type() == roleTypeSTS {
		if len(role.RemotePolicies) > 0 {
			return nil, errors.New("remote_policies must be blank when an arn is present")
		}
		if len(role.InlinePolicies) > 0 {
			return nil, errors.New("inline_policies must be blank when an arn is present")
		}
	} else if len(role.InlinePolicies)+len(role.RemotePolicies) == 0 {
		return nil, errors.New("must include an arn, or at least one of inline_policies or remote_policies")
	}

	entry, err := logical.StorageEntryJSON("role/"+roleName, role)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Let's create a response that we're only going to return if there are warnings.
	resp := &logical.Response{}
	if role.Type() == roleTypeSTS && (role.TTL > 0 || role.MaxTTL > 0) {
		resp.AddWarning("role_arn is set so ttl and max_ttl will be ignored because they're not editable on STS tokens")
	}
	if role.TTL > b.System().MaxLeaseTTL() {
		resp.AddWarning(fmt.Sprintf("ttl of %d exceeds the system max ttl of %d, the latter will be used during login", role.TTL, b.System().MaxLeaseTTL()))
	}
	if len(resp.Warnings) > 0 {
		return resp, nil
	}
	// No warnings, let's return a 204.
	return nil, nil
}

func (b *backend) operationRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return nil, errors.New("name is required")
	}

	role, err := readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"role_arn":        role.RoleARN,
			"remote_policies": role.RemotePolicies,
			"inline_policies": role.InlinePolicies,
			"ttl":             role.TTL / time.Second,
			"max_ttl":         role.MaxTTL / time.Second,
		},
	}, nil
}

func (b *backend) operationRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, "role/"+data.Get("name").(string)); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) operationRolesList(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(entries), nil
}

func readRole(ctx context.Context, s logical.Storage, roleName string) (*roleEntry, error) {
	role, err := s.Get(ctx, "role/"+roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	result := &roleEntry{}
	if err := role.DecodeJSON(result); err != nil {
		return nil, err
	}
	return result, nil
}

type roleType int

const (
	roleTypeUnknown roleType = iota
	roleTypeRAM
	roleTypeSTS
)

func parseRoleType(nameOfRoleType string) (roleType, error) {
	switch nameOfRoleType {
	case "ram":
		return roleTypeRAM, nil
	case "sts":
		return roleTypeSTS, nil
	default:
		return roleTypeUnknown, fmt.Errorf("received unknown role type: %s", nameOfRoleType)
	}
}

func (t roleType) String() string {
	switch t {
	case roleTypeRAM:
		return "ram"
	case roleTypeSTS:
		return "sts"
	}
	return "unknown"
}

type roleEntry struct {
	RoleARN        string          `json:"role_arn"`
	RemotePolicies []*remotePolicy `json:"remote_policies"`
	InlinePolicies []*inlinePolicy `json:"inline_policies"`
	TTL            time.Duration   `json:"ttl"`
	MaxTTL         time.Duration   `json:"max_ttl"`
}

func (r *roleEntry) Type() roleType {
	if r.RoleARN != "" {
		return roleTypeSTS
	}
	return roleTypeRAM
}

// Policies don't have ARNs and instead, their unique combination of their name and type comprise
// their unique identifier.
type remotePolicy struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type inlinePolicy struct {
	// UUID is used in naming the policy. The policy document has no fields
	// that would reliably be there and make a beautiful, human-readable name.
	// So instead, we generate a UUID for it and use that in the policy name,
	// which is likewise returned when roles are read so generated policy names
	// can be tied back to which policy document they're for.
	UUID           string                 `json:"hash"`
	PolicyDocument map[string]interface{} `json:"policy_document"`
}

const pathListRolesHelpSyn = "List the existing roles in this backend."

const pathListRolesHelpDesc = "Roles will be listed by the role name."

const pathRolesHelpSyn = `
Read, write and reference policies and roles that API keys or STS credentials can be made for.
`

const pathRolesHelpDesc = `
This path allows you to read and write roles that are used to
create API keys or STS credentials.

If you supply a role ARN, that role must have been created to allow trusted actors,
and the access key and secret that will be used to call AssumeRole (configured at
the /config path) must qualify as a trusted actor.

If you instead supply inline and/or remote policies to be applied, a user and API
key will be dynamically created. The remote policies will be applied to that user,
and the inline policies will also be dynamically created and applied.

To obtain an API key or STS credential after the role is created, if the
backend is mounted at "alicloud" and you create a role at "alicloud/roles/deploy",
then a user could request access credentials at "alicloud/creds/deploy".

To validate the keys, attempt to read an access key after writing the policy.
`
