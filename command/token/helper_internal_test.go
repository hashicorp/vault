// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package token

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestCommand re-uses the existing Test function to ensure proper behavior of
// the internal token helper
func TestCommand(t *testing.T) {
	helper, err := NewInternalTokenHelper()
	if err != nil {
		t.Fatal(err)
	}
	Test(t, helper)
}

func TestInternalHelperFilePerms(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	helper, err := NewInternalTokenHelper()
	if err != nil {
		t.Fatal(err)
	}
	helper.homeDir = tmpDir

	tmpFile := filepath.Join(tmpDir, ".vault-token")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fi, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if fi.Mode().Perm()&0o04 != 0o04 {
		t.Fatalf("expected world-readable/writable permission bits, got: %o", fi.Mode().Perm())
	}

	err = helper.Store("bogus_token")
	if err != nil {
		t.Fatal(err)
	}

	fi, err = os.Stat(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if fi.Mode().Perm()&0o04 != 0 {
		t.Fatalf("expected no world-readable/writable permission bits, got: %o", fi.Mode().Perm())
	}
}
