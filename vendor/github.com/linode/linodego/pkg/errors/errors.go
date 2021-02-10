package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	// FromString is the Code identifying Errors created by string types
	FromString = iota + 1
	// FromError is the Code identifying Errors created by error types
	FromError
	// FromStringer is the Code identifying Errors created by fmt.Stringer types
	FromStringer
)

// Error represents an Error in the context of an API method call.
type Error struct {
	Response *http.Response
	Code     int
	Message  string
}

func (e Error) Error() string {
	return fmt.Sprintf("[%03d] %s", e.Code, e.Message)
}

// APIErrorReason represents an individual invalid request message returned by the Linode API.
type APIErrorReason struct {
	Reason string `json:"reason"`
	Field  string `json:"field"`
}

func (r APIErrorReason) Error() string {
	if len(r.Field) == 0 {
		return r.Reason
	}

	return fmt.Sprintf("[%s] %s", r.Field, r.Reason)
}

// APIError is the set of errors returned by the Linode API on an invalid request.
type APIError struct {
	Errors []APIErrorReason `json:"errors"`
}

func (e APIError) Error() string {
	x := []string{}
	for _, msg := range e.Errors {
		x = append(x, msg.Error())
	}

	return strings.Join(x, "; ")
}

func CoupleAPIErrors(r *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return nil, New(err)
	}

	if r.Error() != nil {
		// Check that response is of the correct content-type before unmarshalling
		expectedContentType := r.Request.Header.Get("Accept")
		responseContentType := r.Header().Get("Content-Type")

		// If the upstream Linode API server being fronted fails to respond to the request,
		// the http server will respond with a default "Bad Gateway" page with Content-Type
		// "text/html".
		if r.StatusCode() == http.StatusBadGateway && responseContentType == "text/html" {
			return nil, Error{Code: http.StatusBadGateway, Message: http.StatusText(http.StatusBadGateway)}
		}

		if responseContentType != expectedContentType {
			msg := fmt.Sprintf(
				"Unexpected Content-Type: Expected: %v, Received: %v",
				expectedContentType,
				responseContentType,
			)
			return nil, New(msg)
		}

		apiError, ok := r.Error().(*APIError)
		if !ok || (ok && len(apiError.Errors) == 0) {
			return r, nil
		}

		return nil, New(r)
	}

	return r, nil
}

// New creates a Error with a Code identifying the source err type.
// - FromString   (1) from a string
// - FromError    (2) for an error
// - FromStringer (3) for a Stringer
// - HTTP Status Codes (100-600) for a resty.Response object
func New(err interface{}) *Error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *Error:
		return e
	case *resty.Response:
		apiError, ok := e.Error().(*APIError)
		if !ok {
			return nil
		}

		return &Error{
			Code:     e.RawResponse.StatusCode,
			Message:  apiError.Error(),
			Response: e.RawResponse,
		}
	case error:
		return &Error{Code: FromError, Message: e.Error()}
	case string:
		return &Error{Code: FromString, Message: e}
	case fmt.Stringer:
		return &Error{Code: FromStringer, Message: e.String()}
	default:
		return nil
	}
}
