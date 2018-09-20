package http

import (
	"net/http"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/vault"
)

var (
	adjustRequest = func(c *vault.Core, r *http.Request) (*http.Request, int) {
		return r.WithContext(namespace.ContextWithNamespace(r.Context(), namespace.RootNamespace)), 0
	}

	genericWrapping = func(core *vault.Core, in http.Handler, props *vault.HandlerProperties) http.Handler {
		// Wrap the help wrapped handler with another layer with a generic
		// handler
		return wrapGenericHandler(core, in, props.MaxRequestSize, props.MaxRequestDuration)
	}

	additionalRoutes = func(mux *http.ServeMux, core *vault.Core) {}
)
