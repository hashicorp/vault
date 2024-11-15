package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type V3Organization struct {
	Name          string                         `json:"name,omitempty"`
	GUID          string                         `json:"guid,omitempty"`
	Suspended     *bool                          `json:"suspended,omitempty"`
	CreatedAt     string                         `json:"created_at,omitempty"`
	UpdatedAt     string                         `json:"updated_at,omitempty"`
	Relationships map[string]V3ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link                `json:"links,omitempty"`
	Metadata      V3Metadata                     `json:"metadata,omitempty"`
}

type CreateV3OrganizationRequest struct {
	Name      string
	Suspended *bool `json:"suspended,omitempty"`
	Metadata  *V3Metadata
}

type UpdateV3OrganizationRequest struct {
	Name      string
	Suspended *bool `json:"suspended,omitempty"`
	Metadata  *V3Metadata
}

func (c *Client) CreateV3Organization(r CreateV3OrganizationRequest) (*V3Organization, error) {
	req := c.NewRequest("POST", "/v3/organizations")
	params := map[string]interface{}{
		"name": r.Name,
	}
	if r.Suspended != nil {
		params["suspended"] = r.Suspended
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}

	req.obj = params
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating v3 organization")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating v3 organization %s, response code: %d", r.Name, resp.StatusCode)
	}

	var organization V3Organization
	if err := json.NewDecoder(resp.Body).Decode(&organization); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 organization JSON")
	}

	return &organization, nil
}

func (c *Client) GetV3OrganizationByGUID(organizationGUID string) (*V3Organization, error) {
	req := c.NewRequest("GET", "/v3/organizations/"+organizationGUID)

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting v3 organization")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting v3 organization with GUID [%s], response code: %d", organizationGUID, resp.StatusCode)
	}

	var organization V3Organization
	if err := json.NewDecoder(resp.Body).Decode(&organization); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 organization JSON")
	}

	return &organization, nil
}

func (c *Client) DeleteV3Organization(organizationGUID string) error {
	req := c.NewRequest("DELETE", "/v3/organizations/"+organizationGUID)
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting v3 organization")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting v3 organization with GUID [%s], response code: %d", organizationGUID, resp.StatusCode)
	}

	return nil
}

func (c *Client) UpdateV3Organization(organizationGUID string, r UpdateV3OrganizationRequest) (*V3Organization, error) {
	req := c.NewRequest("PATCH", "/v3/organizations/"+organizationGUID)
	params := make(map[string]interface{})
	if r.Name != "" {
		params["name"] = r.Name
	}
	if r.Suspended != nil {
		params["suspended"] = r.Suspended
	}
	if r.Metadata != nil {
		params["metadata"] = r.Metadata
	}
	if len(params) > 0 {
		req.obj = params
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating v3 organization")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error updating v3 organization %s, response code: %d", organizationGUID, resp.StatusCode)
	}

	var organization V3Organization
	if err := json.NewDecoder(resp.Body).Decode(&organization); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 organization JSON")
	}

	return &organization, nil
}

type listV3OrganizationsResponse struct {
	Pagination Pagination       `json:"pagination,omitempty"`
	Resources  []V3Organization `json:"resources,omitempty"`
}

func (c *Client) ListV3OrganizationsByQuery(query url.Values) ([]V3Organization, error) {
	var organizations []V3Organization
	requestURL := "/v3/organizations"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 organizations")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 organizations, response code: %d", resp.StatusCode)
		}

		var data listV3OrganizationsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 organizations")
		}

		organizations = append(organizations, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 organizations")
		}
	}

	return organizations, nil
}
