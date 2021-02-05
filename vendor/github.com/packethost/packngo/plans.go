package packngo

import "path"

const planBasePath = "/plans"

// PlanService interface defines available plan methods
type PlanService interface {
	List(*ListOptions) ([]Plan, *Response, error)
	ProjectList(string, *ListOptions) ([]Plan, *Response, error)
	OrganizationList(string, *ListOptions) ([]Plan, *Response, error)
}

type planRoot struct {
	Plans []Plan `json:"plans"`
}

// Plan represents an Equinix Metal service plan
type Plan struct {
	ID              string     `json:"id"`
	Slug            string     `json:"slug,omitempty"`
	Name            string     `json:"name,omitempty"`
	Description     string     `json:"description,omitempty"`
	Line            string     `json:"line,omitempty"`
	Legacy          bool       `json:"legacy,omitempty"`
	Specs           *Specs     `json:"specs,omitempty"`
	Pricing         *Pricing   `json:"pricing,omitempty"`
	DeploymentTypes []string   `json:"deployment_types"`
	Class           string     `json:"class"`
	AvailableIn     []Facility `json:"available_in"`
}

func (p Plan) String() string {
	return Stringify(p)
}

// Specs - the server specs for a plan
type Specs struct {
	Cpus     []*Cpus   `json:"cpus,omitempty"`
	Memory   *Memory   `json:"memory,omitempty"`
	Drives   []*Drives `json:"drives,omitempty"`
	Nics     []*Nics   `json:"nics,omitempty"`
	Features *Features `json:"features,omitempty"`
}

func (s Specs) String() string {
	return Stringify(s)
}

// Cpus - the CPU config details for specs on a plan
type Cpus struct {
	Count int    `json:"count,omitempty"`
	Type  string `json:"type,omitempty"`
}

func (c Cpus) String() string {
	return Stringify(c)
}

// Memory - the RAM config details for specs on a plan
type Memory struct {
	Total string `json:"total,omitempty"`
}

func (m Memory) String() string {
	return Stringify(m)
}

// Drives - the storage config details for specs on a plan
type Drives struct {
	Count int    `json:"count,omitempty"`
	Size  string `json:"size,omitempty"`
	Type  string `json:"type,omitempty"`
}

func (d Drives) String() string {
	return Stringify(d)
}

// Nics - the network hardware details for specs on a plan
type Nics struct {
	Count int    `json:"count,omitempty"`
	Type  string `json:"type,omitempty"`
}

func (n Nics) String() string {
	return Stringify(n)
}

// Features - other features in the specs for a plan
type Features struct {
	Raid bool `json:"raid,omitempty"`
	Txt  bool `json:"txt,omitempty"`
}

func (f Features) String() string {
	return Stringify(f)
}

// Pricing - the pricing options on a plan
type Pricing struct {
	Hour  float32 `json:"hour,omitempty"`
	Month float32 `json:"month,omitempty"`
}

func (p Pricing) String() string {
	return Stringify(p)
}

// PlanServiceOp implements PlanService
type PlanServiceOp struct {
	client *Client
}

func planList(c *Client, apiPath string, opts *ListOptions) ([]Plan, *Response, error) {
	root := new(planRoot)
	apiPathQuery := opts.WithQuery(apiPath)

	resp, err := c.DoRequest("GET", apiPathQuery, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Plans, resp, err

}

// List method returns all available plans
func (s *PlanServiceOp) List(opts *ListOptions) ([]Plan, *Response, error) {
	return planList(s.client, planBasePath, opts)

}

// ProjectList method returns plans available in a project
func (s *PlanServiceOp) ProjectList(projectID string, opts *ListOptions) ([]Plan, *Response, error) {
	return planList(s.client, path.Join(projectBasePath, projectID, planBasePath), opts)
}

// OrganizationList method returns plans available in an organization
func (s *PlanServiceOp) OrganizationList(organizationID string, opts *ListOptions) ([]Plan, *Response, error) {
	return planList(s.client, path.Join(organizationBasePath, organizationID, planBasePath), opts)
}
