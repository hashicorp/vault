package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// OperatorPolicy represents a configured policy.
type OperatorPolicy struct {
	// Virtual host this policy is in.
	Vhost string `json:"vhost"`
	// Regular expression pattern used to match queues,
	// , e.g. "^ha\..+"
	Pattern string `json:"pattern"`
	// What this policy applies to: "queues".
	ApplyTo  string `json:"apply-to"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	// Additional arguments added to the queues that match a operator policy
	Definition PolicyDefinition `json:"definition"`
}

//
// GET /api/operator-policies
//

// ListOperatorPolicies returns all operator policies (across all virtual hosts).
func (c *Client) ListOperatorPolicies() (rec []OperatorPolicy, err error) {
	req, err := newGETRequest(c, "operator-policies")
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// GET /api/operator-policies/{vhost}
//

// ListOperatorPoliciesIn returns operator policies in a specific virtual host.
func (c *Client) ListOperatorPoliciesIn(vhost string) (rec []OperatorPolicy, err error) {
	req, err := newGETRequest(c, "operator-policies/"+url.PathEscape(vhost))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// GET /api/operator-policies/{vhost}/{name}
//

// GetOperatorPolicy returns individual operator policy in virtual host.
func (c *Client) GetOperatorPolicy(vhost, name string) (rec *OperatorPolicy, err error) {
	req, err := newGETRequest(c, "operator-policies/"+url.PathEscape(vhost)+"/"+url.PathEscape(name))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// PUT /api/operator-policies/{vhost}/{name}
//

// PutOperatorPolicy creates or updates an operator policy.
func (c *Client) PutOperatorPolicy(vhost string, name string, operatorPolicy OperatorPolicy) (res *http.Response, err error) {
	body, err := json.Marshal(operatorPolicy)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "operator-policies/"+url.PathEscape(vhost)+"/"+url.PathEscape(name), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/operator-policies/{vhost}/{name}
//

// DeleteOperatorPolicy deletes an operator policy.
func (c *Client) DeleteOperatorPolicy(vhost, name string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "operator-policies/"+url.PathEscape(vhost)+"/"+url.PathEscape(name), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
