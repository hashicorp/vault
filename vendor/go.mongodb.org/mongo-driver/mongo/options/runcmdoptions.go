// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// RunCmdOptions represents options that can be used to configure a RunCommand operation.
type RunCmdOptions struct {
	// The read preference to use for the operation. The default value is nil, which means that the primary read
	// preference will be used.
	ReadPreference *readpref.ReadPref
}

// RunCmd creates a new RunCmdOptions instance.
func RunCmd() *RunCmdOptions {
	return &RunCmdOptions{}
}

// SetReadPreference sets value for the ReadPreference field.
func (rc *RunCmdOptions) SetReadPreference(rp *readpref.ReadPref) *RunCmdOptions {
	rc.ReadPreference = rp
	return rc
}

// MergeRunCmdOptions combines the given RunCmdOptions instances into one *RunCmdOptions in a last-one-wins fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeRunCmdOptions(opts ...*RunCmdOptions) *RunCmdOptions {
	rc := RunCmd()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.ReadPreference != nil {
			rc.ReadPreference = opt.ReadPreference
		}
	}

	return rc
}
