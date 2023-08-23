// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keysutil

type Cache interface {
	Delete(key interface{})
	Load(key interface{}) (value interface{}, ok bool)
	Store(key, value interface{})
	Size() int
}
