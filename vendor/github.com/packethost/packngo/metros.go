package packngo

const metroBasePath = "/locations/metros"

// MetroService interface defines available metro methods
type MetroService interface {
	List(*ListOptions) ([]Metro, *Response, error)
}

type metroRoot struct {
	Metros []Metro `json:"metros"`
}

// Metro represents an Equinix Metal metro
type Metro struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Code    string `json:"code,omitempty"`
	Country string `json:"country,omitempty"`
}

func (f Metro) String() string {
	return Stringify(f)
}

// MetroServiceOp implements MetroService
type MetroServiceOp struct {
	client *Client
}

// List returns all metros
func (s *MetroServiceOp) List(opts *ListOptions) ([]Metro, *Response, error) {
	root := new(metroRoot)
	apiPathQuery := opts.WithQuery(metroBasePath)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Metros, resp, err
}
