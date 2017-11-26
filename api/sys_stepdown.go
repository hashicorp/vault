package api

import (
	"context"
	"net/http"
)

// StepDown sends an API call to the client to inform it to give up leadership.
func (c *Sys) StepDown() error {
	return c.StepDownWithContext(context.Background())
}

// StepDown sends an API call to the client to inform it to give up leadership,
// with a context.
func (c *Sys) StepDownWithContext(ctx context.Context) error {
	req := c.c.NewRequest(http.MethodPut, "/v1/sys/step-down")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}
