// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package osutil

import (
	"io/fs"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"testing"
)

func TestCheckPathInfo(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Errorf("failed to get details of current process owner. The error is: %v", err)
	}
	uid, err := strconv.ParseInt(currentUser.Uid, 0, 64)
	if err != nil {
		t.Errorf("failed to convert uid to int64. The error is: %v", err)
	}
	uid2, err := strconv.ParseInt(currentUser.Uid+"1", 0, 64)
	if err != nil {
		t.Errorf("failed to convert uid to int64. The error is: %v", err)
	}

	testCases := []struct {
		uid             int
		filepermissions fs.FileMode
		permissions     int
		expectError     bool
	}{
		{
			uid:             0,
			filepermissions: 0o700,
			permissions:     0,
			expectError:     false,
		},
		{
			uid:             int(uid2),
			filepermissions: 0o700,
			permissions:     0,
			expectError:     true,
		},
		{
			uid:             int(uid),
			filepermissions: 0o700,
			permissions:     0,
			expectError:     false,
		},
		{
			uid:             0,
			filepermissions: 0o777,
			permissions:     744,
			expectError:     true,
		},
	}

	for _, tc := range testCases {
		err := os.Mkdir("testFile", tc.filepermissions)
		if err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat("testFile")
		if err != nil {
			t.Errorf("error stating %q: %v", "testFile", err)
		}
		if tc.uid != 0 && runtime.GOOS == "windows" && tc.expectError == true {
			t.Skip("Skipping test in windows environment as no error will be returned in this case")
		}

		err = checkPathInfo(info, "testFile", tc.uid, int(tc.permissions))
		if tc.expectError && err == nil {
			t.Errorf("invalid result. expected error")
		}
		if !tc.expectError && err != nil {
			t.Errorf(err.Error())
		}

		err = os.RemoveAll("testFile")
		if err != nil {
			t.Fatal(err)
		}
	}
}
