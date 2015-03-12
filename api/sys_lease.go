package api

func (c *Sys) Renew(id string) (*Secret, error) {
	r := c.c.NewRequest("PUT", "/sys/renew/"+id)
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: parse secret

	return nil, nil
}

func (c *Sys) Revoke(id string) error {
	r := c.c.NewRequest("PUT", "/sys/revoke/"+id)
	resp, err := c.c.RawRequest(r)
	defer resp.Body.Close()
	return err
}
