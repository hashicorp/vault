package plugin

import (
	"context"
	"errors"
	"github.com/hashicorp/vault-plugin-secrets-ad/plugin/util"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathRotateCredentials() *framework.Path {
	return &framework.Path{
		Pattern: "rotate-root",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathRotateCredentialsUpdate,
		},

		HelpSynopsis:    pathRotateCredentialsUpdateHelpSyn,
		HelpDescription: pathRotateCredentialsUpdateHelpDesc,
	}
}

func (b *backend) pathRotateCredentialsUpdate(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	engineConf, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if engineConf == nil {
		return nil, errors.New("the config is currently unset")
	}

	newPassword, err := util.GeneratePassword(engineConf.PasswordConf.Formatter, engineConf.PasswordConf.Length)
	if err != nil {
		return nil, err
	}

	if err := b.client.UpdateRootPassword(engineConf.ADConf, engineConf.ADConf.BindDN, newPassword); err != nil {
		return nil, err
	}

	engineConf.ADConf.BindPassword = newPassword
	entry, err := logical.StorageEntryJSON(configStorageKey, engineConf)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Respond with a 204.
	return nil, nil
}

const pathRotateCredentialsUpdateHelpSyn = `
Request to rotate the root credentials.
`

const pathRotateCredentialsUpdateHelpDesc = `
This path attempts to rotate the root credentials. 
`
