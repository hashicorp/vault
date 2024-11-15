package packngo

import "path"

// API documentation https://metal.equinix.com/developers/api#tag/Memberships
const membersBasePath = "/members"

// OrganizationService interface defines available organization methods
type MemberService interface {
	List(string, *ListOptions) ([]Member, *Response, error)
	Delete(string, string) (*Response, error)
}

type membersRoot struct {
	Members []Member `json:"members"`
	Meta    meta     `json:"meta"`
}

// Member is the returned from organization/id/members
type Member struct {
	*Href         `json:",inline"`
	ID            string       `json:"id"`
	Roles         []string     `json:"roles"`
	ProjectsCount int          `json:"projects_count"`
	User          User         `json:"user"`
	Organization  Organization `json:"organization"`
	Projects      []Project    `json:"projects"`
}

// MemberServiceOp implements MemberService
type MemberServiceOp struct {
	client *Client
}

var _ MemberService = (*MemberServiceOp)(nil)

// List returns the members in an organization
func (s *MemberServiceOp) List(organizationID string, opts *ListOptions) (orgs []Member, resp *Response, err error) {
	subset := new(membersRoot)
	endpointPath := path.Join(organizationBasePath, organizationID, membersBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		orgs = append(orgs, subset.Members...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// Delete removes the given member from the given organization
func (s *MemberServiceOp) Delete(organizationID, memberID string) (*Response, error) {
	if validateErr := ValidateUUID(organizationID); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(organizationBasePath, organizationID, membersBasePath, memberID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
