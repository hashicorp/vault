// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package capabilities

const (
	APICapability            = "api"
	MetaCapability           = "meta"
	APIPassThroughCapability = "passthrough"
	LinkControlCapability    = "link-control"
)

type Capability interface {
	Start() error
	Stop() error
}
