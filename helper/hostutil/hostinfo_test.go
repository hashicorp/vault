// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hostutil

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/errwrap"
)

func TestCollectHostInfo(t *testing.T) {
	info, err := CollectHostInfo(context.Background())
	if err != nil && !errwrap.ContainsType(err, new(HostInfoError)) {
		t.Fatal(err)
	}

	// Get all the possible HostInfoError errors and check for the resulting
	// stat if the package is able to fetch it for the platform we're testing
	// on.
	errs := errwrap.GetAllType(err, new(HostInfoError))

	if info.Timestamp.IsZero() {
		t.Fatal("expected non-zero Timestamp")
	}
	if !checkErrTypeExists(errs, "cpu") && info.CPU == nil {
		t.Fatal("expected non-nil CPU value")
	}
	if !checkErrTypeExists(errs, "cpu_times") && info.CPUTimes == nil {
		t.Fatal("expected non-nil CPUTimes value")
	}
	if !checkErrTypeExists(errs, "disk") && info.Disk == nil {
		t.Fatal("expected non-nil Disk value")
	}
	if !checkErrTypeExists(errs, "host") && info.Host == nil {
		t.Fatal("expected non-nil Host value")
	}
	if !checkErrTypeExists(errs, "memory") && info.Memory == nil {
		t.Fatal("expected non-nil Memory value")
	}
}

// checkErrTypeExists is a helper that checks whether an particular
// HostInfoError.Type exists within a set of errors.
func checkErrTypeExists(errs []error, errType string) bool {
	for _, e := range errs {
		err, ok := e.(*HostInfoError)
		if !ok {
			return false
		}

		// This is mainly for disk error since the type string can contain an
		// index for the disk.
		parts := strings.SplitN(err.Type, ".", 2)

		if parts[0] == errType {
			return true
		}
	}
	return false
}
