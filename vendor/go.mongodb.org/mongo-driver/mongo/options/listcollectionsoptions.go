// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// ListCollectionsOptions represents all possible options for a listCollections command.
type ListCollectionsOptions struct {
	NameOnly *bool // If true, only the collection names will be returned.
}

// ListCollections creates a new *ListCollectionsOptions
func ListCollections() *ListCollectionsOptions {
	return &ListCollectionsOptions{}
}

// SetNameOnly specifies whether to return only the collection names.
func (lc *ListCollectionsOptions) SetNameOnly(b bool) *ListCollectionsOptions {
	lc.NameOnly = &b
	return lc
}

// MergeListCollectionsOptions combines the given *ListCollectionsOptions into a single *ListCollectionsOptions in a
// last one wins fashion.
func MergeListCollectionsOptions(opts ...*ListCollectionsOptions) *ListCollectionsOptions {
	lc := ListCollections()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.NameOnly != nil {
			lc.NameOnly = opt.NameOnly
		}
	}

	return lc
}
