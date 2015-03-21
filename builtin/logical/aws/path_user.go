package aws

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/aws-sdk-go/aws"
	"github.com/hashicorp/aws-sdk-go/gen/iam"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

func pathUser(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the policy",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathUserRead,
		},
	}
}

func (b *backend) pathUserRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	client, err := clientIAM(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Read the policy
	policy, err := req.Storage.Get("policy/" + d.Get("name").(string))
	if err != nil {
		return nil, fmt.Errorf("error retrieving policy: %s", err)
	}
	if policy == nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Policy '%s' not found", d.Get("name").(string))), nil
	}

	// Generate a random username. We don't put the policy names in the
	// username because the AWS console makes it pretty easy to see that.
	username := fmt.Sprintf("vault-%d-%d", time.Now().Unix(), rand.Int31n(10000))

	// Write to the WAL that this user will be created. We do this before
	// the user is created because if switch the order then the WAL put
	// can fail, which would put us in an awkward position: we have a user
	// we need to rollback but can't put the WAL entry to do the rollback.
	walId, err := framework.PutWAL(req.Storage, "user", &walUser{
		UserName: username,
	})
	if err != nil {
		return nil, fmt.Errorf("Error writing WAL entry: %s", err)
	}

	// Create the user
	_, err = client.CreateUser(&iam.CreateUserRequest{
		UserName: aws.String(username),
	})
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error creating IAM user: %s", err)), nil
	}

	// Add the user to all the groups
	err = client.PutUserPolicy(&iam.PutUserPolicyRequest{
		UserName:       aws.String(username),
		PolicyName:     aws.String(d.Get("name").(string)),
		PolicyDocument: aws.String(string(policy.Value)),
	})
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error adding user to group: %s", err)), nil
	}

	// Create the keys
	keyResp, err := client.CreateAccessKey(&iam.CreateAccessKeyRequest{
		UserName: aws.String(username),
	})
	err = fmt.Errorf("SUCK!")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error creating access keys: %s", err)), nil
	}

	// Remove the WAL entry, we succeeded! If we fail, we don't return
	// the secret because it'll get rolled back anyways, so we have to return
	// an error here.
	if err := framework.DeleteWAL(req.Storage, walId); err != nil {
		return nil, fmt.Errorf("Failed to commit WAL entry: %s", err)
	}

	// Return the info!
	return b.Secret(SecretAccessKeyType).Response(map[string]interface{}{
		"access_key": *keyResp.AccessKey.AccessKeyID,
		"secret_key": *keyResp.AccessKey.SecretAccessKey,
	}, map[string]interface{}{
		"username": username,
	}), nil
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
