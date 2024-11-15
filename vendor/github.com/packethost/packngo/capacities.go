package packngo

const (
	capacityBasePath       = "/capacity"
	capacityBasePathMetros = "/capacity/metros"
)

// CapacityService interface defines available capacity methods
type CapacityService interface {
	List() (*CapacityReport, *Response, error)
	ListMetros() (*CapacityReport, *Response, error)
	Check(*CapacityInput) (*CapacityInput, *Response, error)
	CheckMetros(*CapacityInput) (*CapacityInput, *Response, error)
}

// CapacityInput struct
type CapacityInput struct {
	Servers []ServerInfo `json:"servers,omitempty"`
}

// ServerInfo struct
type ServerInfo struct {
	Facility  string `json:"facility,omitempty"`
	Metro     string `json:"metro,omitempty"`
	Plan      string `json:"plan,omitempty"`
	Quantity  int    `json:"quantity,omitempty"`
	Available bool   `json:"available,omitempty"`
}

type capacityRoot struct {
	Capacity CapacityReport `json:"capacity,omitempty"`
}

// CapacityReport map
type CapacityReport map[string]map[string]CapacityPerBaremetal

// // CapacityPerFacility struct
// type CapacityPerFacility struct {
// 	T1SmallX86  *CapacityPerBaremetal `json:"t1.small.x86,omitempty"`
// 	C1SmallX86  *CapacityPerBaremetal `json:"c1.small.x86,omitempty"`
// 	M1XlargeX86 *CapacityPerBaremetal `json:"m1.xlarge.x86,omitempty"`
// 	C1XlargeX86 *CapacityPerBaremetal `json:"c1.xlarge.x86,omitempty"`

// 	Baremetal0   *CapacityPerBaremetal `json:"baremetal_0,omitempty"`
// 	Baremetal1   *CapacityPerBaremetal `json:"baremetal_1,omitempty"`
// 	Baremetal1e  *CapacityPerBaremetal `json:"baremetal_1e,omitempty"`
// 	Baremetal2   *CapacityPerBaremetal `json:"baremetal_2,omitempty"`
// 	Baremetal2a  *CapacityPerBaremetal `json:"baremetal_2a,omitempty"`
// 	Baremetal2a2 *CapacityPerBaremetal `json:"baremetal_2a2,omitempty"`
// 	Baremetal3   *CapacityPerBaremetal `json:"baremetal_3,omitempty"`
// }

// CapacityPerBaremetal struct
type CapacityPerBaremetal struct {
	Level string `json:"level,omitempty"`
}

// CapacityList struct
type CapacityList struct {
	Capacity CapacityReport `json:"capacity,omitempty"`
}

// CapacityServiceOp implements CapacityService
type CapacityServiceOp struct {
	client *Client
}

func capacityList(client *Client, capUrl string) (*CapacityReport, *Response, error) {
	root := new(capacityRoot)

	resp, err := client.DoRequest("GET", capUrl, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.Capacity, nil, nil
}

// List returns a list of facilities and plans with their current capacity.
func (s *CapacityServiceOp) List() (*CapacityReport, *Response, error) {
	return capacityList(s.client, capacityBasePath)
}

// ListMetros returns a list of metros and plans with their current capacity.
func (s *CapacityServiceOp) ListMetros() (*CapacityReport, *Response, error) {
	return capacityList(s.client, capacityBasePathMetros)
}

func checkCapacity(client *Client, input *CapacityInput, capUrl string) (capInput *CapacityInput, resp *Response, err error) {
	capInput = new(CapacityInput)
	resp, err = client.DoRequest("POST", capUrl, input, capInput)
	return capInput, resp, err
}

// Check validates if a deploy can be fulfilled in a capacity.
func (s *CapacityServiceOp) Check(input *CapacityInput) (capInput *CapacityInput, resp *Response, err error) {
	return checkCapacity(s.client, input, capacityBasePath)
}

// Check validates if a deploy can be fulfilled in a metro.
func (s *CapacityServiceOp) CheckMetros(input *CapacityInput) (capInput *CapacityInput, resp *Response, err error) {
	return checkCapacity(s.client, input, capacityBasePathMetros)
}
