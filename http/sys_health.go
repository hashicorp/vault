package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/version"
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
		core.Logger().Error("error checking health", "error", err)
		respondError(w, code, nil)
		return
	}

	if body == nil {
		respondError(w, code, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Generate the response
	enc := json.NewEncoder(w)
	enc.Encode(body)
}

func handleSysHealthHead(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	code, body, _ := getSysHealth(core, r)

	if body != nil {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(code)
}

func getSysHealth(core *vault.Core, r *http.Request) (int, *HealthResponse, error) {
	var err error

	// Check if being a standby is allowed for the purpose of a 200 OK
	standbyOKStr, standbyOK := r.URL.Query()["standbyok"]
	if standbyOK {
		standbyOK, err = parseutil.ParseBool(standbyOKStr[0])
		if err != nil {
			return http.StatusBadRequest, nil, fmt.Errorf("bad value for standbyok parameter: %w", err)
		}
	}
	perfStandbyOKStr, perfStandbyOK := r.URL.Query()["perfstandbyok"]
	if perfStandbyOK {
		perfStandbyOK, err = parseutil.ParseBool(perfStandbyOKStr[0])
		if err != nil {
			return http.StatusBadRequest, nil, fmt.Errorf("bad value for perfstandbyok parameter: %w", err)
		}
	}

	uninitCode := http.StatusNotImplemented
	if code, found, ok := fetchStatusCode(r, "uninitcode"); !ok {
		return http.StatusBadRequest, nil, nil
	} else if found {
		uninitCode = code
	}

	sealedCode := http.StatusServiceUnavailable
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

	drSecondaryCode := 472 // unofficial 4xx status code
	if code, found, ok := fetchStatusCode(r, "drsecondarycode"); !ok {
		return http.StatusBadRequest, nil, nil
	} else if found {
		drSecondaryCode = code
	}

	perfStandbyCode := 473 // unofficial 4xx status code
	if code, found, ok := fetchStatusCode(r, "performancestandbycode"); !ok {
		return http.StatusBadRequest, nil, nil
	} else if found {
		perfStandbyCode = code
	}

	ctx := context.Background()

	// Check system status
	sealed := core.Sealed()
	standby, perfStandby := core.StandbyStates()
	var replicationState consts.ReplicationState
	if standby {
		replicationState = core.ActiveNodeReplicationState()
	} else {
		replicationState = core.ReplicationState()
	}

	init, err := core.Initialized(ctx)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Determine the status code
	code := activeCode
	switch {
	case !init:
		code = uninitCode
	case sealed:
		code = sealedCode
	case replicationState.HasState(consts.ReplicationDRSecondary):
		code = drSecondaryCode
	case perfStandby:
		if !perfStandbyOK {
			code = perfStandbyCode
		}
	case standby:
		if !standbyOK {
			code = standbyCode
		}
	}

	// Fetch the local cluster name and identifier
	var clusterName, clusterID string
	if !sealed {
		cluster, err := core.Cluster(ctx)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
		if cluster == nil {
			return http.StatusInternalServerError, nil, fmt.Errorf("failed to fetch cluster details")
		}
		clusterName = cluster.Name
		clusterID = cluster.ID
	}

	// Format the body
	body := &HealthResponse{
		Initialized:                init,
		Sealed:                     sealed,
		Standby:                    standby,
		PerformanceStandby:         perfStandby,
		ReplicationPerformanceMode: replicationState.GetPerformanceString(),
		ReplicationDRMode:          replicationState.GetDRString(),
		ServerTimeUTC:              time.Now().UTC().Unix(),
		Version:                    version.GetVersion().VersionNumber(),
		ClusterName:                clusterName,
		ClusterID:                  clusterID,
	}

	licenseState, err := vault.LicenseSummary(core)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if licenseState != nil {
		body.License = &HealthResponseLicense{
			State:      licenseState.State,
			Terminated: licenseState.Terminated,
		}
		if !licenseState.ExpiryTime.IsZero() {
			body.License.ExpiryTime = licenseState.ExpiryTime.Format(time.RFC3339)
		}
	}

	if init && !sealed && !standby {
		body.LastWAL = vault.LastWAL(core)
	}

	return code, body, nil
}

type HealthResponseLicense struct {
	State      string `json:"state"`
	ExpiryTime string `json:"expiry_time"`
	Terminated bool   `json:"terminated"`
}

type HealthResponse struct {
	Initialized                bool                   `json:"initialized"`
	Sealed                     bool                   `json:"sealed"`
	Standby                    bool                   `json:"standby"`
	PerformanceStandby         bool                   `json:"performance_standby"`
	ReplicationPerformanceMode string                 `json:"replication_performance_mode"`
	ReplicationDRMode          string                 `json:"replication_dr_mode"`
	ServerTimeUTC              int64                  `json:"server_time_utc"`
	Version                    string                 `json:"version"`
	ClusterName                string                 `json:"cluster_name,omitempty"`
	ClusterID                  string                 `json:"cluster_id,omitempty"`
	LastWAL                    uint64                 `json:"last_wal,omitempty"`
	License                    *HealthResponseLicense `json:"license,omitempty"`
}
