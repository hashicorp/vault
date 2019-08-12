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

type OrgQuotasResponse struct {
	Count     int                 `json:"total_results"`
	Pages     int                 `json:"total_pages"`
	NextUrl   string              `json:"next_url"`
	Resources []OrgQuotasResource `json:"resources"`
}

type OrgQuotasResource struct {
	Meta   Meta     `json:"metadata"`
	Entity OrgQuota `json:"entity"`
}

type OrgQuotaRequest struct {
	Name                    string `json:"name"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	TotalServices           int    `json:"total_services"`
	TotalRoutes             int    `json:"total_routes"`
	TotalPrivateDomains     int    `json:"total_private_domains"`
	MemoryLimit             int    `json:"memory_limit"`
	TrialDBAllowed          bool   `json:"trial_db_allowed"`
	InstanceMemoryLimit     int    `json:"instance_memory_limit"`
	AppInstanceLimit        int    `json:"app_instance_limit"`
	AppTaskLimit            int    `json:"app_task_limit"`
	TotalServiceKeys        int    `json:"total_service_keys"`
	TotalReservedRoutePorts int    `json:"total_reserved_route_ports"`
}

type OrgQuota struct {
	Guid                    string `json:"guid"`
	Name                    string `json:"name"`
	CreatedAt               string `json:"created_at,omitempty"`
	UpdatedAt               string `json:"updated_at,omitempty"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	TotalServices           int    `json:"total_services"`
	TotalRoutes             int    `json:"total_routes"`
	TotalPrivateDomains     int    `json:"total_private_domains"`
	MemoryLimit             int    `json:"memory_limit"`
	TrialDBAllowed          bool   `json:"trial_db_allowed"`
	InstanceMemoryLimit     int    `json:"instance_memory_limit"`
	AppInstanceLimit        int    `json:"app_instance_limit"`
	AppTaskLimit            int    `json:"app_task_limit"`
	TotalServiceKeys        int    `json:"total_service_keys"`
	TotalReservedRoutePorts int    `json:"total_reserved_route_ports"`
	c                       *Client
}

func (c *Client) ListOrgQuotasByQuery(query url.Values) ([]OrgQuota, error) {
	var orgQuotas []OrgQuota
	requestUrl := "/v2/quota_definitions?" + query.Encode()
	for {
		orgQuotasResp, err := c.getOrgQuotasResponse(requestUrl)
		if err != nil {
			return []OrgQuota{}, err
		}
		for _, org := range orgQuotasResp.Resources {
			org.Entity.Guid = org.Meta.Guid
			org.Entity.CreatedAt = org.Meta.CreatedAt
			org.Entity.UpdatedAt = org.Meta.UpdatedAt
			org.Entity.c = c
			orgQuotas = append(orgQuotas, org.Entity)
		}
		requestUrl = orgQuotasResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return orgQuotas, nil
}

func (c *Client) ListOrgQuotas() ([]OrgQuota, error) {
	return c.ListOrgQuotasByQuery(nil)
}

func (c *Client) GetOrgQuotaByName(name string) (OrgQuota, error) {
	q := url.Values{}
	q.Set("q", "name:"+name)
	orgQuotas, err := c.ListOrgQuotasByQuery(q)
	if err != nil {
		return OrgQuota{}, err
	}
	if len(orgQuotas) != 1 {
		return OrgQuota{}, fmt.Errorf("Unable to find org quota " + name)
	}
	return orgQuotas[0], nil
}

func (c *Client) getOrgQuotasResponse(requestUrl string) (OrgQuotasResponse, error) {
	var orgQuotasResp OrgQuotasResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return OrgQuotasResponse{}, errors.Wrap(err, "Error requesting org quotas")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return OrgQuotasResponse{}, errors.Wrap(err, "Error reading org quotas body")
	}
	err = json.Unmarshal(resBody, &orgQuotasResp)
	if err != nil {
		return OrgQuotasResponse{}, errors.Wrap(err, "Error unmarshalling org quotas")
	}
	return orgQuotasResp, nil
}

func (c *Client) CreateOrgQuota(orgQuote OrgQuotaRequest) (*OrgQuota, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(orgQuote)
	if err != nil {
		return nil, err
	}
	r := c.NewRequestWithBody("POST", "/v2/quota_definitions", buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return c.handleOrgQuotaResp(resp)
}

func (c *Client) UpdateOrgQuota(orgQuotaGUID string, orgQuota OrgQuotaRequest) (*OrgQuota, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(orgQuota)
	if err != nil {
		return nil, err
	}
	r := c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/quota_definitions/%s", orgQuotaGUID), buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return c.handleOrgQuotaResp(resp)
}

func (c *Client) DeleteOrgQuota(guid string, async bool) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/quota_definitions/%s?async=%t", guid, async)))
	if err != nil {
		return err
	}
	if (async && (resp.StatusCode != http.StatusAccepted)) || (!async && (resp.StatusCode != http.StatusNoContent)) {
		return errors.Wrapf(err, "Error deleting organization %s, response code: %d", guid, resp.StatusCode)
	}
	return nil
}

func (c *Client) handleOrgQuotaResp(resp *http.Response) (*OrgQuota, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var orgQuotasResource OrgQuotasResource
	err = json.Unmarshal(body, &orgQuotasResource)
	if err != nil {
		return nil, err
	}
	return c.mergeOrgQuotaResource(orgQuotasResource), nil
}

func (c *Client) mergeOrgQuotaResource(orgQuotaResource OrgQuotasResource) *OrgQuota {
	orgQuotaResource.Entity.Guid = orgQuotaResource.Meta.Guid
	orgQuotaResource.Entity.CreatedAt = orgQuotaResource.Meta.CreatedAt
	orgQuotaResource.Entity.UpdatedAt = orgQuotaResource.Meta.UpdatedAt
	orgQuotaResource.Entity.c = c
	return &orgQuotaResource.Entity
}
