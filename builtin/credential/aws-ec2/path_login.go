package awsec2

import (
	"encoding/pem"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fullsailor/pkcs7"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		Fields: map[string]*framework.FieldSchema{
			"role": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Name of the role against which the login is being attempted.
If 'role' is not specified, then the login endpoint looks for a role
bearing the name of the AMI ID of the EC2 instance that is trying to login.
If a matching role is not found, login fails.`,
			},

			"pkcs7": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "PKCS7 signature of the identity document.",
			},

			"nonce": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `The nonce created by a client of this backend. When 'disallow_reauthentication'
option is enabled on either the role or the role tag, then nonce parameter is
optional. It is a required parameter otherwise.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLoginUpdate,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

// validateInstance queries the status of the EC2 instance using AWS EC2 API and
// checks if the instance is running and is healthy.
func (b *backend) validateInstance(s logical.Storage, instanceID, region string) (*ec2.DescribeInstancesOutput, error) {
	// Create an EC2 client to pull the instance information
	ec2Client, err := b.clientEC2(s, region)
	if err != nil {
		return nil, err
	}

	status, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(instanceID),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching description for instance ID %s: %s\n", instanceID, err)
	}
	if len(status.Reservations) == 0 {
		return nil, fmt.Errorf("no reservations found in instance description")

	}
	if len(status.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("no instance details found in reservations")
	}
	if *status.Reservations[0].Instances[0].InstanceId != instanceID {
		return nil, fmt.Errorf("expected instance ID not matching the instance ID in the instance description")
	}
	if status.Reservations[0].Instances[0].State == nil {
		return nil, fmt.Errorf("instance state in instance description is nil")
	}
	if *status.Reservations[0].Instances[0].State.Code != 16 ||
		*status.Reservations[0].Instances[0].State.Name != "running" {
		return nil, fmt.Errorf("instance is not in 'running' state")
	}
	return status, nil
}

// validateMetadata matches the given client nonce and pending time with the one cached
// in the identity whitelist during the previous login. But, if reauthentication is
// disabled, login attempt is failed immediately.
func validateMetadata(clientNonce, pendingTime string, storedIdentity *whitelistIdentity, roleEntry *awsRoleEntry) error {
	// If reauthentication is disabled, doesn't matter what other metadata is provided,
	// authentication will not succeed.
	if storedIdentity.DisallowReauthentication {
		return fmt.Errorf("reauthentication is disabled")
	}

	givenPendingTime, err := time.Parse(time.RFC3339, pendingTime)
	if err != nil {
		return err
	}

	storedPendingTime, err := time.Parse(time.RFC3339, storedIdentity.PendingTime)
	if err != nil {
		return err
	}

	// When the presented client nonce does not match the cached entry, it is
	// either that a rogue client is trying to login or that a valid client
	// suffered a migration. The migration is detected via pendingTime in the
	// instance metadata, which sadly is only updated when an instance is
	// stopped and started but *not* when the instance is rebooted. If reboot
	// survivability is needed, either instrumentation to delete the instance
	// ID from the whitelist is necessary, or the client must durably store
	// the nonce.
	//
	// If the `allow_instance_migration` property of the registered role is
	// enabled, then the client nonce mismatch is ignored, as long as the
	// pending time in the presented instance identity document is newer than
	// the cached pending time. The new pendingTime is stored and used for
	// future checks.
	//
	// This is a weak criterion and hence the `allow_instance_migration` option
	// should be used with caution.
	if clientNonce != storedIdentity.ClientNonce {
		if !roleEntry.AllowInstanceMigration {
			return fmt.Errorf("client nonce mismatch")
		}
		if roleEntry.AllowInstanceMigration && !givenPendingTime.After(storedPendingTime) {
			return fmt.Errorf("client nonce mismatch and instance meta-data incorrect")
		}
	}

	// Ensure that the 'pendingTime' on the given identity document is not before the
	// 'pendingTime' that was used for previous login. This disallows old metadata documents
	// from being used to perform login.
	if givenPendingTime.Before(storedPendingTime) {
		return fmt.Errorf("instance meta-data is older than the one used for previous login")
	}
	return nil
}

