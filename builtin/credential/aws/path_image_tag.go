package aws

import (
	"crypto/hmac"
	"crypto/sha256"
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
		Pattern: "image/" + framework.GenericNameRegex("name") + "/tag$",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "AMI name to create a tag for.",
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Policies to be associated with the tag.",
			},

			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     0,
				Description: "The maximum allowed lease duration",
			},

			"disallow_reauthentication": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, only allows a single token to be granted per instance ID. This can be cleared with the auth/aws/whitelist/identity endpoint.",
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

	imageID := strings.ToLower(data.Get("name").(string))
	if imageID == "" {
		return logical.ErrorResponse("missing image name"), nil
	}

	// Parse the given policies into a slice and add 'default' if not provided.
	// Remove all other policies if 'root' is present.
	policies := policyutil.ParsePolicies(data.Get("policies").(string))

	disallowReauthentication := data.Get("disallow_reauthentication").(bool)

	// Fetch the image entry corresponding to the AMI name
	imageEntry, err := awsImage(req.Storage, imageID)
	if err != nil {
		return nil, err
	}
	if imageEntry == nil {
		return logical.ErrorResponse("image entry not found"), nil
	}

	// If RoleTag is empty, disallow creation of tag.
	if imageEntry.RoleTag == "" {
		return logical.ErrorResponse("tag creation is not enabled for this image"), nil
	}

	// Create a random nonce
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

	// Attach version, nonce, policies and maxTTL to the role tag value.
	rTagValue, err := prepareRoleTagPlainValue(&roleTag{Version: roleTagVersion,
		ImageID:  imageID,
		Nonce:    nonce,
		Policies: policies,
		MaxTTL:   maxTTL,
		DisallowReauthentication: disallowReauthentication,
	})
	if err != nil {
		return nil, err
	}

	// Get the key used for creating the HMAC
	key, err := hmacKey(req.Storage)
	if err != nil {
		return nil, err
	}

	// Create the HMAC of the value
	hmacB64, err := createRoleTagHMACBase64(key, rTagValue)
	if err != nil {
		return nil, err
	}

	// attach the HMAC to the value
	rTagValue = fmt.Sprintf("%s:%s", rTagValue, hmacB64)
	if len(rTagValue) > 255 {
		return nil, fmt.Errorf("role tag 'value' exceeding the limit of 255 characters")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"tag_key":   imageEntry.RoleTag,
			"tag_value": rTagValue,
		},
	}, nil
}

// verifyRoleTagValue rebuilds the role tag value without the HMAC,
// computes the HMAC from it using the backend specific key and
// compares it with the received HMAC.
func verifyRoleTagValue(s logical.Storage, rTag *roleTag) (bool, error) {
	// Fetch the plaintext part of role tag
	rTagPlainText, err := prepareRoleTagPlainValue(rTag)
	if err != nil {
		return false, err
	}

	// Get the key used for creating the HMAC
	key, err := hmacKey(s)
	if err != nil {
		return false, err
	}

	// Compute the HMAC of the plaintext
	hmacB64, err := createRoleTagHMACBase64(key, rTagPlainText)
	if err != nil {
		return false, err
	}
	return rTag.HMAC == hmacB64, nil
}

// prepareRoleTagPlainValue builds the role tag value without the HMAC in it.
func prepareRoleTagPlainValue(rTag *roleTag) (string, error) {
	if rTag.Version == "" {
		return "", fmt.Errorf("missing version")
	}
	// attach version to the value
	value := rTag.Version

	if rTag.Nonce == "" {
		return "", fmt.Errorf("missing nonce")
	}
	// attach nonce to the value
	value = fmt.Sprintf("%s:%s", value, rTag.Nonce)

	if rTag.ImageID == "" {
		return "", fmt.Errorf("missing ami_name")
	}
	// attach ami_name to the value
	value = fmt.Sprintf("%s:a=%s", value, rTag.ImageID)

	// attach policies to value
	value = fmt.Sprintf("%s:p=%s", value, strings.Join(rTag.Policies, ","))

	// attach disallow_reauthentication field
	value = fmt.Sprintf("%s:d=%s", value, strconv.FormatBool(rTag.DisallowReauthentication))

	// attach max_ttl if it is provided
	if rTag.MaxTTL > time.Duration(0) {
		value = fmt.Sprintf("%s:t=%s", value, rTag.MaxTTL)
	}

	return value, nil
}

