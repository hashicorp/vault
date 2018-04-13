package roles

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/helper/activedirectory"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	BackendPath = "roles"
	StorageKey  = "roles"

	unsetTTL = -1000
)

func NewManager(logger hclog.Logger, configReader config.Reader) *Manager {
	return &Manager{
		logger:       logger,
		configReader: configReader,
	}
}

type Manager struct {
	logger       hclog.Logger
	configReader config.Reader
}

func (m *Manager) Path() *framework.Path {
	return &framework.Path{
		Pattern: `^roles$|^roles/.|^roles/$`,
		Fields: map[string]*framework.FieldSchema{
			"service_account_name": {
				Type:        framework.TypeString,
				Description: "The username/logon name for the service account with which this role will be associated.",
			},

			"ttl": {
				Type:        framework.TypeDurationSecond,
				Default:     unsetTTL,
				Description: "In seconds, the default password time-to-live.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: m.delete,
			logical.ListOperation:   m.list,
			logical.ReadOperation:   m.read,
			logical.UpdateOperation: m.update,
		},
	}
}

type NotFound struct {
	roleName string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("%s not found", e.roleName)
}

// TODO
//func (m *Manager) Role(ctx context.Context, storage logical.Storage, name string) (*Role, error) {
//
//}

func (m *Manager) delete(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	roleName, err := roleName(req.Path)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Delete(ctx, StorageKey+"/"+roleName); err != nil {
		return nil, err
	}
	return nil, nil
}

func (m *Manager) list(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	keys, err := req.Storage.List(ctx, StorageKey+"/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(keys), nil
}

func (m *Manager) read(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	roleName, err := roleName(req.Path)
	if err != nil {
		return nil, err
	}

	entry, err := req.Storage.Get(ctx, StorageKey+"/"+roleName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	role := &Role{}
	if err := entry.DecodeJSON(role); err != nil {
		return nil, err
	}

	passwordLastSet, err := m.getPasswordLastSet(ctx, req.Storage, role.ServiceAccountName)
	if err != nil {
		return nil, err
	}
	role.PasswordLastSet = passwordLastSet

	return &logical.Response{
		Data: role.Map(),
	}, nil
}

func (m *Manager) update(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {

	roleName, err := roleName(req.Path)
	if err != nil {
		return nil, err
	}

	role, err := newRole(m.logger, ctx, req.Storage, m.configReader, roleName, fieldData)
	if err != nil {
		return nil, err
	}

	entry, err := logical.StorageEntryJSON(StorageKey+"/"+roleName, role)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	passwordLastSet, err := m.getPasswordLastSet(ctx, req.Storage, role.ServiceAccountName)
	if err != nil {
		return nil, err
	}
	role.PasswordLastSet = passwordLastSet

	return &logical.Response{
		Data: role.Map(),
	}, nil
}

// TODO this manager is going to need to expose a way for the creds endpoint to read a role
// so it can find out last rotated times

// TODO we'll need to cache roles too because when requests come in quickly, we'll need the same role like 100 times at once

func (m *Manager) getPasswordLastSet(ctx context.Context, storage logical.Storage, serviceAccountName string) (*time.Time, error) {

	engineConf, err := m.configReader.Config(ctx, storage)
	if err != nil {
		return nil, err
	}

	filters := map[*activedirectory.Field][]string{
		activedirectory.FieldRegistry.UserPrincipalName: {serviceAccountName},
	}

	// TODO make a helper method like get by service account name
	adClient := activedirectory.NewClient(m.logger, engineConf.ADConf)
	entries, err := adClient.Search(filters)
	if err != nil {
		return nil, err
	}
	if len(entries) <= 0 {
		return nil, fmt.Errorf("service account \"%s\" is not found by Active Directory", serviceAccountName)
	}
	if len(entries) > 1 {
		return nil, fmt.Errorf("unable to tell which service account to use from %s", entries)
	}

	values, found := entries[0].Get(activedirectory.FieldRegistry.PasswordLastSet)
	if !found {
		return nil, fmt.Errorf("%s lacks a PasswordLastSet field", entries[0])
	}

	if len(values) != 1 {
		return nil, fmt.Errorf("expected only one value for PasswordLastSet, but received %s", values)
	}

	ticks := values[0]
	if ticks == "0" {
		// password has never been rolled in Active Directory, only created
		return nil, nil
	}

	t, err := activedirectory.ParseTime(ticks)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func roleName(reqPath string) (string, error) {

	prefix := BackendPath + "/"
	prefixLen := len(prefix)

	if len(reqPath) <= prefixLen {
		return "", errors.New("role name must be provided")
	}
	return reqPath[prefixLen:], nil
}
