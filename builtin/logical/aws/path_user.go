package aws

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/aws-sdk-go/aws"
	"github.com/hashicorp/aws-sdk-go/gen/iam"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error creating access keys: %s", err)), nil
	}

	// Return the info!
	return b.Secret(SecretAccessKeyType).Response(map[string]interface{}{
		"access_key": *keyResp.AccessKey.AccessKeyID,
		"secret_key": *keyResp.AccessKey.SecretAccessKey,
	}, map[string]interface{}{
		"username": username,
	}), nil
}
