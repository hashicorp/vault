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

type ServiceKeysResponse struct {
	Count     int                  `json:"total_results"`
	Pages     int                  `json:"total_pages"`
	Resources []ServiceKeyResource `json:"resources"`
	NextUrl   string               `json:"next_url"`
}

type ServiceKeyResource struct {
	Meta   Meta       `json:"metadata"`
	Entity ServiceKey `json:"entity"`
}

type CreateServiceKeyRequest struct {
	Name                string      `json:"name"`
	ServiceInstanceGuid string      `json:"service_instance_guid"`
	Parameters          interface{} `json:"parameters,omitempty"`
}

type ServiceKey struct {
	Name                string      `json:"name"`
	Guid                string      `json:"guid"`
	CreatedAt           string      `json:"created_at"`
	UpdatedAt           string      `json:"updated_at"`
	ServiceInstanceGuid string      `json:"service_instance_guid"`
	Credentials         interface{} `json:"credentials"`
	ServiceInstanceUrl  string      `json:"service_instance_url"`
	c                   *Client
}

func (c *Client) ListServiceKeysByQuery(query url.Values) ([]ServiceKey, error) {
	var serviceKeys []ServiceKey
	requestUrl := "/v2/service_keys?" + query.Encode()

	for {
		var serviceKeysResp ServiceKeysResponse

		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting service keys")
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading service keys request:")
		}

		err = json.Unmarshal(resBody, &serviceKeysResp)
		if err != nil {
			return nil, errors.Wrapf(err, "Error unmarshaling service keys: %q", string(resBody))
		}
		for _, serviceKey := range serviceKeysResp.Resources {
			serviceKey.Entity.Guid = serviceKey.Meta.Guid
			serviceKey.Entity.CreatedAt = serviceKey.Meta.CreatedAt
			serviceKey.Entity.UpdatedAt = serviceKey.Meta.UpdatedAt
			serviceKey.Entity.c = c
			serviceKeys = append(serviceKeys, serviceKey.Entity)
		}

		requestUrl = serviceKeysResp.NextUrl
		if requestUrl == "" {
			break
		}
	}

	return serviceKeys, nil
}

func (c *Client) ListServiceKeys() ([]ServiceKey, error) {
	return c.ListServiceKeysByQuery(nil)
}

func (c *Client) GetServiceKeyByName(name string) (ServiceKey, error) {
	var serviceKey ServiceKey
	q := url.Values{}
	q.Set("q", "name:"+name)
	serviceKeys, err := c.ListServiceKeysByQuery(q)
	if err != nil {
		return serviceKey, err
	}
	if len(serviceKeys) == 0 {
		return serviceKey, fmt.Errorf("Unable to find service key %s", name)
	}
	return serviceKeys[0], nil
}

// GetServiceKeyByInstanceGuid is deprecated in favor of GetServiceKeysByInstanceGuid
func (c *Client) GetServiceKeyByInstanceGuid(guid string) (ServiceKey, error) {
	var serviceKey ServiceKey
	q := url.Values{}
	q.Set("q", "service_instance_guid:"+guid)
	serviceKeys, err := c.ListServiceKeysByQuery(q)
	if err != nil {
		return serviceKey, err
	}
	if len(serviceKeys) == 0 {
		return serviceKey, fmt.Errorf("Unable to find service key for guid %s", guid)
	}
	return serviceKeys[0], nil
}

// GetServiceKeysByInstanceGuid returns the service keys for a service instance.
// If none are found, it returns an error.
func (c *Client) GetServiceKeysByInstanceGuid(guid string) ([]ServiceKey, error) {
	q := url.Values{}
	q.Set("q", "service_instance_guid:"+guid)
	serviceKeys, err := c.ListServiceKeysByQuery(q)
	if err != nil {
		return serviceKeys, err
	}
	if len(serviceKeys) == 0 {
		return serviceKeys, fmt.Errorf("Unable to find service key for guid %s", guid)
	}
	return serviceKeys, nil
}

// CreateServiceKey creates a service key from the request. If a service key
// exists already, it returns an error containing `CF-ServiceKeyNameTaken`
func (c *Client) CreateServiceKey(csr CreateServiceKeyRequest) (ServiceKey, error) {
	var serviceKeyResource ServiceKeyResource

	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(csr)
	if err != nil {
		return ServiceKey{}, err
	}
	req := c.NewRequestWithBody("POST", "/v2/service_keys", buf)
	resp, err := c.DoRequest(req)
	if err != nil {
		return ServiceKey{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return ServiceKey{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ServiceKey{}, err
	}
	err = json.Unmarshal(body, &serviceKeyResource)
	if err != nil {
		return ServiceKey{}, err
	}

	return serviceKeyResource.Entity, nil
}

// DeleteServiceKey removes a service key instance
func (c *Client) DeleteServiceKey(guid string) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/service_keys/%s", guid)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleting service instance key %s, response code %d", guid, resp.StatusCode)
	}
	return nil
}
