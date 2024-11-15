package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// VolumeStatus indicates the status of the Volume
type VolumeStatus string

const (
	// VolumeCreating indicates the Volume is being created and is not yet available for use
	VolumeCreating VolumeStatus = "creating"

	// VolumeActive indicates the Volume is online and available for use
	VolumeActive VolumeStatus = "active"

	// VolumeResizing indicates the Volume is in the process of upgrading its current capacity
	VolumeResizing VolumeStatus = "resizing"

	// VolumeContactSupport indicates there is a problem with the Volume. A support ticket must be opened to resolve the issue
	VolumeContactSupport VolumeStatus = "contact_support"
)

// Volume represents a linode volume object
type Volume struct {
	ID             int          `json:"id"`
	Label          string       `json:"label"`
	Status         VolumeStatus `json:"status"`
	Region         string       `json:"region"`
	Size           int          `json:"size"`
	LinodeID       *int         `json:"linode_id"`
	FilesystemPath string       `json:"filesystem_path"`
	Tags           []string     `json:"tags"`
	Created        *time.Time   `json:"-"`
	Updated        *time.Time   `json:"-"`

	// Note: Block Storage Disk Encryption is not currently available to all users.
	Encryption string `json:"encryption"`
}

// VolumeCreateOptions fields are those accepted by CreateVolume
type VolumeCreateOptions struct {
	Label    string `json:"label,omitempty"`
	Region   string `json:"region,omitempty"`
	LinodeID int    `json:"linode_id,omitempty"`
	ConfigID int    `json:"config_id,omitempty"`
	// The Volume's size, in GiB. Minimum size is 10GiB, maximum size is 10240GiB. A "0" value will result in the default size.
	Size int `json:"size,omitempty"`
	// An array of tags applied to this object. Tags are for organizational purposes only.
	Tags               []string `json:"tags"`
	PersistAcrossBoots *bool    `json:"persist_across_boots,omitempty"`
	Encryption         string   `json:"encryption,omitempty"`
}

// VolumeUpdateOptions fields are those accepted by UpdateVolume
type VolumeUpdateOptions struct {
	Label string    `json:"label,omitempty"`
	Tags  *[]string `json:"tags,omitempty"`
}

// VolumeAttachOptions fields are those accepted by AttachVolume
type VolumeAttachOptions struct {
	LinodeID           int   `json:"linode_id"`
	ConfigID           int   `json:"config_id,omitempty"`
	PersistAcrossBoots *bool `json:"persist_across_boots,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Volume) UnmarshalJSON(b []byte) error {
	type Mask Volume

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(v),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	v.Created = (*time.Time)(p.Created)
	v.Updated = (*time.Time)(p.Updated)

	return nil
}

// GetUpdateOptions converts a Volume to VolumeUpdateOptions for use in UpdateVolume
func (v Volume) GetUpdateOptions() (updateOpts VolumeUpdateOptions) {
	updateOpts.Label = v.Label
	updateOpts.Tags = &v.Tags
	return
}

// GetCreateOptions converts a Volume to VolumeCreateOptions for use in CreateVolume
func (v Volume) GetCreateOptions() (createOpts VolumeCreateOptions) {
	createOpts.Label = v.Label
	createOpts.Tags = v.Tags
	createOpts.Region = v.Region
	createOpts.Size = v.Size
	if v.LinodeID != nil && *v.LinodeID > 0 {
		createOpts.LinodeID = *v.LinodeID
	}
	return
}

// ListVolumes lists Volumes
func (c *Client) ListVolumes(ctx context.Context, opts *ListOptions) ([]Volume, error) {
	response, err := getPaginatedResults[Volume](ctx, c, "volumes", opts)
	return response, err
}

// GetVolume gets the template with the provided ID
func (c *Client) GetVolume(ctx context.Context, volumeID int) (*Volume, error) {
	e := formatAPIPath("volumes/%d", volumeID)
	response, err := doGETRequest[Volume](ctx, c, e)
	return response, err
}

// AttachVolume attaches a volume to a Linode instance
func (c *Client) AttachVolume(ctx context.Context, volumeID int, opts *VolumeAttachOptions) (*Volume, error) {
	e := formatAPIPath("volumes/%d/attach", volumeID)
	response, err := doPOSTRequest[Volume](ctx, c, e, opts)
	return response, err
}

// CreateVolume creates a Linode Volume
func (c *Client) CreateVolume(ctx context.Context, opts VolumeCreateOptions) (*Volume, error) {
	e := "volumes"
	response, err := doPOSTRequest[Volume](ctx, c, e, opts)
	return response, err
}

// UpdateVolume updates the Volume with the specified id
func (c *Client) UpdateVolume(ctx context.Context, volumeID int, opts VolumeUpdateOptions) (*Volume, error) {
	e := formatAPIPath("volumes/%d", volumeID)
	response, err := doPUTRequest[Volume](ctx, c, e, opts)
	return response, err
}

// CloneVolume clones a Linode volume
func (c *Client) CloneVolume(ctx context.Context, volumeID int, label string) (*Volume, error) {
	opts := map[string]any{
		"label": label,
	}

	e := formatAPIPath("volumes/%d/clone", volumeID)
	response, err := doPOSTRequest[Volume](ctx, c, e, opts)
	return response, err
}

// DetachVolume detaches a Linode volume
func (c *Client) DetachVolume(ctx context.Context, volumeID int) error {
	e := formatAPIPath("volumes/%d/detach", volumeID)
	_, err := doPOSTRequest[Volume, any](ctx, c, e)
	return err
}

// ResizeVolume resizes an instance to new Linode type
func (c *Client) ResizeVolume(ctx context.Context, volumeID int, size int) error {
	opts := map[string]int{
		"size": size,
	}

	e := formatAPIPath("volumes/%d/resize", volumeID)
	_, err := doPOSTRequest[Volume](ctx, c, e, opts)
	return err
}

// DeleteVolume deletes the Volume with the specified id
func (c *Client) DeleteVolume(ctx context.Context, volumeID int) error {
	e := formatAPIPath("volumes/%d", volumeID)
	err := doDELETERequest(ctx, c, e)
	return err
}
