package driver

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
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
	batchSize            int32
	maxTimeMS            int64
	currentBatch         *bsoncore.DocumentSequence
	firstBatch           bool
	cmdMonitor           *event.CommandMonitor
	postBatchResumeToken bsoncore.Document
	crypt                *Crypt

	// legacy server (< 3.2) fields
	legacy      bool // This field is provided for ListCollectionsBatchCursor.
	limit       int32
	numReturned int32 // number of docs returned by server
}

// CursorResponse represents the response from a command the results in a cursor. A BatchCursor can
// be constructed from a CursorResponse.
type CursorResponse struct {
	Server               Server
	Desc                 description.Server
	FirstBatch           *bsoncore.DocumentSequence
	Database             string
	Collection           string
	ID                   int64
	postBatchResumeToken bsoncore.Document
}

// NewCursorResponse constructs a cursor response from the given response and server. This method
// can be used within the ProcessResponse method for an operation.
func NewCursorResponse(response bsoncore.Document, server Server, desc description.Server) (CursorResponse, error) {
	cur, ok := response.Lookup("cursor").DocumentOK()
	if !ok {
		return CursorResponse{}, fmt.Errorf("cursor should be an embedded document but is of BSON type %s", response.Lookup("cursor").Type)
	}
	elems, err := cur.Elements()
	if err != nil {
		return CursorResponse{}, err
	}
	curresp := CursorResponse{Server: server, Desc: desc}

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
	return curresp, nil
}

// CursorOptions are extra options that are required to construct a BatchCursor.
type CursorOptions struct {
	BatchSize      int32
	MaxTimeMS      int64
	Limit          int32
	CommandMonitor *event.CommandMonitor
	Crypt          *Crypt
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
		batchSize:            opts.BatchSize,
		maxTimeMS:            opts.MaxTimeMS,
		cmdMonitor:           opts.CommandMonitor,
		firstBatch:           true,
		postBatchResumeToken: cr.postBatchResumeToken,
		crypt:                opts.Crypt,
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
		Deployment:     SingleServerDeployment{Server: bc.server},
		Client:         bc.clientSession,
		Clock:          bc.clock,
		Legacy:         LegacyKillCursors,
		CommandMonitor: bc.cmdMonitor,
	}.Execute(ctx, nil)
}

func (bc *BatchCursor) getMore(ctx context.Context) {
	bc.clearBatch()
	if bc.id == 0 {
		return
	}

	// Required for legacy operations which don't support limit.
	numToReturn := bc.batchSize
	if bc.limit != 0 && bc.numReturned+bc.batchSize > bc.limit {
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
		Deployment: SingleServerDeployment{Server: bc.server},
		ProcessResponseFn: func(response bsoncore.Document, srvr Server, desc description.Server, currIndex int) error {
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
	}.Execute(ctx, nil)

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
