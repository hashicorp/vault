package api

import "fmt"

func (c *Sys) Ssh(data map[string]interface{}) (*Secret, error) {
	r := c.c.NewRequest("PUT", fmt.Sprintf("/v1/ssh/creds/web"))
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)

	/*
		result := new(Secret)
		err = resp.DecodeJSON(&result)
		log.Printf("Vishal: api.sys_ssh.Ssh: result:%#v\n", result.Data)

		var oneTimeKey OneTimeKey
		err = result.Data.DecodeJSON(&oneTimeKey)
		log.Printf("Vishal: oneTimeKey:%#v\n", oneTimeKey)
		return &oneTimeKey, err
	*/
	//return result, err
}
