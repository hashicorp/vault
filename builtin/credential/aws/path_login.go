package aws

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/vishalnayak/pkcs7"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		Fields: map[string]*framework.FieldSchema{
			"pkcs7": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "PKCS7 signature of the identity document.",
			},

			"nonce": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The nonce created by a client of this backend.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLoginUpdate,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

// validateInstanceID queries the status of the EC2 instance using AWS EC2 API and
// checks if the instance is running and is healthy.
func (b *backend) validateInstanceID(s logical.Storage, instanceID string) error {
	// Create an EC2 client to pull the instance information
	ec2Client, err := b.clientEC2(s, false)
	if err != nil {
		return err
	}

	// Get the status of the instance
	instanceStatus, err := ec2Client.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		return err
	}

	// Validate the instance through InstanceState, InstanceStatus and SystemStatus
	return validateInstanceStatus(instanceStatus)
}

// validateMetadata matches the given client nonce and pending time with the one cached
// in the identity whitelist during the previous login. But, if reauthentication is
// disabled, login attempt is failed immediately.
func validateMetadata(clientNonce, pendingTime string, storedIdentity *whitelistIdentity, imageEntry *awsImageEntry) error {

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
	// ID is necessary, or the client must durably store the nonce.
	//
	// If the `allow_instance_migration` property of the registered AMI is
	// enabled, then the client nonce mismatch is ignored, as long as the
	// pending time in the presented instance identity document is newer than
	// the cached pending time. The new pendingTime is stored and used for
	// future checks.
	//
	// This is a weak criterion and hence the `allow_instance_migration` option
	// should be used with caution.
	if clientNonce != storedIdentity.ClientNonce {
		if !imageEntry.AllowInstanceMigration {
			return fmt.Errorf("client nonce mismatch")
		}
		if imageEntry.AllowInstanceMigration && !givenPendingTime.After(storedPendingTime) {
			return fmt.Errorf("client nonce mismatch and instance meta-data incorrect")
		}
	}

	// ensure that the 'pendingTime' on the given identity document is not before than the
	// 'pendingTime' that was used for previous login.
	if givenPendingTime.Before(storedPendingTime) {
		return fmt.Errorf("instance meta-data is older than the one used for previous login")
	}
	return nil
}

// Verifies the correctness of the authenticated attributes present in the PKCS#7
// signature. After verification, extracts the instance identity document from the
// signature, parses it and returns it.
func parseIdentityDocument(s logical.Storage, pkcs7B64 string) (*identityDocument, error) {
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

	// Get the public certificate that is used to verify the signature.
	publicCert, err := awsPublicCertificateParsed(s)
	if err != nil {
		return nil, err
	}
	if publicCert == nil {
		return nil, fmt.Errorf("certificate to verify the signature is not found")
	}

	// Before calling Verify() on the PKCS#7 struct, set the certificate to be used
	// to verify the contents in the signer information.
	pkcs7Data.Certificates = []*x509.Certificate{publicCert}

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
	err = json.Unmarshal(pkcs7Data.Content, &identityDoc)
	if err != nil {
		return nil, err
	}

	return &identityDoc, nil
}

