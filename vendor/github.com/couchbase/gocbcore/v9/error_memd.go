package gocbcore

import (
	"errors"
	"log"

	"github.com/couchbase/gocbcore/v9/memd"
)

var statusCodeErrorMap = make(map[memd.StatusCode]error)

func makeKvStatusError(code memd.StatusCode) error {
	err := errors.New(code.KVText())
	if statusCodeErrorMap[code] != nil {
		log.Fatal("error handling setup failure")
	}
	statusCodeErrorMap[code] = err
	return err
}

func getKvStatusCodeError(code memd.StatusCode) error {
	if err := statusCodeErrorMap[code]; err != nil {
		return err
	}
	return errors.New(code.KVText())
}

var (
	// ErrMemdKeyNotFound occurs when an operation is performed on a key that does not exist.
	ErrMemdKeyNotFound = makeKvStatusError(memd.StatusKeyNotFound)

	// ErrMemdKeyExists occurs when an operation is performed on a key that could not be found.
	ErrMemdKeyExists = makeKvStatusError(memd.StatusKeyExists)

	// ErrMemdTooBig occurs when an operation attempts to store more data in a single document
	// than the server is capable of storing (by default, this is a 20MB limit).
	ErrMemdTooBig = makeKvStatusError(memd.StatusTooBig)

	// ErrMemdInvalidArgs occurs when the server receives invalid arguments for an operation.
	ErrMemdInvalidArgs = makeKvStatusError(memd.StatusInvalidArgs)

	// ErrMemdNotStored occurs when the server fails to store a key.
	ErrMemdNotStored = makeKvStatusError(memd.StatusNotStored)

	// ErrMemdBadDelta occurs when an invalid delta value is specified to a counter operation.
	ErrMemdBadDelta = makeKvStatusError(memd.StatusBadDelta)

	// ErrMemdNotMyVBucket occurs when an operation is dispatched to a server which is
	// non-authoritative for a specific vbucket.
	ErrMemdNotMyVBucket = makeKvStatusError(memd.StatusNotMyVBucket)

	// ErrMemdNoBucket occurs when no bucket was selected on a connection.
	ErrMemdNoBucket = makeKvStatusError(memd.StatusNoBucket)

	// ErrMemdLocked occurs when a document is already locked.
	ErrMemdLocked = makeKvStatusError(memd.StatusLocked)

	// ErrMemdAuthStale occurs when authentication credentials have become invalidated.
	ErrMemdAuthStale = makeKvStatusError(memd.StatusAuthStale)

	// ErrMemdAuthError occurs when the authentication information provided was not valid.
	ErrMemdAuthError = makeKvStatusError(memd.StatusAuthError)

	// ErrMemdAuthContinue occurs in multi-step authentication when more authentication
	// work needs to be performed in order to complete the authentication process.
	ErrMemdAuthContinue = makeKvStatusError(memd.StatusAuthContinue)

	// ErrMemdRangeError occurs when the range specified to the server is not valid.
	ErrMemdRangeError = makeKvStatusError(memd.StatusRangeError)

	// ErrMemdRollback occurs when a DCP stream fails to open due to a rollback having
	// previously occurred since the last time the stream was opened.
	ErrMemdRollback = makeKvStatusError(memd.StatusRollback)

	// ErrMemdAccessError occurs when an access error occurs.
	ErrMemdAccessError = makeKvStatusError(memd.StatusAccessError)

	// ErrMemdNotInitialized is sent by servers which are still initializing, and are not
	// yet ready to accept operations on behalf of a particular bucket.
	ErrMemdNotInitialized = makeKvStatusError(memd.StatusNotInitialized)

	// ErrMemdUnknownCommand occurs when an unknown operation is sent to a server.
	ErrMemdUnknownCommand = makeKvStatusError(memd.StatusUnknownCommand)

	// ErrMemdOutOfMemory occurs when the server cannot service a request due to memory
	// limitations.
	ErrMemdOutOfMemory = makeKvStatusError(memd.StatusOutOfMemory)

	// ErrMemdNotSupported occurs when an operation is understood by the server, but that
	// operation is not supported on this server (occurs for a variety of reasons).
	ErrMemdNotSupported = makeKvStatusError(memd.StatusNotSupported)

	// ErrMemdInternalError occurs when internal errors prevent the server from processing
	// your request.
	ErrMemdInternalError = makeKvStatusError(memd.StatusInternalError)

	// ErrMemdBusy occurs when the server is too busy to process your request right away.
	// Attempting the operation at a later time will likely succeed.
	ErrMemdBusy = makeKvStatusError(memd.StatusBusy)

	// ErrMemdTmpFail occurs when a temporary failure is preventing the server from
	// processing your request.
	ErrMemdTmpFail = makeKvStatusError(memd.StatusTmpFail)

	// ErrMemdCollectionNotFound occurs when a Collection cannot be found.
	ErrMemdCollectionNotFound = makeKvStatusError(memd.StatusCollectionUnknown)

	// ErrMemdScopeNotFound occurs when a Collection cannot be found.
	ErrMemdScopeNotFound = makeKvStatusError(memd.StatusScopeUnknown)

	// ErrMemdDurabilityInvalidLevel occurs when an invalid durability level was requested.
	ErrMemdDurabilityInvalidLevel = makeKvStatusError(memd.StatusDurabilityInvalidLevel)

	// ErrMemdDurabilityImpossible occurs when a request is performed with impossible
	//  durability level requirements.
	ErrMemdDurabilityImpossible = makeKvStatusError(memd.StatusDurabilityImpossible)

	// ErrMemdSyncWriteInProgess occurs when an attempt is made to write to a key that has
	//  a SyncWrite pending.
	ErrMemdSyncWriteInProgess = makeKvStatusError(memd.StatusSyncWriteInProgress)

	// ErrMemdSyncWriteAmbiguous occurs when an SyncWrite does not complete in the specified
	// time and the result is ambiguous.
	ErrMemdSyncWriteAmbiguous = makeKvStatusError(memd.StatusSyncWriteAmbiguous)

	// ErrMemdSyncWriteReCommitInProgress occurs when an SyncWrite is being recommitted.
	ErrMemdSyncWriteReCommitInProgress = makeKvStatusError(memd.StatusSyncWriteReCommitInProgress)

	// ErrMemdSubDocPathNotFound occurs when a sub-document operation targets a path
	// which does not exist in the specifie document.
	ErrMemdSubDocPathNotFound = makeKvStatusError(memd.StatusSubDocPathNotFound)

	// ErrMemdSubDocPathMismatch occurs when a sub-document operation specifies a path
	// which does not match the document structure (field access on an array).
	ErrMemdSubDocPathMismatch = makeKvStatusError(memd.StatusSubDocPathMismatch)

	// ErrMemdSubDocPathInvalid occurs when a sub-document path could not be parsed.
	ErrMemdSubDocPathInvalid = makeKvStatusError(memd.StatusSubDocPathInvalid)

	// ErrMemdSubDocPathTooBig occurs when a sub-document path is too big.
	ErrMemdSubDocPathTooBig = makeKvStatusError(memd.StatusSubDocPathTooBig)

	// ErrMemdSubDocDocTooDeep occurs when an operation would cause a document to be
	// nested beyond the depth limits allowed by the sub-document specification.
	ErrMemdSubDocDocTooDeep = makeKvStatusError(memd.StatusSubDocDocTooDeep)

	// ErrMemdSubDocCantInsert occurs when a sub-document operation could not insert.
	ErrMemdSubDocCantInsert = makeKvStatusError(memd.StatusSubDocCantInsert)

	// ErrMemdSubDocNotJSON occurs when a sub-document operation is performed on a
	// document which is not JSON.
	ErrMemdSubDocNotJSON = makeKvStatusError(memd.StatusSubDocNotJSON)

	// ErrMemdSubDocBadRange occurs when a sub-document operation is performed with
	// a bad range.
	ErrMemdSubDocBadRange = makeKvStatusError(memd.StatusSubDocBadRange)

	// ErrMemdSubDocBadDelta occurs when a sub-document counter operation is performed
	// and the specified delta is not valid.
	ErrMemdSubDocBadDelta = makeKvStatusError(memd.StatusSubDocBadDelta)

	// ErrMemdSubDocPathExists occurs when a sub-document operation expects a path not
	// to exists, but the path was found in the document.
	ErrMemdSubDocPathExists = makeKvStatusError(memd.StatusSubDocPathExists)

	// ErrMemdSubDocValueTooDeep occurs when a sub-document operation specifies a value
	// which is deeper than the depth limits of the sub-document specification.
	ErrMemdSubDocValueTooDeep = makeKvStatusError(memd.StatusSubDocValueTooDeep)

	// ErrMemdSubDocBadCombo occurs when a multi-operation sub-document operation is
	// performed and operations within the package of ops conflict with each other.
	ErrMemdSubDocBadCombo = makeKvStatusError(memd.StatusSubDocBadCombo)

	// ErrMemdSubDocBadMulti occurs when a multi-operation sub-document operation is
	// performed and operations within the package of ops conflict with each other.
	ErrMemdSubDocBadMulti = makeKvStatusError(memd.StatusSubDocBadMulti)

	// ErrMemdSubDocSuccessDeleted occurs when a multi-operation sub-document operation
	// is performed on a soft-deleted document.
	ErrMemdSubDocSuccessDeleted = makeKvStatusError(memd.StatusSubDocSuccessDeleted)

	// ErrMemdSubDocXattrInvalidFlagCombo occurs when an invalid set of
	// extended-attribute flags is passed to a sub-document operation.
	ErrMemdSubDocXattrInvalidFlagCombo = makeKvStatusError(memd.StatusSubDocXattrInvalidFlagCombo)

	// ErrMemdSubDocXattrInvalidKeyCombo occurs when an invalid set of key operations
	// are specified for a extended-attribute sub-document operation.
	ErrMemdSubDocXattrInvalidKeyCombo = makeKvStatusError(memd.StatusSubDocXattrInvalidKeyCombo)

	// ErrMemdSubDocXattrUnknownMacro occurs when an invalid macro value is specified.
	ErrMemdSubDocXattrUnknownMacro = makeKvStatusError(memd.StatusSubDocXattrUnknownMacro)

	// ErrMemdSubDocXattrUnknownVAttr occurs when an invalid virtual attribute is specified.
	ErrMemdSubDocXattrUnknownVAttr = makeKvStatusError(memd.StatusSubDocXattrUnknownVAttr)

	// ErrMemdSubDocXattrCannotModifyVAttr occurs when a mutation is attempted upon
	// a virtual attribute (which are immutable by definition).
	ErrMemdSubDocXattrCannotModifyVAttr = makeKvStatusError(memd.StatusSubDocXattrCannotModifyVAttr)

	// ErrMemdSubDocMultiPathFailureDeleted occurs when a Multi Path Failure occurs on
	// a soft-deleted document.
	ErrMemdSubDocMultiPathFailureDeleted = makeKvStatusError(memd.StatusSubDocMultiPathFailureDeleted)
)
