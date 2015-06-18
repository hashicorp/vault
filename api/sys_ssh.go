package api

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func (c *Sys) Ssh(target string) (*OneTimeKey, error) {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/ssh/connect"))
	input := strings.Split(target, "@")
	username := input[0]
	ipAddr := input[1]
	ip4Addr, err := net.ResolveIPAddr("ip4", ipAddr)
	log.Printf("Vishal: ssh.Ssh ipAddr_resolved: %#v\n", ip4Addr.String())
	body := map[string]interface{}{
		"username": username,
		"address":  ip4Addr.String(),
	}
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result OneTimeKey
	err = resp.DecodeJSON(&result)
	return &result, err
}

type OneTimeKey struct {
	Key string
}
