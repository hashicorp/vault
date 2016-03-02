package api

func (c *Sys) CapabilitiesSelf(path string) ([]string, error) {
	body := map[string]string{
		"path": path,
	}

	r := c.c.NewRequest("POST", "/v1/sys/capabilities-self")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result capabilitiesResp
	err = resp.DecodeJSON(&result)
	return result.Capabilities, err
}

func (c *Sys) Capabilities(token, path string) ([]string, error) {
	body := map[string]string{
		"token": token,
		"path":  path,
	}

	r := c.c.NewRequest("POST", "/v1/sys/capabilities")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result capabilitiesResp
	err = resp.DecodeJSON(&result)
	return result.Capabilities, err
}

type capabilitiesResp struct {
	Capabilities []string `json:"capabilities"`
}
