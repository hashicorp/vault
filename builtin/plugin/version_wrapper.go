package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	v5 "github.com/hashicorp/vault/builtin/plugin/v5"
	"github.com/hashicorp/vault/sdk/logical"
)

// NewBackendWrapper figures out which version of the database the pluginName is referring to and returns a wrapper object
// that can be used to make operations on the underlying database plugin.
func NewBackendWrapper(ctx context.Context, conf *logical.BackendConfig) logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		newBackend, err := v5.Factory(ctx, conf)
		if err == nil {
			return newBackend, nil
		}

		merr := &multierror.Error{}
		merr = multierror.Append(merr, err)

		// legacy backends will be lazy-loaded and do not support multiplexing
		legacyBackend, err := Factory(ctx, conf)
		if err == nil {
			return legacyBackend, nil
		}
		merr = multierror.Append(merr, err)

		return nil, fmt.Errorf("invalid backend version: %s", merr)
	}
}
