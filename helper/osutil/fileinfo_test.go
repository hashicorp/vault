// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package osutil

import (
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
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
			t.Error("invalid result. expected error")
		}
		if !tc.expectError && err != nil {
			t.Error(err.Error())
		}

		err = os.RemoveAll("testFile")
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestOwnerPermissionsMatchFile creates a file and verifies that the current user of the process is the owner of the
// file
func TestOwnerPermissionsMatchFile(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatal("failed to get current user", err)
	}
	uid, err := strconv.ParseInt(currentUser.Uid, 0, 64)
	if err != nil {
		t.Fatal("failed to convert uid", err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "foo")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal("failed to create test file", err)
	}
	defer f.Close()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal("failed to stat test file", err)
	}

	if err := OwnerPermissionsMatchFile(f, int(uid), int(info.Mode())); err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}

// TestOwnerPermissionsMatchFile_OtherUser creates a file using the user that started the current process and verifies
// that a different user is not the owner of the file
func TestOwnerPermissionsMatchFile_OtherUser(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatal("failed to get current user", err)
	}
	uid, err := strconv.ParseInt(currentUser.Uid, 0, 64)
	if err != nil {
		t.Fatal("failed to convert uid", err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "foo")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal("failed to create test file", err)
	}
	defer f.Close()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal("failed to stat test file", err)
	}

	if err := OwnerPermissionsMatchFile(f, int(uid)+1, int(info.Mode())); err == nil {
		t.Fatalf("expected error but none")
	}
}

// TestOwnerPermissionsMatchFile_Symlink creates a file and a symlink to that file. The test verifies that the current
// user of the process is the owner of the file
func TestOwnerPermissionsMatchFile_Symlink(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatal("failed to get current user", err)
	}
	uid, err := strconv.ParseInt(currentUser.Uid, 0, 64)
	if err != nil {
		t.Fatal("failed to convert uid", err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "foo")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal("failed to create test file", err)
	}
	defer f.Close()

	symlink := filepath.Join(dir, "symlink")
	err = os.Symlink(path, symlink)
	if err != nil {
		t.Fatal("failed to symlink file", err)
	}
	symlinkedFile, err := os.Open(symlink)
	if err != nil {
		t.Fatal("failed to open file", err)
	}
	info, err := os.Stat(symlink)
	if err != nil {
		t.Fatal("failed to stat test file", err)
	}
	if err := OwnerPermissionsMatchFile(symlinkedFile, int(uid), int(info.Mode())); err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
}
