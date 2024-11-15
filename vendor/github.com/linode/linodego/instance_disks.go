package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// InstanceDisk represents an Instance Disk object
type InstanceDisk struct {
	ID         int            `json:"id"`
	Label      string         `json:"label"`
	Status     DiskStatus     `json:"status"`
	Size       int            `json:"size"`
	Filesystem DiskFilesystem `json:"filesystem"`
	Created    *time.Time     `json:"-"`
	Updated    *time.Time     `json:"-"`

	// NOTE: Disk encryption may not currently be available to all users.
	DiskEncryption InstanceDiskEncryption `json:"disk_encryption"`
}

// DiskFilesystem constants start with Filesystem and include Linode API Filesystems
type DiskFilesystem string

// DiskFilesystem constants represent the filesystems types an Instance Disk may use
const (
	FilesystemRaw    DiskFilesystem = "raw"
	FilesystemSwap   DiskFilesystem = "swap"
	FilesystemExt3   DiskFilesystem = "ext3"
	FilesystemExt4   DiskFilesystem = "ext4"
	FilesystemInitrd DiskFilesystem = "initrd"
)

// DiskStatus constants have the prefix "Disk" and include Linode API Instance Disk Status
type DiskStatus string

// DiskStatus constants represent the status values an Instance Disk may have
const (
	DiskReady    DiskStatus = "ready"
	DiskNotReady DiskStatus = "not ready"
	DiskDeleting DiskStatus = "deleting"
)

// InstanceDiskCreateOptions are InstanceDisk settings that can be used at creation
type InstanceDiskCreateOptions struct {
	Label string `json:"label"`
	Size  int    `json:"size"`

	// Image is optional, but requires RootPass if provided
	Image    string `json:"image,omitempty"`
	RootPass string `json:"root_pass,omitempty"`

	Filesystem      string            `json:"filesystem,omitempty"`
	AuthorizedKeys  []string          `json:"authorized_keys,omitempty"`
	AuthorizedUsers []string          `json:"authorized_users,omitempty"`
	StackscriptID   int               `json:"stackscript_id,omitempty"`
	StackscriptData map[string]string `json:"stackscript_data,omitempty"`
}

// InstanceDiskUpdateOptions are InstanceDisk settings that can be used in updates
type InstanceDiskUpdateOptions struct {
	Label string `json:"label"`
}

// ListInstanceDisks lists InstanceDisks
func (c *Client) ListInstanceDisks(ctx context.Context, linodeID int, opts *ListOptions) ([]InstanceDisk, error) {
	response, err := getPaginatedResults[InstanceDisk](ctx, c, formatAPIPath("linode/instances/%d/disks", linodeID), opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *InstanceDisk) UnmarshalJSON(b []byte) error {
	type Mask InstanceDisk

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Updated *parseabletime.ParseableTime `json:"updated"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)

	return nil
}

// GetInstanceDisk gets the template with the provided ID
func (c *Client) GetInstanceDisk(ctx context.Context, linodeID int, diskID int) (*InstanceDisk, error) {
	e := formatAPIPath("linode/instances/%d/disks/%d", linodeID, diskID)
	response, err := doGETRequest[InstanceDisk](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateInstanceDisk creates a new InstanceDisk for the given Instance
func (c *Client) CreateInstanceDisk(ctx context.Context, linodeID int, opts InstanceDiskCreateOptions) (*InstanceDisk, error) {
	e := formatAPIPath("linode/instances/%d/disks", linodeID)
	response, err := doPOSTRequest[InstanceDisk](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateInstanceDisk creates a new InstanceDisk for the given Instance
func (c *Client) UpdateInstanceDisk(ctx context.Context, linodeID int, diskID int, opts InstanceDiskUpdateOptions) (*InstanceDisk, error) {
	e := formatAPIPath("linode/instances/%d/disks/%d", linodeID, diskID)
	response, err := doPUTRequest[InstanceDisk](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// RenameInstanceDisk renames an InstanceDisk
func (c *Client) RenameInstanceDisk(ctx context.Context, linodeID int, diskID int, label string) (*InstanceDisk, error) {
	return c.UpdateInstanceDisk(ctx, linodeID, diskID, InstanceDiskUpdateOptions{Label: label})
}

// ResizeInstanceDisk resizes the size of the Instance disk
func (c *Client) ResizeInstanceDisk(ctx context.Context, linodeID int, diskID int, size int) error {
	opts := map[string]any{
		"size": size,
	}

	e := formatAPIPath("linode/instances/%d/disks/%d/resize", linodeID, diskID)
	_, err := doPOSTRequest[InstanceDisk](ctx, c, e, opts)

	return err
}

// PasswordResetInstanceDisk resets the "root" account password on the Instance disk
func (c *Client) PasswordResetInstanceDisk(ctx context.Context, linodeID int, diskID int, password string) error {
	opts := map[string]any{
		"password": password,
	}

	e := formatAPIPath("linode/instances/%d/disks/%d/password", linodeID, diskID)
	_, err := doPOSTRequest[InstanceDisk](ctx, c, e, opts)

	return err
}

// DeleteInstanceDisk deletes a Linode Instance Disk
func (c *Client) DeleteInstanceDisk(ctx context.Context, linodeID int, diskID int) error {
	e := formatAPIPath("linode/instances/%d/disks/%d", linodeID, diskID)
	err := doDELETERequest(ctx, c, e)
	return err
}
