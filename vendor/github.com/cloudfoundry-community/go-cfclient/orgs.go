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

type OrgResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	NextUrl   string        `json:"next_url"`
	Resources []OrgResource `json:"resources"`
}

type OrgResource struct {
	Meta   Meta `json:"metadata"`
	Entity Org  `json:"entity"`
}

type OrgUserResponse struct {
	Count     int            `json:"total_results"`
	Pages     int            `json:"total_pages"`
	NextURL   string         `json:"next_url"`
	Resources []UserResource `json:"resources"`
}

type Org struct {
	Guid                        string `json:"guid"`
	CreatedAt                   string `json:"created_at"`
	UpdatedAt                   string `json:"updated_at"`
	Name                        string `json:"name"`
	Status                      string `json:"status"`
	QuotaDefinitionGuid         string `json:"quota_definition_guid"`
	DefaultIsolationSegmentGuid string `json:"default_isolation_segment_guid"`
	c                           *Client
}

type OrgSummary struct {
	Guid   string             `json:"guid"`
	Name   string             `json:"name"`
	Status string             `json:"status"`
	Spaces []OrgSummarySpaces `json:"spaces"`
}

type OrgSummarySpaces struct {
	Guid         string `json:"guid"`
	Name         string `json:"name"`
	ServiceCount int    `json:"service_count"`
	AppCount     int    `json:"app_count"`
	MemDevTotal  int    `json:"mem_dev_total"`
	MemProdTotal int    `json:"mem_prod_total"`
}

type OrgRequest struct {
	Name                        string `json:"name"`
	Status                      string `json:"status,omitempty"`
	QuotaDefinitionGuid         string `json:"quota_definition_guid,omitempty"`
	DefaultIsolationSegmentGuid string `json:"default_isolation_segment_guid,omitempty"`
}

func (c *Client) ListOrgsByQuery(query url.Values) ([]Org, error) {
	var orgs []Org
	requestURL := "/v2/organizations?" + query.Encode()
	for {
		orgResp, err := c.getOrgResponse(requestURL)
		if err != nil {
			return []Org{}, err
		}
		for _, org := range orgResp.Resources {
			orgs = append(orgs, c.mergeOrgResource(org))
		}
		requestURL = orgResp.NextUrl
		if requestURL == "" {
			break
		}
	}
	return orgs, nil
}

func (c *Client) ListOrgs() ([]Org, error) {
	return c.ListOrgsByQuery(nil)
}

func (c *Client) GetOrgByName(name string) (Org, error) {
	var org Org
	q := url.Values{}
	q.Set("q", "name:"+name)
	orgs, err := c.ListOrgsByQuery(q)
	if err != nil {
		return org, err
	}
	if len(orgs) == 0 {
		return org, fmt.Errorf("Unable to find org %s", name)
	}
	return orgs[0], nil
}

func (c *Client) GetOrgByGuid(guid string) (Org, error) {
	var orgRes OrgResource
	r := c.NewRequest("GET", "/v2/organizations/"+guid)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return Org{}, err
	}
	err = json.Unmarshal(body, &orgRes)
	if err != nil {
		return Org{}, err
	}
	return c.mergeOrgResource(orgRes), nil
}

func (c *Client) OrgSpaces(guid string) ([]Space, error) {
	return c.fetchSpaces(fmt.Sprintf("/v2/organizations/%s/spaces", guid))
}

func (o *Org) Summary() (OrgSummary, error) {
	var orgSummary OrgSummary
	requestURL := fmt.Sprintf("/v2/organizations/%s/summary", o.Guid)
	r := o.c.NewRequest("GET", requestURL)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return OrgSummary{}, errors.Wrap(err, "Error requesting org summary")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return OrgSummary{}, errors.Wrap(err, "Error reading org summary body")
	}
	err = json.Unmarshal(resBody, &orgSummary)
	if err != nil {
		return OrgSummary{}, errors.Wrap(err, "Error unmarshalling org summary")
	}
	return orgSummary, nil
}

