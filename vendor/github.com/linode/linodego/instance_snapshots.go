package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// InstanceBackupsResponse response struct for backup snapshot
type InstanceBackupsResponse struct {
	Automatic []*InstanceSnapshot             `json:"automatic"`
	Snapshot  *InstanceBackupSnapshotResponse `json:"snapshot"`
}

// InstanceBackupSnapshotResponse fields are those representing Instance Backup Snapshots
type InstanceBackupSnapshotResponse struct {
	Current    *InstanceSnapshot `json:"current"`
	InProgress *InstanceSnapshot `json:"in_progress"`
}

// RestoreInstanceOptions fields are those accepted by InstanceRestore
type RestoreInstanceOptions struct {
	LinodeID  int  `json:"linode_id"`
	Overwrite bool `json:"overwrite"`
}

// InstanceSnapshot represents a linode backup snapshot
type InstanceSnapshot struct {
	ID        int                     `json:"id"`
	Label     string                  `json:"label"`
	Status    InstanceSnapshotStatus  `json:"status"`
	Type      string                  `json:"type"`
	Created   *time.Time              `json:"-"`
	Updated   *time.Time              `json:"-"`
	Finished  *time.Time              `json:"-"`
	Configs   []string                `json:"configs"`
	Disks     []*InstanceSnapshotDisk `json:"disks"`
	Available bool                    `json:"available"`
}

// InstanceSnapshotDisk fields represent the source disk of a Snapshot
type InstanceSnapshotDisk struct {
	Label      string `json:"label"`
	Size       int    `json:"size"`
	Filesystem string `json:"filesystem"`
}

// InstanceSnapshotStatus constants start with Snapshot and include Linode API Instance Backup Snapshot status values
type InstanceSnapshotStatus string

// InstanceSnapshotStatus constants reflect the current status of an Instance Snapshot
var (
	SnapshotPaused              InstanceSnapshotStatus = "paused"
	SnapshotPending             InstanceSnapshotStatus = "pending"
	SnapshotRunning             InstanceSnapshotStatus = "running"
	SnapshotNeedsPostProcessing InstanceSnapshotStatus = "needsPostProcessing"
	SnapshotSuccessful          InstanceSnapshotStatus = "successful"
	SnapshotFailed              InstanceSnapshotStatus = "failed"
	SnapshotUserAborted         InstanceSnapshotStatus = "userAborted"
)

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *InstanceSnapshot) UnmarshalJSON(b []byte) error {
	type Mask InstanceSnapshot

	p := struct {
		*Mask
		Created  *parseabletime.ParseableTime `json:"created"`
		Updated  *parseabletime.ParseableTime `json:"updated"`
		Finished *parseabletime.ParseableTime `json:"finished"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Updated = (*time.Time)(p.Updated)
	i.Finished = (*time.Time)(p.Finished)

	return nil
}

// GetInstanceSnapshot gets the snapshot with the provided ID
func (c *Client) GetInstanceSnapshot(ctx context.Context, linodeID int, snapshotID int) (*InstanceSnapshot, error) {
	e := formatAPIPath("linode/instances/%d/backups/%d", linodeID, snapshotID)
	response, err := doGETRequest[InstanceSnapshot](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateInstanceSnapshot Creates or Replaces the snapshot Backup of a Linode. If a previous snapshot exists for this Linode, it will be deleted.
func (c *Client) CreateInstanceSnapshot(ctx context.Context, linodeID int, label string) (*InstanceSnapshot, error) {
	opts := map[string]string{"label": label}

	e := formatAPIPath("linode/instances/%d/backups", linodeID)
	response, err := doPOSTRequest[InstanceSnapshot](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetInstanceBackups gets the Instance's available Backups.
// This is not called ListInstanceBackups because a single object is returned, matching the API response.
func (c *Client) GetInstanceBackups(ctx context.Context, linodeID int) (*InstanceBackupsResponse, error) {
	e := formatAPIPath("linode/instances/%d/backups", linodeID)
	response, err := doGETRequest[InstanceBackupsResponse](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// EnableInstanceBackups Enables backups for the specified Linode.
func (c *Client) EnableInstanceBackups(ctx context.Context, linodeID int) error {
	e := formatAPIPath("linode/instances/%d/backups/enable", linodeID)
	_, err := doPOSTRequest[InstanceBackup, any](ctx, c, e)
	return err
}

// CancelInstanceBackups Cancels backups for the specified Linode.
func (c *Client) CancelInstanceBackups(ctx context.Context, linodeID int) error {
	e := formatAPIPath("linode/instances/%d/backups/cancel", linodeID)
	_, err := doPOSTRequest[InstanceBackup, any](ctx, c, e)
	return err
}

// RestoreInstanceBackup Restores a Linode's Backup to the specified Linode.
func (c *Client) RestoreInstanceBackup(ctx context.Context, linodeID int, backupID int, opts RestoreInstanceOptions) error {
	e := formatAPIPath("linode/instances/%d/backups/%d/restore", linodeID, backupID)
	_, err := doPOSTRequest[InstanceBackup](ctx, c, e, opts)
	return err
}
