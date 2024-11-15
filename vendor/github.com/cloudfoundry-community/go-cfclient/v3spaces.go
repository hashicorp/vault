package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type V3Space struct {
	Name          string                         `json:"name,omitempty"`
	GUID          string                         `json:"guid,omitempty"`
	CreatedAt     string                         `json:"created_at,omitempty"`
	UpdatedAt     string                         `json:"updated_at,omitempty"`
	Relationships map[string]V3ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link                `json:"links,omitempty"`
	Metadata      V3Metadata                     `json:"metadata,omitempty"`
}

type CreateV3SpaceRequest struct {
	Name     string
	OrgGUID  string
	Metadata *V3Metadata
}

type UpdateV3SpaceRequest struct {
	Name     string
	Metadata *V3Metadata
}

type V3SpaceUsers struct {
	Name          string                         `json:"name,omitempty"`
	GUID          string                         `json:"guid,omitempty"`
	CreatedAt     string                         `json:"created_at,omitempty"`
	UpdatedAt     string                         `json:"updated_at,omitempty"`
	Relationships map[string]V3ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link                `json:"links,omitempty"`
	Metadata      V3Metadata                     `json:"metadata,omitempty"`
}

func (c *Client) CreateV3Space(r CreateV3SpaceRequest) (*V3Space, error) {
	req := c.NewRequest("POST", "/v3/spaces")
	params := map[string]interface{}{
		"name": r.Name,
		"relationships": map[string]interface{}{
			"organization": V3ToOneRelationship{
				Data: V3Relationship{
					GUID: r.OrgGUID,
				},
			},
		},
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}

	req.obj = params
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating v3 space")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating v3 space %s, response code: %d", r.Name, resp.StatusCode)
	}

	var space V3Space
	if err := json.NewDecoder(resp.Body).Decode(&space); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 space JSON")
	}

	return &space, nil
}

func (c *Client) GetV3SpaceByGUID(spaceGUID string) (*V3Space, error) {
	req := c.NewRequest("GET", "/v3/spaces/"+spaceGUID)

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting v3 space")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting v3 space with GUID [%s], response code: %d", spaceGUID, resp.StatusCode)
	}

	var space V3Space
	if err := json.NewDecoder(resp.Body).Decode(&space); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 space JSON")
	}

	return &space, nil
}

func (c *Client) DeleteV3Space(spaceGUID string) error {
	req := c.NewRequest("DELETE", "/v3/spaces/"+spaceGUID)
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting v3 space")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting v3 space with GUID [%s], response code: %d", spaceGUID, resp.StatusCode)
	}

	return nil
}

func (c *Client) UpdateV3Space(spaceGUID string, r UpdateV3SpaceRequest) (*V3Space, error) {
	req := c.NewRequest("PATCH", "/v3/spaces/"+spaceGUID)
	params := make(map[string]interface{})
	if r.Name != "" {
		params["name"] = r.Name
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}
	if len(params) > 0 {
		req.obj = params
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating v3 space")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error updating v3 space %s, response code: %d", spaceGUID, resp.StatusCode)
	}

	var space V3Space
	if err := json.NewDecoder(resp.Body).Decode(&space); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 space JSON")
	}

	return &space, nil
}

type listV3SpacesResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []V3Space  `json:"resources,omitempty"`
}

func (c *Client) ListV3SpacesByQuery(query url.Values) ([]V3Space, error) {
	var spaces []V3Space
	requestURL := "/v3/spaces"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 spaces")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 spaces, response code: %d", resp.StatusCode)
		}

		var data listV3SpacesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 spaces")
		}

		spaces = append(spaces, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 spaces")
		}
	}

	return spaces, nil
}

type listV3SpaceUsersResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []V3User   `json:"resources,omitempty"`
}

// ListV3SpaceUsers lists users by space GUID
func (c *Client) ListV3SpaceUsers(spaceGUID string) ([]V3User, error) {
	var users []V3User
	requestURL := "/v3/spaces/" + spaceGUID + "/users"

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 space users")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 space users, response code: %d", resp.StatusCode)
		}

		var data listV3SpaceUsersResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 space users")
		}
		users = append(users, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 space users")
		}
	}

	return users, nil
}
