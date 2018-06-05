package kv

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	// configPath is the location where the config is stored
	configPath string = "config"

	// metadataPrefix is the prefix where the key metadata is stored.
	metadataPrefix string = "metadata/"

	// versionPrefix is the prefix where the version data is stored.
	versionPrefix string = "versions/"

	// defaultMaxVersions is the number of versions to keep around unless set by
	// the config or key configuration.
	defaultMaxVersions uint32 = 10
)

// versionedKVBackend implements logical.Backend
type versionedKVBackend struct {
	*framework.Backend

	// keyEncryptedWrapper is a cached version of the EncryptedKeyStorageWrapper
	keyEncryptedWrapper *keysutil.EncryptedKeyStorageWrapper

	// salt is the cached version of the salt used to create paths for version
	// data storage paths.
	salt *salt.Salt

	// l locks the keyPolicy and salt caches.
	l sync.RWMutex

	// locks is a slice of 256 locks that are used to protect key and version
	// updates.
	locks []*locksutil.LockEntry

	// storagePrefix is the prefix given to all the data for a versioned KV
	// store. We prefix this data so that upgrading from a passthrough backend
	// to a versioned backend is easier. This value is passed from Vault core
	// through the backend config.
	storagePrefix string

	// upgrading is an atomic value denoting if the backend is in the process of
	// upgrading its data.
	upgrading *uint32

	// globalConfig is a cached value for fast lookup
	globalConfig     *Configuration
	globalConfigLock *sync.RWMutex
}

// Factory will return a logical backend of type versionedKVBackend or
// PassthroughBackend based on the config passed in.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	version := conf.Config["version"]

	var b logical.Backend
	var err error
	switch version {
	case "1", "":
		return LeaseSwitchedPassthroughBackend(ctx, conf, conf.Config["leased_passthrough"] == "true")
	case "2":
		b, err = VersionedKVFactory(ctx, conf)
	}
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Factory returns a new backend as logical.Backend.
func VersionedKVFactory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &versionedKVBackend{
		upgrading:        new(uint32),
		globalConfigLock: new(sync.RWMutex),
	}
	if conf.BackendUUID == "" {
		return nil, errors.New("could not initialize versioned K/V Store, no UUID was provided")
	}
	b.storagePrefix = conf.BackendUUID

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        backendHelp,

		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				// Seal wrap the versioned data
				path.Join(b.storagePrefix, versionPrefix) + "/",

				// Seal wrap the key policy
				path.Join(b.storagePrefix, "policy") + "/",

				// Seal wrap the archived key policy
				path.Join(b.storagePrefix, "archive") + "/",
			},
		},

		Paths: framework.PathAppend(
			[]*framework.Path{
				pathConfig(b),
				pathData(b),
				pathMetadata(b),
				pathDestroy(b),
			},
			pathsDelete(b),

			// Make sure this stays at the end so that the valid paths are
			// processed first.
			pathInvalid(b),
		),
	}

	b.locks = locksutil.CreateLocks()

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	upgradeDone, err := b.upgradeDone(ctx, conf.StorageView)
	if err != nil {
		return nil, err
	}
	if !upgradeDone {
		err := b.Upgrade(ctx, conf.StorageView)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (b *versionedKVBackend) upgradeDone(ctx context.Context, s logical.Storage) (bool, error) {
	upgradeEntry, err := s.Get(ctx, path.Join(b.storagePrefix, "upgrading"))
	if err != nil {
		return false, err
	}

	var upgradeInfo UpgradeInfo
	if upgradeEntry != nil {
		err := proto.Unmarshal(upgradeEntry.Value, &upgradeInfo)
		if err != nil {
			return false, err
		}
	}

	return upgradeInfo.Done, nil
}

func pathInvalid(b *versionedKVBackend) []*framework.Path {
	handler := func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		var subCommand string
		switch req.Operation {
		case logical.CreateOperation, logical.UpdateOperation:
			subCommand = "put"
		case logical.ReadOperation:
			subCommand = "get"
		case logical.ListOperation:
			subCommand = "list"
		case logical.DeleteOperation:
			subCommand = "delete"
		}
		resp := &logical.Response{}
		resp.AddWarning(fmt.Sprintf("Invalid path for a versioned K/V secrets engine. See the API docs for the appropriate API endpoints to use. If using the Vault CLI, use 'vault kv %s' for this operation.", subCommand))
		return logical.RespondWithStatusCode(resp, req, http.StatusNotFound)
	}

	return []*framework.Path{
		&framework.Path{
			Pattern: ".*",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: handler,
				logical.CreateOperation: handler,
				logical.ReadOperation:   handler,
				logical.DeleteOperation: handler,
				logical.ListOperation:   handler,
			},

			HelpDescription: pathInvalidHelp,
		},
	}
}

