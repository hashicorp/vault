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

type UserRequest struct {
	Guid             string `json:"guid"`
	DefaultSpaceGuid string `json:"default_space_guid,omitempty"`
}

type Users []User

type User struct {
	Guid                  string `json:"guid"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
	Admin                 bool   `json:"admin"`
	Active                bool   `json:"active"`
	DefaultSpaceGUID      string `json:"default_space_guid"`
	Username              string `json:"username"`
	SpacesURL             string `json:"spaces_url"`
	OrgsURL               string `json:"organizations_url"`
	ManagedOrgsURL        string `json:"managed_organizations_url"`
	BillingManagedOrgsURL string `json:"billing_managed_organizations_url"`
	AuditedOrgsURL        string `json:"audited_organizations_url"`
	ManagedSpacesURL      string `json:"managed_spaces_url"`
	AuditedSpacesURL      string `json:"audited_spaces_url"`
	c                     *Client
}

type UserResource struct {
	Meta   Meta `json:"metadata"`
	Entity User `json:"entity"`
}

type UserResponse struct {
	Count     int            `json:"total_results"`
	Pages     int            `json:"total_pages"`
	NextUrl   string         `json:"next_url"`
	Resources []UserResource `json:"resources"`
}

// GetUserByGUID retrieves the user with the provided guid.
func (c *Client) GetUserByGUID(guid string) (User, error) {
	var userRes UserResource
	r := c.NewRequest("GET", "/v2/users/"+guid)
	resp, err := c.DoRequest(r)
	if err != nil {
		return User{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(body, &userRes)
	if err != nil {
		return User{}, err
	}
	return c.mergeUserResource(userRes), nil
}

func (c *Client) ListUsersByQuery(query url.Values) (Users, error) {
	var users []User
	requestUrl := "/v2/users?" + query.Encode()
	for {
		userResp, err := c.getUserResponse(requestUrl)
		if err != nil {
			return []User{}, err
		}
		for _, user := range userResp.Resources {
			user.Entity.Guid = user.Meta.Guid
			user.Entity.CreatedAt = user.Meta.CreatedAt
			user.Entity.UpdatedAt = user.Meta.UpdatedAt
			user.Entity.c = c
			users = append(users, user.Entity)
		}
		requestUrl = userResp.NextUrl
		if requestUrl == "" {
			break
		}
	}
	return users, nil
}

func (c *Client) ListUsers() (Users, error) {
	return c.ListUsersByQuery(nil)
}

func (c *Client) ListUserSpaces(userGuid string) ([]Space, error) {
	return c.fetchSpaces(fmt.Sprintf("/v2/users/%s/spaces", userGuid))
}

func (c *Client) ListUserAuditedSpaces(userGuid string) ([]Space, error) {
	return c.fetchSpaces(fmt.Sprintf("/v2/users/%s/audited_spaces", userGuid))
}

func (c *Client) ListUserManagedSpaces(userGuid string) ([]Space, error) {
	return c.fetchSpaces(fmt.Sprintf("/v2/users/%s/managed_spaces", userGuid))
}

func (c *Client) ListUserOrgs(userGuid string) ([]Org, error) {
	return c.fetchOrgs(fmt.Sprintf("/v2/users/%s/organizations", userGuid))
}

func (c *Client) ListUserManagedOrgs(userGuid string) ([]Org, error) {
	return c.fetchOrgs(fmt.Sprintf("/v2/users/%s/managed_organizations", userGuid))
}

func (c *Client) ListUserAuditedOrgs(userGuid string) ([]Org, error) {
	return c.fetchOrgs(fmt.Sprintf("/v2/users/%s/audited_organizations", userGuid))
}

func (c *Client) ListUserBillingManagedOrgs(userGuid string) ([]Org, error) {
	return c.fetchOrgs(fmt.Sprintf("/v2/users/%s/billing_managed_organizations", userGuid))
}

func (c *Client) CreateUser(req UserRequest) (User, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return User{}, err
	}
	r := c.NewRequestWithBody("POST", "/v2/users", buf)
	resp, err := c.DoRequest(r)
	if err != nil {
		return User{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return User{}, errors.Wrapf(err, "Error creating user, response code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return User{}, err
	}
	var userResource UserResource
	err = json.Unmarshal(body, &userResource)
	if err != nil {
		return User{}, err
	}
	user := userResource.Entity
	user.Guid = userResource.Meta.Guid
	user.c = c
	return user, nil
}

func (c *Client) DeleteUser(userGuid string) error {
	resp, err := c.DoRequest(c.NewRequest("DELETE", fmt.Sprintf("/v2/users/%s", userGuid)))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "Error deleting user %s, response code: %d", userGuid, resp.StatusCode)
	}
	return nil
}

func (u Users) GetUserByUsername(username string) User {
	for _, user := range u {
		if user.Username == username {
			return user
		}
	}
	return User{}
}

func (c *Client) getUserResponse(requestUrl string) (UserResponse, error) {
	var userResp UserResponse
	r := c.NewRequest("GET", requestUrl)
	resp, err := c.DoRequest(r)
	if err != nil {
		return UserResponse{}, errors.Wrap(err, "Error requesting users")
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return UserResponse{}, errors.Wrap(err, "Error reading user request")
	}
	err = json.Unmarshal(resBody, &userResp)
	if err != nil {
		return UserResponse{}, errors.Wrap(err, "Error unmarshalling user")
	}
	return userResp, nil
}

func (c *Client) mergeUserResource(u UserResource) User {
	u.Entity.Guid = u.Meta.Guid
	u.Entity.CreatedAt = u.Meta.CreatedAt
	u.Entity.UpdatedAt = u.Meta.UpdatedAt
	u.Entity.c = c
	return u.Entity
}
