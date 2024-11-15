package packngo

import (
	"fmt"
	"path"
)

const batchBasePath = "/batches"

// BatchService interface defines available batch methods
type BatchService interface {
	Get(batchID string, getOpt *GetOptions) (*Batch, *Response, error)
	List(ProjectID string, listOpt *ListOptions) ([]Batch, *Response, error)
	Create(projectID string, batches *BatchCreateRequest) ([]Batch, *Response, error)
	Delete(string, bool) (*Response, error)
}

// Batch type
type Batch struct {
	ID            string   `json:"id"`
	ErrorMessages []string `json:"error_messages,omitempty"`

	// State may be 'failed' or 'completed'
	State     string     `json:"state,omitempty"`
	Quantity  int32      `json:"quantity,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
	Href      string     `json:"href,omitempty"`
	Project   Href       `json:"project,omitempty"`
	Devices   []Device   `json:"devices,omitempty"`
}

//BatchesList represents collection of batches
type batchesList struct {
	Batches []Batch `json:"batches,omitempty"`
}

// BatchCreateRequest type used to create batch of device instances
type BatchCreateRequest struct {
	Batches []BatchCreateDevice `json:"batches"`
}

// BatchCreateDevice type used to describe batch instances
type BatchCreateDevice struct {
	DeviceCreateRequest
	Quantity               int32   `json:"quantity"`
	FacilityDiversityLevel int32   `json:"facility_diversity_level,omitempty"`
	SpotInstance           bool    `json:"spot_instance,omitempty"`
	SpotPriceMax           float64 `json:"spot_price_max,omitempty"`
}

// BatchServiceOp implements BatchService
type BatchServiceOp struct {
	client *Client
}

// Get returns batch details
func (s *BatchServiceOp) Get(batchID string, opts *GetOptions) (*Batch, *Response, error) {
	if validateErr := ValidateUUID(batchID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(batchBasePath, batchID)
	apiPathQuery := opts.WithQuery(endpointPath)
	batch := new(Batch)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, batch)
	if err != nil {
		return nil, resp, err
	}

	return batch, resp, err
}

// List returns batches on a project
func (s *BatchServiceOp) List(projectID string, opts *ListOptions) (batches []Batch, resp *Response, err error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(projectBasePath, projectID, batchBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	subset := new(batchesList)
	resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
	if err != nil {
		return nil, resp, err
	}

	batches = append(batches, subset.Batches...)
	return batches, resp, err
}

// Create function to create batch of device instances
func (s *BatchServiceOp) Create(projectID string, request *BatchCreateRequest) ([]Batch, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(projectBasePath, projectID, "devices", "batch")

	batches := new(batchesList)
	resp, err := s.client.DoRequest("POST", apiPath, request, batches)

	if err != nil {
		return nil, resp, err
	}

	return batches.Batches, resp, err
}

// Delete function to remove an instance batch
func (s *BatchServiceOp) Delete(id string, removeDevices bool) (*Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, validateErr
	}
	// API doc days the remove_associated_instances params shout be in the body
	// https://metal.equinix.com/developers/api/batches/#delete-the-batch
	// .. does this even work?
	apiPath := fmt.Sprintf("%s/%s?remove_associated_instances=%t", batchBasePath, id, removeDevices)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
