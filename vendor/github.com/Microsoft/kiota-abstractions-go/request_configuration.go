package abstractions

// RequestConfiguration represents a set of options to be used when making HTTP requests.
type RequestConfiguration[T any] struct {
	// Request headers
	Headers *RequestHeaders
	// Request options
	Options []RequestOption
	// Query parameters
	QueryParameters *T
}

// DefaultQueryParameters is a placeholder for operations without any query parameter documented.
type DefaultQueryParameters struct {
}
