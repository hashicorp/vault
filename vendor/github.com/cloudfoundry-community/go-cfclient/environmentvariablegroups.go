package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type EnvironmentVariableGroup map[string]interface{}

func (c *Client) GetRunningEnvironmentVariableGroup() (EnvironmentVariableGroup, error) {
	return c.getEnvironmentVariableGroup(true)
}

func (c *Client) GetStagingEnvironmentVariableGroup() (EnvironmentVariableGroup, error) {
	return c.getEnvironmentVariableGroup(false)
}

func (c *Client) getEnvironmentVariableGroup(running bool) (EnvironmentVariableGroup, error) {
	evgType := "staging"
	if running {
		evgType = "running"
	}

	req := c.NewRequest("GET", fmt.Sprintf("/v2/config/environment_variable_groups/%s", evgType))
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	evg := EnvironmentVariableGroup{}
	err = json.NewDecoder(resp.Body).Decode(&evg)
	return evg, err
}

func (c *Client) SetRunningEnvironmentVariableGroup(evg EnvironmentVariableGroup) error {
	return c.setEnvironmentVariableGroup(evg, true)
}

func (c *Client) SetStagingEnvironmentVariableGroup(evg EnvironmentVariableGroup) error {
	return c.setEnvironmentVariableGroup(evg, false)
}

func (c *Client) setEnvironmentVariableGroup(evg EnvironmentVariableGroup, running bool) error {
	evgType := "staging"
	if running {
		evgType = "running"
	}

	marshalled, err := json.Marshal(evg)
	if err != nil {
		return err
	}

	req := c.NewRequestWithBody("PUT", fmt.Sprintf("/v2/config/environment_variable_groups/%s", evgType), bytes.NewBuffer(marshalled))
	_, err = c.DoRequest(req)
	return err
}
