package api

func (c *Sys) InitStatus() (bool, error) {
	r := c.c.NewRequest("GET", "/v1/sys/init")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result InitStatusResponse
	err = resp.DecodeJSON(&result)
	return result.Initialized, err
}

func (c *Sys) Init(opts *InitRequest) (*InitResponse, error) {
	body := map[string]interface{}{
		"secret_shares":    opts.SecretShares,
		"secret_threshold": opts.SecretThreshold,
	}

	r := c.c.NewRequest("PUT", "/v1/sys/init")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result InitResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

type InitRequest struct {
	SecretShares    int
	SecretThreshold int
}

type InitStatusResponse struct {
	Initialized bool
}

type InitResponse struct {
	Keys      []string
	RootToken string `json:"root_token"`
}
