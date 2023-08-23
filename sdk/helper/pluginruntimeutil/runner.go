package pluginruntimeutil

import "github.com/hashicorp/vault/sdk/helper/consts"

const (
	DefaultOCIRuntime = "runsc"
	DefaultCPU        = 0.5
	DefaultMemory     = 10000
)

// PluginRuntimeRunner defines the metadata needed to run a plugin runtime
type PluginRuntimeConfig struct {
	Name         string                   `json:"name" structs:"name"`
	Type         consts.PluginRuntimeType `json:"type" structs:"type"`
	OCIRuntime   string                   `json:"oci_runtime" structs:"oci_runtime"`
	ParentCGroup string                   `json:"parent_cgroup" structs:"parent_cgroup"`
	CPU          float32                  `json:"cpu" structs:"cpu"`
	Memory       uint64                   `json:"memory" structs:"memory"`
}