// Verifies the correctness of the authenticated attributes present in the PKCS#7
// signature. After verification, extracts the instance identity document from the
// signature, parses it and returns it.
func (b *backend) parseIdentityDocument(s logical.Storage, pkcs7B64 string) (*identityDocument, error) {
	// Insert the header and footer for the signature to be able to pem decode it.
	pkcs7B64 = fmt.Sprintf("-----BEGIN PKCS7-----\n%s\n-----END PKCS7-----", pkcs7B64)

	// Decode the PEM encoded signature.
	pkcs7BER, pkcs7Rest := pem.Decode([]byte(pkcs7B64))
	if len(pkcs7Rest) != 0 {
		return nil, fmt.Errorf("failed to decode the PEM encoded PKCS#7 signature")
	}

	// Parse the signature from asn1 format into a struct.
	pkcs7Data, err := pkcs7.Parse(pkcs7BER.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the BER encoded PKCS#7 signature: %s\n", err)
	}

	// Get the public certificates that are used to verify the signature.
	// This returns a slice of certificates containing the default certificate
	// and all the registered certificates via 'config/certificate/<cert_name>' endpoint
	publicCerts, err := b.awsPublicCertificates(s)
	if err != nil {
		return nil, err
	}
	if publicCerts == nil || len(publicCerts) == 0 {
		return nil, fmt.Errorf("certificates to verify the signature are not found")
	}

	// Before calling Verify() on the PKCS#7 struct, set the certificates to be used
	// to verify the contents in the signer information.
	pkcs7Data.Certificates = publicCerts

	// Verify extracts the authenticated attributes in the PKCS#7 signature, and verifies
	// the authenticity of the content using 'dsa.PublicKey' embedded in the public certificate.
	if pkcs7Data.Verify() != nil {
		return nil, fmt.Errorf("failed to verify the signature")
	}

	// Check if the signature has content inside of it.
	if len(pkcs7Data.Content) == 0 {
		return nil, fmt.Errorf("instance identity document could not be found in the signature")
	}

	var identityDoc identityDocument
	if err := jsonutil.DecodeJSON(pkcs7Data.Content, &identityDoc); err != nil {
		return nil, err
	}

	return &identityDoc, nil
}

