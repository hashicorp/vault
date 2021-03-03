package http

import (
	"encoding/json"
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
		switch r.Method {
		case "GET":
			break
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}

		response := &FeatureFlagsResponse{}

		for _, f := range FeatureFlag_EnvVariables {
			if featureFlagIsSet(f) {
				response.FeatureFlags = append(response.FeatureFlags, f)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Generate the response
		enc := json.NewEncoder(w)
		enc.Encode(response)
	})
}
