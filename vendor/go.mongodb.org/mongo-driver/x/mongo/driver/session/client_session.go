// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package session // import "go.mongodb.org/mongo-driver/x/mongo/driver/session"

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal/uuid"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// ErrSessionEnded is returned when a client session is used after a call to endSession().
var ErrSessionEnded = errors.New("ended session was used")

// ErrNoTransactStarted is returned if a transaction operation is called when no transaction has started.
var ErrNoTransactStarted = errors.New("no transaction started")

// ErrTransactInProgress is returned if startTransaction() is called when a transaction is in progress.
var ErrTransactInProgress = errors.New("transaction already in progress")

// ErrAbortAfterCommit is returned when abort is called after a commit.
var ErrAbortAfterCommit = errors.New("cannot call abortTransaction after calling commitTransaction")

// ErrAbortTwice is returned if abort is called after transaction is already aborted.
var ErrAbortTwice = errors.New("cannot call abortTransaction twice")

// ErrCommitAfterAbort is returned if commit is called after an abort.
var ErrCommitAfterAbort = errors.New("cannot call commitTransaction after calling abortTransaction")

// ErrUnackWCUnsupported is returned if an unacknowledged write concern is supported for a transaction.
var ErrUnackWCUnsupported = errors.New("transactions do not support unacknowledged write concerns")

// ErrSnapshotTransaction is returned if an transaction is started on a snapshot session.
var ErrSnapshotTransaction = errors.New("transactions are not supported in snapshot sessions")

// TransactionState indicates the state of the transactions FSM.
type TransactionState uint8

// Client Session states
const (
	None TransactionState = iota
	Starting
	InProgress
	Committed
	Aborted
)

// String implements the fmt.Stringer interface.
func (s TransactionState) String() string {
	switch s {
	case None:
		return "none"
	case Starting:
		return "starting"
	case InProgress:
		return "in progress"
	case Committed:
		return "committed"
	case Aborted:
		return "aborted"
	default:
		return "unknown"
	}
}

// LoadBalancedTransactionConnection represents a connection that's pinned by a ClientSession because it's being used
// to execute a transaction when running against a load balancer. This interface is a copy of driver.PinnedConnection
// and exists to be able to pin transactions to a connection without causing an import cycle.
type LoadBalancedTransactionConnection interface {
	// Functions copied over from driver.Connection.
	WriteWireMessage(context.Context, []byte) error
	ReadWireMessage(ctx context.Context) ([]byte, error)
	Description() description.Server
	Close() error
	ID() string
	ServerConnectionID() *int64
	DriverConnectionID() uint64 // TODO(GODRIVER-2824): change type to int64.
	Address() address.Address
	Stale() bool
	OIDCTokenGenID() uint64
	SetOIDCTokenGenID(uint64)

	// Functions copied over from driver.PinnedConnection that are not part of Connection or Expirable.
	PinToCursor() error
	PinToTransaction() error
	UnpinFromCursor() error
	UnpinFromTransaction() error
}

// Client is a session for clients to run commands.
type Client struct {
	*Server
	ClientID       uuid.UUID
	ClusterTime    bson.Raw
	Consistent     bool // causal consistency
	OperationTime  *primitive.Timestamp
	IsImplicit     bool
	Terminated     bool
	RetryingCommit bool
	Committing     bool
	Aborting       bool
	RetryWrite     bool
	RetryRead      bool
	Snapshot       bool

	// options for the current transaction
	// most recently set by transactionopt
	CurrentRc  *readconcern.ReadConcern
	CurrentRp  *readpref.ReadPref
	CurrentWc  *writeconcern.WriteConcern
	CurrentMct *time.Duration

	// default transaction options
	transactionRc            *readconcern.ReadConcern
	transactionRp            *readpref.ReadPref
	transactionWc            *writeconcern.WriteConcern
	transactionMaxCommitTime *time.Duration

	pool             *Pool
	TransactionState TransactionState
	PinnedServer     *description.Server
	RecoveryToken    bson.Raw
	PinnedConnection LoadBalancedTransactionConnection
	SnapshotTime     *primitive.Timestamp
}

