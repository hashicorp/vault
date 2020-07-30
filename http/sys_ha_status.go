package http

import (
	"fmt"
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
		conf := core.SanitizedConfig()
		//If Vault is not HA enabled,use the address from the configuration
		address = conf["api_addr"].(string)
		clusterAddr = conf["cluster_addr"].(string)

		//If addresses is not defined at the top level in the config,use the first listener addresses as last resort
		listener := conf["listeners"].([]interface{})
		if address == "" && len(listener) > 0 {
			rawConf := listener[0].(map[string]interface{})["config"].(map[string]interface{})
			address = fmt.Sprintf("%s", rawConf["address"])

		}
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