// pathLoginUpdate is used to create a Vault token by the EC2 instances
// by providing its instance identity document, pkcs7 signature of the document,
// and a client created nonce.
func (b *backend) pathLoginUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	pkcs7B64 := data.Get("pkcs7").(string)

	if pkcs7B64 == "" {
		return logical.ErrorResponse("missing pkcs7"), nil
	}

	// Verify the signature of the identity document.
	identityDoc, err := parseIdentityDocument(req.Storage, pkcs7B64)
	if err != nil {
		return nil, err
	}
	if identityDoc == nil {
		return logical.ErrorResponse("failed to extract instance identity document from PKCS#7 signature"), nil
	}

	// Validate the instance ID.
	if err := b.validateInstanceID(req.Storage, identityDoc.InstanceID); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to verify instance ID: %s", err)), nil
	}

	// Get the entry for the AMI used by the instance.
	imageEntry, err := awsImage(req.Storage, identityDoc.ImageID)
	if err != nil {
		return nil, err
	}
	if imageEntry == nil {
		return logical.ErrorResponse("image entry not found"), nil
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
		// with which previous login was made.
		if err = validateMetadata(clientNonce, identityDoc.PendingTime, storedIdentity, imageEntry); err != nil {
			return nil, err
		}
	}

	// Load the current values for max TTL and policies from the image entry,
	// before checking for overriding by the RoleTag
	maxTTL := b.System().MaxLeaseTTL()
	if imageEntry.MaxTTL > time.Duration(0) && imageEntry.MaxTTL < maxTTL {
		maxTTL = imageEntry.MaxTTL
	}

	policies := imageEntry.Policies
	rTagMaxTTL := time.Duration(0)
	disallowReauthentication := imageEntry.DisallowReauthentication

	// Role tag is enabled for the AMI.
	if imageEntry.RoleTag != "" {
		// Overwrite the policies with the ones returned from processing the role tag.
		resp, err := b.handleRoleTagLogin(req.Storage, identityDoc, imageEntry)
		if err != nil {
			return nil, err
		}
		policies = resp.Policies
		rTagMaxTTL = resp.MaxTTL

		// If imageEntry had disallowReauthentication set to 'true', do not reset it
		// to 'false' based on role tag having it not set. But, if role tag had it set,
		// be sure to override the value.
		if !disallowReauthentication {
			disallowReauthentication = resp.DisallowReauthentication
		}

		if resp.MaxTTL > time.Duration(0) && resp.MaxTTL < maxTTL {
			maxTTL = resp.MaxTTL
		}
	}

	// Save the login attempt in the identity whitelist.
	currentTime := time.Now()
	if storedIdentity == nil {
		// ImageID, ClientNonce and CreationTime of the identity entry,
		// once set, should never change.
		storedIdentity = &whitelistIdentity{
			ImageID:      identityDoc.ImageID,
			ClientNonce:  clientNonce,
			CreationTime: currentTime,
		}
	}

	// DisallowReauthentication, PendingTime, LastUpdatedTime and ExpirationTime may change.
	storedIdentity.LastUpdatedTime = currentTime
	storedIdentity.ExpirationTime = currentTime.Add(maxTTL)
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
				"role_tag_max_ttl": rTagMaxTTL.String(),
			},
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       b.System().DefaultLeaseTTL(),
			},
		},
	}

	// Enforce our image/role tag maximum TTL
	if maxTTL < resp.Auth.TTL {
		resp.Auth.TTL = maxTTL
	}

	return resp, nil

}

