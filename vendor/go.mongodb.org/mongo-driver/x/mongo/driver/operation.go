// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/internal/handshake"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

const defaultLocalThreshold = 15 * time.Millisecond

var (
	// ErrNoDocCommandResponse occurs when the server indicated a response existed, but none was found.
	ErrNoDocCommandResponse = errors.New("command returned no documents")
	// ErrMultiDocCommandResponse occurs when the server sent multiple documents in response to a command.
	ErrMultiDocCommandResponse = errors.New("command returned multiple documents")
	// ErrReplyDocumentMismatch occurs when the number of documents returned in an OP_QUERY does not match the numberReturned field.
	ErrReplyDocumentMismatch = errors.New("number of documents returned does not match numberReturned field")
	// ErrNonPrimaryReadPref is returned when a read is attempted in a transaction with a non-primary read preference.
	ErrNonPrimaryReadPref = errors.New("read preference in a transaction must be primary")
	// errDatabaseNameEmpty occurs when a database name is not provided.
	errDatabaseNameEmpty = errors.New("database name cannot be empty")
)

const (
	// maximum BSON object size when client side encryption is enabled
	cryptMaxBsonObjectSize uint32 = 2097152
	// minimum wire version necessary to use automatic encryption
	cryptMinWireVersion int32 = 8
	// minimum wire version necessary to use read snapshots
	readSnapshotMinWireVersion int32 = 13
)

// RetryablePoolError is a connection pool error that can be retried while executing an operation.
type RetryablePoolError interface {
	Retryable() bool
}

// labeledError is an error that can have error labels added to it.
type labeledError interface {
	error
	HasErrorLabel(string) bool
}

// InvalidOperationError is returned from Validate and indicates that a required field is missing
// from an instance of Operation.
type InvalidOperationError struct{ MissingField string }

func (err InvalidOperationError) Error() string {
	return "the " + err.MissingField + " field must be set on Operation"
}

// opReply stores information returned in an OP_REPLY response from the server.
// The err field stores any error that occurred when decoding or validating the OP_REPLY response.
type opReply struct {
	responseFlags wiremessage.ReplyFlag
	cursorID      int64
	startingFrom  int32
	numReturned   int32
	documents     []bsoncore.Document
	err           error
}

// startedInformation keeps track of all of the information necessary for monitoring started events.
type startedInformation struct {
	cmd                      bsoncore.Document
	requestID                int32
	cmdName                  string
	documentSequenceIncluded bool
	connID                   string
	driverConnectionID       uint64 // TODO(GODRIVER-2824): change type to int64.
	serverConnID             *int64
	redacted                 bool
	serviceID                *primitive.ObjectID
	serverAddress            address.Address
}

// finishedInformation keeps track of all of the information necessary for monitoring success and failure events.
type finishedInformation struct {
	cmdName            string
	requestID          int32
	response           bsoncore.Document
	cmdErr             error
	connID             string
	driverConnectionID uint64 // TODO(GODRIVER-2824): change type to int64.
	serverConnID       *int64
	redacted           bool
	serviceID          *primitive.ObjectID
	serverAddress      address.Address
	duration           time.Duration
}

// convertInt64PtrToInt32Ptr will convert an int64 pointer reference to an int32 pointer
// reference. If the int64 value cannot be converted to int32 without causing
// an overflow, then this function will return nil.
func convertInt64PtrToInt32Ptr(i64 *int64) *int32 {
	if i64 == nil {
		return nil
	}

	if *i64 > math.MaxInt32 || *i64 < math.MinInt32 {
		return nil
	}

	i32 := int32(*i64)
	return &i32
}

// success returns true if there was no command error or the command error is a
// "WriteCommandError". Commands that executed on the server and return a status
// of { ok: 1.0 } are considered successful commands and MUST generate a
// CommandSucceededEvent and "command succeeded" log message. Commands that have
// write errors are included since the actual command did succeed, only writes
// failed.
func (info finishedInformation) success() bool {
	if _, ok := info.cmdErr.(WriteCommandError); ok {
		return true
	}

	return info.cmdErr == nil
}

// ResponseInfo contains the context required to parse a server response.
type ResponseInfo struct {
	ServerResponse        bsoncore.Document
	Server                Server
	Connection            Connection
	ConnectionDescription description.Server
	CurrentIndex          int
}

func redactStartedInformationCmd(op Operation, info startedInformation) bson.Raw {
	var cmdCopy bson.Raw

	// Make a copy of the command. Redact if the command is security
	// sensitive and cannot be monitored. If there was a type 1 payload for
	// the current batch, convert it to a BSON array
	if !info.redacted {
		cmdCopy = make([]byte, len(info.cmd))
		copy(cmdCopy, info.cmd)

		if info.documentSequenceIncluded {
			// remove 0 byte at end
			cmdCopy = cmdCopy[:len(info.cmd)-1]
			cmdCopy = op.addBatchArray(cmdCopy)

			// add back 0 byte and update length
			cmdCopy, _ = bsoncore.AppendDocumentEnd(cmdCopy, 0)
		}
	}

	return cmdCopy
}

func redactFinishedInformationResponse(info finishedInformation) bson.Raw {
	if !info.redacted {
		return bson.Raw(info.response)
	}

	return bson.Raw{}
}

// Operation is used to execute an operation. It contains all of the common code required to
// select a server, transform an operation into a command, write the command to a connection from
// the selected server, read a response from that connection, process the response, and potentially
// retry.
//
// The required fields are Database, CommandFn, and Deployment. All other fields are optional.
//
// While an Operation can be constructed manually, drivergen should be used to generate an
// implementation of an operation instead. This will ensure that there are helpers for constructing
// the operation and that this type isn't configured incorrectly.
type Operation struct {
	// CommandFn is used to create the command that will be wrapped in a wire message and sent to
	// the server. This function should only add the elements of the command and not start or end
	// the enclosing BSON document. Per the command API, the first element must be the name of the
	// command to run. This field is required.
	CommandFn func(dst []byte, desc description.SelectedServer) ([]byte, error)

	// Database is the database that the command will be run against. This field is required.
	Database string

	// Deployment is the MongoDB Deployment to use. While most of the time this will be multiple
	// servers, commands that need to run against a single, preselected server can use the
	// SingleServerDeployment type. Commands that need to run on a preselected connection can use
	// the SingleConnectionDeployment type.
	Deployment Deployment

	// ProcessResponseFn is called after a response to the command is returned. The server is
	// provided for types like Cursor that are required to run subsequent commands using the same
	// server.
	ProcessResponseFn func(ResponseInfo) error

	// Selector is the server selector that's used during both initial server selection and
	// subsequent selection for retries. Depending on the Deployment implementation, the
	// SelectServer method may not actually be called.
	Selector description.ServerSelector

	// ReadPreference is the read preference that will be attached to the command. If this field is
	// not specified a default read preference of primary will be used.
	ReadPreference *readpref.ReadPref

	// ReadConcern is the read concern used when running read commands. This field should not be set
	// for write operations. If this field is set, it will be encoded onto the commands sent to the
	// server.
	ReadConcern *readconcern.ReadConcern

	// MinimumReadConcernWireVersion specifies the minimum wire version to add the read concern to
	// the command being executed.
	MinimumReadConcernWireVersion int32

	// WriteConcern is the write concern used when running write commands. This field should not be
	// set for read operations. If this field is set, it will be encoded onto the commands sent to
	// the server.
	WriteConcern *writeconcern.WriteConcern

	// MinimumWriteConcernWireVersion specifies the minimum wire version to add the write concern to
	// the command being executed.
	MinimumWriteConcernWireVersion int32

	// Client is the session used with this operation. This can be either an implicit or explicit
	// session. If the server selected does not support sessions and Client is specified the
	// behavior depends on the session type. If the session is implicit, the session fields will not
	// be encoded onto the command. If the session is explicit, an error will be returned. The
	// caller is responsible for ensuring that this field is nil if the Deployment does not support
	// sessions.
	Client *session.Client

	// Clock is a cluster clock, different from the one contained within a session.Client. This
	// allows updating cluster times for a global cluster clock while allowing individual session's
	// cluster clocks to be only updated as far as the last command that's been run.
	Clock *session.ClusterClock

	// RetryMode specifies how to retry. There are three modes that enable retry: RetryOnce,
	// RetryOncePerCommand, and RetryContext. For more information about what these modes do, please
	// refer to their definitions. Both RetryMode and Type must be set for retryability to be enabled.
	// If Timeout is set on the Client, the operation will automatically retry as many times as
	// possible unless RetryNone is used.
	RetryMode *RetryMode

	// Type specifies the kind of operation this is. There is only one mode that enables retry: Write.
	// For more information about what this mode does, please refer to it's definition. Both Type and
	// RetryMode must be set for retryability to be enabled.
	Type Type

	// Batches contains the documents that are split when executing a write command that potentially
	// has more documents than can fit in a single command. This should only be specified for
	// commands that are batch compatible. For more information, please refer to the definition of
	// Batches.
	Batches *Batches

	// Legacy sets the legacy type for this operation. There are only 3 types that require legacy
	// support: find, getMore, and killCursors. For more information about LegacyOperationKind,
	// please refer to it's definition.
	Legacy LegacyOperationKind

	// CommandMonitor specifies the monitor to use for APM events. If this field is not set,
	// no events will be reported.
	CommandMonitor *event.CommandMonitor

	// Crypt specifies a Crypt object to use for automatic client side encryption and decryption.
	Crypt Crypt

	// ServerAPI specifies options used to configure the API version sent to the server.
	ServerAPI *ServerAPIOptions

	// IsOutputAggregate specifies whether this operation is an aggregate with an output stage. If true,
	// read preference will not be added to the command on wire versions < 13.
	IsOutputAggregate bool

	// MaxTime specifies the maximum amount of time to allow the operation to run on the server.
	MaxTime *time.Duration

	// Timeout is the amount of time that this operation can execute before returning an error. The default value
	// nil, which means that the timeout of the operation's caller will be used.
	Timeout *time.Duration

	Logger *logger.Logger

	// Name is the name of the operation. This is used when serializing
	// OP_MSG as well as for logging server selection data.
	Name string

	// OmitCSOTMaxTimeMS omits the automatically-calculated "maxTimeMS" from the
	// command when CSOT is enabled. It does not effect "maxTimeMS" set by
	// [Operation.MaxTime].
	OmitCSOTMaxTimeMS bool

	// Authenticator is the authenticator to use for this operation when a reauthentication is
	// required.
	Authenticator Authenticator

	// omitReadPreference is a boolean that indicates whether to omit the
	// read preference from the command. This omition includes the case
	// where a default read preference is used when the operation
	// ReadPreference is not specified.
	omitReadPreference bool
}

