package cfclient

import (
	"encoding/json"
	"io/ioutil"
	"net/url"

	"github.com/pkg/errors"
)

type ServicesResponse struct {
	Count     int                `json:"total_results"`
	Pages     int                `json:"total_pages"`
	NextUrl   string             `json:"next_url"`
	Resources []ServicesResource `json:"resources"`
}

type ServicesResource struct {
	Meta   Meta    `json:"metadata"`
	Entity Service `json:"entity"`
}

type Service struct {
	Guid                 string   `json:"guid"`
	Label                string   `json:"label"`
	CreatedAt            string   `json:"created_at"`
	UpdatedAt            string   `json:"updated_at"`
	Description          string   `json:"description"`
	Active               bool     `json:"active"`
	Bindable             bool     `json:"bindable"`
	ServiceBrokerGuid    string   `json:"service_broker_guid"`
	PlanUpdateable       bool     `json:"plan_updateable"`
	Tags                 []string `json:"tags"`
	UniqueID             string   `json:"unique_id"`
	Extra                string   `json:"extra"`
	Requires             []string `json:"requires"`
	InstancesRetrievable bool     `json:"instances_retrievable"`
	BindingsRetrievable  bool     `json:"bindings_retrievable"`
	c                    *Client
}

type ServiceSummary struct {
	Guid          string `json:"guid"`
	Name          string `json:"name"`
	BoundAppCount int    `json:"bound_app_count"`
}

func (c *Client) GetServiceByGuid(guid string) (Service, error) {
	var serviceRes ServicesResource
	r := c.NewRequest("GET", "/v2/services/"+guid)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Service{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return Service{}, err
	}
	err = json.Unmarshal(body, &serviceRes)
	if err != nil {
		return Service{}, err
	}
	serviceRes.Entity.Guid = serviceRes.Meta.Guid
	serviceRes.Entity.CreatedAt = serviceRes.Meta.CreatedAt
	serviceRes.Entity.UpdatedAt = serviceRes.Meta.UpdatedAt
	return serviceRes.Entity, nil

}

func (c *Client) ListServicesByQuery(query url.Values) ([]Service, error) {
	var services []Service
	requestUrl := "/v2/services?" + query.Encode()
	for {
		var serviceResp ServicesResponse
		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting services")
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading services request:")
		}

		err = json.Unmarshal(resBody, &serviceResp)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshaling services")
		}
		for _, service := range serviceResp.Resources {
			service.Entity.Guid = service.Meta.Guid
			service.Entity.CreatedAt = service.Meta.CreatedAt
			service.Entity.UpdatedAt = service.Meta.UpdatedAt
			service.Entity.c = c
			services = append(services, service.Entity)
		}
		requestUrl = serviceResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return services, nil
}

func (c *Client) ListServices() ([]Service, error) {
	return c.ListServicesByQuery(nil)
}
