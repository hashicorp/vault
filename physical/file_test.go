package physical

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func TestFileBackend_Base64URLEncoding(t *testing.T) {
	backendPath, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(backendPath)

	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("file", logger, map[string]string{
		"path": backendPath,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// List the entries. Length should be zero.
	keys, err := b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: len(keys): expected: 0, actual: %d", len(keys))
	}

	// Create a storage entry without base64 encoding the file name
	rawFullPath := filepath.Join(backendPath, "_foo")
	e := &Entry{Key: "foo", Value: []byte("test")}
	f, err := os.OpenFile(
		rawFullPath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)
	if err != nil {
		t.Fatal(err)
	}
	json.NewEncoder(f).Encode(e)
	f.Close()

	// Get should work
	out, err := b.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Fatalf("bad: %v expected: %v", out, e)
	}

	// List the entries. There should be one entry.
	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("bad: len(keys): expected: 1, actual: %d", len(keys))
	}

	err = b.Put(e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// List the entries again. There should still be one entry.
	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("bad: len(keys): expected: 1, actual: %d", len(keys))
	}

	// Get should work
	out, err = b.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Fatalf("bad: %v expected: %v", out, e)
	}

	err = b.Delete("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err = b.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: entry: expected: nil, actual: %#v", e)
	}

	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: len(keys): expected: 0, actual: %d", len(keys))
	}

	f, err = os.OpenFile(
		rawFullPath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)
	if err != nil {
		t.Fatal(err)
	}
	json.NewEncoder(f).Encode(e)
	f.Close()

	keys, err = b.List("")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("bad: len(keys): expected: 1, actual: %d", len(keys))
	}
}

func TestFileBackend_ValidatePath(t *testing.T) {
	dir, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("file", logger, map[string]string{
		"path": dir,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := b.Delete("foo/bar/../zip"); err == nil {
		t.Fatal("expected error")
	}
	if err := b.Delete("foo/bar/zip"); err != nil {
		t.Fatal("did not expect error")
	}
}

func TestFileBackend(t *testing.T) {
	dir, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("file", logger, map[string]string{
		"path": dir,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
}
