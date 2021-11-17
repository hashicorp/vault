package memd

import "fmt"

// StatusCode represents a memcached response status.
type StatusCode uint16

const (
	// StatusSuccess indicates the operation completed successfully.
	StatusSuccess = StatusCode(0x00)

	// StatusKeyNotFound occurs when an operation is performed on a key that does not exist.
	StatusKeyNotFound = StatusCode(0x01)

	// StatusKeyExists occurs when an operation is performed on a key that could not be found.
	StatusKeyExists = StatusCode(0x02)

	// StatusTooBig occurs when an operation attempts to store more data in a single document
	// than the server is capable of storing (by default, this is a 20MB limit).
	StatusTooBig = StatusCode(0x03)

	// StatusInvalidArgs occurs when the server receives invalid arguments for an operation.
	StatusInvalidArgs = StatusCode(0x04)

	// StatusNotStored occurs when the server fails to store a key.
	StatusNotStored = StatusCode(0x05)

	// StatusBadDelta occurs when an invalid delta value is specified to a counter operation.
	StatusBadDelta = StatusCode(0x06)

	// StatusNotMyVBucket occurs when an operation is dispatched to a server which is
	// non-authoritative for a specific vbucket.
	StatusNotMyVBucket = StatusCode(0x07)

	// StatusNoBucket occurs when no bucket was selected on a connection.
	StatusNoBucket = StatusCode(0x08)

	// StatusLocked occurs when an operation fails due to the document being locked.
	StatusLocked = StatusCode(0x09)

	// StatusAuthStale occurs when authentication credentials have become invalidated.
	StatusAuthStale = StatusCode(0x1f)

	// StatusAuthError occurs when the authentication information provided was not valid.
	StatusAuthError = StatusCode(0x20)

	// StatusAuthContinue occurs in multi-step authentication when more authentication
	// work needs to be performed in order to complete the authentication process.
	StatusAuthContinue = StatusCode(0x21)

	// StatusRangeError occurs when the range specified to the server is not valid.
	StatusRangeError = StatusCode(0x22)

	// StatusRollback occurs when a DCP stream fails to open due to a rollback having
	// previously occurred since the last time the stream was opened.
	StatusRollback = StatusCode(0x23)

	// StatusAccessError occurs when an access error occurs.
	StatusAccessError = StatusCode(0x24)

	// StatusNotInitialized is sent by servers which are still initializing, and are not
	// yet ready to accept operations on behalf of a particular bucket.
	StatusNotInitialized = StatusCode(0x25)

	// StatusUnknownCommand occurs when an unknown operation is sent to a server.
	StatusUnknownCommand = StatusCode(0x81)

	// StatusOutOfMemory occurs when the server cannot service a request due to memory
	// limitations.
	StatusOutOfMemory = StatusCode(0x82)

	// StatusNotSupported occurs when an operation is understood by the server, but that
	// operation is not supported on this server (occurs for a variety of reasons).
	StatusNotSupported = StatusCode(0x83)

	// StatusInternalError occurs when internal errors prevent the server from processing
	// your request.
	StatusInternalError = StatusCode(0x84)

	// StatusBusy occurs when the server is too busy to process your request right away.
	// Attempting the operation at a later time will likely succeed.
	StatusBusy = StatusCode(0x85)

	// StatusTmpFail occurs when a temporary failure is preventing the server from
	// processing your request.
	StatusTmpFail = StatusCode(0x86)

	// StatusCollectionUnknown occurs when a Collection cannot be found.
	StatusCollectionUnknown = StatusCode(0x88)

	// StatusScopeUnknown occurs when a Scope cannot be found.
	StatusScopeUnknown = StatusCode(0x8c)

	// StatusDurabilityInvalidLevel occurs when an invalid durability level was requested.
	StatusDurabilityInvalidLevel = StatusCode(0xa0)

	// StatusDurabilityImpossible occurs when a request is performed with impossible
	// durability level requirements.
	StatusDurabilityImpossible = StatusCode(0xa1)

	// StatusSyncWriteInProgress occurs when an attempt is made to write to a key that has
	// a SyncWrite pending.
	StatusSyncWriteInProgress = StatusCode(0xa2)

	// StatusSyncWriteAmbiguous occurs when an SyncWrite does not complete in the specified
	// time and the result is ambiguous.
	StatusSyncWriteAmbiguous = StatusCode(0xa3)

	// StatusSyncWriteReCommitInProgress occurs when an SyncWrite is being recommitted.
	StatusSyncWriteReCommitInProgress = StatusCode(0xa4)

	// StatusSubDocPathNotFound occurs when a sub-document operation targets a path
	// which does not exist in the specifie document.
	StatusSubDocPathNotFound = StatusCode(0xc0)

	// StatusSubDocPathMismatch occurs when a sub-document operation specifies a path
	// which does not match the document structure (field access on an array).
	StatusSubDocPathMismatch = StatusCode(0xc1)

	// StatusSubDocPathInvalid occurs when a sub-document path could not be parsed.
	StatusSubDocPathInvalid = StatusCode(0xc2)

	// StatusSubDocPathTooBig occurs when a sub-document path is too big.
	StatusSubDocPathTooBig = StatusCode(0xc3)

	// StatusSubDocDocTooDeep occurs when an operation would cause a document to be
	// nested beyond the depth limits allowed by the sub-document specification.
	StatusSubDocDocTooDeep = StatusCode(0xc4)

	// StatusSubDocCantInsert occurs when a sub-document operation could not insert.
	StatusSubDocCantInsert = StatusCode(0xc5)

	// StatusSubDocNotJSON occurs when a sub-document operation is performed on a
	// document which is not JSON.
	StatusSubDocNotJSON = StatusCode(0xc6)

	// StatusSubDocBadRange occurs when a sub-document operation is performed with
	// a bad range.
	StatusSubDocBadRange = StatusCode(0xc7)

	// StatusSubDocBadDelta occurs when a sub-document counter operation is performed
	// and the specified delta is not valid.
	StatusSubDocBadDelta = StatusCode(0xc8)

	// StatusSubDocPathExists occurs when a sub-document operation expects a path not
	// to exists, but the path was found in the document.
	StatusSubDocPathExists = StatusCode(0xc9)

	// StatusSubDocValueTooDeep occurs when a sub-document operation specifies a value
	// which is deeper than the depth limits of the sub-document specification.
	StatusSubDocValueTooDeep = StatusCode(0xca)

	// StatusSubDocBadCombo occurs when a multi-operation sub-document operation is
	// performed and operations within the package of ops conflict with each other.
	StatusSubDocBadCombo = StatusCode(0xcb)

	// StatusSubDocBadMulti occurs when a multi-operation sub-document operation is
	// performed and operations within the package of ops conflict with each other.
	StatusSubDocBadMulti = StatusCode(0xcc)

	// StatusSubDocSuccessDeleted occurs when a multi-operation sub-document operation
	// is performed on a soft-deleted document.
	StatusSubDocSuccessDeleted = StatusCode(0xcd)

	// StatusSubDocXattrInvalidFlagCombo occurs when an invalid set of
	// extended-attribute flags is passed to a sub-document operation.
	StatusSubDocXattrInvalidFlagCombo = StatusCode(0xce)

	// StatusSubDocXattrInvalidKeyCombo occurs when an invalid set of key operations
	// are specified for a extended-attribute sub-document operation.
	StatusSubDocXattrInvalidKeyCombo = StatusCode(0xcf)

	// StatusSubDocXattrUnknownMacro occurs when an invalid macro value is specified.
	StatusSubDocXattrUnknownMacro = StatusCode(0xd0)

	// StatusSubDocXattrUnknownVAttr occurs when an invalid virtual attribute is specified.
	StatusSubDocXattrUnknownVAttr = StatusCode(0xd1)

	// StatusSubDocXattrCannotModifyVAttr occurs when a mutation is attempted upon
	// a virtual attribute (which are immutable by definition).
	StatusSubDocXattrCannotModifyVAttr = StatusCode(0xd2)

	// StatusSubDocMultiPathFailureDeleted occurs when a Multi Path Failure occurs on
	// a soft-deleted document.
	StatusSubDocMultiPathFailureDeleted = StatusCode(0xd3)
)

