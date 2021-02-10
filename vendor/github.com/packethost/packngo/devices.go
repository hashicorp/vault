package packngo

import (
	"fmt"
	"path"
)

const deviceBasePath = "/devices"

const (
	NetworkTypeHybrid       = "hybrid"
	NetworkTypeL2Bonded     = "layer2-bonded"
	NetworkTypeL2Individual = "layer2-individual"
	NetworkTypeL3           = "layer3"
)

// DeviceService interface defines available device methods
type DeviceService interface {
	List(ProjectID string, opts *ListOptions) ([]Device, *Response, error)
	Get(DeviceID string, opts *GetOptions) (*Device, *Response, error)
	Create(*DeviceCreateRequest) (*Device, *Response, error)
	Update(string, *DeviceUpdateRequest) (*Device, *Response, error)
	Delete(string, bool) (*Response, error)
	Reboot(string) (*Response, error)
	PowerOff(string) (*Response, error)
	PowerOn(string) (*Response, error)
	Lock(string) (*Response, error)
	Unlock(string) (*Response, error)
	ListBGPSessions(deviceID string, opts *ListOptions) ([]BGPSession, *Response, error)
	ListBGPNeighbors(deviceID string, opts *ListOptions) ([]BGPNeighbor, *Response, error)
	ListEvents(deviceID string, opts *ListOptions) ([]Event, *Response, error)
}

type devicesRoot struct {
	Devices []Device `json:"devices"`
	Meta    meta     `json:"meta"`
}

// Device represents an Equinix Metal device from API
type Device struct {
	ID                  string                 `json:"id"`
	Href                string                 `json:"href,omitempty"`
	Hostname            string                 `json:"hostname,omitempty"`
	Description         *string                `json:"description,omitempty"`
	State               string                 `json:"state,omitempty"`
	Created             string                 `json:"created_at,omitempty"`
	Updated             string                 `json:"updated_at,omitempty"`
	Locked              bool                   `json:"locked,omitempty"`
	BillingCycle        string                 `json:"billing_cycle,omitempty"`
	Storage             *CPR                   `json:"storage,omitempty"`
	Tags                []string               `json:"tags,omitempty"`
	Network             []*IPAddressAssignment `json:"ip_addresses"`
	Volumes             []*Volume              `json:"volumes"`
	OS                  *OS                    `json:"operating_system,omitempty"`
	Plan                *Plan                  `json:"plan,omitempty"`
	Facility            *Facility              `json:"facility,omitempty"`
	Project             *Project               `json:"project,omitempty"`
	ProvisionEvents     []*Event               `json:"provisioning_events,omitempty"`
	ProvisionPer        float32                `json:"provisioning_percentage,omitempty"`
	UserData            string                 `json:"userdata,omitempty"`
	User                string                 `json:"user,omitempty"`
	RootPassword        string                 `json:"root_password,omitempty"`
	IPXEScriptURL       string                 `json:"ipxe_script_url,omitempty"`
	AlwaysPXE           bool                   `json:"always_pxe,omitempty"`
	HardwareReservation Href                   `json:"hardware_reservation,omitempty"`
	SpotInstance        bool                   `json:"spot_instance,omitempty"`
	SpotPriceMax        float64                `json:"spot_price_max,omitempty"`
	TerminationTime     *Timestamp             `json:"termination_time,omitempty"`
	NetworkPorts        []Port                 `json:"network_ports,omitempty"`
	CustomData          map[string]interface{} `json:"customdata,omitempty"`
	SSHKeys             []SSHKey               `json:"ssh_keys,omitempty"`
	ShortID             string                 `json:"short_id,omitempty"`
	SwitchUUID          string                 `json:"switch_uuid,omitempty"`
}

type NetworkInfo struct {
	PublicIPv4  string
	PublicIPv6  string
	PrivateIPv4 string
}

func (d *Device) GetNetworkInfo() NetworkInfo {
	ni := NetworkInfo{}
	for _, ip := range d.Network {
		// Initial device IPs are fixed and marked as "Management"
		if ip.Management {
			if ip.AddressFamily == 4 {
				if ip.Public {
					ni.PublicIPv4 = ip.Address
				} else {
					ni.PrivateIPv4 = ip.Address
				}
			} else {
				ni.PublicIPv6 = ip.Address
			}
		}
	}
	return ni
}

