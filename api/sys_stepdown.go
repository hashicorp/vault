package api

import "context"

func (c *Sys) StepDown() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.StepDownContext(ctx)
}

func (c *Sys) StepDownContext(ctx context.Context) error {
	r := c.c.NewRequest("PUT", "/v1/sys/step-down")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	return err
}
