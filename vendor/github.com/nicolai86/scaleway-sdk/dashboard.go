package api

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// DashboardResp represents a dashboard received from the API
type DashboardResp struct {
	Dashboard Dashboard
}

// Dashboard represents a dashboard
type Dashboard struct {
	VolumesCount        int `json:"volumes_count"`
	RunningServersCount int `json:"running_servers_count"`
	ImagesCount         int `json:"images_count"`
	SnapshotsCount      int `json:"snapshots_count"`
	ServersCount        int `json:"servers_count"`
	IPsCount            int `json:"ips_count"`
}

// GetDashboard returns the dashboard
func (s *API) GetDashboard() (*Dashboard, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "dashboard", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var dashboard DashboardResp

	if err = json.Unmarshal(body, &dashboard); err != nil {
		return nil, err
	}
	return &dashboard.Dashboard, nil
}
