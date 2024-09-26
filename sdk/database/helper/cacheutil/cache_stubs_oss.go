// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package cacheutil

import "errors"

type EvictionFunc func(key interface{}, value interface{})

type Cache struct{}

func NewCache(_ int, _ EvictionFunc) (*Cache, error) {
	return nil, errors.New("self-managed static roles only available in Vault Enterprise")
}