func (o *Org) Quota() (*OrgQuota, error) {
	var orgQuota *OrgQuota
	var orgQuotaResource OrgQuotasResource
	if o.QuotaDefinitionGuid == "" {
		return nil, nil
	}
	requestURL := fmt.Sprintf("/v2/quota_definitions/%s", o.QuotaDefinitionGuid)
	r := o.c.NewRequest("GET", requestURL)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return &OrgQuota{}, errors.Wrap(err, "Error requesting org quota")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return &OrgQuota{}, errors.Wrap(err, "Error reading org quota body")
	}
	err = json.Unmarshal(resBody, &orgQuotaResource)
	if err != nil {
		return &OrgQuota{}, errors.Wrap(err, "Error unmarshalling org quota")
	}
	orgQuota = &orgQuotaResource.Entity
	orgQuota.Guid = orgQuotaResource.Meta.Guid
	orgQuota.c = o.c
	return orgQuota, nil
}

func (c *Client) ListOrgUsersByQuery(orgGUID string, query url.Values) ([]User, error) {
	var users []User
	requestURL := fmt.Sprintf("/v2/organizations/%s/users?%s", orgGUID, query.Encode())
	for {
		omResp, err := c.getOrgUserResponse(requestURL)
		if err != nil {
			return []User{}, err
		}
		for _, u := range omResp.Resources {
			users = append(users, c.mergeUserResource(u))
		}
		requestURL = omResp.NextURL
		if requestURL == "" {
			break
		}
	}
	return users, nil
}

func (c *Client) ListOrgUsers(orgGUID string) ([]User, error) {
	return c.ListOrgUsersByQuery(orgGUID, nil)
}

func (c *Client) listOrgRolesByQuery(orgGUID, role string, query url.Values) ([]User, error) {
	var users []User
	requestURL := fmt.Sprintf("/v2/organizations/%s/%s?%s", orgGUID, role, query.Encode())
	for {
		omResp, err := c.getOrgUserResponse(requestURL)
		if err != nil {
			return []User{}, err
		}
		for _, u := range omResp.Resources {
			users = append(users, c.mergeUserResource(u))
		}
		requestURL = omResp.NextURL
		if requestURL == "" {
			break
		}
	}
	return users, nil
}

func (c *Client) ListOrgManagersByQuery(orgGUID string, query url.Values) ([]User, error) {
	return c.listOrgRolesByQuery(orgGUID, "managers", query)
}

func (c *Client) ListOrgManagers(orgGUID string) ([]User, error) {
	return c.ListOrgManagersByQuery(orgGUID, nil)
}

func (c *Client) ListOrgAuditorsByQuery(orgGUID string, query url.Values) ([]User, error) {
	return c.listOrgRolesByQuery(orgGUID, "auditors", query)
}

func (c *Client) ListOrgAuditors(orgGUID string) ([]User, error) {
	return c.ListOrgAuditorsByQuery(orgGUID, nil)
}

func (c *Client) ListOrgBillingManagersByQuery(orgGUID string, query url.Values) ([]User, error) {
	return c.listOrgRolesByQuery(orgGUID, "billing_managers", query)
}

func (c *Client) ListOrgBillingManagers(orgGUID string) ([]User, error) {
	return c.ListOrgBillingManagersByQuery(orgGUID, nil)
}

func (c *Client) AssociateOrgManager(orgGUID, userGUID string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateManager(userGUID)
}

func (c *Client) AssociateOrgManagerByUsername(orgGUID, name string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateManagerByUsername(name)
}

func (c *Client) AssociateOrgManagerByUsernameAndOrigin(orgGUID, name, origin string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateManagerByUsernameAndOrigin(name, origin)
}

func (c *Client) AssociateOrgUser(orgGUID, userGUID string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateUser(userGUID)
}

func (c *Client) AssociateOrgAuditor(orgGUID, userGUID string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateAuditor(userGUID)
}

func (c *Client) AssociateOrgUserByUsername(orgGUID, name string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateUserByUsername(name)
}

func (c *Client) AssociateOrgUserByUsernameAndOrigin(orgGUID, name, origin string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateUserByUsernameAndOrigin(name, origin)
}

func (c *Client) AssociateOrgAuditorByUsername(orgGUID, name string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateAuditorByUsername(name)
}

