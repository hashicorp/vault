package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type ServiceBindingsResponse struct {
	Count     int                      `json:"total_results"`
	Pages     int                      `json:"total_pages"`
	Resources []ServiceBindingResource `json:"resources"`
	NextUrl   string                   `json:"next_url"`
}

type ServiceBindingResource struct {
	Meta   Meta           `json:"metadata"`
	Entity ServiceBinding `json:"entity"`
}

type ServiceBinding struct {
	Guid                string      `json:"guid"`
	Name                string      `json:"name"`
	CreatedAt           string      `json:"created_at"`
	UpdatedAt           string      `json:"updated_at"`
	AppGuid             string      `json:"app_guid"`
	ServiceInstanceGuid string      `json:"service_instance_guid"`
	Credentials         interface{} `json:"credentials"`
	BindingOptions      interface{} `json:"binding_options"`
	GatewayData         interface{} `json:"gateway_data"`
	GatewayName         string      `json:"gateway_name"`
	SyslogDrainUrl      string      `json:"syslog_drain_url"`
	VolumeMounts        interface{} `json:"volume_mounts"`
	AppUrl              string      `json:"app_url"`
	ServiceInstanceUrl  string      `json:"service_instance_url"`
	c                   *Client
}

func (c *Client) ListServiceBindingsByQuery(query url.Values) ([]ServiceBinding, error) {
	var serviceBindings []ServiceBinding
	requestUrl := "/v2/service_bindings?" + query.Encode()

	for {
		var serviceBindingsResp ServiceBindingsResponse

		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting service bindings")
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading service bindings request:")
		}

		err = json.Unmarshal(resBody, &serviceBindingsResp)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshaling service bindings")
		}
		for _, serviceBinding := range serviceBindingsResp.Resources {
			serviceBinding.Entity.Guid = serviceBinding.Meta.Guid
			serviceBinding.Entity.CreatedAt = serviceBinding.Meta.CreatedAt
			serviceBinding.Entity.UpdatedAt = serviceBinding.Meta.UpdatedAt
			serviceBinding.Entity.c = c
			serviceBindings = append(serviceBindings, serviceBinding.Entity)
		}
		requestUrl = serviceBindingsResp.NextUrl
		if requestUrl == "" {
			break
		}
	}

	return serviceBindings, nil
}

func (c *Client) ListServiceBindings() ([]ServiceBinding, error) {
	return c.ListServiceBindingsByQuery(nil)
}

func (c *Client) GetServiceBindingByGuid(guid string) (ServiceBinding, error) {
	var serviceBinding ServiceBindingResource
	r := c.NewRequest("GET", "/v2/service_bindings/"+url.QueryEscape(guid))
	resp, err := c.DoRequest(r)
	if err != nil {
		return ServiceBinding{}, errors.Wrap(err, "Error requesting serving binding")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ServiceBinding{}, errors.Wrap(err, "Error reading service binding response body")
	}
	err = json.Unmarshal(resBody, &serviceBinding)
	if err != nil {
		return ServiceBinding{}, errors.Wrap(err, "Error unmarshalling service binding")
	}
	serviceBinding.Entity.Guid = serviceBinding.Meta.Guid
	serviceBinding.Entity.CreatedAt = serviceBinding.Meta.CreatedAt
	serviceBinding.Entity.UpdatedAt = serviceBinding.Meta.UpdatedAt
	serviceBinding.Entity.c = c
	return serviceBinding.Entity, nil
}

func (c *Client) ServiceBindingByGuid(guid string) (ServiceBinding, error) {
	return c.GetServiceBindingByGuid(guid)
}

func (c *Client) DeleteServiceBinding(guid string) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/service_bindings/%s", guid)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleting service binding %s, response code %d", guid, resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateServiceBinding(appGUID, serviceInstanceGUID string) (*ServiceBinding, error) {
	req := c.NewRequest("POST", fmt.Sprintf("/v2/service_bindings"))
	req.obj = map[string]interface{}{
		"app_guid":              appGUID,
		"service_instance_guid": serviceInstanceGUID,
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.Wrapf(err, "Error binding app %s to service instance %s, response code %d", appGUID, serviceInstanceGUID, resp.StatusCode)
	}
	return c.handleServiceBindingResp(resp)
}

func (c *Client) CreateRouteServiceBinding(routeGUID, serviceInstanceGUID string) error {
	req := c.NewRequest("PUT", fmt.Sprintf("/v2/user_provided_service_instances/%s/routes/%s", serviceInstanceGUID, routeGUID))
	resp, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.Wrapf(err, "Error binding route %s to service instance %s, response code %d", routeGUID, serviceInstanceGUID, resp.StatusCode)
	}
	return nil
}

func (c *Client) DeleteRouteServiceBinding(routeGUID, serviceInstanceGUID string) error {
	req := c.NewRequest("DELETE", fmt.Sprintf("/v2/service_instances/%s/routes/%s", serviceInstanceGUID, routeGUID))
	resp, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "Error deleting bound route %s from service instance %s, response code %d", routeGUID, serviceInstanceGUID, resp.StatusCode)
	}
	return nil
}

func (c *Client) handleServiceBindingResp(resp *http.Response) (*ServiceBinding, error) {
	defer resp.Body.Close()
	var sb ServiceBindingResource
	err := json.NewDecoder(resp.Body).Decode(&sb)
	if err != nil {
		return nil, err
	}
	return c.mergeServiceBindingResource(sb), nil
}

func (c *Client) mergeServiceBindingResource(serviceBinding ServiceBindingResource) *ServiceBinding {
	serviceBinding.Entity.Guid = serviceBinding.Meta.Guid
	serviceBinding.Entity.c = c
	return &serviceBinding.Entity
}
