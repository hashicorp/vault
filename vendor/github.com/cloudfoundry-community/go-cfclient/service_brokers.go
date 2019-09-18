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

type ServiceBrokerResponse struct {
	Count     int                     `json:"total_results"`
	Pages     int                     `json:"total_pages"`
	NextUrl   string                  `json:"next_url"`
	Resources []ServiceBrokerResource `json:"resources"`
}

type ServiceBrokerResource struct {
	Meta   Meta          `json:"metadata"`
	Entity ServiceBroker `json:"entity"`
}

type UpdateServiceBrokerRequest struct {
	Name      string `json:"name"`
	BrokerURL string `json:"broker_url"`
	Username  string `json:"auth_username"`
	Password  string `json:"auth_password"`
}

type CreateServiceBrokerRequest struct {
	Name      string `json:"name"`
	BrokerURL string `json:"broker_url"`
	Username  string `json:"auth_username"`
	Password  string `json:"auth_password"`
	SpaceGUID string `json:"space_guid,omitempty"`
}

type ServiceBroker struct {
	Guid      string `json:"guid"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	BrokerURL string `json:"broker_url"`
	Username  string `json:"auth_username"`
	Password  string `json:"auth_password"`
	SpaceGUID string `json:"space_guid,omitempty"`
	c         *Client
}

func (c *Client) DeleteServiceBroker(guid string) error {
	requestUrl := fmt.Sprintf("/v2/service_brokers/%s", guid)
	r := c.NewRequest("DELETE", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleteing service broker %s, response code: %d", guid, resp.StatusCode)
	}
	return nil

}

func (c *Client) UpdateServiceBroker(guid string, usb UpdateServiceBrokerRequest) (ServiceBroker, error) {
	var serviceBrokerResource ServiceBrokerResource

	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(usb)
	if err != nil {
		return ServiceBroker{}, err
	}
	req := c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/service_brokers/%s", guid), buf)
	resp, err := c.DoRequest(req)
	if err != nil {
		return ServiceBroker{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return ServiceBroker{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ServiceBroker{}, err
	}
	err = json.Unmarshal(body, &serviceBrokerResource)
	if err != nil {
		return ServiceBroker{}, err
	}
	serviceBrokerResource.Entity.Guid = serviceBrokerResource.Meta.Guid
	return serviceBrokerResource.Entity, nil
}

func (c *Client) CreateServiceBroker(csb CreateServiceBrokerRequest) (ServiceBroker, error) {
	var serviceBrokerResource ServiceBrokerResource

	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(csb)
	if err != nil {
		return ServiceBroker{}, err
	}
	req := c.NewRequestWithBody("POST", "/v2/service_brokers", buf)
	resp, err := c.DoRequest(req)
	if err != nil {
		return ServiceBroker{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return ServiceBroker{}, fmt.Errorf("CF API returned with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ServiceBroker{}, err
	}
	err = json.Unmarshal(body, &serviceBrokerResource)
	if err != nil {
		return ServiceBroker{}, err
	}

	serviceBrokerResource.Entity.Guid = serviceBrokerResource.Meta.Guid
	return serviceBrokerResource.Entity, nil
}

func (c *Client) ListServiceBrokersByQuery(query url.Values) ([]ServiceBroker, error) {
	var sbs []ServiceBroker
	requestUrl := "/v2/service_brokers?" + query.Encode()
	for {
		serviceBrokerResp, err := c.getServiceBrokerResponse(requestUrl)
		if err != nil {
			return []ServiceBroker{}, err
		}
		for _, sb := range serviceBrokerResp.Resources {
			sb.Entity.Guid = sb.Meta.Guid
			sb.Entity.CreatedAt = sb.Meta.CreatedAt
			sb.Entity.UpdatedAt = sb.Meta.UpdatedAt
			sbs = append(sbs, sb.Entity)
		}
		requestUrl = serviceBrokerResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return sbs, nil
}

func (c *Client) ListServiceBrokers() ([]ServiceBroker, error) {
	return c.ListServiceBrokersByQuery(nil)
}

func (c *Client) GetServiceBrokerByGuid(guid string) (ServiceBroker, error) {
	var serviceBrokerRes ServiceBrokerResource
	r := c.NewRequest("GET", "/v2/service_brokers/"+guid)
	resp, err := c.DoRequest(r)
	if err != nil {
		return ServiceBroker{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ServiceBroker{}, err
	}
	err = json.Unmarshal(body, &serviceBrokerRes)
	if err != nil {
		return ServiceBroker{}, err
	}
	serviceBrokerRes.Entity.Guid = serviceBrokerRes.Meta.Guid
	serviceBrokerRes.Entity.CreatedAt = serviceBrokerRes.Meta.CreatedAt
	serviceBrokerRes.Entity.UpdatedAt = serviceBrokerRes.Meta.UpdatedAt
	return serviceBrokerRes.Entity, nil
}

func (c *Client) GetServiceBrokerByName(name string) (ServiceBroker, error) {
	var sb ServiceBroker
	q := url.Values{}
	q.Set("q", "name:"+name)
	sbs, err := c.ListServiceBrokersByQuery(q)
	if err != nil {
		return sb, err
	}
	if len(sbs) == 0 {
		return sb, fmt.Errorf("Unable to find service broker %s", name)
	}
	return sbs[0], nil
}

func (c *Client) getServiceBrokerResponse(requestUrl string) (ServiceBrokerResponse, error) {
	var serviceBrokerResp ServiceBrokerResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return ServiceBrokerResponse{}, errors.Wrap(err, "Error requesting Service Brokers")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ServiceBrokerResponse{}, errors.Wrap(err, "Error reading Service Broker request")
	}
	err = json.Unmarshal(resBody, &serviceBrokerResp)
	if err != nil {
		return ServiceBrokerResponse{}, errors.Wrap(err, "Error unmarshalling Service Broker")
	}
	return serviceBrokerResp, nil
}
