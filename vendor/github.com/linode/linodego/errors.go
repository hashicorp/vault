package linodego

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	ErrorUnsupported = iota
	// ErrorFromString is the Code identifying Errors created by string types
	ErrorFromString
	// ErrorFromError is the Code identifying Errors created by error types
	ErrorFromError
	// ErrorFromStringer is the Code identifying Errors created by fmt.Stringer types
	ErrorFromStringer
)

// Error wraps the LinodeGo error with the relevant http.Response
type Error struct {
	Response *http.Response
	Code     int
	Message  string
}

// APIErrorReason is an individual invalid request message returned by the Linode API
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

// APIError is the error-set returned by the Linode API when presented with an invalid request
type APIError struct {
	Errors []APIErrorReason `json:"errors"`
}

// String returns the error reason in a formatted string
func (r APIErrorReason) String() string {
	return fmt.Sprintf("[%s] %s", r.Field, r.Reason)
}

func coupleAPIErrors(r *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		// an error was raised in go code, no need to check the resty Response
		return nil, NewError(err)
	}

	if r.Error() == nil {
		// no error in the resty Response
		return r, nil
	}

	// handle the resty Response errors

	// Check that response is of the correct content-type before unmarshalling
	expectedContentType := r.Request.Header.Get("Accept")
	responseContentType := r.Header().Get("Content-Type")

	// If the upstream Linode API server being fronted fails to respond to the request,
	// the http server will respond with a default "Bad Gateway" page with Content-Type
	// "text/html".
	if r.StatusCode() == http.StatusBadGateway && responseContentType == "text/html" { //nolint:goconst
		return nil, Error{Code: http.StatusBadGateway, Message: http.StatusText(http.StatusBadGateway)}
	}

	if responseContentType != expectedContentType {
		msg := fmt.Sprintf(
			"Unexpected Content-Type: Expected: %v, Received: %v\nResponse body: %s",
			expectedContentType,
			responseContentType,
			string(r.Body()),
		)

		return nil, Error{Code: r.StatusCode(), Message: msg}
	}

	apiError, ok := r.Error().(*APIError)
	if !ok || (ok && len(apiError.Errors) == 0) {
		return r, nil
	}

	return nil, NewError(r)
}

//nolint:unused
func coupleAPIErrorsHTTP(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		// an error was raised in go code, no need to check the http.Response
		return nil, NewError(err)
	}

	if resp == nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Check that response is of the correct content-type before unmarshalling
		expectedContentType := resp.Request.Header.Get("Accept")
		responseContentType := resp.Header.Get("Content-Type")

		// If the upstream server fails to respond to the request,
		// the http server will respond with a default error page with Content-Type "text/html".
		if resp.StatusCode == http.StatusBadGateway && responseContentType == "text/html" { //nolint:goconst
			return nil, Error{Code: http.StatusBadGateway, Message: http.StatusText(http.StatusBadGateway)}
		}

		if responseContentType != expectedContentType {
			bodyBytes, _ := io.ReadAll(resp.Body)
			msg := fmt.Sprintf(
				"Unexpected Content-Type: Expected: %v, Received: %v\nResponse body: %s",
				expectedContentType,
				responseContentType,
				string(bodyBytes),
			)

			return nil, Error{Code: resp.StatusCode, Message: msg}
		}

		var apiError APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return nil, NewError(fmt.Errorf("failed to decode response body: %w", err))
		}

		if len(apiError.Errors) == 0 {
			return resp, nil
		}

		return nil, Error{Code: resp.StatusCode, Message: apiError.Errors[0].String()}
	}

	// no error in the http.Response
	return resp, nil
}

func (e APIError) Error() string {
	x := []string{}
	for _, msg := range e.Errors {
		x = append(x, msg.Error())
	}

	return strings.Join(x, "; ")
}

func (err Error) Error() string {
	return fmt.Sprintf("[%03d] %s", err.Code, err.Message)
}

func (err Error) StatusCode() int {
	return err.Code
}

func (err Error) Is(target error) bool {
	if x, ok := target.(interface{ StatusCode() int }); ok || errors.As(target, &x) {
		return err.StatusCode() == x.StatusCode()
	}

	return false
}

// NewError creates a linodego.Error with a Code identifying the source err type,
// - ErrorFromString   (1) from a string
// - ErrorFromError    (2) for an error
// - ErrorFromStringer (3) for a Stringer
// - HTTP Status Codes (100-600) for a resty.Response object
func NewError(err any) *Error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *Error:
		return e
	case *resty.Response:
		apiError, ok := e.Error().(*APIError)

		if !ok {
			return &Error{Code: ErrorUnsupported, Message: "Unexpected Resty Error Response, no error"}
		}

		return &Error{
			Code:     e.RawResponse.StatusCode,
			Message:  apiError.Error(),
			Response: e.RawResponse,
		}
	case error:
		return &Error{Code: ErrorFromError, Message: e.Error()}
	case string:
		return &Error{Code: ErrorFromString, Message: e}
	case fmt.Stringer:
		return &Error{Code: ErrorFromStringer, Message: e.String()}
	default:
		return &Error{Code: ErrorUnsupported, Message: fmt.Sprintf("Unsupported type to linodego.NewError: %s", reflect.TypeOf(e))}
	}
}

// IsNotFound indicates if err indicates a 404 Not Found error from the Linode API.
func IsNotFound(err error) bool {
	return ErrHasStatus(err, http.StatusNotFound)
}

// ErrHasStatus checks if err is an error from the Linode API, and whether it contains the given HTTP status code.
// More than one status code may be given.
// If len(code) == 0, err is nil or is not a [Error], ErrHasStatus will return false.
func ErrHasStatus(err error, code ...int) bool {
	if err == nil {
		return false
	}

	// Short-circuit if the caller did not provide any status codes.
	if len(code) == 0 {
		return false
	}

	var e *Error
	if !errors.As(err, &e) {
		return false
	}
	ec := e.StatusCode()
	for _, c := range code {
		if ec == c {
			return true
		}
	}
	return false
}
