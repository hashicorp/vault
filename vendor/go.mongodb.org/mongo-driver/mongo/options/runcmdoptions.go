// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import "go.mongodb.org/mongo-driver/mongo/readpref"

// RunCmdOptions represents all possible options for a runCommand operation.
type RunCmdOptions struct {
	ReadPreference *readpref.ReadPref // The read preference for the operation.
}

// RunCmd creates a new *RunCmdOptions
func RunCmd() *RunCmdOptions {
	return &RunCmdOptions{}
}

// SetReadPreference sets the read preference for the operation.
func (rc *RunCmdOptions) SetReadPreference(rp *readpref.ReadPref) *RunCmdOptions {
	rc.ReadPreference = rp
	return rc
}

// MergeRunCmdOptions combines the given *RunCmdOptions into one *RunCmdOptions in a last one wins fashion.
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
