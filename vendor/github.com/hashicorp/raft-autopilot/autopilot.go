package autopilot

import (
	"context"
	"sync"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
)

const (
	// These constants were take from what exists in Consul at the time of module extraction.

	DefaultUpdateInterval    = 2 * time.Second
	DefaultReconcileInterval = 10 * time.Second
)

// Option is an option to be used when creating a new Autopilot instance
type Option func(*Autopilot)

// WithUpdateInterval returns an Option to set the Autopilot instance's
// update interval.
func WithUpdateInterval(t time.Duration) Option {
	if t == 0 {
		t = DefaultUpdateInterval
	}
	return func(a *Autopilot) {
		a.updateInterval = t
	}
}

// WithReconcileInterval returns an Option to set the Autopilot instance's
// reconcile interval.
func WithReconcileInterval(t time.Duration) Option {
	if t == 0 {
		t = DefaultReconcileInterval
	}
	return func(a *Autopilot) {
		a.reconcileInterval = t
	}
}

// WithLogger returns an Option to set the Autopilot instance's logger
func WithLogger(logger hclog.Logger) Option {
	if logger == nil {
		logger = hclog.Default()
	}

	return func(a *Autopilot) {
		a.logger = logger.Named("autopilot")
	}
}

// withTimeProvider returns an Option which overrides and Autopilot instance's
// time provider with the given one. This should only be used in tests
// as a means of making some time.Time values in an autopilot state deterministic.
// For real uses the default runtimeTimeProvider should be used.
func withTimeProvider(provider timeProvider) Option {
	return func(a *Autopilot) {
		a.time = provider
	}
}

// WithPromoter returns an option to set the Promoter type that Autpilot will
// use. When the option is not given the default StablePromoter from this package
// will be used.
func WithPromoter(promoter Promoter) Option {
	if promoter == nil {
		promoter = DefaultPromoter()
	}

	return func(a *Autopilot) {
		a.promoter = promoter
	}
}

// ExecutionStatus represents the current status of the autopilot background go routines
type ExecutionStatus string

const (
	NotRunning   ExecutionStatus = "not-running"
	Running      ExecutionStatus = "running"
	ShuttingDown ExecutionStatus = "shutting-down"
)

type execInfo struct {
	// status is the current state of autopilot executation
	status ExecutionStatus

	// shutdown is a function that can be execute to shutdown a running
	// autopilot's go routines.
	shutdown context.CancelFunc

	// done is a chan that will be closed when the running autopilot go
	// routines have exited. Technically closing it is the very last
	// thing done in the go routine but at that point enough state has
	// been cleaned up that we would then allow it to be started
	// immediately afterward
	done chan struct{}
}

// Autopilot is the type to manage a running Raft instance.
//
// Each Raft node in the cluster will have a corresponding Autopilot instance but
// only 1 Autopilot instance should run at a time in the cluster. So when a node
// gains Raft leadership the corresponding Autopilot instance should have it's
// Start method called. Then if leadership is lost that node should call the
// Stop method on the Autopilot instance.
type Autopilot struct {
	logger hclog.Logger
	// delegate is used to get information about the system such as Raft server
	// states, known servers etc.
	delegate ApplicationIntegration
	// promoter is used to calculate promotions, demotions and leadership transfers
	// given a particular autopilot State. The interface also contains methods
	// for filling in parts of the autopilot state that the core module doesn't
	// control such as the Ext fields on the Server and State types.
	promoter Promoter
	// raft is an interface that implements all the parts of the Raft library interface
	// that we use. It is an interface to allow for mocking raft during testing.
	raft Raft
	// time is an interface with a single method for getting the current time - `Now`.
	// In some tests this will be the MockTimeProvider which allows tests to be more
	// deterministic but for running systems this should not be overrided from the
	// default which is the runtimeTimeProvider and is a small shim around calling
	// time.Now.
	time timeProvider

	// reconcileInterval is how long between rounds of performing promotions, demotions
	// and leadership transfers.
	reconcileInterval time.Duration

	// updateInterval is the time between the periodic state updates. These periodic
	// state updates take in known servers from the delegate, request Raft stats be
	// fetched and pull in other inputs such as the Raft configuration to create
	// an updated view of the Autopilot State.
	updateInterval time.Duration

	// state is the structure that autopilot uses to make decisions about what to do.
	// This field should be considered immutable and no modifications to an existing
	// state should be made but instead a new state is created and set to this field
	// while holding the stateLock.
	state *State
	// stateLock is meant to only protect the state field. This just prevents
	// the periodic state update and consumers requesting the autopilot state from
	// racing.
	stateLock sync.RWMutex

	// startTime is recorded so that we can make better determinations about server
	// stability during the initial period of time after autopilot first starts.
	// If autopilot has just started the default behavior to check if a server is
	// stable will not work as it will ensure the server has been healthy for
	// the configured server stabilization time. If that configure time is longer
	// than the amount of time autopilot has been running you can run into issues
	// with leadership flapping during some scenarios where a cluster is being
	// brought up.
	startTime time.Time

	// removeDeadCh is used to trigger the running autopilot go routines to
	// find and remove any dead/failed servers
	removeDeadCh chan struct{}

	// reconcileCh is used to trigger an immediate round of reconciliation.
	reconcileCh chan struct{}

	// leaderLock implements a cancellable mutex that will be used to ensure
	// that only one autopilot go routine is the "leader". The leader is
	// the go routine that is currently responsible for updating the
	// autopilot state and performing raft promotions/demotions.
	leaderLock *mutex

	// execution is the information about the most recent autopilot execution.
	// Start will initialize this with the most recent execution and it will
	// be updated by Stop and by the go routines being executed when they are
	// finished.
	execution *execInfo

	// execLock protects access to the execution field
	execLock sync.Mutex
}

// New will create a new Autopilot instance utilizing the given Raft and Delegate.
// If the WithPromoter option is not provided the default StablePromoter will
// be used.
func New(raft Raft, delegate ApplicationIntegration, options ...Option) *Autopilot {
	a := &Autopilot{
		raft:     raft,
		delegate: delegate,
		state:    &State{},
		promoter: DefaultPromoter(),
		logger:   hclog.Default().Named("autopilot"),
		// should this be buffered?
		removeDeadCh:      make(chan struct{}, 1),
		reconcileInterval: DefaultReconcileInterval,
		updateInterval:    DefaultUpdateInterval,
		time:              &runtimeTimeProvider{},
		leaderLock:        newMutex(),
	}

	for _, opt := range options {
		opt(a)
	}

	return a
}

// RemoveDeadServers will trigger an immediate removal of dead/failed servers.
func (a *Autopilot) RemoveDeadServers() {
	select {
	case a.removeDeadCh <- struct{}{}:
	default:
	}
}

// GetState retrieves the current autopilot State
func (a *Autopilot) GetState() *State {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	return a.state
}

// GetServerHealth returns the latest ServerHealth for a given server.
// The returned struct should not be modified or else it will im
func (a *Autopilot) GetServerHealth(id raft.ServerID) *ServerHealth {
	state := a.GetState()

	srv, ok := state.Servers[id]
	if ok {
		return &srv.Health
	}

	return nil
}
