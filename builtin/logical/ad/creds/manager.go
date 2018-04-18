package creds

import (
	"context"
	"time"

	"fmt"

	"crypto/rand"
	"encoding/base64"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/builtin/logical/ad/roles"
	"github.com/hashicorp/vault/builtin/logical/ad/util"
	"github.com/hashicorp/vault/helper/activedirectory"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/patrickmn/go-cache"
)

// TODO need to work through what does and doesn't need to be exported
const (
	BackendPath = "creds"
	StorageKey  = "creds"

	// Since password TTL can be set to as low as 1 second,
	// we can't cache passwords for an entire second.
	cacheCleanup    = time.Second / 3
	cacheExpiration = time.Second / 2
)

func NewManager(logger hclog.Logger, configHandler config.Handler, roleHandler roles.Handler) *Manager {
	return &Manager{
		logger:        logger,
		configHandler: configHandler,
		roleHandler:   roleHandler,
		cache:         cache.New(cacheExpiration, cacheCleanup),
	}
}

type Manager struct {
	logger        hclog.Logger
	configHandler config.Handler
	roleHandler   roles.Handler
	cache         *cache.Cache
}

func (m *Manager) Invalidate(ctx context.Context, key string) {
	if key == BackendPath {
		m.cache.Flush()
	}
}

func (m *Manager) Path() *framework.Path {
	return &framework.Path{
		Pattern: "^creds/.+$",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: m.read,
		},
	}
}

// TODO implements DeleteHandler in roles, but this code has gotten scary, need to rethink
func (m *Manager) Delete(ctx context.Context, storage logical.Storage, roleName string) error {

	if err := storage.Delete(ctx, StorageKey+"/"+roleName); err != nil {
		return err
	}

	m.cache.Delete(roleName)
	return nil
}

func (m *Manager) read(ctx context.Context, req *logical.Request, fd *framework.FieldData) (*logical.Response, error) {

	roleName, err := util.ParseRoleName(BackendPath+"/", req.Path)
	if err != nil {
		return nil, err
	}

	role, err := m.roleHandler.Role(ctx, req.Storage, roleName)
	if err != nil {
		_, ok := err.(*roles.NotFound)
		if ok {
			return nil, nil
		}
		return nil, err
	}

	// Have we ever managed this cred before?
	// If not, we need to rotate the password so Vault will know it.
	var unset time.Time
	if role.LastVaultRotation == unset {
		return m.generateAndReturnCreds(ctx, req.Storage, role, &credential{})
	}

	// Has anyone manually rotated the password in Active Directory?
	// If so, we need to rotate it now so Vault will know it.
	if role.PasswordLastSet.After(role.LastVaultRotation) {
		return m.generateAndReturnCreds(ctx, req.Storage, role, &credential{})
	}

	// Since we should know the last password, let's retrieve it now so we can return it with the new one.
	cred := &credential{}
	credIfc, found := m.cache.Get(roleName)
	if found {
		cred = credIfc.(*credential)
	} else {
		entry, err := req.Storage.Get(ctx, StorageKey+"/"+roleName)
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
		m.cache.SetDefault(cred.RoleName, cred)
	}

	// Is the password too old?
	// If so, time for a new one!
	if role.LastVaultRotation.Add(time.Duration(role.TTL) * time.Second).Before(time.Now()) {
		return m.generateAndReturnCreds(ctx, req.Storage, role, cred)
	}

	// Current credential is accurate! Return it.
	return &logical.Response{
		Data: cred.Map(),
	}, nil
}

func (m *Manager) generateAndReturnCreds(ctx context.Context, storage logical.Storage, role *roles.Role, previousCred *credential) (*logical.Response, error) {

	engineConf, err := m.configHandler.Config(ctx, storage)
	if err != nil {
		return nil, err
	}

	adClient := activedirectory.NewClient(m.logger, engineConf.ADConf)

	newPassword, err := generatePassword(engineConf.PasswordConf.Length)
	if err != nil {
		return nil, err
	}

	filters := map[*activedirectory.Field][]string{
		activedirectory.FieldRegistry.UserPrincipalName: {role.ServiceAccountName},
	}
	if err := adClient.UpdatePassword(filters, newPassword); err != nil {
		return nil, err
	}

	// Time recorded should be UTC for easier user comparison to AD's last rotated time.
	role.LastVaultRotation = time.Now().UTC()
	if err := m.roleHandler.Update(ctx, storage, role); err != nil {
		return nil, err
	}

	cred := &credential{
		RoleName:        role.Name,
		Username:        role.ServiceAccountName,
		CurrentPassword: newPassword,
	}

	if previousCred.CurrentPassword != "" {
		cred.LastPassword = previousCred.CurrentPassword
	}

	if err := m.cacheAndSave(ctx, storage, cred); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: cred.Map(),
	}, nil
}

func (m *Manager) cacheAndSave(ctx context.Context, storage logical.Storage, cred *credential) error {

	m.cache.SetDefault(cred.RoleName, cred)

	entry, err := logical.StorageEntryJSON(StorageKey+"/"+cred.RoleName, cred)
	if err != nil {
		return err
	}

	if err := storage.Put(ctx, entry); err != nil {
		return err
	}
	return nil
}

func generatePassword(desiredLength int) (string, error) {

	if desiredLength <= 0 {
		return "", fmt.Errorf("it's not possible to generate a password of password_length %d", desiredLength)
	}
	if desiredLength < 14 {
		return "", fmt.Errorf("it's not possible to generate a _secure_ password of length %d, please boost password_length to 14, though Vault recommends higher", desiredLength)
	}

	complexityPrefix := "?@09AZ"

	// First, get some cryptographically secure pseudorandom bytes.
	b := make([]byte, desiredLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	result := ""
	// Though the result should immediately be longer than the desiredLength,
	// do this in a loop to ensure there's absolutely no risk of a panic when slicing it down later.
	for len(result) <= desiredLength {
		// Encode to base64 because it's more complex and performant than base62.
		result += base64.StdEncoding.EncodeToString(b)
	}

	result = complexityPrefix + result
	return result[:desiredLength], nil
}