// Invalidate invalidates the salt and the policy so replication secondaries can
// cache these values.
func (b *versionedKVBackend) Invalidate(ctx context.Context, key string) {
	switch key {
	case path.Join(b.storagePrefix, salt.DefaultLocation):
		b.l.Lock()
		b.salt = nil
		b.l.Unlock()
	case path.Join(b.storagePrefix, "policy/metadata"):
		b.l.Lock()
		b.keyEncryptedWrapper = nil
		b.l.Unlock()
	case path.Join(b.storagePrefix, configPath):
		b.globalConfigLock.Lock()
		b.globalConfig = nil
		b.globalConfigLock.Unlock()
	}
}

// Salt will load a the salt, or if one has not been created yet it will
// generate and store a new salt.
func (b *versionedKVBackend) Salt(ctx context.Context, s logical.Storage) (*salt.Salt, error) {
	b.l.RLock()
	if b.salt != nil {
		defer b.l.RUnlock()
		return b.salt, nil
	}
	b.l.RUnlock()
	b.l.Lock()
	defer b.l.Unlock()
	if b.salt != nil {
		return b.salt, nil
	}
	salt, err := salt.NewSalt(ctx, s, &salt.Config{
		HashFunc: salt.SHA256Hash,
		Location: path.Join(b.storagePrefix, salt.DefaultLocation),
	})
	if err != nil {
		return nil, err
	}
	b.salt = salt
	return salt, nil
}

// policy loads the key policy for this backend, if one has not been created yet
// it will generate and store a new policy. The caller must have the backend lock.
func (b *versionedKVBackend) policy(ctx context.Context, s logical.Storage) (*keysutil.Policy, error) {
	// Try loading policy
	policy, err := keysutil.LoadPolicy(ctx, s, path.Join(b.storagePrefix, "policy/metadata"))
	if err != nil {
		return nil, err
	}
	if policy != nil {
		return policy, nil
	}

	// Policy didn't exist, create it.
	policy = keysutil.NewPolicy(keysutil.PolicyConfig{
		Name:                 "metadata",
		Type:                 keysutil.KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  keysutil.Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		StoragePrefix:        b.storagePrefix,
		VersionTemplate:      keysutil.EncryptedKeyPolicyVersionTpl,
	})

	err = policy.Rotate(ctx, s)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (b *versionedKVBackend) getKeyEncryptor(ctx context.Context, s logical.Storage) (*keysutil.EncryptedKeyStorageWrapper, error) {
	b.l.RLock()
	if b.keyEncryptedWrapper != nil {
		defer b.l.RUnlock()
		return b.keyEncryptedWrapper, nil
	}
	b.l.RUnlock()
	b.l.Lock()
	defer b.l.Unlock()

	if b.keyEncryptedWrapper != nil {
		return b.keyEncryptedWrapper, nil
	}

	policy, err := b.policy(ctx, s)
	if err != nil {
		return nil, err
	}

	e, err := keysutil.NewEncryptedKeyStorageWrapper(keysutil.EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: path.Join(b.storagePrefix, metadataPrefix),
	})
	if err != nil {
		return nil, err
	}

	// Cache the value
	b.keyEncryptedWrapper = e

	return b.keyEncryptedWrapper, nil
}