// pathLoginUpdate is used to create a Vault token by the EC2 instances
// by providing the pkcs7 signature of the instance identity document
// and a client created nonce. Client nonce is optional if 'disallow_reauthentication'
// option is enabled on the registered role.
func (b *backend) pathLoginUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	pkcs7B64 := data.Get("pkcs7").(string)
	if pkcs7B64 == "" {
		return logical.ErrorResponse("missing pkcs7"), nil
	}

	// Verify the signature of the identity document.
	identityDoc, err := b.parseIdentityDocument(req.Storage, pkcs7B64)
	if err != nil {
		return nil, err
	}
	if identityDoc == nil {
		return logical.ErrorResponse("failed to extract instance identity document from PKCS#7 signature"), nil
	}

	roleName := data.Get("role").(string)

	// If roleName is not supplied, a role in the name of the instance's AMI ID will be looked for.
	if roleName == "" {
		roleName = identityDoc.AmiID
	}

	// Validate the instance ID by making a call to AWS EC2 DescribeInstances API
	// and fetching the instance description. Validation succeeds only if the
	// instance is in 'running' state.
	instanceDesc, err := b.validateInstance(req.Storage, identityDoc.InstanceID, identityDoc.Region)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to verify instance ID: %s", err)), nil
	}

	// Get the entry for the role used by the instance.
	roleEntry, err := b.lockedAWSRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("entry for role '%s' not found", roleName)), nil
	}

	// Verify that the AMI ID of the instance trying to login matches the
	// AMI ID specified as a constraint on the role.
	if roleEntry.BoundAmiID != "" && identityDoc.AmiID != roleEntry.BoundAmiID {
		return logical.ErrorResponse(fmt.Sprintf("AMI ID '%s' does not belong to role '%s'", identityDoc.AmiID, roleName)), nil
	}

	// Verify that the AccountID of the instance trying to login matches the
	// AccountID specified as a constraint on the role.
	if roleEntry.BoundAccountID != "" && identityDoc.AccountID != roleEntry.BoundAccountID {
		return logical.ErrorResponse(fmt.Sprintf("Account ID '%s' does not belong to role '%s'", identityDoc.AccountID, roleName)), nil
	}

	// Check if the IAM Role ARN of the instance trying to login matches the
	// IAM Role ARN specified as a constraint on the role.
	if roleEntry.BoundIamARN != "" {
		if instanceDesc.Reservations[0].Instances[0].IamInstanceProfile == nil {
			return nil, fmt.Errorf("IAM Instance Profile in the instance description is nil")
		}
		if instanceDesc.Reservations[0].Instances[0].IamInstanceProfile.Arn == nil {
			return nil, fmt.Errorf("ARN in the instance description is nil")
		}
		iamRoleArn := *instanceDesc.Reservations[0].Instances[0].IamInstanceProfile.Arn
		if iamRoleArn != roleEntry.BoundIamARN {
			return logical.ErrorResponse(fmt.Sprintf("IAM Role ARN %s does not belong to role %s", iamRoleArn, roleName)), nil
		}
	}

	// Get the entry from the identity whitelist, if there is one.
	storedIdentity, err := whitelistIdentityEntry(req.Storage, identityDoc.InstanceID)
	if err != nil {
		return nil, err
	}

	clientNonce := data.Get("nonce").(string)

	// This is NOT a first login attempt from the client.
	if storedIdentity != nil {
		// Check if the client nonce match the cached nonce and if the pending time
		// of the identity document is not before the pending time of the document
		// with which previous login was made. If 'allow_instance_migration' is
		// enabled on the registered role, client nonce requirement is relaxed.
		if err = validateMetadata(clientNonce, identityDoc.PendingTime, storedIdentity, roleEntry); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	// Load the current values for max TTL and policies from the role entry,
	// before checking for overriding max TTL in the role tag.  The shortest
	// max TTL is used to cap the token TTL; the longest max TTL is used to
	// make the whitelist entry as long as possible as it controls for replay
	// attacks.
	shortestMaxTTL := b.System().MaxLeaseTTL()
	longestMaxTTL := b.System().MaxLeaseTTL()
	if roleEntry.MaxTTL > time.Duration(0) && roleEntry.MaxTTL < shortestMaxTTL {
		shortestMaxTTL = roleEntry.MaxTTL
	}
	if roleEntry.MaxTTL > longestMaxTTL {
		longestMaxTTL = roleEntry.MaxTTL
	}

	policies := roleEntry.Policies
	rTagMaxTTL := time.Duration(0)

	// Read this value from the role entry; however, once it's been set, do not
	// allow it to be flipped back. This prevents a role with this set to false
	// to be overridden by a role tag, then have the role tag swapped and have
	// this go back to false.
	disallowReauthentication := roleEntry.DisallowReauthentication
	if storedIdentity != nil {
		if !disallowReauthentication && storedIdentity.DisallowReauthentication {
			disallowReauthentication = true
		}
	}

	if roleEntry.RoleTag != "" {
		// Role tag is enabled on the role.

		// Overwrite the policies with the ones returned from processing the role tag.
		resp, err := b.handleRoleTagLogin(req.Storage, identityDoc, roleName, roleEntry, instanceDesc)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return logical.ErrorResponse("failed to fetch and verify the role tag"), nil
		}

		// If there are no policies on the role tag, policies on the role are inherited.
		// If policies on role tag are set, by this point, it is verified that it is a subset of the
		// policies on the role. So, apply only those.
		if len(resp.Policies) != 0 {
			policies = resp.Policies
		}

		// If roleEntry had disallowReauthentication set to 'true', do not reset it
		// to 'false' based on role tag having it not set. But, if role tag had it set,
		// be sure to override the value.
		if !disallowReauthentication {
			disallowReauthentication = resp.DisallowReauthentication
		}

		// Cache the value of role tag's max_ttl value.
		rTagMaxTTL = resp.MaxTTL

		// Scope the shortestMaxTTL to the value set on the role tag.
		if resp.MaxTTL > time.Duration(0) && resp.MaxTTL < shortestMaxTTL {
			shortestMaxTTL = resp.MaxTTL
		}
		if resp.MaxTTL > longestMaxTTL {
			longestMaxTTL = resp.MaxTTL
		}
	}

	// Save the login attempt in the identity whitelist.
	currentTime := time.Now()
	if storedIdentity == nil {
		// Role, ClientNonce and CreationTime of the identity entry,
		// once set, should never change.
		storedIdentity = &whitelistIdentity{
			Role:         roleName,
			ClientNonce:  clientNonce,
			CreationTime: currentTime,
		}
	}

	// DisallowReauthentication, PendingTime, LastUpdatedTime and ExpirationTime may change.
	storedIdentity.LastUpdatedTime = currentTime
	storedIdentity.ExpirationTime = currentTime.Add(longestMaxTTL)
	storedIdentity.PendingTime = identityDoc.PendingTime
	storedIdentity.DisallowReauthentication = disallowReauthentication

	// Performing the clientNonce empty check after determining the DisallowReauthentication
	// option. This is to make clientNonce optional when DisallowReauthentication is set.
	if clientNonce == "" && !storedIdentity.DisallowReauthentication {
		return logical.ErrorResponse("missing nonce"), nil
	}

	// Limit the nonce to a reasonable length.
	if len(clientNonce) > 128 && !storedIdentity.DisallowReauthentication {
		return logical.ErrorResponse("client nonce exceeding the limit of 128 characters"), nil
	}

	if err = setWhitelistIdentityEntry(req.Storage, identityDoc.InstanceID, storedIdentity); err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Auth: &logical.Auth{
			Policies: policies,
			Metadata: map[string]string{
				"instance_id":      identityDoc.InstanceID,
				"region":           identityDoc.Region,
				"role_tag_max_ttl": rTagMaxTTL.String(),
				"role":             roleName,
				"ami_id":           identityDoc.AmiID,
			},
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       b.System().DefaultLeaseTTL(),
			},
		},
	}

	// Cap the TTL value.
	if shortestMaxTTL < resp.Auth.TTL {
		resp.Auth.TTL = shortestMaxTTL
	}

	return resp, nil

}

