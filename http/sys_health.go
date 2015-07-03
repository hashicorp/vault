package http

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func handleSysHealth(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysHealthGet(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysHealthGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Check if being a standby is allowed for the purpose of a 200 OK
	_, standbyOK := r.URL.Query()["standbyok"]

	// Check system status
	sealed, _ := core.Sealed()
	standby, _ := core.Standby()
	init, err := core.Initialized()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Determine the status code
	code := http.StatusOK
	switch {
	case !init:
		code = http.StatusInternalServerError
	case sealed:
		code = http.StatusInternalServerError
	case !standbyOK && standby:
		code = 429 // Consul warning code
	}

	// Format the body
	body := &HealthResponse{
		Initialized: init,
		Sealed:      sealed,
		Standby:     standby,
	}

	// Generate the response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.Encode(body)
}

type HealthResponse struct {
	Initialized bool `json:"initialized"`
	Sealed      bool `json:"sealed"`
	Standby     bool `json:"standby"`
}
