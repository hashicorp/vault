package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Server represents a  server
type Server struct {
	// Arch is the architecture target of the server
	Arch string `json:"arch,omitempty"`

	// Identifier is a unique identifier for the server
	Identifier string `json:"id,omitempty"`

	// Name is the user-defined name of the server
	Name string `json:"name,omitempty"`

	// CreationDate is the creation date of the server
	CreationDate string `json:"creation_date,omitempty"`

	// ModificationDate is the date of the last modification of the server
	ModificationDate string `json:"modification_date,omitempty"`

	// Image is the image used by the server
	Image Image `json:"image,omitempty"`

	// DynamicIPRequired is a flag that defines a server with a dynamic ip address attached
	DynamicIPRequired *bool `json:"dynamic_ip_required,omitempty"`

	// PublicIP is the public IP address bound to the server
	PublicAddress IPAddress `json:"public_ip,omitempty"`

	// State is the current status of the server
	State string `json:"state,omitempty"`

	// StateDetail is the detailed status of the server
	StateDetail string `json:"state_detail,omitempty"`

	// PrivateIP represents the private IPV4 attached to the server (changes on each boot)
	PrivateIP string `json:"private_ip,omitempty"`

	// Bootscript is the unique identifier of the selected bootscript
	Bootscript *Bootscript `json:"bootscript,omitempty"`

	// BootType defines the type of boot. Can be local or bootscript
	BootType string `json:"boot_type,omitempty"`

	// Hostname represents the ServerName in a format compatible with unix's hostname
	Hostname string `json:"hostname,omitempty"`

	// Tags represents user-defined tags
	Tags []string `json:"tags,omitempty"`

	// Volumes are the attached volumes
	Volumes map[string]Volume `json:"volumes,omitempty"`

	// SecurityGroup is the selected security group object
	SecurityGroup SecurityGroupRef `json:"security_group,omitempty"`

	// Organization is the owner of the server
	Organization string `json:"organization,omitempty"`

	// CommercialType is the commercial type of the server (i.e: C1, C2[SML], VC1S)
	CommercialType string `json:"commercial_type,omitempty"`

	// Location of the server
	Location struct {
		Platform   string `json:"platform_id,omitempty"`
		Chassis    string `json:"chassis_id,omitempty"`
		Cluster    string `json:"cluster_id,omitempty"`
		Hypervisor string `json:"hypervisor_id,omitempty"`
		Blade      string `json:"blade_id,omitempty"`
		Node       string `json:"node_id,omitempty"`
		ZoneID     string `json:"zone_id,omitempty"`
	} `json:"location,omitempty"`

	IPV6 *IPV6 `json:"ipv6,omitempty"`

	EnableIPV6 bool `json:"enable_ipv6,omitempty"`

	// This fields are not returned by the API, we generate it
	DNSPublic  string `json:"dns_public,omitempty"`
	DNSPrivate string `json:"dns_private,omitempty"`
}

// ServerPatchDefinition represents a  server with nullable fields (for PATCH)
type ServerPatchDefinition struct {
	Arch              *string            `json:"arch,omitempty"`
	Name              *string            `json:"name,omitempty"`
	CreationDate      *string            `json:"creation_date,omitempty"`
	ModificationDate  *string            `json:"modification_date,omitempty"`
	Image             *Image             `json:"image,omitempty"`
	DynamicIPRequired *bool              `json:"dynamic_ip_required,omitempty"`
	PublicAddress     *IPAddress         `json:"public_ip,omitempty"`
	State             *string            `json:"state,omitempty"`
	StateDetail       *string            `json:"state_detail,omitempty"`
	PrivateIP         *string            `json:"private_ip,omitempty"`
	Bootscript        *string            `json:"bootscript,omitempty"`
	Hostname          *string            `json:"hostname,omitempty"`
	Volumes           *map[string]Volume `json:"volumes,omitempty"`
	SecurityGroup     *SecurityGroupRef  `json:"security_group,omitempty"`
	Organization      *string            `json:"organization,omitempty"`
	Tags              *[]string          `json:"tags,omitempty"`
	IPV6              *IPV6              `json:"ipv6,omitempty"`
	EnableIPV6        *bool              `json:"enable_ipv6,omitempty"`
}