// shouldEncrypt returns true if this operation should automatically be encrypted.
func (op Operation) shouldEncrypt() bool {
	return op.Crypt != nil && !op.Crypt.BypassAutoEncryption()
}

// filterDeprioritizedServers will filter out the server candidates that have
// been deprioritized by the operation due to failure.
//
// The server selector should try to select a server that is not in the
// deprioritization list. However, if this is not possible (e.g. there are no
// other healthy servers in the cluster), the selector may return a
// deprioritized server.
func filterDeprioritizedServers(candidates, deprioritized []description.Server) []description.Server {
	if len(deprioritized) == 0 {
		return candidates
	}

	dpaSet := make(map[address.Address]*description.Server)
	for i, srv := range deprioritized {
		dpaSet[srv.Addr] = &deprioritized[i]
	}

	allowed := []description.Server{}

	// Iterate over the candidates and append them to the allowdIndexes slice if
	// they are not in the deprioritizedServers list.
	for _, candidate := range candidates {
		if srv, ok := dpaSet[candidate.Addr]; !ok || !srv.Equal(candidate) {
			allowed = append(allowed, candidate)
		}
	}

	// If nothing is allowed, then all available servers must have been
	// deprioritized. In this case, return the candidates list as-is so that the
	// selector can find a suitable server
	if len(allowed) == 0 {
		return candidates
	}

	return allowed
}

// opServerSelector is a wrapper for the server selector that is assigned to the
// operation. The purpose of this wrapper is to filter candidates with
// operation-specific logic, such as deprioritizing failing servers.
type opServerSelector struct {
	selector             description.ServerSelector
	deprioritizedServers []description.Server
}

// SelectServer will filter candidates with operation-specific logic before
// passing them onto the user-defined or default selector.
func (oss *opServerSelector) SelectServer(
	topo description.Topology,
	candidates []description.Server,
) ([]description.Server, error) {
	selectedServers, err := oss.selector.SelectServer(topo, candidates)
	if err != nil {
		return nil, err
	}

	filteredServers := filterDeprioritizedServers(selectedServers, oss.deprioritizedServers)

	return filteredServers, nil
}

// selectServer handles performing server selection for an operation.
func (op Operation) selectServer(
	ctx context.Context,
	requestID int32,
	deprioritized []description.Server,
) (Server, error) {
	if err := op.Validate(); err != nil {
		return nil, err
	}

	selector := op.Selector
	if selector == nil {
		rp := op.ReadPreference
		if rp == nil {
			rp = readpref.Primary()
		}
		selector = description.CompositeSelector([]description.ServerSelector{
			description.ReadPrefSelector(rp),
			description.LatencySelector(defaultLocalThreshold),
		})
	}

	oss := &opServerSelector{
		selector:             selector,
		deprioritizedServers: deprioritized,
	}

	ctx = logger.WithOperationName(ctx, op.Name)
	ctx = logger.WithOperationID(ctx, requestID)

	return op.Deployment.SelectServer(ctx, oss)
}

// getServerAndConnection should be used to retrieve a Server and Connection to execute an operation.
func (op Operation) getServerAndConnection(
	ctx context.Context,
	requestID int32,
	deprioritized []description.Server,
) (Server, Connection, error) {
	server, err := op.selectServer(ctx, requestID, deprioritized)
	if err != nil {
		if op.Client != nil &&
			!(op.Client.Committing || op.Client.Aborting) && op.Client.TransactionRunning() {
			err = Error{
				Message: err.Error(),
				Labels:  []string{TransientTransactionError},
				Wrapped: err,
			}
		}
		return nil, nil, err
	}

	// If the provided client session has a pinned connection, it should be used for the operation because this
	// indicates that we're in a transaction and the target server is behind a load balancer.
	if op.Client != nil && op.Client.PinnedConnection != nil {
		return server, op.Client.PinnedConnection, nil
	}

	// Otherwise, default to checking out a connection from the server's pool.
	conn, err := server.Connection(ctx)
	if err != nil {
		return nil, nil, err
	}

	// If we're in load balanced mode and this is the first operation in a transaction, pin the session to a connection.
	if conn.Description().LoadBalanced() && op.Client != nil && op.Client.TransactionStarting() {
		pinnedConn, ok := conn.(PinnedConnection)
		if !ok {
			// Close the original connection to avoid a leak.
			_ = conn.Close()
			return nil, nil, fmt.Errorf("expected Connection used to start a transaction to be a PinnedConnection, but got %T", conn)
		}
		if err := pinnedConn.PinToTransaction(); err != nil {
			// Close the original connection to avoid a leak.
			_ = conn.Close()
			return nil, nil, fmt.Errorf("error incrementing connection reference count when starting a transaction: %w", err)
		}
		op.Client.PinnedConnection = pinnedConn
	}

	return server, conn, nil
}

// Validate validates this operation, ensuring the fields are set properly.
func (op Operation) Validate() error {
	if op.CommandFn == nil {
		return InvalidOperationError{MissingField: "CommandFn"}
	}
	if op.Deployment == nil {
		return InvalidOperationError{MissingField: "Deployment"}
	}
	if op.Database == "" {
		return errDatabaseNameEmpty
	}
	if op.Client != nil && !writeconcern.AckWrite(op.WriteConcern) {
		return errors.New("session provided for an unacknowledged write")
	}
	return nil
}

var memoryPool = sync.Pool{
	New: func() interface{} {
		// Start with 1kb buffers.
		b := make([]byte, 1024)
		// Return a pointer as the static analysis tool suggests.
		return &b
	},
}

