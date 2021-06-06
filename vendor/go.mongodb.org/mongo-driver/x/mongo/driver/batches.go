package driver

import (
	"errors"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// this is the amount of reserved buffer space in a message that the
// driver reserves for command overhead.
const reservedCommandBufferBytes = 16 * 10 * 10 * 10

// ErrDocumentTooLarge occurs when a document that is larger than the maximum size accepted by a
// server is passed to an insert command.
var ErrDocumentTooLarge = errors.New("an inserted document is too large")

// Batches contains the necessary information to batch split an operation. This is only used for write
// oeprations.
type Batches struct {
	Identifier string
	Documents  []bsoncore.Document
	Current    []bsoncore.Document
	Ordered    *bool
}

// Valid returns true if Batches contains both an identifier and the length of Documents is greater
// than zero.
func (b *Batches) Valid() bool { return b != nil && b.Identifier != "" && len(b.Documents) > 0 }

// ClearBatch clears the Current batch. This must be called before AdvanceBatch will advance to the
// next batch.
func (b *Batches) ClearBatch() { b.Current = b.Current[:0] }

// AdvanceBatch splits the next batch using maxCount and targetBatchSize. This method will do nothing if
// the current batch has not been cleared. We do this so that when this is called during execute we
// can call it without first needing to check if we already have a batch, which makes the code
// simpler and makes retrying easier.
// The maxDocSize parameter is used to check that any one document is not too large. If the first document is bigger
// than targetBatchSize but smaller than maxDocSize, a batch of size 1 containing that document will be created.
func (b *Batches) AdvanceBatch(maxCount, targetBatchSize, maxDocSize int) error {
	if len(b.Current) > 0 {
		return nil
	}

	if maxCount <= 0 {
		maxCount = 1
	}

	splitAfter := 0
	size := 0
	for i, doc := range b.Documents {
		if i == maxCount {
			break
		}
		if len(doc) > maxDocSize {
			return ErrDocumentTooLarge
		}
		if size+len(doc) > targetBatchSize {
			break
		}

		size += len(doc)
		splitAfter++
	}

	// if there are no documents, take the first one.
	// this can happen if there is a document that is smaller than maxDocSize but greater than targetBatchSize.
	if splitAfter == 0 {
		splitAfter = 1
	}

	b.Current, b.Documents = b.Documents[:splitAfter], b.Documents[splitAfter:]
	return nil
}
