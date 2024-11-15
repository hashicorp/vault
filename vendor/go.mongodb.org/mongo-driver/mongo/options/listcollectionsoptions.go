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

	// The maximum number of documents to be included in each batch returned by the server.
	BatchSize *int32

	// If true, and NameOnly is true, limits the documents returned to only contain collections the user is authorized to use. The default value
	// is false. This option is only valid for MongoDB server versions >= 4.0. Server versions < 4.0 ignore this option.
	AuthorizedCollections *bool
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

// SetBatchSize sets the value for the BatchSize field.
func (lc *ListCollectionsOptions) SetBatchSize(size int32) *ListCollectionsOptions {
	lc.BatchSize = &size
	return lc
}

// SetAuthorizedCollections sets the value for the AuthorizedCollections field. This option is only valid for MongoDB server versions >= 4.0. Server
// versions < 4.0 ignore this option.
func (lc *ListCollectionsOptions) SetAuthorizedCollections(b bool) *ListCollectionsOptions {
	lc.AuthorizedCollections = &b
	return lc
}

// MergeListCollectionsOptions combines the given ListCollectionsOptions instances into a single *ListCollectionsOptions
// in a last-one-wins fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeListCollectionsOptions(opts ...*ListCollectionsOptions) *ListCollectionsOptions {
	lc := ListCollections()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.NameOnly != nil {
			lc.NameOnly = opt.NameOnly
		}
		if opt.BatchSize != nil {
			lc.BatchSize = opt.BatchSize
		}
		if opt.AuthorizedCollections != nil {
			lc.AuthorizedCollections = opt.AuthorizedCollections
		}
	}

	return lc
}