// handleRoleTagLogin is used to fetch the role tag of the instance and verifies it to be correct.
// Then the policies for the login request will be set off of the role tag, if certain creteria satisfies.
func (b *backend) handleRoleTagLogin(s logical.Storage, identityDoc *identityDocument, roleName string, roleEntry *awsRoleEntry, instanceDesc *ec2.DescribeInstancesOutput) (*roleTagLoginResponse, error) {
	if identityDoc == nil {
		return nil, fmt.Errorf("nil identityDoc")
	}
	if roleEntry == nil {
		return nil, fmt.Errorf("nil roleEntry")
	}
	if instanceDesc == nil {
		return nil, fmt.Errorf("nil instanceDesc")
	}

	// Input validation on instanceDesc is not performed here considering
	// that it would have been done in validateInstance method.
	tags := instanceDesc.Reservations[0].Instances[0].Tags
	if tags == nil || len(tags) == 0 {
		return nil, fmt.Errorf("missing tag with key %s on the instance", roleEntry.RoleTag)
	}

	// Iterate through the tags attached on the instance and look for
	// a tag with its 'key' matching the expected role tag value.
	rTagValue := ""
	for _, tagItem := range tags {
		if tagItem.Key != nil && *tagItem.Key == roleEntry.RoleTag {
			rTagValue = *tagItem.Value
			break
		}
	}

	// If 'role_tag' is enabled on the role, and if a corresponding tag is not found
	// to be attached to the instance, fail.
	if rTagValue == "" {
		return nil, fmt.Errorf("missing tag with key %s on the instance", roleEntry.RoleTag)
	}

	// Parse the role tag into a struct, extract the plaintext part of it and verify its HMAC.
	rTag, err := b.parseAndVerifyRoleTagValue(s, rTagValue)
	if err != nil {
		return nil, err
	}

	// Check if the role name with which this login is being made is same
	// as the role name embedded in the tag.
	if rTag.Role != roleName {
		return nil, fmt.Errorf("role on the tag is not matching the role supplied")
	}

	// If instance_id was set on the role tag, check if the same instance is attempting to login.
	if rTag.InstanceID != "" && rTag.InstanceID != identityDoc.InstanceID {
		return nil, fmt.Errorf("role tag is being used by an unauthorized instance.")
	}

	// Check if the role tag is blacklisted.
	blacklistEntry, err := b.lockedBlacklistRoleTagEntry(s, rTagValue)
	if err != nil {
		return nil, err
	}
	if blacklistEntry != nil {
		return nil, fmt.Errorf("role tag is blacklisted")
	}

	// Ensure that the policies on the RoleTag is a subset of policies on the role
	if !strutil.StrListSubset(roleEntry.Policies, rTag.Policies) {
		return nil, fmt.Errorf("policies on the role tag must be subset of policies on the role")
	}

	return &roleTagLoginResponse{
		Policies: rTag.Policies,
		MaxTTL:   rTag.MaxTTL,
		DisallowReauthentication: rTag.DisallowReauthentication,
	}, nil
}

