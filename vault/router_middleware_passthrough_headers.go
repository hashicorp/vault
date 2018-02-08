package vault

import (
	"context"

	"github.com/hashicorp/vault/logical"
)

var _ RouterMiddleware = (*PassthroughHeadersMiddleware)(nil)

type PassthroughHeadersMiddleware struct {
	// Map of path to applicable headers
	Headers map[string][]string
}

// Update will update m.Headers from the merged config
func (m *PassthroughHeadersMiddleware) Update(config map[string]interface{}) error {
	return nil
}

func (m *PassthroughHeadersMiddleware) Handler(next LogicalHandlerFunc) LogicalHandlerFunc {
	return LogicalHandlerFunc(func(ctx context.Context, req *logical.Request) (*logical.Response, error) {
		// Alter headers and then call next, for example:

		// Get headers from struct
		// passthroughHeaders = m.filteredHeaders(req.Path)

		// Set headers
		// req.Headers = passthroughHeaders

		return next(ctx, req)

		// If we needed to modity the resp, we could call next(ctx, req),
		// alter the response and then return resp, err.
	})
}

// Returns all headers for particular path
func (m *PassthroughHeadersMiddleware) filteredHeaders(string) []string {
	return []string{}
}
