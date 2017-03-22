package api

func (c *Sys) GenerateShareStatus() (statusResponse *GenerateShareStatusResponse, err error) {
	r := c.c.NewRequest("GET", "/v1/sys/generate-share/attempt")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result GenerateShareStatusResponse
	err = resp.DecodeJSON(&result)
	if err == nil {
		statusResponse = &result
	}
	return
}

func (c *Sys) GenerateShareInit(pgpKey string) (statusResponse *GenerateShareStatusResponse, err error) {
	body := map[string]interface{}{
		"pgp_key": pgpKey,
	}

	r := c.c.NewRequest("PUT", "/v1/sys/generate-share/attempt")
	if err = r.SetJSONBody(body); err != nil {
		return
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result GenerateShareStatusResponse
	err = resp.DecodeJSON(&result)
	if err == nil {
		statusResponse = &result
	}
	return
}

func (c *Sys) GenerateShareCancel() (err error) {
	r := c.c.NewRequest("DELETE", "/v1/sys/generate-share/attempt")
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return
}

func (c *Sys) GenerateShareUpdate(key string) (statusResponse *GenerateShareStatusResponse, err error) {
	body := map[string]interface{}{
		"key": key,
	}

	r := c.c.NewRequest("PUT", "/v1/sys/generate-share/update")
	if err = r.SetJSONBody(body); err != nil {
		return
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result GenerateShareStatusResponse
	err = resp.DecodeJSON(&result)
	if err == nil {
		statusResponse = &result
	}
	return
}

type GenerateShareStatusResponse struct {
	Started        bool
	Progress       int
	Required       int
	Complete       bool
	Key            string `json:"key"`
	PGPFingerprint string `json:"pgp_fingerprint"`
}
