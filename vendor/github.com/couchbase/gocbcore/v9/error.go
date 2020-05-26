package gocbcore

import (
	"errors"
	"io"
)

// dwError is a special error used for the purposes of rewrapping
// another error to provide more detailed information inherently
// with the error type itself.  Mainly used for timeout.
type dwError struct {
	InnerError error
	Message    string
}

func (e dwError) Error() string {
	return e.Message
}

func (e dwError) Unwrap() error {
	return e.InnerError
}

var (
	// ErrNoSupportedMechanisms occurs when the server does not support any of the
	// authentication methods that the client finds suitable.
	ErrNoSupportedMechanisms = errors.New("no supported authentication mechanisms")

	// ErrBadHosts occurs when the list of hosts specified cannot be contacted.
	ErrBadHosts = errors.New("failed to connect to any of the specified hosts")

	// ErrProtocol occurs when the server responds with unexpected or unparseable data.
	ErrProtocol = errors.New("failed to parse server response")

	// ErrNoReplicas occurs when no replicas respond in time
	ErrNoReplicas = errors.New("no replicas responded in time")

	// ErrCliInternalError indicates an internal error occurred within the client.
	ErrCliInternalError = errors.New("client internal error")

	// ErrInvalidCredentials is returned when an invalid set of credentials is provided for a service.
	ErrInvalidCredentials = errors.New("an invalid set of credentials was provided")

	// ErrInvalidServer occurs when an explicit, but invalid server index is specified.
	ErrInvalidServer = errors.New("specific server index is invalid")

	// ErrInvalidVBucket occurs when an explicit, but invalid vbucket index is specified.
	ErrInvalidVBucket = errors.New("specific vbucket index is invalid")

	// ErrInvalidReplica occurs when an explicit, but invalid replica index is specified.
	ErrInvalidReplica = errors.New("specific server index is invalid")

	// ErrInvalidService occurs when an explicit but invalid service type is specified
	ErrInvalidService = errors.New("invalid service")

	// ErrInvalidCertificate occurs when a certificate that is not useable is passed to an Agent.
	ErrInvalidCertificate = errors.New("certificate is invalid")

	// ErrCollectionsUnsupported occurs when collections are used but either server does not support them or the agent
	// was created without them enabled.
	ErrCollectionsUnsupported = errors.New("collections are not enabled")

	// ErrBucketAlreadySelected occurs when SelectBucket is called when a bucket is already selected..
	ErrBucketAlreadySelected = errors.New("bucket already selected")

	// ErrShutdown occurs when operations are performed on a previously closed Agent.
	ErrShutdown = errors.New("connection shut down")

	// ErrOverload occurs when too many operations are dispatched and all queues are full.
	ErrOverload = errors.New("queue overflowed")

	// ErrSocketClosed occurs when a socket closes while an operation is in flight.
	ErrSocketClosed = io.EOF

	// ErrGCCCPInUse occurs when an operation dis performed whilst the client is connect via GCCCP.
	ErrGCCCPInUse = errors.New("connected via gcccp, kv operations are not supported, open a bucket first")

	// ErrNotMyVBucket occurs when an operation is sent to a node which does not own the vbucket.
	ErrNotMyVBucket = errors.New("not my vbucket")
)

// Shared Error Definitions RFC#58@15
var (
	// ErrTimeout occurs when an operation does not receive a response in a timely manner.
	ErrTimeout = errors.New("operation has timed out")

	ErrRequestCanceled       = errors.New("request canceled")
	ErrInvalidArgument       = errors.New("invalid argument")
	ErrServiceNotAvailable   = errors.New("service not available")
	ErrInternalServerFailure = errors.New("internal server failure")
	ErrAuthenticationFailure = errors.New("authentication failure")
	ErrTemporaryFailure      = errors.New("temporary failure")
	ErrParsingFailure        = errors.New("parsing failure")

	ErrCasMismatch          = errors.New("cas mismatch")
	ErrBucketNotFound       = errors.New("bucket not found")
	ErrCollectionNotFound   = errors.New("collection not found")
	ErrEncodingFailure      = errors.New("encoding failure")
	ErrDecodingFailure      = errors.New("decoding failure")
	ErrUnsupportedOperation = errors.New("unsupported operation")
	ErrAmbiguousTimeout     = &dwError{ErrTimeout, "ambiguous timeout"}

	ErrUnambiguousTimeout = &dwError{ErrTimeout, "unambiguous timeout"}

	// ErrFeatureNotAvailable occurs when an operation is performed on a bucket which does not support it.
	ErrFeatureNotAvailable = errors.New("feature is not available")
	ErrScopeNotFound       = errors.New("scope not found")
	ErrIndexNotFound       = errors.New("index not found")

	ErrIndexExists = errors.New("index exists")
)

