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

func fetchStatusCode(r *http.Request, field string) (int, bool, bool) {
	var err error
	statusCode := http.StatusOK
	if statusCodeStr, statusCodeOk := r.URL.Query()[field]; statusCodeOk {
		statusCode, err = strconv.Atoi(statusCodeStr[0])
		if err != nil || len(statusCodeStr) < 1 {
			return http.StatusBadRequest, false, false
		}
		return statusCode, true, true
	}
	return statusCode, false, true
}

func handleSysHealthGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {

	// Check if being a standby is allowed for the purpose of a 200 OK
	_, standbyOK := r.URL.Query()["standbyok"]

	// FIXME: Change the sealed code to http.StatusServiceUnavailable at some
	// point
	sealedCode := http.StatusInternalServerError
	if code, found, ok := fetchStatusCode(r, "sealedcode"); !ok {
		respondError(w, http.StatusBadRequest, nil)
		return
	} else if found {
		sealedCode = code
	}

	standbyCode := http.StatusTooManyRequests // Consul warning code
	if code, found, ok := fetchStatusCode(r, "standbycode"); !ok {
		respondError(w, http.StatusBadRequest, nil)
		return
	} else if found {
		standbyCode = code
	}

	activeCode := http.StatusOK
	if code, found, ok := fetchStatusCode(r, "activecode"); !ok {
		respondError(w, http.StatusBadRequest, nil)
		return
	} else if found {
		activeCode = code
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
