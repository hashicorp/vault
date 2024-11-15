package abstractions

import "fmt"

type ApiErrorable interface {
	SetResponseHeaders(ResponseHeaders *ResponseHeaders)
	SetStatusCode(ResponseStatusCode int)
	GetResponseHeaders() *ResponseHeaders
	GetStatusCode() int
}

// ApiError is the parent type for errors thrown by the client when receiving failed responses to its requests
type ApiError struct {
	Message            string
	ResponseStatusCode int
	ResponseHeaders    *ResponseHeaders
}

func (e *ApiError) Error() string {
	if len(e.Message) > 0 {
		return fmt.Sprint(e.Message)
	} else {
		return "error status code received from the API"
	}
}

// NewApiError creates a new ApiError instance
func NewApiError() *ApiError {
	return &ApiError{ResponseHeaders: NewResponseHeaders()}
}

func (e *ApiError) SetResponseHeaders(ResponseHeaders *ResponseHeaders) {
	e.ResponseHeaders = ResponseHeaders
}

func (e *ApiError) SetStatusCode(ResponseStatusCode int) {
	e.ResponseStatusCode = ResponseStatusCode
}

func (e *ApiError) GetResponseHeaders() *ResponseHeaders {
	return e.ResponseHeaders
}

func (e *ApiError) GetStatusCode() int {
	return e.ResponseStatusCode
}
