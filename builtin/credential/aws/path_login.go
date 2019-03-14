package awsauth

import (
	"context"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/fullsailor/pkcs7"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/awsutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	reauthenticationDisabledNonce = "reauthentication-disabled-nonce"
	iamAuthType                   = "iam"
	ec2AuthType                   = "ec2"
	ec2EntityType                 = "ec2_instance"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type: framework.TypeString,
				Description: `Name of the role against which the login is being attempted.
If 'role' is not specified, then the login endpoint looks for a role
bearing the name of the AMI ID of the EC2 instance that is trying to login.
If a matching role is not found, login fails.`,
			},

			"pkcs7": {
				Type: framework.TypeString,
				Description: `PKCS7 signature of the identity document when using an auth_type
of ec2.`,
			},

			"nonce": {
				Type: framework.TypeString,
				Description: `The nonce to be used for subsequent login requests when
auth_type is ec2.  If this parameter is not specified at
all and if reauthentication is allowed, then the backend will generate a random
nonce, attaches it to the instance's identity-whitelist entry and returns the
nonce back as part of auth metadata.  This value should be used with further
login requests, to establish client authenticity. Clients can choose to set a
custom nonce if preferred, in which case, it is recommended that clients provide
a strong nonce.  If a nonce is provided but with an empty value, it indicates
intent to disable reauthentication. Note that, when 'disallow_reauthentication'
option is enabled on either the role or the role tag, the 'nonce' holds no
significance.`,
			},

			"iam_http_request_method": {
				Type: framework.TypeString,
				Description: `HTTP method to use for the AWS request when auth_type is
iam. This must match what has been signed in the
presigned request. Currently, POST is the only supported value`,
			},

			"iam_request_url": {
				Type: framework.TypeString,
				Description: `Base64-encoded full URL against which to make the AWS request
when using iam auth_type.`,
			},

			"iam_request_body": {
				Type: framework.TypeString,
				Description: `Base64-encoded request body when auth_type is iam.
This must match the request body included in the signature.`,
			},
			"iam_request_headers": {
				Type: framework.TypeHeader,
				Description: `Key/value pairs of headers for use in the
sts:GetCallerIdentity HTTP requests headers when auth_type is iam. Can be either 
a Base64-encoded, JSON-serialized string, or a JSON object of key/value pairs. 
This must at a minimum include the headers over which AWS has included a  signature.`,
			},
			"identity": {
				Type: framework.TypeString,
				Description: `Base64 encoded EC2 instance identity document. This needs to be supplied along
with the 'signature' parameter. If using 'curl' for fetching the identity
document, consider using the option '-w 0' while piping the output to 'base64'
binary.`,
			},
			"signature": {
				Type: framework.TypeString,
				Description: `Base64 encoded SHA256 RSA signature of the instance identity document. This
needs to be supplied along with 'identity' parameter.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLoginUpdate,
			logical.AliasLookaheadOperation: b.pathLoginUpdate,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

// instanceIamRoleARN fetches the IAM role ARN associated with the given
// instance profile name
func (b *backend) instanceIamRoleARN(iamClient *iam.IAM, instanceProfileName string) (string, error) {
	if iamClient == nil {
		return "", fmt.Errorf("nil iamClient")
	}
	if instanceProfileName == "" {
		return "", fmt.Errorf("missing instance profile name")
	}

	profile, err := iamClient.GetInstanceProfile(&iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(instanceProfileName),
	})
	if err != nil {
		return "", awsutil.AppendLogicalError(err)
	}
	if profile == nil {
		return "", fmt.Errorf("nil output while getting instance profile details")
	}

	if profile.InstanceProfile == nil {
		return "", fmt.Errorf("nil instance profile in the output of instance profile details")
	}

	if profile.InstanceProfile.Roles == nil || len(profile.InstanceProfile.Roles) != 1 {
		return "", fmt.Errorf("invalid roles in the output of instance profile details")
	}

	if profile.InstanceProfile.Roles[0].Arn == nil {
		return "", fmt.Errorf("nil role ARN in the output of instance profile details")
	}

	return *profile.InstanceProfile.Roles[0].Arn, nil
}

// validateInstance queries the status of the EC2 instance using AWS EC2 API
// and checks if the instance is running and is healthy
func (b *backend) validateInstance(ctx context.Context, s logical.Storage, instanceID, region, accountID string) (*ec2.Instance, error) {
	// Create an EC2 client to pull the instance information
	ec2Client, err := b.clientEC2(ctx, s, region, accountID)
	if err != nil {
		return nil, err
	}

	status, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	})
	if err != nil {
		errW := errwrap.Wrapf(fmt.Sprintf("error fetching description for instance ID %q: {{err}}", instanceID), err)
		return nil, errwrap.Wrap(errW, awsutil.CheckAWSError(err))
	}
	if status == nil {
		return nil, fmt.Errorf("nil output from describe instances")
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
	if *status.Reservations[0].Instances[0].State.Name != "running" {
		return nil, fmt.Errorf("instance is not in 'running' state")
	}
	return status.Reservations[0].Instances[0], nil
}

// validateMetadata matches the given client nonce and pending time with the
// one cached in the identity whitelist during the previous login. But, if
// reauthentication is disabled, login attempt is failed immediately.
func validateMetadata(clientNonce, pendingTime string, storedIdentity *whitelistIdentity, roleEntry *awsRoleEntry) error {
	// For sanity
	if !storedIdentity.DisallowReauthentication && storedIdentity.ClientNonce == "" {
		return fmt.Errorf("client nonce missing in stored identity")
	}

	// If reauthentication is disabled or if the nonce supplied matches a
	// predefined nonce which indicates reauthentication to be disabled,
	// authentication will not succeed.
	if storedIdentity.DisallowReauthentication ||
		subtle.ConstantTimeCompare([]byte(reauthenticationDisabledNonce), []byte(clientNonce)) == 1 {
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

	// When the presented client nonce does not match the cached entry, it
	// is either that a rogue client is trying to login or that a valid
	// client suffered a migration. The migration is detected via
	// pendingTime in the instance metadata, which sadly is only updated
	// when an instance is stopped and started but *not* when the instance
	// is rebooted. If reboot survivability is needed, either
	// instrumentation to delete the instance ID from the whitelist is
	// necessary, or the client must durably store the nonce.
	//
	// If the `allow_instance_migration` property of the registered role is
	// enabled, then the client nonce mismatch is ignored, as long as the
	// pending time in the presented instance identity document is newer
	// than the cached pending time. The new pendingTime is stored and used
	// for future checks.
	//
	// This is a weak criterion and hence the `allow_instance_migration`
	// option should be used with caution.
	if subtle.ConstantTimeCompare([]byte(clientNonce), []byte(storedIdentity.ClientNonce)) != 1 {
		if !roleEntry.AllowInstanceMigration {
			return fmt.Errorf("client nonce mismatch")
		}
		if roleEntry.AllowInstanceMigration && !givenPendingTime.After(storedPendingTime) {
			return fmt.Errorf("client nonce mismatch and instance meta-data incorrect")
		}
	}

	// Ensure that the 'pendingTime' on the given identity document is not
	// before the 'pendingTime' that was used for previous login. This
	// disallows old metadata documents from being used to perform login.
	if givenPendingTime.Before(storedPendingTime) {
		return fmt.Errorf("instance meta-data is older than the one used for previous login")
	}
	return nil
}

// Verifies the integrity of the instance identity document using its SHA256
// RSA signature. After verification, returns the unmarshaled instance identity
// document.
func (b *backend) verifyInstanceIdentitySignature(ctx context.Context, s logical.Storage, identityBytes, signatureBytes []byte) (*identityDocument, error) {
	if len(identityBytes) == 0 {
		return nil, fmt.Errorf("missing instance identity document")
	}

	if len(signatureBytes) == 0 {
		return nil, fmt.Errorf("missing SHA256 RSA signature of the instance identity document")
	}

	// Get the public certificates that are used to verify the signature.
	// This returns a slice of certificates containing the default
	// certificate and all the registered certificates via
	// 'config/certificate/<cert_name>' endpoint, for verifying the RSA
	// digest.
	publicCerts, err := b.awsPublicCertificates(ctx, s, false)
	if err != nil {
		return nil, err
	}
	if publicCerts == nil || len(publicCerts) == 0 {
		return nil, fmt.Errorf("certificates to verify the signature are not found")
	}

	// Check if any of the certs registered at the backend can verify the
	// signature
	for _, cert := range publicCerts {
		err := cert.CheckSignature(x509.SHA256WithRSA, identityBytes, signatureBytes)
		if err == nil {
			var identityDoc identityDocument
			if decErr := jsonutil.DecodeJSON(identityBytes, &identityDoc); decErr != nil {
				return nil, decErr
			}
			return &identityDoc, nil
		}
	}

	return nil, fmt.Errorf("instance identity verification using SHA256 RSA signature is unsuccessful")
}

// Verifies the correctness of the authenticated attributes present in the PKCS#7
// signature. After verification, extracts the instance identity document from the
// signature, parses it and returns it.
func (b *backend) parseIdentityDocument(ctx context.Context, s logical.Storage, pkcs7B64 string) (*identityDocument, error) {
	// Insert the header and footer for the signature to be able to pem decode it
	pkcs7B64 = fmt.Sprintf("-----BEGIN PKCS7-----\n%s\n-----END PKCS7-----", pkcs7B64)

	// Decode the PEM encoded signature
	pkcs7BER, pkcs7Rest := pem.Decode([]byte(pkcs7B64))
	if len(pkcs7Rest) != 0 {
		return nil, fmt.Errorf("failed to decode the PEM encoded PKCS#7 signature")
	}

	// Parse the signature from asn1 format into a struct
	pkcs7Data, err := pkcs7.Parse(pkcs7BER.Bytes)
	if err != nil {
		return nil, errwrap.Wrapf("failed to parse the BER encoded PKCS#7 signature: {{err}}", err)
	}

	// Get the public certificates that are used to verify the signature.
	// This returns a slice of certificates containing the default certificate
	// and all the registered certificates via 'config/certificate/<cert_name>' endpoint
	publicCerts, err := b.awsPublicCertificates(ctx, s, true)
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

	// Check if the signature has content inside of it
	if len(pkcs7Data.Content) == 0 {
		return nil, fmt.Errorf("instance identity document could not be found in the signature")
	}

	var identityDoc identityDocument
	if err := jsonutil.DecodeJSON(pkcs7Data.Content, &identityDoc); err != nil {
		return nil, err
	}

	return &identityDoc, nil
}

func (b *backend) pathLoginUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	anyEc2, allEc2 := hasValuesForEc2Auth(data)
	anyIam, allIam := hasValuesForIamAuth(data)
	switch {
	case anyEc2 && anyIam:
		return logical.ErrorResponse("supplied auth values for both ec2 and iam auth types"), nil
	case anyEc2 && !allEc2:
		return logical.ErrorResponse("supplied some of the auth values for the ec2 auth type but not all"), nil
	case anyEc2:
		return b.pathLoginUpdateEc2(ctx, req, data)
	case anyIam && !allIam:
		return logical.ErrorResponse("supplied some of the auth values for the iam auth type but not all"), nil
	case anyIam:
		return b.pathLoginUpdateIam(ctx, req, data)
	default:
		return logical.ErrorResponse("didn't supply required authentication values"), nil
	}
}

// Returns whether the EC2 instance meets the requirements of the particular
// AWS role entry.
// The first error return value is whether there's some sort of validation
// error that means the instance doesn't meet the role requirements
// The second error return value indicates whether there's an error in even
// trying to validate those requirements
func (b *backend) verifyInstanceMeetsRoleRequirements(ctx context.Context,
	s logical.Storage, instance *ec2.Instance, roleEntry *awsRoleEntry, roleName string, identityDoc *identityDocument) (error, error) {

	switch {
	case instance == nil:
		return nil, fmt.Errorf("nil instance")
	case roleEntry == nil:
		return nil, fmt.Errorf("nil roleEntry")
	case identityDoc == nil:
		return nil, fmt.Errorf("nil identityDoc")
	}

	// Verify that the instance ID matches one of the ones set by the role
	if len(roleEntry.BoundEc2InstanceIDs) > 0 && !strutil.StrListContains(roleEntry.BoundEc2InstanceIDs, *instance.InstanceId) {
		return fmt.Errorf("instance ID %q does not belong to the role %q", *instance.InstanceId, roleName), nil
	}

	// Verify that the AccountID of the instance trying to login matches the
	// AccountID specified as a constraint on role
	if len(roleEntry.BoundAccountIDs) > 0 && !strutil.StrListContains(roleEntry.BoundAccountIDs, identityDoc.AccountID) {
		return fmt.Errorf("account ID %q does not belong to role %q", identityDoc.AccountID, roleName), nil
	}

	// Verify that the AMI ID of the instance trying to login matches the
	// AMI ID specified as a constraint on the role.
	//
	// Here, we're making a tradeoff and pulling the AMI ID out of the EC2
	// API rather than the signed instance identity doc. They *should* match.
	// This means we require an EC2 API call to retrieve the AMI ID, but we're
	// already calling the API to validate the Instance ID anyway, so it shouldn't
	// matter. The benefit is that we have the exact same code whether auth_type
	// is ec2 or iam.
	if len(roleEntry.BoundAmiIDs) > 0 {
		if instance.ImageId == nil {
			return nil, fmt.Errorf("AMI ID in the instance description is nil")
		}
		if !strutil.StrListContains(roleEntry.BoundAmiIDs, *instance.ImageId) {
			return fmt.Errorf("AMI ID %q does not belong to role %q", *instance.ImageId, roleName), nil
		}
	}

	// Validate the SubnetID if corresponding bound was set on the role
	if len(roleEntry.BoundSubnetIDs) > 0 {
		if instance.SubnetId == nil {
			return nil, fmt.Errorf("subnet ID in the instance description is nil")
		}
		if !strutil.StrListContains(roleEntry.BoundSubnetIDs, *instance.SubnetId) {
			return fmt.Errorf("subnet ID %q does not satisfy the constraint on role %q", *instance.SubnetId, roleName), nil
		}
	}

	// Validate the VpcID if corresponding bound was set on the role
	if len(roleEntry.BoundVpcIDs) > 0 {
		if instance.VpcId == nil {
			return nil, fmt.Errorf("VPC ID in the instance description is nil")
		}
		if !strutil.StrListContains(roleEntry.BoundVpcIDs, *instance.VpcId) {
			return fmt.Errorf("VPC ID %q does not satisfy the constraint on role %q", *instance.VpcId, roleName), nil
		}
	}

	// Check if the IAM instance profile ARN of the instance trying to
	// login, matches the IAM instance profile ARN specified as a constraint
	// on the role
	if len(roleEntry.BoundIamInstanceProfileARNs) > 0 {
		if instance.IamInstanceProfile == nil {
			return nil, fmt.Errorf("IAM instance profile in the instance description is nil")
		}
		if instance.IamInstanceProfile.Arn == nil {
			return nil, fmt.Errorf("IAM instance profile ARN in the instance description is nil")
		}
		iamInstanceProfileARN := *instance.IamInstanceProfile.Arn
		matchesInstanceProfile := false
		// NOTE: Can't use strutil.StrListContainsGlob. A * is a perfectly valid character in the "path" component
		// of an ARN. See, e.g., https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateInstanceProfile.html :
		// The path allows strings "containing any ASCII character from the ! (\u0021) thru the DEL character
		// (\u007F), including most punctuation characters, digits, and upper and lowercased letters."
		// So, e.g., arn:aws:iam::123456789012:instance-profile/Some*Path/MyProfileName is a perfectly valid instance
		// profile ARN, and it wouldn't be correct to expand the * in the middle as a wildcard.
		// If a user wants to match an IAM instance profile arn beginning with arn:aws:iam::123456789012:instance-profile/foo*
		// then bound_iam_instance_profile_arn would need to be arn:aws:iam::123456789012:instance-profile/foo**
		// Wanting to exactly match an ARN that has a * at the end is not a valid use case. The * is only valid in the
		// path; it's not valid in the name. That means no valid ARN can ever end with a *. For example,
		// arn:aws:iam::123456789012:instance-profile/Foo* is NOT valid as an instance profile ARN, so no valid instance
		// profile ARN could ever equal that value.
		for _, boundInstanceProfileARN := range roleEntry.BoundIamInstanceProfileARNs {
			switch {
			case strings.HasSuffix(boundInstanceProfileARN, "*") && strings.HasPrefix(iamInstanceProfileARN, boundInstanceProfileARN[:len(boundInstanceProfileARN)-1]):
				matchesInstanceProfile = true
				break
			case iamInstanceProfileARN == boundInstanceProfileARN:
				matchesInstanceProfile = true
				break
			}
		}
		if !matchesInstanceProfile {
			return fmt.Errorf("IAM instance profile ARN %q does not satisfy the constraint role %q", iamInstanceProfileARN, roleName), nil
		}
	}

	// Check if the IAM role ARN of the instance trying to login, matches
	// the IAM role ARN specified as a constraint on the role.
	if len(roleEntry.BoundIamRoleARNs) > 0 {
		if instance.IamInstanceProfile == nil {
			return nil, fmt.Errorf("IAM instance profile in the instance description is nil")
		}
		if instance.IamInstanceProfile.Arn == nil {
			return nil, fmt.Errorf("IAM instance profile ARN in the instance description is nil")
		}

		// Fetch the instance profile ARN from the instance description
		iamInstanceProfileARN := *instance.IamInstanceProfile.Arn

		if iamInstanceProfileARN == "" {
			return nil, fmt.Errorf("IAM instance profile ARN in the instance description is empty")
		}

		// Extract out the instance profile name from the instance
		// profile ARN
		iamInstanceProfileEntity, err := parseIamArn(iamInstanceProfileARN)

		if err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("failed to parse IAM instance profile ARN %q: {{err}}", iamInstanceProfileARN), err)
		}

		// Use instance profile ARN to fetch the associated role ARN
		iamClient, err := b.clientIAM(ctx, s, identityDoc.Region, identityDoc.AccountID)
		if err != nil {
			return nil, errwrap.Wrapf("could not fetch IAM client: {{err}}", err)
		} else if iamClient == nil {
			return nil, fmt.Errorf("received a nil iamClient")
		}
		iamRoleARN, err := b.instanceIamRoleARN(iamClient, iamInstanceProfileEntity.FriendlyName)
		if err != nil {
			return nil, errwrap.Wrapf("IAM role ARN could not be fetched: {{err}}", err)
		}
		if iamRoleARN == "" {
			return nil, fmt.Errorf("IAM role ARN could not be fetched")
		}

		matchesInstanceRoleARN := false
		for _, boundIamRoleARN := range roleEntry.BoundIamRoleARNs {
			switch {
			// as with boundInstanceProfileARN, can't use strutil.StrListContainsGlob because * can validly exist in the middle of an ARN
			case strings.HasSuffix(boundIamRoleARN, "*") && strings.HasPrefix(iamRoleARN, boundIamRoleARN[:len(boundIamRoleARN)-1]):
				matchesInstanceRoleARN = true
				break
			case iamRoleARN == boundIamRoleARN:
				matchesInstanceRoleARN = true
				break
			}
		}
		if !matchesInstanceRoleARN {
			return fmt.Errorf("IAM role ARN %q does not satisfy the constraint role %q", iamRoleARN, roleName), nil
		}
	}

	return nil, nil
}

// pathLoginUpdateEc2 is used to create a Vault token by the EC2 instances
// by providing the pkcs7 signature of the instance identity document
// and a client created nonce. Client nonce is optional if 'disallow_reauthentication'
// option is enabled on the registered role.
func (b *backend) pathLoginUpdateEc2(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	identityDocB64 := data.Get("identity").(string)
	var identityDocBytes []byte
	var err error
	if identityDocB64 != "" {
		identityDocBytes, err = base64.StdEncoding.DecodeString(identityDocB64)
		if err != nil || len(identityDocBytes) == 0 {
			return logical.ErrorResponse("failed to base64 decode the instance identity document"), nil
		}
	}

	signatureB64 := data.Get("signature").(string)
	var signatureBytes []byte
	if signatureB64 != "" {
		signatureBytes, err = base64.StdEncoding.DecodeString(signatureB64)
		if err != nil {
			return logical.ErrorResponse("failed to base64 decode the SHA256 RSA signature of the instance identity document"), nil
		}
	}

	pkcs7B64 := data.Get("pkcs7").(string)

	// Either the pkcs7 signature of the instance identity document, or
	// the identity document itself along with its SHA256 RSA signature
	// needs to be provided.
	if pkcs7B64 == "" && (len(identityDocBytes) == 0 && len(signatureBytes) == 0) {
		return logical.ErrorResponse("either pkcs7 or a tuple containing the instance identity document and its SHA256 RSA signature needs to be provided"), nil
	} else if pkcs7B64 != "" && (len(identityDocBytes) != 0 && len(signatureBytes) != 0) {
		return logical.ErrorResponse("both pkcs7 and a tuple containing the instance identity document and its SHA256 RSA signature is supplied; provide only one"), nil
	}

	// Verify the signature of the identity document and unmarshal it
	var identityDocParsed *identityDocument
	if pkcs7B64 != "" {
		identityDocParsed, err = b.parseIdentityDocument(ctx, req.Storage, pkcs7B64)
		if err != nil {
			return nil, err
		}
		if identityDocParsed == nil {
			return logical.ErrorResponse("failed to verify the instance identity document using pkcs7"), nil
		}
	} else {
		identityDocParsed, err = b.verifyInstanceIdentitySignature(ctx, req.Storage, identityDocBytes, signatureBytes)
		if err != nil {
			return nil, err
		}
		if identityDocParsed == nil {
			return logical.ErrorResponse("failed to verify the instance identity document using the SHA256 RSA digest"), nil
		}
	}

	roleName := data.Get("role").(string)

	// If roleName is not supplied, a role in the name of the instance's AMI ID will be looked for
	if roleName == "" {
		roleName = identityDocParsed.AmiID
	}

	// Get the entry for the role used by the instance
	roleEntry, err := b.lockedAWSRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("entry for role %q not found", roleName)), nil
	}

	if roleEntry.AuthType != ec2AuthType {
		return logical.ErrorResponse(fmt.Sprintf("auth method ec2 not allowed for role %s", roleName)), nil
	}

	identityConfigEntry, err := identityConfigEntry(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	identityAlias := ""

	switch identityConfigEntry.EC2Alias {
	case identityAliasRoleID:
		identityAlias = roleEntry.RoleID
	case identityAliasEC2InstanceID:
		identityAlias = identityDocParsed.InstanceID
	case identityAliasEC2ImageID:
		identityAlias = identityDocParsed.AmiID
	}

	// If we're just looking up for MFA, return the Alias info
	if req.Operation == logical.AliasLookaheadOperation {
		return &logical.Response{
			Auth: &logical.Auth{
				Alias: &logical.Alias{
					Name: identityAlias,
				},
			},
		}, nil
	}

	// Validate the instance ID by making a call to AWS EC2 DescribeInstances API
	// and fetching the instance description. Validation succeeds only if the
	// instance is in 'running' state.
	instance, err := b.validateInstance(ctx, req.Storage, identityDocParsed.InstanceID, identityDocParsed.Region, identityDocParsed.AccountID)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to verify instance ID: %v", err)), nil
	}

	// Verify that the `Region` of the instance trying to login matches the
	// `Region` specified as a constraint on role
	if len(roleEntry.BoundRegions) > 0 && !strutil.StrListContains(roleEntry.BoundRegions, identityDocParsed.Region) {
		return logical.ErrorResponse(fmt.Sprintf("Region %q does not satisfy the constraint on role %q", identityDocParsed.Region, roleName)), nil
	}

	validationError, err := b.verifyInstanceMeetsRoleRequirements(ctx, req.Storage, instance, roleEntry, roleName, identityDocParsed)
	if err != nil {
		return nil, err
	}
	if validationError != nil {
		return logical.ErrorResponse(fmt.Sprintf("Error validating instance: %v", validationError)), nil
	}

	// Get the entry from the identity whitelist, if there is one
	storedIdentity, err := whitelistIdentityEntry(ctx, req.Storage, identityDocParsed.InstanceID)
	if err != nil {
		return nil, err
	}

	// disallowReauthentication value that gets cached at the stored
	// identity-whitelist entry is determined not just by the role entry.
	// If client explicitly sets nonce to be empty, it implies intent to
	// disable reauthentication. Also, role tag can override the 'false'
	// value with 'true' (the other way around is not allowed).

	// Read the value from the role entry
	disallowReauthentication := roleEntry.DisallowReauthentication

	clientNonce := ""

	// Check if the nonce is supplied by the client
	clientNonceRaw, clientNonceSupplied := data.GetOk("nonce")
	if clientNonceSupplied {
		clientNonce = clientNonceRaw.(string)

		// Nonce explicitly set to empty implies intent to disable
		// reauthentication by the client. Set a predefined nonce which
		// indicates reauthentication being disabled.
		if clientNonce == "" {
			clientNonce = reauthenticationDisabledNonce

			// Ensure that the intent lands in the whitelist
			disallowReauthentication = true
		}
	}

	// This is NOT a first login attempt from the client
	if storedIdentity != nil {
		// Check if the client nonce match the cached nonce and if the pending time
		// of the identity document is not before the pending time of the document
		// with which previous login was made. If 'allow_instance_migration' is
		// enabled on the registered role, client nonce requirement is relaxed.
		if err = validateMetadata(clientNonce, identityDocParsed.PendingTime, storedIdentity, roleEntry); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		// Don't let subsequent login attempts to bypass the initial
		// intent of disabling reauthentication, despite the properties
		// of role getting updated. For example: Role has the value set
		// to 'false', a role-tag login sets the value to 'true', then
		// role gets updated to not use a role-tag, and a login attempt
		// is made with role's value set to 'false'. Removing the entry
		// from the identity-whitelist should be the only way to be
		// able to login from the instance again.
		disallowReauthentication = disallowReauthentication || storedIdentity.DisallowReauthentication
	}

	// If we reach this point without erroring and if the client nonce was
	// not supplied, a first time login is implied and that the client
	// intends that the nonce be generated by the backend. Create a random
	// nonce to be associated for the instance ID.
	if !clientNonceSupplied {
		if clientNonce, err = uuid.GenerateUUID(); err != nil {
			return nil, fmt.Errorf("failed to generate random nonce")
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
	var roleTagResp *roleTagLoginResponse
	if roleEntry.RoleTag != "" {
		roleTagResp, err = b.handleRoleTagLogin(ctx, req.Storage, roleName, roleEntry, instance)
		if err != nil {
			return nil, err
		}
		if roleTagResp == nil {
			return logical.ErrorResponse("failed to fetch and verify the role tag"), nil
		}
	}

	if roleTagResp != nil {
		// Role tag is enabled on the role.

		// Overwrite the policies with the ones returned from processing the role tag
		// If there are no policies on the role tag, policies on the role are inherited.
		// If policies on role tag are set, by this point, it is verified that it is a subset of the
		// policies on the role. So, apply only those.
		if len(roleTagResp.Policies) != 0 {
			policies = roleTagResp.Policies
		}

		// If roleEntry had disallowReauthentication set to 'true', do not reset it
		// to 'false' based on role tag having it not set. But, if role tag had it set,
		// be sure to override the value.
		if !disallowReauthentication {
			disallowReauthentication = roleTagResp.DisallowReauthentication
		}

		// Cache the value of role tag's max_ttl value
		rTagMaxTTL = roleTagResp.MaxTTL

		// Scope the shortestMaxTTL to the value set on the role tag
		if roleTagResp.MaxTTL > time.Duration(0) && roleTagResp.MaxTTL < shortestMaxTTL {
			shortestMaxTTL = roleTagResp.MaxTTL
		}
		if roleTagResp.MaxTTL > longestMaxTTL {
			longestMaxTTL = roleTagResp.MaxTTL
		}
	}

	// Save the login attempt in the identity whitelist
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

	// DisallowReauthentication, PendingTime, LastUpdatedTime and
	// ExpirationTime may change.
	storedIdentity.LastUpdatedTime = currentTime
	storedIdentity.ExpirationTime = currentTime.Add(longestMaxTTL)
	storedIdentity.PendingTime = identityDocParsed.PendingTime
	storedIdentity.DisallowReauthentication = disallowReauthentication

	// Don't cache the nonce if DisallowReauthentication is set
	if storedIdentity.DisallowReauthentication {
		storedIdentity.ClientNonce = ""
	}

	// Sanitize the nonce to a reasonable length
	if len(clientNonce) > 128 && !storedIdentity.DisallowReauthentication {
		return logical.ErrorResponse("client nonce exceeding the limit of 128 characters"), nil
	}

	if err = setWhitelistIdentityEntry(ctx, req.Storage, identityDocParsed.InstanceID, storedIdentity); err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Auth: &logical.Auth{
			Period:   roleEntry.Period,
			Policies: policies,
			Metadata: map[string]string{
				"instance_id":      identityDocParsed.InstanceID,
				"region":           identityDocParsed.Region,
				"account_id":       identityDocParsed.AccountID,
				"role_tag_max_ttl": rTagMaxTTL.String(),
				"role":             roleName,
				"ami_id":           identityDocParsed.AmiID,
			},
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       roleEntry.TTL,
				MaxTTL:    shortestMaxTTL,
			},
			Alias: &logical.Alias{
				Name: identityAlias,
			},
		},
	}

	// Return the nonce only if reauthentication is allowed and if the nonce
	// was not supplied by the user.
	if !disallowReauthentication && !clientNonceSupplied {
		// Echo the client nonce back. If nonce param was not supplied
		// to the endpoint at all (setting it to empty string does not
		// qualify here), callers should extract out the nonce from
		// this field for reauthentication requests.
		resp.Auth.Metadata["nonce"] = clientNonce
	}

	return resp, nil
}

// handleRoleTagLogin is used to fetch the role tag of the instance and
// verifies it to be correct.  Then the policies for the login request will be
// set off of the role tag, if certain criteria satisfies.
func (b *backend) handleRoleTagLogin(ctx context.Context, s logical.Storage, roleName string, roleEntry *awsRoleEntry, instance *ec2.Instance) (*roleTagLoginResponse, error) {
	if roleEntry == nil {
		return nil, fmt.Errorf("nil role entry")
	}
	if instance == nil {
		return nil, fmt.Errorf("nil instance")
	}

	// Input validation on instance is not performed here considering
	// that it would have been done in validateInstance method.
	tags := instance.Tags
	if tags == nil || len(tags) == 0 {
		return nil, fmt.Errorf("missing tag with key %q on the instance", roleEntry.RoleTag)
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
		return nil, fmt.Errorf("missing tag with key %q on the instance", roleEntry.RoleTag)
	}

	// Parse the role tag into a struct, extract the plaintext part of it and verify its HMAC
	rTag, err := b.parseAndVerifyRoleTagValue(ctx, s, rTagValue)
	if err != nil {
		return nil, err
	}

	// Check if the role name with which this login is being made is same
	// as the role name embedded in the tag.
	if rTag.Role != roleName {
		return nil, fmt.Errorf("role on the tag is not matching the role supplied")
	}

	// If instance_id was set on the role tag, check if the same instance is attempting to login
	if rTag.InstanceID != "" && rTag.InstanceID != *instance.InstanceId {
		return nil, fmt.Errorf("role tag is being used by an unauthorized instance")
	}

	// Check if the role tag is blacklisted
	blacklistEntry, err := b.lockedBlacklistRoleTagEntry(ctx, s, rTagValue)
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
		Policies:                 rTag.Policies,
		MaxTTL:                   rTag.MaxTTL,
		DisallowReauthentication: rTag.DisallowReauthentication,
	}, nil
}

// pathLoginRenew is used to renew an authenticated token
func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	authType, ok := req.Auth.Metadata["auth_type"]
	if !ok {
		// backwards compatibility for clients that have leases from before we added auth_type
		authType = ec2AuthType
	}

	if authType == ec2AuthType {
		return b.pathLoginRenewEc2(ctx, req, data)
	} else if authType == iamAuthType {
		return b.pathLoginRenewIam(ctx, req, data)
	} else {
		return nil, fmt.Errorf("unrecognized auth_type: %q", authType)
	}
}

func (b *backend) pathLoginRenewIam(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	canonicalArn := req.Auth.Metadata["canonical_arn"]
	if canonicalArn == "" {
		return nil, fmt.Errorf("unable to retrieve canonical ARN from metadata during renewal")
	}

	roleName := ""
	roleNameIfc, ok := req.Auth.InternalData["role_name"]
	if ok {
		roleName = roleNameIfc.(string)
	}
	if roleName == "" {
		return nil, fmt.Errorf("error retrieving role_name during renewal")
	}
	roleEntry, err := b.lockedAWSRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, fmt.Errorf("role entry not found")
	}

	// we don't really care what the inferred entity type was when the role was initially created. We
	// care about what the role currently requires. However, the metadata's inferred_entity_id is only
	// set when inferencing is turned on at initial login time. So, if inferencing is turned on, any
	// existing roles will NOT be able to renew tokens.
	// This might change later, but authenticating the actual inferred entity ID is NOT done if there
	// is no inferencing requested in the role. The reason is that authenticating the inferred entity
	// ID requires additional AWS IAM permissions that might not be present (e.g.,
	// ec2:DescribeInstances) as well as additional inferencing configuration (the inferred region).
	// So, for now, if you want to turn on inferencing, all clients must re-authenticate and cannot
	// renew existing tokens.
	if roleEntry.InferredEntityType != "" {
		if roleEntry.InferredEntityType == ec2EntityType {
			instanceID, ok := req.Auth.Metadata["inferred_entity_id"]
			if !ok {
				return nil, fmt.Errorf("no inferred entity ID in auth metadata")
			}
			instanceRegion, ok := req.Auth.Metadata["inferred_aws_region"]
			if !ok {
				return nil, fmt.Errorf("no inferred AWS region in auth metadata")
			}
			_, err := b.validateInstance(ctx, req.Storage, instanceID, instanceRegion, req.Auth.Metadata["account_id"])
			if err != nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("failed to verify instance ID %q: {{err}}", instanceID), err)
			}
		} else {
			return nil, fmt.Errorf("unrecognized entity_type in metadata: %q", roleEntry.InferredEntityType)
		}
	}

	// Note that the error messages below can leak a little bit of information about the role information
	// For example, if on renew, the client gets the "error parsing ARN..." error message, the client
	// will know that it's a wildcard bind (but not the actual bind), even if the client can't actually
	// read the role directly to know what the bind is. It's a relatively small amount of leakage, in
	// some fairly corner cases, and in the most likely error case (role has been changed to a new ARN),
	// the error message is identical.
	if len(roleEntry.BoundIamPrincipalARNs) > 0 {
		// We might not get here if all bindings were on the inferred entity, which we've already validated
		// above
		// As with logins, there are three ways to pass this check:
		// 1: clientUserId is in roleEntry.BoundIamPrincipalIDs (entries in roleEntry.BoundIamPrincipalIDs
		//    implies that roleEntry.ResolveAWSUniqueIDs is true)
		// 2: roleEntry.ResolveAWSUniqueIDs is false and canonical_arn is in roleEntry.BoundIamPrincipalARNs
		// 3: Full ARN matches one of the wildcard globs in roleEntry.BoundIamPrincipalARNs
		clientUserId, ok := req.Auth.Metadata["client_user_id"]
		switch {
		case ok && strutil.StrListContains(roleEntry.BoundIamPrincipalIDs, clientUserId): // check 1 passed
		case !roleEntry.ResolveAWSUniqueIDs && strutil.StrListContains(roleEntry.BoundIamPrincipalARNs, canonicalArn): // check 2 passed
		default:
			// check 3 is a bit more complex, so we do it last
			fullArn := b.getCachedUserId(clientUserId)
			if fullArn == "" {
				entity, err := parseIamArn(canonicalArn)
				if err != nil {
					return nil, errwrap.Wrapf(fmt.Sprintf("error parsing ARN %q: {{err}}", canonicalArn), err)
				}
				fullArn, err = b.fullArn(ctx, entity, req.Storage)
				if err != nil {
					return nil, errwrap.Wrapf(fmt.Sprintf("error looking up full ARN of entity %v: {{err}}", entity), err)
				}
				if fullArn == "" {
					return nil, fmt.Errorf("got empty string back when looking up full ARN of entity %v", entity)
				}
				if clientUserId != "" {
					b.setCachedUserId(clientUserId, fullArn)
				}
			}
			matchedWildcardBind := false
			for _, principalARN := range roleEntry.BoundIamPrincipalARNs {
				if strings.HasSuffix(principalARN, "*") && strutil.GlobbedStringsMatch(principalARN, fullArn) {
					matchedWildcardBind = true
					break
				}
			}
			if !matchedWildcardBind {
				return nil, fmt.Errorf("role no longer bound to ARN %q", canonicalArn)
			}
		}
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.TTL = roleEntry.TTL
	resp.Auth.MaxTTL = roleEntry.MaxTTL
	resp.Auth.Period = roleEntry.Period
	return resp, nil
}

func (b *backend) pathLoginRenewEc2(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	instanceID := req.Auth.Metadata["instance_id"]
	if instanceID == "" {
		return nil, fmt.Errorf("unable to fetch instance ID from metadata during renewal")
	}

	region := req.Auth.Metadata["region"]
	if region == "" {
		return nil, fmt.Errorf("unable to fetch region from metadata during renewal")
	}

	// Ensure backwards compatibility for older clients without account_id saved in metadata
	accountID, ok := req.Auth.Metadata["account_id"]
	if ok {
		if accountID == "" {
			return nil, fmt.Errorf("unable to fetch account_id from metadata during renewal")
		}
	}

	// Cross check that the instance is still in 'running' state
	_, err := b.validateInstance(ctx, req.Storage, instanceID, region, accountID)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to verify instance ID %q: {{err}}", instanceID), err)
	}

	storedIdentity, err := whitelistIdentityEntry(ctx, req.Storage, instanceID)
	if err != nil {
		return nil, err
	}
	if storedIdentity == nil {
		return nil, fmt.Errorf("failed to verify the whitelist identity entry for instance ID: %q", instanceID)
	}

	// Ensure that role entry is not deleted
	roleEntry, err := b.lockedAWSRole(ctx, req.Storage, storedIdentity.Role)
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

	// Re-evaluate the maxTTL bounds
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

	// Only LastUpdatedTime and ExpirationTime change and all other fields remain the same
	currentTime := time.Now()
	storedIdentity.LastUpdatedTime = currentTime
	storedIdentity.ExpirationTime = currentTime.Add(longestMaxTTL)

	// Updating the expiration time is required for the tidy operation on the
	// whitelist identity storage items
	if err = setWhitelistIdentityEntry(ctx, req.Storage, instanceID, storedIdentity); err != nil {
		return nil, err
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.TTL = roleEntry.TTL
	resp.Auth.MaxTTL = shortestMaxTTL
	resp.Auth.Period = roleEntry.Period
	return resp, nil
}

func (b *backend) pathLoginUpdateIam(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	method := data.Get("iam_http_request_method").(string)
	if method == "" {
		return logical.ErrorResponse("missing iam_http_request_method"), nil
	}

	// In the future, might consider supporting GET
	if method != "POST" {
		return logical.ErrorResponse("invalid iam_http_request_method; currently only 'POST' is supported"), nil
	}

	rawUrlB64 := data.Get("iam_request_url").(string)
	if rawUrlB64 == "" {
		return logical.ErrorResponse("missing iam_request_url"), nil
	}
	rawUrl, err := base64.StdEncoding.DecodeString(rawUrlB64)
	if err != nil {
		return logical.ErrorResponse("failed to base64 decode iam_request_url"), nil
	}
	parsedUrl, err := url.Parse(string(rawUrl))
	if err != nil {
		return logical.ErrorResponse("error parsing iam_request_url"), nil
	}

	// TODO: There are two potentially valid cases we're not yet supporting that would
	// necessitate this check being changed. First, if we support GET requests.
	// Second if we support presigned POST requests
	bodyB64 := data.Get("iam_request_body").(string)
	if bodyB64 == "" {
		return logical.ErrorResponse("missing iam_request_body"), nil
	}
	bodyRaw, err := base64.StdEncoding.DecodeString(bodyB64)
	if err != nil {
		return logical.ErrorResponse("failed to base64 decode iam_request_body"), nil
	}
	body := string(bodyRaw)

	headers := data.Get("iam_request_headers").(http.Header)
	if len(headers) == 0 {
		return logical.ErrorResponse("missing iam_request_headers"), nil
	}

	config, err := b.lockedClientConfigEntry(ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse("error getting configuration"), nil
	}

	endpoint := "https://sts.amazonaws.com"

	if config != nil {
		if config.IAMServerIdHeaderValue != "" {
			err = validateVaultHeaderValue(headers, parsedUrl, config.IAMServerIdHeaderValue)
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("error validating %s header: %v", iamServerIdHeader, err)), nil
			}
		}
		if config.STSEndpoint != "" {
			endpoint = config.STSEndpoint
		}
	}

	callerID, err := submitCallerIdentityRequest(method, endpoint, parsedUrl, body, headers)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error making upstream request: %v", err)), nil
	}

	entity, err := parseIamArn(callerID.Arn)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error parsing arn %q: %v", callerID.Arn, err)), nil
	}

	roleName := data.Get("role").(string)
	if roleName == "" {
		roleName = entity.FriendlyName
	}

	roleEntry, err := b.lockedAWSRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("entry for role %s not found", roleName)), nil
	}

	if roleEntry.AuthType != iamAuthType {
		return logical.ErrorResponse(fmt.Sprintf("auth method iam not allowed for role %s", roleName)), nil
	}

	identityConfigEntry, err := identityConfigEntry(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// This could either be a "userID:SessionID" (in the case of an assumed role) or just a "userID"
	// (in the case of an IAM user).
	callerUniqueId := strings.Split(callerID.UserId, ":")[0]
	identityAlias := ""
	switch identityConfigEntry.IAMAlias {
	case identityAliasRoleID:
		identityAlias = roleEntry.RoleID
	case identityAliasIAMUniqueID:
		identityAlias = callerUniqueId
	case identityAliasIAMFullArn:
		identityAlias = callerID.Arn
	}

	// If we're just looking up for MFA, return the Alias info
	if req.Operation == logical.AliasLookaheadOperation {
		return &logical.Response{
			Auth: &logical.Auth{
				Alias: &logical.Alias{
					Name: identityAlias,
				},
			},
		}, nil
	}

	// The role creation should ensure that either we're inferring this is an EC2 instance
	// or that we're binding an ARN
	if len(roleEntry.BoundIamPrincipalARNs) > 0 {
		// As with renews, there are three ways to pass this check:
		// 1: callerUniqueId is in roleEntry.BoundIamPrincipalIDs (entries in roleEntry.BoundIamPrincipalIDs
		//    implies that roleEntry.ResolveAWSUniqueIDs is true)
		// 2: roleEntry.ResolveAWSUniqueIDs is false and entity.canonicalArn() is in roleEntry.BoundIamPrincipalARNs
		// 3: Full ARN matches one of the wildcard globs in roleEntry.BoundIamPrincipalARNs
		// Need to be able to handle pathological configurations such as roleEntry.BoundIamPrincipalARNs looking something like:
		// arn:aw:iam::123456789012:{user/UserName,user/path/*,role/RoleName,role/path/*}
		switch {
		case strutil.StrListContains(roleEntry.BoundIamPrincipalIDs, callerUniqueId): // check 1 passed
		case !roleEntry.ResolveAWSUniqueIDs && strutil.StrListContains(roleEntry.BoundIamPrincipalARNs, entity.canonicalArn()): // check 2 passed
		default:
			// evaluate check 3
			fullArn := b.getCachedUserId(callerUniqueId)
			if fullArn == "" {
				fullArn, err = b.fullArn(ctx, entity, req.Storage)
				if err != nil {
					return logical.ErrorResponse(fmt.Sprintf("error looking up full ARN of entity %v: %v", entity, err)), nil
				}
				if fullArn == "" {
					return logical.ErrorResponse(fmt.Sprintf("got empty string back when looking up full ARN of entity %v", entity)), nil
				}
				b.setCachedUserId(callerUniqueId, fullArn)
			}
			matchedWildcardBind := false
			for _, principalARN := range roleEntry.BoundIamPrincipalARNs {
				if strings.HasSuffix(principalARN, "*") && strutil.GlobbedStringsMatch(principalARN, fullArn) {
					matchedWildcardBind = true
					break
				}
			}
			if !matchedWildcardBind {
				return logical.ErrorResponse(fmt.Sprintf("IAM Principal %q does not belong to the role %q", callerID.Arn, roleName)), nil
			}
		}
	}

	policies := roleEntry.Policies

	inferredEntityType := ""
	inferredEntityID := ""
	if roleEntry.InferredEntityType == ec2EntityType {
		instance, err := b.validateInstance(ctx, req.Storage, entity.SessionInfo, roleEntry.InferredAWSRegion, callerID.Account)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to verify %s as a valid EC2 instance in region %s", entity.SessionInfo, roleEntry.InferredAWSRegion)), nil
		}

		// build a fake identity doc to pass on metadata about the instance to verifyInstanceMeetsRoleRequirements
		identityDoc := &identityDocument{
			Tags:        nil, // Don't really need the tags, so not doing the work of converting them from Instance.Tags to identityDocument.Tags
			InstanceID:  *instance.InstanceId,
			AmiID:       *instance.ImageId,
			AccountID:   callerID.Account,
			Region:      roleEntry.InferredAWSRegion,
			PendingTime: instance.LaunchTime.Format(time.RFC3339),
		}

		validationError, err := b.verifyInstanceMeetsRoleRequirements(ctx, req.Storage, instance, roleEntry, roleName, identityDoc)
		if err != nil {
			return nil, err
		}
		if validationError != nil {
			return logical.ErrorResponse(fmt.Sprintf("error validating instance: %s", validationError)), nil
		}

		inferredEntityType = ec2EntityType
		inferredEntityID = entity.SessionInfo
	}

	resp := &logical.Response{
		Auth: &logical.Auth{
			Period:   roleEntry.Period,
			Policies: policies,
			Metadata: map[string]string{
				"client_arn":           callerID.Arn,
				"canonical_arn":        entity.canonicalArn(),
				"client_user_id":       callerUniqueId,
				"auth_type":            iamAuthType,
				"inferred_entity_type": inferredEntityType,
				"inferred_entity_id":   inferredEntityID,
				"inferred_aws_region":  roleEntry.InferredAWSRegion,
				"account_id":           entity.AccountNumber,
				"role_id":              roleEntry.RoleID,
			},
			InternalData: map[string]interface{}{
				"role_name": roleName,
				"role_id":   roleEntry.RoleID,
			},
			DisplayName: entity.FriendlyName,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       roleEntry.TTL,
				MaxTTL:    roleEntry.MaxTTL,
			},
			Alias: &logical.Alias{
				Name: identityAlias,
			},
		},
	}

	return resp, nil
}

// These two methods (hasValuesFor*) return two bools
// The first is a hasAll, that is, does the request have all the values
// necessary for this auth method
// The second is a hasAny, that is, does the request have any of the fields
// exclusive to this auth method
func hasValuesForEc2Auth(data *framework.FieldData) (bool, bool) {
	_, hasPkcs7 := data.GetOk("pkcs7")
	_, hasIdentity := data.GetOk("identity")
	_, hasSignature := data.GetOk("signature")
	return (hasPkcs7 || (hasIdentity && hasSignature)), (hasPkcs7 || hasIdentity || hasSignature)
}

func hasValuesForIamAuth(data *framework.FieldData) (bool, bool) {
	_, hasRequestMethod := data.GetOk("iam_http_request_method")
	_, hasRequestURL := data.GetOk("iam_request_url")
	_, hasRequestBody := data.GetOk("iam_request_body")
	_, hasRequestHeaders := data.GetOk("iam_request_headers")
	return (hasRequestMethod && hasRequestURL && hasRequestBody && hasRequestHeaders),
		(hasRequestMethod || hasRequestURL || hasRequestBody || hasRequestHeaders)
}

func parseIamArn(iamArn string) (*iamEntity, error) {
	// iamArn should look like one of the following:
	// 1. arn:aws:iam::<account_id>:<entity_type>/<UserName>
	// 2. arn:aws:sts::<account_id>:assumed-role/<RoleName>/<RoleSessionName>
	// if we get something like 2, then we want to transform that back to what
	// most people would expect, which is arn:aws:iam::<account_id>:role/<RoleName>
	var entity iamEntity
	fullParts := strings.Split(iamArn, ":")
	if len(fullParts) != 6 {
		return nil, fmt.Errorf("unrecognized arn: contains %d colon-separated parts, expected 6", len(fullParts))
	}
	if fullParts[0] != "arn" {
		return nil, fmt.Errorf("unrecognized arn: does not begin with \"arn:\"")
	}
	// normally aws, but could be aws-cn or aws-us-gov
	entity.Partition = fullParts[1]
	if fullParts[2] != "iam" && fullParts[2] != "sts" {
		return nil, fmt.Errorf("unrecognized service: %v, not one of iam or sts", fullParts[2])
	}
	// fullParts[3] is the region, which doesn't matter for AWS IAM entities
	entity.AccountNumber = fullParts[4]
	// fullParts[5] would now be something like user/<UserName> or assumed-role/<RoleName>/<RoleSessionName>
	parts := strings.Split(fullParts[5], "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("unrecognized arn: %q contains fewer than 2 slash-separated parts", fullParts[5])
	}
	entity.Type = parts[0]
	entity.Path = strings.Join(parts[1:len(parts)-1], "/")
	entity.FriendlyName = parts[len(parts)-1]
	// now, entity.FriendlyName should either be <UserName> or <RoleName>
	switch entity.Type {
	case "assumed-role":
		// Assumed roles don't have paths and have a slightly different format
		// parts[2] is <RoleSessionName>
		entity.Path = ""
		entity.FriendlyName = parts[1]
		entity.SessionInfo = parts[2]
	case "user":
	case "role":
	case "instance-profile":
	default:
		return &iamEntity{}, fmt.Errorf("unrecognized principal type: %q", entity.Type)
	}
	return &entity, nil
}

func validateVaultHeaderValue(headers http.Header, requestUrl *url.URL, requiredHeaderValue string) error {
	providedValue := ""
	for k, v := range headers {
		if strings.EqualFold(iamServerIdHeader, k) {
			providedValue = strings.Join(v, ",")
			break
		}
	}
	if providedValue == "" {
		return fmt.Errorf("missing header %q", iamServerIdHeader)
	}

	// NOT doing a constant time compare here since the value is NOT intended to be secret
	if providedValue != requiredHeaderValue {
		return fmt.Errorf("expected %q but got %q", requiredHeaderValue, providedValue)
	}

	if authzHeaders, ok := headers["Authorization"]; ok {
		// authzHeader looks like AWS4-HMAC-SHA256 Credential=AKI..., SignedHeaders=host;x-amz-date;x-vault-awsiam-id, Signature=...
		// We need to extract out the SignedHeaders
		re := regexp.MustCompile(".*SignedHeaders=([^,]+)")
		authzHeader := strings.Join(authzHeaders, ",")
		matches := re.FindSubmatch([]byte(authzHeader))
		if len(matches) < 1 {
			return fmt.Errorf("vault header wasn't signed")
		}
		if len(matches) > 2 {
			return fmt.Errorf("found multiple SignedHeaders components")
		}
		signedHeaders := string(matches[1])
		return ensureHeaderIsSigned(signedHeaders, iamServerIdHeader)
	}
	// TODO: If we support GET requests, then we need to parse the X-Amz-SignedHeaders
	// argument out of the query string and search in there for the header value
	return fmt.Errorf("missing Authorization header")
}

func buildHttpRequest(method, endpoint string, parsedUrl *url.URL, body string, headers http.Header) *http.Request {
	// This is all a bit complicated because the AWS signature algorithm requires that
	// the Host header be included in the signed headers. See
	// http://docs.aws.amazon.com/general/latest/gr/sigv4-create-canonical-request.html
	// The use cases we want to support, in order of increasing complexity, are:
	// 1. All defaults (client assumes sts.amazonaws.com and server has no override)
	// 2. Alternate STS regions: client wants to go to a specific region, in which case
	//    Vault must be configured with that endpoint as well. The client's signed request
	//    will include a signature over what the client expects the Host header to be,
	//    so we cannot change that and must match.
	// 3. Alternate STS regions with a proxy that is transparent to Vault's clients.
	//    In this case, Vault is aware of the proxy, as the proxy is configured as the
	//    endpoint, but the clients should NOT be aware of the proxy (because STS will
	//    not be aware of the proxy)
	// It's also annoying because:
	// 1. The AWS Sigv4 algorithm requires the Host header to be defined
	// 2. Some of the official SDKs (at least botocore and aws-sdk-go) don't actually
	//    include an explicit Host header in the HTTP requests they generate, relying on
	//    the underlying HTTP library to do that for them.
	// 3. To get a validly signed request, the SDKs check if a Host header has been set
	//    and, if not, add an inferred host header (based on the URI) to the internal
	//    data structure used for calculating the signature, but never actually expose
	//    that to clients. So then they just "hope" that the underlying library actually
	//    adds the right Host header which was included in the signature calculation.
	// We could either explicitly require all Vault clients to explicitly add the Host header
	// in the encoded request, or we could also implicitly infer it from the URI.
	// We choose to support both -- allow you to explicitly set a Host header, but if not,
	// infer one from the URI.
	// HOWEVER, we have to preserve the request URI portion of the client's
	// URL because the GetCallerIdentity Action can be encoded in either the body
	// or the URL. So, we need to rebuild the URL sent to the http library to have the
	// custom, Vault-specified endpoint with the client-side request parameters.
	targetUrl := fmt.Sprintf("%s/%s", endpoint, parsedUrl.RequestURI())
	request, err := http.NewRequest(method, targetUrl, strings.NewReader(body))
	if err != nil {
		return nil
	}
	request.Host = parsedUrl.Host
	for k, vals := range headers {
		for _, val := range vals {
			request.Header.Add(k, val)
		}
	}
	return request
}

func ensureHeaderIsSigned(signedHeaders, headerToSign string) error {
	// Not doing a constant time compare here, the values aren't secret
	for _, header := range strings.Split(signedHeaders, ";") {
		if header == strings.ToLower(headerToSign) {
			return nil
		}
	}
	return fmt.Errorf("vault header wasn't signed")
}

func parseGetCallerIdentityResponse(response string) (GetCallerIdentityResponse, error) {
	decoder := xml.NewDecoder(strings.NewReader(response))
	result := GetCallerIdentityResponse{}
	err := decoder.Decode(&result)
	return result, err
}

func submitCallerIdentityRequest(method, endpoint string, parsedUrl *url.URL, body string, headers http.Header) (*GetCallerIdentityResult, error) {
	// NOTE: We need to ensure we're calling STS, instead of acting as an unintended network proxy
	// The protection against this is that this method will only call the endpoint specified in the
	// client config (defaulting to sts.amazonaws.com), so it would require a Vault admin to override
	// the endpoint to talk to alternate web addresses
	request := buildHttpRequest(method, endpoint, parsedUrl, body, headers)
	client := cleanhttp.DefaultClient()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, errwrap.Wrapf("error making request: {{err}}", err)
	}
	if response != nil {
		defer response.Body.Close()
	}
	// we check for status code afterwards to also print out response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("received error code %d from STS: %s", response.StatusCode, string(responseBody))
	}
	callerIdentityResponse, err := parseGetCallerIdentityResponse(string(responseBody))
	if err != nil {
		return nil, fmt.Errorf("error parsing STS response")
	}
	return &callerIdentityResponse.GetCallerIdentityResult[0], nil
}

type GetCallerIdentityResponse struct {
	XMLName                 xml.Name                  `xml:"GetCallerIdentityResponse"`
	GetCallerIdentityResult []GetCallerIdentityResult `xml:"GetCallerIdentityResult"`
	ResponseMetadata        []ResponseMetadata        `xml:"ResponseMetadata"`
}

type GetCallerIdentityResult struct {
	Arn     string `xml:"Arn"`
	UserId  string `xml:"UserId"`
	Account string `xml:"Account"`
}

type ResponseMetadata struct {
	RequestId string `xml:"RequestId"`
}

// identityDocument represents the items of interest from the EC2 instance
// identity document
type identityDocument struct {
	Tags        map[string]interface{} `json:"tags,omitempty"`
	InstanceID  string                 `json:"instanceId,omitempty"`
	AmiID       string                 `json:"imageId,omitempty"`
	AccountID   string                 `json:"accountId,omitempty"`
	Region      string                 `json:"region,omitempty"`
	PendingTime string                 `json:"pendingTime,omitempty"`
}

// roleTagLoginResponse represents the return values required after the process
// of verifying a role tag login
type roleTagLoginResponse struct {
	Policies                 []string      `json:"policies"`
	MaxTTL                   time.Duration `json:"max_ttl"`
	DisallowReauthentication bool          `json:"disallow_reauthentication"`
}

type iamEntity struct {
	Partition     string
	AccountNumber string
	Type          string
	Path          string
	FriendlyName  string
	SessionInfo   string
}

// Returns a Vault-internal canonical ARN for referring to an IAM entity
func (e *iamEntity) canonicalArn() string {
	entityType := e.Type
	// canonicalize "assumed-role" into "role"
	if entityType == "assumed-role" {
		entityType = "role"
	}
	// Annoyingly, the assumed-role entity type doesn't have the Path of the role which was assumed
	// So, we "canonicalize" it by just completely dropping the path. The other option would be to
	// make an AWS API call to look up the role by FriendlyName, which introduces more complexity to
	// code and test, and it also breaks backwards compatibility in an area where we would really want
	// it
	return fmt.Sprintf("arn:%s:iam::%s:%s/%s", e.Partition, e.AccountNumber, entityType, e.FriendlyName)
}

// This returns the "full" ARN of an iamEntity, how it would be referred to in AWS proper
func (b *backend) fullArn(ctx context.Context, e *iamEntity, s logical.Storage) (string, error) {
	// Not assuming path is reliable for any entity types
	client, err := b.clientIAM(ctx, s, getAnyRegionForAwsPartition(e.Partition).ID(), e.AccountNumber)
	if err != nil {
		return "", errwrap.Wrapf("error creating IAM client: {{err}}", err)
	}

	switch e.Type {
	case "user":
		input := iam.GetUserInput{
			UserName: aws.String(e.FriendlyName),
		}
		resp, err := client.GetUser(&input)
		if err != nil {
			return "", errwrap.Wrapf(fmt.Sprintf("error fetching user %q: {{err}}", e.FriendlyName), err)
		}
		if resp == nil {
			return "", fmt.Errorf("nil response from GetUser")
		}
		return *(resp.User.Arn), nil
	case "assumed-role":
		fallthrough
	case "role":
		input := iam.GetRoleInput{
			RoleName: aws.String(e.FriendlyName),
		}
		resp, err := client.GetRole(&input)
		if err != nil {
			return "", errwrap.Wrapf(fmt.Sprintf("error fetching role %q: {{err}}", e.FriendlyName), err)
		}
		if resp == nil {
			return "", fmt.Errorf("nil response form GetRole")
		}
		return *(resp.Role.Arn), nil
	default:
		return "", fmt.Errorf("unrecognized entity type: %s", e.Type)
	}
}

const iamServerIdHeader = "X-Vault-AWS-IAM-Server-ID"

const pathLoginSyn = `
Authenticates an EC2 instance with Vault.
`

const pathLoginDesc = `
Authenticate AWS entities, either an arbitrary IAM principal or EC2 instances.

IAM principals are authenticated by processing a signed sts:GetCallerIdentity
request and then parsing the response to see who signed the request. Optionally,
the caller can be inferred to be another AWS entity type, with EC2 instances
the only currently supported entity type, and additional filtering can be
implemented based on that inferred type.

An EC2 instance is authenticated using the PKCS#7 signature of the instance identity
document and a client created nonce. This nonce should be unique and should be used by
the instance for all future logins, unless 'disallow_reauthentication' option on the
registered role is enabled, in which case client nonce is optional.

First login attempt, creates a whitelist entry in Vault associating the instance to the nonce
provided. All future logins will succeed only if the client nonce matches the nonce in the
whitelisted entry.

By default, a cron task will periodically look for expired entries in the whitelist
and deletes them. The duration to periodically run this, is one hour by default.
However, this can be configured using the 'config/tidy/identities' endpoint. This tidy
action can be triggered via the API as well, using the 'tidy/identities' endpoint.
`
