package aws

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathImage(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "image/" + framework.GenericNameRegex("ami_id"),
		Fields: map[string]*framework.FieldSchema{
			"ami_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "AMI ID to be mapped.",
			},

			"role_tag": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "If set, enables the RoleTag for this AMI. The value set for this field should be the 'key' of the tag on the EC2 instance using the RoleTag. Defaults to empty string.",
			},

			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     0,
				Description: "The maximum allowed lease duration",
			},

			"policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "default",
				Description: "Policies to be associated with the AMI.",
			},

			"allow_instance_migration": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, allows migration of the underlying instance where the client resides. This keys off of pendingTime in the metadata document, so essentially, this disables the client nonce check whenever the instance is migrated to a new host and pendingTime is newer than the previously-remembered time. Use with caution.",
			},

			"disallow_reauthentication": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "If set, only allows a single token to be granted per instance ID. This can be cleared with the auth/aws/whitelist/identity endpoint.",
			},
		},

		ExistenceCheck: b.pathImageExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathImageCreateUpdate,
			logical.UpdateOperation: b.pathImageCreateUpdate,
			logical.ReadOperation:   b.pathImageRead,
			logical.DeleteOperation: b.pathImageDelete,
		},

		HelpSynopsis:    pathImageSyn,
		HelpDescription: pathImageDesc,
	}
}

// pathListImages createa a path that enables listing of all the AMIs that are
// registered with Vault.
func pathListImages(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "images/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathImageList,
		},

		HelpSynopsis:    pathListImagesHelpSyn,
		HelpDescription: pathListImagesHelpDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathImageExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := awsImage(req.Storage, strings.ToLower(data.Get("ami_id").(string)))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// awsImage is used to get the information registered for the given AMI ID.
func awsImage(s logical.Storage, amiID string) (*awsImageEntry, error) {
	entry, err := s.Get("image/" + strings.ToLower(amiID))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result awsImageEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// pathImageDelete is used to delete the information registered for a given AMI ID.
func (b *backend) pathImageDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("image/" + strings.ToLower(data.Get("ami_id").(string)))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// pathImageList is used to list all the AMI IDs registered with Vault.
func (b *backend) pathImageList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	images, err := req.Storage.List("image/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(images), nil
}

// pathImageRead is used to view the information registered for a given AMI ID.
func (b *backend) pathImageRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	imageEntry, err := awsImage(req.Storage, strings.ToLower(data.Get("ami_id").(string)))
	if err != nil {
		return nil, err
	}
	if imageEntry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"role_tag":                  imageEntry.RoleTag,
			"policies":                  imageEntry.Policies,
			"max_ttl":                   imageEntry.MaxTTL / time.Second,
			"allow_instance_migration":  imageEntry.AllowInstanceMigration,
			"disallow_reauthentication": imageEntry.DisallowReauthentication,
		},
	}, nil
}

// pathImageCreateUpdate is used to associate Vault policies to a given AMI ID.
func (b *backend) pathImageCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	imageID := strings.ToLower(data.Get("ami_id").(string))
	if imageID == "" {
		return logical.ErrorResponse("missing ami_id"), nil
	}

	imageEntry, err := awsImage(req.Storage, imageID)
	if err != nil {
		return nil, err
	}
	if imageEntry == nil {
		imageEntry = &awsImageEntry{}
	}

	policiesStr, ok := data.GetOk("policies")
	if ok {
		imageEntry.Policies = policyutil.ParsePolicies(policiesStr.(string))
	} else if req.Operation == logical.CreateOperation {
		imageEntry.Policies = []string{"default"}
	}

	disallowReauthenticationBool, ok := data.GetOk("disallow_reauthentication")
	if ok {
		imageEntry.DisallowReauthentication = disallowReauthenticationBool.(bool)
	} else if req.Operation == logical.CreateOperation {
		imageEntry.DisallowReauthentication = data.Get("disallow_reauthentication").(bool)
	}

	allowInstanceMigrationBool, ok := data.GetOk("allow_instance_migration")
	if ok {
		imageEntry.AllowInstanceMigration = allowInstanceMigrationBool.(bool)
	} else if req.Operation == logical.CreateOperation {
		imageEntry.AllowInstanceMigration = data.Get("allow_instance_migration").(bool)
	}

	maxTTLInt, ok := data.GetOk("max_ttl")
	if ok {
		maxTTL := time.Duration(maxTTLInt.(int)) * time.Second
		systemMaxTTL := b.System().MaxLeaseTTL()
		if maxTTL > systemMaxTTL {
			return logical.ErrorResponse(fmt.Sprintf("Given TTL of %d seconds greater than current mount/system default of %d seconds", maxTTL/time.Second, systemMaxTTL/time.Second)), nil
		}

		if maxTTL < time.Duration(0) {
			return logical.ErrorResponse("max_ttl cannot be negative"), nil
		}

		imageEntry.MaxTTL = maxTTL
	} else if req.Operation == logical.CreateOperation {
		imageEntry.MaxTTL = time.Duration(data.Get("max_ttl").(int)) * time.Second
	}

	roleTagStr, ok := data.GetOk("role_tag")
	if ok {
		imageEntry.RoleTag = roleTagStr.(string)
		if len(imageEntry.RoleTag) > 127 {
			return logical.ErrorResponse("role tag 'key' is exceeding the limit of 127 characters"), nil
		}
	} else if req.Operation == logical.CreateOperation {
		imageEntry.RoleTag = data.Get("role_tag").(string)
	}

	entry, err := logical.StorageEntryJSON("image/"+imageID, imageEntry)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}
	return nil, nil
}

// Struct to hold the information associated with an AMI ID in Vault.
type awsImageEntry struct {
	RoleTag                  string        `json:"role_tag" structs:"role_tag" mapstructure:"role_tag"`
	AllowInstanceMigration   bool          `json:"allow_instance_migration" structs:"allow_instance_migration" mapstructure:"allow_instance_migration"`
	MaxTTL                   time.Duration `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	Policies                 []string      `json:"policies" structs:"policies" mapstructure:"policies"`
	DisallowReauthentication bool          `json:"disallow_reauthentication" structs:"disallow_reauthentication" mapstructure:"disallow_reauthentication"`
}

const pathImageSyn = `
Associate an AMI to Vault's policies.
`

const pathImageDesc = `
A precondition for login is that the AMI used by the EC2 instance, needs to
be registered with Vault. After the authentication of the instance, the
authorization for the instance to access Vault's resources is determined
by the policies that are associated to the AMI through this endpoint.

In case the AMI is shared by many instances, then a role tag can be created
through the endpoint 'image/<ami_id>/tag'. This tag needs to be applied on the
instance before it attempts to login to Vault. The policies on the tag should
be a subset of policies that are associated to the AMI in this endpoint. In
order to enable login using tags, RoleTag needs to be enabled in this endpoint.

Also, a 'max_ttl' can be configured in this endpoint that determines the maximum
duration for which a login can be renewed. Note that the 'max_ttl' has a upper
limit of the 'max_ttl' value that is applicable to the backend.
`

const pathListImagesHelpSyn = `
Lists all the AMIs that are registered with Vault.
`

const pathListImagesHelpDesc = `
AMIs will be listed by their respective AMI ID.
`
