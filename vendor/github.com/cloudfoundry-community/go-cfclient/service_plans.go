package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/pkg/errors"
)

type ServicePlansResponse struct {
	Count     int                   `json:"total_results"`
	Pages     int                   `json:"total_pages"`
	NextUrl   string                `json:"next_url"`
	Resources []ServicePlanResource `json:"resources"`
}

type ServicePlanResource struct {
	Meta   Meta        `json:"metadata"`
	Entity ServicePlan `json:"entity"`
}

type ServicePlan struct {
	Name                string      `json:"name"`
	Guid                string      `json:"guid"`
	CreatedAt           string      `json:"created_at"`
	UpdatedAt           string      `json:"updated_at"`
	Free                bool        `json:"free"`
	Description         string      `json:"description"`
	ServiceGuid         string      `json:"service_guid"`
	Extra               interface{} `json:"extra"`
	UniqueId            string      `json:"unique_id"`
	Public              bool        `json:"public"`
	Active              bool        `json:"active"`
	Bindable            bool        `json:"bindable"`
	ServiceUrl          string      `json:"service_url"`
	ServiceInstancesUrl string      `json:"service_instances_url"`
	c                   *Client
}

func (c *Client) ListServicePlansByQuery(query url.Values) ([]ServicePlan, error) {
	var servicePlans []ServicePlan
	requestUrl := "/v2/service_plans?" + query.Encode()
	for {
		var servicePlansResp ServicePlansResponse
		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting service plans")
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading service plans request:")
		}
		err = json.Unmarshal(resBody, &servicePlansResp)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshaling service plans")
		}
		for _, servicePlan := range servicePlansResp.Resources {
			servicePlan.Entity.Guid = servicePlan.Meta.Guid
			servicePlan.Entity.CreatedAt = servicePlan.Meta.CreatedAt
			servicePlan.Entity.UpdatedAt = servicePlan.Meta.UpdatedAt
			servicePlan.Entity.c = c
			servicePlans = append(servicePlans, servicePlan.Entity)
		}
		requestUrl = servicePlansResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return servicePlans, nil
}

func (c *Client) ListServicePlans() ([]ServicePlan, error) {
	return c.ListServicePlansByQuery(nil)
}

func (c *Client) GetServicePlanByGUID(guid string) (*ServicePlan, error) {
	var (
		plan         *ServicePlan
		planResponse ServicePlanResource
	)

	r := c.NewRequest("GET", "/v2/service_plans/"+guid)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &planResponse)
	if err != nil {
		return nil, err
	}

	planResponse.Entity.Guid = planResponse.Meta.Guid
	planResponse.Entity.CreatedAt = planResponse.Meta.CreatedAt
	planResponse.Entity.UpdatedAt = planResponse.Meta.UpdatedAt
	plan = &planResponse.Entity

	return plan, nil
}

func (c *Client) MakeServicePlanPublic(servicePlanGUID string) error {
	return c.setPlanGlobalVisibility(servicePlanGUID, true)
}

func (c *Client) MakeServicePlanPrivate(servicePlanGUID string) error {
	return c.setPlanGlobalVisibility(servicePlanGUID, false)
}

func (c *Client) setPlanGlobalVisibility(servicePlanGUID string, public bool) error {
	bodyString := fmt.Sprintf(`{"public": %t}`, public)
	req := c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/service_plans/%s", servicePlanGUID), bytes.NewBufferString(bodyString))

	resp, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
