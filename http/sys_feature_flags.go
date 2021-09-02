package http

import (
	"encoding/json"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"net/http"
	"os"

	"github.com/hashicorp/vault/vault"
)

type FeatureFlagsResponse struct {
	FeatureFlags []string `json:"feature_flags"`
}

var FeatureFlag_EnvVariables = [...]string{
	"VAULT_CLOUD_ADMIN_NAMESPACE",
}

func featureFlagIsSet(name string) bool {
	switch os.Getenv(name) {
	case "", "0":
		return false
	default:
		return true
	}
}

func handleSysInternalFeatureFlags(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Getting custom headers from listener's config
		la := w.Header().Get("X-Vault-Listener-Add")
		lc, err := core.GetCustomResponseHeaders(la)
		if err != nil {
			core.Logger().Debug("failed to get custom headers from listener config")
		}
		switch r.Method {
		case "GET":
			break
		default:
			respondError(w, http.StatusMethodNotAllowed, nil, lc)
		}

		response := &FeatureFlagsResponse{}

		for _, f := range FeatureFlag_EnvVariables {
			if featureFlagIsSet(f) {
				response.FeatureFlags = append(response.FeatureFlags, f)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		status := http.StatusOK
		listenerutil.SetCustomResponseHeaders(lc, w, status)
		w.WriteHeader(status)

		// Generate the response
		enc := json.NewEncoder(w)
		enc.Encode(response)
	})
}
