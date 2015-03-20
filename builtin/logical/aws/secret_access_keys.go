package aws

import (
	"fmt"

	"github.com/hashicorp/aws-sdk-go/aws"
	"github.com/hashicorp/aws-sdk-go/gen/iam"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretAccessKeyType = "access_keys"

func secretAccessKeys() *framework.Secret {
	return &framework.Secret{
		Type: SecretAccessKeyType,
		Fields: map[string]*framework.FieldSchema{
			"access_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Access Key",
			},

			"secret_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Secret Key",
			},
		},

		Revoke: secretAccessKeysRevoke,
	}
}

func secretAccessKeysRevoke(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}

	// Get the client
	client, err := clientIAM(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Get information about this user
	groupsResp, err := client.ListGroupsForUser(&iam.ListGroupsForUserRequest{
		UserName: aws.String(username),
		MaxItems: aws.Integer(1000),
	})
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	groups := groupsResp.Groups

	policiesResp, err := client.ListUserPolicies(&iam.ListUserPoliciesRequest{
		UserName: aws.String(username),
		MaxItems: aws.Integer(1000),
	})
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	policies := policiesResp.PolicyNames

	// Revoke it!
	err = client.DeleteAccessKey(&iam.DeleteAccessKeyRequest{
		AccessKeyID: aws.String(d.Get("access_key").(string)),
		UserName:    aws.String(username),
	})
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Delete any policies
	for _, p := range policies {
		err = client.DeleteUserPolicy(&iam.DeleteUserPolicyRequest{
			UserName:   aws.String(username),
			PolicyName: aws.String(p),
		})
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	// Remove the user from all their groups
	for _, g := range groups {
		err = client.RemoveUserFromGroup(&iam.RemoveUserFromGroupRequest{
			GroupName: g.GroupName,
			UserName:  aws.String(username),
		})
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	// Delete the user
	err = client.DeleteUser(&iam.DeleteUserRequest{
		UserName: aws.String(username),
	})
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	return nil, nil
}
