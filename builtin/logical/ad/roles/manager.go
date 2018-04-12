package roles

import (
	"context"
	"errors"

	"regexp"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/ad/config"
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

	keys, err := req.Storage.List(ctx, StorageKey)
	if err != nil {
		return nil, err
	}
	// TODO strip me, just for getting started
	for _, key := range keys {
		m.logger.Info(key)
	}

	return nil, nil
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

	return &logical.Response{
		Data: role.Map(),
	}, nil
}

func (m *Manager) update(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {

	// TODO stripme
	m.logger.Info("reqPath: " + req.Path)

	roleName, err := roleName(req.Path)
	if err != nil {
		return nil, err
	}
	// TODO stripme
	m.logger.Info("roleName: " + roleName)

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

	return &logical.Response{
		Data: role.Map(),
	}, nil
}

// TODO test with some real reqPaths, this algorithm assumes a path like "role/kibana"
func roleName(reqPath string) (string, error) {

	prefix := BackendPath + "/"
	prefixLen := len(prefix)

	if len(reqPath) <= prefixLen {
		return "", errors.New("role name must be provided")
	}

	roleName := reqPath[prefixLen:]

	// ensure it's a valid name per Vault standards
	// TODO throw in some junk names and see if this is actually helpful
	matched, err := regexp.MatchString(framework.GenericNameRegex("name"), roleName)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", errors.New(roleName + " has an unacceptable character")
	}

	return roleName, nil
}
