package awsauth

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// awsStsEntry is used to store details of an STS role for assumption
type awsStsEntry struct {
	StsRole string `json:"sts_role"`
}

func pathListSts(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/sts/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathStsList,
		},

		HelpSynopsis:    pathListStsHelpSyn,
		HelpDescription: pathListStsHelpDesc,
	}
}

func pathConfigSts(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/sts/" + framework.GenericNameRegex("account_id"),
		Fields: map[string]*framework.FieldSchema{
			"account_id": {
				Type: framework.TypeString,
				Description: `AWS account ID to be associated with STS role. If set,
Vault will use assumed credentials to verify any login attempts from EC2
instances in this account.`,
			},
			"sts_role": {
				Type: framework.TypeString,
				Description: `AWS ARN for STS role to be assumed when interacting with the account specified.
The Vault server must have permissions to assume this role.`,
			},
		},

		ExistenceCheck: b.pathConfigStsExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigStsCreateUpdate,
			logical.UpdateOperation: b.pathConfigStsCreateUpdate,
			logical.ReadOperation:   b.pathConfigStsRead,
			logical.DeleteOperation: b.pathConfigStsDelete,
		},

		HelpSynopsis:    pathConfigStsSyn,
		HelpDescription: pathConfigStsDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigStsExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	accountID := data.Get("account_id").(string)
	if accountID == "" {
		return false, fmt.Errorf("missing account_id")
	}

	entry, err := b.lockedAwsStsEntry(ctx, req.Storage, accountID)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

// pathStsList is used to list all the AWS STS role configurations
func (b *backend) pathStsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()
	sts, err := req.Storage.List(ctx, "config/sts/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(sts), nil
}

// nonLockedSetAwsStsEntry creates or updates an STS role association with the given accountID
// This method does not acquire the write lock before creating or updating. If locking is
// desired, use lockedSetAwsStsEntry instead
func (b *backend) nonLockedSetAwsStsEntry(ctx context.Context, s logical.Storage, accountID string, stsEntry *awsStsEntry) error {
	if accountID == "" {
		return fmt.Errorf("missing AWS account ID")
	}

	if stsEntry == nil {
		return fmt.Errorf("missing AWS STS Role ARN")
	}

	entry, err := logical.StorageEntryJSON("config/sts/"+accountID, stsEntry)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("failed to create storage entry for AWS STS configuration")
	}

	return s.Put(ctx, entry)
}

// lockedSetAwsStsEntry creates or updates an STS role association with the given accountID
// This method acquires the write lock before creating or updating the STS entry.
func (b *backend) lockedSetAwsStsEntry(ctx context.Context, s logical.Storage, accountID string, stsEntry *awsStsEntry) error {
	if accountID == "" {
		return fmt.Errorf("missing AWS account ID")
	}

	if stsEntry == nil {
		return fmt.Errorf("missing sts entry")
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return b.nonLockedSetAwsStsEntry(ctx, s, accountID, stsEntry)
}

// nonLockedAwsStsEntry returns the STS role associated with the given accountID.
// This method does not acquire the read lock before returning information. If locking is
// desired, use lockedAwsStsEntry instead
func (b *backend) nonLockedAwsStsEntry(ctx context.Context, s logical.Storage, accountID string) (*awsStsEntry, error) {
	entry, err := s.Get(ctx, "config/sts/"+accountID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	var stsEntry awsStsEntry
	if err := entry.DecodeJSON(&stsEntry); err != nil {
		return nil, err
	}

	return &stsEntry, nil
}

// lockedAwsStsEntry returns the STS role associated with the given accountID.
// This method acquires the read lock before returning the association.
func (b *backend) lockedAwsStsEntry(ctx context.Context, s logical.Storage, accountID string) (*awsStsEntry, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedAwsStsEntry(ctx, s, accountID)
}

// pathConfigStsRead is used to return information about an STS role/AWS accountID association
func (b *backend) pathConfigStsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	accountID := data.Get("account_id").(string)
	if accountID == "" {
		return logical.ErrorResponse("missing account id"), nil
	}

	stsEntry, err := b.lockedAwsStsEntry(ctx, req.Storage, accountID)
	if err != nil {
		return nil, err
	}
	if stsEntry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"sts_role": stsEntry.StsRole,
		},
	}, nil
}

// pathConfigStsCreateUpdate is used to associate an STS role with a given AWS accountID
func (b *backend) pathConfigStsCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	accountID := data.Get("account_id").(string)
	if accountID == "" {
		return logical.ErrorResponse("missing AWS account ID"), nil
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// Check if an STS role is already registered
	stsEntry, err := b.nonLockedAwsStsEntry(ctx, req.Storage, accountID)
	if err != nil {
		return nil, err
	}
	if stsEntry == nil {
		stsEntry = &awsStsEntry{}
	}

	// Check that an STS role has actually been provided
	stsRole, ok := data.GetOk("sts_role")
	if ok {
		stsEntry.StsRole = stsRole.(string)
	} else if req.Operation == logical.CreateOperation {
		return logical.ErrorResponse("missing sts role"), nil
	}

	if stsEntry.StsRole == "" {
		return logical.ErrorResponse("sts role cannot be empty"), nil
	}

	// save the provided STS role
	if err := b.nonLockedSetAwsStsEntry(ctx, req.Storage, accountID, stsEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathConfigStsDelete is used to delete a previously configured STS configuration
func (b *backend) pathConfigStsDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	accountID := data.Get("account_id").(string)
	if accountID == "" {
		return logical.ErrorResponse("missing account id"), nil
	}

	return nil, req.Storage.Delete(ctx, "config/sts/"+accountID)
}

const pathConfigStsSyn = `
Specify STS roles to be assumed for certain AWS accounts.
`

const pathConfigStsDesc = `
Allows the explicit association of STS roles to satellite AWS accounts (i.e. those
which are not the account in which the Vault server is running.) Login attempts from
EC2 instances running in these accounts will be verified using credentials obtained
by assumption of these STS roles.

The environment in which the Vault server resides must have access to assume the
given STS roles.
`
const pathListStsHelpSyn = `
List all the AWS account/STS role relationships registered with Vault.
`

const pathListStsHelpDesc = `
AWS accounts will be listed by account ID, along with their respective role names.
`
