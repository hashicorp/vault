package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const (
	teamsBasePath = "orgs/%s/teams"
)

// TeamsService is an interface for interfacing with the Teams
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/teams/
type TeamsService interface {
	List(context.Context, string, *ListOptions) ([]Team, *Response, error)
	Get(context.Context, string, string) (*Team, *Response, error)
	GetOneTeamByName(context.Context, string, string) (*Team, *Response, error)
	GetTeamUsersAssigned(context.Context, string, string) ([]AtlasUser, *Response, error)
	Create(context.Context, string, *Team) (*Team, *Response, error)
	Rename(context.Context, string, string, string) (*Team, *Response, error)
	UpdateTeamRoles(context.Context, string, string, *TeamUpdateRoles) ([]TeamRoles, *Response, error)
	AddUserToTeam(context.Context, string, string, string) ([]AtlasUser, *Response, error)
	RemoveUserToTeam(context.Context, string, string, string) (*Response, error)
	RemoveTeamFromOrganization(context.Context, string, string) (*Response, error)
	RemoveTeamFromProject(context.Context, string, string) (*Response, error)
}

//TeamsServiceOp handles communication with the Teams related methos of the
//MongoDB Atlas API
type TeamsServiceOp struct {
	client *Client
}

var _ TeamsService = &TeamsServiceOp{}

// Teams represents a array of project
type TeamsResponse struct {
	Links      []*Link `json:"links"`
	Results    []Team  `json:"results"`
	TotalCount int     `json:"totalCount"`
}

type Team struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name"`
	Usernames []string `json:"usernames,omitempty"`
}

//AtlasUserAssigned represents the user assigned to the project.
type AtlasUserAssigned struct {
	Links      []*Link     `json:"links"`
	Results    []AtlasUser `json:"results"`
	TotalCount int         `json:"totalCount"`
}

type TeamUpdateRoles struct {
	RoleNames []string `json:"roleNames"`
}

type TeamUpdateRolesResponse struct {
	Links      []*Link     `json:"links"`
	Results    []TeamRoles `json:"results"`
	TotalCount int         `json:"totalCount"`
}

type TeamRoles struct {
	Links     []*Link  `json:"links"`
	RoleNames []string `json:"roleNames"`
	TeamID    string   `json:"teamId"`
}

//GetAllTeams gets all teams.
//See more: https://docs.atlas.mongodb.com/reference/api/project-get-all/
func (s *TeamsServiceOp) List(ctx context.Context, orgID string, listOptions *ListOptions) ([]Team, *Response, error) {
	path := fmt.Sprintf(teamsBasePath, orgID)

	//Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(TeamsResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

//Get gets a single team in the organization by team ID.
//See more: https://docs.atlas.mongodb.com/reference/api/teams-get-one-by-id/
func (s *TeamsServiceOp) Get(ctx context.Context, orgID string, teamID string) (*Team, *Response, error) {
	if teamID == "" {
		return nil, nil, NewArgError("teamID", "must be set")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, teamID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Team)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//GetOneTeamByName gets a single project by its name.
//See more: https://docs.atlas.mongodb.com/reference/api/project-get-one-by-name/
func (s *TeamsServiceOp) GetOneTeamByName(ctx context.Context, orgID, teamName string) (*Team, *Response, error) {
	if teamName == "" {
		return nil, nil, NewArgError("teamName", "must be set")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/byName/%s", basePath, teamName)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Team)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//GetTeamUsersAssigned gets all the users assigned to a team.
//See more: https://docs.atlas.mongodb.com/reference/api/teams-get-all-users/
func (s *TeamsServiceOp) GetTeamUsersAssigned(ctx context.Context, orgID, teamID string) ([]AtlasUser, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s/users", basePath, teamID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AtlasUserAssigned)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

//Create creates a team.
//See more: https://docs.atlas.mongodb.com/reference/api/teams-create-one/
func (s *TeamsServiceOp) Create(ctx context.Context, orgID string, createRequest *Team) (*Team, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, fmt.Sprintf(teamsBasePath, orgID), createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Team)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//RenameTeam renames a team
//See more: https://docs.atlas.mongodb.com/reference/api/teams-rename-one/
func (s *TeamsServiceOp) Rename(ctx context.Context, orgID, teamID, teamName string) (*Team, *Response, error) {
	if teamName == "" {
		return nil, nil, NewArgError("teamName", "cannot be nil")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, teamID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, map[string]interface{}{
		"name": teamName,
	})
	if err != nil {
		return nil, nil, err
	}

	root := new(Team)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//UpdateTeamRoles Update the roles of a team in an Atlas project.
//See more: https://docs.atlas.mongodb.com/reference/api/teams-update-roles/
func (s *TeamsServiceOp) UpdateTeamRoles(ctx context.Context, orgID string, teamID string, updateTeamRolesRequest *TeamUpdateRoles) ([]TeamRoles, *Response, error) {
	if updateTeamRolesRequest == nil {
		return nil, nil, NewArgError("updateTeamRolesRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, teamID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateTeamRolesRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(TeamUpdateRolesResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

//AddUserToTeam adds a user from the organization associated with {ORG-ID} to the team with ID {TEAM-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/teams-add-user/
func (s *TeamsServiceOp) AddUserToTeam(ctx context.Context, orgID, teamID, userID string) ([]AtlasUser, *Response, error) {
	if userID == "" {
		return nil, nil, NewArgError("userID", "cannot be nil")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s/users", basePath, teamID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, map[string]interface{}{
		"id": userID,
	})

	if err != nil {
		return nil, nil, err
	}

	root := new(AtlasUserAssigned)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

//RemoveUserToTeam removes the specified user from the specified team.
//See more: https://docs.atlas.mongodb.com/reference/api/teams-remove-user/
func (s *TeamsServiceOp) RemoveUserToTeam(ctx context.Context, orgID, teamID, userID string) (*Response, error) {
	if userID == "" {
		return nil, NewArgError("userID", "cannot be nil")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s/users/%s", basePath, teamID, userID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

//RemoveTeamFromOrganization deletes the team with ID {TEAM-ID} from the organization specified to {ORG-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/teams-delete-one/
func (s *TeamsServiceOp) RemoveTeamFromOrganization(ctx context.Context, orgID, teamID string) (*Response, error) {
	if teamID == "" {
		return nil, NewArgError("teamID", "cannot be nil")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, teamID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

//RemoveTeamFromProject removes the specified team from the specified project.
//See more: https://docs.atlas.mongodb.com/reference/api/teams-remove-from-project/
func (s *TeamsServiceOp) RemoveTeamFromProject(ctx context.Context, groupID, teamID string) (*Response, error) {
	if teamID == "" {
		return nil, NewArgError("teamID", "cannot be nil")
	}

	path := fmt.Sprintf("groups/%s/teams/%s", groupID, teamID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
