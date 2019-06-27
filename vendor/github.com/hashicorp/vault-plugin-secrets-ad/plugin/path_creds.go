package plugin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/hashicorp/vault-plugin-secrets-ad/plugin/util"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	credPrefix = "creds/"
	storageKey = "creds"

	// Since password TTL can be set to as low as 1 second,
	// we can't cache passwords for an entire second.
	credCacheCleanup    = time.Second / 3
	credCacheExpiration = time.Second / 2
)

// deleteCred fulfills the DeleteWatcher interface in roles.
// It allows the roleHandler to let us know when a role's been deleted so we can delete its associated creds too.
func (b *backend) deleteCred(ctx context.Context, storage logical.Storage, roleName string) error {
	if err := storage.Delete(ctx, storageKey+"/"+roleName); err != nil {
		return err
	}
	b.credCache.Delete(roleName)
	return nil
}

func (b *backend) invalidateCred(ctx context.Context, key string) {
	if strings.HasPrefix(key, credPrefix) {
		roleName := key[len(credPrefix):]
		b.credCache.Delete(roleName)
	}
}

func (b *backend) pathCreds() *framework.Path {
	return &framework.Path{
		Pattern: credPrefix + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.credReadOperation,
		},
		HelpSynopsis:    credHelpSynopsis,
		HelpDescription: credHelpDescription,
	}
}

func (b *backend) credReadOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	cred := make(map[string]interface{})

	engineConf, err := b.readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if engineConf == nil {
		return nil, errors.New("the config is currently unset")
	}

	roleName := fieldData.Get("name").(string)

	// We act upon quite a few things below that could be racy if not locked:
	// 		- Roles. If a new cred is created, the role is updated to include the new LastVaultRotation time,
	//		  effecting role storage (and the role cache, but that's already thread-safe).
	//		- Creds. New creds involve writing to cred storage and the cred cache (also already thread-safe).
	// Rather than setting read locks of different types, and upgrading them to write locks, let's keep complexity
	// low and use one simple mutex.
	b.credLock.Lock()
	defer b.credLock.Unlock()

	role, err := b.readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	b.Logger().Debug(fmt.Sprintf("role is: %+v", role))

	var resp *logical.Response
	var respErr error
	var unset time.Time

	switch {

	case role.LastVaultRotation == unset:
		b.Logger().Info("rotating password for the first time so Vault will know it")
		resp, respErr = b.generateAndReturnCreds(ctx, engineConf, req.Storage, roleName, role, cred)

	case role.PasswordLastSet.After(role.LastVaultRotation.Add(time.Second * time.Duration(engineConf.LastRotationTolerance))):
		b.Logger().Warn(fmt.Sprintf(
			"Vault rotated the password at %s, but it was rotated in AD later at %s, so rotating it again so Vault will know it",
			role.LastVaultRotation.String(), role.PasswordLastSet.String()),
		)
		resp, respErr = b.generateAndReturnCreds(ctx, engineConf, req.Storage, roleName, role, cred)

	default:
		b.Logger().Debug("determining whether to rotate credential")
		credIfc, found := b.credCache.Get(roleName)
		if found {
			b.Logger().Debug("checking cached credential")
			cred = credIfc.(map[string]interface{})
		} else {
			b.Logger().Debug("checking stored credential")
			entry, err := req.Storage.Get(ctx, storageKey+"/"+roleName)
			if err != nil {
				return nil, err
			}
			if entry == nil {
				// If the creds aren't in storage, but roles are and we've created creds before,
				// this is an unexpected state and something has gone wrong.
				// Let's be explicit and error about this.
				return nil, fmt.Errorf("should have the creds for %+v but they're not found", role)
			}
			if err := entry.DecodeJSON(&cred); err != nil {
				return nil, err
			}
			b.credCache.SetDefault(roleName, cred)
		}

		now := time.Now().UTC()
		shouldBeRolled := role.LastVaultRotation.Add(time.Duration(role.TTL) * time.Second) // already in UTC
		if now.After(shouldBeRolled) {
			b.Logger().Info(fmt.Sprintf(
				"last Vault rotation was at %s, and since the TTL is %d and it's now %s, it's time to rotate it",
				role.LastVaultRotation.String(), role.TTL, now.String()),
			)
			resp, respErr = b.generateAndReturnCreds(ctx, engineConf, req.Storage, roleName, role, cred)
		} else {
			b.Logger().Debug("returning previous credential")
			resp = &logical.Response{
				Data: cred,
			}
		}
	}
	if respErr != nil {
		return nil, respErr
	}
	return resp, nil
}

func (b *backend) generateAndReturnCreds(ctx context.Context, engineConf *configuration, storage logical.Storage, roleName string, role *backendRole, previousCred map[string]interface{}) (*logical.Response, error) {
	newPassword, err := util.GeneratePassword(engineConf.PasswordConf.Formatter, engineConf.PasswordConf.Length)
	if err != nil {
		return nil, err
	}

	if err := b.client.UpdatePassword(engineConf.ADConf, role.ServiceAccountName, newPassword); err != nil {
		return nil, err
	}

	// Time recorded is in UTC for easier user comparison to AD's last rotated time, which is set to UTC by Microsoft.
	role.LastVaultRotation = time.Now().UTC()
	if err := b.writeRoleToStorage(ctx, storage, roleName, role); err != nil {
		return nil, err
	}
	// Cache the full role to minimize Vault storage calls.
	b.roleCache.SetDefault(roleName, role)

	// Although a service account name is typically my_app@example.com,
	// the username it uses is just my_app, or everything before the @.
	var username string
	fields := strings.Split(role.ServiceAccountName, "@")
	if len(fields) > 0 {
		username = fields[0]
	} else {
		return nil, fmt.Errorf("unable to infer username from service account name: %s", role.ServiceAccountName)
	}

	cred := map[string]interface{}{
		"username":         username,
		"current_password": newPassword,
	}
	if previousCred["current_password"] != nil {
		cred["last_password"] = previousCred["current_password"]
	}

	// Cache and save the cred.
	entry, err := logical.StorageEntryJSON(storageKey+"/"+roleName, cred)
	if err != nil {
		return nil, err
	}
	if err := storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	b.credCache.SetDefault(roleName, cred)

	return &logical.Response{
		Data: cred,
	}, nil
}

const (
	credHelpSynopsis = `
Retrieve a role's creds by role name.
`
	credHelpDescription = `
Read creds using a role's name to view the login, current password, and last password.
`
)
