// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

type Inspectable interface {
	// Returns a record view of a particular subsystem
	GetRecords(tag string) ([]map[string]interface{}, error)
}

type Deserializable interface {
	// Converts a structure into a consummable map
	Deserialize() map[string]interface{}
}
