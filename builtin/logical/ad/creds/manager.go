package creds

import (
	"context"
	"errors"
	"time"

	"fmt"

	"crypto/rand"
	"encoding/base64"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/builtin/logical/ad/roles"
	"github.com/hashicorp/vault/helper/activedirectory"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/patrickmn/go-cache"
)

const (
	BackendPath = "creds"
	StorageKey  = "creds"
)

func NewManager(logger hclog.Logger, configReader config.Reader, roleReader roles.Reader) *Manager {
	return &Manager{
		logger:       logger,
		configReader: configReader,
		roleReader:   roleReader,
	}
}

type Manager struct {
	logger       hclog.Logger
	configReader config.Reader
	roleReader   roles.Reader
	cache        *cache.Cache
}

func (m *Manager) Invalidate(ctx context.Context, key string) {
	if key == BackendPath {
		m.cache.Flush()
	}
}

func (m *Manager) Path() *framework.Path {
	return &framework.Path{
		Pattern: framework.GenericNameRegex(BackendPath), // TODO if this doesn't immediately work, borrow from the roles regexes
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: m.read,
		},
	}
}

func (m *Manager) read(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	roleName, err := roleName(req.Path)
	if err != nil {
		return nil, err
	}

	role, err := m.roleReader.Role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	// Have we ever managed this cred before?
	// If not, we need to rotate the password so Vault will know it.
	var unset time.Time
	if role.LastVaultRotation == unset {
		return m.generateAndReturnCreds(ctx, req.Storage, role)
	}

	// Has anyone manually rotated the password in Active Directory?
	// If so, we need to rotate it now so Vault will know it.
	if role.PasswordLastSet.After(role.LastVaultRotation) {
		return m.generateAndReturnCreds(ctx, req.Storage, role)
	}

	// Is the password too old?
	// If so, time for a new one!
	if role.LastVaultRotation.Add(time.Duration(role.TTL) * time.Second).Before(time.Now()) {
		return m.generateAndReturnCreds(ctx, req.Storage, role)
	}

	// The current password should be accurate. Can we just return it from the cache?
	credIfc, found := m.cache.Get(roleName)
	if found {
		cred := credIfc.(*Cred)
		return &logical.Response{
			Data: cred.Map(),
		}, nil
	}

	// It's not cached but it should be in storage.
	entry, err := req.Storage.Get(ctx, StorageKey+"/"+roleName)
	if err != nil {
		return nil, err
	}
	if entry != nil {
		cred := &Cred{}
		if err := entry.DecodeJSON(cred); err != nil {
			return nil, err
		}
		if err := m.cacheAndSave(ctx, req.Storage, cred); err != nil {
			return nil, err
		}
		return &logical.Response{
			Data: cred.Map(),
		}, nil
	}

	// If the creds aren't in storage, but roles are and we've created creds before,
	// this is an unexpected state and something has gone wrong.
	// Let's be explicit and error about this.
	return nil, fmt.Errorf("should have the creds for %+v but they're not found", role)
}

func (m *Manager) generateAndReturnCreds(ctx context.Context, storage logical.Storage, role *roles.Role) (*logical.Response, error) {

	engineConf, err := m.configReader.Config(ctx, storage)
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

	cred := &Cred{
		RoleName:        role.Name,
		Username:        role.ServiceAccountName,
		CurrentPassword: newPassword,
	}

	if err := m.cacheAndSave(ctx, storage, cred); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: cred.Map(),
	}, nil
}

func (m *Manager) cacheAndSave(ctx context.Context, storage logical.Storage, cred *Cred) error {

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

// TODO can this be moved to util.go in the next package up, or would that be a circular dependency?
// or is there already something for this somewhere?
func roleName(reqPath string) (string, error) {

	prefix := BackendPath + "/"
	prefixLen := len(prefix)

	if len(reqPath) <= prefixLen {
		return "", errors.New("role name must be provided")
	}
	return reqPath[prefixLen:], nil
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
