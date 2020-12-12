// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// ListDatabasesOptions represents options that can be used to configure a ListDatabases operation.
type ListDatabasesOptions struct {
	// If true, only the Name field of the returned DatabaseSpecification objects will be populated. The default value
	// is false.
	NameOnly *bool

	// If true, only the databases which the user is authorized to see will be returned. For more information about
	// the behavior of this option, see https://docs.mongodb.com/manual/reference/privilege-actions/#find. The default
	// value is true.
	AuthorizedDatabases *bool
}

// ListDatabases creates a new ListDatabasesOptions instance.
func ListDatabases() *ListDatabasesOptions {
	return &ListDatabasesOptions{}
}

// SetNameOnly sets the value for the NameOnly field.
func (ld *ListDatabasesOptions) SetNameOnly(b bool) *ListDatabasesOptions {
	ld.NameOnly = &b
	return ld
}

// SetAuthorizedDatabases sets the value for the AuthorizedDatabases field.
func (ld *ListDatabasesOptions) SetAuthorizedDatabases(b bool) *ListDatabasesOptions {
	ld.AuthorizedDatabases = &b
	return ld
}

// MergeListDatabasesOptions combines the given ListDatabasesOptions instances into a single *ListDatabasesOptions in a
// last-one-wins fashion.
func MergeListDatabasesOptions(opts ...*ListDatabasesOptions) *ListDatabasesOptions {
	ld := ListDatabases()
	for _, opt := range opts {
		if opts == nil {
			continue
		}
		if opt.NameOnly != nil {
			ld.NameOnly = opt.NameOnly
		}
		if opt.AuthorizedDatabases != nil {
			ld.AuthorizedDatabases = opt.AuthorizedDatabases
		}
	}

	return ld
}
