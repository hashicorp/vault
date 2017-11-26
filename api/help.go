package api

import (
	"context"
	"fmt"
	"net/http"
)

// Help returns help information about the given path.
func (c *Client) Help(path string) (*Help, error) {
	return c.HelpWithContext(context.Background(), path)
}

// HelpWithContext returns help information about the given path, with a context.
func (c *Client) HelpWithContext(ctx context.Context, path string) (*Help, error) {
	req := c.NewRequest(http.MethodGet, fmt.Sprintf("/v1/%s", path))
	req = req.WithContext(ctx)

	req.Params.Add("help", "1")

	resp, err := c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Help
	err = resp.DecodeJSON(&result)
	return &result, err
}

// Help is the response from a help request.
type Help struct {
	// Help is the raw help.
	Help string `json:"help"`

	// SeeAlso is a list of other methods to see for help.
	SeeAlso []string `json:"see_also"`
}
