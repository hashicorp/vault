// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
)

func handleSysHealth(core *vault.Core, opt ...ListenerConfigOption) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysHealthGet(core, w, r, opt...)
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

func handleSysHealthGet(core *vault.Core, w http.ResponseWriter, r *http.Request, opt ...ListenerConfigOption) {
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

	var tokenPresent bool
	token := r.Header.Get(consts.AuthHeaderName)

	if token != "" {
		// We don't care about the error, we just want to know if the token exists
		lock := core.HALock()
		lock.Lock()
		tokenEntry, err := core.LookupToken(r.Context(), token)
		lock.Unlock()
		tokenPresent = err == nil && tokenEntry != nil
	}
	opts, _ := getOpts(opt...)

	if !tokenPresent {
		if opts.withRedactVersion {
			body.Version = opts.withRedactionValue
		}

		if opts.withRedactClusterName {
			body.ClusterName = opts.withRedactionValue
		}
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

	haUnhealthyCode := 474
	if code, found, ok := fetchStatusCode(r, "haunhealthycode"); !ok {
		return http.StatusBadRequest, nil, nil
	} else if found {
		haUnhealthyCode = code
	}

	removedCode := 530
	if code, found, ok := fetchStatusCode(r, "removedcode"); !ok {
		return http.StatusBadRequest, nil, nil
	} else if found {
		removedCode = code
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

	removed, shouldIncludeRemoved := core.IsRemovedFromCluster()

	haHealthy, lastHeartbeat := core.GetHAHeartbeatHealth()

	// Determine the status code
	code := activeCode
	switch {
	case !init:
		code = uninitCode
	case removed:
		code = removedCode
	case sealed:
		code = sealedCode
	case !haHealthy && lastHeartbeat != nil:
		code = haUnhealthyCode
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
		Enterprise:                 constants.IsEnterprise,
		ClusterName:                clusterName,
		ClusterID:                  clusterID,
		ClockSkewMillis:            core.ActiveNodeClockSkewMillis(),
		EchoDurationMillis:         core.EchoDuration().Milliseconds(),
	}
	if standby {
		body.ReplicationPrimaryCanaryAgeMillis = core.GetReplicationLagMillisIgnoreErrs()
	}

	licenseState, err := core.EntGetLicenseState()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	if shouldIncludeRemoved {
		body.RemovedFromCluster = &removed
	}

	if lastHeartbeat != nil {
		body.LastRequestForwardingHeartbeatMillis = lastHeartbeat.Milliseconds()
		body.HAConnectionHealthy = &haHealthy
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
		body.LastWAL = core.EntLastWAL()
	}

	return code, body, nil
}

type HealthResponseLicense struct {
	State      string `json:"state"`
	ExpiryTime string `json:"expiry_time"`
	Terminated bool   `json:"terminated"`
}

type HealthResponse struct {
	Initialized                          bool                   `json:"initialized"`
	Sealed                               bool                   `json:"sealed"`
	Standby                              bool                   `json:"standby"`
	PerformanceStandby                   bool                   `json:"performance_standby"`
	ReplicationPerformanceMode           string                 `json:"replication_performance_mode"`
	ReplicationDRMode                    string                 `json:"replication_dr_mode"`
	ServerTimeUTC                        int64                  `json:"server_time_utc"`
	Version                              string                 `json:"version"`
	Enterprise                           bool                   `json:"enterprise"`
	ClusterName                          string                 `json:"cluster_name,omitempty"`
	ClusterID                            string                 `json:"cluster_id,omitempty"`
	LastWAL                              uint64                 `json:"last_wal,omitempty"`
	License                              *HealthResponseLicense `json:"license,omitempty"`
	EchoDurationMillis                   int64                  `json:"echo_duration_ms"`
	ClockSkewMillis                      int64                  `json:"clock_skew_ms"`
	ReplicationPrimaryCanaryAgeMillis    int64                  `json:"replication_primary_canary_age_ms"`
	RemovedFromCluster                   *bool                  `json:"removed_from_cluster,omitempty"`
	HAConnectionHealthy                  *bool                  `json:"ha_connection_healthy,omitempty"`
	LastRequestForwardingHeartbeatMillis int64                  `json:"last_request_forwarding_heartbeat_ms,omitempty"`
}