func (d Device) String() string {
	return Stringify(d)
}

func (d *Device) NumOfBonds() int {
	numOfBonds := 0
	for _, p := range d.NetworkPorts {
		if p.Type == "NetworkBondPort" {
			numOfBonds++
		}
	}
	return numOfBonds
}

func (d *Device) GetPortsInBond(name string) map[string]*Port {
	ports := map[string]*Port{}
	for _, port := range d.NetworkPorts {
		if port.Bond != nil && port.Bond.Name == name {
			p := port
			ports[p.Name] = &p
		}
	}
	return ports
}

func (d *Device) GetBondPorts() map[string]*Port {
	ports := map[string]*Port{}
	for _, port := range d.NetworkPorts {
		if port.Type == "NetworkBondPort" {
			p := port
			ports[p.Name] = &p
		}
	}
	return ports
}

func (d *Device) GetPhysicalPorts() map[string]*Port {
	ports := map[string]*Port{}
	for _, port := range d.NetworkPorts {
		if port.Type == "NetworkPort" {
			p := port
			ports[p.Name] = &p
		}
	}
	return ports
}

func (d *Device) GetPortByName(name string) (*Port, error) {
	for _, port := range d.NetworkPorts {
		if port.Name == name {
			return &port, nil
		}
	}
	return nil, fmt.Errorf("Port %s not found in device %s", name, d.ID)
}

type ports map[string]*Port

func (ports ports) allBonded() bool {
	if ports == nil {
		return false
	}

	if len(ports) == 0 {
		return false
	}

	for _, p := range ports {
		if (p == nil) || (!p.Data.Bonded) {
			return false
		}
	}
	return true
}

func (d *Device) HasManagementIPs() bool {
	for _, ip := range d.Network {
		if ip.Management {
			return true
		}
	}
	return false
}

// GetNetworkType returns a composite network type identification for a device
// based on the plan, network_type, and IP management state of the device.
// GetNetworkType provides the same composite state rendered in the Packet
// Portal for a given device.
func (d *Device) GetNetworkType() string {
	if d.Plan != nil {
		if d.Plan.Slug == "baremetal_0" || d.Plan.Slug == "baremetal_1" {
			return NetworkTypeL3
		}
		if d.Plan.Slug == "baremetal_1e" {
			return NetworkTypeHybrid
		}
	}

	bonds := ports(d.GetBondPorts())
	phys := ports(d.GetPhysicalPorts())

	if bonds.allBonded() {
		if phys.allBonded() {
			if !d.HasManagementIPs() {
				return NetworkTypeL2Bonded
			}
			return NetworkTypeL3
		}
		return NetworkTypeHybrid
	}
	return NetworkTypeL2Individual
}

type IPAddressCreateRequest struct {
	// Address Family for IP Address
	AddressFamily int `json:"address_family"`

	// Address Type for IP Address
	Public bool `json:"public"`

	// CIDR Size for the IP Block created. Valid values depends on the operating system provisioned.
	CIDR int `json:"cidr,omitempty"`

	// Reservations are UUIDs of any IP reservations to use when assigning IPs
	Reservations []string `json:"ip_reservations,omitempty"`
}

// CPR is a struct for custom partitioning and RAID
// If you don't want to bother writing the struct, just write the CPR conf to
// a string and then do
//
// 	var cpr CPR
//  err := json.Unmarshal([]byte(cprString), &cpr)
//	if err != nil {
//		log.Fatal(err)
//	}
type CPR struct {
	Disks []struct {
		Device     string `json:"device"`
		WipeTable  bool   `json:"wipeTable"`
		Partitions []struct {
			Label  string `json:"label"`
			Number int    `json:"number"`
			Size   string `json:"size"`
		} `json:"partitions"`
	} `json:"disks"`
	Raid []struct {
		Devices []string `json:"devices"`
		Level   string   `json:"level"`
		Name    string   `json:"name"`
	} `json:"raid,omitempty"`
	Filesystems []struct {
		Mount struct {
			Device string `json:"device"`
			Format string `json:"format"`
			Point  string `json:"point"`
			Create struct {
				Options []string `json:"options"`
			} `json:"create"`
		} `json:"mount"`
	} `json:"filesystems"`
}

