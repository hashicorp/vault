package gocb

import (
	"errors"
	"fmt"

	gocbcore "github.com/couchbase/gocbcore/v9"
)

type wrappedError struct {
	Message    string
	InnerError error
}

func (e wrappedError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.InnerError.Error())
}

func (e wrappedError) Unwrap() error {
	return e.InnerError
}

func wrapError(err error, message string) error {
	return wrappedError{
		Message:    message,
		InnerError: err,
	}
}

type invalidArgumentsError struct {
	message string
}

func (e invalidArgumentsError) Error() string {
	return fmt.Sprintf("invalid arguments: %s", e.message)
}

func (e invalidArgumentsError) Unwrap() error {
	return ErrInvalidArgument
}

func makeInvalidArgumentsError(message string) error {
	return invalidArgumentsError{
		message: message,
	}
}

// Shared Error Definitions RFC#58@15
var (
	// ErrTimeout occurs when an operation does not receive a response in a timely manner.
	ErrTimeout = gocbcore.ErrTimeout

	// ErrRequestCanceled occurs when an operation has been canceled.
	ErrRequestCanceled = gocbcore.ErrRequestCanceled

	// ErrInvalidArgument occurs when an invalid argument is provided for an operation.
	ErrInvalidArgument = gocbcore.ErrInvalidArgument

	// ErrServiceNotAvailable occurs when the requested service is not available.
	ErrServiceNotAvailable = gocbcore.ErrServiceNotAvailable

	// ErrInternalServerFailure occurs when the server encounters an internal server error.
	ErrInternalServerFailure = gocbcore.ErrInternalServerFailure

	// ErrAuthenticationFailure occurs when authentication has failed.
	ErrAuthenticationFailure = gocbcore.ErrAuthenticationFailure

	// ErrTemporaryFailure occurs when an operation has failed for a reason that is temporary.
	ErrTemporaryFailure = gocbcore.ErrTemporaryFailure

	// ErrParsingFailure occurs when a query has failed to be parsed by the server.
	ErrParsingFailure = gocbcore.ErrParsingFailure

	// ErrCasMismatch occurs when an operation has been performed with a cas value that does not the value on the server.
	ErrCasMismatch = gocbcore.ErrCasMismatch

	// ErrBucketNotFound occurs when the requested bucket could not be found.
	ErrBucketNotFound = gocbcore.ErrBucketNotFound

	// ErrCollectionNotFound occurs when the requested collection could not be found.
	ErrCollectionNotFound = gocbcore.ErrCollectionNotFound

	// ErrEncodingFailure occurs when encoding of a value failed.
	ErrEncodingFailure = gocbcore.ErrEncodingFailure

	// ErrDecodingFailure occurs when decoding of a value failed.
	ErrDecodingFailure = gocbcore.ErrDecodingFailure

	// ErrUnsupportedOperation occurs when an operation that is unsupported or unknown is performed against the server.
	ErrUnsupportedOperation = gocbcore.ErrUnsupportedOperation

	// ErrAmbiguousTimeout occurs when an operation does not receive a response in a timely manner for a reason that
	//
	ErrAmbiguousTimeout = gocbcore.ErrAmbiguousTimeout

	// ErrAmbiguousTimeout occurs when an operation does not receive a response in a timely manner for a reason that
	// it can be safely established that
	ErrUnambiguousTimeout = gocbcore.ErrUnambiguousTimeout

	// ErrFeatureNotAvailable occurs when an operation is performed on a bucket which does not support it.
	ErrFeatureNotAvailable = gocbcore.ErrFeatureNotAvailable

	// ErrScopeNotFound occurs when the requested scope could not be found.
	ErrScopeNotFound = gocbcore.ErrScopeNotFound

	// ErrIndexNotFound occurs when the requested index could not be found.
	ErrIndexNotFound = gocbcore.ErrIndexNotFound

	// ErrIndexExists occurs when creating an index that already exists.
	ErrIndexExists = gocbcore.ErrIndexExists
)

