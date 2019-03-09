package cache

import (
	"context"
	"net/http"

	"github.com/hashicorp/vault/api"
)

// SendRequest is the input for Proxier.Send.
type SendRequest struct {
	Token       string
	Request     *http.Request
	RequestBody []byte
}

// SendResponse is the output from Proxier.Send.
type SendResponse struct {
	Response     *api.Response
	ResponseBody []byte
}

// Proxier is the interface implemented by different components that are
// responsible for performing specific tasks, such as caching and proxying. All
// these tasks combined together would serve the request received by the agent.
type Proxier interface {
	Send(ctx context.Context, req *SendRequest) (*SendResponse, error)
}