// DeviceCreateRequest type used to create an Equinix Metal device
type DeviceCreateRequest struct {
	Hostname              string     `json:"hostname"`
	Plan                  string     `json:"plan"`
	Facility              []string   `json:"facility"`
	OS                    string     `json:"operating_system"`
	BillingCycle          string     `json:"billing_cycle"`
	ProjectID             string     `json:"project_id"`
	UserData              string     `json:"userdata"`
	Storage               *CPR       `json:"storage,omitempty"`
	Tags                  []string   `json:"tags"`
	Description           string     `json:"description,omitempty"`
	IPXEScriptURL         string     `json:"ipxe_script_url,omitempty"`
	PublicIPv4SubnetSize  int        `json:"public_ipv4_subnet_size,omitempty"`
	AlwaysPXE             bool       `json:"always_pxe,omitempty"`
	HardwareReservationID string     `json:"hardware_reservation_id,omitempty"`
	SpotInstance          bool       `json:"spot_instance,omitempty"`
	SpotPriceMax          float64    `json:"spot_price_max,omitempty,string"`
	TerminationTime       *Timestamp `json:"termination_time,omitempty"`
	CustomData            string     `json:"customdata,omitempty"`
	// UserSSHKeys is a list of user UUIDs - essentialy a list of
	// collaborators. The users must be a collaborator in the same project
	// where the device is created. The user's SSH keys then go to the
	// device
	UserSSHKeys []string `json:"user_ssh_keys,omitempty"`
	// Project SSHKeys is a list of SSHKeys resource UUIDs. If this param
	// is supplied, only the listed SSHKeys will go to the device.
	// Any other Project SSHKeys and any User SSHKeys will not be present
	// in the device.
	ProjectSSHKeys []string                 `json:"project_ssh_keys,omitempty"`
	Features       map[string]string        `json:"features,omitempty"`
	IPAddresses    []IPAddressCreateRequest `json:"ip_addresses,omitempty"`
}

// DeviceUpdateRequest type used to update an Equinix Metal device
type DeviceUpdateRequest struct {
	Hostname      *string   `json:"hostname,omitempty"`
	Description   *string   `json:"description,omitempty"`
	UserData      *string   `json:"userdata,omitempty"`
	Locked        *bool     `json:"locked,omitempty"`
	Tags          *[]string `json:"tags,omitempty"`
	AlwaysPXE     *bool     `json:"always_pxe,omitempty"`
	IPXEScriptURL *string   `json:"ipxe_script_url,omitempty"`
	CustomData    *string   `json:"customdata,omitempty"`
}

func (d DeviceCreateRequest) String() string {
	return Stringify(d)
}

// DeviceActionRequest type used to execute actions on devices
type DeviceActionRequest struct {
	Type string `json:"type"`
}

type DeviceDeleteRequest struct {
	Force bool `json:"force_delete"`
}

func (d DeviceActionRequest) String() string {
	return Stringify(d)
}

// DeviceServiceOp implements DeviceService
type DeviceServiceOp struct {
	client *Client
}