// KVText returns the textual representation of this StatusCode.
func (code StatusCode) KVText() string {
	switch code {
	case StatusSuccess:
		return "success"
	case StatusKeyNotFound:
		return "key not found"
	case StatusKeyExists:
		return "key already exists, if a cas was provided the key exists with a different cas"
	case StatusTooBig:
		return "document value was too large"
	case StatusInvalidArgs:
		return "invalid arguments"
	case StatusNotStored:
		return "document could not be stored"
	case StatusBadDelta:
		return "invalid delta was passed"
	case StatusNotMyVBucket:
		return "operation sent to incorrect server"
	case StatusNoBucket:
		return "not connected to a bucket"
	case StatusAuthStale:
		return "authentication context is stale, try re-authenticating"
	case StatusAuthError:
		return "authentication error"
	case StatusAuthContinue:
		return "more authentication steps needed"
	case StatusRangeError:
		return "requested value is outside range"
	case StatusAccessError:
		return "no access"
	case StatusNotInitialized:
		return "cluster is being initialized, requests are blocked"
	case StatusRollback:
		return "rollback is required"
	case StatusUnknownCommand:
		return "unknown command was received"
	case StatusOutOfMemory:
		return "server is out of memory"
	case StatusNotSupported:
		return "server does not support this command"
	case StatusInternalError:
		return "internal server error"
	case StatusBusy:
		return "server is busy, try again later"
	case StatusTmpFail:
		return "temporary failure occurred, try again later"
	case StatusCollectionUnknown:
		return "the requested collection cannot be found"
	case StatusScopeUnknown:
		return "the requested scope cannot be found."
	case StatusDurabilityInvalidLevel:
		return "invalid request, invalid durability level specified."
	case StatusDurabilityImpossible:
		return "the requested durability requirements are impossible."
	case StatusSyncWriteInProgress:
		return "key already has syncwrite pending."
	case StatusSyncWriteAmbiguous:
		return "the syncwrite request did not complete in time."
	case StatusSubDocPathNotFound:
		return "sub-document path does not exist"
	case StatusSubDocPathMismatch:
		return "type of element in sub-document path conflicts with type in document"
	case StatusSubDocPathInvalid:
		return "malformed sub-document path"
	case StatusSubDocPathTooBig:
		return "sub-document contains too many components"
	case StatusSubDocDocTooDeep:
		return "existing document contains too many levels of nesting"
	case StatusSubDocCantInsert:
		return "subdocument operation would invalidate the JSON"
	case StatusSubDocNotJSON:
		return "existing document is not valid JSON"
	case StatusSubDocBadRange:
		return "existing numeric value is too large"
	case StatusSubDocBadDelta:
		return "numeric operation would yield a number that is too large, or " +
			"a zero delta was specified"
	case StatusSubDocPathExists:
		return "given path already exists in the document"
	case StatusSubDocValueTooDeep:
		return "value is too deep to insert"
	case StatusSubDocBadCombo:
		return "incorrectly matched subdocument operation types"
	case StatusSubDocBadMulti:
		return "could not execute one or more multi lookups or mutations"
	case StatusSubDocSuccessDeleted:
		return "document is soft-deleted"
	case StatusSubDocXattrInvalidFlagCombo:
		return "invalid xattr flag combination"
	case StatusSubDocXattrInvalidKeyCombo:
		return "invalid xattr key combination"
	case StatusSubDocXattrUnknownMacro:
		return "unknown xattr macro"
	case StatusSubDocXattrUnknownVAttr:
		return "unknown xattr virtual attribute"
	case StatusSubDocXattrCannotModifyVAttr:
		return "cannot modify virtual attributes"
	case StatusSubDocMultiPathFailureDeleted:
		return "sub-document multi-path error"
	default:
		return fmt.Sprintf("unknown kv status code (%d)", code)
	}
}
