package cfclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type DomainsResponse struct {
	Count     int              `json:"total_results"`
	Pages     int              `json:"total_pages"`
	NextUrl   string           `json:"next_url"`
	Resources []DomainResource `json:"resources"`
}

type SharedDomainsResponse struct {
	Count     int                    `json:"total_results"`
	Pages     int                    `json:"total_pages"`
	NextUrl   string                 `json:"next_url"`
	Resources []SharedDomainResource `json:"resources"`
}

type DomainResource struct {
	Meta   Meta   `json:"metadata"`
	Entity Domain `json:"entity"`
}

type SharedDomainResource struct {
	Meta   Meta         `json:"metadata"`
	Entity SharedDomain `json:"entity"`
}

type Domain struct {
	Guid                   string `json:"guid"`
	Name                   string `json:"name"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
	OwningOrganizationGuid string `json:"owning_organization_guid"`
	OwningOrganizationUrl  string `json:"owning_organization_url"`
	SharedOrganizationsUrl string `json:"shared_organizations_url"`
	c                      *Client
}

type SharedDomain struct {
	Guid            string `json:"guid"`
	Name            string `json:"name"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	RouterGroupGuid string `json:"router_group_guid"`
	RouterGroupType string `json:"router_group_type"`
	Internal        bool   `json:"internal"`
	c               *Client
}

func (c *Client) ListDomainsByQuery(query url.Values) ([]Domain, error) {
	var domains []Domain
	requestUrl := "/v2/private_domains?" + query.Encode()
	for {
		var domainResp DomainsResponse
		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting domains")
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading domains request")
		}

		err = json.Unmarshal(resBody, &domainResp)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshaling domains")
		}
		for _, domain := range domainResp.Resources {
			domain.Entity.Guid = domain.Meta.Guid
			domain.Entity.CreatedAt = domain.Meta.CreatedAt
			domain.Entity.UpdatedAt = domain.Meta.UpdatedAt
			domain.Entity.c = c
			domains = append(domains, domain.Entity)
		}
		requestUrl = domainResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return domains, nil
}

func (c *Client) ListDomains() ([]Domain, error) {
	return c.ListDomainsByQuery(nil)
}

func (c *Client) ListSharedDomainsByQuery(query url.Values) ([]SharedDomain, error) {
	var domains []SharedDomain
	requestUrl := "/v2/shared_domains?" + query.Encode()
	for {
		var domainResp SharedDomainsResponse
		r := c.NewRequest("GET", requestUrl)
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting shared domains")
		}
		resBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading shared domains request")
		}

		err = json.Unmarshal(resBody, &domainResp)
		if err != nil {
			return nil, errors.Wrap(err, "Error unmarshaling shared domains")
		}
		for _, domain := range domainResp.Resources {
			domain.Entity.Guid = domain.Meta.Guid
			domain.Entity.CreatedAt = domain.Meta.CreatedAt
			domain.Entity.UpdatedAt = domain.Meta.UpdatedAt
			domain.Entity.c = c
			domains = append(domains, domain.Entity)
		}
		requestUrl = domainResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return domains, nil
}

func (c *Client) ListSharedDomains() ([]SharedDomain, error) {
	return c.ListSharedDomainsByQuery(nil)
}

func (c *Client) GetSharedDomainByGuid(guid string) (SharedDomain, error) {
       r := c.NewRequest("GET", "/v2/shared_domains/"+guid)
       resp, err := c.DoRequest(r)
       if err != nil {
               return SharedDomain{}, errors.Wrap(err, "Error requesting shared domain")
       }
       defer resp.Body.Close()
       retval, err := c.handleSharedDomainResp(resp)
       return *retval, err
}

func (c *Client) CreateSharedDomain(name string, internal bool, router_group_guid string) (*SharedDomain, error) {
	req := c.NewRequest("POST", "/v2/shared_domains")
	params := map[string]interface{}{
		"name":     name,
		"internal": internal,
	}

	if strings.TrimSpace(router_group_guid) != "" {
		params["router_group_guid"] = router_group_guid
	}

	req.obj = params

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.Wrapf(err, "Error creating shared domain %s, response code: %d", name, resp.StatusCode)
	}
	return c.handleSharedDomainResp(resp)
}

func (c *Client) DeleteSharedDomain(guid string, async bool) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/shared_domains/%s?async=%t", guid, async)))
	if err != nil {
		return err
	}
	if (async && (resp.StatusCode != http.StatusAccepted)) || (!async && (resp.StatusCode != http.StatusNoContent)) {
		return errors.Wrapf(err, "Error deleting organization %s, response code: %d", guid, resp.StatusCode)
	}
	return nil
}

func (c *Client) GetDomainByName(name string) (Domain, error) {
	q := url.Values{}
	q.Set("q", "name:"+name)
	domains, err := c.ListDomainsByQuery(q)
	if err != nil {
		return Domain{}, errors.Wrapf(err, "Error during domain lookup %s", name)
	}
	if len(domains) == 0 {
		return Domain{}, fmt.Errorf("Unable to find domain %s", name)
	}
	return domains[0], nil
}

func (c *Client) GetSharedDomainByName(name string) (SharedDomain, error) {
	q := url.Values{}
	q.Set("q", "name:"+name)
	domains, err := c.ListSharedDomainsByQuery(q)
	if err != nil {
		return SharedDomain{}, errors.Wrapf(err, "Error during shared domain lookup %s", name)
	}
	if len(domains) == 0 {
		return SharedDomain{}, fmt.Errorf("Unable to find shared domain %s", name)
	}
	return domains[0], nil
}

func (c *Client) CreateDomain(name, orgGuid string) (*Domain, error) {
	req := c.NewRequest("POST", "/v2/private_domains")
	req.obj = map[string]interface{}{
		"name": name,
		"owning_organization_guid": orgGuid,
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.Wrapf(err, "Error creating domain %s, response code: %d", name, resp.StatusCode)
	}
	return c.handleDomainResp(resp)
}

func (c *Client) DeleteDomain(guid string) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/private_domains/%s", guid)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleting domain %s, response code: %d", guid, resp.StatusCode)
	}
	return nil
}
func (c *Client) handleDomainResp(resp *http.Response) (*Domain, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var domainResource DomainResource
	err = json.Unmarshal(body, &domainResource)
	if err != nil {
		return nil, err
	}
	return c.mergeDomainResource(domainResource), nil
}

func (c *Client) handleSharedDomainResp(resp *http.Response) (*SharedDomain, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var domainResource SharedDomainResource
	err = json.Unmarshal(body, &domainResource)
	if err != nil {
		return nil, err
	}
	return c.mergeSharedDomainResource(domainResource), nil
}

func (c *Client) getDomainsResponse(requestUrl string) (DomainsResponse, error) {
	var domainResp DomainsResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return DomainsResponse{}, errors.Wrap(err, "Error requesting domains")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return DomainsResponse{}, errors.Wrap(err, "Error reading domains request")
	}
	err = json.Unmarshal(resBody, &domainResp)
	if err != nil {
		return DomainsResponse{}, errors.Wrap(err, "Error unmarshalling org")
	}
	return domainResp, nil
}

func (c *Client) mergeDomainResource(domainResource DomainResource) *Domain {
	domainResource.Entity.Guid = domainResource.Meta.Guid
	domainResource.Entity.c = c
	return &domainResource.Entity
}

func (c *Client) mergeSharedDomainResource(domainResource SharedDomainResource) *SharedDomain {
	domainResource.Entity.Guid = domainResource.Meta.Guid
	domainResource.Entity.c = c
	return &domainResource.Entity
}
