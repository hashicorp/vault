package http

import (
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/vault"
)

func handleSysLeader(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysLeaderGet(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysLeaderGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	haEnabled := true
	isLeader, address, clusterAddr, err := core.Leader()
	if errwrap.Contains(err, vault.ErrHANotEnabled.Error()) {
		haEnabled = false
		err = nil
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	resp := &LeaderResponse{
		HAEnabled:            haEnabled,
		IsSelf:               isLeader,
		LeaderAddress:        address,
		LeaderClusterAddress: clusterAddr,
		PerfStandby:          core.PerfStandby(),
	}
	if resp.PerfStandby {
		resp.PerfStandbyLastRemoteWAL = vault.LastRemoteWAL(core)
	} else if isLeader || !haEnabled {
		resp.LastWAL = vault.LastWAL(core)
	}

	respondOk(w, resp)
}

type LeaderResponse struct {
	HAEnabled                bool   `json:"ha_enabled"`
	IsSelf                   bool   `json:"is_self"`
	LeaderAddress            string `json:"leader_address"`
	LeaderClusterAddress     string `json:"leader_cluster_address"`
	PerfStandby              bool   `json:"performance_standby"`
	PerfStandbyLastRemoteWAL uint64 `json:"performance_standby_last_remote_wal"`
	LastWAL                  uint64 `json:"last_wal,omitempty"`
}
