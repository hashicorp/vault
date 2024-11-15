package packngo

import (
	"math"
	"path"
)

const spotMarketRequestBasePath = "/spot-market-requests"

type SpotMarketRequestService interface {
	List(string, *ListOptions) ([]SpotMarketRequest, *Response, error)
	Create(*SpotMarketRequestCreateRequest, string) (*SpotMarketRequest, *Response, error)
	Delete(string, bool) (*Response, error)
	Get(string, *GetOptions) (*SpotMarketRequest, *Response, error)
}

type SpotMarketRequestCreateRequest struct {
	DevicesMax  int        `json:"devices_max"`
	DevicesMin  int        `json:"devices_min"`
	EndAt       *Timestamp `json:"end_at,omitempty"`
	FacilityIDs []string   `json:"facilities,omitempty"`
	Metro       string     `json:"metro,omitempty"`
	MaxBidPrice float64    `json:"max_bid_price"`

	Parameters SpotMarketRequestInstanceParameters `json:"instance_parameters"`
}

type SpotMarketRequest struct {
	SpotMarketRequestCreateRequest
	ID         string     `json:"id"`
	Devices    []Device   `json:"devices"`
	Facilities []Facility `json:"facilities,omitempty"`
	Metro      *Metro     `json:"metro,omitempty"`
	Project    Project    `json:"project"`
	Href       string     `json:"href"`
	Plan       Plan       `json:"plan"`
}

type SpotMarketRequestInstanceParameters struct {
	AlwaysPXE       bool       `json:"always_pxe,omitempty"`
	IPXEScriptURL   string     `json:"ipxe_script_url,omitempty"`
	BillingCycle    string     `json:"billing_cycle"`
	CustomData      string     `json:"customdata,omitempty"`
	Description     string     `json:"description,omitempty"`
	Features        []string   `json:"features,omitempty"`
	Hostname        string     `json:"hostname,omitempty"`
	Hostnames       []string   `json:"hostnames,omitempty"`
	Locked          bool       `json:"locked,omitempty"`
	OperatingSystem string     `json:"operating_system"`
	Plan            string     `json:"plan"`
	ProjectSSHKeys  []string   `json:"project_ssh_keys,omitempty"`
	Tags            []string   `json:"tags"`
	TerminationTime *Timestamp `json:"termination_time,omitempty"`
	UserSSHKeys     []string   `json:"user_ssh_keys,omitempty"`
	UserData        string     `json:"userdata"`
}

type SpotMarketRequestServiceOp struct {
	client *Client
}

func roundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

func (s *SpotMarketRequestServiceOp) Create(cr *SpotMarketRequestCreateRequest, pID string) (*SpotMarketRequest, *Response, error) {
	if validateErr := ValidateUUID(pID); validateErr != nil {
		return nil, nil, validateErr
	}
	opts := (&GetOptions{}).Including("devices", "project", "plan")
	endpointPath := path.Join(projectBasePath, pID, spotMarketRequestBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	cr.MaxBidPrice = roundPlus(cr.MaxBidPrice, 2)
	smr := new(SpotMarketRequest)

	resp, err := s.client.DoRequest("POST", apiPathQuery, cr, smr)
	if err != nil {
		return nil, resp, err
	}

	return smr, resp, err
}

func (s *SpotMarketRequestServiceOp) List(pID string, opts *ListOptions) ([]SpotMarketRequest, *Response, error) {
	if validateErr := ValidateUUID(pID); validateErr != nil {
		return nil, nil, validateErr
	}
	type smrRoot struct {
		SMRs []SpotMarketRequest `json:"spot_market_requests"`
	}

	endpointPath := path.Join(projectBasePath, pID, spotMarketRequestBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	output := new(smrRoot)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, output)
	if err != nil {
		return nil, nil, err
	}

	return output.SMRs, resp, nil
}

func (s *SpotMarketRequestServiceOp) Get(id string, opts *GetOptions) (*SpotMarketRequest, *Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(spotMarketRequestBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	smr := new(SpotMarketRequest)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, &smr)
	if err != nil {
		return nil, resp, err
	}

	return smr, resp, err
}

func (s *SpotMarketRequestServiceOp) Delete(id string, forceDelete bool) (*Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(spotMarketRequestBasePath, id)
	var params *map[string]bool
	if forceDelete {
		params = &map[string]bool{"force_termination": true}
	}
	return s.client.DoRequest("DELETE", apiPath, params, nil)
}
