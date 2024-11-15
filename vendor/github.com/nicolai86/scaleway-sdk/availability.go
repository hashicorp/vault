package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type InstanceTypeAvailability string

var (
	InstanceTypeAvailable InstanceTypeAvailability = "available"
	InstanceTypeScarce    InstanceTypeAvailability = "scarce"
	InstanceTypeShortage  InstanceTypeAvailability = "shortage"
)

type ServerAvailability struct {
	Availability InstanceTypeAvailability `json:"availability"`
}

type ServerAvailabilities map[string]ServerAvailability

func (a ServerAvailabilities) CommercialTypes() []string {
	types := []string{}
	for k, _ := range a {
		types = append(types, k)
	}
	return types
}

type availabilityResponse struct {
	Servers ServerAvailabilities
}

func (s *API) GetServerAvailabilities() (ServerAvailabilities, error) {
	resp, err := s.response("GET", fmt.Sprintf("%s/products/servers/availability", s.computeAPI), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	content := availabilityResponse{}
	if err := json.Unmarshal(bs, &content); err != nil {
		return nil, err
	}
	return content.Servers, nil
}