func getClusterTime(clusterTime bson.Raw) (uint32, uint32) {
	if clusterTime == nil {
		return 0, 0
	}

	clusterTimeVal, err := clusterTime.LookupErr("$clusterTime")
	if err != nil {
		return 0, 0
	}

	timestampVal, err := bson.Raw(clusterTimeVal.Value).LookupErr("clusterTime")
	if err != nil {
		return 0, 0
	}

	return timestampVal.Timestamp()
}

// MaxClusterTime compares 2 clusterTime documents and returns the document representing the highest cluster time.
func MaxClusterTime(ct1, ct2 bson.Raw) bson.Raw {
	epoch1, ord1 := getClusterTime(ct1)
	epoch2, ord2 := getClusterTime(ct2)

	switch {
	case epoch1 > epoch2:
		return ct1
	case epoch1 < epoch2:
		return ct2
	case ord1 > ord2:
		return ct1
	case ord1 < ord2:
		return ct2
	}

	return ct1
}

// NewImplicitClientSession creates a new implicit client-side session.
func NewImplicitClientSession(pool *Pool, clientID uuid.UUID) *Client {
	// Server-side session checkout for implicit sessions is deferred until after checking out a
	// connection, so don't check out a server-side session right now. This will limit the number of
	// implicit sessions to no greater than an application's maxPoolSize.

	return &Client{
		pool:       pool,
		ClientID:   clientID,
		IsImplicit: true,
	}
}

// NewClientSession creates a new explicit client-side session.
func NewClientSession(pool *Pool, clientID uuid.UUID, opts ...*ClientOptions) (*Client, error) {
	c := &Client{
		pool:     pool,
		ClientID: clientID,
	}

	mergedOpts := mergeClientOptions(opts...)
	if mergedOpts.DefaultReadPreference != nil {
		c.transactionRp = mergedOpts.DefaultReadPreference
	}
	if mergedOpts.DefaultReadConcern != nil {
		c.transactionRc = mergedOpts.DefaultReadConcern
	}
	if mergedOpts.DefaultWriteConcern != nil {
		c.transactionWc = mergedOpts.DefaultWriteConcern
	}
	if mergedOpts.DefaultMaxCommitTime != nil {
		c.transactionMaxCommitTime = mergedOpts.DefaultMaxCommitTime
	}
	if mergedOpts.Snapshot != nil {
		c.Snapshot = *mergedOpts.Snapshot
	}

	// For explicit sessions, the default for causalConsistency is true, unless Snapshot is
	// enabled, then it's false. Set the default and then allow any explicit causalConsistency
	// setting to override it.
	c.Consistent = !c.Snapshot
	if mergedOpts.CausalConsistency != nil {
		c.Consistent = *mergedOpts.CausalConsistency
	}

	if c.Consistent && c.Snapshot {
		return nil, errors.New("causal consistency and snapshot cannot both be set for a session")
	}

	if err := c.SetServer(); err != nil {
		return nil, err
	}

	return c, nil
}

// SetServer will check out a session from the client session pool.
func (c *Client) SetServer() error {
	var err error
	c.Server, err = c.pool.GetSession()
	return err
}

// AdvanceClusterTime updates the session's cluster time.
func (c *Client) AdvanceClusterTime(clusterTime bson.Raw) error {
	if c.Terminated {
		return ErrSessionEnded
	}
	c.ClusterTime = MaxClusterTime(c.ClusterTime, clusterTime)
	return nil
}

// AdvanceOperationTime updates the session's operation time.
func (c *Client) AdvanceOperationTime(opTime *primitive.Timestamp) error {
	if c.Terminated {
		return ErrSessionEnded
	}

	if c.OperationTime == nil {
		c.OperationTime = opTime
		return nil
	}

	if opTime.T > c.OperationTime.T {
		c.OperationTime = opTime
	} else if (opTime.T == c.OperationTime.T) && (opTime.I > c.OperationTime.I) {
		c.OperationTime = opTime
	}

	return nil
}

