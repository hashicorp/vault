package packngo

import (
	"path"
)

type ConnectionRedundancy string
type ConnectionType string
type ConnectionPortRole string

const (
	connectionBasePath                           = "/connections"
	ConnectionShared        ConnectionType       = "shared"
	ConnectionDedicated     ConnectionType       = "dedicated"
	ConnectionRedundant     ConnectionRedundancy = "redundant"
	ConnectionPrimary       ConnectionRedundancy = "primary"
	ConnectionPortPrimary   ConnectionPortRole   = "primary"
	ConnectionPortSecondary ConnectionPortRole   = "secondary"
)

type ConnectionService interface {
	OrganizationCreate(string, *ConnectionCreateRequest) (*Connection, *Response, error)
	ProjectCreate(string, *ConnectionCreateRequest) (*Connection, *Response, error)
	OrganizationList(string, *GetOptions) ([]Connection, *Response, error)
	ProjectList(string, *GetOptions) ([]Connection, *Response, error)
	Delete(string) (*Response, error)
	Get(string, *GetOptions) (*Connection, *Response, error)
	Events(string, *GetOptions) ([]Event, *Response, error)
	PortEvents(string, string, *GetOptions) ([]Event, *Response, error)
	Ports(string, *GetOptions) ([]ConnectionPort, *Response, error)
	Port(string, string, *GetOptions) (*ConnectionPort, *Response, error)
	VirtualCircuits(string, string, *GetOptions) ([]VirtualCircuit, *Response, error)
}

type ConnectionServiceOp struct {
	client *Client
}

type connectionPortsRoot struct {
	Ports []ConnectionPort `json:"ports"`
}

type connectionsRoot struct {
	Connections []Connection `json:"interconnections"`
	Meta        meta         `json:"meta"`
}

type ConnectionPort struct {
	ID              string             `json:"id"`
	Name            string             `json:"name,omitempty"`
	Status          string             `json:"status,omitempty"`
	Role            ConnectionPortRole `json:"role,omitempty"`
	Speed           int                `json:"speed,omitempty"`
	Organization    *Organization      `json:"organization,omitempty"`
	VirtualCircuits []VirtualCircuit   `json:"virtual_circuits,omitempty"`
	LinkStatus      string             `json:"link_status,omitempty"`
	Href            string             `json:"href,omitempty"`
}

type Connection struct {
	ID           string               `json:"id"`
	Name         string               `json:"name,omitempty"`
	Status       string               `json:"status,omitempty"`
	Redundancy   ConnectionRedundancy `json:"redundancy,omitempty"`
	Facility     *Facility            `json:"facility,omitempty"`
	Type         ConnectionType       `json:"type,omitempty"`
	Description  string               `json:"description,omitempty"`
	Project      *Project             `json:"project,omitempty"`
	Organization *Organization        `json:"organization,omitempty"`
	Speed        int                  `json:"speed,omitempty"`
	Token        string               `json:"token,omitempty"`
	Tags         []string             `json:"tags,omitempty"`
	Ports        []ConnectionPort     `json:"ports,omitempty"`
}

type ConnectionCreateRequest struct {
	Name        string               `json:"name,omitempty"`
	Redundancy  ConnectionRedundancy `json:"redundancy,omitempty"`
	Facility    string               `json:"facility,omitempty"`
	Type        ConnectionType       `json:"type,omitempty"`
	Description *string              `json:"description,omitempty"`
	Project     string               `json:"project,omitempty"`
	Speed       int                  `json:"speed,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
}

func (c *Connection) PortByRole(r ConnectionPortRole) *ConnectionPort {
	for _, p := range c.Ports {
		if p.Role == r {
			return &p
		}
	}
	return nil
}

func (s *ConnectionServiceOp) create(apiUrl string, createRequest *ConnectionCreateRequest) (*Connection, *Response, error) {
	connection := new(Connection)
	resp, err := s.client.DoRequest("POST", apiUrl, createRequest, connection)
	if err != nil {
		return nil, resp, err
	}

	return connection, resp, err
}

func (s *ConnectionServiceOp) OrganizationCreate(id string, createRequest *ConnectionCreateRequest) (*Connection, *Response, error) {
	apiUrl := path.Join(organizationBasePath, id, connectionBasePath)
	return s.create(apiUrl, createRequest)
}

func (s *ConnectionServiceOp) ProjectCreate(id string, createRequest *ConnectionCreateRequest) (*Connection, *Response, error) {
	apiUrl := path.Join(projectBasePath, id, connectionBasePath)
	return s.create(apiUrl, createRequest)
}

func (s *ConnectionServiceOp) list(url string, opts *GetOptions) (connections []Connection, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(url)

	for {
		subset := new(connectionsRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		connections = append(connections, subset.Connections...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}

}

func (s *ConnectionServiceOp) OrganizationList(id string, opts *GetOptions) ([]Connection, *Response, error) {
	apiUrl := path.Join(organizationBasePath, id, connectionBasePath)
	return s.list(apiUrl, opts)
}

func (s *ConnectionServiceOp) ProjectList(id string, opts *GetOptions) ([]Connection, *Response, error) {
	apiUrl := path.Join(projectBasePath, id, connectionBasePath)
	return s.list(apiUrl, opts)
}

func (s *ConnectionServiceOp) Delete(id string) (*Response, error) {
	apiPath := path.Join(connectionBasePath, id)
	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}

func (s *ConnectionServiceOp) Port(connID, portID string, opts *GetOptions) (*ConnectionPort, *Response, error) {
	endpointPath := path.Join(connectionBasePath, connID, portBasePath, portID)
	apiPathQuery := opts.WithQuery(endpointPath)
	port := new(ConnectionPort)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, port)
	if err != nil {
		return nil, resp, err
	}
	return port, resp, err
}

func (s *ConnectionServiceOp) Get(id string, opts *GetOptions) (*Connection, *Response, error) {
	endpointPath := path.Join(connectionBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	connection := new(Connection)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, connection)
	if err != nil {
		return nil, resp, err
	}
	return connection, resp, err
}

func (s *ConnectionServiceOp) Ports(connID string, opts *GetOptions) ([]ConnectionPort, *Response, error) {
	endpointPath := path.Join(connectionBasePath, connID, portBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	ports := new(connectionPortsRoot)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, ports)
	if err != nil {
		return nil, resp, err
	}
	return ports.Ports, resp, nil

}

func (s *ConnectionServiceOp) Events(id string, opts *GetOptions) ([]Event, *Response, error) {
	apiPath := path.Join(connectionBasePath, id, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *ConnectionServiceOp) PortEvents(connID, portID string, opts *GetOptions) ([]Event, *Response, error) {
	apiPath := path.Join(connectionBasePath, connID, portBasePath, portID, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *ConnectionServiceOp) VirtualCircuits(connID, portID string, opts *GetOptions) (vcs []VirtualCircuit, resp *Response, err error) {
	endpointPath := path.Join(connectionBasePath, connID, portBasePath, portID, virtualCircuitBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	for {
		subset := new(virtualCircuitsRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		vcs = append(vcs, subset.VirtualCircuits...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}
