package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the policy",
				DisplayName: "Policy Name",
			},

			"credential_type": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: fmt.Sprintf("Type of credential to retrieve. Must be one of %s, %s, or %s", assumedRoleCred, iamUserCred, federationTokenCred),
			},

			"role_arns": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "ARNs of AWS roles allowed to be assumed. Only valid when credential_type is " + assumedRoleCred,
				DisplayName: "Role ARNs",
			},

			"policy_arns": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "ARNs of AWS policies to attach to IAM users. Only valid when credential_type is " + iamUserCred,
				DisplayName: "Policy ARNs",
			},

			"policy_document": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `JSON-encoded IAM policy document. Behavior varies by credential_type. When credential_type is
iam_user, then it will attach the contents of the policy_document to the IAM
user generated. When credential_type is assumed_role or federation_token, this
will be passed in as the Policy parameter to the AssumeRole or
GetFederationToken API call, acting as a filter on permissions available.`,
			},

			"default_sts_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: fmt.Sprintf("Default TTL for %s and %s credential types when no TTL is explicitly requested with the credentials", assumedRoleCred, federationTokenCred),
				DisplayName: "Default TTL",
			},

			"max_sts_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: fmt.Sprintf("Max allowed TTL for %s and %s credential types", assumedRoleCred, federationTokenCred),
				DisplayName: "Max TTL",
			},

			"arn": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Deprecated; use role_arns or policy_arns instead. ARN Reference to a managed policy
