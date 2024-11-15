package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// GlobalRuntimeParameter represents a vhost-scoped parameter.
// Value is interface{} to support creating parameters directly from types such as
// FederationUpstream and ShovelInfo.
type GlobalRuntimeParameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

//
// GET /api/global-parameters
//

// ListGlobalParameters returns a list of all global parameters.
func (c *Client) ListGlobalParameters() (params []GlobalRuntimeParameter, err error) {
	req, err := newGETRequest(c, "global-parameters")
	if err != nil {
		return []GlobalRuntimeParameter{}, err
	}

	if err = executeAndParseRequest(c, req, &params); err != nil {
		return []GlobalRuntimeParameter{}, err
	}

	return params, nil
}

//
// GET /api/global-parameters/name
//

// GetGlobalParameter returns information about a global parameter.
func (c *Client) GetGlobalParameter(name string) (p *GlobalRuntimeParameter, err error) {
	req, err := newGETRequest(c, "global-parameters/"+url.PathEscape(name))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &p); err != nil {
		return nil, err
	}

	return p, nil
}

//
// PUT /api/global-parameters/name/{value}
//

// PutRuntimeParameter creates or updates a runtime parameter.
func (c *Client) PutGlobalParameter(name string, value interface{}) (res *http.Response, err error) {
	p := GlobalRuntimeParameter{
		Name:  name,
		Value: value,
	}

	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "global-parameters/"+url.PathEscape(name), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/global-parameters/name
//

// DeleteRuntimeParameter removes a runtime parameter.
func (c *Client) DeleteGlobalParameter(name string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "global-parameters/"+url.PathEscape(name), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
