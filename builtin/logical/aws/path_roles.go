// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

var userPathRegex = regexp.MustCompile(`^\/([\x21-\x7F]{0,510}\/)?$`)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationSuffix: "roles",
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameWithAtRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationSuffix: "role",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role",
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Role Name",
				},
			},

			"credential_type": {
				Type:        framework.TypeString,
				Description: fmt.Sprintf("Type of credential to retrieve. Must be one of %s, %s, %s, or %s", assumedRoleCred, iamUserCred, federationTokenCred, sessionTokenCred),
			},

			"role_arns": {
				Type:        framework.TypeCommaStringSlice,
				Description: "ARNs of AWS roles allowed to be assumed. Only valid when credential_type is " + assumedRoleCred,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Role ARNs",
				},
			},

			"policy_arns": {
				Type: framework.TypeCommaStringSlice,
				Description: fmt.Sprintf(`ARNs of AWS policies. Behavior varies by credential_type. When credential_type is
%s, then it will attach the specified policies to the generated IAM user.
When credential_type is %s or %s, the policies will be passed as the
PolicyArns parameter, acting as a filter on permissions available.`, iamUserCred, assumedRoleCred, federationTokenCred),
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Policy ARNs",
				},
			},

			"policy_document": {
				Type: framework.TypeString,
				Description: `JSON-encoded IAM policy document. Behavior varies by credential_type. When credential_type is
iam_user, then it will attach the contents of the policy_document to the IAM
user generated. When credential_type is assumed_role or federation_token, this
will be passed in as the Policy parameter to the AssumeRole or
GetFederationToken API call, acting as a filter on permissions available.`,
			},

			"iam_groups": {
				Type: framework.TypeCommaStringSlice,
				Description: `Names of IAM groups that generated IAM users will be added to. For a credential
type of assumed_role or federation_token, the policies sent to the
corresponding AWS call (sts:AssumeRole or sts:GetFederation) will be the
policies from each group in iam_groups combined with the policy_document
and policy_arns parameters.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "IAM Groups",
					Value: "group1,group2",
				},
			},

			"iam_tags": {
				Type: framework.TypeKVPairs,
				Description: `IAM tags to be set for any users created by this role. These must be presented
as Key-Value pairs. This can be represented as a map or a list of equal sign
delimited key pairs.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "IAM Tags",
					Value: "[key1=value1, key2=value2]",
				},
			},
			"session_tags": {
				Type: framework.TypeKVPairs,
				Description: fmt.Sprintf(`Session tags to be set for %q creds created by this role. These must be presented
as Key-Value pairs. This can be represented as a map or a list of equal sign
delimited key pairs.`, assumedRoleCred),
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Session Tags",
					Value: "[key1=value1, key2=value2]",
				},
			},
			"external_id": {
				Type:        framework.TypeString,
				Description: "External ID to set when assuming the role; only valid when credential_type is " + assumedRoleCred,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "External ID",
				},
			},
			"default_sts_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: fmt.Sprintf("Default TTL for %s, %s, and %s credential types when no TTL is explicitly requested with the credentials", assumedRoleCred, federationTokenCred, sessionTokenCred),
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Default STS TTL",
				},
			},

			"max_sts_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: fmt.Sprintf("Max allowed TTL for %s, %s, and %s credential types", assumedRoleCred, federationTokenCred, sessionTokenCred),
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Max STS TTL",
				},
			},

			"permissions_boundary_arn": {
				Type:        framework.TypeString,
				Description: "ARN of an IAM policy to attach as a permissions boundary on IAM user credentials; only valid when credential_type is" + iamUserCred,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Permissions Boundary ARN",
				},
			},

			"arn": {
				Type:        framework.TypeString,
				Description: `Use role_arns or policy_arns instead.`,
				Deprecated:  true,
			},

			"policy": {
				Type:        framework.TypeString,
				Description: "Use policy_document instead.",
				Deprecated:  true,
			},

			"user_path": {
				Type:        framework.TypeString,
				Description: "Path for IAM User. Only valid when credential_type is " + iamUserCred,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "User Path",
					Value: "/",
				},
				Default: "/",
			},

			"mfa_serial_number": {
				Type: framework.TypeString,
				Description: fmt.Sprintf(`Identification number or ARN of the MFA device associated with the root config user. Only valid
when credential_type is %s. This is only required when the IAM user has an MFA device configured.`, sessionTokenCred),
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "MFA Device Serial Number",
				},
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
		roleEntry.CredentialTypes = []string{credentialTypeRaw.(string)}
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
		roleEntry.DefaultSTSTTL = time.Duration(defaultSTSTTLRaw.(int)) * time.Second
	}

	if maxSTSTTLRaw, ok := d.GetOk("max_sts_ttl"); ok {
		if legacyRole != "" {
			return logical.ErrorResponse("cannot supply deprecated role or policy parameters with max_sts_ttl"), nil
		}
		roleEntry.MaxSTSTTL = time.Duration(maxSTSTTLRaw.(int)) * time.Second
	}

	if userPathRaw, ok := d.GetOk("user_path"); ok {
		if legacyRole != "" {
			return logical.ErrorResponse("cannot supply deprecated role or policy parameters with user_path"), nil
		}

		roleEntry.UserPath = userPathRaw.(string)
	}

	if permissionsBoundaryARNRaw, ok := d.GetOk("permissions_boundary_arn"); ok {
		if legacyRole != "" {
			return logical.ErrorResponse("cannot supply deprecated role or policy parameters with permissions_boundary_arn"), nil
		}
		roleEntry.PermissionsBoundaryARN = permissionsBoundaryARNRaw.(string)
	}

	if iamGroups, ok := d.GetOk("iam_groups"); ok {
		roleEntry.IAMGroups = iamGroups.([]string)
	}

	if iamTags, ok := d.GetOk("iam_tags"); ok {
		roleEntry.IAMTags = iamTags.(map[string]string)
	}

	if serialNumber, ok := d.GetOk("mfa_serial_number"); ok {
		roleEntry.SerialNumber = serialNumber.(string)
	}

	if sessionTags, ok := d.GetOk("session_tags"); ok {
		roleEntry.SessionTags = sessionTags.(map[string]string)
	}

	if externalID, ok := d.GetOk("external_id"); ok {
		roleEntry.ExternalID = externalID.(string)
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

	err = roleEntry.validate()
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error(s) validating supplied role data: %q", err)), nil
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

func validateAWSManagedPolicy(policyARN string) error {
	parsedARN, err := arn.Parse(policyARN)
	if err != nil {
		return err
	}
	if parsedARN.Service != "iam" {
		return fmt.Errorf("expected a service of iam but got %s", parsedARN.Service)
	}
	if !strings.HasPrefix(parsedARN.Resource, "policy/") {
		return fmt.Errorf("expected a resource type of policy but got %s", parsedARN.Resource)
	}
	return nil
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
	CredentialTypes          []string          `json:"credential_types"`                      // Entries must all be in the set of ("iam_user", "assumed_role", "federation_token")
	PolicyArns               []string          `json:"policy_arns"`                           // ARNs of managed policies to attach to an IAM user
	RoleArns                 []string          `json:"role_arns"`                             // ARNs of roles to assume for AssumedRole credentials
	PolicyDocument           string            `json:"policy_document"`                       // JSON-serialized inline policy to attach to IAM users and/or to specify as the Policy parameter in AssumeRole calls
	IAMGroups                []string          `json:"iam_groups"`                            // Names of IAM groups that generated IAM users will be added to
	IAMTags                  map[string]string `json:"iam_tags"`                              // IAM tags that will be added to the generated IAM users
	SessionTags              map[string]string `json:"session_tags"`                          // Session tags that will be added as Tags parameter in AssumedRole calls
	ExternalID               string            `json:"external_id"`                           // External ID to added as ExternalID in AssumeRole calls
	InvalidData              string            `json:"invalid_data,omitempty"`                // Invalid role data. Exists to support converting the legacy role data into the new format
	ProhibitFlexibleCredPath bool              `json:"prohibit_flexible_cred_path,omitempty"` // Disallow accessing STS credentials via the creds path and vice verse
	Version                  int               `json:"version"`                               // Version number of the role format
	DefaultSTSTTL            time.Duration     `json:"default_sts_ttl"`                       // Default TTL for STS credentials
	MaxSTSTTL                time.Duration     `json:"max_sts_ttl"`                           // Max allowed TTL for STS credentials
	UserPath                 string            `json:"user_path"`                             // The path for the IAM user when using "iam_user" credential type
	PermissionsBoundaryARN   string            `json:"permissions_boundary_arn"`              // ARN of an IAM policy to attach as a permissions boundary
	SerialNumber             string            `json:"mfa_serial_number"`                     // Serial number or ARN of the MFA device
}

func (r *awsRoleEntry) toResponseData() map[string]interface{} {
	respData := map[string]interface{}{
		"credential_type":          strings.Join(r.CredentialTypes, ","),
		"policy_arns":              r.PolicyArns,
		"role_arns":                r.RoleArns,
		"policy_document":          r.PolicyDocument,
		"iam_groups":               r.IAMGroups,
		"iam_tags":                 r.IAMTags,
		"session_tags":             r.SessionTags,
		"external_id":              r.ExternalID,
		"default_sts_ttl":          int64(r.DefaultSTSTTL.Seconds()),
		"max_sts_ttl":              int64(r.MaxSTSTTL.Seconds()),
		"user_path":                r.UserPath,
		"permissions_boundary_arn": r.PermissionsBoundaryARN,
		"mfa_serial_number":        r.SerialNumber,
	}

	if r.InvalidData != "" {
		respData["invalid_data"] = r.InvalidData
	}
	return respData
}

func (r *awsRoleEntry) validate() error {
	var errors *multierror.Error

	if len(r.CredentialTypes) == 0 {
		errors = multierror.Append(errors, fmt.Errorf("did not supply credential_type"))
	}

	allowedCredentialTypes := []string{iamUserCred, assumedRoleCred, federationTokenCred, sessionTokenCred}
	for _, credType := range r.CredentialTypes {
		if !strutil.StrListContains(allowedCredentialTypes, credType) {
			errors = multierror.Append(errors, fmt.Errorf("unrecognized credential type: %s", credType))
		}
	}

	if r.DefaultSTSTTL != 0 && !strutil.StrListContains(r.CredentialTypes, assumedRoleCred) && !strutil.StrListContains(r.CredentialTypes, federationTokenCred) && !strutil.StrListContains(r.CredentialTypes, sessionTokenCred) {
		errors = multierror.Append(errors, fmt.Errorf("default_sts_ttl parameter only valid for %s, %s, and %s credential types", assumedRoleCred, federationTokenCred, sessionTokenCred))
	}

	if r.MaxSTSTTL != 0 && !strutil.StrListContains(r.CredentialTypes, assumedRoleCred) && !strutil.StrListContains(r.CredentialTypes, federationTokenCred) && !strutil.StrListContains(r.CredentialTypes, sessionTokenCred) {
		errors = multierror.Append(errors, fmt.Errorf("max_sts_ttl parameter only valid for %s, %s, and %s credential types", assumedRoleCred, federationTokenCred, sessionTokenCred))
	}

	if r.MaxSTSTTL > 0 &&
		r.DefaultSTSTTL > 0 &&
		r.DefaultSTSTTL > r.MaxSTSTTL {
		errors = multierror.Append(errors, fmt.Errorf(`"default_sts_ttl" value must be less than or equal to "max_sts_ttl" value`))
	}

	if r.UserPath != "" {
		if !strutil.StrListContains(r.CredentialTypes, iamUserCred) {
			errors = multierror.Append(errors, fmt.Errorf("user_path parameter only valid for %s credential type", iamUserCred))
		}
		if !userPathRegex.MatchString(r.UserPath) {
			errors = multierror.Append(errors, fmt.Errorf("the specified value for user_path is invalid. It must match %q regexp", userPathRegex.String()))
		}
	}

	if r.PermissionsBoundaryARN != "" {
		if !strutil.StrListContains(r.CredentialTypes, iamUserCred) {
			errors = multierror.Append(errors, fmt.Errorf("cannot supply permissions_boundary_arn when credential_type isn't %s", iamUserCred))
		}
		if err := validateAWSManagedPolicy(r.PermissionsBoundaryARN); err != nil {
			errors = multierror.Append(fmt.Errorf("invalid permissions_boundary_arn parameter: %v", err))
		}
	}

	if (r.PolicyDocument != "" || len(r.PolicyArns) != 0) && strutil.StrListContains(r.CredentialTypes, sessionTokenCred) {
		errors = multierror.Append(errors, fmt.Errorf("cannot supply a policy or role when using credential_type %s", sessionTokenCred))
	}

	if len(r.RoleArns) > 0 && !strutil.StrListContains(r.CredentialTypes, assumedRoleCred) {
		errors = multierror.Append(errors, fmt.Errorf("cannot supply role_arns when credential_type isn't %s", assumedRoleCred))
	}

	if len(r.SessionTags) > 0 && !strutil.StrListContains(r.CredentialTypes, assumedRoleCred) {
		errors = multierror.Append(errors, fmt.Errorf("cannot supply session_tags when credential_type isn't %s", assumedRoleCred))
	}

	if r.ExternalID != "" && !strutil.StrListContains(r.CredentialTypes, assumedRoleCred) {
		errors = multierror.Append(errors, fmt.Errorf("cannot supply external_id when credential_type isn't %s", assumedRoleCred))
	}

	return errors.ErrorOrNil()
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
	sessionTokenCred    = "session_token"
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
