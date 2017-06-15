package api

func (c *Sys) Health() (*HealthResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/health")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result HealthResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

type HealthResponse struct {
	Initialized   bool   `json:"initialized"`
	Sealed        bool   `json:"sealed"`
	Standby       bool   `json:"standby"`
	ServerTimeUTC int64  `json:"server_time_utc"`
	Version       string `json:"version"`
	ClusterName   string `json:"cluster_name,omitempty"`
	ClusterID     string `json:"cluster_id,omitempty"`
}
