package driver

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// BatchCursor is a batch implementation of a cursor. It returns documents in entire batches instead
// of one at a time. An individual document cursor can be built on top of this batch cursor.
type BatchCursor struct {
	clientSession        *session.Client
	clock                *session.ClusterClock
	database             string
	collection           string
	id                   int64
	err                  error
	server               Server
	errorProcessor       ErrorProcessor // This will only be set when pinning to a connection.
	connection           PinnedConnection
	batchSize            int32
	maxTimeMS            int64
	currentBatch         *bsoncore.DocumentSequence
	firstBatch           bool
	cmdMonitor           *event.CommandMonitor
	postBatchResumeToken bsoncore.Document
	crypt                *Crypt
	serverAPI            *ServerAPIOptions

	// legacy server (< 3.2) fields
	legacy      bool // This field is provided for ListCollectionsBatchCursor.
	limit       int32
	numReturned int32 // number of docs returned by server
}

// CursorResponse represents the response from a command the results in a cursor. A BatchCursor can
// be constructed from a CursorResponse.
type CursorResponse struct {
	Server               Server
	ErrorProcessor       ErrorProcessor // This will only be set when pinning to a connection.
	Connection           PinnedConnection
	Desc                 description.Server
	FirstBatch           *bsoncore.DocumentSequence
	Database             string
	Collection           string
	ID                   int64
	postBatchResumeToken bsoncore.Document
}

// NewCursorResponse constructs a cursor response from the given response and server. This method
// can be used within the ProcessResponse method for an operation.
func NewCursorResponse(info ResponseInfo) (CursorResponse, error) {
	response := info.ServerResponse
	cur, ok := response.Lookup("cursor").DocumentOK()
	if !ok {
		return CursorResponse{}, fmt.Errorf("cursor should be an embedded document but is of BSON type %s", response.Lookup("cursor").Type)
	}
	elems, err := cur.Elements()
	if err != nil {
		return CursorResponse{}, err
	}
	curresp := CursorResponse{Server: info.Server, Desc: info.ConnectionDescription}

	for _, elem := range elems {
		switch elem.Key() {
		case "firstBatch":
			arr, ok := elem.Value().ArrayOK()
			if !ok {
				return CursorResponse{}, fmt.Errorf("firstBatch should be an array but is a BSON %s", elem.Value().Type)
			}
			curresp.FirstBatch = &bsoncore.DocumentSequence{Style: bsoncore.ArrayStyle, Data: arr}
		case "ns":
			ns, ok := elem.Value().StringValueOK()
			if !ok {
				return CursorResponse{}, fmt.Errorf("ns should be a string but is a BSON %s", elem.Value().Type)
			}
			index := strings.Index(ns, ".")
			if index == -1 {
				return CursorResponse{}, errors.New("ns field must contain a valid namespace, but is missing '.'")
			}
			curresp.Database = ns[:index]
			curresp.Collection = ns[index+1:]
		case "id":
			curresp.ID, ok = elem.Value().Int64OK()
			if !ok {
				return CursorResponse{}, fmt.Errorf("id should be an int64 but it is a BSON %s", elem.Value().Type)
			}
		case "postBatchResumeToken":
			curresp.postBatchResumeToken, ok = elem.Value().DocumentOK()
			if !ok {
				return CursorResponse{}, fmt.Errorf("post batch resume token should be a document but it is a BSON %s", elem.Value().Type)
			}
		}
	}

	// If the deployment is behind a load balancer and the cursor has a non-zero ID, pin the cursor to a connection and
	// use the same connection to execute getMore and killCursors commands.
	if curresp.Desc.LoadBalanced() && curresp.ID != 0 {
		// Cache the server as an ErrorProcessor to use when constructing deployments for cursor commands.
		ep, ok := curresp.Server.(ErrorProcessor)
		if !ok {
			return CursorResponse{}, fmt.Errorf("expected Server used to establish a cursor to implement ErrorProcessor, but got %T", curresp.Server)
		}
		curresp.ErrorProcessor = ep

		refConn, ok := info.Connection.(PinnedConnection)
		if !ok {
			return CursorResponse{}, fmt.Errorf("expected Connection used to establish a cursor to implement PinnedConnection, but got %T", info.Connection)
		}
		if err := refConn.PinToCursor(); err != nil {
			return CursorResponse{}, fmt.Errorf("error incrementing connection reference count when creating a cursor: %v", err)
		}
		curresp.Connection = refConn
	}

	return curresp, nil
}

// CursorOptions are extra options that are required to construct a BatchCursor.
type CursorOptions struct {
	BatchSize      int32
	MaxTimeMS      int64
	Limit          int32
	CommandMonitor *event.CommandMonitor
	Crypt          *Crypt
	ServerAPI      *ServerAPIOptions
}

