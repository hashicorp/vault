package packngo

import (
	"errors"
	"fmt"
	"path"
	"time"
)

type ConnectionRedundancy string
type ConnectionType string
type ConnectionPortRole string
type ConnectionMode string

const (
	connectionBasePath                           = "/connections"
	ConnectionShared        ConnectionType       = "shared"
	ConnectionDedicated     ConnectionType       = "dedicated"
	ConnectionRedundant     ConnectionRedundancy = "redundant"
	ConnectionPrimary       ConnectionRedundancy = "primary"
	ConnectionPortPrimary   ConnectionPortRole   = "primary"
	ConnectionPortSecondary ConnectionPortRole   = "secondary"
	ConnectionModeStandard  ConnectionMode       = "standard"
	ConnectionModeTunnel    ConnectionMode       = "tunnel"
	ConnectionDeleteTimeout                      = 60 * time.Second
	ConnectionDeleteCheck                        = 2 * time.Second
)

type ConnectionService interface {
	OrganizationCreate(string, *ConnectionCreateRequest) (*Connection, *Response, error)
	ProjectCreate(string, *ConnectionCreateRequest) (*Connection, *Response, error)
	Update(string, *ConnectionUpdateRequest, *GetOptions) (*Connection, *Response, error)
	OrganizationList(string, *GetOptions) ([]Connection, *Response, error)
	ProjectList(string, *GetOptions) ([]Connection, *Response, error)
	Delete(string, bool) (*Response, error)
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
	*Href        `json:",inline"`
	ID           string             `json:"id"`
	LinkStatus   string             `json:"link_status,omitempty"`
	Name         string             `json:"name,omitempty"`
	Organization *Organization      `json:"organization,omitempty"`
	Role         ConnectionPortRole `json:"role,omitempty"`
	// Speed is the maximum allowed throughput. This value inherits changes made in the Equinix Fabric API.
	Speed           uint64               `json:"speed,omitempty"`
	Status          string               `json:"status,omitempty"`
	Tokens          []FabricServiceToken `json:"tokens,omitempty"`
	VirtualCircuits []VirtualCircuit     `json:"virtual_circuits,omitempty"`
}

type Connection struct {
	ID               string               `json:"id"`
	ContactEmail     string               `json:"contact_email,omitempty"`
	Name             string               `json:"name,omitempty"`
	Status           string               `json:"status,omitempty"`
	Redundancy       ConnectionRedundancy `json:"redundancy,omitempty"`
	Facility         *Facility            `json:"facility,omitempty"`
	Metro            *Metro               `json:"metro,omitempty"`
	Type             ConnectionType       `json:"type,omitempty"`
	Mode             *ConnectionMode      `json:"mode,omitempty"`
	Description      string               `json:"description,omitempty"`
	Project          *Project             `json:"project,omitempty"`
	Organization     *Organization        `json:"organization,omitempty"`
	Speed            uint64               `json:"speed,omitempty"`
	Token            string               `json:"token,omitempty"`
	Tokens           []FabricServiceToken `json:"service_tokens,omitempty"`
	Tags             []string             `json:"tags,omitempty"`
	Ports            []ConnectionPort     `json:"ports,omitempty"`
	ServiceTokenType string               `json:"service_token_type,omitempty"`
}

type ConnectionCreateRequest struct {
	ContactEmail     string                 `json:"contact_email,omitempty"`
	Description      *string                `json:"description,omitempty"`
	Facility         string                 `json:"facility,omitempty"`
	Metro            string                 `json:"metro,omitempty"`
	Mode             ConnectionMode         `json:"mode,omitempty"`
	Name             string                 `json:"name,omitempty"`
	Project          string                 `json:"project,omitempty"`
	Redundancy       ConnectionRedundancy   `json:"redundancy,omitempty"`
	ServiceTokenType FabricServiceTokenType `json:"service_token_type,omitempty"`
	Speed            uint64                 `json:"speed,omitempty"`
	Tags             []string               `json:"tags,omitempty"`
	Type             ConnectionType         `json:"type,omitempty"`
	VLANs            []int                  `json:"vlans,omitempty"`
}

type ConnectionUpdateRequest struct {
	Redundancy  ConnectionRedundancy `json:"redundancy,omitempty"`
	Mode        *ConnectionMode      `json:"mode,omitempty"`
	Description *string              `json:"description,omitempty"`
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
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	apiUrl := path.Join(organizationBasePath, id, connectionBasePath)
	return s.create(apiUrl, createRequest)
}

func (s *ConnectionServiceOp) ProjectCreate(id string, createRequest *ConnectionCreateRequest) (*Connection, *Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
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
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	apiUrl := path.Join(organizationBasePath, id, connectionBasePath)
	return s.list(apiUrl, opts)
}

func (s *ConnectionServiceOp) ProjectList(id string, opts *GetOptions) ([]Connection, *Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	apiUrl := path.Join(projectBasePath, id, connectionBasePath)
	return s.list(apiUrl, opts)
}

func (s *ConnectionServiceOp) Delete(id string, wait bool) (*Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(connectionBasePath, id)
	connection := new(Connection)

	resp, err := s.client.DoRequest("DELETE", apiPath, nil, connection)
	if err != nil {
		return resp, err
	}
	// We really miss context here.
	if wait {
		timeout := time.After(ConnectionDeleteTimeout)
		//ticker := time.Tick(ConnectionDeleteCheck)
		ticker := time.NewTicker(ConnectionDeleteCheck)
		for {
			select {
			case <-ticker.C:
				c, resp2, err := s.Get(id, nil)
				if resp2.StatusCode == 404 {
					// Connection has been deleted
					return resp, nil
				}
				if err != nil {
					return resp, err
				}
				if c.Status != "deleting" {
					return resp, fmt.Errorf("Connection %s is in undexpected state %s", id, c.Status)
				}
			case <-timeout:
				return resp, errors.New("Timeout waiting for connection to be deleted")
			}
		}
	}
	return resp, nil
}

func (s *ConnectionServiceOp) Port(connID, portID string, opts *GetOptions) (*ConnectionPort, *Response, error) {
	if validateErr := ValidateUUID(connID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
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
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(connectionBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	connection := new(Connection)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, connection)
	if err != nil {
		return nil, resp, err
	}
	return connection, resp, err
}

func (s *ConnectionServiceOp) Update(id string, updateRequest *ConnectionUpdateRequest, opts *GetOptions) (*Connection, *Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(connectionBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	connection := new(Connection)
	resp, err := s.client.DoRequest("PUT", apiPathQuery, updateRequest, connection)
	if err != nil {
		return nil, resp, err
	}

	return connection, resp, err
}

func (s *ConnectionServiceOp) Ports(connID string, opts *GetOptions) ([]ConnectionPort, *Response, error) {
	if validateErr := ValidateUUID(connID); validateErr != nil {
		return nil, nil, validateErr
	}
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
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(connectionBasePath, id, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *ConnectionServiceOp) PortEvents(connID, portID string, opts *GetOptions) ([]Event, *Response, error) {
	if validateErr := ValidateUUID(connID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(connectionBasePath, connID, portBasePath, portID, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *ConnectionServiceOp) VirtualCircuits(connID, portID string, opts *GetOptions) (vcs []VirtualCircuit, resp *Response, err error) {
	if validateErr := ValidateUUID(connID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
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