// List returns devices on a project
//
// Regarding ListOptions.Search: The API documentation does not provide guidance
// on the fields that will be searched using this parameter, so this behavior is
// undefined and prone to change.
//
// As of 2020-10-20, ListOptions.Search will look for matches in the following
// Device properties: Hostname, Description, Tags, ID, ShortID, Network.Address,
// Plan.Name, Plan.Slug, Facility.Code, Facility.Name, OS.Name, OS.Slug,
// HardwareReservation.ID, HardwareReservation.ShortID
func (s *DeviceServiceOp) List(projectID string, opts *ListOptions) (devices []Device, resp *Response, err error) {
	opts = opts.Including("facility")
	endpointPath := path.Join(projectBasePath, projectID, deviceBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(devicesRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		devices = append(devices, subset.Devices...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns a device by id
func (s *DeviceServiceOp) Get(deviceID string, opts *GetOptions) (*Device, *Response, error) {
	opts = opts.Including("facility")
	endpointPath := path.Join(deviceBasePath, deviceID)
	apiPathQuery := opts.WithQuery(endpointPath)
	device := new(Device)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, device)
	if err != nil {
		return nil, resp, err
	}
	return device, resp, err
}

// Create creates a new device
func (s *DeviceServiceOp) Create(createRequest *DeviceCreateRequest) (*Device, *Response, error) {
	apiPath := path.Join(projectBasePath, createRequest.ProjectID, deviceBasePath)
	device := new(Device)

	resp, err := s.client.DoRequest("POST", apiPath, createRequest, device)
	if err != nil {
		return nil, resp, err
	}
	return device, resp, err
}

// Update updates an existing device
func (s *DeviceServiceOp) Update(deviceID string, updateRequest *DeviceUpdateRequest) (*Device, *Response, error) {
	opts := &GetOptions{}
	opts = opts.Including("facility")
	endpointPath := path.Join(deviceBasePath, deviceID)
	apiPathQuery := opts.WithQuery(endpointPath)
	device := new(Device)

	resp, err := s.client.DoRequest("PUT", apiPathQuery, updateRequest, device)
	if err != nil {
		return nil, resp, err
	}

	return device, resp, err
}

// Delete deletes a device
func (s *DeviceServiceOp) Delete(deviceID string, force bool) (*Response, error) {
	apiPath := path.Join(deviceBasePath, deviceID)
	req := &DeviceDeleteRequest{Force: force}

	return s.client.DoRequest("DELETE", apiPath, req, nil)
}

// Reboot reboots on a device
func (s *DeviceServiceOp) Reboot(deviceID string) (*Response, error) {
	apiPath := path.Join(deviceBasePath, deviceID, "actions")
	action := &DeviceActionRequest{Type: "reboot"}

	return s.client.DoRequest("POST", apiPath, action, nil)
}

// PowerOff powers on a device
func (s *DeviceServiceOp) PowerOff(deviceID string) (*Response, error) {
	apiPath := path.Join(deviceBasePath, deviceID, "actions")
	action := &DeviceActionRequest{Type: "power_off"}

	return s.client.DoRequest("POST", apiPath, action, nil)
}

// PowerOn powers on a device
func (s *DeviceServiceOp) PowerOn(deviceID string) (*Response, error) {
	apiPath := path.Join(deviceBasePath, deviceID, "actions")
	action := &DeviceActionRequest{Type: "power_on"}

	return s.client.DoRequest("POST", apiPath, action, nil)
}

type lockType struct {
	Locked bool `json:"locked"`
}

// Lock sets a device to "locked"
func (s *DeviceServiceOp) Lock(deviceID string) (*Response, error) {
	apiPath := path.Join(deviceBasePath, deviceID)
	action := lockType{Locked: true}

	return s.client.DoRequest("PATCH", apiPath, action, nil)
}

// Unlock sets a device to "unlocked"
func (s *DeviceServiceOp) Unlock(deviceID string) (*Response, error) {
	apiPath := path.Join(deviceBasePath, deviceID)
	action := lockType{Locked: false}

	return s.client.DoRequest("PATCH", apiPath, action, nil)
}

func (s *DeviceServiceOp) ListBGPNeighbors(deviceID string, opts *ListOptions) ([]BGPNeighbor, *Response, error) {
	root := new(bgpNeighborsRoot)
	endpointPath := path.Join(deviceBasePath, deviceID, bgpNeighborsBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.BGPNeighbors, resp, err
}

// ListBGPSessions returns all BGP Sessions associated with the device
func (s *DeviceServiceOp) ListBGPSessions(deviceID string, opts *ListOptions) (bgpSessions []BGPSession, resp *Response, err error) {

	endpointPath := path.Join(deviceBasePath, deviceID, bgpSessionBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(bgpSessionsRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		bgpSessions = append(bgpSessions, subset.Sessions...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// ListEvents returns list of device events
func (s *DeviceServiceOp) ListEvents(deviceID string, opts *ListOptions) ([]Event, *Response, error) {
	apiPath := path.Join(deviceBasePath, deviceID, eventBasePath)

	return listEvents(s.client, apiPath, opts)
}