// UpdateUseTime sets the session's last used time to the current time. This must be called whenever the session is
// used to send a command to the server to ensure that the session is not prematurely marked expired in the driver's
// session pool. If the session has already been ended, this method will return ErrSessionEnded.
func (c *Client) UpdateUseTime() error {
	if c.Terminated {
		return ErrSessionEnded
	}
	c.updateUseTime()
	return nil
}

// UpdateRecoveryToken updates the session's recovery token from the server response.
func (c *Client) UpdateRecoveryToken(response bson.Raw) {
	if c == nil {
		return
	}

	token, err := response.LookupErr("recoveryToken")
	if err != nil {
		return
	}

	c.RecoveryToken = token.Document()
}

// UpdateSnapshotTime updates the session's value for the atClusterTime field of ReadConcern.
func (c *Client) UpdateSnapshotTime(response bsoncore.Document) {
	if c == nil {
		return
	}

	subDoc := response
	if cur, ok := response.Lookup("cursor").DocumentOK(); ok {
		subDoc = cur
	}

	ssTimeElem, err := subDoc.LookupErr("atClusterTime")
	if err != nil {
		// atClusterTime not included by the server
		return
	}

	t, i := ssTimeElem.Timestamp()
	c.SnapshotTime = &primitive.Timestamp{
		T: t,
		I: i,
	}
}

// ClearPinnedResources clears the pinned server and/or connection associated with the session.
func (c *Client) ClearPinnedResources() error {
	if c == nil {
		return nil
	}

	c.PinnedServer = nil
	if c.PinnedConnection != nil {
		if err := c.PinnedConnection.UnpinFromTransaction(); err != nil {
			return err
		}
		if err := c.PinnedConnection.Close(); err != nil {
			return err
		}
	}
	c.PinnedConnection = nil
	return nil
}

// unpinConnection gracefully unpins the connection associated with the session
// if there is one. This is done via the pinned connection's
// UnpinFromTransaction function.
func (c *Client) unpinConnection() error {
	if c == nil || c.PinnedConnection == nil {
		return nil
	}

	err := c.PinnedConnection.UnpinFromTransaction()
	closeErr := c.PinnedConnection.Close()
	if err == nil && closeErr != nil {
		err = closeErr
	}
	c.PinnedConnection = nil
	return err
}

// EndSession ends the session.
func (c *Client) EndSession() {
	if c.Terminated {
		return
	}
	c.Terminated = true

	// Ignore the error when unpinning the connection because we can't do
	// anything about it if it doesn't work. Typically the only errors that can
	// happen here indicate that something went wrong with the connection state,
	// like it wasn't marked as pinned or attempted to return to the wrong pool.
	_ = c.unpinConnection()
	c.pool.ReturnSession(c.Server)
}

// TransactionInProgress returns true if the client session is in an active transaction.
func (c *Client) TransactionInProgress() bool {
	return c.TransactionState == InProgress
}

// TransactionStarting returns true if the client session is starting a transaction.
func (c *Client) TransactionStarting() bool {
	return c.TransactionState == Starting
}

// TransactionRunning returns true if the client session has started the transaction
// and it hasn't been committed or aborted
func (c *Client) TransactionRunning() bool {
	return c != nil && (c.TransactionState == Starting || c.TransactionState == InProgress)
}

// TransactionCommitted returns true of the client session just committed a transaction.
func (c *Client) TransactionCommitted() bool {
	return c.TransactionState == Committed
}

// CheckStartTransaction checks to see if allowed to start transaction and returns
// an error if not allowed
func (c *Client) CheckStartTransaction() error {
	if c.TransactionState == InProgress || c.TransactionState == Starting {
		return ErrTransactInProgress
	}
	if c.Snapshot {
		return ErrSnapshotTransaction
	}
	return nil
}

