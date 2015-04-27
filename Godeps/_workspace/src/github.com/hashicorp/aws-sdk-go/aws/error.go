package aws

// An APIError is an error returned by an AWS API.
type APIError struct {
	StatusCode int // HTTP status code e.g. 200
	Type       string
	Code       string
	Message    string
	RequestID  string
	HostID     string
	Specifics  map[string]string
}

func (e APIError) Error() string {
	return e.Message
}
