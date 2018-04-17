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
	"github.com/patrickmn/go-cache"
)

const (
	BackendPath = "roles"
	StorageKey  = "roles"

	cacheCleanup    = time.Second / 2
	cacheExpiration = time.Second

	unsetTTL = -1000
)

func NewManager(logger hclog.Logger, configReader config.Reader) *Manager {
	return &Manager{
		logger:       logger,
		configReader: configReader,
		cache:        cache.New(cacheExpiration, cacheCleanup),
	}
}

type Manager struct {
	logger       hclog.Logger
	configReader config.Reader
	cache        *cache.Cache
}

func (m *Manager) Invalidate(ctx context.Context, key string) {
	if key == BackendPath {
		m.cache.Flush()
	}
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

type Reader interface {
	Role(ctx context.Context, storage logical.Storage, name string) (*Role, error)
}

type NotFound struct {
	roleName string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("%s not found", e.roleName)
}

func (m *Manager) Role(ctx context.Context, storage logical.Storage, name string) (*Role, error) {

	roleIfc, found := m.cache.Get(name)
	if found {
		return roleIfc.(*Role), nil
	}

	role, err := m.readFromStorage(ctx, storage, name)
	if err != nil {
		return nil, err
	}

	m.cache.SetDefault(name, role)

	return role, nil
}

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

	role, err := m.readFromStorage(ctx, req.Storage, roleName)
	if err != nil {
		_, ok := err.(*NotFound)
		if ok {
			return nil, nil
		}
	}

	return &logical.Response{
		Data: role.Map(),
	}, nil
}

func (m *Manager) update(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {

	roleName, err := roleName(req.Path)
	if err != nil {
		return nil, err
	}

	engineConf, err := m.configReader.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	adClient := activedirectory.NewClient(m.logger, engineConf.ADConf)

	role, err := newRole(adClient, engineConf.PasswordConf, roleName, fieldData)
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

func (m *Manager) getPasswordLastSet(ctx context.Context, storage logical.Storage, serviceAccountName string) (time.Time, error) {

	engineConf, err := m.configReader.Config(ctx, storage)
	if err != nil {
		return time.Time{}, err
	}

	adClient := activedirectory.NewClient(m.logger, engineConf.ADConf)
	entry, err := getServiceAccountByName(adClient, serviceAccountName)
	if err != nil {
		return time.Time{}, err
	}

	values, found := entry.Get(activedirectory.FieldRegistry.PasswordLastSet)
	if !found {
		return time.Time{}, fmt.Errorf("%s lacks a PasswordLastSet field", entry)
	}

	if len(values) != 1 {
		return time.Time{}, fmt.Errorf("expected only one value for PasswordLastSet, but received %s", values)
	}

	ticks := values[0]
	if ticks == "0" {
		// password has never been rolled in Active Directory, only created
		return time.Time{}, nil
	}

	t, err := activedirectory.ParseTime(ticks)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (m *Manager) readFromStorage(ctx context.Context, storage logical.Storage, roleName string) (*Role, error) {

	entry, err := storage.Get(ctx, StorageKey+"/"+roleName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, &NotFound{}
	}

	role := &Role{}
	if err := entry.DecodeJSON(role); err != nil {
		return nil, err
	}

	passwordLastSet, err := m.getPasswordLastSet(ctx, storage, role.ServiceAccountName)
	if err != nil {
		return nil, err
	}
	role.PasswordLastSet = passwordLastSet

	return role, nil
}

func roleName(reqPath string) (string, error) {

	prefix := BackendPath + "/"
	prefixLen := len(prefix)

	if len(reqPath) <= prefixLen {
		return "", errors.New("role name must be provided")
	}
	return reqPath[prefixLen:], nil
}
