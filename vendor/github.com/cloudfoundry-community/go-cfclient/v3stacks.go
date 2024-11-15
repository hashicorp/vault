package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// V3Stack implements stack object. Stacks are the base operating system and file system that your application will execute in. A stack is how you configure applications to run against different operating systems (like Windows or Linux) and different versions of those operating systems.
type V3Stack struct {
	Name        string          `json:"name,omitempty"`
	GUID        string          `json:"guid,omitempty"`
	CreatedAt   string          `json:"created_at,omitempty"`
	UpdatedAt   string          `json:"updated_at,omitempty"`
	Description string          `json:"description,omitempty"`
	Links       map[string]Link `json:"links,omitempty"`
	Metadata    V3Metadata      `json:"metadata,omitempty"`
}

type listV3StacksResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []V3Stack  `json:"resources,omitempty"`
}

// ListV3StacksByQuery retrieves stacks based on query
func (c *Client) ListV3StacksByQuery(query url.Values) ([]V3Stack, error) {
	var stacks []V3Stack
	requestURL, err := url.Parse("/v3/stacks")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 stacks")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 stacks, response code: %d", resp.StatusCode)
		}

		var data listV3StacksResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 stacks")
		}

		stacks = append(stacks, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return stacks, nil
}
