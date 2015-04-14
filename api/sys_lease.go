package api

func (c *Sys) Renew(id string, increment int) (*Secret, error) {
	r := c.c.NewRequest("PUT", "/v1/sys/renew/"+id)

	body := map[string]interface{}{"increment": increment}
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *Sys) Revoke(id string) error {
	r := c.c.NewRequest("PUT", "/v1/sys/revoke/"+id)
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RevokePrefix(id string) error {
	r := c.c.NewRequest("PUT", "/v1/sys/revoke-prefix/"+id)
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}