// Key Value Error Definitions RFC#58@15
var (
	// ErrDocumentNotFound occurs when the requested document could not be found.
	ErrDocumentNotFound = gocbcore.ErrDocumentNotFound

	// ErrDocumentUnretrievable occurs when GetAnyReplica cannot find the document on any replica.
	ErrDocumentUnretrievable = gocbcore.ErrDocumentUnretrievable

	// ErrDocumentLocked occurs when a mutation operation is attempted against a document that is locked.
	ErrDocumentLocked = gocbcore.ErrDocumentLocked

	// ErrValueTooLarge occurs when a document has gone over the maximum size allowed by the server.
	ErrValueTooLarge = gocbcore.ErrValueTooLarge

	// ErrDocumentExists occurs when an attempt is made to insert a document but a document with that key already exists.
	ErrDocumentExists = gocbcore.ErrDocumentExists

	// ErrValueNotJSON occurs when a sub-document operation is performed on a
	// document which is not JSON.
	ErrValueNotJSON = gocbcore.ErrValueNotJSON

	// ErrDurabilityLevelNotAvailable occurs when an invalid durability level was requested.
	ErrDurabilityLevelNotAvailable = gocbcore.ErrDurabilityLevelNotAvailable

	// ErrDurabilityImpossible occurs when a request is performed with impossible
	// durability level requirements.
	ErrDurabilityImpossible = gocbcore.ErrDurabilityImpossible

	// ErrDurabilityAmbiguous occurs when an SyncWrite does not complete in the specified
	// time and the result is ambiguous.
	ErrDurabilityAmbiguous = gocbcore.ErrDurabilityAmbiguous

	// ErrDurableWriteInProgress occurs when an attempt is made to write to a key that has
	// a SyncWrite pending.
	ErrDurableWriteInProgress = gocbcore.ErrDurableWriteInProgress

	// ErrDurableWriteReCommitInProgress occurs when an SyncWrite is being recommitted.
	ErrDurableWriteReCommitInProgress = gocbcore.ErrDurableWriteReCommitInProgress

	// ErrMutationLost occurs when a mutation was lost.
	ErrMutationLost = gocbcore.ErrMutationLost

	// ErrPathNotFound occurs when a sub-document operation targets a path
	// which does not exist in the specified document.
	ErrPathNotFound = gocbcore.ErrPathNotFound

	// ErrPathMismatch occurs when a sub-document operation specifies a path
	// which does not match the document structure (field access on an array).
	ErrPathMismatch = gocbcore.ErrPathMismatch

	// ErrPathInvalid occurs when a sub-document path could not be parsed.
	ErrPathInvalid = gocbcore.ErrPathInvalid

	// ErrPathTooBig occurs when a sub-document path is too big.
	ErrPathTooBig = gocbcore.ErrPathTooBig

	// ErrPathTooDeep occurs when an operation would cause a document to be
	// nested beyond the depth limits allowed by the sub-document specification.
	ErrPathTooDeep = gocbcore.ErrPathTooDeep

	// ErrValueTooDeep occurs when a sub-document operation specifies a value
	// which is deeper than the depth limits of the sub-document specification.
	ErrValueTooDeep = gocbcore.ErrValueTooDeep

	// ErrValueInvalid occurs when a sub-document operation could not insert.
	ErrValueInvalid = gocbcore.ErrValueInvalid

	// ErrDocumentNotJSON occurs when a sub-document operation is performed on a
	// document which is not JSON.
	ErrDocumentNotJSON = gocbcore.ErrDocumentNotJSON

	// ErrNumberTooBig occurs when a sub-document operation is performed with
	// a bad range.
	ErrNumberTooBig = gocbcore.ErrNumberTooBig

	// ErrDeltaInvalid occurs when a sub-document counter operation is performed
	// and the specified delta is not valid.
	ErrDeltaInvalid = gocbcore.ErrDeltaInvalid

	// ErrPathExists occurs when a sub-document operation expects a path not
	// to exists, but the path was found in the document.
	ErrPathExists = gocbcore.ErrPathExists

	// ErrXattrUnknownMacro occurs when an invalid macro value is specified.
	ErrXattrUnknownMacro = gocbcore.ErrXattrUnknownMacro

	// ErrXattrInvalidFlagCombo occurs when an invalid set of
	// extended-attribute flags is passed to a sub-document operation.
	ErrXattrInvalidFlagCombo = gocbcore.ErrXattrInvalidFlagCombo

	// ErrXattrInvalidKeyCombo occurs when an invalid set of key operations
	// are specified for a extended-attribute sub-document operation.
	ErrXattrInvalidKeyCombo = gocbcore.ErrXattrInvalidKeyCombo

	// ErrXattrUnknownVirtualAttribute occurs when an invalid virtual attribute is specified.
	ErrXattrUnknownVirtualAttribute = gocbcore.ErrXattrUnknownVirtualAttribute

	// ErrXattrCannotModifyVirtualAttribute occurs when a mutation is attempted upon
	// a virtual attribute (which are immutable by definition).
	ErrXattrCannotModifyVirtualAttribute = gocbcore.ErrXattrCannotModifyVirtualAttribute

	// ErrXattrInvalidOrder occurs when a set key key operations are specified for a extended-attribute sub-document
	// operation in the incorrect order.
	ErrXattrInvalidOrder = gocbcore.ErrXattrInvalidOrder
)

