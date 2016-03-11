package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

	// FIXME: Change the sealed code to http.StatusServiceUnavailable at some
	// point
	sealedCode := http.StatusInternalServerError
	standbyCode := http.StatusTooManyRequests // Consul warning code
	activeCode := http.StatusOK

	var err error
	sealedCodeStr, sealedCodeOk := r.URL.Query()["sealedcode"]
	if sealedCodeOk {
		if len(sealedCodeStr) < 1 {
			respondError(w, http.StatusBadRequest, nil)
			return
		}
		sealedCode, err = strconv.Atoi(sealedCodeStr[0])
		if err != nil {
			respondError(w, http.StatusBadRequest, nil)
			return
		}
	}
	standbyCodeStr, standbyCodeOk := r.URL.Query()["standbycode"]
	if standbyCodeOk {
		if len(standbyCodeStr) < 1 {
			respondError(w, http.StatusBadRequest, nil)
			return
		}
		standbyCode, err = strconv.Atoi(standbyCodeStr[0])
		if err != nil {
			respondError(w, http.StatusBadRequest, nil)
			return
		}
	}

	activeCodeStr, activeCodeOk := r.URL.Query()["activecode"]
	if activeCodeOk {
		if len(activeCodeStr) < 1 {
			respondError(w, http.StatusBadRequest, nil)
			return
		}

		activeCode, err = strconv.Atoi(activeCodeStr[0])
		if err != nil {
			respondError(w, http.StatusBadRequest, nil)
			return
		}
	}

	// Check system status
	sealed, _ := core.Sealed()
	standby, _ := core.Standby()
	init, err := core.Initialized()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Determine the status code
	code := activeCode
	switch {
	case !init:
		code = http.StatusInternalServerError
	case sealed:
		code = sealedCode
	case !standbyOK && standby:
		code = standbyCode
	}

	// Format the body
	body := &HealthResponse{
		Initialized:   init,
		Sealed:        sealed,
		Standby:       standby,
		ServerTimeUTC: time.Now().UTC().Unix(),
	}

	// Generate the response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.Encode(body)
}

type HealthResponse struct {
	Initialized   bool  `json:"initialized"`
	Sealed        bool  `json:"sealed"`
	Standby       bool  `json:"standby"`
	ServerTimeUTC int64 `json:"server_time_utc"`
}
