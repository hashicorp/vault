package builder

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configStoragePath = "config"
)

var userClientConfig any

// pathConfig extends the Vault API with a `/config`
// endpoint for the backend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. For example, password
// is marked as sensitive and will not be output
// when you read the configuration.
func (gb *GenericBackend[CC, C, R]) pathConfig(inputConfig *ClientConfig[CC, C, R]) *framework.Path {
	// userClientConfig = inputConfig

	return &framework.Path{
		Pattern: "config",
		Fields:  inputConfig.Fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: gb.pathConfigRead,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: gb.pathConfigWrite,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: gb.pathConfigWrite,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: gb.pathConfigDelete,
			},
		},
		ExistenceCheck:  gb.pathConfigExistenceCheck,
		HelpSynopsis:    pathConfigHelpSynopsis,
		HelpDescription: pathConfigHelpDescription,
	}
}

func (gb *GenericBackend[CC, C, R]) getConfig(ctx context.Context, s logical.Storage) (*CC, error) {
	entry, err := s.Get(ctx, configStoragePath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	config := new(CC)
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %w", err)
	}

	// return the config, we are done
	return config, nil
}

// pathConfigExistenceCheck verifies if the configuration exists.
func (gb *GenericBackend[CC, C, R]) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	return out != nil, nil
}

// pathConfigRead reads the configuration and outputs non-sensitive information.
func (gb *GenericBackend[CC, C, R]) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := gb.getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	responseData := structs.Map(config)

	return &logical.Response{
		Data: responseData,
	}, nil
}

// pathConfigWrite updates the configuration for the backend
func (gb *GenericBackend[CC, C, R]) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := gb.getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	createOperation := (req.Operation == logical.CreateOperation)
	if config == nil {
		if !createOperation {
			return nil, errors.New("config not found during update operation")
		}
		config = new(CC)
	}

	writeData := structs.Map(config)

	for k := range writeData {
		if userInput, ok := data.GetOk(k); ok {
			writeData[k] = userInput
		} else if createOperation {
			return nil, fmt.Errorf("missing %s in configuration", k)
		}
	}

	entry, err := logical.StorageEntryJSON(configStoragePath, writeData)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// reset the client so the next invocation will pick up the new configuration
	gb.reset()

	return nil, nil
}

// pathConfigDelete removes the configuration for the backend
func (gb *GenericBackend[CC, C, R]) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, configStoragePath)

	if err == nil {
		gb.reset()
	}

	return nil, err
}

//func (gb *GenericBackend[CC, C, R]) validateConfig(writeData map[string]any) error {
//	fmt.Println("made it")
//	inputConfig := userClientConfig.(ClientConfig[CC, C, R])
//	result := new(CC)
//	err := mapstructure.Decode(writeData, result)
//	if err != nil {
//		return err
//	}
//
//	return inputConfig.ValidateFunc(result)
//}

// pathConfigHelpSynopsis summarizes the help text for the configuration
const pathConfigHelpSynopsis = `Configure the HashiCups backend.`

// pathConfigHelpDescription describes the help text for the configuration
const pathConfigHelpDescription = `
The HashiCups secret backend requires credentials for managing
JWTs issued to users working with the products API.

You must sign up with a username and password and
specify the HashiCups address for the products API
before using this secrets backend.
`
