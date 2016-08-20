package duo

import (
	"errors"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathDuoConfig() *framework.Path {
	return &framework.Path{
		Pattern: `duo/config`,
		Fields: map[string]*framework.FieldSchema{
			"user_agent": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "User agent to connect to Duo (default \"\")",
			},
			"username_format": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Format string given auth backend username as argument to create Duo username (default '%s')",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: pathDuoConfigWrite,
			logical.ReadOperation:   pathDuoConfigRead,
		},

		HelpSynopsis:    pathDuoConfigHelpSyn,
		HelpDescription: pathDuoConfigHelpDesc,
	}
}

func GetDuoConfig(req *logical.Request) (*DuoConfig, error) {
	var result DuoConfig
	// all config parameters are optional, so path need not exist
	entry, err := req.Storage.Get("duo/config")
	if err == nil && entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, err
		}
	}
	if result.UsernameFormat == "" {
		result.UsernameFormat = "%s"
	}
	return &result, nil
}

func pathDuoConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username_format := d.Get("username_format").(string)
	if !strings.Contains(username_format, "%s") {
		return nil, errors.New("username_format must include username ('%s')")
	}
	entry, err := logical.StorageEntryJSON("duo/config", DuoConfig{
		UsernameFormat: username_format,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func pathDuoConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	config, err := GetDuoConfig(req)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"username_format": config.UsernameFormat,
		},
	}, nil
}

type DuoConfig struct {
	UsernameFormat string `json:"username_format"`
	UserAgent      string `json:"user_agent"`
}

const pathDuoConfigHelpSyn = `
Configure Duo second factor behavior. 
`

const pathDuoConfigHelpDesc = `
This endpoint allows you to configure how the original auth backend username maps to
the Duo username by providing a template format string.
`
