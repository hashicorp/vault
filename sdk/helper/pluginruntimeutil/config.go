package pluginruntimeutil

import "github.com/hashicorp/vault/sdk/helper/consts"

const (
	DefaultOCIRuntime = "runsc"
	DefaultCPU        = 0.1
	DefaultMemory     = 100000000
)

// PluginRuntimeConfig defines the metadata needed to run a plugin runtime
type PluginRuntimeConfig struct {
	Name         string                   `json:"name" structs:"name"`
	Type         consts.PluginRuntimeType `json:"type" structs:"type"`
	OCIRuntime   string                   `json:"oci_runtime" structs:"oci_runtime"`
	CgroupParent string                   `json:"cgroup_parent" structs:"cgroup_parent"`
	CPU          float32                  `json:"cpu" structs:"cpu"`
	Memory       uint64                   `json:"memory" structs:"memory"`
}
