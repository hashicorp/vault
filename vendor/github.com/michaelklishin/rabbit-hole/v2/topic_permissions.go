package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

//
// GET /api/topic-permissions
//

// TopicPermissionInfo represents a user's permissions on a topic.
type TopicPermissionInfo struct {
	User  string `json:"user"`
	Vhost string `json:"vhost"`

	// Configuration topic-permissions
	Exchange string `json:"exchange"`
	// Write topic-permissions
	Write string `json:"write"`
	// Read topic-permissions
	Read string `json:"read"`
}

// ListTopicPermissions returns topic-permissions for all users and virtual hosts.
func (c *Client) ListTopicPermissions() (rec []TopicPermissionInfo, err error) {
	req, err := newGETRequest(c, "topic-permissions/")
	if err != nil {
		return []TopicPermissionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []TopicPermissionInfo{}, err
	}

	return rec, nil
}

//
// GET /api/users/{user}/topic-permissions
//

// ListTopicPermissionsOf returns topic-permissions of a specific user.
func (c *Client) ListTopicPermissionsOf(username string) (rec []TopicPermissionInfo, err error) {
	req, err := newGETRequest(c, "users/"+url.PathEscape(username)+"/topic-permissions")
	if err != nil {
		return []TopicPermissionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []TopicPermissionInfo{}, err
	}

	return rec, nil
}

//
// GET /api/topic-permissions/{vhost}/{user}
//

// GetTopicPermissionsIn returns topic-permissions of user in virtual host.
func (c *Client) GetTopicPermissionsIn(vhost, username string) (rec []TopicPermissionInfo, err error) {
	req, err := newGETRequest(c, "topic-permissions/"+url.PathEscape(vhost)+"/"+url.PathEscape(username))
	if err != nil {
		return []TopicPermissionInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []TopicPermissionInfo{}, err
	}

	return rec, nil
}

//
// PUT /api/topic-permissions/{vhost}/{user}
//

// TopicPermissions represents a user's permissions on a topic.
type TopicPermissions struct {
	Exchange string `json:"exchange"`
	Write    string `json:"write"`
	Read     string `json:"read"`
}

// UpdateTopicPermissionsIn updates topic-permissions of user in virtual host.
func (c *Client) UpdateTopicPermissionsIn(vhost, username string, TopicPermissions TopicPermissions) (res *http.Response, err error) {
	body, err := json.Marshal(TopicPermissions)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "topic-permissions/"+url.PathEscape(vhost)+"/"+url.PathEscape(username), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/topic-permissions/{vhost}/{user}
//

// ClearTopicPermissionsIn clears (deletes) topic-permissions of user in virtual host.
func (c *Client) ClearTopicPermissionsIn(vhost, username string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "topic-permissions/"+url.PathEscape(vhost)+"/"+url.PathEscape(username), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteTopicPermissionsIn delete topic-permissions of exchange for user in virtual host.
func (c *Client) DeleteTopicPermissionsIn(vhost, username string, exchange string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "topic-permissions/"+url.PathEscape(vhost)+"/"+url.PathEscape(username)+"/"+url.PathEscape(exchange), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
