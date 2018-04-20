package creds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/builtin/logical/ad/roles"
	"github.com/hashicorp/vault/builtin/logical/ad/util"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/patrickmn/go-cache"
)

const (
	InboundPath = "creds"
	storageKey  = "creds"

	// Since Active Directory offers eventual consistency, in testing we found that sometimes
	// Active Directory returned "password last set" times that were _later_ than our own,
	// even though ours were captured after synchronously completing a password update operation.
	//
	// An example we captured was:
	// 		last_vault_rotation     2018-04-18T22:29:57.385454779Z
	// 		password_last_set       2018-04-18T22:29:57.3902786Z
	//
	// Thus we add a short time buffer when checking whether anyone _else_ updated the AD password
	// since Vault last rotated it.
	passwordLastSetBuffer = time.Second

	// Since password TTL can be set to as low as 1 second,
	// we can't cache passwords for an entire second.
	cacheCleanup    = time.Second / 3
	cacheExpiration = time.Second / 2
)

func Handler(logger hclog.Logger, config config.Reader, roleRW roles.ReadWriter) *handler {
	return &handler{
		logger: logger,
		config: config,
		roleRW: roleRW,
		cache:  cache.New(cacheExpiration, cacheCleanup),
	}
}

type handler struct {
	logger hclog.Logger
	config config.Reader
	roleRW roles.ReadWriter
	cache  *cache.Cache
}

// Delete fulfills the DeleteWatcher interface in roles.
// It allows the roleHandler to let us know when a role's been deleted so we can delete its associated creds too.
func (h *handler) Delete(ctx context.Context, storage logical.Storage, roleName string) error {
	if err := storage.Delete(ctx, storageKey+"/"+roleName); err != nil {
		return err
	}
	h.cache.Delete(roleName)
	return nil
}

func (h *handler) Invalidate(ctx context.Context, key string) {
	prefix := InboundPath + "/"
	if strings.HasPrefix(key, prefix) {
		roleName, err := util.ParseRoleName(prefix, key)
		if err != nil {
			// The role name is invalid, so it's not in storage anyways.
			// Only roles with valid names can get put in storage.
			return
		}
		h.cache.Delete(roleName)
	}
}

func (h *handler) Path() *framework.Path {
	return &framework.Path{
		Pattern: "^creds/.+$",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: h.readOperation,
		},
		HelpSynopsis:    helpSynopsis,
		HelpDescription: helpDescription,
	}
}

func (h *handler) readOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	resp, err := h.readOperationLogic(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	resp.AddWarning("Read access to this endpoint should be controlled via ACLs as it will return the creds information as-is, including any passwords.")
	return resp, nil
}

func (h *handler) readOperationLogic(ctx context.Context, req *logical.Request) (*logical.Response, error) {

	cred := &credential{}

	roleName, err := util.ParseRoleName(InboundPath+"/", req.Path)
	if err != nil {
		return nil, err
	}

	role, err := h.roleRW.Read(ctx, req.Storage, roleName)
	if err != nil {
		if _, ok := err.(*roles.NotFound); ok {
			// The role was deleted, and we've also already deleted the cred.
			return nil, nil
		}
		return nil, err
	}

	// Have we ever managed this cred before?
	// If not, we need to rotate the password so Vault will know it.
	var unset time.Time
	if role.LastVaultRotation == unset {
		return h.generateAndReturnCreds(ctx, req.Storage, roleName, role, cred)
	}

	// Has anyone manually rotated the password in Active Directory?
	// If so, we need to rotate it now so Vault will know it.
	if role.PasswordLastSet.After(role.LastVaultRotation.Add(passwordLastSetBuffer)) {
		return h.generateAndReturnCreds(ctx, req.Storage, roleName, role, cred)
	}

	// Since we should know the last password, let's retrieve it now so we can return it with the new one.
	credIfc, found := h.cache.Get(roleName)
	if found {
		cred = credIfc.(*credential)
	} else {
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
		if err := entry.DecodeJSON(cred); err != nil {
			return nil, err
		}
		h.cache.SetDefault(roleName, cred)
	}

	// Is the password too old?
	// If so, time for a new one!
	now := time.Now().UTC()
	shouldBeRolled := role.LastVaultRotation.Add(time.Duration(role.TTL) * time.Second) // already in UTC
	if now.After(shouldBeRolled) {
		return h.generateAndReturnCreds(ctx, req.Storage, roleName, role, cred)
	}

	// Current credential is accurate! Return it.
	return &logical.Response{
		Data: cred.Map(),
	}, nil
}

func (h *handler) generateAndReturnCreds(ctx context.Context, storage logical.Storage, roleName string, role *roles.Role, previousCred *credential) (*logical.Response, error) {

	engineConf, err := h.config.Read(ctx, storage)
	if err != nil {
		return nil, err
	}

	newPassword, err := util.GeneratePassword(engineConf.PasswordConf.Length)
	if err != nil {
		return nil, err
	}

	secretsClient := util.NewSecretsClient(h.logger, engineConf.ADConf)
	if err := secretsClient.UpdatePassword(role.ServiceAccountName, newPassword); err != nil {
		return nil, err
	}

	// Time recorded is in UTC for easier user comparison to AD's last rotated time, which is set to UTC by Microsoft.
	role.LastVaultRotation = time.Now().UTC()
	if err := h.roleRW.Write(ctx, storage, roleName, role); err != nil {
		return nil, err
	}

	cred := &credential{
		Username:        role.ServiceAccountName,
		CurrentPassword: newPassword,
		LastPassword:    previousCred.CurrentPassword,
	}

	// Cache and save the cred.
	entry, err := logical.StorageEntryJSON(storageKey+"/"+roleName, cred)
	if err != nil {
		return nil, err
	}
	if err := storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	h.cache.SetDefault(roleName, cred)

	return &logical.Response{
		Data: cred.Map(),
	}, nil
}

const (
	helpSynopsis    = ``
	helpDescription = ``
)
