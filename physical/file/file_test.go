package file

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
)

func TestFileBackend_Base64URLEncoding(t *testing.T) {
	backendPath, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(backendPath)

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewFileBackend(map[string]string{
		"path": backendPath,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// List the entries. Length should be zero.
	keys, err := b.List(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: len(keys): expected: 0, actual: %d", len(keys))
	}

	// Create a storage entry without base64 encoding the file name
	rawFullPath := filepath.Join(backendPath, "_foo")
	e := &physical.Entry{Key: "foo", Value: []byte("test")}
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
	out, err := b.Get(context.Background(), "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Fatalf("bad: %v expected: %v", out, e)
	}

	// List the entries. There should be one entry.
	keys, err = b.List(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("bad: len(keys): expected: 1, actual: %d", len(keys))
	}

	err = b.Put(context.Background(), e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// List the entries again. There should still be one entry.
	keys, err = b.List(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("bad: len(keys): expected: 1, actual: %d", len(keys))
	}

	// Get should work
	out, err = b.Get(context.Background(), "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, e) {
		t.Fatalf("bad: %v expected: %v", out, e)
	}

	err = b.Delete(context.Background(), "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err = b.Get(context.Background(), "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: entry: expected: nil, actual: %#v", e)
	}

	keys, err = b.List(context.Background(), "")
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

	keys, err = b.List(context.Background(), "")
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

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewFileBackend(map[string]string{
		"path": dir,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := b.Delete(context.Background(), "foo/bar/../zip"); err == nil {
		t.Fatal("expected error")
	}
	if err := b.Delete(context.Background(), "foo/bar/zip"); err != nil {
		t.Fatal("did not expect error")
	}
}

func TestFileBackend(t *testing.T) {
	dir, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewFileBackend(map[string]string{
		"path": dir,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)

	// Underscores should not trip things up; ref GH-3476
	e := &physical.Entry{Key: "_zip", Value: []byte("foobar")}
	err = b.Put(context.Background(), e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	e = &physical.Entry{Key: "_zip/_zap", Value: []byte("boofar")}
	err = b.Put(context.Background(), e)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	e, err = b.Get(context.Background(), "_zip/_zap")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if e == nil {
		t.Fatal("got nil entry")
	}
	vals, err := b.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 || vals[0] == vals[1] {
		t.Fatalf("bad: %v", vals)
	}
	for _, val := range vals {
		if val != "_zip/" && val != "_zip" {
			t.Fatalf("bad val: %v", val)
		}
	}
	vals, err = b.List(context.Background(), "_zip/")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 1 || vals[0] != "_zap" {
		t.Fatalf("bad: %v", vals)
	}
	err = b.Delete(context.Background(), "_zip/_zap")
	if err != nil {
		t.Fatal(err)
	}
	vals, err = b.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 1 || vals[0] != "_zip" {
		t.Fatalf("bad: %v", vals)
	}
	err = b.Delete(context.Background(), "_zip")
	if err != nil {
		t.Fatal(err)
	}
	vals, err = b.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 0 {
		t.Fatalf("bad: %v", vals)
	}

	physical.ExerciseBackend_ListPrefix(t, b)
}
