package aws

import (
	"fmt"

	"github.com/hashicorp/aws-sdk-go/aws"
	"github.com/hashicorp/aws-sdk-go/gen/iam"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

func pathUser(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `creds/(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathUserRead,
		},

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}
}

func (b *backend) pathUserRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	policyName := d.Get("name").(string)

	// Read the policy
	policy, err := req.Storage.Get("policy/" + policyName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %s", err)
	}
	if policy == nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Role '%s' not found", policyName)), nil
	}

	// Use the helper to create the secret
	return b.secretAccessKeysCreate(
		req.Storage, req.DisplayName, policyName, string(policy.Value))
}

func pathUserRollback(req *logical.Request, _kind string, data interface{}) error {
	var entry walUser
	if err := mapstructure.Decode(data, &entry); err != nil {
		return err
	}
	username := entry.UserName

	// Get the client
	client, err := clientIAM(req.Storage)
	if err != nil {
		return err
	}

	// Get information about this user
	groupsResp, err := client.ListGroupsForUser(&iam.ListGroupsForUserRequest{
		UserName: aws.String(username),
		MaxItems: aws.Integer(1000),
	})
	if err != nil {
		return err
	}
	groups := groupsResp.Groups

	policiesResp, err := client.ListUserPolicies(&iam.ListUserPoliciesRequest{
		UserName: aws.String(username),
		MaxItems: aws.Integer(1000),
	})
	if err != nil {
		return err
	}
	policies := policiesResp.PolicyNames

	keysResp, err := client.ListAccessKeys(&iam.ListAccessKeysRequest{
		UserName: aws.String(username),
		MaxItems: aws.Integer(1000),
	})
	if err != nil {
		return err
	}
	keys := keysResp.AccessKeyMetadata

	// Revoke all keys
	for _, k := range keys {
		err = client.DeleteAccessKey(&iam.DeleteAccessKeyRequest{
			AccessKeyID: k.AccessKeyID,
			UserName:    aws.String(username),
		})
		if err != nil {
			return err
		}
	}

	// Delete any policies
	for _, p := range policies {
		err = client.DeleteUserPolicy(&iam.DeleteUserPolicyRequest{
			UserName:   aws.String(username),
			PolicyName: aws.String(p),
		})
		if err != nil {
			return err
		}
	}

	// Remove the user from all their groups
	for _, g := range groups {
		err = client.RemoveUserFromGroup(&iam.RemoveUserFromGroupRequest{
			GroupName: g.GroupName,
			UserName:  aws.String(username),
		})
		if err != nil {
			return err
		}
	}

	// Delete the user
	err = client.DeleteUser(&iam.DeleteUserRequest{
		UserName: aws.String(username),
	})
	if err != nil {
		return err
	}

	return nil
}

type walUser struct {
	UserName string
}

const pathUserHelpSyn = `
Generate an access key pair for a specific role.
`

const pathUserHelpDesc = `
This path will generate a new, never before used key pair for
accessing AWS. The IAM policy used to back this key pair will be
the "name" parameter. For example, if this backend is mounted at "aws",
then "aws/creds/deploy" would generate access keys for the "deploy" role.

The access keys will have a lease associated with them. The access keys
can be revoked by using the lease ID.
`
