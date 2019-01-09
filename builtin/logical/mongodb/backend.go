package mongodb

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	mgo "gopkg.in/mgo.v2"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config/connection",
			},
		},

		Paths: []*framework.Path{
			pathConfigConnection(&b),
			pathConfigLease(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},

		Clean: b.ResetSession,

		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	return b.Backend
}

type backend struct {
	*framework.Backend

	session *mgo.Session
	lock    sync.Mutex
}

// Session returns the database connection.
func (b *backend) Session(ctx context.Context, s logical.Storage) (*mgo.Session, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.session != nil {
		if err := b.session.Ping(); err == nil {
			return b.session, nil
		}
		b.session.Close()
	}

	connConfigJSON, err := s.Get(ctx, "config/connection")
	if err != nil {
		return nil, err
	}
	if connConfigJSON == nil {
		return nil, fmt.Errorf("configure the MongoDB connection with config/connection first")
	}

	var connConfig connectionConfig
	if err := connConfigJSON.DecodeJSON(&connConfig); err != nil {
		return nil, err
	}

	dialInfo, err := parseMongoURI(connConfig.URI)
	if err != nil {
		return nil, err
	}

	b.session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}
	b.session.SetSyncTimeout(1 * time.Minute)
	b.session.SetSocketTimeout(1 * time.Minute)

	return b.session, nil
}

// ResetSession forces creation of a new connection next time Session() is called.
func (b *backend) ResetSession(_ context.Context) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.session != nil {
		b.session.Close()
	}

	b.session = nil
}

func (b *backend) invalidate(ctx context.Context, key string) {
	switch key {
	case "config/connection":
		b.ResetSession(ctx)
	}
}

// LeaseConfig returns the lease configuration
func (b *backend) LeaseConfig(ctx context.Context, s logical.Storage) (*configLease, error) {
	entry, err := s.Get(ctx, "config/lease")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result configLease
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

const backendHelp = `
The mongodb backend dynamically generates MongoDB credentials.

After mounting this backend, configure it using the endpoints within
the "config/" path.
`