// Execute runs this operation.
func (op Operation) Execute(ctx context.Context) error {
	err := op.Validate()
	if err != nil {
		return err
	}

	// If op.Timeout is set, and context is not already a Timeout context, honor
	// op.Timeout in new Timeout context for operation execution.
	if op.Timeout != nil && !csot.IsTimeoutContext(ctx) {
		newCtx, cancelFunc := csot.MakeTimeoutContext(ctx, *op.Timeout)
		// Redefine ctx to be the new timeout-derived context.
		ctx = newCtx
		// Cancel the timeout-derived context at the end of Execute to avoid a context leak.
		defer cancelFunc()
	}

	if op.Client != nil {
		if err := op.Client.StartCommand(); err != nil {
			return err
		}
	}

	var retries int
	if op.RetryMode != nil {
		switch op.Type {
		case Write:
			if op.Client == nil {
				break
			}
			switch *op.RetryMode {
			case RetryOnce, RetryOncePerCommand:
				retries = 1
			case RetryContext:
				retries = -1
			}
		case Read:
			switch *op.RetryMode {
			case RetryOnce, RetryOncePerCommand:
				retries = 1
			case RetryContext:
				retries = -1
			}
		}
	}
	// If context is a Timeout context, automatically set retries to -1 (infinite) if retrying is
	// enabled.
	retryEnabled := op.RetryMode != nil && op.RetryMode.Enabled()
	if csot.IsTimeoutContext(ctx) && retryEnabled {
		retries = -1
	}

	var srvr Server
	var conn Connection
	var res bsoncore.Document
	var operationErr WriteCommandError
	var prevErr error
	var prevIndefiniteErr error
	batching := op.Batches.Valid()
	retrySupported := false
	first := true
	currIndex := 0

	// deprioritizedServers are a running list of servers that should be
	// deprioritized during server selection. Per the specifications, we should
	// only ever deprioritize the "previous server".
	var deprioritizedServers []description.Server

	// resetForRetry records the error that caused the retry, decrements retries, and resets the
	// retry loop variables to request a new server and a new connection for the next attempt.
	resetForRetry := func(err error) {
		retries--
		prevErr = err

		// Set the previous indefinite error to be returned in any case where a retryable write error does not have a
		// NoWritesPerfomed label (the definite case).
		if err, ok := err.(labeledError); ok {
			// If the "prevIndefiniteErr" is nil, then the current error is the first error encountered
			// during the retry attempt cycle. We must persist the first error in the case where all
			// following errors are labeled "NoWritesPerformed", which would otherwise raise nil as the
			// error.
			if prevIndefiniteErr == nil {
				prevIndefiniteErr = err
			}

			// If the error is not labeled NoWritesPerformed and is retryable, then set the previous
			// indefinite error to be the current error.
			if !err.HasErrorLabel(NoWritesPerformed) && err.HasErrorLabel(RetryableWriteError) {
				prevIndefiniteErr = err
			}
		}

		// If we got a connection, close it immediately to release pool resources
		// for subsequent retries.
		if conn != nil {
			// If we are dealing with a sharded cluster, then mark the failed server
			// as "deprioritized".
			if desc := conn.Description; desc != nil && op.Deployment.Kind() == description.Sharded {
				deprioritizedServers = []description.Server{conn.Description()}
			}

			conn.Close()
		}

		// Set the server and connection to nil to request a new server and connection.
		srvr = nil
		conn = nil
	}

	wm := memoryPool.Get().(*[]byte)
	defer func() {
		// Proper usage of a sync.Pool requires each entry to have approximately the same memory
		// cost. To obtain this property when the stored type contains a variably-sized buffer,
		// we add a hard limit on the maximum buffer to place back in the pool. We limit the
		// size to 16MiB because that's the maximum wire message size supported by MongoDB.
		//
		// Comment copied from https://cs.opensource.google/go/go/+/refs/tags/go1.19:src/fmt/print.go;l=147
		//
		// Recycle byte slices that are smaller than 16MiB and at least half occupied.
		if c := cap(*wm); c < 16*1024*1024 && c/2 < len(*wm) {
			memoryPool.Put(wm)
		}
	}()
	for {
		// If we're starting a retry and the error from the previous try was
		// a context canceled or deadline exceeded error, stop retrying and
		// return that error.
		if errors.Is(prevErr, context.Canceled) || errors.Is(prevErr, context.DeadlineExceeded) {
			return prevErr
		}

		requestID := wiremessage.NextRequestID()

		// If the server or connection are nil, try to select a new server and get a new connection.
		if srvr == nil || conn == nil {
			srvr, conn, err = op.getServerAndConnection(ctx, requestID, deprioritizedServers)
			if err != nil {
				// If the returned error is retryable and there are retries remaining (negative
				// retries means retry indefinitely), then retry the operation. Set the server
				// and connection to nil to request a new server and connection.
				if rerr, ok := err.(RetryablePoolError); ok && rerr.Retryable() && retries != 0 {
					resetForRetry(err)
					continue
				}

				// If this is a retry and there's an error from a previous attempt, return the previous
				// error instead of the current connection error.
				if prevErr != nil {
					return prevErr
				}
				return err
			}
			defer conn.Close()

			// Set the server if it has not already been set and the session type is implicit. This will
			// limit the number of implicit sessions to no greater than an application's maxPoolSize
			// (ignoring operations that hold on to the session like cursors).
			if op.Client != nil && op.Client.Server == nil && op.Client.IsImplicit {
				if op.Client.Terminated {
					return fmt.Errorf("unexpected nil session for a terminated implicit session")
				}
				if err := op.Client.SetServer(); err != nil {
					return err
				}
			}
		}

		// Run steps that must only be run on the first attempt, but not again for retries.
		if first {
			// Determine if retries are supported for the current operation on the current server
			// description. Per the retryable writes specification, only determine this for the
			// first server selected:
			//
			//   If the server selected for the first attempt of a retryable write operation does
			//   not support retryable writes, drivers MUST execute the write as if retryable writes
			//   were not enabled.
			retrySupported = op.retryable(conn.Description())

			// If retries are supported for the current operation on the current server description,
			// client retries are enabled, the operation type is write, and we haven't incremented
			// the txn number yet, enable retry writes on the session and increment the txn number.
			// Calling IncrementTxnNumber() for server descriptions or topologies that do not
			// support retries (e.g. standalone topologies) will cause server errors. Only do this
			// check for the first attempt to keep retried writes in the same transaction.
			if retrySupported && op.RetryMode != nil && op.Type == Write && op.Client != nil {
				op.Client.RetryWrite = false
				if op.RetryMode.Enabled() {
					op.Client.RetryWrite = true
					if !op.Client.Committing && !op.Client.Aborting {
						op.Client.IncrementTxnNumber()
					}
				}
			}

			first = false
		}

		maxTimeMS, err := op.calculateMaxTimeMS(ctx, srvr.RTTMonitor())
		if err != nil {
			return err
		}

		// Set maxTimeMS to 0 if connected to mongocryptd to avoid appending the field. The final
		// encrypted command may contain multiple maxTimeMS fields otherwise.
		if conn.Description().IsCryptd {
			maxTimeMS = 0
		}

		desc := description.SelectedServer{Server: conn.Description(), Kind: op.Deployment.Kind()}

		if batching {
			targetBatchSize := desc.MaxDocumentSize
			maxDocSize := desc.MaxDocumentSize
			if op.shouldEncrypt() {
				// For client-side encryption, we want the batch to be split at 2 MiB instead of 16MiB.
				// If there's only one document in the batch, it can be up to 16MiB, so we set target batch size to
				// 2MiB but max document size to 16MiB. This will allow the AdvanceBatch call to create a batch
				// with a single large document.
				targetBatchSize = cryptMaxBsonObjectSize
			}

			err = op.Batches.AdvanceBatch(int(desc.MaxBatchCount), int(targetBatchSize), int(maxDocSize))
			if err != nil {
				// TODO(GODRIVER-982): Should we also be returning operationErr?
				return err
			}
		}

		var startedInfo startedInformation
		*wm, startedInfo, err = op.createWireMessage(ctx, maxTimeMS, (*wm)[:0], desc, conn, requestID)

		if err != nil {
			return err
		}

		// set extra data and send event if possible
		startedInfo.connID = conn.ID()
		startedInfo.driverConnectionID = conn.DriverConnectionID()
		startedInfo.cmdName = op.getCommandName(startedInfo.cmd)

		// If the command name does not match the operation name, update
		// the operation name as a sanity check. It's more correct to
		// be aligned with the data passed to the server via the
		// wire message.
		if startedInfo.cmdName != op.Name {
			op.Name = startedInfo.cmdName
		}

		startedInfo.redacted = op.redactCommand(startedInfo.cmdName, startedInfo.cmd)
		startedInfo.serviceID = conn.Description().ServiceID
		startedInfo.serverConnID = conn.ServerConnectionID()
		startedInfo.serverAddress = conn.Description().Addr

		op.publishStartedEvent(ctx, startedInfo)

		// get the moreToCome flag information before we compress
		moreToCome := wiremessage.IsMsgMoreToCome(*wm)

		// compress wiremessage if allowed
		if compressor, ok := conn.(Compressor); ok && op.canCompress(startedInfo.cmdName) {
			b := memoryPool.Get().(*[]byte)
			*b, err = compressor.CompressWireMessage(*wm, (*b)[:0])
			memoryPool.Put(wm)
			wm = b
			if err != nil {
				return err
			}
		}

		finishedInfo := finishedInformation{
			cmdName:            startedInfo.cmdName,
			driverConnectionID: startedInfo.driverConnectionID,
			requestID:          startedInfo.requestID,
			connID:             startedInfo.connID,
			serverConnID:       startedInfo.serverConnID,
			redacted:           startedInfo.redacted,
			serviceID:          startedInfo.serviceID,
			serverAddress:      desc.Server.Addr,
		}

		startedTime := time.Now()

		// Check for possible context error. If no context error, check if there's enough time to perform a
		// round trip before the Context deadline. If ctx is a Timeout Context, use the 90th percentile RTT
		// as a threshold. Otherwise, use the minimum observed RTT.
		if ctx.Err() != nil {
			err = ctx.Err()
		} else if deadline, ok := ctx.Deadline(); ok {
			if csot.IsTimeoutContext(ctx) && time.Now().Add(srvr.RTTMonitor().P90()).After(deadline) {
				err = fmt.Errorf(
					"remaining time %v until context deadline is less than 90th percentile network round-trip time: %w\n%v",
					time.Until(deadline),
					ErrDeadlineWouldBeExceeded,
					srvr.RTTMonitor().Stats())
			} else if time.Now().Add(srvr.RTTMonitor().Min()).After(deadline) {
				err = context.DeadlineExceeded
			}
		}

		if err == nil {
			// roundtrip using either the full roundTripper or a special one for when the moreToCome
			// flag is set
			roundTrip := op.roundTrip
			if moreToCome {
				roundTrip = op.moreToComeRoundTrip
			}
			res, err = roundTrip(ctx, conn, *wm)

			if ep, ok := srvr.(ErrorProcessor); ok {
				_ = ep.ProcessError(err, conn)
			}
		}

		finishedInfo.response = res
		finishedInfo.cmdErr = err
		finishedInfo.duration = time.Since(startedTime)

		op.publishFinishedEvent(ctx, finishedInfo)

		// prevIndefiniteErrorIsSet is "true" if the "err" variable has been set to the "prevIndefiniteErr" in
		// a case in the switch statement below.
		var prevIndefiniteErrIsSet bool

		// TODO(GODRIVER-2579): When refactoring the "Execute" method, consider creating a separate method for the
		// error handling logic below. This will remove the necessity of the "checkError" goto label.
	checkError:
		var perr error
		switch tt := err.(type) {
		case WriteCommandError:
			if e := err.(WriteCommandError); retrySupported && op.Type == Write && e.UnsupportedStorageEngine() {
				return ErrUnsupportedStorageEngine
			}

			connDesc := conn.Description()
			retryableErr := tt.Retryable(connDesc.WireVersion)
			preRetryWriteLabelVersion := connDesc.WireVersion != nil && connDesc.WireVersion.Max < 9
			inTransaction := op.Client != nil &&
				!(op.Client.Committing || op.Client.Aborting) && op.Client.TransactionRunning()
			// If retry is enabled and the operation isn't in a transaction, add a RetryableWriteError label for
			// retryable errors from pre-4.4 servers
			if retryableErr && preRetryWriteLabelVersion && retryEnabled && !inTransaction {
				tt.Labels = append(tt.Labels, RetryableWriteError)
			}

			// If retries are supported for the current operation on the first server description,
			// the error is considered retryable, and there are retries remaining (negative retries
			// means retry indefinitely), then retry the operation.
			if retrySupported && retryableErr && retries != 0 {
				if op.Client != nil && op.Client.Committing {
					// Apply majority write concern for retries
					op.Client.UpdateCommitTransactionWriteConcern()
					op.WriteConcern = op.Client.CurrentWc
				}
				resetForRetry(tt)
				continue
			}

			// If the error is no longer retryable and has the NoWritesPerformed label, then we should
			// set the error to the "previous indefinite error" unless the current error is already the
			// "previous indefinite error". After resetting, repeat the error check.
			if tt.HasErrorLabel(NoWritesPerformed) && !prevIndefiniteErrIsSet {
				err = prevIndefiniteErr
				prevIndefiniteErrIsSet = true

				goto checkError
			}

			// If the operation isn't being retried, process the response
			if op.ProcessResponseFn != nil {
				info := ResponseInfo{
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				}
				_ = op.ProcessResponseFn(info)
			}

			if batching && len(tt.WriteErrors) > 0 && currIndex > 0 {
				for i := range tt.WriteErrors {
					tt.WriteErrors[i].Index += int64(currIndex)
				}
			}

			// If batching is enabled and either ordered is the default (which is true) or
			// explicitly set to true and we have write errors, return the errors.
			if batching && (op.Batches.Ordered == nil || *op.Batches.Ordered) && len(tt.WriteErrors) > 0 {
				return tt
			}
			if op.Client != nil && op.Client.Committing && tt.WriteConcernError != nil {
				// When running commitTransaction we return WriteConcernErrors as an Error.
				err := Error{
					Name:    tt.WriteConcernError.Name,
					Code:    int32(tt.WriteConcernError.Code),
					Message: tt.WriteConcernError.Message,
					Labels:  tt.Labels,
					Raw:     tt.Raw,
				}
				// The UnknownTransactionCommitResult label is added to all writeConcernErrors besides unknownReplWriteConcernCode
				// and unsatisfiableWriteConcernCode
				if err.Code != unknownReplWriteConcernCode && err.Code != unsatisfiableWriteConcernCode {
					err.Labels = append(err.Labels, UnknownTransactionCommitResult)
				}
				if retryableErr && retryEnabled {
					err.Labels = append(err.Labels, RetryableWriteError)
				}
				return err
			}
			operationErr.WriteConcernError = tt.WriteConcernError
			operationErr.WriteErrors = append(operationErr.WriteErrors, tt.WriteErrors...)
			operationErr.Labels = tt.Labels
			operationErr.Raw = tt.Raw
		case Error:
			// 391 is the reauthentication required error code, so we will attempt a reauth and
			// retry the operation, if it is successful.
			if tt.Code == 391 {
				if op.Authenticator != nil {
					cfg := AuthConfig{
						Description:  conn.Description(),
						Connection:   conn,
						ClusterClock: op.Clock,
						ServerAPI:    op.ServerAPI,
					}
					if err := op.Authenticator.Reauth(ctx, &cfg); err != nil {
						return fmt.Errorf("error reauthenticating: %w", err)
					}
					if op.Client != nil && op.Client.Committing {
						// Apply majority write concern for retries
						op.Client.UpdateCommitTransactionWriteConcern()
						op.WriteConcern = op.Client.CurrentWc
					}
					resetForRetry(tt)
					continue
				}
			}
			if tt.HasErrorLabel(TransientTransactionError) || tt.HasErrorLabel(UnknownTransactionCommitResult) {
				if err := op.Client.ClearPinnedResources(); err != nil {
					return err
				}
			}

			if e := err.(Error); retrySupported && op.Type == Write && e.UnsupportedStorageEngine() {
				return ErrUnsupportedStorageEngine
			}

			connDesc := conn.Description()
			var retryableErr bool
			if op.Type == Write {
				retryableErr = tt.RetryableWrite(connDesc.WireVersion)
				preRetryWriteLabelVersion := connDesc.WireVersion != nil && connDesc.WireVersion.Max < 9
				inTransaction := op.Client != nil &&
					!(op.Client.Committing || op.Client.Aborting) && op.Client.TransactionRunning()
				// If retryWrites is enabled and the operation isn't in a transaction, add a RetryableWriteError label
				// for network errors and retryable errors from pre-4.4 servers
				if retryEnabled && !inTransaction &&
					(tt.HasErrorLabel(NetworkError) || (retryableErr && preRetryWriteLabelVersion)) {
					tt.Labels = append(tt.Labels, RetryableWriteError)
				}
			} else {
				retryableErr = tt.RetryableRead()
			}

			// If retries are supported for the current operation on the first server description,
			// the error is considered retryable, and there are retries remaining (negative retries
			// means retry indefinitely), then retry the operation.
			if retrySupported && retryableErr && retries != 0 {
				if op.Client != nil && op.Client.Committing {
					// Apply majority write concern for retries
					op.Client.UpdateCommitTransactionWriteConcern()
					op.WriteConcern = op.Client.CurrentWc
				}
				resetForRetry(tt)
				continue
			}

			// If the error is no longer retryable and has the NoWritesPerformed label, then we should
			// set the error to the "previous indefinite error" unless the current error is already the
			// "previous indefinite error". After resetting, repeat the error check.
			if tt.HasErrorLabel(NoWritesPerformed) && !prevIndefiniteErrIsSet {
				err = prevIndefiniteErr
				prevIndefiniteErrIsSet = true

				goto checkError
			}

			// If the operation isn't being retried, process the response
			if op.ProcessResponseFn != nil {
				info := ResponseInfo{
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				}
				_ = op.ProcessResponseFn(info)
			}

			if op.Client != nil && op.Client.Committing && (retryableErr || tt.Code == 50) {
				// If we got a retryable error or MaxTimeMSExpired error, we add UnknownTransactionCommitResult.
				tt.Labels = append(tt.Labels, UnknownTransactionCommitResult)
			}
			return tt
		case nil:
			if moreToCome {
				return ErrUnacknowledgedWrite
			}
			if op.ProcessResponseFn != nil {
				info := ResponseInfo{
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				}
				perr = op.ProcessResponseFn(info)
			}
			if perr != nil {
				return perr
			}
		default:
			if op.ProcessResponseFn != nil {
				info := ResponseInfo{
					ServerResponse:        res,
					Server:                srvr,
					Connection:            conn,
					ConnectionDescription: desc.Server,
					CurrentIndex:          currIndex,
				}
				_ = op.ProcessResponseFn(info)
			}
			return err
		}

		// If we're batching and there are batches remaining, advance to the next batch. This isn't
		// a retry, so increment the transaction number, reset the retries number, and don't set
		// server or connection to nil to continue using the same connection.
		if batching && len(op.Batches.Documents) > 0 {
			// If retries are supported for the current operation on the current server description,
			// the session isn't nil, and client retries are enabled, increment the txn number.
			// Calling IncrementTxnNumber() for server descriptions or topologies that do not
			// support retries (e.g. standalone topologies) will cause server errors.
			if retrySupported && op.Client != nil && op.RetryMode != nil {
				if op.RetryMode.Enabled() {
					op.Client.IncrementTxnNumber()
				}
				// Reset the retries number for RetryOncePerCommand unless context is a Timeout context, in
				// which case retries should remain as -1 (as many times as possible).
				if *op.RetryMode == RetryOncePerCommand && !csot.IsTimeoutContext(ctx) {
					retries = 1
				}
			}
			currIndex += len(op.Batches.Current)
			op.Batches.ClearBatch()
			continue
		}
		break
	}
	if len(operationErr.WriteErrors) > 0 || operationErr.WriteConcernError != nil {
		return operationErr
	}
	return nil
}

