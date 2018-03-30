package framework

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
)

// DEPRECATED

// LeaseExtend is left for backwards compatibility for plugins. This function
// now just passes back the data that was passed into it to be processed in core.
func LeaseExtend(backendIncrement, backendMax time.Duration, systemView logical.SystemView) OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
		switch {
		case req.Auth != nil:
			req.Auth.TTL = backendIncrement
			req.Auth.MaxTTL = backendMax
			return &logical.Response{Auth: req.Auth}, nil
		case req.Secret != nil:
			req.Secret.TTL = backendIncrement
			return &logical.Response{Secret: req.Secret}, nil
		}
		return nil, fmt.Errorf("no lease options for request")
	}
}
