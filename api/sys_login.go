package api

import (
	"fmt"
)

// Login performs the /sys/login API call.
//
// This API call is stateful: it will set the access token on the client
// for future API calls to be authenticated. The access token can be retrieved
// at any time from the client using `client.Token()` and it can be cleared
// with `sys.Logout()`.
func (c *Sys) Login(vars map[string]string) error {
	r := c.c.NewRequest("PUT", "/v1/sys/login")
	if err := r.SetJSONBody(vars); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.c.Token() == "" {
		return fmt.Errorf(
			"Login had status code %d, but token cookie was not set!",
			resp.StatusCode)
	}

	return nil
}
