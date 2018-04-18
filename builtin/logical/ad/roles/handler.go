package roles

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/builtin/logical/ad/util"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/patrickmn/go-cache"
)

const (
	BackendPath = "roles"
	storageKey  = "roles"

	cacheCleanup    = time.Second / 2
	cacheExpiration = time.Second

	unsetTTL = -1000
)

func Handler(logger hclog.Logger, configReader config.Reader) *handler {
	return &handler{
		logger:         logger,
		configReader:   configReader,
		cache:          cache.New(cacheExpiration, cacheCleanup),
		deleteWatchers: []DeleteWatcher{},
	}
}

type handler struct {
	logger         hclog.Logger
	configReader   config.Reader
	cache          *cache.Cache
	deleteWatchers []DeleteWatcher
}

func (h *handler) AddDeleteWatcher(w DeleteWatcher) {
	h.deleteWatchers = append(h.deleteWatchers, w)
}

func (h *handler) Invalidate(ctx context.Context, key string) {
	prefix := BackendPath + "/"
	if strings.HasPrefix(key, prefix) {
		roleName, err := util.ParseRoleName(prefix, key)
		if err != nil {
			// An invalid roleName was provided so we can't have this in storage anyways.
			return
		}
		h.cache.Delete(roleName)
	}
}

func (h *handler) Path() *framework.Path {
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
			logical.UpdateOperation: h.updateOperation,
			logical.ReadOperation:   h.readOperation,
			logical.ListOperation:   h.listOperation,
			logical.DeleteOperation: h.deleteOperation,
		},
	}
}

func (h *handler) Read(ctx context.Context, storage logical.Storage, roleName string) (*Role, error) {

	// If it's cached, return it from there.
	roleIfc, found := h.cache.Get(roleName)
	if found {
		return roleIfc.(*Role), nil
	}

	// It's not, read it from storage.
	entry, err := storage.Get(ctx, storageKey+"/"+roleName)
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

	// Always check when ActiveDirectory shows the password as last set on the fly.
	engineConf, err := h.configReader.Read(ctx, storage)
	if err != nil {
		return nil, err
	}

	passwordLastSet, err := util.NewSecretsClient(h.logger, engineConf.ADConf).GetPasswordLastSet(role.ServiceAccountName)
	if err != nil {
		return nil, err
	}
	role.PasswordLastSet = passwordLastSet

	// Cache it.
	h.cache.SetDefault(roleName, role)

	return role, nil
}

func (h *handler) Write(ctx context.Context, storage logical.Storage, role *Role) error {
	entry, err := logical.StorageEntryJSON(storageKey+"/"+role.Name, role)
	if err != nil {
		return err
	}
	if err := storage.Put(ctx, entry); err != nil {
		return err
	}
	h.cache.SetDefault(role.Name, role)
	return nil
}

func (h *handler) updateOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {

	// Get everything we need to construct the role.
	roleName, err := util.ParseRoleName(BackendPath+"/", req.Path)
	if err != nil {
		return nil, err
	}

	engineConf, err := h.configReader.Read(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// Actually construct it.
	role, err := newRole(h.logger, engineConf, roleName, fieldData)
	if err != nil {
		return nil, err
	}

	// Write it to storage and the cache.
	if err := h.Write(ctx, req.Storage, role); err != nil {
		return nil, err
	}

	// We don't want to store when AD shows the password as last set because we're going to pull it on the fly regularly.
	passwordLastSet, err := util.NewSecretsClient(h.logger, engineConf.ADConf).GetPasswordLastSet(role.ServiceAccountName)
	if err != nil {
		return nil, err
	}
	role.PasswordLastSet = passwordLastSet

	return &logical.Response{
		Data: role.Map(),
	}, nil
}

func (h *handler) readOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	roleName, err := util.ParseRoleName(BackendPath+"/", req.Path)
	if err != nil {
		return nil, err
	}

	role, err := h.Read(ctx, req.Storage, roleName)
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

func (h *handler) listOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	keys, err := req.Storage.List(ctx, storageKey+"/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(keys), nil
}

func (h *handler) deleteOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	roleName, err := util.ParseRoleName(BackendPath+"/", req.Path)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Delete(ctx, storageKey+"/"+roleName); err != nil {
		return nil, err
	}

	h.cache.Delete(roleName)

	for _, h := range h.deleteWatchers {
		if err := h.Delete(ctx, req.Storage, roleName); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

type ReadWriter interface {
	// Read returns the present *Read for a by a given name.
	// If error == nil, *Read != nil.
	//
	// The returned error may be due to issues reaching storage,
	// or it may be because the *Read is unset by the user - they have no role by the given name.
	//
	// If knowing the error is useful to the caller, it may be inspected like so:
	//
	// 		role, err := Read(ctx, storage, name)
	//		if err != nil {
	// 			_, ok := err.(*NotFound)
	//			...
	// 		}
	Read(ctx context.Context, storage logical.Storage, roleName string) (*Role, error)

	// Write allows the role to be updated outside the package.
	// This is useful for updating the role.LastVaultRotation time when new creds are generated.
	Write(ctx context.Context, storage logical.Storage, role *Role) error
}

type DeleteWatcher interface {
	Delete(ctx context.Context, storage logical.Storage, roleName string) error
}

type NotFound struct {
	roleName string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("%s not found", e.roleName)
}