func (c *Client) AssociateOrgAuditorByUsernameAndOrigin(orgGUID, name, origin string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateAuditorByUsernameAndOrigin(name, origin)
}

func (c *Client) AssociateOrgBillingManager(orgGUID, userGUID string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateBillingManager(userGUID)
}

func (c *Client) AssociateOrgBillingManagerByUsername(orgGUID, name string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateBillingManagerByUsername(name)
}

func (c *Client) AssociateOrgBillingManagerByUsernameAndOrigin(orgGUID, name, origin string) (Org, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.AssociateBillingManagerByUsernameAndOrigin(name, origin)
}

func (c *Client) RemoveOrgManager(orgGUID, userGUID string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveManager(userGUID)
}

func (c *Client) RemoveOrgManagerByUsername(orgGUID, name string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveManagerByUsername(name)
}

func (c *Client) RemoveOrgManagerByUsernameAndOrigin(orgGUID, name, origin string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveManagerByUsernameAndOrigin(name, origin)
}

func (c *Client) RemoveOrgUser(orgGUID, userGUID string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveUser(userGUID)
}

func (c *Client) RemoveOrgAuditor(orgGUID, userGUID string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveAuditor(userGUID)
}

func (c *Client) RemoveOrgUserByUsername(orgGUID, name string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveUserByUsername(name)
}

func (c *Client) RemoveOrgUserByUsernameAndOrigin(orgGUID, name, origin string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveUserByUsernameAndOrigin(name, origin)
}

func (c *Client) RemoveOrgAuditorByUsername(orgGUID, name string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveAuditorByUsername(name)
}

func (c *Client) RemoveOrgAuditorByUsernameAndOrigin(orgGUID, name, origin string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveAuditorByUsernameAndOrigin(name, origin)
}

func (c *Client) RemoveOrgBillingManager(orgGUID, userGUID string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveBillingManager(userGUID)
}

func (c *Client) RemoveOrgBillingManagerByUsername(orgGUID, name string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveBillingManagerByUsername(name)
}

func (c *Client) RemoveOrgBillingManagerByUsernameAndOrigin(orgGUID, name, origin string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.RemoveBillingManagerByUsernameAndOrigin(name, origin)
}

func (c *Client) ListOrgSpaceQuotas(orgGUID string) ([]SpaceQuota, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.ListSpaceQuotas()
}

func (c *Client) ListOrgPrivateDomains(orgGUID string) ([]Domain, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.ListPrivateDomains()
}

func (c *Client) ShareOrgPrivateDomain(orgGUID, privateDomainGUID string) (*Domain, error) {
	org := Org{Guid: orgGUID, c: c}
	return org.SharePrivateDomain(privateDomainGUID)
}

func (c *Client) UnshareOrgPrivateDomain(orgGUID, privateDomainGUID string) error {
	org := Org{Guid: orgGUID, c: c}
	return org.UnsharePrivateDomain(privateDomainGUID)
}

func (o *Org) ListSpaceQuotas() ([]SpaceQuota, error) {
	var spaceQuotas []SpaceQuota
	requestURL := fmt.Sprintf("/v2/organizations/%s/space_quota_definitions", o.Guid)
	for {
		spaceQuotasResp, err := o.c.getSpaceQuotasResponse(requestURL)
		if err != nil {
			return []SpaceQuota{}, err
		}
		for _, resource := range spaceQuotasResp.Resources {
			spaceQuotas = append(spaceQuotas, *o.c.mergeSpaceQuotaResource(resource))
		}
		requestURL = spaceQuotasResp.NextUrl
		if requestURL == "" {
			break
		}
	}
	return spaceQuotas, nil
}

func (o *Org) ListPrivateDomains() ([]Domain, error) {
	var domains []Domain
	requestURL := fmt.Sprintf("/v2/organizations/%s/private_domains", o.Guid)
	for {
		domainsResp, err := o.c.getDomainsResponse(requestURL)
		if err != nil {
			return []Domain{}, err
		}
		for _, resource := range domainsResp.Resources {
			domains = append(domains, *o.c.mergeDomainResource(resource))
		}
		requestURL = domainsResp.NextUrl
		if requestURL == "" {
			break
		}
	}
	return domains, nil
}

