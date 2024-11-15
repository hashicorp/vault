package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

type DomainRelationships struct {
	Organization        V3ToOneRelationship   `json:"organization"`
	SharedOrganizations V3ToManyRelationships `json:"shared_organizations"`
}

type V3Domain struct {
	Guid          string              `json:"guid"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Name          string              `json:"name"`
	Internal      bool                `json:"internal"`
	Metadata      Metadata            `json:"metadata"`
	Relationships DomainRelationships `json:"relationships"`
	Links         map[string]Link     `json:"links"`
}

type listV3DomainsResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []V3Domain `json:"resources,omitempty"`
}

func (c *Client) ListV3Domains(query url.Values) ([]V3Domain, error) {
	var domains []V3Domain
	requestURL := "/v3/domains"
	if e := query.Encode(); len(e) > 0 {
		requestURL += "?" + e
	}

	for {
		resp, err := c.DoRequest(c.NewRequest("GET", requestURL))
		if err != nil {
			return nil, errors.Wrapf(err, "Error getting domains")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 app domains, response code: %d", resp.StatusCode)
		}

		var data listV3DomainsResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 app domains")
		}

		domains = append(domains, data.Resources...)
		requestURL = data.Pagination.Next.Href
		if requestURL == "" || query.Get("page") != "" {
			break
		}
		requestURL, err = extractPathFromURL(requestURL)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing the next page request url for v3 domains")
		}
	}
	return domains, nil
}
