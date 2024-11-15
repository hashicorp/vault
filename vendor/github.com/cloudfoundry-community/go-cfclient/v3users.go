package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// V3User implements the user object
type V3User struct {
	GUID             string          `json:"guid,omitempty"`
	CreatedAt        string          `json:"created_at,omitempty"`
	UpdatedAt        string          `json:"updated_at,omitempty"`
	Username         string          `json:"username,omitempty"`
	PresentationName string          `json:"presentation_name,omitempty"`
	Origin           string          `json:"origin,omitempty"`
	Links            map[string]Link `json:"links,omitempty"`
	Metadata         V3Metadata      `json:"metadata,omitempty"`
}

type listV3UsersResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []V3User   `json:"resources,omitempty"`
}

// ListV3UsersByQuery by query
func (c *Client) ListV3UsersByQuery(query url.Values) ([]V3User, error) {
	var users []V3User
	requestURL, err := url.Parse("/v3/users")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 users")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 users, response code: %d", resp.StatusCode)
		}

		var data listV3UsersResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 users")
		}

		users = append(users, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return users, nil
}
