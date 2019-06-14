package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/vault/sdk/helper/jsonutil"
)

// Response is a raw response that wraps an HTTP response.
type Response struct {
	*http.Response
}

// DecodeJSON will decode the response body to a JSON structure. This
// will consume the response body, but will not close it. Close must
// still be called.
func (r *Response) DecodeJSON(out interface{}) error {
	return jsonutil.DecodeJSONFromReader(r.Body, out)
}

// Error returns an error response if there is one. If there is an error,
// this will fully consume the response body, but will not close it. The
// body must still be closed manually.
func (r *Response) Error() error {
	// 200 to 399 are okay status codes. 429 is the code for health status of
	// standby nodes.
	if (r.StatusCode >= 200 && r.StatusCode < 400) || r.StatusCode == 429 {
		return nil
	}

	// We have an error. Let's copy the body into our own buffer first,
	// so that if we can't decode JSON, we can at least copy it raw.
	bodyBuf := &bytes.Buffer{}
	if _, err := io.Copy(bodyBuf, r.Body); err != nil {
		return err
	}

	r.Body.Close()
	r.Body = ioutil.NopCloser(bodyBuf)

	// Decode the error response if we can. Note that we wrap the bodyBuf
	// in a bytes.Reader here so that the JSON decoder doesn't move the
	// read pointer for the original buffer.
	var resp ErrorResponse
	if err := jsonutil.DecodeJSON(bodyBuf.Bytes(), &resp); err != nil {
		// Ignore the decoding error and just drop the raw response
		return &ResponseError{r.Request.Method, r.Request.URL.String(), r.StatusCode, bodyBuf.String(), nil}
	}
	return &ResponseError{r.Request.Method, r.Request.URL.String(), r.StatusCode, "", resp.Errors}
}

// ErrorResponse is the raw structure of errors when they're returned by the
// HTTP API.
type ErrorResponse struct {
	Errors []string
}

// ResponseError represents a Vault client response error.
type ResponseError struct {
	Method     string   // The request method.
	URL        string   // The request URL.
	StatusCode int      // The response status code.
	RawMessage string   // The raw response body if decoding fails.
	Errors     []string // The decoded response body errors.
}

// Error implements the error interface.
func (e *ResponseError) Error() string {
	if e.RawMessage != "" {
		return fmt.Sprintf(
			"Error making API request.\n\n"+
				"URL: %s %s\n"+
				"Code: %d. Raw Message:\n\n%s",
			e.Method, e.URL,
			e.StatusCode, e.RawMessage)
	}
	var errBody bytes.Buffer
	errBody.WriteString(fmt.Sprintf(
		"Error making API request.\n\n"+
			"URL: %s %s\n"+
			"Code: %d. Errors:\n\n",
		e.Method, e.URL,
		e.StatusCode))
	for _, err := range e.Errors {
		errBody.WriteString(fmt.Sprintf("* %s", err))
	}
	return errBody.String()
}
