// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
