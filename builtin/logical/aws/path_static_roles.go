package aws

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	pathStaticRole = "static-roles"

	paramRoleName       = "name"
	paramUsername       = "username"
	paramRotationPeriod = "rotation_period"
)

type staticRoleConfig struct {
	Name           string        `json:"name" structs:"name" mapstructure:"name"`
	Username       string        `json:"username" structs:"username" mapstructure:"username"`
	RotationPeriod time.Duration `json:"rotation_period" structs:"rotation_period" mapstructure:"rotation_period"`
}

func pathStaticRoles(b *backend) *framework.Path {
	roleResponse := map[int][]framework.Response{
		http.StatusOK: {{
			Description: http.StatusText(http.StatusOK),
			Fields: map[string]*framework.FieldSchema{
				paramRoleName: {
					Type:        framework.TypeString,
					Description: descRoleName,
				},
				paramUsername: {
					Type:        framework.TypeString,
					Description: descUsername,
				},
				paramRotationPeriod: {
					Type:        framework.TypeDurationSecond,
					Description: descRotationPeriod,
				},
			},
		}},
	}

	return &framework.Path{
		Pattern: fmt.Sprintf("%s/%s", pathStaticRole, framework.GenericNameWithAtRegex(paramRoleName)),
		Fields: map[string]*framework.FieldSchema{
			paramRoleName: {
				Type:        framework.TypeString,
				Description: descRoleName,
			},
			paramUsername: {
				Type:        framework.TypeString,
				Description: descUsername,
			},
			paramRotationPeriod: {
				Type:        framework.TypeDurationSecond,
				Description: descRotationPeriod,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:  b.pathStaticRolesRead,
				Responses: roleResponse,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathStaticRolesWrite,
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
				Responses:                   roleResponse,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback:                    b.pathStaticRolesDelete,
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
				Responses: map[int][]framework.Response{
					http.StatusNoContent: {{
						Description: http.StatusText(http.StatusNoContent),
					}},
				},
			},
		},

		HelpSynopsis:    pathStaticRolesHelpSyn,
		HelpDescription: pathStaticRolesHelpDesc,
	}
}

func (b *backend) pathStaticRolesRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName, ok := data.GetOk(paramRoleName)
	if !ok {
		return nil, fmt.Errorf("missing '%s' parameter", paramRoleName)
	}

	entry, err := req.Storage.Get(ctx, formatRoleStoragePath(roleName.(string)))
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration for static role '%q': %w", roleName, err)
	}
	if entry == nil {
		return nil, nil
	}

	var config staticRoleConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("failed to decode configuration for tatic role '%q': %w", roleName, err)
	}

	return &logical.Response{
		Data: formatResponse(config),
	}, nil
}

func (b *backend) pathStaticRolesWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Create & validate config from request parameters
	config := staticRoleConfig{}

	if rawRoleName, ok := data.GetOk(paramRoleName); ok {
		config.Name = rawRoleName.(string)

		if err := b.validateRoleName(config.Name); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("missing '%s' parameter", paramRoleName)
	}

	if rawUsername, ok := data.GetOk(paramUsername); ok {
		config.Username = rawUsername.(string)

		if err := b.validateIAMUserExists(ctx, req, config.Username); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("missing '%s' parameter", paramUsername)
	}

	if rawRotationPeriod, ok := data.GetOk(paramRotationPeriod); ok {
		config.RotationPeriod = time.Duration(rawRotationPeriod.(int)) * time.Second

		if err := b.validateRotationPeriod(config.RotationPeriod); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("missing '%s' parameter", paramRotationPeriod)
	}

	// Upsert role config
	newRole, err := logical.StorageEntryJSON(formatRoleStoragePath(config.Name), config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object to JSON: %w", err)
	}
	err = req.Storage.Put(ctx, newRole)
	if err != nil {
		return nil, fmt.Errorf("failed to save object in storage: %w", err)
	}

	// Bootstrap initial set of keys if they did not exist before. AWS Secret Access Keys can only be obtained on creation,
	// so we need to boostrap new roles with a new initial set of keys to be able to serve valid credentials to Vault clients.
	existingCreds, err := req.Storage.Get(ctx, formatCredsStoragePath(config.Name))
	if err != nil {
		return nil, fmt.Errorf("unable to verify if credentials already exist for role '%q': %w", config.Name, err)
	}
	if existingCreds == nil {
		err := b.createCredential(ctx, req.Storage, config)
		if err != nil {
			return nil, fmt.Errorf("failed to create new credentials for role '%q': %w", config.Name, err)
		}

		err = b.credRotationQueue.Push(&queue.Item{
			Key:      config.Name,
			Value:    config,
			Priority: time.Now().Add(config.RotationPeriod).Unix(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add item into the rotation queue for role '%q': %w", config.Name, err)
		}
	}

	return &logical.Response{
		Data: formatResponse(config),
	}, nil
}

func (b *backend) pathStaticRolesDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName, ok := data.GetOk(paramRoleName)
	if !ok {
		return nil, fmt.Errorf("missing '%s' parameter", paramRoleName)
	}

	err := b.deleteCredential(ctx, req.Storage, roleName.(string))
	if err != nil {
		return nil, fmt.Errorf("failed to clean credentials while deleting role '%q': %w", roleName.(string), err)
	}

	return nil, req.Storage.Delete(ctx, formatRoleStoragePath(roleName.(string)))
}

func (b *backend) validateRoleName(name string) error {
	return nil
}

func (b *backend) validateIAMUserExists(ctx context.Context, req *logical.Request, username string) error {
	return nil
}

func (b *backend) validateRotationPeriod(period time.Duration) error {
	return nil
}

func formatResponse(cfg staticRoleConfig) map[string]interface{} {
	response := structs.New(cfg).Map()
	response[paramRotationPeriod] = int64(cfg.RotationPeriod.Seconds())

	return response
}

func formatRoleStoragePath(roleName string) string {
	return fmt.Sprintf("%s/%s", pathStaticRole, roleName)
}

const pathStaticRolesHelpSyn = `
`

const pathStaticRolesHelpDesc = `
`

const (
	descRoleName       = ""
	descUsername       = ""
	descRotationPeriod = ""
)
