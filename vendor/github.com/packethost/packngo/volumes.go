package packngo

import (
	"path"
)

const (
	volumeBasePath      = "/storage"
	attachmentsBasePath = "/attachments"
)

// VolumeService interface defines available Volume methods
type VolumeService interface {
	List(string, *ListOptions) ([]Volume, *Response, error)
	Get(string, *GetOptions) (*Volume, *Response, error)
	Update(string, *VolumeUpdateRequest) (*Volume, *Response, error)
	Delete(string) (*Response, error)
	Create(*VolumeCreateRequest, string) (*Volume, *Response, error)
	Lock(string) (*Response, error)
	Unlock(string) (*Response, error)
}

// VolumeAttachmentService defines attachment methdods
type VolumeAttachmentService interface {
	Get(string, *GetOptions) (*VolumeAttachment, *Response, error)
	Create(string, string) (*VolumeAttachment, *Response, error)
	Delete(string) (*Response, error)
}

type volumesRoot struct {
	Volumes []Volume `json:"volumes"`
	Meta    meta     `json:"meta"`
}

// Volume represents a volume
type Volume struct {
	Attachments      []*VolumeAttachment `json:"attachments,omitempty"`
	BillingCycle     string              `json:"billing_cycle,omitempty"`
	Created          string              `json:"created_at,omitempty"`
	Description      string              `json:"description,omitempty"`
	Facility         *Facility           `json:"facility,omitempty"`
	Href             string              `json:"href,omitempty"`
	ID               string              `json:"id"`
	Locked           bool                `json:"locked,omitempty"`
	Name             string              `json:"name,omitempty"`
	Plan             *Plan               `json:"plan,omitempty"`
	Project          *Project            `json:"project,omitempty"`
	Size             int                 `json:"size,omitempty"`
	SnapshotPolicies []*SnapshotPolicy   `json:"snapshot_policies,omitempty"`
	State            string              `json:"state,omitempty"`
	Updated          string              `json:"updated_at,omitempty"`
}

// SnapshotPolicy used to execute actions on volume
type SnapshotPolicy struct {
	ID                string `json:"id"`
	Href              string `json:"href"`
	SnapshotFrequency string `json:"snapshot_frequency,omitempty"`
	SnapshotCount     int    `json:"snapshot_count,omitempty"`
}

func (v Volume) String() string {
	return Stringify(v)
}

// VolumeCreateRequest type used to create an Equinix Metal volume
type VolumeCreateRequest struct {
	BillingCycle     string            `json:"billing_cycle"`
	Description      string            `json:"description,omitempty"`
	Locked           bool              `json:"locked,omitempty"`
	Size             int               `json:"size"`
	PlanID           string            `json:"plan_id"`
	FacilityID       string            `json:"facility_id"`
	SnapshotPolicies []*SnapshotPolicy `json:"snapshot_policies,omitempty"`
}

func (v VolumeCreateRequest) String() string {
	return Stringify(v)
}

// VolumeUpdateRequest type used to update an Equinix Metal volume
type VolumeUpdateRequest struct {
	Description  *string `json:"description,omitempty"`
	PlanID       *string `json:"plan_id,omitempty"`
	Size         *int    `json:"size,omitempty"`
	BillingCycle *string `json:"billing_cycle,omitempty"`
}

// VolumeAttachment is a type from Equinix Metal API
type VolumeAttachment struct {
	Href   string `json:"href"`
	ID     string `json:"id"`
	Volume Volume `json:"volume"`
	Device Device `json:"device"`
}

func (v VolumeUpdateRequest) String() string {
	return Stringify(v)
}

// VolumeAttachmentServiceOp implements VolumeService
type VolumeAttachmentServiceOp struct {
	client *Client
}

// VolumeServiceOp implements VolumeService
type VolumeServiceOp struct {
	client *Client
}

// List returns the volumes for a project
func (v *VolumeServiceOp) List(projectID string, opts *ListOptions) (volumes []Volume, resp *Response, err error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(projectBasePath, projectID, volumeBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	for {
		subset := new(volumesRoot)

		resp, err = v.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		volumes = append(volumes, subset.Volumes...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// Get returns a volume by id
func (v *VolumeServiceOp) Get(volumeID string, opts *GetOptions) (*Volume, *Response, error) {
	if validateErr := ValidateUUID(volumeID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(volumeBasePath, volumeID)
	apiPathQuery := opts.WithQuery(endpointPath)
	volume := new(Volume)

	resp, err := v.client.DoRequest("GET", apiPathQuery, nil, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Update updates a volume
func (v *VolumeServiceOp) Update(id string, updateRequest *VolumeUpdateRequest) (*Volume, *Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(volumeBasePath, id)
	volume := new(Volume)

	resp, err := v.client.DoRequest("PATCH", apiPath, updateRequest, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Delete deletes a volume
func (v *VolumeServiceOp) Delete(volumeID string) (*Response, error) {
	if validateErr := ValidateUUID(volumeID); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(volumeBasePath, volumeID)

	return v.client.DoRequest("DELETE", apiPath, nil, nil)
}

// Create creates a new volume for a project
func (v *VolumeServiceOp) Create(createRequest *VolumeCreateRequest, projectID string) (*Volume, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	url := path.Join(projectBasePath, projectID, volumeBasePath)
	volume := new(Volume)

	resp, err := v.client.DoRequest("POST", url, createRequest, volume)
	if err != nil {
		return nil, resp, err
	}

	return volume, resp, err
}

// Attachments

// Create Attachment, i.e. attach volume to a device
func (v *VolumeAttachmentServiceOp) Create(volumeID, deviceID string) (*VolumeAttachment, *Response, error) {
	if validateErr := ValidateUUID(volumeID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(deviceID); validateErr != nil {
		return nil, nil, validateErr
	}
	url := path.Join(volumeBasePath, volumeID, attachmentsBasePath)
	volAttachParam := map[string]string{
		"device_id": deviceID,
	}
	volumeAttachment := new(VolumeAttachment)

	resp, err := v.client.DoRequest("POST", url, volAttachParam, volumeAttachment)
	if err != nil {
		return nil, resp, err
	}
	return volumeAttachment, resp, nil
}

// Get gets attachment by id
func (v *VolumeAttachmentServiceOp) Get(attachmentID string, opts *GetOptions) (*VolumeAttachment, *Response, error) {
	if validateErr := ValidateUUID(attachmentID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(volumeBasePath, attachmentsBasePath, attachmentID)
	apiPathQuery := opts.WithQuery(endpointPath)
	volumeAttachment := new(VolumeAttachment)

	resp, err := v.client.DoRequest("GET", apiPathQuery, nil, volumeAttachment)
	if err != nil {
		return nil, resp, err
	}

	return volumeAttachment, resp, nil
}

// Delete deletes attachment by id
func (v *VolumeAttachmentServiceOp) Delete(attachmentID string) (*Response, error) {
	if validateErr := ValidateUUID(attachmentID); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(volumeBasePath, attachmentsBasePath, attachmentID)

	return v.client.DoRequest("DELETE", apiPath, nil, nil)
}

// Lock sets a volume to "locked"
func (v *VolumeServiceOp) Lock(id string) (*Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(volumeBasePath, id)
	action := lockType{Locked: true}

	return v.client.DoRequest("PATCH", apiPath, action, nil)
}

// Unlock sets a volume to "unlocked"
func (v *VolumeServiceOp) Unlock(id string) (*Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(volumeBasePath, id)
	action := lockType{Locked: false}

	return v.client.DoRequest("PATCH", apiPath, action, nil)
}
