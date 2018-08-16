package api

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

func (c *Sys) GetGithubTeamMap(name string) (string, error) {
	r := c.c.NewRequest("GET", fmt.Sprintf("/v1/auth/github/map/teams/%s", name))
	resp, err := c.c.RawRequest(r)
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	 panic(err.Error())
   	}
	bodyString := string(body)

	policy := gjson.Get(bodyString, "data.value")	
	return policy.String(), nil
}

func (c *Sys) PostGithubTeamMap(name, policies string) error {
	body := map[string]string{
		"value": policies,
	}

	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/auth/github/map/teams/%s", name))
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


func (c *Sys) DeleteGithubTeamMap(name string) error {
	r := c.c.NewRequest("DELETE", fmt.Sprintf("/v1/auth/github/map/teams/%s", name))
	resp, err := c.c.RawRequest(r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}