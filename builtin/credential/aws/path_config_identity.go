package awsauth

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigIdentity(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/identity$",
		Fields: map[string]*framework.FieldSchema{
			"iam_alias": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     identityAliasIAMUniqueID,
				Description: fmt.Sprintf("Configure how the AWS auth method generates entity aliases when using IAM auth. Valid values are %q and %q", identityAliasIAMUniqueID, identityAliasIAMFullArn),
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

func pathConfigIdentityRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get(ctx, "config/identity")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	var result identityConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"iam_alias": result.IAMAlias,
		},
	}, nil
}

func pathConfigIdentityUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var configEntry identityConfig

	iamAliasRaw, ok := data.GetOk("iam_alias")
	if ok {
		iamAlias := iamAliasRaw.(string)
		allowedIAMAliasValues := []string{identityAliasIAMUniqueID, identityAliasIAMFullArn}
		if !strutil.StrListContains(allowedIAMAliasValues, iamAlias) {
			return logical.ErrorResponse(fmt.Sprintf("iam_alias of %q not in set of allowed values: %v", iamAlias, allowedIAMAliasValues)), nil
		}
		configEntry.IAMAlias = iamAlias
		entry, err := logical.StorageEntryJSON("config/identity", configEntry)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

type identityConfig struct {
	IAMAlias string `json:"iam_alias"`
}

const identityAliasIAMUniqueID = "unique_id"
const identityAliasIAMFullArn = "full_arn"

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
