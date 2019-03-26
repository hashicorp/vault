package cache

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
)

// SendRequest is the input for Proxier.Send.
type SendRequest struct {
	Token   string
	Request *http.Request

	// RequestBody is the stored body bytes from Request.Body. It is set here to
	// avoid reading and re-setting the stream multiple times.
	RequestBody []byte
}

// SendResponse is the output from Proxier.Send.
type SendResponse struct {
	Response *api.Response

	// ResponseBody is the stored body bytes from Response.Body. It is set here to
	// avoid reading and re-setting the stream multiple times.
	ResponseBody []byte
	CacheMeta    *CacheMeta
}

// CacheMeta contains metadata information about the response,
// such as whether it was a cache hit or miss, and the age of the
// cached entry.
type CacheMeta struct {
	Hit bool
	Age time.Duration
}

// Proxier is the interface implemented by different components that are
// responsible for performing specific tasks, such as caching and proxying. All
// these tasks combined together would serve the request received by the agent.
type Proxier interface {
	Send(ctx context.Context, req *SendRequest) (*SendResponse, error)
}

// NewSendResponse creates a new SendResponse and takes care of initializing its
// fields properly.
func NewSendResponse(apiResponse *api.Response, responseBody []byte) (*SendResponse, error) {
	if apiResponse == nil {
		return nil, fmt.Errorf("nil api response provided")
	}

	resp := &SendResponse{
		Response:  apiResponse,
		CacheMeta: &CacheMeta{},
	}

	// If a response body is separately provided we set that as the SendResponse.ResponseBody,
	// otherwise we will do an ioutil.ReadAll to extract the response body from apiResponse.
	switch {
	case len(responseBody) > 0:
		resp.ResponseBody = responseBody
	case apiResponse.Body != nil:
		respBody, err := ioutil.ReadAll(apiResponse.Body)
		if err != nil {
			return nil, err
		}
		// Close the old body
		apiResponse.Body.Close()

		// Re-set the response body after reading from the Reader
		apiResponse.Body = ioutil.NopCloser(bytes.NewReader(respBody))

		resp.ResponseBody = respBody
	}

	return resp, nil
}
