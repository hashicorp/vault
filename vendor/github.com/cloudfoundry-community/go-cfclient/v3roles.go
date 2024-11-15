package cfclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// V3Role implements role object. Roles control access to resources in organizations and spaces. Roles are assigned to users.
type V3Role struct {
	GUID          string                         `json:"guid,omitempty"`
	CreatedAt     string                         `json:"created_at,omitempty"`
	UpdatedAt     string                         `json:"updated_at,omitempty"`
	Type          string                         `json:"type,omitempty"`
	Relationships map[string]V3ToOneRelationship `json:"relationships,omitempty"`
	Links         map[string]Link                `json:"links,omitempty"`
}

type Included struct {
	Users         []V3User         `json:"users,omitempty"`
	Organizations []V3Organization `json:"organizations,omitempty"`
	Spaces        []V3Space        `json:"spaces,omitempty"`
}

type listV3RolesResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Resources  []V3Role   `json:"resources,omitempty"`
	Included   Included   `json:"included,omitempty"`
}

type createV3SpaceRoleRequest struct {
	RoleType      string                 `json:"type"`
	Relationships spaceUserRelationships `json:"relationships"`
}

type createV3OrganizationRoleRequest struct {
	RoleType      string               `json:"type"`
	Relationships orgUserRelationships `json:"relationships"`
}

type spaceUserRelationships struct {
	Space V3ToOneRelationship `json:"space"`
	User  V3ToOneRelationship `json:"user"`
}

type orgUserRelationships struct {
	Org  V3ToOneRelationship `json:"organization"`
	User V3ToOneRelationship `json:"user"`
}

func (c *Client) CreateV3SpaceRole(spaceGUID, userGUID, roleType string) (*V3Role, error) {
	spaceRel := V3ToOneRelationship{Data: V3Relationship{GUID: spaceGUID}}
	userRel := V3ToOneRelationship{Data: V3Relationship{GUID: userGUID}}
	req := c.NewRequest("POST", "/v3/roles")
	req.obj = createV3SpaceRoleRequest{
		RoleType:      roleType,
		Relationships: spaceUserRelationships{Space: spaceRel, User: userRel},
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating v3 role")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating v3 role, response code: %d", resp.StatusCode)
	}

	var role V3Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 role")
	}

	return &role, nil
}

func (c *Client) CreateV3OrganizationRole(orgGUID, userGUID, roleType string) (*V3Role, error) {
	orgRel := V3ToOneRelationship{Data: V3Relationship{GUID: orgGUID}}
	userRel := V3ToOneRelationship{Data: V3Relationship{GUID: userGUID}}
	req := c.NewRequest("POST", "/v3/roles")
	req.obj = createV3OrganizationRoleRequest{
		RoleType:      roleType,
		Relationships: orgUserRelationships{Org: orgRel, User: userRel},
	}
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating v3 role")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error creating v3 role, response code: %d", resp.StatusCode)
	}

	var role V3Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 role")
	}

	return &role, nil
}

// ListV3RolesByQuery retrieves roles based on query
func (c *Client) ListV3RolesByQuery(query url.Values) ([]V3Role, error) {
	var roles []V3Role
	requestURL, err := url.Parse("/v3/roles")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 space roles")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 space roles, response code: %d", resp.StatusCode)
		}

		var data listV3RolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 space roles")
		}

		roles = append(roles, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return roles, nil
}

func (c *Client) ListV3RoleUsersByQuery(query url.Values) ([]V3User, error) {
	var users []V3User
	requestURL, err := url.Parse("/v3/roles")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 roles")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 roles, response code: %d", resp.StatusCode)
		}

		var data listV3RolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 roles")
		}

		users = append(users, data.Included.Users...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return users, nil
}

func (c *Client) ListV3RoleAndUsersByQuery(query url.Values) ([]V3Role, []V3User, error) {
	var roles []V3Role
	var users []V3User
	requestURL, err := url.Parse("/v3/roles")
	if err != nil {
		return nil, nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Error requesting v3 roles")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, nil, fmt.Errorf("Error listing v3 roles, response code: %d", resp.StatusCode)
		}

		var data listV3RolesResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, nil, errors.Wrap(err, "Error parsing JSON from list v3 roles")
		}

		roles = append(roles, data.Resources...)
		users = append(users, data.Included.Users...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return roles, users, nil
}

// ListV3SpaceRolesByGUID retrieves roles based on query
func (c *Client) ListV3SpaceRolesByGUID(spaceGUID string) ([]V3Role, []V3User, error) {
	query := url.Values{}
	query["space_guids"] = []string{spaceGUID}
	query["include"] = []string{"user"}
	return c.ListV3RoleAndUsersByQuery(query)
}

// ListV3SpaceRolesByGUIDAndType retrieves roles based on query
func (c *Client) ListV3SpaceRolesByGUIDAndType(spaceGUID string, roleType string) ([]V3User, error) {
	query := url.Values{}
	query["space_guids"] = []string{spaceGUID}
	query["types"] = []string{roleType}
	query["include"] = []string{"user"}
	return c.ListV3RoleUsersByQuery(query)
}

// ListV3SpaceRolesByGUIDAndType retrieves roles based on query
func (c *Client) ListV3OrganizationRolesByGUIDAndType(orgGUID string, roleType string) ([]V3User, error) {
	query := url.Values{}
	query["organization_guids"] = []string{orgGUID}
	query["types"] = []string{roleType}
	query["include"] = []string{"user"}
	return c.ListV3RoleUsersByQuery(query)
}

// ListV3OrganizationRolesByGUID retrieves roles based on query
func (c *Client) ListV3OrganizationRolesByGUID(orgGUID string) ([]V3Role, []V3User, error) {
	query := url.Values{}
	query["organization_guids"] = []string{orgGUID}
	query["include"] = []string{"user"}
	return c.ListV3RoleAndUsersByQuery(query)
}

func (c *Client) DeleteV3Role(roleGUID string) error {
	req := c.NewRequest("DELETE", "/v3/roles/"+roleGUID)
	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting v3 role")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting v3 role with GUID [%s], response code: %d", roleGUID, resp.StatusCode)
	}

	return nil
}