// ServerDefinition represents a  server with image definition
type ServerDefinition struct {
	// Name is the user-defined name of the server
	Name string `json:"name"`

	// Image is the image used by the server
	Image *string `json:"image,omitempty"`

	// Volumes are the attached volumes
	Volumes map[string]string `json:"volumes,omitempty"`

	// DynamicIPRequired is a flag that defines a server with a dynamic ip address attached
	DynamicIPRequired *bool `json:"dynamic_ip_required,omitempty"`

	// Bootscript is the bootscript used by the server
	Bootscript *string `json:"bootscript"`

	// Tags are the metadata tags attached to the server
	Tags []string `json:"tags,omitempty"`

	// Organization is the owner of the server
	Organization string `json:"organization"`

	// CommercialType is the commercial type of the server (i.e: C1, C2[SML], VC1S)
	CommercialType string `json:"commercial_type"`

	// BootType defines the type of boot. Can be local or bootscript
	BootType string `json:"boot_type,omitempty"`

	PublicIP string `json:"public_ip,omitempty"`

	EnableIPV6 bool `json:"enable_ipv6,omitempty"`

	SecurityGroup string `json:"security_group,omitempty"`
}

// Servers represents a group of  servers
type Servers struct {
	// Servers holds  servers of the response
	Servers []Server `json:"servers,omitempty"`
}

// ServerAction represents an action to perform on a  server
type ServerAction struct {
	// Action is the name of the action to trigger
	Action string `json:"action,omitempty"`
}

// OneServer represents the response of a GET /servers/UUID API call
type OneServer struct {
	Server Server `json:"server,omitempty"`
}

// PatchServer updates a server
func (s *API) PatchServer(serverID string, definition ServerPatchDefinition) error {
	resp, err := s.PatchResponse(s.computeAPI, fmt.Sprintf("servers/%s", serverID), definition)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := s.handleHTTPError([]int{http.StatusOK}, resp); err != nil {
		return err
	}
	return nil
}

// GetServers gets the list of servers from the API
func (s *API) GetServers(all bool, limit int) ([]Server, error) {
	query := url.Values{}
	if !all {
		query.Set("state", "running")
	}
	// TODO per_page=20&page=2&state=running
	if limit > 0 {
		// FIXME: wait for the API to be ready
		// query.Set("per_page", strconv.Itoa(limit))
		panic("Not implemented yet")
	}

	servers, err := s.fetchServers(query)
	if err != nil {
		return nil, err
	}

	for i, server := range servers.Servers {
		servers.Servers[i].DNSPublic = server.Identifier + URLPublicDNS
		servers.Servers[i].DNSPrivate = server.Identifier + URLPrivateDNS
	}
	return servers.Servers, nil
}

// SortServers represents a wrapper to sort by CreationDate the servers
type SortServers []Server

func (s SortServers) Len() int {
	return len(s)
}

func (s SortServers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortServers) Less(i, j int) bool {
	date1, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", s[i].CreationDate)
	date2, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", s[j].CreationDate)
	return date2.Before(date1)
}

// GetServer gets a server from the API
func (s *API) GetServer(serverID string) (*Server, error) {
	if serverID == "" {
		return nil, fmt.Errorf("cannot get server without serverID")
	}
	resp, err := s.GetResponsePaginate(s.computeAPI, "servers/"+serverID, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}

	var oneServer OneServer

	if err = json.Unmarshal(body, &oneServer); err != nil {
		return nil, err
	}
	// FIXME arch, owner, title
	oneServer.Server.DNSPublic = oneServer.Server.Identifier + URLPublicDNS
	oneServer.Server.DNSPrivate = oneServer.Server.Identifier + URLPrivateDNS
	return &oneServer.Server, nil
}

// PostServerAction posts an action on a server
func (s *API) PostServerAction(serverID, action string) (*Task, error) {
	data := ServerAction{
		Action: action,
	}
	resp, err := s.PostResponse(s.computeAPI, fmt.Sprintf("servers/%s/action", serverID), data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusAccepted}, resp)
	if err != nil {
		return nil, err
	}

	var t oneTask
	if err = json.Unmarshal(body, &t); err != nil {
		return nil, err
	}
	return &t.Task, err
}

func (s *API) fetchServers(query url.Values) (*Servers, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "servers", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var servers Servers

	if err = json.Unmarshal(body, &servers); err != nil {
		return nil, err
	}
	return &servers, nil
}

// DeleteServer deletes a server
func (s *API) DeleteServer(serverID string) error {
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("servers/%s", serverID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = s.handleHTTPError([]int{http.StatusNoContent}, resp); err != nil {
		return err
	}
	return nil
}

// CreateServer creates a new server
func (s *API) CreateServer(definition ServerDefinition) (*Server, error) {
	definition.Organization = s.Organization

	resp, err := s.PostResponse(s.computeAPI, "servers", definition)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return nil, err
	}
	var data OneServer

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data.Server, nil
}