// config takes a storage object and returns a configuration object
func (b *versionedKVBackend) config(ctx context.Context, s logical.Storage) (*Configuration, error) {
	b.globalConfigLock.RLock()
	if b.globalConfig != nil {
		defer b.globalConfigLock.RUnlock()
		return &Configuration{
			CasRequired: b.globalConfig.CasRequired,
			MaxVersions: b.globalConfig.MaxVersions,
		}, nil
	}

	b.globalConfigLock.RUnlock()
	b.globalConfigLock.Lock()
	defer b.globalConfigLock.Unlock()

	// Verify this hasn't already changed
	if b.globalConfig != nil {
		return &Configuration{
			CasRequired: b.globalConfig.CasRequired,
			MaxVersions: b.globalConfig.MaxVersions,
		}, nil
	}

	raw, err := s.Get(ctx, path.Join(b.storagePrefix, configPath))
	if err != nil {
		return nil, err
	}

	conf := &Configuration{}
	if raw != nil {
		if err := proto.Unmarshal(raw.Value, conf); err != nil {
			return nil, err
		}
	}

	b.globalConfig = conf

	return conf, nil
}

// getVersionKey uses the salt to generate the version key for a specific
// version of a key.
func (b *versionedKVBackend) getVersionKey(ctx context.Context, key string, version uint64, s logical.Storage) (string, error) {
	salt, err := b.Salt(ctx, s)
	if err != nil {
		return "", err
	}

	salted := salt.SaltID(fmt.Sprintf("%s|%d", key, version))

	return path.Join(b.storagePrefix, versionPrefix, salted[0:3], salted[3:]), nil
}

// getKeyMetadata returns the metadata object for the provided key, if no object
// exits it will return nil.
func (b *versionedKVBackend) getKeyMetadata(ctx context.Context, s logical.Storage, key string) (*KeyMetadata, error) {

	wrapper, err := b.getKeyEncryptor(ctx, s)
	if err != nil {
		return nil, err
	}

	es := wrapper.Wrap(s)

	item, err := es.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}

	meta := &KeyMetadata{}
	err = proto.Unmarshal(item.Value, meta)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key metadata from storage: %v", err)
	}

	return meta, nil
}

// writeKeyMetadata writes a metadata object to storage.
func (b *versionedKVBackend) writeKeyMetadata(ctx context.Context, s logical.Storage, meta *KeyMetadata) error {
	wrapper, err := b.getKeyEncryptor(ctx, s)
	if err != nil {
		return err
	}

	es := wrapper.Wrap(s)

	bytes, err := proto.Marshal(meta)
	if err != nil {
		return err
	}

	err = es.Put(ctx, &logical.StorageEntry{
		Key:   meta.Key,
		Value: bytes,
	})
	if err != nil {
		return err
	}

	return nil
}

func ptypesTimestampToString(t *timestamp.Timestamp) string {
	if t == nil {
		return ""
	}

	return ptypes.TimestampString(t)
}

var backendHelp string = `
This backend provides a versioned key-value store. The kv backend reads and
writes arbitrary secrets to the storage backend. The secrets are
encrypted/decrypted by Vault: they are never stored unencrypted in the backend
and the backend never has an opportunity to see the unencrypted value. Each key
can have a configured number of versions, and versions can be retrieved based on
their version numbers.
`

var pathInvalidHelp string = backendHelp + `

## PATHS

The following paths are supported by this backend. To view help for
any of the paths below, use the help command with any route matching
the path pattern. Note that depending on the policy of your auth token,
you may or may not be able to access certain paths.

    ^config$
        Configures settings for the KV store

    ^data/.*$
        Write, Read, and Delete data in the Key-Value Store.

    ^delete/.*$
        Marks one or more versions as deleted in the KV store.

    ^destroy/.*$
        Permanently removes one or more versions in the KV store

    ^metadata/.*$
        Configures settings for the KV store

    ^undelete/.*$
        Undeletes one or more versions from the KV store.
`
