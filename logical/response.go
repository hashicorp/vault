package logical

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/mitchellh/copystructure"
)

const (
	// HTTPContentType can be specified in the Data field of a Response
	// so that the HTTP front end can specify a custom Content-Type associated
	// with the HTTPRawBody. This can only be used for non-secrets, and should
	// be avoided unless absolutely necessary, such as implementing a specification.
	// The value must be a string.
	HTTPContentType = "http_content_type"

	// HTTPRawBody is the raw content of the HTTP body that goes with the HTTPContentType.
	// This can only be specified for non-secrets, and should should be similarly
	// avoided like the HTTPContentType. The value must be a byte slice.
	HTTPRawBody = "http_raw_body"

	// HTTPStatusCode is the response code of the HTTP body that goes with the HTTPContentType.
	// This can only be specified for non-secrets, and should should be similarly
	// avoided like the HTTPContentType. The value must be an integer.
	HTTPStatusCode = "http_status_code"
)

type WrapInfo struct {
	// Setting to non-zero specifies that the response should be wrapped.
	// Specifies the desired TTL of the wrapping token.
	TTL time.Duration

	// The token containing the wrapped response
	Token string

	// The creation time. This can be used with the TTL to figure out an
	// expected expiration.
	CreationTime time.Time

	// If the contained response is the output of a token creation call, the
	// created token's accessor will be accessible here
	WrappedAccessor string
}

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

	// Warnings allow operations or backends to return warnings in response
	// to user actions without failing the action outright.
	// Making it private helps ensure that it is easy for various parts of
	// Vault (backend, core, etc.) to add warnings without accidentally
	// replacing what exists.
	warnings []string

	// Information for wrapping the response in a cubbyhole
	WrapInfo *WrapInfo
}

func init() {
	copystructure.Copiers[reflect.TypeOf(Response{})] = func(v interface{}) (interface{}, error) {
		input := v.(Response)
		ret := Response{
			Redirect: input.Redirect,
		}

		if input.Secret != nil {
			retSec, err := copystructure.Copy(input.Secret)
			if err != nil {
				return nil, fmt.Errorf("error copying Secret: %v", err)
			}
			ret.Secret = retSec.(*Secret)
		}

		if input.Auth != nil {
			retAuth, err := copystructure.Copy(input.Auth)
			if err != nil {
				return nil, fmt.Errorf("error copying Auth: %v", err)
			}
			ret.Auth = retAuth.(*Auth)
		}

		if input.Data != nil {
			retData, err := copystructure.Copy(&input.Data)
			if err != nil {
				return nil, fmt.Errorf("error copying Data: %v", err)
			}
			ret.Data = retData.(map[string]interface{})
		}

		if input.Warnings() != nil {
			for _, warning := range input.Warnings() {
				ret.AddWarning(warning)
			}
		}

		if input.WrapInfo != nil {
			retWrapInfo, err := copystructure.Copy(input.WrapInfo)
			if err != nil {
				return nil, fmt.Errorf("error copying WrapInfo: %v", err)
			}
			ret.WrapInfo = retWrapInfo.(*WrapInfo)
		}

		return &ret, nil
	}
}

// AddWarning adds a warning into the response's warning list
func (r *Response) AddWarning(warning string) {
	if r.warnings == nil {
		r.warnings = make([]string, 0, 1)
	}
	r.warnings = append(r.warnings, warning)
}

// Warnings returns the list of warnings set on the response
func (r *Response) Warnings() []string {
	return r.warnings
}

// ClearWarnings clears the response's warning list
func (r *Response) ClearWarnings() {
	r.warnings = make([]string, 0, 1)
}

// Copies the warnings from the other response to this one
func (r *Response) CloneWarnings(other *Response) {
	r.warnings = other.warnings
}

// IsError returns true if this response seems to indicate an error.
func (r *Response) IsError() bool {
	return r != nil && len(r.Data) == 1 && r.Data["error"] != nil
}

func (r *Response) Error() error {
	if !r.IsError() {
		return nil
	}
	switch r.Data["error"].(type) {
	case string:
		return errors.New(r.Data["error"].(string))
	case error:
		return r.Data["error"].(error)
	}
	return nil
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
	resp := &Response{
		Data: map[string]interface{}{},
	}
	if len(keys) != 0 {
		resp.Data["keys"] = keys
	}
	return resp
}
