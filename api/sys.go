package api

// Sys is used to perform system-related operations on Vault.
type Sys struct {
	c *Client
}

// Sys is used to return the client for sys-related API calls.
func (c *Client) Sys() *Sys {
	return &Sys{c: c}
}

func (c *Sys) SealStatus() (*SealStatusResponse, error) {
	r := c.c.NewRequest("GET", "/sys/seal-status")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SealStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) Seal() error {
	r := c.c.NewRequest("PUT", "/sys/seal")
	resp, err := c.c.RawRequest(r)
	defer resp.Body.Close()
	return err
}

func (c *Sys) Unseal(shard string) (*SealStatusResponse, error) {
	body := map[string]interface{}{"key": shard}

	r := c.c.NewRequest("PUT", "/sys/unseal")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SealStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// Structures for the requests/resposne are all down here. They aren't
// individually documentd because the map almost directly to the raw HTTP API
// documentation. Please refer to that documentation for more details.

type SealStatusResponse struct {
	Sealed   bool
	T        int
	N        int
	Progress int
}
