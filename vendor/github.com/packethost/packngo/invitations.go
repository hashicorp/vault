package packngo

import (
	"path"
)

const (
	invitationsBasePath = "/invitations"
)

// InvitationService interface defines available invitation methods
type InvitationService interface {
	Create(string, *InvitationCreateRequest, *GetOptions) (*Invitation, *Response, error)
	List(string, *ListOptions) ([]Invitation, *Response, error)
	Get(string, *GetOptions) (*Invitation, *Response, error)
	Accept(string, *InvitationUpdateRequest) (*Invitation, *Response, error)
	Resend(string) (*Invitation, *Response, error)
	Delete(string) (*Response, error)
}

type invitationsRoot struct {
	Invitations []Invitation `json:"invitations"`
	Meta        meta         `json:"meta"`
}

// Invitation represents an Equinix Metal invitation
type Invitation struct {
	*Href        `json:",inline"`
	ID           string     `json:"id,omitempty"`
	Invitation   Href       `json:"invitation,omitempty"`
	InvitedBy    Href       `json:"invited_by,omitempty"`
	Invitee      string     `json:"invitee,omitempty"`
	Nonce        string     `json:"nonce,omitempty"`
	Organization Href       `json:"organization,omitempty"`
	Projects     []Href     `json:"projects,omitempty"`
	Roles        []string   `json:"roles,omitempty"`
	CreatedAt    *Timestamp `json:"created_at,omitempty"`
	UpdatedAt    *Timestamp `json:"updated_at,omitempty"`
}

// InvitationCreateRequest struct for InvitationService.Create
type InvitationCreateRequest struct {
	// Invitee is the email address of the recipient
	Invitee     string   `json:"invitee"`
	Message     string   `json:"message,omitempty"`
	ProjectsIDs []string `json:"projects_ids,omitempty"`
	Roles       []string `json:"roles,omitempty"`
}

// InvitationUpdateRequest struct for InvitationService.Update
type InvitationUpdateRequest struct{}

func (u Invitation) String() string {
	return Stringify(u)
}

// InvitationServiceOp implements InvitationService
type InvitationServiceOp struct {
	client *Client
}

var _ InvitationService = (*InvitationServiceOp)(nil)

// Lists open invitations to the project
func (s *InvitationServiceOp) List(organizationID string, opts *ListOptions) (invitations []Invitation, resp *Response, err error) {
	endpointPath := path.Join(organizationBasePath, organizationID, invitationsBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(invitationsRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		invitations = append(invitations, subset.Invitations...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// Create a Invitation with the given InvitationCreateRequest. New invitation VerificationStage
// will be AccountCreated, unless InvitationCreateRequest contains an valid
// InvitationID and Nonce in which case the VerificationStage will be Verified.
func (s *InvitationServiceOp) Create(organizationID string, createRequest *InvitationCreateRequest, opts *GetOptions) (*Invitation, *Response, error) {
	endpointPath := path.Join(organizationBasePath, organizationID, invitationsBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	invitation := new(Invitation)

	resp, err := s.client.DoRequest("POST", apiPathQuery, createRequest, invitation)
	if err != nil {
		return nil, resp, err
	}

	return invitation, resp, err
}

func (s *InvitationServiceOp) Get(invitationID string, opts *GetOptions) (*Invitation, *Response, error) {
	if validateErr := ValidateUUID(invitationID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(invitationsBasePath, invitationID)
	apiPathQuery := opts.WithQuery(endpointPath)
	invitation := new(Invitation)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, invitation)
	if err != nil {
		return nil, resp, err
	}

	return invitation, resp, err
}

// Update updates the current invitation
func (s *InvitationServiceOp) Delete(id string) (*Response, error) {
	opts := &GetOptions{}
	endpointPath := path.Join(invitationsBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)

	return s.client.DoRequest("DELETE", apiPathQuery, nil, nil)
}

// Update updates the current invitation
func (s *InvitationServiceOp) Accept(id string, updateRequest *InvitationUpdateRequest) (*Invitation, *Response, error) {
	opts := &GetOptions{}
	endpointPath := path.Join(invitationsBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	invitation := new(Invitation)

	resp, err := s.client.DoRequest("PUT", apiPathQuery, updateRequest, invitation)
	if err != nil {
		return nil, resp, err
	}

	return invitation, resp, err
}

// Update updates the current invitation
func (s *InvitationServiceOp) Resend(id string) (*Invitation, *Response, error) {
	opts := &GetOptions{}
	endpointPath := path.Join(invitationsBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	invitation := new(Invitation)

	resp, err := s.client.DoRequest("POST", apiPathQuery, nil, invitation)
	if err != nil {
		return nil, resp, err
	}

	return invitation, resp, err
}
