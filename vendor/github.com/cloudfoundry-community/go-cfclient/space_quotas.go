package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type SpaceQuotasResponse struct {
	Count     int                   `json:"total_results"`
	Pages     int                   `json:"total_pages"`
	NextUrl   string                `json:"next_url"`
	Resources []SpaceQuotasResource `json:"resources"`
}

type SpaceQuotasResource struct {
	Meta   Meta       `json:"metadata"`
	Entity SpaceQuota `json:"entity"`
}

type SpaceQuotaRequest struct {
	Name                    string `json:"name"`
	OrganizationGuid        string `json:"organization_guid"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	TotalServices           int    `json:"total_services"`
	TotalRoutes             int    `json:"total_routes"`
	MemoryLimit             int    `json:"memory_limit"`
	InstanceMemoryLimit     int    `json:"instance_memory_limit"`
	AppInstanceLimit        int    `json:"app_instance_limit"`
	AppTaskLimit            int    `json:"app_task_limit"`
	TotalServiceKeys        int    `json:"total_service_keys"`
	TotalReservedRoutePorts int    `json:"total_reserved_route_ports"`
}

type SpaceQuota struct {
	Guid                    string `json:"guid"`
	CreatedAt               string `json:"created_at,omitempty"`
	UpdatedAt               string `json:"updated_at,omitempty"`
	Name                    string `json:"name"`
	OrganizationGuid        string `json:"organization_guid"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	TotalServices           int    `json:"total_services"`
	TotalRoutes             int    `json:"total_routes"`
	MemoryLimit             int    `json:"memory_limit"`
	InstanceMemoryLimit     int    `json:"instance_memory_limit"`
	AppInstanceLimit        int    `json:"app_instance_limit"`
	AppTaskLimit            int    `json:"app_task_limit"`
	TotalServiceKeys        int    `json:"total_service_keys"`
	TotalReservedRoutePorts int    `json:"total_reserved_route_ports"`
	c                       *Client
}

func (c *Client) ListSpaceQuotasByQuery(query url.Values) ([]SpaceQuota, error) {
	var spaceQuotas []SpaceQuota
	requestUrl := "/v2/space_quota_definitions?" + query.Encode()
	for {
		spaceQuotasResp, err := c.getSpaceQuotasResponse(requestUrl)
		if err != nil {
			return []SpaceQuota{}, err
		}
		for _, space := range spaceQuotasResp.Resources {
			space.Entity.Guid = space.Meta.Guid
			space.Entity.CreatedAt = space.Meta.CreatedAt
			space.Entity.UpdatedAt = space.Meta.UpdatedAt
			space.Entity.c = c
			spaceQuotas = append(spaceQuotas, space.Entity)
		}
		requestUrl = spaceQuotasResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return spaceQuotas, nil
}

func (c *Client) ListSpaceQuotas() ([]SpaceQuota, error) {
	return c.ListSpaceQuotasByQuery(nil)
}

func (c *Client) GetSpaceQuotaByName(name string) (SpaceQuota, error) {
	q := url.Values{}
	q.Set("q", "name:"+name)
	spaceQuotas, err := c.ListSpaceQuotasByQuery(q)
	if err != nil {
		return SpaceQuota{}, err
	}
	if len(spaceQuotas) != 1 {
		return SpaceQuota{}, fmt.Errorf("Unable to find space quota " + name)
	}
	return spaceQuotas[0], nil
}

func (c *Client) getSpaceQuotasResponse(requestUrl string) (SpaceQuotasResponse, error) {
	var spaceQuotasResp SpaceQuotasResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return SpaceQuotasResponse{}, errors.Wrap(err, "Error requesting space quotas")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return SpaceQuotasResponse{}, errors.Wrap(err, "Error reading space quotas body")
	}
	err = json.Unmarshal(resBody, &spaceQuotasResp)
	if err != nil {
		return SpaceQuotasResponse{}, errors.Wrap(err, "Error unmarshalling space quotas")
	}
	return spaceQuotasResp, nil
}

func (c *Client) AssignSpaceQuota(quotaGUID, spaceGUID string) error {
	//Perform the PUT and check for errors
	resp, err := c.DoRequest(c.NewRequest("PUT", fmt.Sprintf("/v2/space_quota_definitions/%s/spaces/%s", quotaGUID, spaceGUID)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated { //201
		return fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateSpaceQuota(spaceQuote SpaceQuotaRequest) (*SpaceQuota, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(spaceQuote)
	if err != nil {
		return nil, err
	}
	r := c.NewRequestWithBody("POST", "/v2/space_quota_definitions", buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return c.handleSpaceQuotaResp(resp)
}

func (c *Client) UpdateSpaceQuota(spaceQuotaGUID string, spaceQuote SpaceQuotaRequest) (*SpaceQuota, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(spaceQuote)
	if err != nil {
		return nil, err
	}
	r := c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/space_quota_definitions/%s", spaceQuotaGUID), buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return c.handleSpaceQuotaResp(resp)
}

func (c *Client) handleSpaceQuotaResp(resp *http.Response) (*SpaceQuota, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var spaceQuotasResource SpaceQuotasResource
	err = json.Unmarshal(body, &spaceQuotasResource)
	if err != nil {
		return nil, err
	}
	return c.mergeSpaceQuotaResource(spaceQuotasResource), nil
}

func (c *Client) mergeSpaceQuotaResource(spaceQuote SpaceQuotasResource) *SpaceQuota {
	spaceQuote.Entity.Guid = spaceQuote.Meta.Guid
	spaceQuote.Entity.CreatedAt = spaceQuote.Meta.CreatedAt
	spaceQuote.Entity.UpdatedAt = spaceQuote.Meta.UpdatedAt
	spaceQuote.Entity.c = c
	return &spaceQuote.Entity
}
