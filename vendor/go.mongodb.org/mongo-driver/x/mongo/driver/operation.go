package driver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

const defaultLocalThreshold = 15 * time.Millisecond

var dollarCmd = [...]byte{'.', '$', 'c', 'm', 'd'}

var (
	// ErrNoDocCommandResponse occurs when the server indicated a response existed, but none was found.
	ErrNoDocCommandResponse = errors.New("command returned no documents")
	// ErrMultiDocCommandResponse occurs when the server sent multiple documents in response to a command.
	ErrMultiDocCommandResponse = errors.New("command returned multiple documents")
	// ErrReplyDocumentMismatch occurs when the number of documents returned in an OP_QUERY does not match the numberReturned field.
	ErrReplyDocumentMismatch = errors.New("number of documents returned does not match numberReturned field")
	// ErrNonPrimaryReadPref is returned when a read is attempted in a transaction with a non-primary read preference.
	ErrNonPrimaryReadPref = errors.New("read preference in a transaction must be primary")
)

const (
	// maximum BSON object size when client side encryption is enabled
	cryptMaxBsonObjectSize uint32 = 2097152
	// minimum wire version necessary to use automatic encryption
	cryptMinWireVersion int32 = 8
)

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
}

// finishedInformation keeps track of all of the information necessary for monitoring success and failure events.
type finishedInformation struct {
	cmdName   string
	requestID int32
	response  bsoncore.Document
	cmdErr    error
	connID    string
	startTime time.Time
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
	ProcessResponseFn func(response bsoncore.Document, srvr Server, desc description.Server) error

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
	Crypt *Crypt
}

// shouldEncrypt returns true if this operation should automatically be encrypted.
func (op Operation) shouldEncrypt() bool {
	return op.Crypt != nil && !op.Crypt.BypassAutoEncryption
}

// selectServer handles performing server selection for an operation.
func (op Operation) selectServer(ctx context.Context) (Server, error) {
	if err := op.Validate(); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
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

	return op.Deployment.SelectServer(ctx, selector)
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
		return InvalidOperationError{MissingField: "Database"}
	}
	if op.Client != nil && !writeconcern.AckWrite(op.WriteConcern) {
		return errors.New("session provided for an unacknowledged write")
	}
	return nil
}

