// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// ListCollectionsOptions represents options that can be used to configure a ListCollections operation.
type ListCollectionsOptions struct {
	// If true, each collection document will only contain a field for the collection name. The default value is false.
	NameOnly *bool
}

// ListCollections creates a new ListCollectionsOptions instance.
func ListCollections() *ListCollectionsOptions {
	return &ListCollectionsOptions{}
}

// SetNameOnly sets the value for the NameOnly field.
func (lc *ListCollectionsOptions) SetNameOnly(b bool) *ListCollectionsOptions {
	lc.NameOnly = &b
	return lc
}

// MergeListCollectionsOptions combines the given ListCollectionsOptions instances into a single *ListCollectionsOptions
// in a last-one-wins fashion.
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
