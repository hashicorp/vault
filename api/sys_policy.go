package api

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func (c *Sys) ListPolicies() ([]string, error) {
	r := c.c.NewRequest("GET", "/v1/sys/policy")
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}

	if secret == nil || secret.Data == nil || len(secret.Data) == 0 {
		return nil, nil
	}

	var result listPoliciesResp
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return result.Policies, err
}

func (c *Sys) GetPolicy(name string) (string, error) {
	r := c.c.NewRequest("GET", fmt.Sprintf("/v1/sys/policy/%s", name))
	resp, err := c.c.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return "", nil
		}
	}
	if err != nil {
		return "", err
	}

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return "", err
	}

	if secret == nil || secret.Data == nil || len(secret.Data) == 0 {
		return "", nil
	}

	var result getPoliciesResp
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return "", err
	}

	return result.Rules, err
}

func (c *Sys) PutPolicy(name, rules string) error {
	body := map[string]string{
		"rules": rules,
	}

	r := c.c.NewRequest("PUT", fmt.Sprintf("/v1/sys/policy/%s", name))
	if err := r.SetJSONBody(body); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Sys) DeletePolicy(name string) error {
	r := c.c.NewRequest("DELETE", fmt.Sprintf("/v1/sys/policy/%s", name))
	resp, err := c.c.RawRequest(r)
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
