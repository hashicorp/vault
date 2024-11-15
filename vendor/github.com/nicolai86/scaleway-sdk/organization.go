package api

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// Organization represents a  Organization
type Organization struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Users []User `json:"users"`
}

// organizationsDefinition represents a  Organizations
type organizationsDefinition struct {
	Organizations []Organization `json:"organizations"`
}

// GetOrganization returns Organization
func (s *API) GetOrganization() ([]Organization, error) {
	resp, err := s.GetResponsePaginate(AccountAPI, "organizations", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var data organizationsDefinition

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data.Organizations, nil
}
