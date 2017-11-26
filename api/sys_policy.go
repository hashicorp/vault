package api

import (
	"context"
	"fmt"
	"net/http"
)

// ListPolicies returns a string list of all policies.
func (c *Sys) ListPolicies() ([]string, error) {
	return c.ListPoliciesWithContext(context.Background())
}

// ListPoliciesWithContext returns a string list of all policies, with a
// context.
func (c *Sys) ListPoliciesWithContext(ctx context.Context) ([]string, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/policy")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return nil, err
	}

	var ok bool
	if _, ok = result["policies"]; !ok {
		return nil, fmt.Errorf("policies not found in response")
	}

	listRaw := result["policies"].([]interface{})
	var policies []string

	for _, val := range listRaw {
		policies = append(policies, val.(string))
	}

	return policies, err
}

// GetPolicy retrieves the contents of the given policy by name.
func (c *Sys) GetPolicy(name string) (string, error) {
	return c.GetPolicyWithContext(context.Background(), name)
}

// GetPolicyWithContext retrieves the contents of the given policy by name, with
// a context.
func (c *Sys) GetPolicyWithContext(ctx context.Context, name string) (string, error) {
	req := c.c.NewRequest(http.MethodGet, fmt.Sprintf("/v1/sys/policy/%s", name))
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return "", nil
		}
	}
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return "", err
	}

	if rulesRaw, ok := result["rules"]; ok {
		return rulesRaw.(string), nil
	}
	if policyRaw, ok := result["policy"]; ok {
		return policyRaw.(string), nil
	}

	return "", fmt.Errorf("no policy found in response")
}

// PutPolicy creates a new or updates an existing policy.
func (c *Sys) PutPolicy(name, rules string) error {
	return c.PutPolicyWithContext(context.Background(), name, rules)
}

// PutPolicyWithContext creates a new or updates an existing policy, with a
// context.
func (c *Sys) PutPolicyWithContext(ctx context.Context, name, rules string) error {
	body := map[string]string{
		"rules": rules,
	}

	req := c.c.NewRequest(http.MethodPut, fmt.Sprintf("/v1/sys/policy/%s", name))
	req = req.WithContext(ctx)
	if err := req.SetJSONBody(body); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeletePolicy deletes the policy by the given name, if it exists.
func (c *Sys) DeletePolicy(name string) error {
	return c.DeletePolicyWithContext(context.Background(), name)
}

// DeletePolicyWithContext deletes the policy by the given name, if it exists,
// with a context.
func (c *Sys) DeletePolicyWithContext(ctx context.Context, name string) error {
	req := c.c.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/sys/policy/%s", name))
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

type getPoliciesResp struct {
	Rules string `json:"rules"`
}

type listPoliciesResp struct {
	Policies []string `json:"policies"`
}
