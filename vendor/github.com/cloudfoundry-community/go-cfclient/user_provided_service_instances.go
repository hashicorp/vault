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

type UserProvidedServiceInstancesResponse struct {
	Count     int                                   `json:"total_results"`
	Pages     int                                   `json:"total_pages"`
	NextUrl   string                                `json:"next_url"`
	Resources []UserProvidedServiceInstanceResource `json:"resources"`
}

type UserProvidedServiceInstanceResource struct {
	Meta   Meta                        `json:"metadata"`
	Entity UserProvidedServiceInstance `json:"entity"`
}

type UserProvidedServiceInstance struct {
	Guid               string                 `json:"guid"`
	Name               string                 `json:"name"`
	CreatedAt          string                 `json:"created_at"`
	UpdatedAt          string                 `json:"updated_at"`
	Credentials        map[string]interface{} `json:"credentials"`
	SpaceGuid          string                 `json:"space_guid"`
	Type               string                 `json:"type"`
	Tags               []string               `json:"tags"`
	SpaceUrl           string                 `json:"space_url"`
	ServiceBindingsUrl string                 `json:"service_bindings_url"`
	RoutesUrl          string                 `json:"routes_url"`
	RouteServiceUrl    string                 `json:"route_service_url"`
	SyslogDrainUrl     string                 `json:"syslog_drain_url"`
	c                  *Client
}

type UserProvidedServiceInstanceRequest struct {
	Name            string                 `json:"name"`
	Credentials     map[string]interface{} `json:"credentials"`
	SpaceGuid       string                 `json:"space_guid"`
	Tags            []string               `json:"tags"`
	RouteServiceUrl string                 `json:"route_service_url"`
	SyslogDrainUrl  string                 `json:"syslog_drain_url"`
}

func (c *Client) ListUserProvidedServiceInstancesByQuery(query url.Values) ([]UserProvidedServiceInstance, error) {
	var instances []UserProvidedServiceInstance

	requestUrl := "/v2/user_provided_service_instances?" + query.Encode()
	for {
		var sir UserProvidedServiceInstancesResponse
		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting user provided service instances")
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading user provided service instances request:")
		}

		err = json.Unmarshal(resBody, &sir)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshaling user provided service instances")
		}
		for _, instance := range sir.Resources {
			instance.Entity.Guid = instance.Meta.Guid
			instance.Entity.CreatedAt = instance.Meta.CreatedAt
			instance.Entity.UpdatedAt = instance.Meta.UpdatedAt
			instance.Entity.c = c
			instances = append(instances, instance.Entity)
		}

		requestUrl = sir.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return instances, nil
}

func (c *Client) ListUserProvidedServiceInstances() ([]UserProvidedServiceInstance, error) {
	return c.ListUserProvidedServiceInstancesByQuery(nil)
}

func (c *Client) GetUserProvidedServiceInstanceByGuid(guid string) (UserProvidedServiceInstance, error) {
	var sir UserProvidedServiceInstanceResource
	req := c.NewRequest("GET", "/v2/user_provided_service_instances/"+guid)
	res, err := c.DoRequest(req)
	if err != nil {
		return UserProvidedServiceInstance{}, errors.Wrap(err, "Error requesting user provided service instance")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return UserProvidedServiceInstance{}, errors.Wrap(err, "Error reading user provided service instance response")
	}
	err = json.Unmarshal(data, &sir)
	if err != nil {
		return UserProvidedServiceInstance{}, errors.Wrap(err, "Error JSON parsing user provided service instance response")
	}
	sir.Entity.Guid = sir.Meta.Guid
	sir.Entity.CreatedAt = sir.Meta.CreatedAt
	sir.Entity.UpdatedAt = sir.Meta.UpdatedAt
	sir.Entity.c = c
	return sir.Entity, nil
}

func (c *Client) UserProvidedServiceInstanceByGuid(guid string) (UserProvidedServiceInstance, error) {
	return c.GetUserProvidedServiceInstanceByGuid(guid)
}

func (c *Client) CreateUserProvidedServiceInstance(req UserProvidedServiceInstanceRequest) (*UserProvidedServiceInstance, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}
	r := c.NewRequestWithBody("POST", "/v2/user_provided_service_instances", buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}

	return c.handleUserProvidedServiceInstanceResp(resp)
}

func (c *Client) DeleteUserProvidedServiceInstance(guid string) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/user_provided_service_instances/%s", guid)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleting user provided service instance %s, response code %d", guid, resp.StatusCode)
	}
	return nil
}

func (c *Client) UpdateUserProvidedServiceInstance(guid string, req UserProvidedServiceInstanceRequest) (*UserProvidedServiceInstance, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}
	r := c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/user_provided_service_instances/%s", guid), buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}
	return c.handleUserProvidedServiceInstanceResp(resp)
}

func (c *Client) handleUserProvidedServiceInstanceResp(resp *http.Response) (*UserProvidedServiceInstance, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var upsResource UserProvidedServiceInstanceResource
	err = json.Unmarshal(body, &upsResource)
	if err != nil {
		return nil, err
	}
	return c.mergeUserProvidedServiceInstanceResource(upsResource), nil
}

func (c *Client) mergeUserProvidedServiceInstanceResource(ups UserProvidedServiceInstanceResource) *UserProvidedServiceInstance {
	ups.Entity.Guid = ups.Meta.Guid
	ups.Entity.CreatedAt = ups.Meta.CreatedAt
	ups.Entity.UpdatedAt = ups.Meta.UpdatedAt
	ups.Entity.c = c
	return &ups.Entity
}
