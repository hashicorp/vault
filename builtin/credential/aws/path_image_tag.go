package aws

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
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const roleTagVersion = "v1"

func pathImageTag(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "image/" + framework.GenericNameRegex("ami_id") + "/roletag$",
		Fields: map[string]*framework.FieldSchema{
			"ami_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "AMI ID to create a tag for.",
			},

			"instance_id": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Instance ID for which this tag is intended for.
This is an optional field, but if set, the created tag can only be used by the instance with the given ID.`,
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Policies to be associated with the tag.",
			},

			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     0,
				Description: "The maximum allowed lease duration.",
			},

			"allow_instance_migration": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, allows migration of the underlying instance where the client resides. This keys off of pendingTime in the metadata document, so essentially, this disables the client nonce check whenever the instance is migrated to a new host and pendingTime is newer than the previously-remembered time. Use with caution.",
			},

			"disallow_reauthentication": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, only allows a single token to be granted per instance ID. In order to perform a fresh login, the entry in whitelist for the instance ID needs to be cleared using 'auth/aws/whitelist/identity/<instance_id>' endpoint.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathImageTagUpdate,
		},

		HelpSynopsis:    pathImageTagSyn,
		HelpDescription: pathImageTagDesc,
	}
}

// pathImageTagUpdate is used to create an EC2 instance tag which will
// identify the Vault resources that the instance will be authorized for.
func (b *backend) pathImageTagUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	amiID := strings.ToLower(data.Get("ami_id").(string))
	if amiID == "" {
		return logical.ErrorResponse("missing ami_id"), nil
	}

	// Instance ID is an optional field.
	instanceID := strings.ToLower(data.Get("instance_id").(string))

	// Parse the given policies into a slice and add 'default' if not provided.
	// Remove all other policies if 'root' is present.
	policies := policyutil.ParsePolicies(data.Get("policies").(string))

	// This is an optional field.
	disallowReauthentication := data.Get("disallow_reauthentication").(bool)

	// This is an optional field.
	allowInstanceMigration := data.Get("allow_instance_migration").(bool)

	// Fetch the image entry corresponding to the AMI ID
	imageEntry, err := awsImage(req.Storage, amiID)
	if err != nil {
		return nil, err
	}
	if imageEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("entry not found for AMI %s", amiID)), nil
	}

	// If RoleTag is empty, disallow creation of tag.
	if imageEntry.RoleTag == "" {
		return logical.ErrorResponse("tag creation is not enabled for this image"), nil
	}

	// There should be a HMAC key present in the image entry
	if imageEntry.HMACKey == "" {
		// Not being able to find the HMACKey is an internal error
		return nil, fmt.Errorf("failed to find the HMAC key")
	}

	// Create a random nonce.
	nonce, err := createRoleTagNonce()
	if err != nil {
		return nil, err
	}

	// max_ttl for the role tag should be less than the max_ttl set on the image.
	maxTTL := time.Duration(data.Get("max_ttl").(int)) * time.Second

	// max_ttl on the tag should not be greater than the system view's max_ttl value.
	if maxTTL > b.System().MaxLeaseTTL() {
		return logical.ErrorResponse(fmt.Sprintf("Registered AMI does not have a max_ttl set. So, the given TTL of %d seconds should be less than the max_ttl set for the corresponding backend mount of %d seconds.", maxTTL/time.Second, b.System().MaxLeaseTTL()/time.Second)), nil
	}

	// If max_ttl is set for the image, check the bounds for tag's max_ttl value using that.
	if imageEntry.MaxTTL != time.Duration(0) && maxTTL > imageEntry.MaxTTL {
		return logical.ErrorResponse(fmt.Sprintf("Given TTL of %d seconds greater than the max_ttl set for the corresponding image of %d seconds", maxTTL/time.Second, imageEntry.MaxTTL/time.Second)), nil
	}

	if maxTTL < time.Duration(0) {
		return logical.ErrorResponse("max_ttl cannot be negative"), nil
	}

	// Create a role tag out of all the information provided.
	rTagValue, err := createRoleTagValue(&roleTag{
		Version:                  roleTagVersion,
		AmiID:                    amiID,
		Nonce:                    nonce,
		Policies:                 policies,
		MaxTTL:                   maxTTL,
		InstanceID:               instanceID,
		DisallowReauthentication: disallowReauthentication,
		AllowInstanceMigration:   allowInstanceMigration,
	}, imageEntry)
	if err != nil {
		return nil, err
	}

	// Return the key to be used for the tag and the value to be used for that tag key.
	// This key value pair should be set on the EC2 instance.
	return &logical.Response{
		Data: map[string]interface{}{
			"tag_key":   imageEntry.RoleTag,
			"tag_value": rTagValue,
		},
	}, nil
}

// createRoleTagValue prepares the plaintext version of the role tag,
// and appends a HMAC of the plaintext value to it, before returning.
func createRoleTagValue(rTag *roleTag, imageEntry *awsImageEntry) (string, error) {
	if rTag == nil {
		return "", fmt.Errorf("nil role tag")
	}

	if imageEntry == nil {
		return "", fmt.Errorf("nil image entry")
	}

	// Attach version, nonce, policies and maxTTL to the role tag value.
	rTagPlaintext, err := prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		return "", err
	}

	// Attach HMAC to tag's plaintext and return.
	return appendHMAC(rTagPlaintext, imageEntry)
}

// Takes in the plaintext part of the role tag, creates a HMAC of it and returns
// a role tag value containing both the plaintext part and the HMAC part.
func appendHMAC(rTagPlaintext string, imageEntry *awsImageEntry) (string, error) {
	if rTagPlaintext == "" {
		return "", fmt.Errorf("empty role tag plaintext string")
	}

	if imageEntry == nil {
		return "", fmt.Errorf("nil image entry")
	}

	// Create the HMAC of the value
	hmacB64, err := createRoleTagHMACBase64(imageEntry.HMACKey, rTagPlaintext)
	if err != nil {
		return "", err
	}

	// attach the HMAC to the value
	rTagValue := fmt.Sprintf("%s:%s", rTagPlaintext, hmacB64)

	// This limit of 255 is enforced on the EC2 instance. Hence complying to it here.
	if len(rTagValue) > 255 {
		return "", fmt.Errorf("role tag 'value' exceeding the limit of 255 characters")
	}

	return rTagValue, nil
}

// verifyRoleTagValue rebuilds the role tag value without the HMAC,
// computes the HMAC from it using the backend specific key and
// compares it with the received HMAC.
func verifyRoleTagValue(rTag *roleTag, imageEntry *awsImageEntry) (bool, error) {
	if rTag == nil {
		return false, fmt.Errorf("nil role tag")
	}

	if imageEntry == nil {
		return false, fmt.Errorf("nil image entry")
	}

	// Fetch the plaintext part of role tag
	rTagPlaintext, err := prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		return false, err
	}

	// Compute the HMAC of the plaintext
	hmacB64, err := createRoleTagHMACBase64(imageEntry.HMACKey, rTagPlaintext)
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
	if rTag.AmiID == "" {
		return "", fmt.Errorf("missing ami_id")
	}

	// This avoids an empty policy, ":p=:" in the role tag.
	if rTag.Policies == nil || len(rTag.Policies) == 0 {
		rTag.Policies = []string{"default"}
	}

	// Attach Version, Nonce, AMI ID, Policies, DisallowReauthentication fields.
	value := fmt.Sprintf("%s:%s:a=%s:p=%s:d=%s:m=%s", rTag.Version, rTag.Nonce, rTag.AmiID, strings.Join(rTag.Policies, ","), strconv.FormatBool(rTag.DisallowReauthentication), strconv.FormatBool(rTag.AllowInstanceMigration))

	// Attach instance_id if set.
	if rTag.InstanceID != "" {
		value = fmt.Sprintf("%s:i=%s", value, rTag.InstanceID)
	}

	// Attach max_ttl if it is provided.
	if rTag.MaxTTL > time.Duration(0) {
		value = fmt.Sprintf("%s:t=%s", value, rTag.MaxTTL)
	}

	return value, nil
}

// Parses the tag from string form into a struct form. This method
// also verifies the correctness of the parsed role tag.
func parseAndVerifyRoleTagValue(s logical.Storage, tag string) (*roleTag, error) {
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
		case strings.Contains(tagItem, "a="):
			rTag.AmiID = strings.TrimPrefix(tagItem, "a=")
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
			rTag.MaxTTL, err = time.ParseDuration(strings.TrimPrefix(tagItem, "t="))
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unrecognized item %s in tag", tagItem)
		}
	}

	if rTag.AmiID == "" {
		return nil, fmt.Errorf("missing image ID")
	}

	imageEntry, err := awsImage(s, rTag.AmiID)
	if err != nil {
		return nil, err
	}
	if imageEntry == nil {
		return nil, fmt.Errorf("entry not found for AMI %s", rTag.AmiID)
	}

	// Create a HMAC of the plaintext value of role tag and compare it with the given value.
	verified, err := verifyRoleTagValue(rTag, imageEntry)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, fmt.Errorf("role tag signature verification failed")
	}

	return rTag, nil
}

// Creates base64 encoded HMAC using a backend specific key.
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
	AmiID                    string        `json:"ami_id" structs:"ami_id" mapstructure:"ami_id"`
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
		rTag1.AmiID == rTag2.AmiID &&
		rTag1.HMAC == rTag2.HMAC &&
		rTag1.InstanceID == rTag2.InstanceID &&
		rTag1.DisallowReauthentication == rTag2.DisallowReauthentication &&
		rTag1.AllowInstanceMigration == rTag2.AllowInstanceMigration
}

const pathImageTagSyn = `
Create a tag for an EC2 instance.
`

const pathImageTagDesc = `
When an AMI is used by more than one EC2 instance and there is a need
to apply only a subset of AMI's policies on the instance, create a 
role tag using this endpoint and apply it on the instance.

A RoleTag setting needs to be enabled in 'image/<ami_id>' endpoint, to be able
to create a tag. Also, the policies to be associated with the tag should be
a subset of the policies associated with the regisred AMI.

This endpoint will return both the 'key' and the 'value' to be set for the
EC2 instance tag.
`
