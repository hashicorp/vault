package logical

// Response is a struct that stores the response of a request.
// It is used to abstract the details of the higher level request protocol.
type Response struct {
	// Secret, if not nil, denotes that this response represents a secret.
	Secret *Secret

	// Auth, if not nil, contains the authentication information for
	// this response. This is only checked and means something for
	// credential backends.
	Auth *Auth

	// Response data is an opaque map that must have string keys. For
	// secrets, this data is sent down to the user as-is. To store internal
	// data that you don't want the user to see, store it in
	// Secret.InternalData.
	Data map[string]interface{}

	// Redirect is an HTTP URL to redirect to for further authentication.
	// This is only valid for credential backends. This will be blanked
	// for any logical backend and ignored.
	Redirect string
}

// IsError returns true if this response seems to indicate an error.
func (r *Response) IsError() bool {
	return r != nil && len(r.Data) == 1 && r.Data["error"] != nil
}

// HelpResponse is used to format a help response
func HelpResponse(text string, seeAlso []string) *Response {
	return &Response{
		Data: map[string]interface{}{
			"help":     text,
			"see_also": seeAlso,
		},
	}
}

// ErrorResponse is used to format an error response
func ErrorResponse(text string) *Response {
	return &Response{
		Data: map[string]interface{}{
			"error": text,
		},
	}
}

// ListResponse is used to format a response to a list operation.
func ListResponse(keys []string) *Response {
	return &Response{
		Data: map[string]interface{}{
			"keys": keys,
		},
	}
}
