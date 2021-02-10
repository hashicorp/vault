package packngo

import (
	"path"
)

var bgpSessionBasePath = "/bgp/sessions"
var bgpNeighborsBasePath = "/bgp/neighbors"

// BGPSessionService interface defines available BGP session methods
type BGPSessionService interface {
	Get(string, *GetOptions) (*BGPSession, *Response, error)
	Create(string, CreateBGPSessionRequest) (*BGPSession, *Response, error)
	Delete(string) (*Response, error)
}

type bgpSessionsRoot struct {
	Sessions []BGPSession `json:"bgp_sessions"`
	Meta     meta         `json:"meta"`
}

// BGPSessionServiceOp implements BgpSessionService
type BGPSessionServiceOp struct {
	client *Client
}

// BGPSession represents an Equinix Metal BGP Session
type BGPSession struct {
	ID            string   `json:"id,omitempty"`
	Status        string   `json:"status,omitempty"`
	LearnedRoutes []string `json:"learned_routes,omitempty"`
	AddressFamily string   `json:"address_family,omitempty"`
	Device        Device   `json:"device,omitempty"`
	Href          string   `json:"href,omitempty"`
	DefaultRoute  *bool    `json:"default_route,omitempty"`
}

type bgpNeighborsRoot struct {
	BGPNeighbors []BGPNeighbor `json:"bgp_neighbors"`
}

// BGPNeighor is struct for listing BGP neighbors of a device
type BGPNeighbor struct {
	AddressFamily int        `json:"address_family"`
	CustomerAs    int        `json:"customer_as"`
	CustomerIP    string     `json:"customer_ip"`
	Md5Enabled    bool       `json:"md5_enabled"`
	Md5Password   string     `json:"md5_password"`
	Multihop      bool       `json:"multihop"`
	PeerAs        int        `json:"peer_as"`
	PeerIps       []string   `json:"peer_ips"`
	RoutesIn      []BGPRoute `json:"routes_in"`
	RoutesOut     []BGPRoute `json:"routes_out"`
}

// BGPRoute is a struct for Route in BGP neighbor listing
type BGPRoute struct {
	Route string `json:"route"`
	Exact bool   `json:"exact"`
}

// CreateBGPSessionRequest struct
type CreateBGPSessionRequest struct {
	AddressFamily string `json:"address_family"`
	DefaultRoute  *bool  `json:"default_route,omitempty"`
}

// Create function
func (s *BGPSessionServiceOp) Create(deviceID string, request CreateBGPSessionRequest) (*BGPSession, *Response, error) {
	apiPath := path.Join(deviceBasePath, deviceID, bgpSessionBasePath)
	session := new(BGPSession)

	resp, err := s.client.DoRequest("POST", apiPath, request, session)
	if err != nil {
		return nil, resp, err
	}

	return session, resp, err
}

// Delete function
func (s *BGPSessionServiceOp) Delete(id string) (*Response, error) {
	apiPath := path.Join(bgpSessionBasePath, id)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}

// Get function
func (s *BGPSessionServiceOp) Get(id string, opts *GetOptions) (session *BGPSession, response *Response, err error) {
	endpointPath := path.Join(bgpSessionBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	session = new(BGPSession)
	response, err = s.client.DoRequest("GET", apiPathQuery, nil, session)
	if err != nil {
		return nil, response, err
	}

	return session, response, err
}
