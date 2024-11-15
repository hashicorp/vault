package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SecurityGroup definition
type SecurityGroup struct {
	Description           string      `json:"description"`
	ID                    string      `json:"id"`
	Organization          string      `json:"organization"`
	Name                  string      `json:"name"`
	Servers               []ServerRef `json:"servers"`
	EnableDefaultSecurity bool        `json:"enable_default_security"`
	OrganizationDefault   bool        `json:"organization_default"`
}

type SecurityGroupRef struct {
	Identifier string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
}

type ServerRef struct {
	Identifier string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
}

// getSecurityGroups represents the response of a GET /security_groups/
type getSecurityGroups struct {
	SecurityGroups []SecurityGroup `json:"security_groups"`
}

// getSecurityGroup represents the response of a GET /security_groups/{groupID}
type getSecurityGroup struct {
	SecurityGroup SecurityGroup `json:"security_group"`
}

// NewSecurityGroup definition POST request /security_groups
type NewSecurityGroup struct {
	Organization          string `json:"organization"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	EnableDefaultSecurity bool   `json:"enable_default_security"`
}

// UpdateSecurityGroup definition PUT request /security_groups
type UpdateSecurityGroup struct {
	Organization        string `json:"organization"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	OrganizationDefault bool   `json:"organization_default"`
}

// DeleteSecurityGroup deletes a SecurityGroup
func (s *API) DeleteSecurityGroup(securityGroupID string) error {
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("security_groups/%s", securityGroupID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusNoContent}, resp)
	return err
}

// UpdateSecurityGroup updates a SecurityGroup
func (s *API) UpdateSecurityGroup(group UpdateSecurityGroup, securityGroupID string) (*SecurityGroup, error) {
	resp, err := s.PutResponse(s.computeAPI, fmt.Sprintf("security_groups/%s", securityGroupID), group)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var data getSecurityGroup
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data.SecurityGroup, err
}

// GetSecurityGroup returns a SecurityGroup
func (s *API) GetSecurityGroup(groupsID string) (*SecurityGroup, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, fmt.Sprintf("security_groups/%s", groupsID), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroups getSecurityGroup

	if err = json.Unmarshal(body, &securityGroups); err != nil {
		return nil, err
	}
	return &securityGroups.SecurityGroup, nil
}

// CreateSecurityGroup posts a group on a server
func (s *API) CreateSecurityGroup(group NewSecurityGroup) (*SecurityGroup, error) {
	resp, err := s.PostResponse(s.computeAPI, "security_groups", group)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroups getSecurityGroup

	if err = json.Unmarshal(body, &securityGroups); err != nil {
		return nil, err
	}
	return &securityGroups.SecurityGroup, nil
}

// GetSecurityGroups returns a SecurityGroups
func (s *API) GetSecurityGroups() ([]SecurityGroup, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "security_groups", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroups getSecurityGroups

	if err = json.Unmarshal(body, &securityGroups); err != nil {
		return nil, err
	}
	return securityGroups.SecurityGroups, nil
}
