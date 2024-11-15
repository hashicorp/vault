package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// UserLimitsValues are properties used to modify virtual host limits (max-connections, max-channels)
type UserLimitsValues map[string]int

// UserLimits are properties used to delete virtual host limits (max-connections, max-channels)
type UserLimits []string

// UserLimitsInfo holds information about the current user limits
type UserLimitsInfo struct {
	User  string           `json:"user"`
	Value UserLimitsValues `json:"value"`
}

// GetAllUserLimits gets all users limits.
func (c *Client) GetAllUserLimits() (rec []UserLimitsInfo, err error) {
	req, err := newGETRequest(c, "user-limits")
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

// GetUserLimits gets a user limits.
func (c *Client) GetUserLimits(username string) (rec []UserLimitsInfo, err error) {
	req, err := newGETRequest(c, "user-limits/"+url.PathEscape(username))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

// PutUserLimits puts limits of a user.
func (c *Client) PutUserLimits(username string, limits UserLimitsValues) (res *http.Response, err error) {
	for limitName, limitValue := range limits {
		body, err := json.Marshal(struct {
			Value int `json:"value"`
		}{Value: limitValue})
		if err != nil {
			return nil, err
		}

		req, err := newRequestWithBody(c, "PUT", "user-limits/"+url.PathEscape(username)+"/"+limitName, body)
		if err != nil {
			return nil, err
		}

		if res, err = executeRequest(c, req); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// DeleteUserLimits deletes limits of a user.
func (c *Client) DeleteUserLimits(username string, limits UserLimits) (res *http.Response, err error) {
	for _, limit := range limits {
		req, err := newRequestWithBody(c, "DELETE", "user-limits/"+url.PathEscape(username)+"/"+limit, nil)
		if err != nil {
			return nil, err
		}

		if res, err = executeRequest(c, req); err != nil {
			return nil, err
		}
	}

	return res, nil
}
