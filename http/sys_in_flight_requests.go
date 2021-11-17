package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/vault"
)

func handleUnAuthenticatedInFlightRequest(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}
		now := time.Now()

		currentInFlightReqMap := core.LoadInFlightReqData()

		for _, v := range currentInFlightReqMap {
			v.SnapshotTime = now
		}

		content, err := json.Marshal(currentInFlightReqMap)
		if err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Errorf("error while marshalling the in-flight requests data: %w", err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(content)

	})
}
