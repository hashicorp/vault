package logical

// CodedError is an error which includes an HTTP status code. Various places in
// Vault will "unwrap" this error to return the correct response to the user.
type CodedError struct {
	Status  int
	Message string
}

// Error implements the standard error interface and returns the given message.
func (e *CodedError) Error() string {
	return e.Message
}

// Code returns the HTTP status code of this error. This exists for backwards
// compatability. New implementations should use .Status directly.
func (e *CodedError) Code() int {
	return e.Status
}

// NewCodedError creates a new HTTP CodedError from the given status and
// message.
func NewCodedError(status int, message string) *CodedError {
	return &CodedError{
		Status:  status,
		Message: message,
	}
}

// Struct to identify user input errors.  This is helpful in responding the
// appropriate status codes to clients from the HTTP endpoints.
type StatusBadRequest struct {
	Err string
}

// Implementing error interface
func (s *StatusBadRequest) Error() string {
	return s.Err
}

// This is a new type declared to not cause potential compatibility problems if
// the logic around the CodedError changes; in particular for logical request
// paths it is basically ignored, and changing that behavior might cause
// unforseen issues.
type ReplicationCodedError struct {
	Msg  string
	Code int
}

func (r *ReplicationCodedError) Error() string {
	return r.Msg
}
