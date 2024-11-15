package radix

import (
	"context"
	"math/rand"
	"net"

	"github.com/mediocregopher/radix/v4/internal/proc"
)

type rawAddr struct {
	network, addr string
}

var _ net.Addr = rawAddr{}

func (a rawAddr) Network() string { return a.network }

func (a rawAddr) String() string { return a.addr }

// Client describes an entity which can carry out Actions on a single redis
// instance. Conn and Pool are Clients.
//
// Implementations of Client are expected to be thread-safe.
type Client interface {
	// Addr returns the address of the redis instance which the Client was
	// initialized against.
	Addr() net.Addr

	// Do performs an Action on a Conn connected to the redis instance.
	Do(context.Context, Action) error

	// Once Close() is called all future method calls on the Client will return
	// an error
	Close() error
}

// ReplicaSet holds the Clients of a redis replica set, consisting of a single
// primary (read+write) instance and zero or more secondary (read-only)
// instances.
type ReplicaSet struct {
	Primary     Client
	Secondaries []Client
}

// MultiClient wraps one or more underlying Clients for different redis
// instances. MultiClient methods are thread-safe and may return the same Client
// instance to different callers at the same time. All returned Clients should
// _not_ have Close called on them.
//
// If the topology backing a MultiClient changes (e.g. a failover occurs) while
// the Clients it returned are still being used then those Clients may return
// errors related to that change.
//
// Sentinel and Cluster are both MultiClients.
type MultiClient interface {
	// Do performs an Action on a Conn from a primary instance.
	Do(context.Context, Action) error

	// DoSecondary performs the Action on a Conn from a secondary instance. If
	// no secondary instance is available then this is equivalent to Do.
	DoSecondary(context.Context, Action) error

	// Clients returns all Clients held by MultiClient, formatted as a mapping
	// of primary redis instance address to a ReplicaSet instance for that
	// primary.
	Clients() (map[string]ReplicaSet, error)

	// Once Close() is called all future method calls on the Client will return
	// an error
	Close() error
}

type replicaSetMultiClient struct {
	ReplicaSet
	proc *proc.Proc
}

// NewMultiClient wraps a ReplicaSet such that it implements MultiClient.
func NewMultiClient(rs ReplicaSet) MultiClient {
	return &replicaSetMultiClient{
		ReplicaSet: rs,
		proc:       proc.New(),
	}
}

func (r *replicaSetMultiClient) Do(ctx context.Context, a Action) error {
	return r.proc.WithRLock(func() error {
		return r.Primary.Do(ctx, a)
	})
}

func (r *replicaSetMultiClient) DoSecondary(ctx context.Context, a Action) error {
	return r.proc.WithRLock(func() error {
		if len(r.Secondaries) == 0 {
			return r.Do(ctx, a)
		}
		return r.Secondaries[rand.Intn(len(r.Secondaries))].Do(ctx, a)
	})
}

func (r *replicaSetMultiClient) Clients() (map[string]ReplicaSet, error) {
	m := make(map[string]ReplicaSet, 1)
	err := r.proc.WithRLock(func() error {
		m[r.Primary.Addr().String()] = r.ReplicaSet
		return nil
	})
	return m, err
}

func (r *replicaSetMultiClient) Close() error {
	return r.proc.Close(func() error {
		err := r.Primary.Close()
		for _, secondary := range r.Secondaries {
			if secErr := secondary.Close(); err == nil {
				err = secErr
			}
		}
		return err
	})
}
