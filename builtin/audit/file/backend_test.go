// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package file

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestAuditFile_fileModeNew(t *testing.T) {
	modeStr := "0777"
	mode, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		t.Fatal(err)
	}

	path, err := ioutil.TempDir("", "vault-test_audit_file-file_mode_new")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(path)

	file := filepath.Join(path, "auditTest.txt")

	config := map[string]string{
		"path": file,
		"mode": modeStr,
	}

	_, err = Factory(context.Background(), &audit.BackendConfig{
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
		Config:     config,
	})
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(file)
	if err != nil {
		t.Fatalf("Cannot retrieve file mode from `Stat`")
	}
	if info.Mode() != os.FileMode(mode) {
		t.Fatalf("File mode does not match.")
	}
}

func TestAuditFile_fileModeExisting(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("Failure to create test file.")
	}
	defer os.Remove(f.Name())

	err = os.Chmod(f.Name(), 0o777)
	if err != nil {
		t.Fatalf("Failure to chmod temp file for testing.")
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("Failure to close temp file for test.")
	}

	config := map[string]string{
		"path": f.Name(),
	}

	_, err = Factory(context.Background(), &audit.BackendConfig{
		Config:     config,
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
	})
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(f.Name())
	if err != nil {
		t.Fatalf("cannot retrieve file mode from `Stat`")
	}
	if info.Mode() != os.FileMode(0o600) {
		t.Fatalf("File mode does not match.")
	}
}

func TestAuditFile_fileMode0000(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("Failure to create test file. The error is %v", err)
	}
	defer os.Remove(f.Name())

	err = os.Chmod(f.Name(), 0o777)
	if err != nil {
		t.Fatalf("Failure to chmod temp file for testing. The error is %v", err)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("Failure to close temp file for test. The error is %v", err)
	}

	config := map[string]string{
		"path": f.Name(),
		"mode": "0000",
	}

	_, err = Factory(context.Background(), &audit.BackendConfig{
		Config:     config,
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
	})
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(f.Name())
	if err != nil {
		t.Fatalf("cannot retrieve file mode from `Stat`. The error is %v", err)
	}
	if info.Mode() != os.FileMode(0o777) {
		t.Fatalf("File mode does not match.")
	}
}

func BenchmarkAuditFile_request(b *testing.B) {
	config := map[string]string{
		"path": "/dev/null",
	}
	sink, err := Factory(context.Background(), &audit.BackendConfig{
		Config:     config,
		SaltConfig: &salt.Config{},
		SaltView:   &logical.InmemStorage{},
	})
	if err != nil {
		b.Fatal(err)
	}

	in := &logical.LogInput{
		Auth: &logical.Auth{
			ClientToken:     "foo",
			Accessor:        "bar",
			EntityID:        "foobarentity",
			DisplayName:     "testtoken",
			NoDefaultPolicy: true,
			Policies:        []string{"root"},
			TokenType:       logical.TokenTypeService,
		},
		Request: &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "/foo",
			Connection: &logical.Connection{
				RemoteAddr: "127.0.0.1",
			},
			WrapInfo: &logical.RequestWrapInfo{
				TTL: 60 * time.Second,
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
	}

	ctx := namespace.RootContext(nil)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := sink.LogRequest(ctx, in); err != nil {
				panic(err)
			}
		}
	})
}
