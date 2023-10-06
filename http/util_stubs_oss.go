//go:build !enterprise

package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func entGenericWrapping(core *vault.Core, in http.Handler, props *vault.HandlerProperties) http.Handler {
	// Wrap the help wrapped handler with another layer with a generic
	// handler
	return wrapGenericHandler(core, in, props)
}

func entAdditionalRoutes(mux *http.ServeMux, core *vault.Core) {}

func entAdjustResponse() {}
