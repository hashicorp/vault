package awsec2

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type awsStsEntry struct {
	StsRole string `json:"sts_role" structs:"sts_role" mapstructure:"sts_role"`
}

func pathListSts(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/sts/?",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathStsList,
		},
	}
}

func pathConfigSts(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/sts/" + framework.GenericNameRegex("account_id"),
		Fields: map[string]*framework.FieldSchema{
			"account_id": {
				Type:        framework.TypeString,
				Description: "AWS account ID for account for which STS role will be assumed",
			},
			"sts_role": {
				Type:        framework.TypeString,
				Description: "AWS ARN for STS role to be assumed",
			},
		},

		ExistenceCheck: b.pathConfigStsExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigStsCreateUpdate,
			logical.UpdateOperation: b.pathConfigStsCreateUpdate,
			logical.ReadOperation:   b.pathConfigStsRead,
		},
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigStsExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	accountId := data.Get("account_id").(string)
	if accountId == "" {
		return false, fmt.Errorf("missing account_id")
	}

	entry, err := b.lockedAwsStsEntry(req.Storage, accountId)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

// pathStsList is used to list all the AWS STS role configurations
func (b *backend) pathStsList(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()
	sts, err := req.Storage.List("config/sts/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(sts), nil
}

func (b *backend) nonLockedSetAwsStsEntry(s logical.Storage, accountID string, stsEntry *awsStsEntry) error {
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

	return s.Put(entry)
}

func (b *backend) lockedSetAwsStsEntry(s logical.Storage, accountID string, stsEntry *awsStsEntry) error {
	if accountID == "" {
		return fmt.Errorf("missing AWS account ID")
	}

	if stsEntry == nil {
		return fmt.Errorf("missing AWS STS Role ARN")
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	return b.nonLockedSetAwsStsEntry(s, accountID, stsEntry)
}

func (b *backend) nonLockedAwsStsEntry(s logical.Storage, accountID string) (*awsStsEntry, error) {
	entry, err := s.Get("config/sts/" + accountID)
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

func (b *backend) lockedAwsStsEntry(s logical.Storage, accountID string) (*awsStsEntry, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedAwsStsEntry(s, accountID)
}

func (b *backend) pathConfigStsRead(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	accountId := data.Get("account_id").(string)
	if accountId == "" {
		return logical.ErrorResponse("missing account id"), nil
	}

	stsEntry, err := b.lockedAwsStsEntry(req.Storage, accountId)
	if err != nil {
		return nil, err
	}
	if stsEntry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: structs.New(stsEntry).Map(),
	}, nil
}

func (b *backend) pathConfigStsCreateUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	accountId := data.Get("account_id").(string)
	if accountId == "" {
		return logical.ErrorResponse("missing AWS account ID"), nil
	}

	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	// Check if an STS role is already registered
	stsEntry, err := b.nonLockedAwsStsEntry(req.Storage, accountId)
	if err != nil {
		return nil, err
	}
	if stsEntry == nil {
		stsEntry = &awsStsEntry{}
	}

	// Check that an STS role has actually been provided
	stsRole, ok := data.GetOk("sts_role")
	if ok {
		if stsRole != "" {
			stsEntry.StsRole = stsRole.(string)
		} else {
			return logical.ErrorResponse("missing sts role"), nil
		}
	} else {
		return logical.ErrorResponse("missing sts role"), nil
	}

	// save the provided STS role
	if err := b.nonLockedSetAwsStsEntry(req.Storage, accountId, stsEntry); err != nil {
		return nil, err
	}

	return nil, nil
}