// pathLoginRenew is used to renew an authenticated token.
func (b *backend) pathLoginRenew(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	instanceID := req.Auth.Metadata["instance_id"]
	if instanceID == "" {
		return nil, fmt.Errorf("unable to fetch instance ID from metadata during renewal")
	}

	region := req.Auth.Metadata["region"]
	if region == "" {
		return nil, fmt.Errorf("unable to fetch region from metadata during renewal")
	}

	// Cross check that the instance is still in 'running' state
	_, err := b.validateInstance(req.Storage, instanceID, region)
	if err != nil {
		return nil, fmt.Errorf("failed to verify instance ID '%s': %s", instanceID, err)
	}

	storedIdentity, err := whitelistIdentityEntry(req.Storage, instanceID)
	if err != nil {
		return nil, err
	}
	if storedIdentity == nil {
		return nil, fmt.Errorf("failed to verify the whitelist identity entry for instance ID: %s", instanceID)
	}

	// Ensure that role entry is not deleted.
	roleEntry, err := b.lockedAWSRole(req.Storage, storedIdentity.Role)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, fmt.Errorf("role entry not found")
	}

	// If the login was made using the role tag, then max_ttl from tag
	// is cached in internal data during login and used here to cap the
	// max_ttl of renewal.
	rTagMaxTTL, err := time.ParseDuration(req.Auth.Metadata["role_tag_max_ttl"])
	if err != nil {
		return nil, err
	}

	// Re-evaluate the maxTTL bounds.
	shortestMaxTTL := b.System().MaxLeaseTTL()
	longestMaxTTL := b.System().MaxLeaseTTL()
	if roleEntry.MaxTTL > time.Duration(0) && roleEntry.MaxTTL < shortestMaxTTL {
		shortestMaxTTL = roleEntry.MaxTTL
	}
	if roleEntry.MaxTTL > longestMaxTTL {
		longestMaxTTL = roleEntry.MaxTTL
	}
	if rTagMaxTTL > time.Duration(0) && rTagMaxTTL < shortestMaxTTL {
		shortestMaxTTL = rTagMaxTTL
	}
	if rTagMaxTTL > longestMaxTTL {
		longestMaxTTL = rTagMaxTTL
	}

	// Only LastUpdatedTime and ExpirationTime change and all other fields remain the same.
	currentTime := time.Now()
	storedIdentity.LastUpdatedTime = currentTime
	storedIdentity.ExpirationTime = currentTime.Add(longestMaxTTL)

	if err = setWhitelistIdentityEntry(req.Storage, instanceID, storedIdentity); err != nil {
		return nil, err
	}

	return framework.LeaseExtend(req.Auth.TTL, shortestMaxTTL, b.System())(req, data)
}

// Struct to represent items of interest from the EC2 instance identity document.
type identityDocument struct {
	Tags        map[string]interface{} `json:"tags,omitempty" structs:"tags" mapstructure:"tags"`
	InstanceID  string                 `json:"instanceId,omitempty" structs:"instanceId" mapstructure:"instanceId"`
	AmiID       string                 `json:"imageId,omitempty" structs:"imageId" mapstructure:"imageId"`
	AccountID   string                 `json:"accountId,omitempty" structs:"accountId" mapstructure:"accountId"`
	Region      string                 `json:"region,omitempty" structs:"region" mapstructure:"region"`
	PendingTime string                 `json:"pendingTime,omitempty" structs:"pendingTime" mapstructure:"pendingTime"`
}

type roleTagLoginResponse struct {
	Policies                 []string      `json:"policies" structs:"policies" mapstructure:"policies"`
	MaxTTL                   time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	DisallowReauthentication bool          `json:"disallow_reauthentication" structs:"disallow_reauthentication" mapstructure:"disallow_reauthentication"`
}

const pathLoginSyn = `
Authenticates an EC2 instance with Vault.
`

const pathLoginDesc = `
An EC2 instance is authenticated using the PKCS#7 signature of the instance identity
document and a client created nonce. This nonce should be unique and should be used by
the instance for all future logins, unless 'disallow_reauthenitcation' option on the
registered role is enabled, in which case client nonce is optional.

First login attempt, creates a whitelist entry in Vault associating the instance to the nonce
provided. All future logins will succeed only if the client nonce matches the nonce in the
whitelisted entry.

By default, a cron task will periodically look for expired entries in the whitelist
and deletes them. The duration to periodically run this, is one hour by default.
However, this can be configured using the 'config/tidy/identities' endpoint. This tidy
action can be triggered via the API as well, using the 'tidy/identities' endpoint.
`