// Retryable writes are supported if the server supports sessions, the operation is not
// within a transaction, and the write is acknowledged
func (op Operation) retryable(desc description.Server) bool {
	switch op.Type {
	case Write:
		if op.Client != nil && (op.Client.Committing || op.Client.Aborting) {
			return true
		}
		if retryWritesSupported(desc) &&
			op.Client != nil && !(op.Client.TransactionInProgress() || op.Client.TransactionStarting()) &&
			writeconcern.AckWrite(op.WriteConcern) {
			return true
		}
	case Read:
		if op.Client != nil && (op.Client.Committing || op.Client.Aborting) {
			return true
		}
		if op.Client == nil || !(op.Client.TransactionInProgress() || op.Client.TransactionStarting()) {
			return true
		}
	}
	return false
}

// roundTrip writes a wiremessage to the connection and then reads a wiremessage. The wm parameter
// is reused when reading the wiremessage.
func (op Operation) roundTrip(ctx context.Context, conn Connection, wm []byte) ([]byte, error) {
	err := conn.WriteWireMessage(ctx, wm)
	if err != nil {
		return nil, op.networkError(err)
	}
	return op.readWireMessage(ctx, conn)
}

func (op Operation) readWireMessage(ctx context.Context, conn Connection) (result []byte, err error) {
	wm, err := conn.ReadWireMessage(ctx)
	if err != nil {
		return nil, op.networkError(err)
	}

	// If we're using a streamable connection, we set its streaming state based on the moreToCome flag in the server
	// response.
	if streamer, ok := conn.(StreamerConnection); ok {
		streamer.SetStreaming(wiremessage.IsMsgMoreToCome(wm))
	}

	length, _, _, opcode, rem, ok := wiremessage.ReadHeader(wm)
	if !ok || len(wm) < int(length) {
		return nil, errors.New("malformed wire message: insufficient bytes")
	}
	if opcode == wiremessage.OpCompressed {
		rawsize := length - 16 // remove header size
		// decompress wiremessage
		opcode, rem, err = op.decompressWireMessage(rem[:rawsize])
		if err != nil {
			return nil, err
		}
	}

	// decode
	res, err := op.decodeResult(ctx, opcode, rem)
	// Update cluster/operation time and recovery tokens before handling the error to ensure we're properly updating
	// everything.
	op.updateClusterTimes(res)
	op.updateOperationTime(res)
	op.Client.UpdateRecoveryToken(bson.Raw(res))

	// Update snapshot time if operation was a "find", "aggregate" or "distinct".
	if op.Name == driverutil.FindOp || op.Name == driverutil.AggregateOp || op.Name == driverutil.DistinctOp {
		op.Client.UpdateSnapshotTime(res)
	}

	if err != nil {
		return res, err
	}

	// If there is no error, automatically attempt to decrypt all results if client side encryption is enabled.
	if op.Crypt != nil {
		res, err = op.Crypt.Decrypt(ctx, res)
	}
	return res, err
}