// Query Error Definitions RFC#58@15
var (
	// ErrPlanningFailure occurs when the query service was unable to create a query plan.
	ErrPlanningFailure = gocbcore.ErrPlanningFailure

	// ErrIndexFailure occurs when there was an issue with the index specified.
	ErrIndexFailure = gocbcore.ErrIndexFailure

	// ErrPreparedStatementFailure occurs when there was an issue with the prepared statement.
	ErrPreparedStatementFailure = gocbcore.ErrPreparedStatementFailure
)

// Analytics Error Definitions RFC#58@15
var (
	// ErrCompilationFailure occurs when there was an issue executing the analytics query because it could not
	// be compiled.
	ErrCompilationFailure = gocbcore.ErrCompilationFailure

	// ErrJobQueueFull occurs when the analytics service job queue is full.
	ErrJobQueueFull = gocbcore.ErrJobQueueFull

	// ErrDatasetNotFound occurs when the analytics dataset requested could not be found.
	ErrDatasetNotFound = gocbcore.ErrDatasetNotFound

	// ErrDataverseNotFound occurs when the analytics dataverse requested could not be found.
	ErrDataverseNotFound = gocbcore.ErrDataverseNotFound

	// ErrDatasetExists occurs when creating an analytics dataset failed because it already exists.
	ErrDatasetExists = gocbcore.ErrDatasetExists

	// ErrDataverseExists occurs when creating an analytics dataverse failed because it already exists.
	ErrDataverseExists = gocbcore.ErrDataverseExists

	// ErrLinkNotFound occurs when the analytics link requested could not be found.
	ErrLinkNotFound = gocbcore.ErrLinkNotFound
)

// Search Error Definitions RFC#58@15
var ()

// View Error Definitions RFC#58@15
var (
	// ErrViewNotFound occurs when the view requested could not be found.
	ErrViewNotFound = gocbcore.ErrViewNotFound

	// ErrDesignDocumentNotFound occurs when the design document requested could not be found.
	ErrDesignDocumentNotFound = gocbcore.ErrDesignDocumentNotFound
)

// Management Error Definitions RFC#58@15
var (
	// ErrCollectionExists occurs when creating a collection failed because it already exists.
	ErrCollectionExists = gocbcore.ErrCollectionExists

	// ErrScopeExists occurs when creating a scope failed because it already exists.
	ErrScopeExists = gocbcore.ErrScopeExists

	// ErrUserNotFound occurs when the user requested could not be found.
	ErrUserNotFound = gocbcore.ErrUserNotFound

	// ErrGroupNotFound occurs when the group requested could not be found.
	ErrGroupNotFound = gocbcore.ErrGroupNotFound

	// ErrBucketExists occurs when creating a bucket failed because it already exists.
	ErrBucketExists = gocbcore.ErrBucketExists

	// ErrUserExists occurs when creating a user failed because it already exists.
	ErrUserExists = gocbcore.ErrUserExists

	// ErrBucketNotFlushable occurs when a bucket could not be flushed because flushing is not enabled.
	ErrBucketNotFlushable = gocbcore.ErrBucketNotFlushable
)

// SDK specific error definitions
var (
	// ErrOverload occurs when too many operations are dispatched and all queues are full.
	ErrOverload = gocbcore.ErrOverload

	// ErrNoResult occurs when no results are available to a query.
	ErrNoResult = errors.New("no result was available")
)
