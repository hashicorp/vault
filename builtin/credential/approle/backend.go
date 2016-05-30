package approle

import (
	"fmt"
	"sync"

	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type backend struct {
	*framework.Backend

	// The salt value to be used by the information to be accessed only
	// by this backend.
	salt *salt.Salt

	// Guard to clean-up the expired SecretID entries
	tidySecretIDCASGuard uint32

	// Lock to make changes to Role entries. This is a low-traffic
	// operation. So, using a single lock would suffice.
	roleLock *sync.RWMutex

	// Map of locks to make changes to the storage entries of RoleIDs
	// generated. This will be initiated to a predefined number of locks
	// when the backend is created, and will be indexed based on the salted
	// RoleIDs.
	roleIDLocksMap map[string]*sync.RWMutex

	// Map of locks to make changes to the storage entries of SecretIDs
	// generated. This will be initiated to a predefined number of locks
	// when the backend is created, and will be indexed based on the HMAC-ed
	// SecretIDs.
	secretIDLocksMap map[string]*sync.RWMutex
}

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}
	return b.Setup(conf)
}

func Backend(conf *logical.BackendConfig) (*backend, error) {
	// Initialize the salt
	salt, err := salt.NewSalt(conf.StorageView, &salt.Config{
		HashFunc: salt.SHA256Hash,
	})
	if err != nil {
		return nil, err
	}

	// Create a backend object
	b := &backend{
		// Set the salt object for the backend
		salt: salt,

		// Create the lock for making changes to the Roles registered with the backend
		roleLock: &sync.RWMutex{},

		// Create the map of locks to modify the generated RoleIDs.
		roleIDLocksMap: map[string]*sync.RWMutex{},

		// Create the map of locks to modify the generated SecretIDs.
		secretIDLocksMap: map[string]*sync.RWMutex{},
	}

	// Create 256 locks each for managing RoleID and SecretIDs. This will avoid
	// a superfluous number of locks directly proportional to the number of RoleID
	// and SecretIDs. These locks can be accessed by indexing based on the first two
	// characters of a randomly generated UUID.
	if err = locksutil.CreateLocks(b.roleIDLocksMap, 256); err != nil {
		return nil, fmt.Errorf("failed to create role ID locks: %v", err)
	}

	if err = locksutil.CreateLocks(b.secretIDLocksMap, 256); err != nil {
		return nil, fmt.Errorf("failed to create secret ID locks: %v", err)
	}

	// Have an extra lock to use in case the indexing does not result in a lock.
	// This happens if the indexing value is not beginning with hex characters.
	// These locks can be used for listing purposes as well.
	b.secretIDLocksMap["custom"] = &sync.RWMutex{}
	b.roleIDLocksMap["custom"] = &sync.RWMutex{}

	// Attach the paths and secrets that are to be handled by the backend
	b.Backend = &framework.Backend{
		// Register a periodic function that deletes the expired SecretID entries
		PeriodicFunc: b.periodicFunc,
		Help:         backendHelp,
		AuthRenew:    b.pathLoginRenew,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},
		Paths: framework.PathAppend(
			rolePaths(b),
			[]*framework.Path{
				pathLogin(b),
				pathTidySecretID(b),
			},
		),
	}
	return b, nil
}

// periodicFunc of the backend will be invoked once a minute by the RollbackManager.
// RoleRole backend utilizes this function to delete expired SecretID entries.
// This could mean that the SecretID may live in the backend upto 1 min after its
// expiration. The deletion of SecretIDs are not security sensitive and it is okay
// to delay the removal of SecretIDs by a minute.
func (b *backend) periodicFunc(req *logical.Request) error {
	// Initiate clean-up of expired SecretID entries
	b.tidySecretID(req.Storage)
	return nil
}

const backendHelp = `
Any registered Role can authenticate itself with Vault. The credentials
depends on the constraints that are set on the Role. One common required
credential is the 'role_id' which is a unique identifier of the Role.
It can be retrieved from the 'role/<appname>/role-id' endpoint.

The default constraint configuration is 'bind_secret_id', which requires
the credential 'secret_id' to be presented during login. Refer to the
documentation for other types of constraints.`