// NewBatchCursor creates a new BatchCursor from the provided parameters.
func NewBatchCursor(cr CursorResponse, clientSession *session.Client, clock *session.ClusterClock, opts CursorOptions) (*BatchCursor, error) {
	ds := cr.FirstBatch
	bc := &BatchCursor{
		clientSession:        clientSession,
		clock:                clock,
		database:             cr.Database,
		collection:           cr.Collection,
		id:                   cr.ID,
		server:               cr.Server,
		connection:           cr.Connection,
		errorProcessor:       cr.ErrorProcessor,
		batchSize:            opts.BatchSize,
		maxTimeMS:            opts.MaxTimeMS,
		cmdMonitor:           opts.CommandMonitor,
		firstBatch:           true,
		postBatchResumeToken: cr.postBatchResumeToken,
		crypt:                opts.Crypt,
		serverAPI:            opts.ServerAPI,
	}

	if ds != nil {
		bc.numReturned = int32(ds.DocumentCount())
	}
	if cr.Desc.WireVersion == nil || cr.Desc.WireVersion.Max < 4 {
		bc.legacy = true
		bc.limit = opts.Limit

		// Take as many documents from the batch as needed.
		if bc.limit != 0 && bc.limit < bc.numReturned {
			for i := int32(0); i < bc.limit; i++ {
				_, err := ds.Next()
				if err != nil {
					return nil, err
				}
			}
			ds.Data = ds.Data[:ds.Pos]
			ds.ResetIterator()
		}
	}

	bc.currentBatch = ds
	return bc, nil
}

// NewEmptyBatchCursor returns a batch cursor that is empty.
func NewEmptyBatchCursor() *BatchCursor {
	return &BatchCursor{currentBatch: new(bsoncore.DocumentSequence)}
}

// ID returns the cursor ID for this batch cursor.
func (bc *BatchCursor) ID() int64 {
	return bc.id
}

// Next indicates if there is another batch available. Returning false does not necessarily indicate
// that the cursor is closed. This method will return false when an empty batch is returned.
//
// If Next returns true, there is a valid batch of documents available. If Next returns false, there
// is not a valid batch of documents available.
func (bc *BatchCursor) Next(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}

	if bc.firstBatch {
		bc.firstBatch = false
		return !bc.currentBatch.Empty()
	}

	if bc.id == 0 || bc.server == nil {
		return false
	}

	bc.getMore(ctx)

	return !bc.currentBatch.Empty()
}

// Batch will return a DocumentSequence for the current batch of documents. The returned
// DocumentSequence is only valid until the next call to Next or Close.
func (bc *BatchCursor) Batch() *bsoncore.DocumentSequence { return bc.currentBatch }

// Err returns the latest error encountered.
func (bc *BatchCursor) Err() error { return bc.err }

// Close closes this batch cursor.
func (bc *BatchCursor) Close(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	err := bc.KillCursor(ctx)
	bc.id = 0
	bc.currentBatch.Data = nil
	bc.currentBatch.Style = 0
	bc.currentBatch.ResetIterator()

	connErr := bc.unpinConnection()
	if err == nil {
		err = connErr
	}
	return err
}

func (bc *BatchCursor) unpinConnection() error {
	if bc.connection == nil {
		return nil
	}

	err := bc.connection.UnpinFromCursor()
	closeErr := bc.connection.Close()
	if err == nil && closeErr != nil {
		err = closeErr
	}
	bc.connection = nil
	return err
}

// Server returns the server for this cursor.
func (bc *BatchCursor) Server() Server {
	return bc.server
}

func (bc *BatchCursor) clearBatch() {
	bc.currentBatch.Data = bc.currentBatch.Data[:0]
}

// KillCursor kills cursor on server without closing batch cursor
func (bc *BatchCursor) KillCursor(ctx context.Context) error {
	if bc.server == nil || bc.id == 0 {
		return nil
	}

	return Operation{
		CommandFn: func(dst []byte, desc description.SelectedServer) ([]byte, error) {
			dst = bsoncore.AppendStringElement(dst, "killCursors", bc.collection)
			dst = bsoncore.BuildArrayElement(dst, "cursors", bsoncore.Value{Type: bsontype.Int64, Data: bsoncore.AppendInt64(nil, bc.id)})
			return dst, nil
		},
		Database:       bc.database,
		Deployment:     bc.getOperationDeployment(),
		Client:         bc.clientSession,
		Clock:          bc.clock,
		Legacy:         LegacyKillCursors,
		CommandMonitor: bc.cmdMonitor,
		ServerAPI:      bc.serverAPI,
	}.Execute(ctx, nil)
}

