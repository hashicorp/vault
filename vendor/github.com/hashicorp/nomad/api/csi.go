// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// CSIVolumes is used to access Container Storage Interface (CSI) endpoints.
type CSIVolumes struct {
	client *Client
}

// CSIVolumes returns a handle on the CSIVolumes endpoint.
func (c *Client) CSIVolumes() *CSIVolumes {
	return &CSIVolumes{client: c}
}

// List returns all CSI volumes.
func (v *CSIVolumes) List(q *QueryOptions) ([]*CSIVolumeListStub, *QueryMeta, error) {
	var resp []*CSIVolumeListStub
	qm, err := v.client.query("/v1/volumes?type=csi", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(CSIVolumeIndexSort(resp))
	return resp, qm, nil
}

// ListExternal returns all CSI volumes, as understood by the external storage
// provider. These volumes may or may not be currently registered with Nomad.
// The response is paginated by the plugin and accepts the
// QueryOptions.PerPage and QueryOptions.NextToken fields.
func (v *CSIVolumes) ListExternal(pluginID string, q *QueryOptions) (*CSIVolumeListExternalResponse, *QueryMeta, error) {
	var resp *CSIVolumeListExternalResponse

	qp := url.Values{}
	qp.Set("plugin_id", pluginID)
	if q.NextToken != "" {
		qp.Set("next_token", q.NextToken)
	}
	if q.PerPage != 0 {
		qp.Set("per_page", fmt.Sprint(q.PerPage))
	}

	qm, err := v.client.query("/v1/volumes/external?"+qp.Encode(), &resp, q)
	if err != nil {
		return nil, nil, err
	}

	sort.Sort(CSIVolumeExternalStubSort(resp.Volumes))
	return resp, qm, nil
}

// PluginList returns all CSI volumes for the specified plugin id
func (v *CSIVolumes) PluginList(pluginID string) ([]*CSIVolumeListStub, *QueryMeta, error) {
	return v.List(&QueryOptions{Prefix: pluginID})
}

// Info is used to retrieve a single CSIVolume
func (v *CSIVolumes) Info(id string, q *QueryOptions) (*CSIVolume, *QueryMeta, error) {
	var resp CSIVolume
	qm, err := v.client.query("/v1/volume/csi/"+id, &resp, q)
	if err != nil {
		return nil, nil, err
	}

	return &resp, qm, nil
}

// Register registers a single CSIVolume with Nomad. The volume must already
// exist in the external storage provider.
func (v *CSIVolumes) Register(vol *CSIVolume, w *WriteOptions) (*WriteMeta, error) {
	req := CSIVolumeRegisterRequest{
		Volumes: []*CSIVolume{vol},
	}
	meta, err := v.client.put("/v1/volume/csi/"+vol.ID, req, nil, w)
	return meta, err
}

// Deregister deregisters a single CSIVolume from Nomad. The volume will not be deleted from the external storage provider.
func (v *CSIVolumes) Deregister(id string, force bool, w *WriteOptions) error {
	_, err := v.client.delete(fmt.Sprintf("/v1/volume/csi/%v?force=%t", url.PathEscape(id), force), nil, nil, w)
	return err
}

// Create creates a single CSIVolume in an external storage provider and
// registers it with Nomad. You do not need to call Register if this call is
// successful.
func (v *CSIVolumes) Create(vol *CSIVolume, w *WriteOptions) ([]*CSIVolume, *WriteMeta, error) {
	req := CSIVolumeCreateRequest{
		Volumes: []*CSIVolume{vol},
	}

	resp := &CSIVolumeCreateResponse{}
	meta, err := v.client.put(fmt.Sprintf("/v1/volume/csi/%v/create", vol.ID), req, resp, w)
	return resp.Volumes, meta, err
}

// DEPRECATED: will be removed in Nomad 1.4.0
// Delete deletes a CSI volume from an external storage provider. The ID
// passed as an argument here is for the storage provider's ID, so a volume
// that's already been deregistered can be deleted.
func (v *CSIVolumes) Delete(externalVolID string, w *WriteOptions) error {
	_, err := v.client.delete(fmt.Sprintf("/v1/volume/csi/%v/delete", url.PathEscape(externalVolID)), nil, nil, w)
	return err
}

// DeleteOpts deletes a CSI volume from an external storage
// provider. The ID passed in the request is for the storage
// provider's ID, so a volume that's already been deregistered can be
// deleted.
func (v *CSIVolumes) DeleteOpts(req *CSIVolumeDeleteRequest, w *WriteOptions) error {
	if w == nil {
		w = &WriteOptions{}
	}
	w.SetHeadersFromCSISecrets(req.Secrets)
	_, err := v.client.delete(fmt.Sprintf("/v1/volume/csi/%v/delete", url.PathEscape(req.ExternalVolumeID)), nil, nil, w)
	return err
}

// Detach causes Nomad to attempt to detach a CSI volume from a client
// node. This is used in the case that the node is temporarily lost and the
// allocations are unable to drop their claims automatically.
func (v *CSIVolumes) Detach(volID, nodeID string, w *WriteOptions) error {
	_, err := v.client.delete(fmt.Sprintf("/v1/volume/csi/%v/detach?node=%v", url.PathEscape(volID), nodeID), nil, nil, w)
	return err
}

// CreateSnapshot snapshots an external storage volume.
func (v *CSIVolumes) CreateSnapshot(snap *CSISnapshot, w *WriteOptions) (*CSISnapshotCreateResponse, *WriteMeta, error) {
	req := &CSISnapshotCreateRequest{
		Snapshots: []*CSISnapshot{snap},
	}
	if w == nil {
		w = &WriteOptions{}
	}
	w.SetHeadersFromCSISecrets(snap.Secrets)
	resp := &CSISnapshotCreateResponse{}
	meta, err := v.client.put("/v1/volumes/snapshot", req, resp, w)
	return resp, meta, err
}

// DeleteSnapshot deletes an external storage volume snapshot.
func (v *CSIVolumes) DeleteSnapshot(snap *CSISnapshot, w *WriteOptions) error {
	qp := url.Values{}
	qp.Set("snapshot_id", snap.ID)
	qp.Set("plugin_id", snap.PluginID)
	if w == nil {
		w = &WriteOptions{}
	}
	w.SetHeadersFromCSISecrets(snap.Secrets)
	_, err := v.client.delete("/v1/volumes/snapshot?"+qp.Encode(), nil, nil, w)
	return err
}

// ListSnapshotsOpts lists external storage volume snapshots.
func (v *CSIVolumes) ListSnapshotsOpts(req *CSISnapshotListRequest) (*CSISnapshotListResponse, *QueryMeta, error) {
	var resp *CSISnapshotListResponse

	qp := url.Values{}
	if req.PluginID != "" {
		qp.Set("plugin_id", req.PluginID)
	}
	if req.NextToken != "" {
		qp.Set("next_token", req.NextToken)
	}
	if req.PerPage != 0 {
		qp.Set("per_page", fmt.Sprint(req.PerPage))
	}
	req.QueryOptions.SetHeadersFromCSISecrets(req.Secrets)

	qm, err := v.client.query("/v1/volumes/snapshot?"+qp.Encode(), &resp, &req.QueryOptions)
	if err != nil {
		return nil, nil, err
	}

	sort.Sort(CSISnapshotSort(resp.Snapshots))
	return resp, qm, nil
}

// DEPRECATED: will be removed in Nomad 1.4.0
// ListSnapshots lists external storage volume snapshots.
func (v *CSIVolumes) ListSnapshots(pluginID string, secrets string, q *QueryOptions) (*CSISnapshotListResponse, *QueryMeta, error) {
	var resp *CSISnapshotListResponse

	qp := url.Values{}
	if pluginID != "" {
		qp.Set("plugin_id", pluginID)
	}
	if q.NextToken != "" {
		qp.Set("next_token", q.NextToken)
	}
	if q.PerPage != 0 {
		qp.Set("per_page", fmt.Sprint(q.PerPage))
	}

	qm, err := v.client.query("/v1/volumes/snapshot?"+qp.Encode(), &resp, q)
	if err != nil {
		return nil, nil, err
	}

	sort.Sort(CSISnapshotSort(resp.Snapshots))
	return resp, qm, nil
}

// CSIVolumeAttachmentMode chooses the type of storage api that will be used to
// interact with the device. (Duplicated in nomad/structs/csi.go)
type CSIVolumeAttachmentMode string

const (
	CSIVolumeAttachmentModeUnknown     CSIVolumeAttachmentMode = ""
	CSIVolumeAttachmentModeBlockDevice CSIVolumeAttachmentMode = "block-device"
	CSIVolumeAttachmentModeFilesystem  CSIVolumeAttachmentMode = "file-system"
)

// CSIVolumeAccessMode indicates how a volume should be used in a storage topology
// e.g whether the provider should make the volume available concurrently. (Duplicated in nomad/structs/csi.go)
type CSIVolumeAccessMode string

const (
	CSIVolumeAccessModeUnknown               CSIVolumeAccessMode = ""
	CSIVolumeAccessModeSingleNodeReader      CSIVolumeAccessMode = "single-node-reader-only"
	CSIVolumeAccessModeSingleNodeWriter      CSIVolumeAccessMode = "single-node-writer"
	CSIVolumeAccessModeMultiNodeReader       CSIVolumeAccessMode = "multi-node-reader-only"
	CSIVolumeAccessModeMultiNodeSingleWriter CSIVolumeAccessMode = "multi-node-single-writer"
	CSIVolumeAccessModeMultiNodeMultiWriter  CSIVolumeAccessMode = "multi-node-multi-writer"
)

const (
	CSIVolumeTypeHost = "host"
	CSIVolumeTypeCSI  = "csi"
)

// CSIMountOptions contain optional additional configuration that can be used
// when specifying that a Volume should be used with VolumeAccessTypeMount.
type CSIMountOptions struct {
	// FSType is an optional field that allows an operator to specify the type
	// of the filesystem.
	FSType string `hcl:"fs_type,optional"`

	// MountFlags contains additional options that may be used when mounting the
	// volume by the plugin. This may contain sensitive data and should not be
	// leaked.
	MountFlags []string `hcl:"mount_flags,optional"`

	ExtraKeysHCL []string `hcl1:",unusedKeys" json:"-"` // report unexpected keys
}

func (o *CSIMountOptions) Merge(p *CSIMountOptions) {
	if p == nil {
		return
	}
	if p.FSType != "" {
		o.FSType = p.FSType
	}
	if p.MountFlags != nil {
		o.MountFlags = p.MountFlags
	}
}

// CSISecrets contain optional additional credentials that may be needed by
// the storage provider. These values will be redacted when reported in the
// API or in Nomad's logs.
type CSISecrets map[string]string

func (q *QueryOptions) SetHeadersFromCSISecrets(secrets CSISecrets) {
	pairs := []string{}
	for k, v := range secrets {
		pairs = append(pairs, fmt.Sprintf("%v=%v", k, v))
	}
	if q.Headers == nil {
		q.Headers = map[string]string{}
	}
	q.Headers["X-Nomad-CSI-Secrets"] = strings.Join(pairs, ",")
}

func (w *WriteOptions) SetHeadersFromCSISecrets(secrets CSISecrets) {
	pairs := []string{}
	for k, v := range secrets {
		pairs = append(pairs, fmt.Sprintf("%v=%v", k, v))
	}
	if w.Headers == nil {
		w.Headers = map[string]string{}
	}
	w.Headers["X-Nomad-CSI-Secrets"] = strings.Join(pairs, ",")
}

// CSIVolume is used for serialization, see also nomad/structs/csi.go
type CSIVolume struct {
	ID         string
	Name       string
	ExternalID string `mapstructure:"external_id" hcl:"external_id"`
	Namespace  string

	// RequestedTopologies are the topologies submitted as options to
	// the storage provider at the time the volume was created. After
	// volumes are created, this field is ignored.
	RequestedTopologies *CSITopologyRequest `hcl:"topology_request"`

	// Topologies are the topologies returned by the storage provider,
	// based on the RequestedTopologies and what the storage provider
	// could support. This value cannot be set by the user.
	Topologies []*CSITopology

	AccessMode     CSIVolumeAccessMode     `hcl:"access_mode"`
	AttachmentMode CSIVolumeAttachmentMode `hcl:"attachment_mode"`
	MountOptions   *CSIMountOptions        `hcl:"mount_options"`
	Secrets        CSISecrets              `mapstructure:"secrets" hcl:"secrets"`
	Parameters     map[string]string       `mapstructure:"parameters" hcl:"parameters"`
	Context        map[string]string       `mapstructure:"context" hcl:"context"`
	Capacity       int64                   `hcl:"-"`

	// These fields are used as part of the volume creation request
	RequestedCapacityMin  int64                  `hcl:"capacity_min"`
	RequestedCapacityMax  int64                  `hcl:"capacity_max"`
	RequestedCapabilities []*CSIVolumeCapability `hcl:"capability"`
	CloneID               string                 `mapstructure:"clone_id" hcl:"clone_id"`
	SnapshotID            string                 `mapstructure:"snapshot_id" hcl:"snapshot_id"`

	// ReadAllocs is a map of allocation IDs for tracking reader claim status.
	// The Allocation value will always be nil; clients can populate this data
	// by iterating over the Allocations field.
	ReadAllocs map[string]*Allocation

	// WriteAllocs is a map of allocation IDs for tracking writer claim
	// status. The Allocation value will always be nil; clients can populate
	// this data by iterating over the Allocations field.
	WriteAllocs map[string]*Allocation

	// Allocations is a combined list of readers and writers
	Allocations []*AllocationListStub

	// Schedulable is true if all the denormalized plugin health fields are true
	Schedulable         bool
	PluginID            string `mapstructure:"plugin_id" hcl:"plugin_id"`
	Provider            string
	ProviderVersion     string
	ControllerRequired  bool
	ControllersHealthy  int
	ControllersExpected int
	NodesHealthy        int
	NodesExpected       int
	ResourceExhausted   time.Time

	CreateIndex uint64
	ModifyIndex uint64

	// CreateTime stored as UnixNano
	CreateTime int64
	// ModifyTime stored as UnixNano
	ModifyTime int64

	// ExtraKeysHCL is used by the hcl parser to report unexpected keys
	ExtraKeysHCL []string `hcl1:",unusedKeys" json:"-"`
}

// CSIVolumeCapability is a requested attachment and access mode for a
// volume
type CSIVolumeCapability struct {
	AccessMode     CSIVolumeAccessMode     `mapstructure:"access_mode" hcl:"access_mode"`
	AttachmentMode CSIVolumeAttachmentMode `mapstructure:"attachment_mode" hcl:"attachment_mode"`
}

// CSIVolumeIndexSort is a helper used for sorting volume stubs by creation
// time.
type CSIVolumeIndexSort []*CSIVolumeListStub

func (v CSIVolumeIndexSort) Len() int {
	return len(v)
}

func (v CSIVolumeIndexSort) Less(i, j int) bool {
	return v[i].CreateIndex > v[j].CreateIndex
}

func (v CSIVolumeIndexSort) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// CSIVolumeListStub omits allocations. See also nomad/structs/csi.go
type CSIVolumeListStub struct {
	ID                  string
	Namespace           string
	Name                string
	ExternalID          string
	Topologies          []*CSITopology
	AccessMode          CSIVolumeAccessMode
	AttachmentMode      CSIVolumeAttachmentMode
	CurrentReaders      int
	CurrentWriters      int
	Schedulable         bool
	PluginID            string
	Provider            string
	ControllerRequired  bool
	ControllersHealthy  int
	ControllersExpected int
	NodesHealthy        int
	NodesExpected       int
	ResourceExhausted   time.Time

	CreateIndex uint64
	ModifyIndex uint64

	// CreateTime stored as UnixNano
	CreateTime int64
	// ModifyTime stored as UnixNano
	ModifyTime int64
}

type CSIVolumeListExternalResponse struct {
	Volumes   []*CSIVolumeExternalStub
	NextToken string
}

// CSIVolumeExternalStub is the storage provider's view of a volume, as
// returned from the controller plugin; all IDs are for external resources
type CSIVolumeExternalStub struct {
	ExternalID               string
	CapacityBytes            int64
	VolumeContext            map[string]string
	CloneID                  string
	SnapshotID               string
	PublishedExternalNodeIDs []string
	IsAbnormal               bool
	Status                   string
}

// CSIVolumeExternalStubSort is a sorting helper for external volumes. We
// can't sort these by creation time because we don't get that data back from
// the storage provider. Sort by External ID within this page.
type CSIVolumeExternalStubSort []*CSIVolumeExternalStub

func (v CSIVolumeExternalStubSort) Len() int {
	return len(v)
}

func (v CSIVolumeExternalStubSort) Less(i, j int) bool {
	return v[i].ExternalID > v[j].ExternalID
}

func (v CSIVolumeExternalStubSort) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

type CSIVolumeCreateRequest struct {
	Volumes []*CSIVolume
	WriteRequest
}

type CSIVolumeCreateResponse struct {
	Volumes []*CSIVolume
	QueryMeta
}

type CSIVolumeRegisterRequest struct {
	Volumes []*CSIVolume
	WriteRequest
}

type CSIVolumeDeregisterRequest struct {
	VolumeIDs []string
	WriteRequest
}

type CSIVolumeDeleteRequest struct {
	ExternalVolumeID string
	Secrets          CSISecrets
	WriteRequest
}

// CSISnapshot is the storage provider's view of a volume snapshot
type CSISnapshot struct {
	ID                     string // storage provider's ID
	ExternalSourceVolumeID string // storage provider's ID for volume
	SizeBytes              int64  // value from storage provider
	CreateTime             int64  // value from storage provider
	IsReady                bool   // value from storage provider
	SourceVolumeID         string // Nomad volume ID
	PluginID               string // CSI plugin ID

	// These field are only used during snapshot creation and will not be
	// populated when the snapshot is returned
	Name       string            // suggested name of the snapshot, used for creation
	Secrets    CSISecrets        // secrets needed to create snapshot
	Parameters map[string]string // secrets needed to create snapshot
}

// CSISnapshotSort is a helper used for sorting snapshots by creation time.
type CSISnapshotSort []*CSISnapshot

func (v CSISnapshotSort) Len() int {
	return len(v)
}

func (v CSISnapshotSort) Less(i, j int) bool {
	return v[i].CreateTime > v[j].CreateTime
}

func (v CSISnapshotSort) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

type CSISnapshotCreateRequest struct {
	Snapshots []*CSISnapshot
	WriteRequest
}

type CSISnapshotCreateResponse struct {
	Snapshots []*CSISnapshot
	QueryMeta
}

// CSISnapshotListRequest is a request to a controller plugin to list all the
// snapshot known to the storage provider. This request is paginated by
// the plugin and accepts the QueryOptions.PerPage and QueryOptions.NextToken
// fields
type CSISnapshotListRequest struct {
	PluginID string
	Secrets  CSISecrets
	QueryOptions
}

type CSISnapshotListResponse struct {
	Snapshots []*CSISnapshot
	NextToken string
	QueryMeta
}

// CSI Plugins are jobs with plugin specific data
type CSIPlugins struct {
	client *Client
}

// CSIPlugin is used for serialization, see also nomad/structs/csi.go
type CSIPlugin struct {
	ID                 string
	Provider           string
	Version            string
	ControllerRequired bool
	// Map Node.ID to CSIInfo fingerprint results
	Controllers         map[string]*CSIInfo
	Nodes               map[string]*CSIInfo
	Allocations         []*AllocationListStub
	ControllersHealthy  int
	ControllersExpected int
	NodesHealthy        int
	NodesExpected       int
	CreateIndex         uint64
	ModifyIndex         uint64

	// CreateTime stored as UnixNano
	CreateTime int64
	// ModifyTime stored as UnixNano
	ModifyTime int64
}

type CSIPluginListStub struct {
	ID                  string
	Provider            string
	ControllerRequired  bool
	ControllersHealthy  int
	ControllersExpected int
	NodesHealthy        int
	NodesExpected       int
	CreateIndex         uint64
	ModifyIndex         uint64
}

// CSIPluginIndexSort is a helper used for sorting plugin stubs by creation
// time.
type CSIPluginIndexSort []*CSIPluginListStub

func (v CSIPluginIndexSort) Len() int {
	return len(v)
}

func (v CSIPluginIndexSort) Less(i, j int) bool {
	return v[i].CreateIndex > v[j].CreateIndex
}

func (v CSIPluginIndexSort) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

// CSIPlugins returns a handle on the CSIPlugins endpoint
func (c *Client) CSIPlugins() *CSIPlugins {
	return &CSIPlugins{client: c}
}

// List returns all CSI plugins
func (v *CSIPlugins) List(q *QueryOptions) ([]*CSIPluginListStub, *QueryMeta, error) {
	var resp []*CSIPluginListStub
	qm, err := v.client.query("/v1/plugins?type=csi", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(CSIPluginIndexSort(resp))
	return resp, qm, nil
}

// Info is used to retrieve a single CSI Plugin Job
func (v *CSIPlugins) Info(id string, q *QueryOptions) (*CSIPlugin, *QueryMeta, error) {
	var resp *CSIPlugin
	qm, err := v.client.query("/v1/plugin/csi/"+id, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}
