package api

func (c *Sys) GenerateShareStatus() (*GenerateShareStatusResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/generate-share/attempt")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenerateShareStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) GenerateShareInit(pgpKey string) (*GenerateShareStatusResponse, error) {
	body := map[string]interface{}{
		"pgp_key": pgpKey,
	}

	r := c.c.NewRequest("PUT", "/v1/sys/generate-share/attempt")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenerateShareStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) GenerateShareCancel() error {
	r := c.c.NewRequest("DELETE", "/v1/sys/generate-share/attempt")
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) GenerateShareUpdate(share string) (*GenerateShareStatusResponse, error) {
	body := map[string]interface{}{
		"key": share,
	}

	r := c.c.NewRequest("PUT", "/v1/sys/generate-share/update")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenerateShareStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

type GenerateShareStatusResponse struct {
	Started        bool
	Progress       int
	Required       int
	Complete       bool
	Key            string `json:"key"`
	PGPFingerprint string `json:"pgp_fingerprint"`
}