func (bc *BatchCursor) getMore(ctx context.Context) {
	bc.clearBatch()
	if bc.id == 0 {
		return
	}

	// Required for legacy operations which don't support limit.
	numToReturn := bc.batchSize
	if bc.limit != 0 && bc.numReturned+bc.batchSize >= bc.limit {
		numToReturn = bc.limit - bc.numReturned
		if numToReturn <= 0 {
			err := bc.Close(ctx)
			if err != nil {
				bc.err = err
			}
			return
		}
	}

	bc.err = Operation{
		CommandFn: func(dst []byte, desc description.SelectedServer) ([]byte, error) {
			dst = bsoncore.AppendInt64Element(dst, "getMore", bc.id)
			dst = bsoncore.AppendStringElement(dst, "collection", bc.collection)
			if numToReturn > 0 {
				dst = bsoncore.AppendInt32Element(dst, "batchSize", numToReturn)
			}
			if bc.maxTimeMS > 0 {
				dst = bsoncore.AppendInt64Element(dst, "maxTimeMS", bc.maxTimeMS)
			}
			return dst, nil
		},
		Database:   bc.database,
		Deployment: bc.getOperationDeployment(),
		ProcessResponseFn: func(info ResponseInfo) error {
			response := info.ServerResponse
			id, ok := response.Lookup("cursor", "id").Int64OK()
			if !ok {
				return fmt.Errorf("cursor.id should be an int64 but is a BSON %s", response.Lookup("cursor", "id").Type)
			}
			bc.id = id

			batch, ok := response.Lookup("cursor", "nextBatch").ArrayOK()
			if !ok {
				return fmt.Errorf("cursor.nextBatch should be an array but is a BSON %s", response.Lookup("cursor", "nextBatch").Type)
			}
			bc.currentBatch.Style = bsoncore.ArrayStyle
			bc.currentBatch.Data = batch
			bc.currentBatch.ResetIterator()
			bc.numReturned += int32(bc.currentBatch.DocumentCount()) // Required for legacy operations which don't support limit.

			pbrt, err := response.LookupErr("cursor", "postBatchResumeToken")
			if err != nil {
				// I don't really understand why we don't set bc.err here
				return nil
			}

			pbrtDoc, ok := pbrt.DocumentOK()
			if !ok {
				bc.err = fmt.Errorf("expected BSON type for post batch resume token to be EmbeddedDocument but got %s", pbrt.Type)
				return nil
			}

			bc.postBatchResumeToken = bsoncore.Document(pbrtDoc)

			return nil
		},
		Client:         bc.clientSession,
		Clock:          bc.clock,
		Legacy:         LegacyGetMore,
		CommandMonitor: bc.cmdMonitor,
		Crypt:          bc.crypt,
		ServerAPI:      bc.serverAPI,
	}.Execute(ctx, nil)

	// Once the cursor has been drained, we can unpin the connection if one is currently pinned.
	if bc.id == 0 {
		err := bc.unpinConnection()
		if err != nil && bc.err == nil {
			bc.err = err
		}
	}

	// If we're in load balanced mode and the pinned connection encounters a network error, we should not use it for
	// future commands. Per the spec, the connection will not be unpinned until the cursor is actually closed, but
	// we set the cursor ID to 0 to ensure the Close() call will not execute a killCursors command.
	if driverErr, ok := bc.err.(Error); ok && driverErr.NetworkError() && bc.connection != nil {
		bc.id = 0
	}

	// Required for legacy operations which don't support limit.
	if bc.limit != 0 && bc.numReturned >= bc.limit {
		// call KillCursor instead of Close because Close will clear out the data for the current batch.
		err := bc.KillCursor(ctx)
		if err != nil && bc.err == nil {
			bc.err = err
		}
	}
	return
}

// PostBatchResumeToken returns the latest seen post batch resume token.
func (bc *BatchCursor) PostBatchResumeToken() bsoncore.Document {
	return bc.postBatchResumeToken
}

// SetBatchSize sets the batchSize for future getMores.
func (bc *BatchCursor) SetBatchSize(size int32) {
	bc.batchSize = size
}

func (bc *BatchCursor) getOperationDeployment() Deployment {
	if bc.connection != nil {
		return &loadBalancedCursorDeployment{
			errorProcessor: bc.errorProcessor,
			conn:           bc.connection,
		}
	}
	return SingleServerDeployment{bc.server}
}

// loadBalancedCursorDeployment is used as a Deployment for getMore and killCursors commands when pinning to a
// connection in load balanced mode. This type also functions as an ErrorProcessor to ensure that SDAM errors are
// handled for these commands in this mode.
type loadBalancedCursorDeployment struct {
	errorProcessor ErrorProcessor
	conn           PinnedConnection
}

var _ Deployment = (*loadBalancedCursorDeployment)(nil)
var _ Server = (*loadBalancedCursorDeployment)(nil)
var _ ErrorProcessor = (*loadBalancedCursorDeployment)(nil)

func (lbcd *loadBalancedCursorDeployment) SelectServer(_ context.Context, _ description.ServerSelector) (Server, error) {
	return lbcd, nil
}

func (lbcd *loadBalancedCursorDeployment) Kind() description.TopologyKind {
	return description.LoadBalanced
}

func (lbcd *loadBalancedCursorDeployment) Connection(_ context.Context) (Connection, error) {
	return lbcd.conn, nil
}

func (lbcd *loadBalancedCursorDeployment) ProcessError(err error, conn Connection) ProcessErrorResult {
	return lbcd.errorProcessor.ProcessError(err, conn)
}
