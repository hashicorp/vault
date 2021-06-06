package vault

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

type RouterTestHandlerFunc func(context.Context, *logical.Request) (*logical.Response, error)

type NoopBackend struct {
	sync.Mutex

	Root            []string
	Login           []string
	Paths           []string
	Requests        []*logical.Request
	Response        *logical.Response
	RequestHandler  RouterTestHandlerFunc
	Invalidations   []string
	DefaultLeaseTTL time.Duration
	MaxLeaseTTL     time.Duration
	BackendType     logical.BackendType
}

func NoopBackendFactory(_ context.Context, _ *logical.BackendConfig) (logical.Backend, error) {
	return &NoopBackend{}, nil
}

func (n *NoopBackend) HandleRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	if req.TokenEntry() != nil {
		panic("got a non-nil TokenEntry")
	}

	var err error
	resp := n.Response
	if n.RequestHandler != nil {
		resp, err = n.RequestHandler(ctx, req)
	}

	n.Lock()
	defer n.Unlock()

	requestCopy := *req
	n.Paths = append(n.Paths, req.Path)
	n.Requests = append(n.Requests, &requestCopy)
	if req.Storage == nil {
		return nil, fmt.Errorf("missing view")
	}

	if req.Path == "panic" {
		panic("as you command")
	}

	return resp, err
}

func (n *NoopBackend) HandleExistenceCheck(ctx context.Context, req *logical.Request) (bool, bool, error) {
	return false, false, nil
}

func (n *NoopBackend) SpecialPaths() *logical.Paths {
	return &logical.Paths{
		Root:            n.Root,
		Unauthenticated: n.Login,
	}
}

func (n *NoopBackend) System() logical.SystemView {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	if n.DefaultLeaseTTL > 0 {
		defaultLeaseTTLVal = n.DefaultLeaseTTL
	}

	if n.MaxLeaseTTL > 0 {
		maxLeaseTTLVal = n.MaxLeaseTTL
	}

	return logical.StaticSystemView{
		DefaultLeaseTTLVal: defaultLeaseTTLVal,
		MaxLeaseTTLVal:     maxLeaseTTLVal,
	}
}

func (n *NoopBackend) Cleanup(ctx context.Context) {
	// noop
}

func (n *NoopBackend) InvalidateKey(ctx context.Context, k string) {
	n.Invalidations = append(n.Invalidations, k)
}

func (n *NoopBackend) Setup(ctx context.Context, config *logical.BackendConfig) error {
	return nil
}

func (n *NoopBackend) Logger() log.Logger {
	return log.NewNullLogger()
}

func (n *NoopBackend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	return nil
}

func (n *NoopBackend) Type() logical.BackendType {
	if n.BackendType == logical.TypeUnknown {
		return logical.TypeLogical
	}
	return n.BackendType
}

// InitializableBackend is a backend that knows whether it has been initialized
// properly.
type InitializableBackend struct {
	*NoopBackend
	isInitialized bool
}

func (b *InitializableBackend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	if b.isInitialized {
		return errors.New("already initialized")
	}

	// do a dummy write, to prove that the storage is not readonly
	entry := &logical.StorageEntry{
		Key:   "initialize/zork",
		Value: []byte("quux"),
	}
	err := req.Storage.Put(ctx, entry)
	if err != nil {
		return err
	}

	b.isInitialized = true
	return nil
}
