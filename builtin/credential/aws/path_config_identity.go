package awsauth

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfigIdentity(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/identity$",
		Fields: map[string]*framework.FieldSchema{
			"iam_alias": {
				Type:        framework.TypeString,
				Default:     identityAliasIAMUniqueID,
				Description: fmt.Sprintf("Configure how the AWS auth method generates entity aliases when using IAM auth. Valid values are %q, %q, and %q. Defaults to %q.", identityAliasRoleID, identityAliasIAMUniqueID, identityAliasIAMFullArn, identityAliasRoleID),
			},
			"ec2_alias": {
				Type:        framework.TypeString,
				Default:     identityAliasEC2InstanceID,
				Description: fmt.Sprintf("Configure how the AWS auth method generates entity alias when using EC2 auth. Valid values are %q, %q, and %q. Defaults to %q.", identityAliasRoleID, identityAliasEC2InstanceID, identityAliasEC2ImageID, identityAliasRoleID),
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   pathConfigIdentityRead,
			logical.UpdateOperation: pathConfigIdentityUpdate,
		},

		HelpSynopsis:    pathConfigIdentityHelpSyn,
		HelpDescription: pathConfigIdentityHelpDesc,
	}
}

func identityConfigEntry(ctx context.Context, s logical.Storage) (*identityConfig, error) {
	entryRaw, err := s.Get(ctx, "config/identity")
	if err != nil {
		return nil, err
	}

	var entry identityConfig
	if entryRaw != nil {
		if err := entryRaw.DecodeJSON(&entry); err != nil {
			return nil, err
		}
	}

	if entry.IAMAlias == "" {
		entry.IAMAlias = identityAliasRoleID
	}

	if entry.EC2Alias == "" {
		entry.EC2Alias = identityAliasRoleID
	}

	return &entry, nil
}

func pathConfigIdentityRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := identityConfigEntry(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"iam_alias": config.IAMAlias,
			"ec2_alias": config.EC2Alias,
		},
	}, nil
}

func pathConfigIdentityUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := identityConfigEntry(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	iamAliasRaw, ok := data.GetOk("iam_alias")
	if ok {
		iamAlias := iamAliasRaw.(string)
		allowedIAMAliasValues := []string{identityAliasRoleID, identityAliasIAMUniqueID, identityAliasIAMFullArn}
		if !strutil.StrListContains(allowedIAMAliasValues, iamAlias) {
			return logical.ErrorResponse(fmt.Sprintf("iam_alias of %q not in set of allowed values: %v", iamAlias, allowedIAMAliasValues)), nil
		}
		config.IAMAlias = iamAlias
	}

	ec2AliasRaw, ok := data.GetOk("ec2_alias")
	if ok {
		ec2Alias := ec2AliasRaw.(string)
		allowedEC2AliasValues := []string{identityAliasRoleID, identityAliasEC2InstanceID, identityAliasEC2ImageID}
		if !strutil.StrListContains(allowedEC2AliasValues, ec2Alias) {
			return logical.ErrorResponse(fmt.Sprintf("ec2_alias of %q not in set of allowed values: %v", ec2Alias, allowedEC2AliasValues)), nil
		}
		config.EC2Alias = ec2Alias
	}

	entry, err := logical.StorageEntryJSON("config/identity", config)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

type identityConfig struct {
	IAMAlias string `json:"iam_alias"`
	EC2Alias string `json:"ec2_alias"`
}

const identityAliasIAMUniqueID = "unique_id"
const identityAliasIAMFullArn = "full_arn"
const identityAliasEC2InstanceID = "instance_id"
const identityAliasEC2ImageID = "image_id"
const identityAliasRoleID = "role_id"

const pathConfigIdentityHelpSyn = `
Configure the way the AWS auth method interacts with the identity store
`

const pathConfigIdentityHelpDesc = `
The AWS auth backend defaults to aliasing an IAM principal's unique ID to the
identity store. This path allows users to change how Vault configures the
mapping to Identity aliases for more flexibility.

You can set the iam_alias parameter to one of the following values:

* 'unique_id': This retains Vault's default behavior
* 'full_arn': This maps the full authenticated ARN to the identity alias, e.g.,
   "arn:aws:sts::<account_id>:assumed-role/<role_name>/<role_session_name>
   This is useful where you have an identity provder that sets role_session_name
   to a known value of a person, such as a username or email address, and allows
   you to map those roles back to entries in your identity store.
`