func (o *Org) SharePrivateDomain(privateDomainGUID string) (*Domain, error) {
	requestURL := fmt.Sprintf("/v2/organizations/%s/private_domains/%s", o.Guid, privateDomainGUID)
	r := o.c.NewRequest("PUT", requestURL)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.Wrapf(err, "Error sharing domain %s for org %s, response code: %d", privateDomainGUID, o.Guid, resp.StatusCode)
	}
	return o.c.handleDomainResp(resp)
}

func (o *Org) UnsharePrivateDomain(privateDomainGUID string) error {
	requestURL := fmt.Sprintf("/v2/organizations/%s/private_domains/%s", o.Guid, privateDomainGUID)
	r := o.c.NewRequest("DELETE", requestURL)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error unsharing domain %s for org %s, response code: %d", privateDomainGUID, o.Guid, resp.StatusCode)
	}
	return nil
}

func (o *Org) associateRole(userGUID, role string) (Org, error) {
	requestURL := fmt.Sprintf("/v2/organizations/%s/%s/%s", o.Guid, role, userGUID)
	r := o.c.NewRequest("PUT", requestURL)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Org{}, errors.Wrapf(err, "Error associating %s %s, response code: %d", role, userGUID, resp.StatusCode)
	}
	return o.c.handleOrgResp(resp)
}

func (o *Org) associateRoleByUsernameAndOrigin(name, role, origin string) (Org, error) {
	requestURL := fmt.Sprintf("/v2/organizations/%s/%s", o.Guid, role)
	buf := bytes.NewBuffer(nil)
	payload := make(map[string]string)
	payload["username"] = name
	if origin != "" {
		payload["origin"] = origin
	}
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return Org{}, err
	}
	r := o.c.NewRequestWithBody("PUT", requestURL, buf)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Org{}, errors.Wrapf(err, "Error associating %s %s, response code: %d", role, name, resp.StatusCode)
	}
	return o.c.handleOrgResp(resp)
}

func (o *Org) AssociateManager(userGUID string) (Org, error) {
	return o.associateRole(userGUID, "managers")
}

func (o *Org) AssociateManagerByUsername(name string) (Org, error) {
	return o.associateRoleByUsernameAndOrigin(name, "managers", "")
}

func (o *Org) AssociateManagerByUsernameAndOrigin(name, origin string) (Org, error) {
	return o.associateRoleByUsernameAndOrigin(name, "managers", origin)
}

func (o *Org) AssociateUser(userGUID string) (Org, error) {
	requestURL := fmt.Sprintf("/v2/organizations/%s/users/%s", o.Guid, userGUID)
	r := o.c.NewRequest("PUT", requestURL)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Org{}, errors.Wrapf(err, "Error associating user %s, response code: %d", userGUID, resp.StatusCode)
	}
	return o.c.handleOrgResp(resp)
}

func (o *Org) AssociateAuditor(userGUID string) (Org, error) {
	return o.associateRole(userGUID, "auditors")
}

func (o *Org) AssociateAuditorByUsername(name string) (Org, error) {
	return o.associateRoleByUsernameAndOrigin(name, "auditors", "")
}

func (o *Org) AssociateAuditorByUsernameAndOrigin(name, origin string) (Org, error) {
	return o.associateRoleByUsernameAndOrigin(name, "auditors", origin)
}

func (o *Org) AssociateBillingManager(userGUID string) (Org, error) {
	return o.associateRole(userGUID, "billing_managers")
}

func (o *Org) AssociateBillingManagerByUsername(name string) (Org, error) {
	return o.associateRoleByUsernameAndOrigin(name, "billing_managers", "")
}
func (o *Org) AssociateBillingManagerByUsernameAndOrigin(name, origin string) (Org, error) {
	return o.associateRoleByUsernameAndOrigin(name, "billing_managers", origin)
}

func (o *Org) AssociateUserByUsername(name string) (Org, error) {
	return o.associateUserByUsernameAndOrigin(name, "")
}

func (o *Org) AssociateUserByUsernameAndOrigin(name, origin string) (Org, error) {
	return o.associateUserByUsernameAndOrigin(name, origin)
}

