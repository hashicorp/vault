package awsauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const roleTagVersion = "v1"

func pathRoleTag(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("role") + "/tag$",
		Fields: map[string]*framework.FieldSchema{
			"role": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},

			"instance_id": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Instance ID for which this tag is intended for.
If set, the created tag can only be used by the instance with the given ID.`,
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Policies to be associated with the tag. If set, must be a subset of the role's policies. If set, but set to an empty value, only the 'default' policy will be given to issued tokens.",
			},

			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     0,
				Description: "If set, specifies the maximum allowed token lifetime.",
			},

			"allow_instance_migration": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, allows migration of the underlying instance where the client resides. This keys off of pendingTime in the metadata document, so essentially, this disables the client nonce check whenever the instance is migrated to a new host and pendingTime is newer than the previously-remembered time. Use with caution.",
			},

			"disallow_reauthentication": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, only allows a single token to be granted per instance ID. In order to perform a fresh login, the entry in whitelist for the instance ID needs to be cleared using the 'auth/aws-ec2/identity-whitelist/<instance_id>' endpoint.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRoleTagUpdate,
		},

		HelpSynopsis:    pathRoleTagSyn,
		HelpDescription: pathRoleTagDesc,
	}
}

// pathRoleTagUpdate is used to create an EC2 instance tag which will
// identify the Vault resources that the instance will be authorized for.
func (b *backend) pathRoleTagUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	roleName := strings.ToLower(data.Get("role").(string))
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	// Fetch the role entry
	roleEntry, err := b.lockedAWSRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("entry not found for role %s", roleName)), nil
	}

	// If RoleTag is empty, disallow creation of tag.
	if roleEntry.RoleTag == "" {
		return logical.ErrorResponse("tag creation is not enabled for this role"), nil
	}

	// There should be a HMAC key present in the role entry
	if roleEntry.HMACKey == "" {
		// Not being able to find the HMACKey is an internal error
		return nil, fmt.Errorf("failed to find the HMAC key")
	}

	resp := &logical.Response{}

	// Instance ID is an optional field.
	instanceID := strings.ToLower(data.Get("instance_id").(string))

	// If no policies field was not supplied, then the tag should inherit all the policies
	// on the role. But, it was provided, but set to empty explicitly, only "default" policy
	// should be inherited. So, by leaving the policies var unset to anything when it is not
	// supplied, we ensure that it inherits all the policies on the role.
	var policies []string
	policiesRaw, ok := data.GetOk("policies")
	if ok {
		policies = policyutil.ParsePolicies(policiesRaw)
	}
	if !strutil.StrListSubset(roleEntry.Policies, policies) {
		resp.AddWarning("Policies on the tag are not a subset of the policies set on the role. Login will not be allowed with this tag unless the role policies are updated.")
	}

	// This is an optional field.
	disallowReauthentication := data.Get("disallow_reauthentication").(bool)

	// This is an optional field.
	allowInstanceMigration := data.Get("allow_instance_migration").(bool)
	if allowInstanceMigration && !roleEntry.AllowInstanceMigration {
		resp.AddWarning("Role does not allow instance migration. Login will not be allowed with this tag unless the role value is updated.")
	}

	if disallowReauthentication && allowInstanceMigration {
		return logical.ErrorResponse("cannot set both disallow_reauthentication and allow_instance_migration"), nil
	}

	// max_ttl for the role tag should be less than the max_ttl set on the role.
	maxTTL := time.Duration(data.Get("max_ttl").(int)) * time.Second

	// max_ttl on the tag should not be greater than the system view's max_ttl value.
	if maxTTL > b.System().MaxLeaseTTL() {
		resp.AddWarning(fmt.Sprintf("Given max TTL of %d is greater than the mount maximum of %d seconds, and will be capped at login time.", maxTTL/time.Second, b.System().MaxLeaseTTL()/time.Second))
	}
	// If max_ttl is set for the role, check the bounds for tag's max_ttl value using that.
	if roleEntry.MaxTTL != time.Duration(0) && maxTTL > roleEntry.MaxTTL {
		resp.AddWarning(fmt.Sprintf("Given max TTL of %d is greater than the role maximum of %d seconds, and will be capped at login time.", maxTTL/time.Second, roleEntry.MaxTTL/time.Second))
	}

	if maxTTL < time.Duration(0) {
		return logical.ErrorResponse("max_ttl cannot be negative"), nil
	}

	// Create a random nonce.
	nonce, err := createRoleTagNonce()
	if err != nil {
		return nil, err
	}

	// Create a role tag out of all the information provided.
	rTagValue, err := createRoleTagValue(&roleTag{
		Version:                  roleTagVersion,
		Role:                     roleName,
		Nonce:                    nonce,
		Policies:                 policies,
		MaxTTL:                   maxTTL,
		InstanceID:               instanceID,
		DisallowReauthentication: disallowReauthentication,
		AllowInstanceMigration:   allowInstanceMigration,
	}, roleEntry)
	if err != nil {
		return nil, err
	}

	// Return the key to be used for the tag and the value to be used for that tag key.
	// This key value pair should be set on the EC2 instance.
	resp.Data = map[string]interface{}{
		"tag_key":   roleEntry.RoleTag,
		"tag_value": rTagValue,
	}

	return resp, nil
}

// createRoleTagValue prepares the plaintext version of the role tag,
// and appends a HMAC of the plaintext value to it, before returning.
func createRoleTagValue(rTag *roleTag, roleEntry *awsRoleEntry) (string, error) {
	if rTag == nil {
		return "", fmt.Errorf("nil role tag")
	}

	if roleEntry == nil {
		return "", fmt.Errorf("nil role entry")
	}

	// Attach version, nonce, policies and maxTTL to the role tag value.
	rTagPlaintext, err := prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		return "", err
	}

	// Attach HMAC to tag's plaintext and return.
	return appendHMAC(rTagPlaintext, roleEntry)
}

// Takes in the plaintext part of the role tag, creates a HMAC of it and returns
// a role tag value containing both the plaintext part and the HMAC part.
func appendHMAC(rTagPlaintext string, roleEntry *awsRoleEntry) (string, error) {
	if rTagPlaintext == "" {
		return "", fmt.Errorf("empty role tag plaintext string")
	}

	if roleEntry == nil {
		return "", fmt.Errorf("nil role entry")
	}

	// Create the HMAC of the value
	hmacB64, err := createRoleTagHMACBase64(roleEntry.HMACKey, rTagPlaintext)
	if err != nil {
		return "", err
	}

	// attach the HMAC to the value
	rTagValue := fmt.Sprintf("%s:%s", rTagPlaintext, hmacB64)

	// This limit of 255 is enforced on the EC2 instance. Hence complying to that here.
	if len(rTagValue) > 255 {
		return "", fmt.Errorf("role tag 'value' exceeding the limit of 255 characters")
	}

	return rTagValue, nil
}

// verifyRoleTagValue rebuilds the role tag's plaintext part, computes the HMAC
// from it using the role specific HMAC key and compares it with the received HMAC.
func verifyRoleTagValue(rTag *roleTag, roleEntry *awsRoleEntry) (bool, error) {
	if rTag == nil {
		return false, fmt.Errorf("nil role tag")
	}

	if roleEntry == nil {
		return false, fmt.Errorf("nil role entry")
	}

	// Fetch the plaintext part of role tag
	rTagPlaintext, err := prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		return false, err
	}

	// Compute the HMAC of the plaintext
	hmacB64, err := createRoleTagHMACBase64(roleEntry.HMACKey, rTagPlaintext)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare([]byte(rTag.HMAC), []byte(hmacB64)) == 1, nil
}

// prepareRoleTagPlaintextValue builds the role tag value without the HMAC in it.
func prepareRoleTagPlaintextValue(rTag *roleTag) (string, error) {
	if rTag == nil {
		return "", fmt.Errorf("nil role tag")
	}
	if rTag.Version == "" {
		return "", fmt.Errorf("missing version")
	}
	if rTag.Nonce == "" {
		return "", fmt.Errorf("missing nonce")
	}
	if rTag.Role == "" {
		return "", fmt.Errorf("missing role")
	}

	// Attach Version, Nonce, Role, DisallowReauthentication and AllowInstanceMigration
	// fields to the role tag.
	value := fmt.Sprintf("%s:%s:r=%s:d=%s:m=%s", rTag.Version, rTag.Nonce, rTag.Role, strconv.FormatBool(rTag.DisallowReauthentication), strconv.FormatBool(rTag.AllowInstanceMigration))

	// Attach the policies only if they are specified.
	if len(rTag.Policies) != 0 {
		value = fmt.Sprintf("%s:p=%s", value, strings.Join(rTag.Policies, ","))
	}

	// Attach instance_id if set.
	if rTag.InstanceID != "" {
		value = fmt.Sprintf("%s:i=%s", value, rTag.InstanceID)
	}

	// Attach max_ttl if it is provided.
	if int(rTag.MaxTTL.Seconds()) > 0 {
		value = fmt.Sprintf("%s:t=%d", value, int(rTag.MaxTTL.Seconds()))
	}

	return value, nil
}

// Parses the tag from string form into a struct form. This method
// also verifies the correctness of the parsed role tag.
func (b *backend) parseAndVerifyRoleTagValue(s logical.Storage, tag string) (*roleTag, error) {
	tagItems := strings.Split(tag, ":")

	// Tag must contain version, nonce, policies and HMAC
	if len(tagItems) < 4 {
		return nil, fmt.Errorf("invalid tag")
	}

	rTag := &roleTag{}

	// Cache the HMAC value. The last item in the collection.
	rTag.HMAC = tagItems[len(tagItems)-1]

	// Remove the HMAC from the list.
	tagItems = tagItems[:len(tagItems)-1]

	// Version will be the first element.
	rTag.Version = tagItems[0]
	if rTag.Version != roleTagVersion {
		return nil, fmt.Errorf("invalid role tag version")
	}

	// Nonce will be the second element.
	rTag.Nonce = tagItems[1]

	// Delete the version and nonce from the list.
	tagItems = tagItems[2:]

	for _, tagItem := range tagItems {
		var err error
		switch {
		case strings.Contains(tagItem, "i="):
			rTag.InstanceID = strings.TrimPrefix(tagItem, "i=")
		case strings.Contains(tagItem, "r="):
			rTag.Role = strings.TrimPrefix(tagItem, "r=")
		case strings.Contains(tagItem, "p="):
			rTag.Policies = strings.Split(strings.TrimPrefix(tagItem, "p="), ",")
		case strings.Contains(tagItem, "d="):
			rTag.DisallowReauthentication, err = strconv.ParseBool(strings.TrimPrefix(tagItem, "d="))
			if err != nil {
				return nil, err
			}
		case strings.Contains(tagItem, "m="):
			rTag.AllowInstanceMigration, err = strconv.ParseBool(strings.TrimPrefix(tagItem, "m="))
			if err != nil {
				return nil, err
			}
		case strings.Contains(tagItem, "t="):
			rTag.MaxTTL, err = time.ParseDuration(fmt.Sprintf("%ss", strings.TrimPrefix(tagItem, "t=")))
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unrecognized item %s in tag", tagItem)
		}
	}

	if rTag.Role == "" {
		return nil, fmt.Errorf("missing role name")
	}

	roleEntry, err := b.lockedAWSRole(s, rTag.Role)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, fmt.Errorf("entry not found for %s", rTag.Role)
	}

	// Create a HMAC of the plaintext value of role tag and compare it with the given value.
	verified, err := verifyRoleTagValue(rTag, roleEntry)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, fmt.Errorf("role tag signature verification failed")
	}

	return rTag, nil
}

// Creates base64 encoded HMAC using a per-role key.
func createRoleTagHMACBase64(key, value string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("invalid HMAC key")
	}
	hm := hmac.New(sha256.New, []byte(key))
	hm.Write([]byte(value))

	// base64 encode the hmac bytes.
	return base64.StdEncoding.EncodeToString(hm.Sum(nil)), nil
}

// Creates a base64 encoded random nonce.
func createRoleTagNonce() (string, error) {
	if uuidBytes, err := uuid.GenerateRandomBytes(8); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(uuidBytes), nil
	}
}

// Struct roleTag represents a role tag in a struc form.
type roleTag struct {
	Version                  string        `json:"version" structs:"version" mapstructure:"version"`
	InstanceID               string        `json:"instance_id" structs:"instance_id" mapstructure:"instance_id"`
	Nonce                    string        `json:"nonce" structs:"nonce" mapstructure:"nonce"`
	Policies                 []string      `json:"policies" structs:"policies" mapstructure:"policies"`
	MaxTTL                   time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	Role                     string        `json:"role" structs:"role" mapstructure:"role"`
	HMAC                     string        `json:"hmac" structs:"hmac" mapstructure:"hmac"`
	DisallowReauthentication bool          `json:"disallow_reauthentication" structs:"disallow_reauthentication" mapstructure:"disallow_reauthentication"`
	AllowInstanceMigration   bool          `json:"allow_instance_migration" structs:"allow_instance_migration" mapstructure:"allow_instance_migration"`
}

func (rTag1 *roleTag) Equal(rTag2 *roleTag) bool {
	return rTag1 != nil &&
		rTag2 != nil &&
		rTag1.Version == rTag2.Version &&
		rTag1.Nonce == rTag2.Nonce &&
		policyutil.EquivalentPolicies(rTag1.Policies, rTag2.Policies) &&
		rTag1.MaxTTL == rTag2.MaxTTL &&
		rTag1.Role == rTag2.Role &&
		rTag1.HMAC == rTag2.HMAC &&
		rTag1.InstanceID == rTag2.InstanceID &&
		rTag1.DisallowReauthentication == rTag2.DisallowReauthentication &&
		rTag1.AllowInstanceMigration == rTag2.AllowInstanceMigration
}

const pathRoleTagSyn = `
Create a tag on a role in order to be able to further restrict the capabilities of a role.
`

const pathRoleTagDesc = `
If there are needs to apply only a subset of role's capabilities to any specific
instance, create a role tag using this endpoint and attach the tag on the instance
before performing login.

To be able to create a role tag, the 'role_tag' option on the role should be
enabled via the endpoint 'role/<role>'. Also, the policies to be associated
with the tag should be a subset of the policies associated with the registered role.

This endpoint will return both the 'key' and the 'value' of the tag to be set
on the EC2 instance.
`
