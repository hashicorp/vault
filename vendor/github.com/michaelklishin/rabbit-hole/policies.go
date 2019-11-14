package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// PolicyDefinition is a map of additional arguments
// added to the entities (queues, exchanges or both)
// that match a policy.
type PolicyDefinition map[string]interface{}

// NodeNames is a list of node names.
type NodeNames []string

// Policy represents a configured policy.
type Policy struct {
	// Virtual host this policy is in.
	Vhost string `json:"vhost"`
	// Regular expression pattern used to match queues and exchanges,
	// , e.g. "^ha\..+"
	Pattern string `json:"pattern"`
	// What this policy applies to: "queues", "exchanges", etc.
	ApplyTo  string `json:"apply-to"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	// Additional arguments added to the entities (queues,
	// exchanges or both) that match a policy
	Definition PolicyDefinition `json:"definition"`
}

//
// GET /api/policies
//

// ListPolicies returns all policies (across all virtual hosts).
func (c *Client) ListPolicies() (rec []Policy, err error) {
	req, err := newGETRequest(c, "policies")
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// GET /api/policies/{vhost}
//

// ListPoliciesIn returns policies in a specific virtual host.
func (c *Client) ListPoliciesIn(vhost string) (rec []Policy, err error) {
	req, err := newGETRequest(c, "policies/"+url.PathEscape(vhost))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// GET /api/policies/{vhost}/{name}
//

// GetPolicy returns individual policy in virtual host.
func (c *Client) GetPolicy(vhost, name string) (rec *Policy, err error) {
	req, err := newGETRequest(c, "policies/"+url.PathEscape(vhost)+"/"+url.PathEscape(name))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// PUT /api/policies/{vhost}/{name}
//

// PutPolicy creates or updates a policy.
func (c *Client) PutPolicy(vhost string, name string, policy Policy) (res *http.Response, err error) {
	body, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "policies/"+url.PathEscape(vhost)+"/"+url.PathEscape(name), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/policies/{vhost}/{name}
//

// DeletePolicy deletes a policy.
func (c *Client) DeletePolicy(vhost, name string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "policies/"+url.PathEscape(vhost)+"/"+url.PathEscape(name), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