func (o *Org) associateUserByUsernameAndOrigin(name, origin string) (Org, error) {
	requestURL := fmt.Sprintf("/v2/organizations/%s/users", o.Guid)
	buf := bytes.NewBuffer(nil)
	payload := make(map[string]string)
	payload["username"] = name
	if origin != "" {
		payload["origin"] = origin
	}
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return Org{}, err
	}
	r := o.c.NewRequestWithBody("PUT", requestURL, buf)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Org{}, errors.Wrapf(err, "Error associating user %s, response code: %d", name, resp.StatusCode)
	}
	return o.c.handleOrgResp(resp)
}

func (o *Org) removeRole(userGUID, role string) error {
	requestURL := fmt.Sprintf("/v2/organizations/%s/%s/%s", o.Guid, role, userGUID)
	r := o.c.NewRequest("DELETE", requestURL)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error removing %s %s, response code: %d", role, userGUID, resp.StatusCode)
	}
	return nil
}

func (o *Org) removeRoleByUsernameAndOrigin(name, role, origin string) error {
	var requestURL string
	var method string
	buf := bytes.NewBuffer(nil)
	payload := make(map[string]string)
	payload["username"] = name
	if origin != "" {
		requestURL = fmt.Sprintf("/v2/organizations/%s/%s/remove", o.Guid, role)
		method = "POST"
		payload["origin"] = origin
	} else {
		requestURL = fmt.Sprintf("/v2/organizations/%s/%s", o.Guid, role)
		method = "DELETE"
	}
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return err
	}

	r := o.c.NewRequestWithBody(method, requestURL, buf)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error removing manager %s, response code: %d", name, resp.StatusCode)
	}
	return nil
}

func (o *Org) RemoveManager(userGUID string) error {
	return o.removeRole(userGUID, "managers")
}

func (o *Org) RemoveManagerByUsername(name string) error {
	return o.removeRoleByUsernameAndOrigin(name, "managers", "")
}
func (o *Org) RemoveManagerByUsernameAndOrigin(name, origin string) error {
	return o.removeRoleByUsernameAndOrigin(name, "managers", origin)
}

func (o *Org) RemoveAuditor(userGUID string) error {
	return o.removeRole(userGUID, "auditors")
}

func (o *Org) RemoveAuditorByUsername(name string) error {
	return o.removeRoleByUsernameAndOrigin(name, "auditors", "")
}
func (o *Org) RemoveAuditorByUsernameAndOrigin(name, origin string) error {
	return o.removeRoleByUsernameAndOrigin(name, "auditors", origin)
}

func (o *Org) RemoveBillingManager(userGUID string) error {
	return o.removeRole(userGUID, "billing_managers")
}

func (o *Org) RemoveBillingManagerByUsername(name string) error {
	return o.removeRoleByUsernameAndOrigin(name, "billing_managers", "")
}

func (o *Org) RemoveBillingManagerByUsernameAndOrigin(name, origin string) error {
	return o.removeRoleByUsernameAndOrigin(name, "billing_managers", origin)
}

func (o *Org) RemoveUser(userGUID string) error {
	requestURL := fmt.Sprintf("/v2/organizations/%s/users/%s", o.Guid, userGUID)
	r := o.c.NewRequest("DELETE", requestURL)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error removing user %s, response code: %d", userGUID, resp.StatusCode)
	}
	return nil
}

func (o *Org) RemoveUserByUsername(name string) error {
	return o.removeUserByUsernameAndOrigin(name, "")
}

func (o *Org) RemoveUserByUsernameAndOrigin(name, origin string) error {
	return o.removeUserByUsernameAndOrigin(name, origin)
}

func (o *Org) removeUserByUsernameAndOrigin(name, origin string) error {
	var requestURL string
	var method string
	buf := bytes.NewBuffer(nil)
	payload := make(map[string]string)
	payload["username"] = name
	if origin != "" {
		payload["origin"] = origin
		requestURL = fmt.Sprintf("/v2/organizations/%s/users/remove", o.Guid)
		method = "POST"
	} else {
		requestURL = fmt.Sprintf("/v2/organizations/%s/users", o.Guid)
		method = "DELETE"
	}
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return err
	}
	r := o.c.NewRequestWithBody(method, requestURL, buf)
	resp, err := o.c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error removing user %s, response code: %d", name, resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateOrg(req OrgRequest) (Org, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return Org{}, err
	}
	r := c.NewRequestWithBody("POST", "/v2/organizations", buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Org{}, errors.Wrapf(err, "Error creating organization, response code: %d", resp.StatusCode)
	}
	return c.handleOrgResp(resp)
}

