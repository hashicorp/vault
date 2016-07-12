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
		case "HEAD":
			handleSysHealthHead(core, w, r)
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
	code, body, err := getSysHealth(core, r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	if body == nil {
		respondError(w, code, nil)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	// Generate the response
	enc := json.NewEncoder(w)
	enc.Encode(body)
}

func handleSysHealthHead(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	code, body, err := getSysHealth(core, r)
	if err != nil {
		code = http.StatusInternalServerError
	}

	if body != nil {
		w.Header().Add("Content-Type", "application/json")
	}
	w.WriteHeader(code)
}

func getSysHealth(core *vault.Core, r *http.Request) (int, *HealthResponse, error) {
	// Check if being a standby is allowed for the purpose of a 200 OK
	_, standbyOK := r.URL.Query()["standbyok"]

	// FIXME: Change the sealed code to http.StatusServiceUnavailable at some
	// point
	sealedCode := http.StatusInternalServerError
	if code, found, ok := fetchStatusCode(r, "sealedcode"); !ok {
		return http.StatusBadRequest, nil, nil
	} else if found {
		sealedCode = code
	}

	standbyCode := http.StatusTooManyRequests // Consul warning code
	if code, found, ok := fetchStatusCode(r, "standbycode"); !ok {
		return http.StatusBadRequest, nil, nil
	} else if found {
		standbyCode = code
	}

	activeCode := http.StatusOK
	if code, found, ok := fetchStatusCode(r, "activecode"); !ok {
		return http.StatusBadRequest, nil, nil
	} else if found {
		activeCode = code
	}

	// Check system status
	sealed, _ := core.Sealed()
	standby, _ := core.Standby()
	init, err := core.Initialized()
	if err != nil {
		return http.StatusInternalServerError, nil, err
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
	return code, body, nil
}

type HealthResponse struct {
	Initialized   bool  `json:"initialized"`
	Sealed        bool  `json:"sealed"`
	Standby       bool  `json:"standby"`
	ServerTimeUTC int64 `json:"server_time_utc"`
}
