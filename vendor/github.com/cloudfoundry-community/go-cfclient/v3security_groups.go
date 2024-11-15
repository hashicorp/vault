package cfclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// V3SecurityGroup implements the security group object. Security groups are collections of egress traffic rules that can be applied to the staging or running state of applications.
type V3SecurityGroup struct {
	Name            string                           `json:"name,omitempty"`
	GUID            string                           `json:"guid,omitempty"`
	CreatedAt       string                           `json:"created_at,omitempty"`
	UpdatedAt       string                           `json:"updated_at,omitempty"`
	GloballyEnabled V3GloballyEnabled                `json:"globally_enabled,omitempty"`
	Rules           []V3Rule                         `json:"rules,omitempty"`
	Relationships   map[string]V3ToManyRelationships `json:"relationships,omitempty"`
	Links           map[string]Link                  `json:"links,omitempty"`
}

// V3GloballyEnabled object controls if the group is applied globally to the lifecycle of all applications
type V3GloballyEnabled struct {
	Running bool `json:"running,omitempty"`
	Staging bool `json:"staging,omitempty"`
}

// V3Rule is an object that provide a rule that will be applied by a security group
type V3Rule struct {
	Protocol    string `json:"protocol,omitempty"`
	Destination string `json:"destination,omitempty"`
	Ports       string `json:"ports,omitempty"`
	Type        *int   `json:"type,omitempty"`
	Code        *int   `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Log         bool   `json:"log,omitempty"`
}

type listV3SecurityGroupResponse struct {
	Pagination Pagination        `json:"pagination,omitempty"`
	Resources  []V3SecurityGroup `json:"resources,omitempty"`
}

// ListV3SecurityGroupsByQuery retrieves security groups based on query
func (c *Client) ListV3SecurityGroupsByQuery(query url.Values) ([]V3SecurityGroup, error) {
	var securityGroups []V3SecurityGroup
	requestURL, err := url.Parse("/v3/security_groups")
	if err != nil {
		return nil, err
	}
	requestURL.RawQuery = query.Encode()

	for {
		r := c.NewRequest("GET", fmt.Sprintf("%s?%s", requestURL.Path, requestURL.RawQuery))
		resp, err := c.DoRequest(r)
		if err != nil {
			return nil, errors.Wrap(err, "Error requesting v3 security groups")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error listing v3 security groups, response code: %d", resp.StatusCode)
		}

		var data listV3SecurityGroupResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, errors.Wrap(err, "Error parsing JSON from list v3 security groups")
		}

		securityGroups = append(securityGroups, data.Resources...)

		requestURL, err = url.Parse(data.Pagination.Next.Href)
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing next page URL")
		}
		if requestURL.String() == "" {
			break
		}
	}

	return securityGroups, nil
}

// CreateV3SecurityGroupRequest implements an object that is passed to CreateV3SecurityGroup method
type CreateV3SecurityGroupRequest struct {
	Name            string                           `json:"name"`
	GloballyEnabled *V3GloballyEnabled               `json:"globally_enabled,omitempty"`
	Rules           []*V3Rule                        `json:"rules,omitempty"`
	Relationships   map[string]V3ToManyRelationships `json:"relationships,omitempty"`
}

// CreateV3SecurityGroup creates security group from CreateV3SecurityGroupRequest
func (c *Client) CreateV3SecurityGroup(r CreateV3SecurityGroupRequest) (*V3SecurityGroup, error) {
	req := c.NewRequest("POST", "/v3/security_groups")

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	req.body = buf

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating v3 security group")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error creating v3 security group %s, response code: %d", r.Name, resp.StatusCode)
	}

	var securitygroup V3SecurityGroup
	if err := json.NewDecoder(resp.Body).Decode(&securitygroup); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 security group JSON")
	}

	return &securitygroup, nil
}

// DeleteV3SecurityGroup deletes security group by GUID
func (c *Client) DeleteV3SecurityGroup(GUID string) error {
	req := c.NewRequest("DELETE", "/v3/security_groups/"+GUID)

	resp, err := c.DoRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error while deleting v3 security group")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Error deleting v3 security group with GUID [%s], response code: %d", GUID, resp.StatusCode)
	}
	return nil
}

// UpdateV3SecurityGroupRequest implements an object that is passed to UpdateV3SecurityGroup method
type UpdateV3SecurityGroupRequest struct {
	Name            string             `json:"name,omitempty"`
	GloballyEnabled *V3GloballyEnabled `json:"globally_enabled,omitempty"`
	Rules           []*V3Rule          `json:"rules,omitempty"`
}

// UpdateV3SecurityGroup updates security group by GUID and from UpdateV3SecurityGroupRequest
func (c *Client) UpdateV3SecurityGroup(GUID string, r UpdateV3SecurityGroupRequest) (*V3SecurityGroup, error) {
	req := c.NewRequest("PATCH", "/v3/security_groups/"+GUID)
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(r); err != nil {
		return nil, err
	}
	req.body = buf

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while updating v3 security group")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error updating v3 security group %s, response code: %d", GUID, resp.StatusCode)
	}

	var securityGroup V3SecurityGroup
	if err := json.NewDecoder(resp.Body).Decode(&securityGroup); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 security group JSON")
	}

	return &securityGroup, nil
}

// GetV3SecurityGroupByGUID retrieves security group base on provided GUID
func (c *Client) GetV3SecurityGroupByGUID(GUID string) (*V3SecurityGroup, error) {
	req := c.NewRequest("GET", "/v3/security_groups/"+GUID)

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting v3 security group")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting v3 security group with GUID [%s], response code: %d", GUID, resp.StatusCode)
	}

	var securityGroup V3SecurityGroup
	if err := json.NewDecoder(resp.Body).Decode(&securityGroup); err != nil {
		return nil, errors.Wrap(err, "Error reading v3 security group JSON")
	}

	return &securityGroup, nil
}