// Key Value Error Definitions RFC#58@15
var (
	ErrDocumentNotFound                  = errors.New("document not found")
	ErrDocumentUnretrievable             = errors.New("document unretrievable")
	ErrDocumentLocked                    = errors.New("document locked")
	ErrValueTooLarge                     = errors.New("value too large")
	ErrDocumentExists                    = errors.New("document exists")
	ErrValueNotJSON                      = errors.New("value not json")
	ErrDurabilityLevelNotAvailable       = errors.New("durability level not available")
	ErrDurabilityImpossible              = errors.New("durability impossible")
	ErrDurabilityAmbiguous               = errors.New("durability ambiguous")
	ErrDurableWriteInProgress            = errors.New("durable write in progress")
	ErrDurableWriteReCommitInProgress    = errors.New("durable write recommit in progress")
	ErrMutationLost                      = errors.New("mutation lost")
	ErrPathNotFound                      = errors.New("path not found")
	ErrPathMismatch                      = errors.New("path mismatch")
	ErrPathInvalid                       = errors.New("path invalid")
	ErrPathTooBig                        = errors.New("path too big")
	ErrPathTooDeep                       = errors.New("path too deep")
	ErrValueTooDeep                      = errors.New("value too deep")
	ErrValueInvalid                      = errors.New("value invalid")
	ErrDocumentNotJSON                   = errors.New("document not json")
	ErrNumberTooBig                      = errors.New("number too big")
	ErrDeltaInvalid                      = errors.New("delta invalid")
	ErrPathExists                        = errors.New("path exists")
	ErrXattrUnknownMacro                 = errors.New("xattr unknown macro")
	ErrXattrInvalidFlagCombo             = errors.New("xattr invalid flag combination")
	ErrXattrInvalidKeyCombo              = errors.New("xattr invalid key combination")
	ErrXattrUnknownVirtualAttribute      = errors.New("xattr unknown virtual attribute")
	ErrXattrCannotModifyVirtualAttribute = errors.New("xattr cannot modify virtual attribute")
	ErrXattrInvalidOrder                 = errors.New("xattr invalid order")
)

// Query Error Definitions RFC#58@15
var (
	ErrPlanningFailure = errors.New("planning failure")

	ErrIndexFailure = errors.New("index failure")

	ErrPreparedStatementFailure = errors.New("prepared statement failure")
)

// Analytics Error Definitions RFC#58@15
var (
	ErrCompilationFailure = errors.New("compilation failure")

	ErrJobQueueFull = errors.New("job queue full")

	ErrDatasetNotFound = errors.New("dataset not found")

	ErrDataverseNotFound = errors.New("dataverse not found")

	ErrDatasetExists = errors.New("dataset exists")

	ErrDataverseExists = errors.New("dataverse exists")

	ErrLinkNotFound = errors.New("link not found")
)

// Search Error Definitions RFC#58@15
var ()

// View Error Definitions RFC#58@15
var (
	ErrViewNotFound = errors.New("view not found")

	ErrDesignDocumentNotFound = errors.New("design document not found")
)

// Management Error Definitions RFC#58@15
var (
	ErrCollectionExists   = errors.New("collection exists")
	ErrScopeExists        = errors.New("scope exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrGroupNotFound      = errors.New("group not found")
	ErrBucketExists       = errors.New("bucket exists")
	ErrUserExists         = errors.New("user exists")
	ErrBucketNotFlushable = errors.New("bucket not flushable")
)
