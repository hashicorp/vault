package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// V3ServiceCredentialBindings implements the service credential binding object. a credential binding can be a binding between apps and a service instance or a service key
type V3ServiceCredentialBindings struct {
	GUID          string                         `json:"guid"`
	CreatedAt     time.Time                      `json:"created_at"`
	UpdatedAt     time.Time                      `json:"updated_at"`
	Name          string                         `json:"name"`
	Type          string                         `json:"type"`
	LastOperation LastOperation                  `json:"last_operation"`
	Metadata      Metadata                       `json:"metadata"`
	Relationships map[string]V3ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link                `json:"links"`
}

type listV3ServiceCredentialBindingsResponse struct {
	Pagination Pagination                    `json:"pagination,omitempty"`
	Resources  []V3ServiceCredentialBindings `json:"resources,omitempty"`
}

// ListV3ServiceCredentialBindings retrieves all service credential bindings
func (c *Client) ListV3ServiceCredentialBindings() ([]V3ServiceCredentialBindings, error) {
	return c.ListV3ServiceCredentialBindingsByQuery(nil)
}

// ListV3ServiceCredentialBindingsByQuery retrieves service credential bindings using a query
func (c *Client) ListV3ServiceCredentialBindingsByQuery(query url.Values) ([]V3ServiceCredentialBindings, error) {
	var svcCredentialBindings []V3ServiceCredentialBindings
	requestURL := "/v3/service_credential_bindings"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		r := c.NewRequest("GET", requestURL)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 service credential bindings")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error listing v3 service credential bindings, response code: %d", resp.StatusCode)
		}

		var data listV3ServiceCredentialBindingsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 service credential bindings")
		}

		svcCredentialBindings = append(svcCredentialBindings, data.Resources...)

		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 service credential bindings")
		}
	}

	return svcCredentialBindings, nil
}

// GetV3ServiceCredentialBindingsByGUID retrieves the service credential binding based on the provided guid
func (c *Client) GetV3ServiceCredentialBindingsByGUID(GUID string) (*V3ServiceCredentialBindings, error) {
	requestURL := fmt.Sprintf("/v3/service_credential_bindings/%s", GUID)
	req := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(req)

	if err != nil {
		return nil, errors.Wrap(err, "Error while getting v3 service credential binding")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting v3 service credential binding with GUID [%s], response code: %d", GUID, resp.StatusCode)
	}

	var svcCredentialBindings V3ServiceCredentialBindings
	if err := json.NewDecoder(resp.Body).Decode(&svcCredentialBindings); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 service credential binding JSON")
	}

	return &svcCredentialBindings, nil
}
