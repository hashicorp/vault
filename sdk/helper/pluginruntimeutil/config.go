// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginruntimeutil

import "github.com/hashicorp/vault/sdk/helper/consts"

// PluginRuntimeConfig defines the metadata needed to run a plugin runtime
type PluginRuntimeConfig struct {
	Name         string                   `json:"name" structs:"name"`
	Type         consts.PluginRuntimeType `json:"type" structs:"type"`
	OCIRuntime   string                   `json:"oci_runtime" structs:"oci_runtime"`
	CgroupParent string                   `json:"cgroup_parent" structs:"cgroup_parent"`
	CPU          int64                    `json:"cpu" structs:"cpu"`
	Memory       int64                    `json:"memory" structs:"memory"`
}
