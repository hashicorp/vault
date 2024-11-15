package abstractions

// Represents a request option.
type RequestOption interface {
	// GetKey returns the key to store the current option under.
	GetKey() RequestOptionKey
}

// RequestOptionKey represents a key to store a request option under.
type RequestOptionKey struct {
	// The unique key for the option.
	Key string
}