// networkError wraps the provided error in an Error with label "NetworkError" and, if a transaction
// is running or committing, the appropriate transaction state labels. The returned error indicates
// the operation should be retried for reads and writes. If err is nil, networkError returns nil.
func (op Operation) networkError(err error) error {
	if err == nil {
		return nil
	}

	labels := []string{NetworkError}
	if op.Client != nil {
		op.Client.MarkDirty()
	}
	if op.Client != nil && op.Client.TransactionRunning() && !op.Client.Committing {
		labels = append(labels, TransientTransactionError)
	}
	if op.Client != nil && op.Client.Committing {
		labels = append(labels, UnknownTransactionCommitResult)
	}
	return Error{Message: err.Error(), Labels: labels, Wrapped: err}
}

// moreToComeRoundTrip writes a wiremessage to the provided connection. This is used when an OP_MSG is
// being sent with  the moreToCome bit set.
func (op *Operation) moreToComeRoundTrip(ctx context.Context, conn Connection, wm []byte) (result []byte, err error) {
	err = conn.WriteWireMessage(ctx, wm)
	if err != nil {
		if op.Client != nil {
			op.Client.MarkDirty()
		}
		err = Error{Message: err.Error(), Labels: []string{TransientTransactionError, NetworkError}, Wrapped: err}
	}
	return bsoncore.BuildDocument(nil, bsoncore.AppendInt32Element(nil, "ok", 1)), err
}

// decompressWireMessage handles decompressing a wiremessage without the header.
func (Operation) decompressWireMessage(wm []byte) (wiremessage.OpCode, []byte, error) {
	// get the original opcode and uncompressed size
	opcode, rem, ok := wiremessage.ReadCompressedOriginalOpCode(wm)
	if !ok {
		return 0, nil, errors.New("malformed OP_COMPRESSED: missing original opcode")
	}
	uncompressedSize, rem, ok := wiremessage.ReadCompressedUncompressedSize(rem)
	if !ok {
		return 0, nil, errors.New("malformed OP_COMPRESSED: missing uncompressed size")
	}
	// get the compressor ID and decompress the message
	compressorID, rem, ok := wiremessage.ReadCompressedCompressorID(rem)
	if !ok {
		return 0, nil, errors.New("malformed OP_COMPRESSED: missing compressor ID")
	}

	opts := CompressionOpts{
		Compressor:       compressorID,
		UncompressedSize: uncompressedSize,
	}
	uncompressed, err := DecompressPayload(rem, opts)
	if err != nil {
		return 0, nil, err
	}

	return opcode, uncompressed, nil
}

func (op Operation) addBatchArray(dst []byte) []byte {
	aidx, dst := bsoncore.AppendArrayElementStart(dst, op.Batches.Identifier)
	for i, doc := range op.Batches.Current {
		dst = bsoncore.AppendDocumentElement(dst, strconv.Itoa(i), doc)
	}
	dst, _ = bsoncore.AppendArrayEnd(dst, aidx)
	return dst
}

func (op Operation) createLegacyHandshakeWireMessage(
	maxTimeMS uint64,
	dst []byte,
	desc description.SelectedServer,
) ([]byte, startedInformation, error) {
	var info startedInformation
	flags := op.secondaryOK(desc)
	var wmindex int32
	info.requestID = wiremessage.NextRequestID()
	wmindex, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpQuery)
	dst = wiremessage.AppendQueryFlags(dst, flags)

	dollarCmd := [...]byte{'.', '$', 'c', 'm', 'd'}

	// FullCollectionName
	dst = append(dst, op.Database...)
	dst = append(dst, dollarCmd[:]...)
	dst = append(dst, 0x00)
	dst = wiremessage.AppendQueryNumberToSkip(dst, 0)
	dst = wiremessage.AppendQueryNumberToReturn(dst, -1)

	wrapper := int32(-1)
	rp, err := op.createReadPref(desc, true)
	if err != nil {
		return dst, info, err
	}
	if len(rp) > 0 {
		wrapper, dst = bsoncore.AppendDocumentStart(dst)
		dst = bsoncore.AppendHeader(dst, bsontype.EmbeddedDocument, "$query")
	}
	idx, dst := bsoncore.AppendDocumentStart(dst)
	dst, err = op.CommandFn(dst, desc)
	if err != nil {
		return dst, info, err
	}

	if op.Batches != nil && len(op.Batches.Current) > 0 {
		dst = op.addBatchArray(dst)
	}

	dst, err = op.addReadConcern(dst, desc)
	if err != nil {
		return dst, info, err
	}

	dst, err = op.addWriteConcern(dst, desc)
	if err != nil {
		return dst, info, err
	}

	dst, err = op.addSession(dst, desc)
	if err != nil {
		return dst, info, err
	}

	dst = op.addClusterTime(dst, desc)
	dst = op.addServerAPI(dst)
	// If maxTimeMS is greater than 0 append it to wire message. A maxTimeMS value of 0 only explicitly
	// specifies the default behavior of no timeout server-side.
	if maxTimeMS > 0 {
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", int64(maxTimeMS))
	}

	dst, _ = bsoncore.AppendDocumentEnd(dst, idx)
	// Command monitoring only reports the document inside $query
	info.cmd = dst[idx:]

	if len(rp) > 0 {
		var err error
		dst = bsoncore.AppendDocumentElement(dst, "$readPreference", rp)
		dst, err = bsoncore.AppendDocumentEnd(dst, wrapper)
		if err != nil {
			return dst, info, err
		}
	}

	return bsoncore.UpdateLength(dst, wmindex, int32(len(dst[wmindex:]))), info, nil
}

func (op Operation) createMsgWireMessage(
	ctx context.Context,
	maxTimeMS uint64,
	dst []byte,
	desc description.SelectedServer,
	conn Connection,
	requestID int32,
) ([]byte, startedInformation, error) {
	var info startedInformation
	var flags wiremessage.MsgFlag
	var wmindex int32
	// We set the MoreToCome bit if we have a write concern, it's unacknowledged, and we either
	// aren't batching or we are encoding the last batch.
	if op.WriteConcern != nil && !writeconcern.AckWrite(op.WriteConcern) && (op.Batches == nil || len(op.Batches.Documents) == 0) {
		flags = wiremessage.MoreToCome
	}
	// Set the ExhaustAllowed flag if the connection supports streaming. This will tell the server that it can
	// respond with the MoreToCome flag and then stream responses over this connection.
	if streamer, ok := conn.(StreamerConnection); ok && streamer.SupportsStreaming() {
		flags |= wiremessage.ExhaustAllowed
	}

	info.requestID = requestID
	wmindex, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpMsg)
	dst = wiremessage.AppendMsgFlags(dst, flags)
	// Body
	dst = wiremessage.AppendMsgSectionType(dst, wiremessage.SingleDocument)

	idx, dst := bsoncore.AppendDocumentStart(dst)

	dst, err := op.addCommandFields(ctx, dst, desc)
	if err != nil {
		return dst, info, err
	}
	dst, err = op.addReadConcern(dst, desc)
	if err != nil {
		return dst, info, err
	}
	dst, err = op.addWriteConcern(dst, desc)
	if err != nil {
		return dst, info, err
	}
	dst, err = op.addSession(dst, desc)
	if err != nil {
		return dst, info, err
	}

	dst = op.addClusterTime(dst, desc)
	dst = op.addServerAPI(dst)
	// If maxTimeMS is greater than 0 append it to wire message. A maxTimeMS value of 0 only explicitly
	// specifies the default behavior of no timeout server-side.
	if maxTimeMS > 0 {
		dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", int64(maxTimeMS))
	}

	dst = bsoncore.AppendStringElement(dst, "$db", op.Database)
	rp, err := op.createReadPref(desc, false)
	if err != nil {
		return dst, info, err
	}
	if len(rp) > 0 {
		dst = bsoncore.AppendDocumentElement(dst, "$readPreference", rp)
	}

	dst, _ = bsoncore.AppendDocumentEnd(dst, idx)
	// The command document for monitoring shouldn't include the type 1 payload as a document sequence
	info.cmd = dst[idx:]

	// add batch as a document sequence if auto encryption is not enabled
	// if auto encryption is enabled, the batch will already be an array in the command document
	if !op.shouldEncrypt() && op.Batches != nil && len(op.Batches.Current) > 0 {
		info.documentSequenceIncluded = true
		dst = wiremessage.AppendMsgSectionType(dst, wiremessage.DocumentSequence)
		idx, dst = bsoncore.ReserveLength(dst)

		dst = append(dst, op.Batches.Identifier...)
		dst = append(dst, 0x00)

		for _, doc := range op.Batches.Current {
			dst = append(dst, doc...)
		}

		dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	}

	return bsoncore.UpdateLength(dst, wmindex, int32(len(dst[wmindex:]))), info, nil
}

