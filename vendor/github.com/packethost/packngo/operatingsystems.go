package packngo

const osBasePath = "/operating-systems"

// OSService interface defines available operating_systems methods
type OSService interface {
	List() ([]OS, *Response, error)
}

type osRoot struct {
	OperatingSystems []OS `json:"operating_systems"`
}

// OS represents an Equinix Metal operating system
type OS struct {
	Name            string   `json:"name"`
	Slug            string   `json:"slug"`
	Distro          string   `json:"distro"`
	Version         string   `json:"version"`
	ProvisionableOn []string `json:"provisionable_on"`
}

func (o OS) String() string {
	return Stringify(o)
}

// OSServiceOp implements OSService
type OSServiceOp struct {
	client *Client
}

// List returns all available operating systems
func (s *OSServiceOp) List() ([]OS, *Response, error) {
	root := new(osRoot)

	resp, err := s.client.DoRequest("GET", osBasePath, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.OperatingSystems, resp, err
}
