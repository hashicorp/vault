package packngo

import (
	"path"
)

const projectBasePath = "/projects"

// ProjectService interface defines available project methods
type ProjectService interface {
	List(listOpt *ListOptions) ([]Project, *Response, error)
	Get(string, *GetOptions) (*Project, *Response, error)
	Create(*ProjectCreateRequest) (*Project, *Response, error)
	Update(string, *ProjectUpdateRequest) (*Project, *Response, error)
	Delete(string) (*Response, error)
	ListBGPSessions(projectID string, listOpt *ListOptions) ([]BGPSession, *Response, error)
	ListEvents(string, *ListOptions) ([]Event, *Response, error)
	ListSSHKeys(projectID string, searchOpt *SearchOptions) ([]SSHKey, *Response, error)
}

type projectsRoot struct {
	Projects []Project `json:"projects"`
	Meta     meta      `json:"meta"`
}

// Project represents an Equinix Metal project
type Project struct {
	ID              string        `json:"id"`
	Name            string        `json:"name,omitempty"`
	Organization    Organization  `json:"organization,omitempty"`
	Created         string        `json:"created_at,omitempty"`
	Updated         string        `json:"updated_at,omitempty"`
	Users           []User        `json:"members,omitempty"`
	Devices         []Device      `json:"devices,omitempty"`
	SSHKeys         []SSHKey      `json:"ssh_keys,omitempty"`
	URL             string        `json:"href,omitempty"`
	PaymentMethod   PaymentMethod `json:"payment_method,omitempty"`
	BackendTransfer bool          `json:"backend_transfer_enabled"`
}

func (p Project) String() string {
	return Stringify(p)
}

// ProjectCreateRequest type used to create an Equinix Metal project
type ProjectCreateRequest struct {
	Name            string `json:"name"`
	PaymentMethodID string `json:"payment_method_id,omitempty"`
	OrganizationID  string `json:"organization_id,omitempty"`
}

func (p ProjectCreateRequest) String() string {
	return Stringify(p)
}

// ProjectUpdateRequest type used to update an Equinix Metal project
type ProjectUpdateRequest struct {
	Name            *string `json:"name,omitempty"`
	PaymentMethodID *string `json:"payment_method_id,omitempty"`
	BackendTransfer *bool   `json:"backend_transfer_enabled,omitempty"`
}

func (p ProjectUpdateRequest) String() string {
	return Stringify(p)
}

// ProjectServiceOp implements ProjectService
type ProjectServiceOp struct {
	client requestDoer
}

// List returns the user's projects
func (s *ProjectServiceOp) List(opts *ListOptions) (projects []Project, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(projectBasePath)

	for {
		subset := new(projectsRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		projects = append(projects, subset.Projects...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns a project by id
func (s *ProjectServiceOp) Get(projectID string, opts *GetOptions) (*Project, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID)
	apiPathQuery := opts.WithQuery(endpointPath)
	project := new(Project)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, err
}

// Create creates a new project
func (s *ProjectServiceOp) Create(createRequest *ProjectCreateRequest) (*Project, *Response, error) {
	project := new(Project)

	resp, err := s.client.DoRequest("POST", projectBasePath, createRequest, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

// Update updates a project
func (s *ProjectServiceOp) Update(id string, updateRequest *ProjectUpdateRequest) (*Project, *Response, error) {
	apiPath := path.Join(projectBasePath, id)
	project := new(Project)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

// Delete deletes a project
func (s *ProjectServiceOp) Delete(projectID string) (*Response, error) {
	apiPath := path.Join(projectBasePath, projectID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}

// ListBGPSessions returns all BGP Sessions associated with the project
func (s *ProjectServiceOp) ListBGPSessions(projectID string, opts *ListOptions) (bgpSessions []BGPSession, resp *Response, err error) {
	endpointPath := path.Join(projectBasePath, projectID, bgpSessionBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(bgpSessionsRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		bgpSessions = append(bgpSessions, subset.Sessions...)
		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// ListSSHKeys returns all SSH Keys associated with the project
func (s *ProjectServiceOp) ListSSHKeys(projectID string, opts *SearchOptions) (sshKeys []SSHKey, resp *Response, err error) {

	endpointPath := path.Join(projectBasePath, projectID, sshKeyBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	subset := new(sshKeyRoot)

	resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
	if err != nil {
		return nil, resp, err
	}

	sshKeys = append(sshKeys, subset.SSHKeys...)

	return
}

// ListEvents returns list of project events
func (s *ProjectServiceOp) ListEvents(projectID string, listOpt *ListOptions) ([]Event, *Response, error) {
	apiPath := path.Join(projectBasePath, projectID, eventBasePath)

	return listEvents(s.client, apiPath, listOpt)
}
