package packngo

import (
	"path"
)

const (
	portVLANAssignmentsPath      = "vlan-assignments"
	portVLANAssignmentsBatchPath = "batches"
)

type vlanAssignmentsRoot struct {
	VLANAssignments []VLANAssignment `json:"vlan_assignments"`
	Meta            meta             `json:"meta"`
}

type vlanAssignmentBatchesRoot struct {
	VLANAssignmentBatches []VLANAssignmentBatch `json:"batches"`
	Meta                  meta                  `json:"meta"`
}

// VLANAssignmentService handles operations on a VLANAssignment
type VLANAssignmentService interface {
	Get(string, string, *GetOptions) (*VLANAssignment, *Response, error)
	List(string, *ListOptions) ([]VLANAssignment, *Response, error)

	GetBatch(string, string, *GetOptions) (*VLANAssignmentBatch, *Response, error)
	ListBatch(string, *ListOptions) ([]VLANAssignmentBatch, *Response, error)
	CreateBatch(string, *VLANAssignmentBatchCreateRequest, *GetOptions) (*VLANAssignmentBatch, *Response, error)
}

type VLANAssignmentServiceOp struct {
	client requestDoer
}

var _ VLANAssignmentService = (*VLANAssignmentServiceOp)(nil)

type VLANAssignmentBatchState string

const (
	VLANAssignmentBatchQueued     VLANAssignmentBatchState = "queued"
	VLANAssignmentBatchInProgress VLANAssignmentBatchState = "in_progress"
	VLANAssignmentBatchCompleted  VLANAssignmentBatchState = "completed"
	VLANAssignmentBatchFailed     VLANAssignmentBatchState = "failed"
)

type VLANAssignmentState string

const (
	VLANAssignmentAssigned   VLANAssignmentState = "assigned"
	VLANAssignmentUnassigned VLANAssignmentState = "unassigned"
)

// VLANAssignment struct for VLANAssignmentService.Get represents a port VLAN assignment that has been enacted
type VLANAssignment struct {
	// ID of the assignment
	ID string `json:"id,omitempty"`

	CreatedAt Timestamp `json:"created_at,omitempty"`
	UpdatedAt Timestamp `json:"updated_at,omitempty"`

	// Native indicates the VLAN is the native VLAN on the port and packets for this vlan will be untagged
	Native bool `json:"native,omitempty"`

	// State of the assignment
	State VLANAssignmentState `json:"state,omitempty"`

	// VLAN is the VirtualNetwork.VXLAN of the VLAN the assignment was made to
	VLAN int `json:"vlan,omitempty"`

	// Port is a reference to the Port the assignment was made on
	Port *Port `json:"port,omitempty"`

	// VirtualNetwork is a reference to the VLAN the assignment was made to
	VirtualNetwork *VirtualNetwork `json:"virtual_network,omitempty"`
}

// BatchedVLANAssignment represents the data requested in the batch before being processed. ID represents the VLAN ID, not the VLAN Assignment ID.
type BatchedVLANAssignment struct {
	// VirtualNetworkID is the VirtualNetwork.ID of the VLAN the assignment was made to
	VirtualNetworkID string `json:"id,omitempty"`

	// Native indicates the VLAN is the native VLAN on the port and packets for this vlan will be untagged
	Native bool `json:"native,omitempty"`

	// State of the assignment
	State VLANAssignmentState `json:"state,omitempty"`

	// VLAN is the VirtualNetwork.VXLAN of the VLAN the assignment was made to
	VLAN int `json:"vlan,omitempty"`
}

// VLANAssignmentBatch struct for VLANAssignmentBatch
type VLANAssignmentBatch struct {
	ID              string                   `json:"id,omitempty"`
	ErrorMessages   []string                 `json:"error_messages,omitempty"`
	Quantity        int                      `json:"quantity,omitempty"`
	State           VLANAssignmentBatchState `json:"state,omitempty"`
	CreatedAt       Timestamp                `json:"created_at,omitempty"`
	UpdatedAt       Timestamp                `json:"updated_at,omitempty"`
	Port            *Port                    `json:"port,omitempty"`
	Project         *Project                 `json:"project,omitempty"`
	VLANAssignments []BatchedVLANAssignment  `json:"vlan_assignments,omitempty"`
}

// VLANAssignmentBatchCreateRequest struct for VLANAssignmentBatch Create
type VLANAssignmentBatchCreateRequest struct {
	VLANAssignments []VLANAssignmentCreateRequest `json:"vlan_assignments"`
}

// VLANAssignmentCreateRequest struct for VLANAssignmentBatchCreateRequest
type VLANAssignmentCreateRequest struct {
	VLAN   string              `json:"vlan,omitempty"`
	State  VLANAssignmentState `json:"state,omitempty"`
	Native *bool               `json:"native,omitempty"`
}

// List returns VLANAssignmentBatches
func (s *VLANAssignmentServiceOp) ListBatch(portID string, opts *ListOptions) (results []VLANAssignmentBatch, resp *Response, err error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(portBasePath, portID, portVLANAssignmentsPath, portVLANAssignmentsBatchPath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(vlanAssignmentBatchesRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		results = append(results, subset.VLANAssignmentBatches...)
		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// Get returns a VLANAssignmentBatch by id
func (s *VLANAssignmentServiceOp) GetBatch(portID, batchID string, opts *GetOptions) (*VLANAssignmentBatch, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(batchID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(portBasePath, portID, portVLANAssignmentsPath, portVLANAssignmentsBatchPath, batchID)
	apiPathQuery := opts.WithQuery(endpointPath)
	batch := new(VLANAssignmentBatch)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, batch)
	if err != nil {
		return nil, resp, err
	}
	return batch, resp, err
}

// Create creates VLANAssignmentBatch objects
func (s *VLANAssignmentServiceOp) CreateBatch(portID string, request *VLANAssignmentBatchCreateRequest, opts *GetOptions) (*VLANAssignmentBatch, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(portBasePath, portID, portVLANAssignmentsPath, portVLANAssignmentsBatchPath)
	apiPathQuery := opts.WithQuery(endpointPath)
	batch := new(VLANAssignmentBatch)
	resp, err := s.client.DoRequest("POST", apiPathQuery, request, batch)
	if err != nil {
		return nil, resp, err
	}
	return batch, resp, err
}

// List returns VLANAssignment
func (s *VLANAssignmentServiceOp) List(portID string, opts *ListOptions) (results []VLANAssignment, resp *Response, err error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(portBasePath, portID, portVLANAssignmentsPath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(vlanAssignmentsRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		results = append(results, subset.VLANAssignments...)
		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// Get returns a VLANAssignment by id
func (s *VLANAssignmentServiceOp) Get(portID, assignmentID string, opts *GetOptions) (*VLANAssignment, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(assignmentID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(portBasePath, portID, portVLANAssignmentsPath, assignmentID)
	apiPathQuery := opts.WithQuery(endpointPath)
	VLANAssignment := new(VLANAssignment)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, VLANAssignment)
	if err != nil {
		return nil, resp, err
	}
	return VLANAssignment, resp, err
}
