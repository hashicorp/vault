package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

func pathUser(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "(creds|sts)/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
			"role_arn": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "ARN of role to assume when credential_type is " + assumedRoleCred,
			},
			"ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: "Lifetime of the returned credentials in seconds",
				Default:     3600,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathCredsRead,
			logical.UpdateOperation: b.pathCredsRead,
		},

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}
}

func (b *backend) pathCredsRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)

	// Read the policy
	role, err := b.roleRead(ctx, req.Storage, roleName, true)
	if err != nil {
		return nil, errwrap.Wrapf("error retrieving role: {{err}}", err)
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Role '%s' not found", roleName)), nil
	}

	var ttl int64
	ttlRaw, ok := d.GetOk("ttl")
	switch {
	case ok:
		ttl = int64(ttlRaw.(int))
	case role.DefaultSTSTTL > 0:
		ttl = int64(role.DefaultSTSTTL.Seconds())
	default:
		ttl = int64(d.Get("ttl").(int))
	}

	var maxTTL int64
	if role.MaxSTSTTL > 0 {
		maxTTL = int64(role.MaxSTSTTL.Seconds())
	} else {
		maxTTL = int64(b.System().MaxLeaseTTL().Seconds())
	}

	if ttl > maxTTL {
		ttl = maxTTL
	}

	roleArn := d.Get("role_arn").(string)

	var credentialType string
	switch {
	case len(role.CredentialTypes) == 1:
		credentialType = role.CredentialTypes[0]
	// There is only one way for the CredentialTypes to contain more than one entry, and that's an upgrade path
	// where it contains iamUserCred and federationTokenCred
	// This ambiguity can be resolved based on req.Path, so resolve it assuming CredentialTypes only has those values
	case len(role.CredentialTypes) > 1:
		if strings.HasPrefix(req.Path, "creds") {
			credentialType = iamUserCred
		} else {
			credentialType = federationTokenCred
		}
		// sanity check on the assumption above
		if !strutil.StrListContains(role.CredentialTypes, credentialType) {
			return logical.ErrorResponse(fmt.Sprintf("requested credential type %q not in allowed credential types %#v", credentialType, role.CredentialTypes)), nil
		}
	}

	// creds requested through the sts path shouldn't be allowed to get iamUserCred type creds
	// when the role is created from legacy data because they might have more privileges in AWS.
	// See https://github.com/hashicorp/vault/issues/4229#issuecomment-380316788 for details.
	if role.ProhibitFlexibleCredPath {
		if credentialType == iamUserCred && strings.HasPrefix(req.Path, "sts") {
			return logical.ErrorResponse(fmt.Sprintf("attempted to retrieve %s credentials through the sts path; this is not allowed for legacy roles", iamUserCred)), nil
		}
		if credentialType != iamUserCred && strings.HasPrefix(req.Path, "creds") {
			return logical.ErrorResponse(fmt.Sprintf("attempted to retrieve %s credentials through the creds path; this is not allowed for legacy roles", credentialType)), nil
		}
	}

	switch credentialType {
	case iamUserCred:
		return b.secretAccessKeysCreate(ctx, req.Storage, req.DisplayName, roleName, role)
	case assumedRoleCred:
		switch {
		case roleArn == "":
			if len(role.RoleArns) != 1 {
				return logical.ErrorResponse("did not supply a role_arn parameter and unable to determine one"), nil
			}
			roleArn = role.RoleArns[0]
		case !strutil.StrListContains(role.RoleArns, roleArn):
			return logical.ErrorResponse(fmt.Sprintf("role_arn %q not in allowed role arns for Vault role %q", roleArn, roleName)), nil
		}
		return b.assumeRole(ctx, req.Storage, req.DisplayName, roleName, roleArn, role.PolicyDocument, ttl)
	case federationTokenCred:
		return b.secretTokenCreate(ctx, req.Storage, req.DisplayName, roleName, role.PolicyDocument, ttl)
	default:
		return logical.ErrorResponse(fmt.Sprintf("unknown credential_type: %q", credentialType)), nil
	}
}

func (b *backend) pathUserRollback(ctx context.Context, req *logical.Request, _kind string, data interface{}) error {
	var entry walUser
	if err := mapstructure.Decode(data, &entry); err != nil {
		return err
	}
	username := entry.UserName

	// Get the client
	client, err := b.clientIAM(ctx, req.Storage)
	if err != nil {
		return err
	}

	// Get information about this user
	groupsResp, err := client.ListGroupsForUser(&iam.ListGroupsForUserInput{
		UserName: aws.String(username),
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		// This isn't guaranteed to be perfect; for example, an IAM user
		// might have gotten put into the WAL but then the IAM user creation
		// failed (e.g., Vault didn't have permissions) and then the WAL
		// deletion failed as well. Then, if Vault doesn't have access to
		// call iam:ListGroupsForUser, AWS will return an access denied error
		// and the WAL will never get cleaned up. But this is better than
		// just having Vault "forget" about a user it actually created.
		//
		// BEWARE a potential race condition -- where this is called
		// immediately after a user is created. AWS eventual consistency
		// might say the user doesn't exist when the user does in fact
		// exist, and this could cause Vault to forget about the user.
		// This won't happen if the user creation fails (because the WAL
		// minimum age is 5 minutes, and AWS eventual consistency is, in
		// practice, never that long), but it could happen if a lease holder
		// asks immediately after getting a user to revoke the lease, causing
		// Vault to leak the secret, which would be a Very Bad Thing to allow.
		// So we make sure that, if there's an associated lease, it must be at
		// least 5 minutes old as well.
		if aerr, ok := err.(awserr.Error); ok {
			acceptMissingIamUsers := false
			if req.Secret == nil || time.Since(req.Secret.IssueTime) > time.Duration(minAwsUserRollbackAge) {
				// WAL rollback
				acceptMissingIamUsers = true
			}
			if aerr.Code() == iam.ErrCodeNoSuchEntityException && acceptMissingIamUsers {
				return nil
			}
		}
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
Generate AWS credentials from a specific Vault role.
`

const pathUserHelpDesc = `
This path will generate new, never before used AWS credentials for
accessing AWS. The IAM policy used to back this key pair will be
the "name" parameter. For example, if this backend is mounted at "aws",
then "aws/creds/deploy" would generate access keys for the "deploy" role.

The access keys will have a lease associated with them. The access keys
can be revoked by using the lease ID when using the iam_user credential type.
When using AWS STS credential types (assumed_role or federation_token),
revoking the lease does not revoke the access keys.
`
