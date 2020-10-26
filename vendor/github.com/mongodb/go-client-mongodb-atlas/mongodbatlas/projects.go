package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const (
	GROUP_OWNER                  = "GROUP_OWNER"                  //GROUP_OWNER - Project Owner
	GROUP_READ_ONLY              = "GROUP_READ_ONLY"              //GROUP_READ_ONLY - Project Read Only
	GROUP_DATA_ACCESS_ADMIN      = "GROUP_DATA_ACCESS_ADMIN"      //GROUP_DATA_ACCESS_ADMIN - Project Data Access Admin
	GROUP_DATA_ACCESS_READ_WRITE = "GROUP_DATA_ACCESS_READ_WRITE" //GROUP_DATA_ACCESS_READ_WRITE - Project Data Access Read/Write
	GROUP_DATA_ACCESS_READ_ONLY  = "GROUP_DATA_ACCESS_READ_ONLY"  //GROUP_DATA_ACCESS_READ_ONLY - Project Data Access Read Only
	projectBasePath              = "groups"
)

// ProjectsService is an interface for interfacing with the Projects
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/projects/
type ProjectsService interface {
	GetAllProjects(context.Context) (*Projects, *Response, error)
	GetOneProject(context.Context, string) (*Project, *Response, error)
	GetOneProjectByName(context.Context, string) (*Project, *Response, error)
	Create(context.Context, *Project) (*Project, *Response, error)
	Delete(context.Context, string) (*Response, error)
	GetProjectTeamsAssigned(context.Context, string) (*TeamsAssigned, *Response, error)
	AddTeamsToProject(context.Context, string, *ProjectTeam) (*TeamsAssigned, *Response, error)
}

//ProjectsServiceOp handles communication with the Projects related methos of the
//MongoDB Atlas API
type ProjectsServiceOp struct {
	client *Client
}

var _ ProjectsService = &ProjectsServiceOp{}

// Project represents the structure of a project.
type Project struct {
	ID           string  `json:"id,omitempty"`
	OrgID        string  `json:"orgId,omitempty"`
	Name         string  `json:"name,omitempty"`
	ClusterCount int     `json:"clusterCount,omitempty"`
	Created      string  `json:"created,omitempty"`
	Links        []*Link `json:"links,omitempty"`
}

// Projects represents a array of project
type Projects struct {
	Links      []*Link    `json:"links"`
	Results    []*Project `json:"results"`
	TotalCount int        `json:"totalCount"`
}

// Result is part og TeamsAssigned structure
type Result struct {
	Links     []*Link  `json:"links"`
	RoleNames []string `json:"roleNames"`
	TeamID    string   `json:"teamId"`
}

// RoleName represents the kind of user role in your project
type RoleName struct {
	RoleName string `json:"rolesNames"`
}

// ProjectTeam reperesents the kind of role that has the team
type ProjectTeam struct {
	TeamID string      `json:"teamId"`
	Roles  []*RoleName `json:"roles"`
}

// TeamsAssigned represents the one team assigned to the project.
type TeamsAssigned struct {
	Links      []*Link   `json:"links"`
	Results    []*Result `json:"results"`
	TotalCount int       `json:"totalCount"`
}

//GetAllProjects gets all project.
//See more: https://docs.atlas.mongodb.com/reference/api/project-get-all/
func (s *ProjectsServiceOp) GetAllProjects(ctx context.Context) (*Projects, *Response, error) {

	req, err := s.client.NewRequest(ctx, http.MethodGet, projectBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Projects)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

//GetOneProject gets a single project.
//See more: https://docs.atlas.mongodb.com/reference/api/project-get-one/
func (s *ProjectsServiceOp) GetOneProject(ctx context.Context, projectID string) (*Project, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", projectBasePath, projectID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Project)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//GetOneProjectByName gets a single project by its name.
//See more: https://docs.atlas.mongodb.com/reference/api/project-get-one-by-name/
func (s *ProjectsServiceOp) GetOneProjectByName(ctx context.Context, projectName string) (*Project, *Response, error) {
	if projectName == "" {
		return nil, nil, NewArgError("projectName", "must be set")
	}

	path := fmt.Sprintf("%s/byName/%s", projectBasePath, projectName)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Project)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Create creates a project.
//See more: https://docs.atlas.mongodb.com/reference/api/project-create-one/
func (s *ProjectsServiceOp) Create(ctx context.Context, createRequest *Project) (*Project, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, projectBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Project)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Delete deletes a project.
// See more: https://docs.atlas.mongodb.com/reference/api/project-delete-one/
func (s *ProjectsServiceOp) Delete(ctx context.Context, projectID string) (*Response, error) {
	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}

	basePath := fmt.Sprintf("%s/%s", projectBasePath, projectID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, basePath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

//GetProjectTeamsAssigned gets all the teams assigned to a project.
//See more: https://docs.atlas.mongodb.com/reference/api/project-get-teams/
func (s *ProjectsServiceOp) GetProjectTeamsAssigned(ctx context.Context, projectID string) (*TeamsAssigned, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}

	path := fmt.Sprintf("%s/%s/teams", projectBasePath, projectID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(TeamsAssigned)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//AddTeamsToProject adds teams to a project
//See more: https://docs.atlas.mongodb.com/reference/api/project-add-team/
func (s *ProjectsServiceOp) AddTeamsToProject(ctx context.Context, projectID string, createRequest *ProjectTeam) (*TeamsAssigned, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/teams", projectBasePath, projectID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(TeamsAssigned)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
