package aws

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fullsailor/pkcs7"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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
				Description: "The nonce created by a client of this backend. Nonce is used to avoid replay attacks. When the instances are configured to be allowed to login only once, nonce parameter is of no use and hence can be skipped.",
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
func (b *backend) validateInstance(s logical.Storage, identityDoc *identityDocument) (*ec2.DescribeInstancesOutput, error) {
	// Create an EC2 client to pull the instance information
	ec2Client, err := b.clientEC2(s, identityDoc.Region, false)
	if err != nil {
		return nil, err
	}

	status, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(identityDoc.InstanceID),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching description for instance ID %s: %s\n", identityDoc.InstanceID, err)
	}
	if len(status.Reservations) == 0 {
		return nil, fmt.Errorf("no reservations found in instance description")

	}
	if len(status.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("no instance details found in reservations")
	}
	if *status.Reservations[0].Instances[0].InstanceId != identityDoc.InstanceID {
		return nil, fmt.Errorf("expected instance ID does not match the instance ID in the instance description")
	}
	if status.Reservations[0].Instances[0].State == nil {
		return nil, fmt.Errorf("instance state in instance description is nil")
	}
	if *status.Reservations[0].Instances[0].State.Code != 16 ||
		*status.Reservations[0].Instances[0].State.Name != "running" {
		return nil, fmt.Errorf("instance is not in 'running' state")
	}
	// Validate the instance through InstanceState, InstanceStatus and SystemStatus
	return status, nil
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
func (b *backend) parseIdentityDocument(s logical.Storage, pkcs7B64 string) (*identityDocument, error) {
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
	publicCerts, err := b.awsPublicCertificates(s)
	if err != nil {
		return nil, err
	}
	if publicCerts == nil || len(publicCerts) == 0 {
		return nil, fmt.Errorf("certificates to verify the signature are not found")
	}

	// Before calling Verify() on the PKCS#7 struct, set the certificate to be used
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
	identityDoc, err := b.parseIdentityDocument(req.Storage, pkcs7B64)
	if err != nil {
		return nil, err
	}
	if identityDoc == nil {
		return logical.ErrorResponse("failed to extract instance identity document from PKCS#7 signature"), nil
	}

	// Validate the instance ID.
	instanceDesc, err := b.validateInstance(req.Storage, identityDoc)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to verify instance ID: %s", err)), nil
	}

	// Get the entry for the AMI used by the instance.
	imageEntry, err := awsImage(req.Storage, identityDoc.AmiID)
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
			return logical.ErrorResponse(err.Error()), nil
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
		resp, err := b.handleRoleTagLogin(req.Storage, identityDoc, imageEntry, instanceDesc)
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
		// AmiID, ClientNonce and CreationTime of the identity entry,
		// once set, should never change.
		storedIdentity = &whitelistIdentity{
			AmiID:        identityDoc.AmiID,
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

// handleRoleTagLogin is used to fetch the role tag of the instance and verifies it to be correct.
// Then the policies for the login request will be set off of the role tag, if certain creteria satisfies.
func (b *backend) handleRoleTagLogin(s logical.Storage, identityDoc *identityDocument, imageEntry *awsImageEntry, instanceDesc *ec2.DescribeInstancesOutput) (*roleTagLoginResponse, error) {

	tags := instanceDesc.Reservations[0].Instances[0].Tags
	if tags == nil || len(tags) == 0 {
		return nil, fmt.Errorf("missing tag with key %s on the instance", imageEntry.RoleTag)
	}

	rTagValue := ""
	for _, tagItem := range tags {
		if tagItem.Key != nil && *tagItem.Key == imageEntry.RoleTag {
			rTagValue = *tagItem.Value
			break
		}
	}

	if rTagValue == "" {
		return nil, fmt.Errorf("missing tag with key %s on the instance", imageEntry.RoleTag)
	}

	// Parse the role tag into a struct, extract the plaintext part of it and verify its HMAC.
	rTag, err := parseRoleTagValue(s, rTagValue)
	if err != nil {
		return nil, err
	}

	// Check if the role tag belongs to the AMI ID of the instance.
	if rTag.AmiID != identityDoc.AmiID {
		return nil, fmt.Errorf("role tag does not belong to the instance's AMI ID.")
	}

	// If instance_id was set on the role tag, check if the same instance is attempting to login.
	if rTag.InstanceID != "" && rTag.InstanceID != identityDoc.InstanceID {
		return nil, fmt.Errorf("role tag is being used by an unauthorized instance.")
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

	imageEntry, err := awsImage(req.Storage, storedIdentity.AmiID)
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

// Struct to represent items of interest from the EC2 instance identity document.
type identityDocument struct {
	Tags        map[string]interface{} `json:"tags,omitempty" structs:"tags" mapstructure:"tags"`
	InstanceID  string                 `json:"instanceId,omitempty" structs:"instanceId" mapstructure:"instanceId"`
	AmiID       string                 `json:"imageId,omitempty" structs:"imageId" mapstructure:"imageId"`
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
