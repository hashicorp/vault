package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

func pathUser(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
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
	groupsResp, err := client.ListGroupsForUser(&iam.ListGroupsForUserInput{
		UserName: aws.String(username),
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return err
	}
	groups := groupsResp.Groups

	// Inline (user) policies
	policiesResp, err := client.ListUserPolicies(&iam.ListUserPoliciesInput{
		UserName: aws.String(username),
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return err
	}
	policies := policiesResp.PolicyNames

	// Attached managed policies
	manPoliciesResp, err := client.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
		UserName: aws.String(username),
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return err
	}
	manPolicies := manPoliciesResp.AttachedPolicies

	keysResp, err := client.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: aws.String(username),
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return err
	}
	keys := keysResp.AccessKeyMetadata

	// Revoke all keys
	for _, k := range keys {
		_, err = client.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: k.AccessKeyId,
			UserName:    aws.String(username),
		})
		if err != nil {
			return err
		}
	}

	// Detach managed policies
	for _, p := range manPolicies {
		_, err = client.DetachUserPolicy(&iam.DetachUserPolicyInput{
			UserName:  aws.String(username),
			PolicyArn: p.PolicyArn,
		})
		if err != nil {
			return err
		}
	}

	// Delete any inline (user) policies
	for _, p := range policies {
		_, err = client.DeleteUserPolicy(&iam.DeleteUserPolicyInput{
			UserName:   aws.String(username),
			PolicyName: p,
		})
		if err != nil {
			return err
		}
	}

	// Remove the user from all their groups
	for _, g := range groups {
		_, err = client.RemoveUserFromGroup(&iam.RemoveUserFromGroupInput{
			GroupName: g.GroupName,
			UserName:  aws.String(username),
		})
		if err != nil {
			return err
		}
	}

	// Delete the user
	_, err = client.DeleteUser(&iam.DeleteUserInput{
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