or IAM role to assume`,
				Deprecated: true,
			},

			"policy": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Deprecated; use policy_document instead. IAM policy document",
				Deprecated:  true,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathRolesDelete,
			logical.ReadOperation:   b.pathRolesRead,
			logical.UpdateOperation: b.pathRolesWrite,
		},

		HelpSynopsis:    pathRolesHelpSyn,
		HelpDescription: pathRolesHelpDesc,
	}
}

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	b.roleMutex.RLock()
	defer b.roleMutex.RUnlock()
	entries, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}
	legacyEntries, err := req.Storage.List(ctx, "policy/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(append(entries, legacyEntries...)), nil
}

func (b *backend) pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	for _, prefix := range []string{"policy/", "role/"} {
		err := req.Storage.Delete(ctx, prefix+d.Get("name").(string))
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (b *backend) pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.roleRead(ctx, req.Storage, d.Get("name").(string), true)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: entry.toResponseData(),
	}, nil
}

func legacyRoleData(d *framework.FieldData) (string, error) {
	policy := d.Get("policy").(string)
	arn := d.Get("arn").(string)

	switch {
	case policy == "" && arn == "":
		return "", nil
	case policy != "" && arn != "":
		return "", errors.New("only one of policy or arn should be provided")
	case policy != "":
		return policy, nil
	default:
		return arn, nil
	}
}

func (b *backend) pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var resp logical.Response

	roleName := d.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()
	roleEntry, err := b.roleRead(ctx, req.Storage, roleName, false)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		roleEntry = &awsRoleEntry{}
	} else if roleEntry.InvalidData != "" {
		resp.AddWarning(fmt.Sprintf("Invalid data of %q cleared out of role", roleEntry.InvalidData))
		roleEntry.InvalidData = ""
	}

	legacyRole, err := legacyRoleData(d)
	if err != nil {
		return nil, err
	}

	if credentialTypeRaw, ok := d.GetOk("credential_type"); ok {
		if legacyRole != "" {
			return logical.ErrorResponse("cannot supply deprecated role or policy parameters with an explicit credential_type"), nil
		}
		credentialType := credentialTypeRaw.(string)
		allowedCredentialTypes := []string{iamUserCred, assumedRoleCred, federationTokenCred}
		if !strutil.StrListContains(allowedCredentialTypes, credentialType) {
			return logical.ErrorResponse(fmt.Sprintf("unrecognized credential_type: %q, not one of %#v", credentialType, allowedCredentialTypes)), nil
		}
		roleEntry.CredentialTypes = []string{credentialType}
	}

	if roleArnsRaw, ok := d.GetOk("role_arns"); ok {
		if legacyRole != "" {
			return logical.ErrorResponse("cannot supply deprecated role or policy parameters with role_arns"), nil
		}
		roleEntry.RoleArns = roleArnsRaw.([]string)
	}

	if policyArnsRaw, ok := d.GetOk("policy_arns"); ok {
		if legacyRole != "" {
			return logical.ErrorResponse("cannot supply deprecated role or policy parameters with policy_arns"), nil
		}
		roleEntry.PolicyArns = policyArnsRaw.([]string)
	}

	if policyDocumentRaw, ok := d.GetOk("policy_document"); ok {
		if legacyRole != "" {
			return logical.ErrorResponse("cannot supply deprecated role or policy parameters with policy_document"), nil
		}
		compacted := policyDocumentRaw.(string)
		if len(compacted) > 0 {
			compacted, err = compactJSON(policyDocumentRaw.(string))
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("cannot parse policy document: %q", policyDocumentRaw.(string))), nil
			}
		}
		roleEntry.PolicyDocument = compacted
	}

	if defaultSTSTTLRaw, ok := d.GetOk("default_sts_ttl"); ok {
		if legacyRole != "" {
			return logical.ErrorResponse("cannot supply deprecated role or policy parameters with default_sts_ttl"), nil
		}
		if !strutil.StrListContains(roleEntry.CredentialTypes, assumedRoleCred) && !strutil.StrListContains(roleEntry.CredentialTypes, federationTokenCred) {
			return logical.ErrorResponse(fmt.Sprintf("default_sts_ttl parameter only valid for %s and %s credential types", assumedRoleCred, federationTokenCred)), nil
		}
		roleEntry.DefaultSTSTTL = time.Duration(defaultSTSTTLRaw.(int)) * time.Second
	}

	if maxSTSTTLRaw, ok := d.GetOk("max_sts_ttl"); ok {
		if legacyRole != "" {
			return logical.ErrorResponse("cannot supply deprecated role or policy parameters with max_sts_ttl"), nil
		}
		if !strutil.StrListContains(roleEntry.CredentialTypes, assumedRoleCred) && !strutil.StrListContains(roleEntry.CredentialTypes, federationTokenCred) {
			return logical.ErrorResponse(fmt.Sprintf("max_sts_ttl parameter only valid for %s and %s credential types", assumedRoleCred, federationTokenCred)), nil
		}

		roleEntry.MaxSTSTTL = time.Duration(maxSTSTTLRaw.(int)) * time.Second
	}

	if roleEntry.MaxSTSTTL > 0 &&
		roleEntry.DefaultSTSTTL > 0 &&
		roleEntry.DefaultSTSTTL > roleEntry.MaxSTSTTL {
		return logical.ErrorResponse(`"default_sts_ttl" value must be less than or equal to "max_sts_ttl" value`), nil
	}

	if legacyRole != "" {
		roleEntry = upgradeLegacyPolicyEntry(legacyRole)
		if roleEntry.InvalidData != "" {
			return logical.ErrorResponse(fmt.Sprintf("unable to parse supplied data: %q", roleEntry.InvalidData)), nil
		}
		resp.AddWarning("Detected use of legacy role or policy parameter. Please upgrade to use the new parameters.")
	} else {
		roleEntry.ProhibitFlexibleCredPath = false
	}

	if len(roleEntry.CredentialTypes) == 0 {
		return logical.ErrorResponse("did not supply credential_type"), nil
	}

	if len(roleEntry.RoleArns) > 0 && !strutil.StrListContains(roleEntry.CredentialTypes, assumedRoleCred) {
		return logical.ErrorResponse(fmt.Sprintf("cannot supply role_arns when credential_type isn't %s", assumedRoleCred)), nil
	}
	if len(roleEntry.PolicyArns) > 0 && !strutil.StrListContains(roleEntry.CredentialTypes, iamUserCred) {
		return logical.ErrorResponse(fmt.Sprintf("cannot supply policy_arns when credential_type isn't %s", iamUserCred)), nil
	}

	err = setAwsRole(ctx, req.Storage, roleName, roleEntry)
	if err != nil {
		return nil, err
	}

	if len(resp.Warnings) == 0 {
		return nil, nil
	}

	return &resp, nil
}

func (b *backend) roleRead(ctx context.Context, s logical.Storage, roleName string, shouldLock bool) (*awsRoleEntry, error) {
	if roleName == "" {
		return nil, fmt.Errorf("missing role name")
	}
	if shouldLock {
		b.roleMutex.RLock()
	}
	entry, err := s.Get(ctx, "role/"+roleName)
	if shouldLock {
		b.roleMutex.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	var roleEntry awsRoleEntry
	if entry != nil {
		if err := entry.DecodeJSON(&roleEntry); err != nil {
			return nil, err
		}
		return &roleEntry, nil
	}

	if shouldLock {
		b.roleMutex.Lock()
		defer b.roleMutex.Unlock()
	}
	entry, err = s.Get(ctx, "role/"+roleName)
	if err != nil {
		return nil, err
	}

	if entry != nil {
		if err := entry.DecodeJSON(&roleEntry); err != nil {
			return nil, err
		}
		return &roleEntry, nil
	}

	legacyEntry, err := s.Get(ctx, "policy/"+roleName)
	if err != nil {
		return nil, err
	}
	if legacyEntry == nil {
		return nil, nil
	}

	newRoleEntry := upgradeLegacyPolicyEntry(string(legacyEntry.Value))
	if b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby) {
		err = setAwsRole(ctx, s, roleName, newRoleEntry)
		if err != nil {
			return nil, err
		}
		// This can leave legacy data around in the policy/ path if it fails for some reason,
		// but should be pretty rare for this to fail but prior writes to succeed, so not worrying
		// about cleaning it up in case of error
		err = s.Delete(ctx, "policy/"+roleName)
		if err != nil {
			return nil, err
		}
	}
	return newRoleEntry, nil
}

func upgradeLegacyPolicyEntry(entry string) *awsRoleEntry {
	var newRoleEntry *awsRoleEntry
	if strings.HasPrefix(entry, "arn:") {
		parsedArn, err := arn.Parse(entry)
		if err != nil {
			newRoleEntry = &awsRoleEntry{
				InvalidData: entry,
				Version:     1,
			}
			return newRoleEntry
		}
		resourceParts := strings.Split(parsedArn.Resource, "/")
		resourceType := resourceParts[0]
		switch resourceType {
		case "role":
			newRoleEntry = &awsRoleEntry{
				CredentialTypes:          []string{assumedRoleCred},
				RoleArns:                 []string{entry},
				ProhibitFlexibleCredPath: true,
				Version:                  1,
			}
		case "policy":
			newRoleEntry = &awsRoleEntry{
				CredentialTypes:          []string{iamUserCred},
				PolicyArns:               []string{entry},
				ProhibitFlexibleCredPath: true,
				Version:                  1,
			}
		default:
			newRoleEntry = &awsRoleEntry{
				InvalidData: entry,
				Version:     1,
			}
		}
	} else {
		compacted, err := compactJSON(entry)
		if err != nil {
			newRoleEntry = &awsRoleEntry{
				InvalidData: entry,
				Version:     1,
			}
		} else {
			// unfortunately, this is ambiguous between the cred types, so allow both
			newRoleEntry = &awsRoleEntry{
				CredentialTypes:          []string{iamUserCred, federationTokenCred},
				PolicyDocument:           compacted,
				ProhibitFlexibleCredPath: true,
				Version:                  1,
			}
		}
	}

	return newRoleEntry
}

func setAwsRole(ctx context.Context, s logical.Storage, roleName string, roleEntry *awsRoleEntry) error {
	if roleName == "" {
		return fmt.Errorf("empty role name")
	}
	if roleEntry == nil {
		return fmt.Errorf("nil roleEntry")
	}
	entry, err := logical.StorageEntryJSON("role/"+roleName, roleEntry)
	if err != nil {
		return err
	}
	if entry == nil {
		return fmt.Errorf("nil result when writing to storage")
	}
	if err := s.Put(ctx, entry); err != nil {
		return err
	}
	return nil
}

type awsRoleEntry struct {
	CredentialTypes          []string      `json:"credential_types"`                      // Entries must all be in the set of ("iam_user", "assumed_role", "federation_token")
	PolicyArns               []string      `json:"policy_arns"`                           // ARNs of managed policies to attach to an IAM user
	RoleArns                 []string      `json:"role_arns"`                             // ARNs of roles to assume for AssumedRole credentials
	PolicyDocument           string        `json:"policy_document"`                       // JSON-serialized inline policy to attach to IAM users and/or to specify as the Policy parameter in AssumeRole calls
	InvalidData              string        `json:"invalid_data,omitempty"`                // Invalid role data. Exists to support converting the legacy role data into the new format
	ProhibitFlexibleCredPath bool          `json:"prohibit_flexible_cred_path,omitempty"` // Disallow accessing STS credentials via the creds path and vice verse
	Version                  int           `json:"version"`                               // Version number of the role format
	DefaultSTSTTL            time.Duration `json:"default_sts_ttl"`                       // Default TTL for STS credentials
	MaxSTSTTL                time.Duration `json:"max_sts_ttl"`                           // Max allowed TTL for STS credentials
}

func (r *awsRoleEntry) toResponseData() map[string]interface{} {
	respData := map[string]interface{}{
		"credential_type": strings.Join(r.CredentialTypes, ","),
		"policy_arns":     r.PolicyArns,
		"role_arns":       r.RoleArns,
		"policy_document": r.PolicyDocument,
		"default_sts_ttl": int64(r.DefaultSTSTTL.Seconds()),
		"max_sts_ttl":     int64(r.MaxSTSTTL.Seconds()),
	}
	if r.InvalidData != "" {
		respData["invalid_data"] = r.InvalidData
	}
	return respData
}

func compactJSON(input string) (string, error) {
	var compacted bytes.Buffer
	err := json.Compact(&compacted, []byte(input))
	return compacted.String(), err
}

const (
	assumedRoleCred     = "assumed_role"
	iamUserCred         = "iam_user"
	federationTokenCred = "federation_token"
)

const pathListRolesHelpSyn = `List the existing roles in this backend`

const pathListRolesHelpDesc = `Roles will be listed by the role name.`

const pathRolesHelpSyn = `
Read, write and reference IAM policies that access keys can be made for.
`

const pathRolesHelpDesc = `
This path allows you to read and write roles that are used to
create access keys. These roles are associated with IAM policies that
map directly to the route to read the access keys. For example, if the
backend is mounted at "aws" and you create a role at "aws/roles/deploy"
then a user could request access credentials at "aws/creds/deploy".

You can either supply a user inline policy (via the policy argument), or
provide a reference to an existing AWS policy by supplying the full arn
reference (via the arn argument). Inline user policies written are normal
IAM policies. Vault will not attempt to parse these except to validate
that they're basic JSON. No validation is performed on arn references.

To validate the keys, attempt to read an access key after writing the policy.
`
