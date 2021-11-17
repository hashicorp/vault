package cfclient

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

// StatsGetResponse is the json body returned from the API
type StatsGetResponse struct {
	Stats []Stats `json:"resources"`
}

// Stats represents the stats of a process
type Stats struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	State string `json:"state"`
	Usage struct {
		Time string  `json:"time"`
		CPU  float64 `json:"cpu"`
		Mem  int     `json:"mem"`
		Disk int     `json:"disk"`
	} `json:"usage"`
	Host          string `json:"host"`
	InstancePorts []struct {
		External             int `json:"external"`
		Internal             int `json:"internal"`
		ExternalTLSProxyPort int `json:"external_tls_proxy_port"`
		InternalTLSProxyPort int `json:"internal_tls_proxy_port"`
	} `json:"instance_ports"`
	Uptime           int    `json:"uptime"`
	MemQuota         int    `json:"mem_quota"`
	DiskQuota        int    `json:"disk_quota"`
	FdsQuota         int    `json:"fds_quota"`
	IsolationSegment string `json:"isolation_segment"`
	Details          string `json:"details"`
}

func (c *Client) GetProcessStats(processGUID string) ([]Stats, error) {
	req := c.NewRequest("GET", "/v3/processes/"+processGUID+"/stats")
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "error getting stats for v3 process")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting stats with GUID [%s], response code: %d", processGUID, resp.StatusCode)
	}

	statsResp := new(StatsGetResponse)
	err = json.NewDecoder(resp.Body).Decode(statsResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding stats with GUID [%s], response code: %d", processGUID, resp.StatusCode)
	}
	return statsResp.Stats, nil
}