func (c *Client) UpdateOrg(orgGUID string, orgRequest OrgRequest) (Org, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(orgRequest)
	if err != nil {
		return Org{}, err
	}
	r := c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/organizations/%s", orgGUID), buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return Org{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Org{}, errors.Wrapf(err, "Error updating organization, response code: %d", resp.StatusCode)
	}
	return c.handleOrgResp(resp)
}

func (c *Client) DeleteOrg(guid string, recursive, async bool) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/organizations/%s?recursive=%t&async=%t", guid, recursive, async)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleting organization %s, response code: %d", guid, resp.StatusCode)
	}
	return nil
}

func (c *Client) getOrgResponse(requestURL string) (OrgResponse, error) {
	var orgResp OrgResponse
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return OrgResponse{}, errors.Wrap(err, "Error requesting orgs")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return OrgResponse{}, errors.Wrap(err, "Error reading org request")
	}
	err = json.Unmarshal(resBody, &orgResp)
	if err != nil {
		return OrgResponse{}, errors.Wrap(err, "Error unmarshalling org")
	}
	return orgResp, nil
}

func (c *Client) fetchOrgs(requestURL string) ([]Org, error) {
	var orgs []Org
	for {
		orgResp, err := c.getOrgResponse(requestURL)
		if err != nil {
			return []Org{}, err
		}
		for _, org := range orgResp.Resources {
			orgs = append(orgs, c.mergeOrgResource(org))
		}
		requestURL = orgResp.NextUrl
		if requestURL == "" {
			break
		}
	}
	return orgs, nil
}

func (c *Client) handleOrgResp(resp *http.Response) (Org, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return Org{}, err
	}
	var orgResource OrgResource
	err = json.Unmarshal(body, &orgResource)
	if err != nil {
		return Org{}, err
	}
	return c.mergeOrgResource(orgResource), nil
}

func (c *Client) getOrgUserResponse(requestURL string) (OrgUserResponse, error) {
	var omResp OrgUserResponse
	r := c.NewRequest("GET", requestURL)
	resp, err := c.DoRequest(r)
	if err != nil {
		return OrgUserResponse{}, errors.Wrap(err, "error requesting org managers")
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return OrgUserResponse{}, errors.Wrap(err, "error reading org managers response body")
	}
	if err := json.Unmarshal(resBody, &omResp); err != nil {
		return OrgUserResponse{}, errors.Wrap(err, "error unmarshaling org managers")
	}
	return omResp, nil
}

func (c *Client) mergeOrgResource(org OrgResource) Org {
	org.Entity.Guid = org.Meta.Guid
	org.Entity.CreatedAt = org.Meta.CreatedAt
	org.Entity.UpdatedAt = org.Meta.UpdatedAt
	org.Entity.c = c
	return org.Entity
}

func (c *Client) DefaultIsolationSegmentForOrg(orgGUID, isolationSegmentGUID string) error {
	return c.updateOrgDefaultIsolationSegment(orgGUID, map[string]interface{}{"guid": isolationSegmentGUID})
}

func (c *Client) ResetDefaultIsolationSegmentForOrg(orgGUID string) error {
	return c.updateOrgDefaultIsolationSegment(orgGUID, nil)
}

func (c *Client) updateOrgDefaultIsolationSegment(orgGUID string, data interface{}) error {
	requestURL := fmt.Sprintf("/v3/organizations/%s/relationships/default_isolation_segment", orgGUID)
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(map[string]interface{}{"data": data})
	if err != nil {
		return err
	}
	r := c.NewRequestWithBody("PATCH", requestURL, buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(err, "Error setting default isolation segment for org %s, response code: %d", orgGUID, resp.StatusCode)
	}
	return nil
}
