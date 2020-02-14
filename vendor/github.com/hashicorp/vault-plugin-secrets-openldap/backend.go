package openldap

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	ldapClient := NewClient()
	b := Backend(ldapClient)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	b.credRotationQueue = queue.New()
	// Create a context with a cancel method for processing any WAL entries and
	// populating the queue
	initCtx := context.Background()
	ictx, cancel := context.WithCancel(initCtx)
	b.cancelQueue = cancel

	// Load queue and kickoff new periodic ticker
	go b.initQueue(ictx, conf)

	return b, nil
}

func Backend(client ldapClient) *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				"config",
				"static-role/*",
			},
		},
		Paths: framework.PathAppend(
			b.pathListRoles(),
			b.pathRoles(),
			b.pathCredsCreate(),
			b.pathRotateCredentials(),
			b.pathConfig(),
		),

		Secrets:     []*framework.Secret{},
		Clean:       b.clean,
		BackendType: logical.TypeLogical,
	}
	b.client = client
	b.roleLocks = locksutil.CreateLocks()

	return &b
}

func (b *backend) clean(ctx context.Context) {
	b.invalidateQueue()
}

// invalidateQueue cancels any background queue loading and destroys the queue.
func (b *backend) invalidateQueue() {
	b.Lock()
	defer b.Unlock()

	if b.cancelQueue != nil {
		b.cancelQueue()
	}
	b.credRotationQueue = nil
}

type backend struct {
	*framework.Backend
	sync.RWMutex
	// CredRotationQueue is an in-memory priority queue used to track Static Roles
	// that require periodic rotation. Backends will have a PriorityQueue
	// initialized on setup, but only backends that are mounted by a primary
	// server or mounted as a local mount will perform the rotations.
	//
	// cancelQueue is used to remove the priority queue and terminate the
	// background ticker.
	credRotationQueue *queue.PriorityQueue
	cancelQueue       context.CancelFunc

	// roleLocks is used to lock modifications to roles in the queue, to ensure
	// concurrent requests are not modifying the same role and possibly causing
	// issues with the priority queue.
	roleLocks []*locksutil.LockEntry
	client    ldapClient
}

const backendHelp = `
The OpenLDAP backend supports managing existing LDAP entry passwords by providing:

 * end points to add entries
 * manual rotation of entry passwords
 * auto rotation of entry passwords
 
The OpenLDAP secret engine is limited to OpenLDAP and does not support any other 
implementations of LDAP.

After mounting this secret backend, configure it using the "openldap/config" path.
`
