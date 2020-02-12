// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// ListDatabasesOptions represents all possible options for a listDatabases command.
type ListDatabasesOptions struct {
	NameOnly *bool // If true, only the database names will be returned.
}

// ListDatabases creates a new *ListDatabasesOptions
func ListDatabases() *ListDatabasesOptions {
	return &ListDatabasesOptions{}
}

// SetNameOnly specifies whether to return only the database names.
func (ld *ListDatabasesOptions) SetNameOnly(b bool) *ListDatabasesOptions {
	ld.NameOnly = &b
	return ld
}

// MergeListDatabasesOptions combines the given *ListDatabasesOptions into a single *ListDatabasesOptions in a last one
// wins fashion.
func MergeListDatabasesOptions(opts ...*ListDatabasesOptions) *ListDatabasesOptions {
	ld := ListDatabases()
	for _, opt := range opts {
		if opts == nil {
			continue
		}
		if opt.NameOnly != nil {
			ld.NameOnly = opt.NameOnly
		}
	}

	return ld
}
