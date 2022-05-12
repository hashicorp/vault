package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/logical"
)

type backendVersionWrapper struct {
	*logical.Backend
}

// NewBackendWrapper figures out which version of the database the pluginName is referring to and returns a wrapper object
// that can be used to make operations on the underlying database plugin.
func NewBackendWrapper(ctx context.Context, conf *logical.BackendConfig) (wrapper backendVersionWrapper, err error) {
	newBackend, err := v5.Factory(ctx, conf)
	if err == nil {
		wrapper = backendVersionWrapper{
			v5: newBackend,
		}
		return wrapper, nil
	}

	merr := &multierror.Error{}
	merr = multierror.Append(merr, err)

	legacyBackend, err := Factory(ctx, conf)
	if err == nil {
		wrapper = backendVersionWrapper{
			v4: legacyBackend,
		}
		return wrapper, nil
	}
	merr = multierror.Append(merr, err)

	return wrapper, fmt.Errorf("invalid backend version: %s", merr)
}