// Execute runs this operation. The scratch parameter will be used and overwritten (potentially many
// times), this should mainly be used to enable pooling of byte slices.
func (op Operation) Execute(ctx context.Context, scratch []byte) error {
	err := op.Validate()
	if err != nil {
		return err
	}

	srvr, err := op.selectServer(ctx)
	if err != nil {
		return err
	}

	conn, err := srvr.Connection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	desc := description.SelectedServer{Server: conn.Description(), Kind: op.Deployment.Kind()}
	scratch = scratch[:0]

	if desc.WireVersion == nil || desc.WireVersion.Max < 4 {
		switch op.Legacy {
		case LegacyFind:
			return op.legacyFind(ctx, scratch, srvr, conn, desc)
		case LegacyGetMore:
			return op.legacyGetMore(ctx, scratch, srvr, conn, desc)
		case LegacyKillCursors:
			return op.legacyKillCursors(ctx, scratch, srvr, conn, desc)
		}
	}
	if desc.WireVersion == nil || desc.WireVersion.Max < 3 {
		switch op.Legacy {
		case LegacyListCollections:
			return op.legacyListCollections(ctx, scratch, srvr, conn, desc)
		case LegacyListIndexes:
			return op.legacyListIndexes(ctx, scratch, srvr, conn, desc)
		}
	}

	var res bsoncore.Document
	var operationErr WriteCommandError
	var original error
	var retries int
	retryable := op.retryable(desc.Server)
	if retryable && op.RetryMode != nil {
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

			op.Client.RetryWrite = false
			if *op.RetryMode > RetryNone {
				op.Client.RetryWrite = true
				if !op.Client.Committing && !op.Client.Aborting {
					op.Client.IncrementTxnNumber()
				}
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
	batching := op.Batches.Valid()
	for {
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

		// convert to wire message
		if len(scratch) > 0 {
			scratch = scratch[:0]
		}
		wm, startedInfo, err := op.createWireMessage(ctx, scratch, desc)
		if err != nil {
			return err
		}

		// set extra data and send event if possible
		startedInfo.connID = conn.ID()
		startedInfo.cmdName = op.getCommandName(startedInfo.cmd)
		op.publishStartedEvent(ctx, startedInfo)

		// get the moreToCome flag information before we compress
		moreToCome := wiremessage.IsMsgMoreToCome(wm)

		// compress wiremessage if allowed
		if compressor, ok := conn.(Compressor); ok && op.canCompress(startedInfo.cmdName) {
			wm, err = compressor.CompressWireMessage(wm, nil)
			if err != nil {
				return err
			}
		}

		finishedInfo := finishedInformation{
			cmdName:   startedInfo.cmdName,
			requestID: startedInfo.requestID,
			startTime: time.Now(),
			connID:    startedInfo.connID,
		}

		// roundtrip using either the full roundTripper or a special one for when the moreToCome
		// flag is set
		var roundTrip = op.roundTrip
		if moreToCome {
			roundTrip = op.moreToComeRoundTrip
		}
		res, err = roundTrip(ctx, conn, wm)
		if ep, ok := srvr.(ErrorProcessor); ok {
			ep.ProcessError(err)
		}

		finishedInfo.response = res
		finishedInfo.cmdErr = err
		op.publishFinishedEvent(ctx, finishedInfo)

		// Pull out $clusterTime and operationTime and update session and clock. We handle this before
		// handling the error to ensure we are properly gossiping the cluster time.
		op.updateClusterTimes(res)
		op.updateOperationTime(res)
		op.Client.UpdateRecoveryToken(bson.Raw(res))

		// automatically attempt to decrypt all results if client side encryption enabled
		if op.Crypt != nil {
			// use decryptErr isntead of err because err is used below for retrying
			var decryptErr error
			res, decryptErr = op.Crypt.Decrypt(ctx, res)
			if decryptErr != nil {
				return decryptErr
			}
		}
		var perr error
		if op.ProcessResponseFn != nil {
			perr = op.ProcessResponseFn(res, srvr, desc.Server)
		}
		switch tt := err.(type) {
		case WriteCommandError:
			if e := err.(WriteCommandError); retryable && op.Type == Write && e.UnsupportedStorageEngine() {
				return ErrUnsupportedStorageEngine
			}
			if retryable && tt.Retryable() && retries != 0 {
				retries--
				original, err = err, nil
				conn.Close() // Avoid leaking the connection.
				srvr, err = op.selectServer(ctx)
				if err != nil {
					return original
				}
				conn, err = srvr.Connection(ctx)
				if err != nil || conn == nil || !op.retryable(conn.Description()) {
					if conn != nil {
						conn.Close()
					}
					return original
				}
				defer conn.Close() // Avoid leaking the new connection.
				if op.Client != nil && op.Client.Committing {
					// Apply majority write concern for retries
					op.Client.UpdateCommitTransactionWriteConcern()
					op.WriteConcern = op.Client.CurrentWc
				}
				continue
			}
			// If batching is enabled and either ordered is the default (which is true) or
			// explicitly set to true and we have write errors, return the errors.
			if batching && (op.Batches.Ordered == nil || *op.Batches.Ordered == true) && len(tt.WriteErrors) > 0 {
				return tt
			}
			if op.Client != nil && op.Client.Committing && tt.WriteConcernError != nil {
				// When running commitTransaction we return WriteConcernErrors as an Error.
				err := Error{
					Name:    tt.WriteConcernError.Name,
					Code:    int32(tt.WriteConcernError.Code),
					Message: tt.WriteConcernError.Message,
				}
				if err.Code == 64 || err.Code == 50 || tt.WriteConcernError.Retryable() {
					err.Labels = []string{UnknownTransactionCommitResult}
				}
				return err
			}
			operationErr.WriteConcernError = tt.WriteConcernError
			operationErr.WriteErrors = append(operationErr.WriteErrors, tt.WriteErrors...)
		case Error:
			if tt.HasErrorLabel(TransientTransactionError) || tt.HasErrorLabel(UnknownTransactionCommitResult) {
				op.Client.ClearPinnedServer()
			}
			if e := err.(Error); retryable && op.Type == Write && e.UnsupportedStorageEngine() {
				return ErrUnsupportedStorageEngine
			}
			if retryable && tt.Retryable() && retries != 0 {
				retries--
				original, err = err, nil
				conn.Close() // Avoid leaking the connection.
				srvr, err = op.selectServer(ctx)
				if err != nil {
					return original
				}
				conn, err = srvr.Connection(ctx)
				if err != nil || conn == nil || !op.retryable(conn.Description()) {
					if conn != nil {
						conn.Close()
					}
					return original
				}
				defer conn.Close() // Avoid leaking the new connection.
				if op.Client != nil && op.Client.Committing {
					// Apply majority write concern for retries
					op.Client.UpdateCommitTransactionWriteConcern()
					op.WriteConcern = op.Client.CurrentWc
				}
				continue
			}
			if op.Client != nil && op.Client.Committing && (tt.Retryable() || tt.Code == 50) {
				// If we got a retryable error or MaxTimeMSExpired error, we add UnknownTransactionCommitResult.
				tt.Labels = append(tt.Labels, UnknownTransactionCommitResult)
			}
			return tt
		case nil:
			if moreToCome {
				return ErrUnacknowledgedWrite
			}
			if perr != nil {
				return perr
			}
		default:
			return err
		}

		if batching && len(op.Batches.Documents) > 0 {
			if retryable && op.Client != nil && op.RetryMode != nil {
				if *op.RetryMode > RetryNone {
					op.Client.IncrementTxnNumber()
				}
				if *op.RetryMode == RetryOncePerCommand {
					retries = 1
				}
			}
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
		if op.Deployment.SupportsRetryWrites() &&
			desc.WireVersion != nil && desc.WireVersion.Max >= 6 &&
			op.Client != nil && !(op.Client.TransactionInProgress() || op.Client.TransactionStarting()) &&
			writeconcern.AckWrite(op.WriteConcern) {
			return true
		}
	case Read:
		if op.Client != nil && (op.Client.Committing || op.Client.Aborting) {
			return true
		}
		if desc.WireVersion != nil && desc.WireVersion.Max >= 6 &&
			(op.Client == nil || !(op.Client.TransactionInProgress() || op.Client.TransactionStarting())) {
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
		return nil, Error{Message: err.Error(), Labels: labels, Wrapped: err}
	}

	wm, err = conn.ReadWireMessage(ctx, wm[:0])
	if err != nil {
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
		return nil, Error{Message: err.Error(), Labels: labels, Wrapped: err}
	}

	// decompress wiremessage
	wm, err = op.decompressWireMessage(wm)
	if err != nil {
		return nil, err
	}

	// decode
	res, err := op.decodeResult(wm)
	// Pull out $clusterTime and operationTime and update session and clock. We handle this before
	// handling the error to ensure we are properly gossiping the cluster time.
	op.updateClusterTimes(res)
	op.updateOperationTime(res)

	return res, err
}

// moreToComeRoundTrip writes a wiremessage to the provided connection. This is used when an OP_MSG is
// being sent with  the moreToCome bit set.
func (op *Operation) moreToComeRoundTrip(ctx context.Context, conn Connection, wm []byte) ([]byte, error) {
	err := conn.WriteWireMessage(ctx, wm)
	if err != nil {
		if op.Client != nil {
			op.Client.MarkDirty()
		}
		err = Error{Message: err.Error(), Labels: []string{TransientTransactionError, NetworkError}, Wrapped: err}
	}
	return bsoncore.BuildDocument(nil, bsoncore.AppendInt32Element(nil, "ok", 1)), err
}

// decompressWireMessage handles decompressing a wiremessage. If the wiremessage
// is not compressed, this method will return the wiremessage.
func (Operation) decompressWireMessage(wm []byte) ([]byte, error) {
	// read the header and ensure this is a compressed wire message
	length, reqid, respto, opcode, rem, ok := wiremessage.ReadHeader(wm)
	if !ok || len(wm) < int(length) {
		return nil, errors.New("malformed wire message: insufficient bytes")
	}
	if opcode != wiremessage.OpCompressed {
		return wm, nil
	}
	// get the original opcode and uncompressed size
	opcode, rem, ok = wiremessage.ReadCompressedOriginalOpCode(rem)
	if !ok {
		return nil, errors.New("malformed OP_COMPRESSED: missing original opcode")
	}
	uncompressedSize, rem, ok := wiremessage.ReadCompressedUncompressedSize(rem)
	if !ok {
		return nil, errors.New("malformed OP_COMPRESSED: missing uncompressed size")
	}
	// get the compressor ID and decompress the message
	compressorID, rem, ok := wiremessage.ReadCompressedCompressorID(rem)
	if !ok {
		return nil, errors.New("malformed OP_COMPRESSED: missing compressor ID")
	}
	compressedSize := length - 25 // header (16) + original opcode (4) + uncompressed size (4) + compressor ID (1)
	// return the original wiremessage
	msg, rem, ok := wiremessage.ReadCompressedCompressedMessage(rem, compressedSize)
	if !ok {
		return nil, errors.New("malformed OP_COMPRESSED: insufficient bytes for compressed wiremessage")
	}

	header := make([]byte, 0, uncompressedSize+16)
	header = wiremessage.AppendHeader(header, uncompressedSize+16, reqid, respto, opcode)
	opts := CompressionOpts{
		Compressor:       compressorID,
		UncompressedSize: uncompressedSize,
	}
	uncompressed, err := DecompressPayload(msg, opts)
	if err != nil {
		return nil, err
	}

	return append(header, uncompressed...), nil
}

func (op Operation) createWireMessage(ctx context.Context, dst []byte,
	desc description.SelectedServer) ([]byte, startedInformation, error) {

	if desc.WireVersion == nil || desc.WireVersion.Max < wiremessage.OpmsgWireVersion {
		return op.createQueryWireMessage(dst, desc)
	}
	return op.createMsgWireMessage(ctx, dst, desc)
}

func (op Operation) addBatchArray(dst []byte) []byte {
	aidx, dst := bsoncore.AppendArrayElementStart(dst, op.Batches.Identifier)
	for i, doc := range op.Batches.Current {
		dst = bsoncore.AppendDocumentElement(dst, strconv.Itoa(i), doc)
	}
	dst, _ = bsoncore.AppendArrayEnd(dst, aidx)
	return dst
}

func (op Operation) createQueryWireMessage(dst []byte, desc description.SelectedServer) ([]byte, startedInformation, error) {
	var info startedInformation
	flags := op.slaveOK(desc)
	var wmindex int32
	info.requestID = wiremessage.NextRequestID()
	wmindex, dst = wiremessage.AppendHeaderStart(dst, info.requestID, 0, wiremessage.OpQuery)
	dst = wiremessage.AppendQueryFlags(dst, flags)
	// FullCollectionName
	dst = append(dst, op.Database...)
	dst = append(dst, dollarCmd[:]...)
	dst = append(dst, 0x00)
	dst = wiremessage.AppendQueryNumberToSkip(dst, 0)
	dst = wiremessage.AppendQueryNumberToReturn(dst, -1)

	wrapper := int32(-1)
	rp, err := op.createReadPref(desc.Server.Kind, desc.Kind, true)
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

func (op Operation) createMsgWireMessage(ctx context.Context, dst []byte, desc description.SelectedServer) ([]byte, startedInformation, error) {
	var info startedInformation
	var flags wiremessage.MsgFlag
	var wmindex int32
	// We set the MoreToCome bit if we have a write concern, it's unacknowledged, and we either
	// aren't batching or we are encoding the last batch.
	if op.WriteConcern != nil && !writeconcern.AckWrite(op.WriteConcern) && (op.Batches == nil || len(op.Batches.Documents) == 0) {
		flags = wiremessage.MoreToCome
	}
	info.requestID = wiremessage.NextRequestID()
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

	dst = bsoncore.AppendStringElement(dst, "$db", op.Database)
	rp, err := op.createReadPref(desc.Server.Kind, desc.Kind, false)
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

	if rc == nil {
		return dst, nil
	}

	_, data, err := rc.MarshalBSONValue() // always returns a document
	if err != nil {
		return dst, err
	}

	if description.SessionsSupported(desc.WireVersion) && client != nil && client.Consistent && client.OperationTime != nil {
		data = data[:len(data)-1] // remove the null byte
		data = bsoncore.AppendTimestampElement(data, "afterClusterTime", client.OperationTime.T, client.OperationTime.I)
		data, _ = bsoncore.AppendDocumentEnd(data, 0)
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
	if err == writeconcern.ErrEmptyWriteConcern {
		return dst, nil
	}
	if err != nil {
		return dst, err
	}

	return append(bsoncore.AppendHeader(dst, t, "writeConcern"), data...), nil
}

func (op Operation) addSession(dst []byte, desc description.SelectedServer) ([]byte, error) {
	client := op.Client
	if client == nil || !description.SessionsSupported(desc.WireVersion) || desc.SessionTimeoutMinutes == 0 {
		return dst, nil
	}
	if client.Terminated {
		return dst, session.ErrSessionEnded
	}
	lsid, _ := client.SessionID.MarshalBSON()
	dst = bsoncore.AppendDocumentElement(dst, "lsid", lsid)

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

	client.ApplyCommand(desc.Server)

	return dst, nil
}

func (op Operation) addClusterTime(dst []byte, desc description.SelectedServer) []byte {
	client, clock := op.Client, op.Clock
	if (clock == nil && client == nil) || !description.SessionsSupported(desc.WireVersion) {
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

func (op Operation) createReadPref(serverKind description.ServerKind, topologyKind description.TopologyKind, isOpQuery bool) (bsoncore.Document, error) {
	if serverKind == description.Standalone || (isOpQuery && serverKind != description.Mongos) || op.Type == Write {
		// Don't send read preference for non-mongos when using OP_QUERY or for all standalones
		return nil, nil
	}

	idx, doc := bsoncore.AppendDocumentStart(nil)
	rp, err := op.getReadPrefBasedOnTransaction()
	if err != nil {
		return nil, err
	}

	if rp == nil {
		if topologyKind == description.Single && serverKind != description.Mongos {
			doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
			doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
			return doc, nil
		}
		return nil, nil
	}

	switch rp.Mode() {
	case readpref.PrimaryMode:
		if serverKind == description.Mongos {
			return nil, nil
		}
		if topologyKind == description.Single {
			doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
			doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
			return doc, nil
		}
		doc = bsoncore.AppendStringElement(doc, "mode", "primary")
	case readpref.PrimaryPreferredMode:
		doc = bsoncore.AppendStringElement(doc, "mode", "primaryPreferred")
	case readpref.SecondaryPreferredMode:
		_, ok := rp.MaxStaleness()
		if serverKind == description.Mongos && isOpQuery && !ok && len(rp.TagSets()) == 0 {
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
		if len(ts) == 0 {
			continue
		}
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

	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
	return doc, nil
}

func (op Operation) slaveOK(desc description.SelectedServer) wiremessage.QueryFlag {
	if desc.Kind == description.Single && desc.Server.Kind != description.Mongos {
		return wiremessage.SlaveOK
	}

	if rp := op.ReadPreference; rp != nil && rp.Mode() != readpref.PrimaryMode {
		return wiremessage.SlaveOK
	}

	return 0
}

func (Operation) canCompress(cmd string) bool {
	if cmd == "isMaster" || cmd == "saslStart" || cmd == "saslContinue" || cmd == "getnonce" || cmd == "authenticate" ||
		cmd == "createUser" || cmd == "updateUser" || cmd == "copydbSaslStart" || cmd == "copydbgetnonce" || cmd == "copydb" {
		return false
	}
	return true
}

// decodeOpReply extracts the necessary information from an OP_REPLY wire message.
// includesHeader: specifies whether or not wm includes the message header
// Returns the decoded OP_REPLY. If the err field of the returned opReply is non-nil, an error occurred while decoding
// or validating the response and the other fields are undefined.
func (Operation) decodeOpReply(wm []byte, includesHeader bool) opReply {
	var reply opReply
	var ok bool

	if includesHeader {
		wmLength := len(wm)
		var length int32
		var opcode wiremessage.OpCode
		length, _, _, opcode, wm, ok = wiremessage.ReadHeader(wm)
		if !ok || int(length) > wmLength {
			reply.err = errors.New("malformed wire message: insufficient bytes")
			return reply
		}
		if opcode != wiremessage.OpReply {
			reply.err = errors.New("malformed wire message: incorrect opcode")
			return reply
		}
	}

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
	reply.documents, wm, ok = wiremessage.ReadReplyDocuments(wm)
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

func (op Operation) decodeResult(wm []byte) (bsoncore.Document, error) {
	wmLength := len(wm)
	length, _, _, opcode, wm, ok := wiremessage.ReadHeader(wm)
	if !ok || int(length) > wmLength {
		return nil, errors.New("malformed wire message: insufficient bytes")
	}

	wm = wm[:wmLength-16] // constrain to just this wiremessage, incase there are multiple in the slice

	switch opcode {
	case wiremessage.OpReply:
		reply := op.decodeOpReply(wm, false)
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

		return rdr, extractError(rdr)
	case wiremessage.OpMsg:
		_, wm, ok = wiremessage.ReadMsgFlags(wm)
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
				// TODO(GODRIVER-617): Implement document sequence returns.
				_, _, wm, ok = wiremessage.ReadMsgSectionDocumentSequence(wm)
				if !ok {
					return nil, errors.New("malformed wire message: insufficient bytes to read document sequence")
				}
			default:
				return nil, fmt.Errorf("malformed wire message: uknown section type %v", stype)
			}
		}

		err := res.Validate()
		if err != nil {
			return nil, NewCommandResponseError("malformed OP_MSG: invalid document", err)
		}

		return res, extractError(res)
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

func (op *Operation) canMonitor(cmd string) bool {
	return !(cmd == "authenticate" || cmd == "saslStart" || cmd == "saslContinue" || cmd == "getnonce" || cmd == "createUser" ||
		cmd == "updateUser" || cmd == "copydbgetnonce" || cmd == "copydbsaslstart" || cmd == "copydb")
}

// publishStartedEvent publishes a CommandStartedEvent to the operation's command monitor if possible. If the command is
// an unacknowledged write, a CommandSucceededEvent will be published as well. If started events are not being monitored,
// no events are published.
func (op Operation) publishStartedEvent(ctx context.Context, info startedInformation) {
	if op.CommandMonitor == nil || op.CommandMonitor.Started == nil {
		return
	}

	// Make a copy of the command. Redact if the command is security sensitive and cannot be monitored.
	// If there was a type 1 payload for the current batch, convert it to a BSON array.
	var cmdCopy []byte
	if op.canMonitor(info.cmdName) {
		cmdCopy = make([]byte, len(info.cmd))
		copy(cmdCopy, info.cmd)
		if info.documentSequenceIncluded {
			cmdCopy = cmdCopy[:len(info.cmd)-1] // remove 0 byte at end
			cmdCopy = op.addBatchArray(cmdCopy)
			cmdCopy, _ = bsoncore.AppendDocumentEnd(cmdCopy, 0) // add back 0 byte and update length
		}
	}

	started := &event.CommandStartedEvent{
		Command:      cmdCopy,
		DatabaseName: op.Database,
		CommandName:  info.cmdName,
		RequestID:    int64(info.requestID),
		ConnectionID: info.connID,
	}
	op.CommandMonitor.Started(ctx, started)
}

// publishFinishedEvent publishes either a CommandSucceededEvent or a CommandFailedEvent to the operation's command
// monitor if possible. If success/failure events aren't being monitored, no events are published.
func (op Operation) publishFinishedEvent(ctx context.Context, info finishedInformation) {
	success := info.cmdErr == nil
	if _, ok := info.cmdErr.(WriteCommandError); ok {
		success = true
	}
	if op.CommandMonitor == nil || (success && op.CommandMonitor.Succeeded == nil) || (!success && op.CommandMonitor.Failed == nil) {
		return
	}

	var durationNanos int64
	var emptyTime time.Time
	if info.startTime != emptyTime {
		durationNanos = time.Now().Sub(info.startTime).Nanoseconds()
	}

	finished := event.CommandFinishedEvent{
		CommandName:   info.cmdName,
		RequestID:     int64(info.requestID),
		ConnectionID:  info.connID,
		DurationNanos: durationNanos,
	}

	if success {
		res := bson.Raw{}
		// Only copy the reply for commands that are not security sensitive
		if op.canMonitor(info.cmdName) {
			res = make([]byte, len(info.response))
			copy(res, info.response)
		}
		successEvent := &event.CommandSucceededEvent{
			Reply:                res,
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
