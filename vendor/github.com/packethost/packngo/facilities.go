package packngo

const facilityBasePath = "/facilities"

// FacilityService interface defines available facility methods
type FacilityService interface {
	List(*ListOptions) ([]Facility, *Response, error)
}

type facilityRoot struct {
	Facilities []Facility `json:"facilities"`
}

// Facility represents an Equinix Metal facility
type Facility struct {
	ID       string   `json:"id"`
	Name     string   `json:"name,omitempty"`
	Code     string   `json:"code,omitempty"`
	Features []string `json:"features,omitempty"`
	Address  *Address `json:"address,omitempty"`
	Metro    *Metro   `json:"metro,omitempty"`
	URL      string   `json:"href,omitempty"`
}

func (f Facility) String() string {
	return Stringify(f)
}

// Coordinates struct for Coordinates
type Coordinates struct {
	Latitude  *string `json:"latitude,omitempty"`
	Longitude *string `json:"longitude,omitempty"`
}

// Address - the physical address of the facility
type Address struct {
	ID          string       `json:"id,omitempty"`
	Address     string       `json:"address"`
	Address2    *string      `json:"address2,omitempty"`
	City        *string      `json:"city,omitempty"`
	State       *string      `json:"state,omitempty"`
	ZipCode     string       `json:"zip_code"`
	Country     string       `json:"country"`
	Coordinates *Coordinates `json:"coordinates,omitempty"`
}

func (a Address) String() string {
	return Stringify(a)
}

// FacilityServiceOp implements FacilityService
type FacilityServiceOp struct {
	client *Client
}

// List returns all facilities
func (s *FacilityServiceOp) List(opts *ListOptions) ([]Facility, *Response, error) {
	root := new(facilityRoot)
	apiPathQuery := opts.WithQuery(facilityBasePath)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Facilities, resp, err
}