// Parses the tag from string form into a struct form.
func parseRoleTagValue(s logical.Storage, tag string) (*roleTag, error) {
	tagItems := strings.Split(tag, ":")
	// Tag must contain version, nonce, policies and HMAC
	if len(tagItems) < 4 {
		return nil, fmt.Errorf("invalid tag")
	}

	rTag := &roleTag{}

	// Cache the HMAC value. The last item in the collection.
	rTag.HMAC = tagItems[len(tagItems)-1]

	// Delete the HMAC from the list.
	tagItems = tagItems[:len(tagItems)-1]

	// Version is the first element.
	rTag.Version = tagItems[0]
	if rTag.Version != roleTagVersion {
		return nil, fmt.Errorf("invalid role tag version")
	}

	// Nonce is the second element.
	rTag.Nonce = tagItems[1]

	if len(tagItems) > 2 {
		// Delete the version and nonce from the list.
		tagItems = tagItems[2:]
		for _, tagItem := range tagItems {
			var err error
			switch {
			case strings.Contains(tagItem, "a="):
				rTag.ImageID = strings.TrimPrefix(tagItem, "a=")
			case strings.Contains(tagItem, "p="):
				rTag.Policies = strings.Split(strings.TrimPrefix(tagItem, "p="), ",")
			case strings.Contains(tagItem, "d="):
				rTag.DisallowReauthentication, err = strconv.ParseBool(strings.TrimPrefix(tagItem, "d="))
				if err != nil {
					return nil, err
				}
			case strings.Contains(tagItem, "t="):
				rTag.MaxTTL, err = time.ParseDuration(strings.TrimPrefix(tagItem, "t="))
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unrecognized item in tag")
			}
		}
	}
	if rTag.ImageID == "" {
		return nil, fmt.Errorf("missing image ID")
	}

	// Create a HMAC of the plaintext value of role tag and compare it with the given value.
	verified, err := verifyRoleTagValue(s, rTag)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, fmt.Errorf("role tag signature mismatch")
	}
	return rTag, nil
}

// Creates base64 encoded HMAC using a backend specific key.
func createRoleTagHMACBase64(key, value string) (string, error) {
	hm := hmac.New(sha256.New, []byte(key))
	hm.Write([]byte(value))

	// base64 encode the hmac bytes.
	return base64.StdEncoding.EncodeToString(hm.Sum(nil)), nil
}

// Creates a base64 encoded random nonce.
func createRoleTagNonce() (string, error) {
	uuidBytes, err := uuid.GenerateRandomBytes(8)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(uuidBytes), nil
}

// Struct roleTag represents a role tag in a struc form.
type roleTag struct {
	Version                  string        `json:"version" structs:"version" mapstructure:"version"`
	Nonce                    string        `json:"nonce" structs:"nonce" mapstructure:"nonce"`
	Policies                 []string      `json:"policies" structs:"policies" mapstructure:"policies"`
	MaxTTL                   time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	ImageID                  string        `json:"image_id" structs:"image_id" mapstructure:"image_id"`
	HMAC                     string        `json:"hmac" structs:"hmac" mapstructure:"hmac"`
	DisallowReauthentication bool          `json:"disallow_reauthentication" structs:"disallow_reauthentication" mapstructure:"disallow_reauthentication"`
}

func (rTag1 *roleTag) Equal(rTag2 *roleTag) bool {
	return rTag1.Version == rTag2.Version &&
		rTag1.Nonce == rTag2.Nonce &&
		policyutil.EquivalentPolicies(rTag1.Policies, rTag2.Policies) &&
		rTag1.MaxTTL == rTag2.MaxTTL &&
		rTag1.ImageID == rTag2.ImageID &&
		rTag1.HMAC == rTag2.HMAC &&
		rTag1.DisallowReauthentication == rTag2.DisallowReauthentication
}

const pathImageTagSyn = `
Create a tag for an EC2 instance.
`

const pathImageTagDesc = `
When an AMI is used by more than one EC2 instance, policies to be associated
during login are determined by a particular tag on the instance. This tag
can be created using this endpoint.

A RoleTag setting needs to be enabled in 'image/<name>' endpoint, to be able
to create a tag. Also, the policies to be associated with the tag should be
a subset of the policies associated with the regisred AMI.

This endpoint will return both the 'key' and the 'value' to be set for the
instance tag.
`
