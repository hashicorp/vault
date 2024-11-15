package abstractions

// Represents the HTTP method used by a request.
type HttpMethod int

const (
	// The HTTP GET method.
	GET HttpMethod = iota
	// The HTTP POST method.
	POST
	// The HTTP PATCH method.
	PATCH
	// The HTTP DELETE method.
	DELETE
	// The HTTP OPTIONS method.
	OPTIONS
	// The HTTP CONNECT method.
	CONNECT
	// The HTTP PUT method.
	PUT
	// The HTTP TRACE method.
	TRACE
	// The HTTP HEAD method.
	HEAD
)

// String returns the string representation of the HTTP method.
func (m HttpMethod) String() string {
	return []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS", "CONNECT", "PUT", "TRACE", "HEAD"}[m]
}
