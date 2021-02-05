package api

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// CSIVolumes is used to query the top level csi volumes
type CSIVolumes struct {
	client *Client
}

// CSIVolumes returns a handle on the CSIVolumes endpoint
func (c *Client) CSIVolumes() *CSIVolumes {
	return &CSIVolumes{client: c}
}

// List returns all CSI volumes
func (v *CSIVolumes) List(q *QueryOptions) ([]*CSIVolumeListStub, *QueryMeta, error) {
	var resp []*CSIVolumeListStub
	qm, err := v.client.query("/v1/volumes?type=csi", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	sort.Sort(CSIVolumeIndexSort(resp))
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

func (v *CSIVolumes) Register(vol *CSIVolume, w *WriteOptions) (*WriteMeta, error) {
	req := CSIVolumeRegisterRequest{
		Volumes: []*CSIVolume{vol},
	}
	meta, err := v.client.write("/v1/volume/csi/"+vol.ID, req, nil, w)
	return meta, err
}

func (v *CSIVolumes) Deregister(id string, force bool, w *WriteOptions) error {
	_, err := v.client.delete(fmt.Sprintf("/v1/volume/csi/%v?force=%t", url.PathEscape(id), force), nil, w)
	return err
}

func (v *CSIVolumes) Detach(volID, nodeID string, w *WriteOptions) error {
	_, err := v.client.delete(fmt.Sprintf("/v1/volume/csi/%v/detach?node=%v", url.PathEscape(volID), nodeID), nil, w)
	return err
}

// CSIVolumeAttachmentMode duplicated in nomad/structs/csi.go
type CSIVolumeAttachmentMode string

const (
	CSIVolumeAttachmentModeUnknown     CSIVolumeAttachmentMode = ""
	CSIVolumeAttachmentModeBlockDevice CSIVolumeAttachmentMode = "block-device"
	CSIVolumeAttachmentModeFilesystem  CSIVolumeAttachmentMode = "file-system"
)

// CSIVolumeAccessMode duplicated in nomad/structs/csi.go
type CSIVolumeAccessMode string

const (
	CSIVolumeAccessModeUnknown CSIVolumeAccessMode = ""

	CSIVolumeAccessModeSingleNodeReader CSIVolumeAccessMode = "single-node-reader-only"
	CSIVolumeAccessModeSingleNodeWriter CSIVolumeAccessMode = "single-node-writer"

	CSIVolumeAccessModeMultiNodeReader       CSIVolumeAccessMode = "multi-node-reader-only"
	CSIVolumeAccessModeMultiNodeSingleWriter CSIVolumeAccessMode = "multi-node-single-writer"
	CSIVolumeAccessModeMultiNodeMultiWriter  CSIVolumeAccessMode = "multi-node-multi-writer"
)

type CSIMountOptions struct {
	FSType       string   `hcl:"fs_type,optional"`
	MountFlags   []string `hcl:"mount_flags,optional"`
	ExtraKeysHCL []string `hcl1:",unusedKeys" json:"-"` // report unexpected keys
}

type CSISecrets map[string]string

// CSIVolume is used for serialization, see also nomad/structs/csi.go
type CSIVolume struct {
	ID             string
	Name           string
	ExternalID     string `hcl:"external_id"`
	Namespace      string
	Topologies     []*CSITopology
	AccessMode     CSIVolumeAccessMode     `hcl:"access_mode"`
	AttachmentMode CSIVolumeAttachmentMode `hcl:"attachment_mode"`
	MountOptions   *CSIMountOptions        `hcl:"mount_options"`
	Secrets        CSISecrets              `hcl:"secrets"`
	Parameters     map[string]string       `hcl:"parameters"`
	Context        map[string]string       `hcl:"context"`

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
	PluginID            string `hcl:"plugin_id"`
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

	// ExtraKeysHCL is used by the hcl parser to report unexpected keys
	ExtraKeysHCL []string `hcl1:",unusedKeys" json:"-"`
}

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
}

type CSIVolumeRegisterRequest struct {
	Volumes []*CSIVolume
	WriteRequest
}

type CSIVolumeDeregisterRequest struct {
	VolumeIDs []string
	WriteRequest
}

// CSI Plugins are jobs with plugin specific data
type CSIPlugins struct {
	client *Client
}

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