// isLegacyHandshake returns True if the operation is the first message of
// the initial handshake and should use a legacy hello.
func isLegacyHandshake(op Operation, desc description.SelectedServer) bool {
	isInitialHandshake := desc.WireVersion == nil || desc.WireVersion.Max == 0

	return op.Legacy == LegacyHandshake && isInitialHandshake
}

func (op Operation) createWireMessage(
	ctx context.Context,
	maxTimeMS uint64,
	dst []byte,
	desc description.SelectedServer,
	conn Connection,
	requestID int32,
) ([]byte, startedInformation, error) {
	if isLegacyHandshake(op, desc) {
		return op.createLegacyHandshakeWireMessage(maxTimeMS, dst, desc)
	}

	return op.createMsgWireMessage(ctx, maxTimeMS, dst, desc, conn, requestID)
}

// addCommandFields adds the fields for a command to the wire message in dst. This assumes that the start of the document
// has already been added and does not add the final 0 byte.
func (op Operation) addCommandFields(ctx context.Context, dst []byte, desc description.SelectedServer) ([]byte, error) {
	if !op.shouldEncrypt() {
		return op.CommandFn(dst, desc)
	}

	if desc.WireVersion.Max < cryptMinWireVersion {
		return dst, errors.New("auto-encryption requires a MongoDB version of 4.2")
	}

	// create temporary command document
	cidx, cmdDst := bsoncore.AppendDocumentStart(nil)
	var err error
	cmdDst, err = op.CommandFn(cmdDst, desc)
	if err != nil {
		return dst, err
	}
	// use a BSON array instead of a type 1 payload because mongocryptd will convert to arrays regardless
	if op.Batches != nil && len(op.Batches.Current) > 0 {
		cmdDst = op.addBatchArray(cmdDst)
	}
	cmdDst, _ = bsoncore.AppendDocumentEnd(cmdDst, cidx)

	// encrypt the command
	encrypted, err := op.Crypt.Encrypt(ctx, op.Database, cmdDst)
	if err != nil {
		return dst, err
	}
	// append encrypted command to original destination, removing the first 4 bytes (length) and final byte (terminator)
	dst = append(dst, encrypted[4:len(encrypted)-1]...)
	return dst, nil
}

// addServerAPI adds the relevant fields for server API specification to the wire message in dst.
func (op Operation) addServerAPI(dst []byte) []byte {
	sa := op.ServerAPI
	if sa == nil {
		return dst
	}

	dst = bsoncore.AppendStringElement(dst, "apiVersion", sa.ServerAPIVersion)
	if sa.Strict != nil {
		dst = bsoncore.AppendBooleanElement(dst, "apiStrict", *sa.Strict)
	}
	if sa.DeprecationErrors != nil {
		dst = bsoncore.AppendBooleanElement(dst, "apiDeprecationErrors", *sa.DeprecationErrors)
	}
	return dst
}

func (op Operation) addReadConcern(dst []byte, desc description.SelectedServer) ([]byte, error) {
	if op.MinimumReadConcernWireVersion > 0 && (desc.WireVersion == nil || !desc.WireVersion.Includes(op.MinimumReadConcernWireVersion)) {
		return dst, nil
	}
	rc := op.ReadConcern
	client := op.Client
	// Starting transaction's read concern overrides all others
	if client != nil && client.TransactionStarting() && client.CurrentRc != nil {
		rc = client.CurrentRc
	}

	// start transaction must append afterclustertime IF causally consistent and operation time exists
	if rc == nil && client != nil && client.TransactionStarting() && client.Consistent && client.OperationTime != nil {
		rc = readconcern.New()
	}

	if client != nil && client.Snapshot {
		if desc.WireVersion.Max < readSnapshotMinWireVersion {
			return dst, errors.New("snapshot reads require MongoDB 5.0 or later")
		}
		rc = readconcern.Snapshot()
	}

	if rc == nil {
		return dst, nil
	}

	_, data, err := rc.MarshalBSONValue() // always returns a document
	if err != nil {
		return dst, err
	}

	if sessionsSupported(desc.WireVersion) && client != nil {
		if client.Consistent && client.OperationTime != nil {
			data = data[:len(data)-1] // remove the null byte
			data = bsoncore.AppendTimestampElement(data, "afterClusterTime", client.OperationTime.T, client.OperationTime.I)
			data, _ = bsoncore.AppendDocumentEnd(data, 0)
		}
		if client.Snapshot && client.SnapshotTime != nil {
			data = data[:len(data)-1] // remove the null byte
			data = bsoncore.AppendTimestampElement(data, "atClusterTime", client.SnapshotTime.T, client.SnapshotTime.I)
			data, _ = bsoncore.AppendDocumentEnd(data, 0)
		}
	}

	if len(data) == bsoncore.EmptyDocumentLength {
		return dst, nil
	}
	return bsoncore.AppendDocumentElement(dst, "readConcern", data), nil
}

func (op Operation) addWriteConcern(dst []byte, desc description.SelectedServer) ([]byte, error) {
	if op.MinimumWriteConcernWireVersion > 0 && (desc.WireVersion == nil || !desc.WireVersion.Includes(op.MinimumWriteConcernWireVersion)) {
		return dst, nil
	}
	wc := op.WriteConcern
	if wc == nil {
		return dst, nil
	}

	t, data, err := wc.MarshalBSONValue()
	if errors.Is(err, writeconcern.ErrEmptyWriteConcern) {
		return dst, nil
	}
	if err != nil {
		return dst, err
	}

	return append(bsoncore.AppendHeader(dst, t, "writeConcern"), data...), nil
}

func (op Operation) addSession(dst []byte, desc description.SelectedServer) ([]byte, error) {
	client := op.Client

	// If the operation is defined for an explicit session but the server
	// does not support sessions, then throw an error.
	if client != nil && !client.IsImplicit && desc.SessionTimeoutMinutesPtr == nil {
		return nil, fmt.Errorf("current topology does not support sessions")
	}

	if client == nil || !sessionsSupported(desc.WireVersion) || desc.SessionTimeoutMinutesPtr == nil {
		return dst, nil
	}
	if err := client.UpdateUseTime(); err != nil {
		return dst, err
	}
	dst = bsoncore.AppendDocumentElement(dst, "lsid", client.SessionID)

	var addedTxnNumber bool
	if op.Type == Write && client.RetryWrite {
		addedTxnNumber = true
		dst = bsoncore.AppendInt64Element(dst, "txnNumber", op.Client.TxnNumber)
	}
	if client.TransactionRunning() || client.RetryingCommit {
		if !addedTxnNumber {
			dst = bsoncore.AppendInt64Element(dst, "txnNumber", op.Client.TxnNumber)
		}
		if client.TransactionStarting() {
			dst = bsoncore.AppendBooleanElement(dst, "startTransaction", true)
		}
		dst = bsoncore.AppendBooleanElement(dst, "autocommit", false)
	}

	return dst, client.ApplyCommand(desc.Server)
}

func (op Operation) addClusterTime(dst []byte, desc description.SelectedServer) []byte {
	client, clock := op.Client, op.Clock
	if (clock == nil && client == nil) || !sessionsSupported(desc.WireVersion) {
		return dst
	}
	clusterTime := clock.GetClusterTime()
	if client != nil {
		clusterTime = session.MaxClusterTime(clusterTime, client.ClusterTime)
	}
	if clusterTime == nil {
		return dst
	}
	val, err := clusterTime.LookupErr("$clusterTime")
	if err != nil {
		return dst
	}
	return append(bsoncore.AppendHeader(dst, val.Type, "$clusterTime"), val.Value...)
	// return bsoncore.AppendDocumentElement(dst, "$clusterTime", clusterTime)
}

