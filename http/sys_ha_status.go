package http

import (
	"net/http"
	"time"

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
	nodes := []Node{
		{
			Hostname:       h.Hostname,
			APIAddress:     address,
			ClusterAddress: clusterAddr,
			ActiveNode:     true,
		},
	}

	for _, peerNode := range core.GetHAPeerNodesCached() {
		lastEcho := peerNode.LastEcho
		nodes = append(nodes, Node{
			Hostname:       peerNode.Hostname,
			APIAddress:     peerNode.APIAddress,
			ClusterAddress: peerNode.ClusterAddress,
			LastEcho:       &lastEcho,
		})
	}
	resp := &HaStatusResponse{
		Nodes: nodes,
	}

	respondOk(w, resp)
}

type Node struct {
	Hostname       string     `json:"hostname"`
	APIAddress     string     `json:"api_address"`
	ClusterAddress string     `json:"cluster_address"`
	ActiveNode     bool       `json:"active_node"`
	LastEcho       *time.Time `json:"last_echo"`
}

type HaStatusResponse struct {
	Nodes []Node
}
