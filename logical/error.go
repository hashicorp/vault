package logical

type HTTPCodedError interface {
	Error() string
	Code() int
}

func CodedError(status int, msg string) HTTPCodedError {
	return &codedError{
		Status:  status,
		Message: msg,
	}
}

type codedError struct {
	Status  int
	Message string
}

func (e *codedError) Error() string {
	return e.Message
}

func (e *codedError) Code() int {
	return e.Status
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