// calculateMaxTimeMS calculates the value of the 'maxTimeMS' field to potentially append
// to the wire message based on the current context's deadline and the 90th percentile RTT
// if the ctx is a Timeout context. If the context is not a Timeout context, it uses the
// operation's MaxTimeMS if set. If no MaxTimeMS is set on the operation, and context is
// not a Timeout context, calculateMaxTimeMS returns 0.
func (op Operation) calculateMaxTimeMS(ctx context.Context, mon RTTMonitor) (uint64, error) {
	// If CSOT is enabled and we're not omitting the CSOT-calculated maxTimeMS
	// value, then calculate maxTimeMS.
	//
	// This allows commands that do not currently send CSOT-calculated maxTimeMS
	// (e.g. Find and Aggregate) to still use a manually-provided maxTimeMS
	// value.
	//
	// TODO(GODRIVER-2944): Remove or refactor this logic when we add the
	// "timeoutMode" option, which will allow users to opt-in to the
	// CSOT-calculated maxTimeMS values if that's the behavior they want.
	if csot.IsTimeoutContext(ctx) && !op.OmitCSOTMaxTimeMS {
		if deadline, ok := ctx.Deadline(); ok {
			remainingTimeout := time.Until(deadline)
			rtt90 := mon.P90()
			maxTime := remainingTimeout - rtt90

			// Always round up to the next millisecond value so we never truncate the calculated
			// maxTimeMS value (e.g. 400 microseconds evaluates to 1ms, not 0ms).
			maxTimeMS := int64((maxTime + (time.Millisecond - 1)) / time.Millisecond)
			if maxTimeMS <= 0 {
				return 0, fmt.Errorf(
					"negative maxTimeMS: remaining time %v until context deadline is less than 90th percentile network round-trip time (%v): %w",
					remainingTimeout,
					mon.Stats(),
					ErrDeadlineWouldBeExceeded)
			}

			// The server will return a "BadValue" error if maxTimeMS is greater
			// than the maximum positive int32 value (about 24.9 days). If the
			// user specified a timeout value greater than that,  omit maxTimeMS
			// and let the client-side timeout handle cancelling the op if the
			// timeout is ever reached.
			if maxTimeMS > math.MaxInt32 {
				return 0, nil
			}

			return uint64(maxTimeMS), nil
		}
	} else if op.MaxTime != nil {
		// Users are not allowed to pass a negative value as MaxTime. A value of 0 would indicate
		// no timeout and is allowed.
		if *op.MaxTime < 0 {
			return 0, ErrNegativeMaxTime
		}
		// Always round up to the next millisecond value so we never truncate the requested
		// MaxTime value (e.g. 400 microseconds evaluates to 1ms, not 0ms).
		return uint64((*op.MaxTime + (time.Millisecond - 1)) / time.Millisecond), nil
	}
	return 0, nil
}

// updateClusterTimes updates the cluster times for the session and cluster clock attached to this
// operation. While the session's AdvanceClusterTime may return an error, this method does not
// because an error being returned from this method will not be returned further up.
func (op Operation) updateClusterTimes(response bsoncore.Document) {
	// Extract cluster time.
	value, err := response.LookupErr("$clusterTime")
	if err != nil {
		// $clusterTime not included by the server
		return
	}
	clusterTime := bsoncore.BuildDocumentFromElements(nil, bsoncore.AppendValueElement(nil, "$clusterTime", value))

	sess, clock := op.Client, op.Clock

	if sess != nil {
		_ = sess.AdvanceClusterTime(bson.Raw(clusterTime))
	}

	if clock != nil {
		clock.AdvanceClusterTime(bson.Raw(clusterTime))
	}
}

// updateOperationTime updates the operation time on the session attached to this operation. While
// the session's AdvanceOperationTime method may return an error, this method does not because an
// error being returned from this method will not be returned further up.
func (op Operation) updateOperationTime(response bsoncore.Document) {
	sess := op.Client
	if sess == nil {
		return
	}

	opTimeElem, err := response.LookupErr("operationTime")
	if err != nil {
		// operationTime not included by the server
		return
	}

	t, i := opTimeElem.Timestamp()
	_ = sess.AdvanceOperationTime(&primitive.Timestamp{
		T: t,
		I: i,
	})
}

func (op Operation) getReadPrefBasedOnTransaction() (*readpref.ReadPref, error) {
	if op.Client != nil && op.Client.TransactionRunning() {
		// Transaction's read preference always takes priority
		rp := op.Client.CurrentRp
		// Reads in a transaction must have read preference primary
		// This must not be checked in startTransaction
		if rp != nil && !op.Client.TransactionStarting() && rp.Mode() != readpref.PrimaryMode {
			return nil, ErrNonPrimaryReadPref
		}
		return rp, nil
	}
	return op.ReadPreference, nil
}

// createReadPref will attempt to create a document with the "readPreference"
// object and various related fields such as "mode", "tags", and
// "maxStalenessSeconds".
func (op Operation) createReadPref(desc description.SelectedServer, isOpQuery bool) (bsoncore.Document, error) {
	if op.omitReadPreference {
		return nil, nil
	}

	// TODO(GODRIVER-2231): Instead of checking if isOutputAggregate and desc.Server.WireVersion.Max < 13, somehow check
	// TODO if supplied readPreference was "overwritten" with primary in description.selectForReplicaSet.
	if desc.Server.Kind == description.Standalone || (isOpQuery && desc.Server.Kind != description.Mongos) ||
		op.Type == Write || (op.IsOutputAggregate && desc.Server.WireVersion.Max < 13) {
		// Don't send read preference for:
		// 1. all standalones
		// 2. non-mongos when using OP_QUERY
		// 3. all writes
		// 4. when operation is an aggregate with an output stage, and selected server's wire
		//    version is < 13
		return nil, nil
	}

	idx, doc := bsoncore.AppendDocumentStart(nil)
	rp, err := op.getReadPrefBasedOnTransaction()
	if err != nil {
		return nil, err
	}

	if rp == nil {
		if desc.Kind == description.Single && desc.Server.Kind != description.Mongos {
			doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
			doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
			return doc, nil
		}
		return nil, nil
	}

	switch rp.Mode() {
	case readpref.PrimaryMode:
		if desc.Server.Kind == description.Mongos {
			return nil, nil
		}
		if desc.Kind == description.Single {
			doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
			doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
			return doc, nil
		}

		// OP_MSG requires never sending read preference "primary"
		// except for topology "single".
		//
		// It is important to note that although the Go Driver does not
		// support legacy opcodes, OP_QUERY has different rules for
		// adding read preference to commands.
		return nil, nil
	case readpref.PrimaryPreferredMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
	case readpref.SecondaryPreferredMode:
		_, ok := rp.MaxStaleness()
		if desc.Server.Kind == description.Mongos && isOpQuery && !ok && len(rp.TagSets()) == 0 && rp.HedgeEnabled() == nil {
			return nil, nil
		}
		doc = bsoncore.AppendStringElement(doc, "mode", "secondaryPreferred")
	case readpref.SecondaryMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "secondary")
	case readpref.NearestMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "nearest")
	}

	sets := make([]bsoncore.Document, 0, len(rp.TagSets()))
	for _, ts := range rp.TagSets() {
		i, set := bsoncore.AppendDocumentStart(nil)
		for _, t := range ts {
			set = bsoncore.AppendStringElement(set, t.Name, t.Value)
		}
		set, _ = bsoncore.AppendDocumentEnd(set, i)
		sets = append(sets, set)
	}
	if len(sets) > 0 {
		var aidx int32
		aidx, doc = bsoncore.AppendArrayElementStart(doc, "tags")
		for i, set := range sets {
			doc = bsoncore.AppendDocumentElement(doc, strconv.Itoa(i), set)
		}
		doc, _ = bsoncore.AppendArrayEnd(doc, aidx)
	}

	if d, ok := rp.MaxStaleness(); ok {
		doc = bsoncore.AppendInt32Element(doc, "maxStalenessSeconds", int32(d.Seconds()))
	}

	if hedgeEnabled := rp.HedgeEnabled(); hedgeEnabled != nil {
		var hedgeIdx int32
		hedgeIdx, doc = bsoncore.AppendDocumentElementStart(doc, "hedge")
		doc = bsoncore.AppendBooleanElement(doc, "enabled", *hedgeEnabled)
		doc, err = bsoncore.AppendDocumentEnd(doc, hedgeIdx)
		if err != nil {
			return nil, fmt.Errorf("error creating hedge document: %w", err)
		}
	}

	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
	return doc, nil
}

func (op Operation) secondaryOK(desc description.SelectedServer) wiremessage.QueryFlag {
	if desc.Kind == description.Single && desc.Server.Kind != description.Mongos {
		return wiremessage.SecondaryOK
	}

	if rp := op.ReadPreference; rp != nil && rp.Mode() != readpref.PrimaryMode {
		return wiremessage.SecondaryOK
	}

	return 0
}

func (Operation) canCompress(cmd string) bool {
	if cmd == handshake.LegacyHello || cmd == "hello" || cmd == "saslStart" || cmd == "saslContinue" || cmd == "getnonce" || cmd == "authenticate" ||
		cmd == "createUser" || cmd == "updateUser" || cmd == "copydbSaslStart" || cmd == "copydbgetnonce" || cmd == "copydb" {
		return false
	}
	return true
}

// decodeOpReply extracts the necessary information from an OP_REPLY wire message.
// Returns the decoded OP_REPLY. If the err field of the returned opReply is non-nil, an error occurred while decoding
// or validating the response and the other fields are undefined.
func (Operation) decodeOpReply(wm []byte) opReply {
	var reply opReply
	var ok bool

	reply.responseFlags, wm, ok = wiremessage.ReadReplyFlags(wm)
	if !ok {
		reply.err = errors.New("malformed OP_REPLY: missing flags")
		return reply
	}
	reply.cursorID, wm, ok = wiremessage.ReadReplyCursorID(wm)
	if !ok {
		reply.err = errors.New("malformed OP_REPLY: missing cursorID")
		return reply
	}
	reply.startingFrom, wm, ok = wiremessage.ReadReplyStartingFrom(wm)
	if !ok {
		reply.err = errors.New("malformed OP_REPLY: missing startingFrom")
		return reply
	}
	reply.numReturned, wm, ok = wiremessage.ReadReplyNumberReturned(wm)
	if !ok {
		reply.err = errors.New("malformed OP_REPLY: missing numberReturned")
		return reply
	}
	reply.documents, _, ok = wiremessage.ReadReplyDocuments(wm)
	if !ok {
		reply.err = errors.New("malformed OP_REPLY: could not read documents from reply")
	}

	if reply.responseFlags&wiremessage.QueryFailure == wiremessage.QueryFailure {
		reply.err = QueryFailureError{
			Message:  "command failure",
			Response: reply.documents[0],
		}
		return reply
	}
	if reply.responseFlags&wiremessage.CursorNotFound == wiremessage.CursorNotFound {
		reply.err = ErrCursorNotFound
		return reply
	}
	if reply.numReturned != int32(len(reply.documents)) {
		reply.err = ErrReplyDocumentMismatch
		return reply
	}

	return reply
}

