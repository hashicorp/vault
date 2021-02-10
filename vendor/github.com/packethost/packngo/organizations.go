package packngo

import (
	"path"
)

// API documentation https://metal.equinix.com/developers/api/organizations/
const organizationBasePath = "/organizations"

// OrganizationService interface defines available organization methods
type OrganizationService interface {
	List(*ListOptions) ([]Organization, *Response, error)
	Get(string, *GetOptions) (*Organization, *Response, error)
	Create(*OrganizationCreateRequest) (*Organization, *Response, error)
	Update(string, *OrganizationUpdateRequest) (*Organization, *Response, error)
	Delete(string) (*Response, error)
	ListPaymentMethods(string) ([]PaymentMethod, *Response, error)
	ListEvents(string, *ListOptions) ([]Event, *Response, error)
}

type organizationsRoot struct {
	Organizations []Organization `json:"organizations"`
	Meta          meta           `json:"meta"`
}

// Organization represents an Equinix Metal organization
type Organization struct {
	ID           string    `json:"id"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	Website      string    `json:"website,omitempty"`
	Twitter      string    `json:"twitter,omitempty"`
	Created      string    `json:"created_at,omitempty"`
	Updated      string    `json:"updated_at,omitempty"`
	Address      Address   `json:"address,omitempty"`
	TaxID        string    `json:"tax_id,omitempty"`
	MainPhone    string    `json:"main_phone,omitempty"`
	BillingPhone string    `json:"billing_phone,omitempty"`
	CreditAmount float64   `json:"credit_amount,omitempty"`
	Logo         string    `json:"logo,omitempty"`
	LogoThumb    string    `json:"logo_thumb,omitempty"`
	Projects     []Project `json:"projects,omitempty"`
	URL          string    `json:"href,omitempty"`
	Users        []User    `json:"members,omitempty"`
	Owners       []User    `json:"owners,omitempty"`
}

func (o Organization) String() string {
	return Stringify(o)
}

// OrganizationCreateRequest type used to create an Equinix Metal organization
type OrganizationCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Twitter     string `json:"twitter"`
	Logo        string `json:"logo"`
}

func (o OrganizationCreateRequest) String() string {
	return Stringify(o)
}

// OrganizationUpdateRequest type used to update an Equinix Metal organization
type OrganizationUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Website     *string `json:"website,omitempty"`
	Twitter     *string `json:"twitter,omitempty"`
	Logo        *string `json:"logo,omitempty"`
}

func (o OrganizationUpdateRequest) String() string {
	return Stringify(o)
}

// OrganizationServiceOp implements OrganizationService
type OrganizationServiceOp struct {
	client *Client
}

// List returns the user's organizations
func (s *OrganizationServiceOp) List(opts *ListOptions) (orgs []Organization, resp *Response, err error) {
	subset := new(organizationsRoot)

	apiPathQuery := opts.WithQuery(organizationBasePath)

	for {
		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		orgs = append(orgs, subset.Organizations...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// Get returns a organization by id
func (s *OrganizationServiceOp) Get(organizationID string, opts *GetOptions) (*Organization, *Response, error) {
	endpointPath := path.Join(organizationBasePath, organizationID)
	apiPathQuery := opts.WithQuery(endpointPath)
	organization := new(Organization)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, err
}

// Create creates a new organization
func (s *OrganizationServiceOp) Create(createRequest *OrganizationCreateRequest) (*Organization, *Response, error) {
	organization := new(Organization)

	resp, err := s.client.DoRequest("POST", organizationBasePath, createRequest, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, err
}

// Update updates an organization
func (s *OrganizationServiceOp) Update(id string, updateRequest *OrganizationUpdateRequest) (*Organization, *Response, error) {
	apiPath := path.Join(organizationBasePath, id)
	organization := new(Organization)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, err
}

// Delete deletes an organizationID
func (s *OrganizationServiceOp) Delete(organizationID string) (*Response, error) {
	apiPath := path.Join(organizationBasePath, organizationID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}

// ListPaymentMethods returns PaymentMethods for an organization
func (s *OrganizationServiceOp) ListPaymentMethods(organizationID string) ([]PaymentMethod, *Response, error) {
	apiPath := path.Join(organizationBasePath, organizationID, paymentMethodBasePath)
	root := new(paymentMethodsRoot)

	resp, err := s.client.DoRequest("GET", apiPath, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.PaymentMethods, resp, err
}

// ListEvents returns list of organization events
func (s *OrganizationServiceOp) ListEvents(organizationID string, listOpt *ListOptions) ([]Event, *Response, error) {
	apiPath := path.Join(organizationBasePath, organizationID, eventBasePath)

	return listEvents(s.client, apiPath, listOpt)
}