// StartTransaction initializes the transaction options and advances the state machine.
// It does not contact the server to start the transaction.
func (c *Client) StartTransaction(opts *TransactionOptions) error {
	err := c.CheckStartTransaction()
	if err != nil {
		return err
	}

	c.IncrementTxnNumber()
	c.RetryingCommit = false

	if opts != nil {
		c.CurrentRc = opts.ReadConcern
		c.CurrentRp = opts.ReadPreference
		c.CurrentWc = opts.WriteConcern
		c.CurrentMct = opts.MaxCommitTime
	}

	if c.CurrentRc == nil {
		c.CurrentRc = c.transactionRc
	}

	if c.CurrentRp == nil {
		c.CurrentRp = c.transactionRp
	}

	if c.CurrentWc == nil {
		c.CurrentWc = c.transactionWc
	}

	if c.CurrentMct == nil {
		c.CurrentMct = c.transactionMaxCommitTime
	}

	if !writeconcern.AckWrite(c.CurrentWc) {
		_ = c.clearTransactionOpts()
		return ErrUnackWCUnsupported
	}

	c.TransactionState = Starting
	return c.ClearPinnedResources()
}

// CheckCommitTransaction checks to see if allowed to commit transaction and returns
// an error if not allowed.
func (c *Client) CheckCommitTransaction() error {
	if c.TransactionState == None {
		return ErrNoTransactStarted
	} else if c.TransactionState == Aborted {
		return ErrCommitAfterAbort
	}
	return nil
}

// CommitTransaction updates the state for a successfully committed transaction and returns
// an error if not permissible.  It does not actually perform the commit.
func (c *Client) CommitTransaction() error {
	err := c.CheckCommitTransaction()
	if err != nil {
		return err
	}
	c.TransactionState = Committed
	return nil
}

// UpdateCommitTransactionWriteConcern will set the write concern to majority and potentially set  a
// w timeout of 10 seconds. This should be called after a commit transaction operation fails with a
// retryable error or after a successful commit transaction operation.
func (c *Client) UpdateCommitTransactionWriteConcern() {
	wc := c.CurrentWc
	timeout := 10 * time.Second
	if wc != nil && wc.GetWTimeout() != 0 {
		timeout = wc.GetWTimeout()
	}
	c.CurrentWc = wc.WithOptions(writeconcern.WMajority(), writeconcern.WTimeout(timeout))
}

// CheckAbortTransaction checks to see if allowed to abort transaction and returns
// an error if not allowed.
func (c *Client) CheckAbortTransaction() error {
	switch {
	case c.TransactionState == None:
		return ErrNoTransactStarted
	case c.TransactionState == Committed:
		return ErrAbortAfterCommit
	case c.TransactionState == Aborted:
		return ErrAbortTwice
	}
	return nil
}

// AbortTransaction updates the state for a successfully aborted transaction and returns
// an error if not permissible.  It does not actually perform the abort.
func (c *Client) AbortTransaction() error {
	err := c.CheckAbortTransaction()
	if err != nil {
		return err
	}
	c.TransactionState = Aborted
	return c.clearTransactionOpts()
}

// StartCommand updates the session's internal state at the beginning of an operation. This must be called before
// server selection is done for the operation as the session's state can impact the result of that process.
func (c *Client) StartCommand() error {
	if c == nil {
		return nil
	}

	// If we're executing the first operation using this session after a transaction, we must ensure that the session
	// is not pinned to any resources.
	if !c.TransactionRunning() && !c.Committing && !c.Aborting {
		return c.ClearPinnedResources()
	}
	return nil
}

// ApplyCommand advances the state machine upon command execution. This must be called after server selection is
// complete.
func (c *Client) ApplyCommand(desc description.Server) error {
	if c.Committing {
		// Do not change state if committing after already committed
		return nil
	}
	if c.TransactionState == Starting {
		c.TransactionState = InProgress
		// If this is in a transaction and the server is a mongos, pin it
		if desc.Kind == description.Mongos {
			c.PinnedServer = &desc
		}
	} else if c.TransactionState == Committed || c.TransactionState == Aborted {
		c.TransactionState = None
		return c.clearTransactionOpts()
	}

	return nil
}

func (c *Client) clearTransactionOpts() error {
	c.RetryingCommit = false
	c.Aborting = false
	c.Committing = false
	c.CurrentWc = nil
	c.CurrentRp = nil
	c.CurrentRc = nil
	c.RecoveryToken = nil

	return c.ClearPinnedResources()
}