func (op Operation) decodeResult(ctx context.Context, opcode wiremessage.OpCode, wm []byte) (bsoncore.Document, error) {
	switch opcode {
	case wiremessage.OpReply:
		reply := op.decodeOpReply(wm)
		if reply.err != nil {
			return nil, reply.err
		}
		if reply.numReturned == 0 {
			return nil, ErrNoDocCommandResponse
		}
		if reply.numReturned > 1 {
			return nil, ErrMultiDocCommandResponse
		}
		rdr := reply.documents[0]
		if err := rdr.Validate(); err != nil {
			return nil, NewCommandResponseError("malformed OP_REPLY: invalid document", err)
		}

		return rdr, ExtractErrorFromServerResponse(ctx, rdr)
	case wiremessage.OpMsg:
		_, wm, ok := wiremessage.ReadMsgFlags(wm)
		if !ok {
			return nil, errors.New("malformed wire message: missing OP_MSG flags")
		}

		var res bsoncore.Document
		for len(wm) > 0 {
			var stype wiremessage.SectionType
			stype, wm, ok = wiremessage.ReadMsgSectionType(wm)
			if !ok {
				return nil, errors.New("malformed wire message: insuffienct bytes to read section type")
			}

			switch stype {
			case wiremessage.SingleDocument:
				res, wm, ok = wiremessage.ReadMsgSectionSingleDocument(wm)
				if !ok {
					return nil, errors.New("malformed wire message: insufficient bytes to read single document")
				}
			case wiremessage.DocumentSequence:
				_, _, wm, ok = wiremessage.ReadMsgSectionDocumentSequence(wm)
				if !ok {
					return nil, errors.New("malformed wire message: insufficient bytes to read document sequence")
				}
			default:
				return nil, fmt.Errorf("malformed wire message: unknown section type %v", stype)
			}
		}

		err := res.Validate()
		if err != nil {
			return nil, NewCommandResponseError("malformed OP_MSG: invalid document", err)
		}

		return res, ExtractErrorFromServerResponse(ctx, res)
	default:
		return nil, fmt.Errorf("cannot decode result from %s", opcode)
	}
}

// getCommandName returns the name of the command from the given BSON document.
func (op Operation) getCommandName(doc []byte) string {
	// skip 4 bytes for document length and 1 byte for element type
	idx := bytes.IndexByte(doc[5:], 0x00) // look for the 0 byte after the command name
	return string(doc[5 : idx+5])
}

func (op *Operation) redactCommand(cmd string, doc bsoncore.Document) bool {
	if cmd == "authenticate" || cmd == "saslStart" || cmd == "saslContinue" || cmd == "getnonce" || cmd == "createUser" ||
		cmd == "updateUser" || cmd == "copydbgetnonce" || cmd == "copydbsaslstart" || cmd == "copydb" {

		return true
	}
	if strings.ToLower(cmd) != handshake.LegacyHelloLowercase && cmd != "hello" {
		return false
	}

	// A hello without speculative authentication can be monitored.
	_, err := doc.LookupErr("speculativeAuthenticate")
	return err == nil
}

// canLogCommandMessage returns true if the command can be logged.
func (op Operation) canLogCommandMessage() bool {
	return op.Logger != nil && op.Logger.LevelComponentEnabled(logger.LevelDebug, logger.ComponentCommand)
}

func (op Operation) canPublishStartedEvent() bool {
	return op.CommandMonitor != nil && op.CommandMonitor.Started != nil
}

// publishStartedEvent publishes a CommandStartedEvent to the operation's command monitor if possible. If the command is
// an unacknowledged write, a CommandSucceededEvent will be published as well. If started events are not being monitored,
// no events are published.
func (op Operation) publishStartedEvent(ctx context.Context, info startedInformation) {
	// If logging is enabled for the command component at the debug level, log the command response.
	if op.canLogCommandMessage() {
		host, port, _ := net.SplitHostPort(info.serverAddress.String())

		redactedCmd := redactStartedInformationCmd(op, info).String()
		formattedCmd := logger.FormatMessage(redactedCmd, op.Logger.MaxDocumentLength)

		op.Logger.Print(logger.LevelDebug,
			logger.ComponentCommand,
			logger.CommandStarted,
			logger.SerializeCommand(logger.Command{
				DriverConnectionID: info.driverConnectionID,
				Message:            logger.CommandStarted,
				Name:               info.cmdName,
				DatabaseName:       op.Database,
				RequestID:          int64(info.requestID),
				ServerConnectionID: info.serverConnID,
				ServerHost:         host,
				ServerPort:         port,
				ServiceID:          info.serviceID,
			},
				logger.KeyCommand, formattedCmd)...)

	}

	if op.canPublishStartedEvent() {
		started := &event.CommandStartedEvent{
			Command:              redactStartedInformationCmd(op, info),
			DatabaseName:         op.Database,
			CommandName:          info.cmdName,
			RequestID:            int64(info.requestID),
			ConnectionID:         info.connID,
			ServerConnectionID:   convertInt64PtrToInt32Ptr(info.serverConnID),
			ServerConnectionID64: info.serverConnID,
			ServiceID:            info.serviceID,
		}
		op.CommandMonitor.Started(ctx, started)
	}
}

// canPublishFinishedEvent returns true if a CommandSucceededEvent can be
// published for the given command. This is true if the command is not an
// unacknowledged write and the command monitor is monitoring succeeded events.
func (op Operation) canPublishFinishedEvent(info finishedInformation) bool {
	success := info.success()

	return op.CommandMonitor != nil &&
		(!success || op.CommandMonitor.Succeeded != nil) &&
		(success || op.CommandMonitor.Failed != nil)
}

// publishFinishedEvent publishes either a CommandSucceededEvent or a CommandFailedEvent to the operation's command
// monitor if possible. If success/failure events aren't being monitored, no events are published.
func (op Operation) publishFinishedEvent(ctx context.Context, info finishedInformation) {
	if op.canLogCommandMessage() && info.success() {
		host, port, _ := net.SplitHostPort(info.serverAddress.String())

		redactedReply := redactFinishedInformationResponse(info).String()
		formattedReply := logger.FormatMessage(redactedReply, op.Logger.MaxDocumentLength)

		op.Logger.Print(logger.LevelDebug,
			logger.ComponentCommand,
			logger.CommandSucceeded,
			logger.SerializeCommand(logger.Command{
				DriverConnectionID: info.driverConnectionID,
				Message:            logger.CommandSucceeded,
				Name:               info.cmdName,
				DatabaseName:       op.Database,
				RequestID:          int64(info.requestID),
				ServerConnectionID: info.serverConnID,
				ServerHost:         host,
				ServerPort:         port,
				ServiceID:          info.serviceID,
			},
				logger.KeyDurationMS, info.duration.Milliseconds(),
				logger.KeyReply, formattedReply)...)
	}

	if op.canLogCommandMessage() && !info.success() {
		host, port, _ := net.SplitHostPort(info.serverAddress.String())

		formattedReply := logger.FormatMessage(info.cmdErr.Error(), op.Logger.MaxDocumentLength)

		op.Logger.Print(logger.LevelDebug,
			logger.ComponentCommand,
			logger.CommandFailed,
			logger.SerializeCommand(logger.Command{
				DriverConnectionID: info.driverConnectionID,
				Message:            logger.CommandFailed,
				Name:               info.cmdName,
				DatabaseName:       op.Database,
				RequestID:          int64(info.requestID),
				ServerConnectionID: info.serverConnID,
				ServerHost:         host,
				ServerPort:         port,
				ServiceID:          info.serviceID,
			},
				logger.KeyDurationMS, info.duration.Milliseconds(),
				logger.KeyFailure, formattedReply)...)
	}

	// If the finished event cannot be published, return early.
	if !op.canPublishFinishedEvent(info) {
		return
	}

	finished := event.CommandFinishedEvent{
		CommandName:          info.cmdName,
		DatabaseName:         op.Database,
		RequestID:            int64(info.requestID),
		ConnectionID:         info.connID,
		Duration:             info.duration,
		DurationNanos:        info.duration.Nanoseconds(),
		ServerConnectionID:   convertInt64PtrToInt32Ptr(info.serverConnID),
		ServerConnectionID64: info.serverConnID,
		ServiceID:            info.serviceID,
	}

	if info.success() {
		successEvent := &event.CommandSucceededEvent{
			Reply:                redactFinishedInformationResponse(info),
			CommandFinishedEvent: finished,
		}
		op.CommandMonitor.Succeeded(ctx, successEvent)

		return
	}

	failedEvent := &event.CommandFailedEvent{
		Failure:              info.cmdErr.Error(),
		CommandFinishedEvent: finished,
	}
	op.CommandMonitor.Failed(ctx, failedEvent)
}

// sessionsSupported returns true of the given server version indicates that it supports sessions.
func sessionsSupported(wireVersion *description.VersionRange) bool {
	return wireVersion != nil
}

// retryWritesSupported returns true if this description represents a server that supports retryable writes.
func retryWritesSupported(s description.Server) bool {
	return s.SessionTimeoutMinutesPtr != nil && s.Kind != description.Standalone
}
