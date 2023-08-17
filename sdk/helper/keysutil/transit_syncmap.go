// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keysutil

import (
	"sync"
)

type TransitSyncMap struct {
	syncmap sync.Map
}

func NewTransitSyncMap() *TransitSyncMap {
	return &TransitSyncMap{syncmap: sync.Map{}}
}

func (c *TransitSyncMap) Delete(key interface{}) {
	c.syncmap.Delete(key)
}

func (c *TransitSyncMap) Load(key interface{}) (value interface{}, ok bool) {
	return c.syncmap.Load(key)
}

func (c *TransitSyncMap) Store(key, value interface{}) {
	c.syncmap.Store(key, value)
}

func (c *TransitSyncMap) Size() int {
	return 0
}
