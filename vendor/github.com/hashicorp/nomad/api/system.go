package api

// Status is used to query the status-related endpoints.
type System struct {
	client *Client
}

// System returns a handle on the system endpoints.
func (c *Client) System() *System {
	return &System{client: c}
}

func (s *System) GarbageCollect() error {
	var req struct{}
	_, err := s.client.write("/v1/system/gc", &req, nil, nil)
	return err
}

func (s *System) ReconcileSummaries() error {
	var req struct{}
	_, err := s.client.write("/v1/system/reconcile/summaries", &req, nil, nil)
	return err
}
