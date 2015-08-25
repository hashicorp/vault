package aws

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const SecretAccessKeyType = "access_keys"

func secretAccessKeys(b *backend) *framework.Secret {
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

		DefaultDuration:    1 * time.Hour,
		DefaultGracePeriod: 10 * time.Minute,

		Renew:  b.secretAccessKeysRenew,
		Revoke: secretAccessKeysRevoke,
	}
}

func (b *backend) secretAccessKeysCreate(
	s logical.Storage,
	displayName, policyName string, policy string) (*logical.Response, error) {
	client, err := clientIAM(s)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Generate a random username. We don't put the policy names in the
	// username because the AWS console makes it pretty easy to see that.
	username := fmt.Sprintf("vault-%s-%d-%d", normalizeDisplayName(displayName), time.Now().Unix(), rand.Int31n(10000))

	// Write to the WAL that this user will be created. We do this before
	// the user is created because if switch the order then the WAL put
	// can fail, which would put us in an awkward position: we have a user
	// we need to rollback but can't put the WAL entry to do the rollback.
	walId, err := framework.PutWAL(s, "user", &walUser{
		UserName: username,
	})
	if err != nil {
		return nil, fmt.Errorf("Error writing WAL entry: %s", err)
	}

	// Create the user
	_, err = client.CreateUser(&iam.CreateUserInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error creating IAM user: %s", err)), nil
	}

	// Add the user to all the groups
	_, err = client.PutUserPolicy(&iam.PutUserPolicyInput{
		UserName:       aws.String(username),
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policy),
	})
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error adding user to group: %s", err)), nil
	}

	// Create the keys
	keyResp, err := client.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error creating access keys: %s", err)), nil
	}

	// Remove the WAL entry, we succeeded! If we fail, we don't return
	// the secret because it'll get rolled back anyways, so we have to return
	// an error here.
	if err := framework.DeleteWAL(s, walId); err != nil {
		return nil, fmt.Errorf("Failed to commit WAL entry: %s", err)
	}

	// Return the info!
	return b.Secret(SecretAccessKeyType).Response(map[string]interface{}{
		"access_key": *keyResp.AccessKey.AccessKeyId,
		"secret_key": *keyResp.AccessKey.SecretAccessKey,
	}, map[string]interface{}{
		"username": username,
		"policy":   policy,
	}), nil
}

func (b *backend) secretAccessKeysRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	lease, err := b.Lease(req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{Lease: 1 * time.Hour}
	}

	f := framework.LeaseExtend(lease.Lease, lease.LeaseMax, false)
	return f(req, d)
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

	// Use the user rollback mechanism to delete this user
	err := pathUserRollback(req, "user", map[string]interface{}{
		"username": username,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func normalizeDisplayName(displayName string) string {
	re := regexp.MustCompile("[^a-zA-Z+=,.@_-]")
	return re.ReplaceAllString(displayName, "_")
}