// fetchRoleTagValue creates an AWS EC2 client and queries the tags
// attached to the instance identified by the given instanceID.
func (b *backend) fetchRoleTagValue(s logical.Storage, tagKey string) (string, error) {
	ec2Client, err := b.clientEC2(s, false)
	if err != nil {
		return "", err
	}

	// Retrieve the instance tag with a "key" filter matching tagKey.
	tagsOutput, err := ec2Client.DescribeTags(&ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("key"),
				Values: []*string{
					aws.String(tagKey),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	if tagsOutput.Tags == nil ||
		len(tagsOutput.Tags) != 1 ||
		*tagsOutput.Tags[0].Key != tagKey ||
		*tagsOutput.Tags[0].ResourceType != "instance" {
		return "", nil
	}

	return *tagsOutput.Tags[0].Value, nil
}

// handleRoleTagLogin is used to fetch the role tag of the instance and verifies it to be correct.
// Then the policies for the login request will be set off of the role tag, if certain creteria satisfies.
func (b *backend) handleRoleTagLogin(s logical.Storage, identityDoc *identityDocument, imageEntry *awsImageEntry) (*roleTagLoginResponse, error) {

	// Make a secondary call to the AWS instance to see if the desired tag is set.
	// NOTE: If AWS adds the instance tags as meta-data in the instance identity
	// document, then it is better to look this information there instead of making
	// another API call. Currently, we don't have an option but make this call.
	rTagValue, err := b.fetchRoleTagValue(s, imageEntry.RoleTag)
	if err != nil {
		return nil, err
	}

	if rTagValue == "" {
		return nil, fmt.Errorf("missing tag with key %s on the instance", imageEntry.RoleTag)
	}

	// Parse the role tag into a struct, extract the plaintext part of it and verify its HMAC.
	rTag, err := parseRoleTagValue(s, rTagValue)
	if err != nil {
		return nil, err
	}

	// Check if the role tag is blacklisted.
	blacklistEntry, err := blacklistRoleTagEntry(s, rTagValue)
	if err != nil {
		return nil, err
	}
	if blacklistEntry != nil {
		return nil, fmt.Errorf("role tag is blacklisted")
	}

	// Ensure that the policies on the RoleTag is a subset of policies on the image
	if !strutil.StrListSubset(imageEntry.Policies, rTag.Policies) {
		return nil, fmt.Errorf("policies on the role tag must be subset of policies on the image")
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

	storedIdentity, err := whitelistIdentityEntry(req.Storage, req.Auth.Metadata["instance_id"])
	if err != nil {
		return nil, err
	}

	// For now, rTagMaxTTL is cached in internal data during login and used in renewal for
	// setting the MaxTTL for the stored login identity entry.
	// If `instance_id` can be used to fetch the role tag again (through an API), it would be good.
	// For accessing the max_ttl, storing the entire identity document is too heavy.
	rTagMaxTTL, err := time.ParseDuration(req.Auth.Metadata["role_tag_max_ttl"])
	if err != nil {
		return nil, err
	}

	imageEntry, err := awsImage(req.Storage, storedIdentity.ImageID)
	if err != nil {
		return nil, err
	}
	if imageEntry == nil {
		return logical.ErrorResponse("image entry not found"), nil
	}

	maxTTL := b.System().MaxLeaseTTL()
	if imageEntry.MaxTTL > time.Duration(0) && imageEntry.MaxTTL < maxTTL {
		maxTTL = imageEntry.MaxTTL
	}
	if rTagMaxTTL > time.Duration(0) && maxTTL > rTagMaxTTL {
		maxTTL = rTagMaxTTL
	}

	// Only LastUpdatedTime and ExpirationTime change, none else.
	currentTime := time.Now()
	storedIdentity.LastUpdatedTime = currentTime
	storedIdentity.ExpirationTime = currentTime.Add(maxTTL)

	if err = setWhitelistIdentityEntry(req.Storage, req.Auth.Metadata["instance_id"], storedIdentity); err != nil {
		return nil, err
	}

	return framework.LeaseExtend(req.Auth.TTL, maxTTL, b.System())(req, data)
}

// Validates the instance by checking the InstanceState, InstanceStatus and SystemStatus
func validateInstanceStatus(instanceStatus *ec2.DescribeInstanceStatusOutput) error {

	if instanceStatus.InstanceStatuses == nil {
		return fmt.Errorf("instance statuses not found")
	}

	if len(instanceStatus.InstanceStatuses) != 1 {
		return fmt.Errorf("length of instance statuses is more than 1")
	}

	if instanceStatus.InstanceStatuses[0].InstanceState == nil {
		return fmt.Errorf("instance state not found")
	}

	// Instance should be in 'running'(code 16) state.
	if *instanceStatus.InstanceStatuses[0].InstanceState.Code != 16 {
		return fmt.Errorf("instance state is not 'running'")
	}

	if instanceStatus.InstanceStatuses[0].InstanceStatus == nil {
		return fmt.Errorf("instance status not found")
	}

	// InstanceStatus should be 'ok'
	if *instanceStatus.InstanceStatuses[0].InstanceStatus.Status != "ok" {
		return fmt.Errorf("instance status is not 'ok'")
	}

	if instanceStatus.InstanceStatuses[0].SystemStatus == nil {
		return fmt.Errorf("system status not found")
	}

	// SystemStatus should be 'ok'
	if *instanceStatus.InstanceStatuses[0].SystemStatus.Status != "ok" {
		return fmt.Errorf("system status is not 'ok'")
	}

	return nil
}

// Struct to represent items of interest from the EC2 instance identity document.
type identityDocument struct {
	Tags        map[string]interface{} `json:"tags,omitempty" structs:"tags" mapstructure:"tags"`
	InstanceID  string                 `json:"instanceId,omitempty" structs:"instanceId" mapstructure:"instanceId"`
	ImageID     string                 `json:"imageId,omitempty" structs:"imageId" mapstructure:"imageId"`
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
An EC2 instance is authenticated using the instance identity document, the identity document's
PKCS#7 signature and a client created nonce. This nonce should be unique and should be used by
the instance for all future logins.

First login attempt, creates a whitelist entry in Vault associating the instance to the nonce
provided. All future logins will succeed only if the client nonce matches the nonce in the
whitelisted entry.

The entries in the whitelist are not automatically deleted. Although, they will have an
expiration time set on the entry. There is a separate endpoint 'whitelist/identity/tidy',
that needs to be invoked to clean-up all the expired entries in the whitelist.
`
