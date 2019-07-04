package cfclient

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type ServicePlanVisibilitiesResponse struct {
	Count     int                             `json:"total_results"`
	Pages     int                             `json:"total_pages"`
	NextUrl   string                          `json:"next_url"`
	Resources []ServicePlanVisibilityResource `json:"resources"`
}

type ServicePlanVisibilityResource struct {
	Meta   Meta                  `json:"metadata"`
	Entity ServicePlanVisibility `json:"entity"`
}

type ServicePlanVisibility struct {
	Guid             string `json:"guid"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	ServicePlanGuid  string `json:"service_plan_guid"`
	OrganizationGuid string `json:"organization_guid"`
	ServicePlanUrl   string `json:"service_plan_url"`
	OrganizationUrl  string `json:"organization_url"`
	c                *Client
}

func (c *Client) ListServicePlanVisibilitiesByQuery(query url.Values) ([]ServicePlanVisibility, error) {
	var servicePlanVisibilities []ServicePlanVisibility
	requestUrl := "/v2/service_plan_visibilities?" + query.Encode()
	for {
		var servicePlanVisibilitiesResp ServicePlanVisibilitiesResponse
		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting service plan visibilities")
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading service plan visibilities request:")
		}

		err = json.Unmarshal(resBody, &servicePlanVisibilitiesResp)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshaling service plan visibilities")
		}
		for _, servicePlanVisibility := range servicePlanVisibilitiesResp.Resources {
			servicePlanVisibility.Entity.Guid = servicePlanVisibility.Meta.Guid
			servicePlanVisibility.Entity.CreatedAt = servicePlanVisibility.Meta.CreatedAt
			servicePlanVisibility.Entity.UpdatedAt = servicePlanVisibility.Meta.UpdatedAt
			servicePlanVisibility.Entity.c = c
			servicePlanVisibilities = append(servicePlanVisibilities, servicePlanVisibility.Entity)
		}
		requestUrl = servicePlanVisibilitiesResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return servicePlanVisibilities, nil
}

func (c *Client) ListServicePlanVisibilities() ([]ServicePlanVisibility, error) {
	return c.ListServicePlanVisibilitiesByQuery(nil)
}

func (c *Client) GetServicePlanVisibilityByGuid(guid string) (ServicePlanVisibility, error) {
	r := c.NewRequest("GET", "/v2/service_plan_visibilities/"+guid)
	resp, err := c.DoRequest(r)
	if err != nil {
		return ServicePlanVisibility{}, err
	}
	return respBodyToServicePlanVisibility(resp.Body, c)
}

//a uniqueID is the id of the service in the catalog and not in cf internal db
func (c *Client) CreateServicePlanVisibilityByUniqueId(uniqueId string, organizationGuid string) (ServicePlanVisibility, error) {
	q := url.Values{}
	q.Set("q", fmt.Sprintf("unique_id:%s", uniqueId))
	plans, err := c.ListServicePlansByQuery(q)
	if err != nil {
		return ServicePlanVisibility{}, errors.Wrap(err, fmt.Sprintf("Couldn't find a service plan with unique_id: %s", uniqueId))
	}
	return c.CreateServicePlanVisibility(plans[0].Guid, organizationGuid)
}

func (c *Client) CreateServicePlanVisibility(servicePlanGuid string, organizationGuid string) (ServicePlanVisibility, error) {
	req := c.NewRequest("POST", "/v2/service_plan_visibilities")
	req.obj = map[string]interface{}{
		"service_plan_guid": servicePlanGuid,
		"organization_guid": organizationGuid,
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return ServicePlanVisibility{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return ServicePlanVisibility{}, errors.Wrapf(err, "Error creating service plan visibility, response code: %d", resp.StatusCode)
	}
	return respBodyToServicePlanVisibility(resp.Body, c)
}

func (c *Client) DeleteServicePlanVisibilityByPlanAndOrg(servicePlanGuid string, organizationGuid string, async bool) error {
	q := url.Values{}
	q.Set("q", fmt.Sprintf("organization_guid:%s;service_plan_guid:%s", organizationGuid, servicePlanGuid))
	plans, err := c.ListServicePlanVisibilitiesByQuery(q)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Couldn't find a service plan visibility for service plan %s and org %s", servicePlanGuid, organizationGuid))
	}
	if len(plans) != 1 {
		return fmt.Errorf("Query for a service plan visibility did not return exactly one result when searching for a service plan visibility for service plan %s and org %s",
			servicePlanGuid, organizationGuid)
	}
	return c.DeleteServicePlanVisibility(plans[0].Guid, async)
}

func (c *Client) DeleteServicePlanVisibility(guid string, async bool) error {
	req := c.NewRequest("DELETE", fmt.Sprintf("/v2/service_plan_visibilities/%s?async=%v", guid, async))
	resp, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleting service plan visibility, response code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) UpdateServicePlanVisibility(guid string, servicePlanGuid string, organizationGuid string) (ServicePlanVisibility, error) {
	req := c.NewRequest("PUT", "/v2/service_plan_visibilities/"+guid)
	req.obj = map[string]interface{}{
		"service_plan_guid": servicePlanGuid,
		"organization_guid": organizationGuid,
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return ServicePlanVisibility{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return ServicePlanVisibility{}, errors.Wrapf(err, "Error updating service plan visibility, response code: %d", resp.StatusCode)
	}
	return respBodyToServicePlanVisibility(resp.Body, c)
}

func respBodyToServicePlanVisibility(body io.ReadCloser, c *Client) (ServicePlanVisibility, error) {
	bodyRaw, err := ioutil.ReadAll(body)
	if err != nil {
		return ServicePlanVisibility{}, err
	}
	servicePlanVisibilityRes := ServicePlanVisibilityResource{}
	err = json.Unmarshal(bodyRaw, &servicePlanVisibilityRes)
	if err != nil {
		return ServicePlanVisibility{}, err
	}
	servicePlanVisibility := servicePlanVisibilityRes.Entity
	servicePlanVisibility.Guid = servicePlanVisibilityRes.Meta.Guid
	servicePlanVisibility.CreatedAt = servicePlanVisibilityRes.Meta.CreatedAt
	servicePlanVisibility.UpdatedAt = servicePlanVisibilityRes.Meta.UpdatedAt
	servicePlanVisibility.c = c
	return servicePlanVisibility, nil
}
