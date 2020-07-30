package http

import (
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/vault"
	"github.com/shirou/gopsutil/host"
)

func handleSysHaStatus(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysHaStatusGet(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysHaStatusGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	_, address, clusterAddr, err := core.Leader()
	if errwrap.Contains(err, vault.ErrHANotEnabled.Error()) {
		err = nil
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	h, _ := host.Info()
	leaderAddr := NodeAddr{
		HostName:       h.Hostname,
		ApiAddress:     address,
		ClusterAddress: clusterAddr,
	}

	addrsCache := core.GetclusterPeerClusterAddrsCache()
	standbyAddr := []NodeAddr{}
	for itemClusterAddr, item := range addrsCache {
		standbyAddr = append(standbyAddr, NodeAddr{
			ClusterAddress: itemClusterAddr,
			ApiAddress:     item.Object.(*vault.NodeInformation).ApiAddr,
			HostName:       item.Object.(*vault.NodeInformation).NodeID,
		})
	}
	resp := &HaStatusResponse{
		Leader:  &leaderAddr,
		Standby: &standbyAddr,
	}

	respondOk(w, resp)
}

type NodeAddr struct {
	HostName       string `json:"host_name"`
	ApiAddress     string `json:"api_addr"`
	ClusterAddress string `json:"cluster_addr"`
}

type HaStatusResponse struct {
	Leader      *NodeAddr   `json:"leader"`
	PerfStandby *[]NodeAddr `json:"performance_standby"`
	Standby     *[]NodeAddr `json:"standby"`
}
