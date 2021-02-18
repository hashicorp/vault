package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

//
// GET /api/permissions
//

// PermissionInfo represents a user's permission in a virtual host.
type PermissionInfo struct {
	User  string `json:"user"`
	Vhost string `json:"vhost"`

	// Configuration permissions
	Configure string `json:"configure"`
	// Write permissions
	Write string `json:"write"`
	// Read permissions
	Read string `json:"read"`
}

// ListPermissions returns permissions for all users and virtual hosts.
func (c *Client) ListPermissions() (rec []PermissionInfo, err error) {
	req, err := newGETRequest(c, "permissions/")
	if err != nil {
		return []PermissionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []PermissionInfo{}, err
	}

	return rec, nil
}

//
// GET /api/users/{user}/permissions
//

// ListPermissionsOf returns permissions of a specific user.
func (c *Client) ListPermissionsOf(username string) (rec []PermissionInfo, err error) {
	req, err := newGETRequest(c, "users/"+url.PathEscape(username)+"/permissions")
	if err != nil {
		return []PermissionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []PermissionInfo{}, err
	}

	return rec, nil
}

//
// GET /api/permissions/{vhost}/{user}
//

// GetPermissionsIn returns permissions of user in virtual host.
func (c *Client) GetPermissionsIn(vhost, username string) (rec PermissionInfo, err error) {
	req, err := newGETRequest(c, "permissions/"+url.PathEscape(vhost)+"/"+url.PathEscape(username))
	if err != nil {
		return PermissionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return PermissionInfo{}, err
	}

	return rec, nil
}

//
// PUT /api/permissions/{vhost}/{user}
//

// Permissions represents permissions of a user in a virtual host. Use this type to set
// permissions of the user.
type Permissions struct {
	Configure string `json:"configure"`
	Write     string `json:"write"`
	Read      string `json:"read"`
}

// UpdatePermissionsIn sets permissions of a user in a virtual host.
func (c *Client) UpdatePermissionsIn(vhost, username string, permissions Permissions) (res *http.Response, err error) {
	body, err := json.Marshal(permissions)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "permissions/"+url.PathEscape(vhost)+"/"+url.PathEscape(username), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/permissions/{vhost}/{user}
//

// ClearPermissionsIn clears (deletes) permissions of user in virtual host.
func (c *Client) ClearPermissionsIn(vhost, username string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "permissions/"+url.PathEscape(vhost)+"/"+url.PathEscape(username), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
