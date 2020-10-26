package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

// batchCursor is the interface implemented by types that can provide batches of document results.
// The Cursor type is built on top of this type.
type batchCursor interface {
	// ID returns the ID of the cursor.
	ID() int64

	// Next returns true if there is a batch available.
	Next(context.Context) bool

	// Batch will return a DocumentSequence for the current batch of documents. The returned
	// DocumentSequence is only valid until the next call to Next or Close.
	Batch() *bsoncore.DocumentSequence

	// Server returns a pointer to the cursor's server.
	Server() driver.Server

	// Err returns the last error encountered.
	Err() error

	// Close closes the cursor.
	Close(context.Context) error
}

// changeStreamCursor is the interface implemented by batch cursors that also provide the functionality for retrieving
// a postBatchResumeToken from commands and allows for the cursor to be killed rather than closed
type changeStreamCursor interface {
	batchCursor
	// PostBatchResumeToken returns the latest seen post batch resume token.
	PostBatchResumeToken() bsoncore.Document

	// KillCursor kills cursor on server without closing batch cursor
	KillCursor(context.Context) error
}
