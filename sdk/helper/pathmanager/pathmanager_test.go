package pathmanager

import (
	"reflect"
	"testing"
)

func TestPathManager(t *testing.T) {
	m := New()

	if m.Len() != 0 {
		t.Fatalf("bad: path length expect 0, got %d", len(m.Paths()))
	}

	paths := []string{
		"path1/",
		"path2/",
		"path3/",
	}

	for _, path := range paths {
		if m.HasPath(path) {
			t.Fatalf("path should not exist in filtered paths '%s'", path)
		}
	}

	// add paths
	m.AddPaths(paths)
	if m.Len() != 3 {
		t.Fatalf("bad: path length expect 3, got %d", len(m.Paths()))
	}
	if !reflect.DeepEqual(paths, m.Paths()) {
		t.Fatalf("mismatch in paths")
	}
	for _, path := range paths {
		if !m.HasPath(path) {
			t.Fatalf("path should exist in filtered paths '%s'", path)
		}
	}

	// remove the paths
	m.RemovePaths(paths)

	for _, path := range paths {
		if m.HasPath(path) {
			t.Fatalf("path should not exist in filtered paths '%s'", path)
		}
	}
}

func TestPathManager_RemovePrefix(t *testing.T) {
	m := New()

	if m.Len() != 0 {
		t.Fatalf("bad: path length expect 0, got %d", len(m.Paths()))
	}

	paths := []string{
		"path1/",
		"path2/",
		"path3/",
	}

	for _, path := range paths {
		if m.HasPath(path) {
			t.Fatalf("path should not exist in filtered paths '%s'", path)
		}
	}

	// add paths
	m.AddPaths(paths)
	if m.Len() != 3 {
		t.Fatalf("bad: path length expect 3, got %d", len(m.Paths()))
	}
	if !reflect.DeepEqual(paths, m.Paths()) {
		t.Fatalf("mismatch in paths")
	}
	for _, path := range paths {
		if !m.HasPath(path) {
			t.Fatalf("path should exist in filtered paths '%s'", path)
		}
	}

	// remove the paths
	m.RemovePathPrefix("path")

	if m.Len() != 0 {
		t.Fatalf("bad: path length expect 0, got %d", len(m.Paths()))
	}

	for _, path := range paths {
		if m.HasPath(path) {
			t.Fatalf("path should not exist in filtered paths '%s'", path)
		}
	}
}

func TestPathManager_HasExactPath(t *testing.T) {
	m := New()
	paths := []string{
		"path1/key1",
		"path1/key1/subkey1",
		"path1/key1/subkey2",
		"path1/key1/subkey3",
		"path2/*",
		"path3/",
		"!path4/key1",
		"!path5/*",
	}
	m.AddPaths(paths)
	if m.Len() != len(paths) {
		t.Fatalf("path count does not match: expected %d, got %d", len(paths), m.Len())
	}

	type tCase struct {
		key    string
		expect bool
	}

	tcases := []tCase{
		tCase{"path1/key1", true},
		tCase{"path2/key1", true},
		tCase{"path3/key1", true},
		tCase{"path1/key1/subkey1", true},
		tCase{"path1/key1/subkey99", false},
		tCase{"path2/key1/subkey1", true},
		tCase{"path1/key1/subkey1/subkey1", false},
		tCase{"nonexistentpath/key1", false},
		tCase{"path4/key1", false},
		tCase{"path5/key1/subkey1", false},
	}

	for _, tc := range tcases {
		if match := m.HasExactPath(tc.key); match != tc.expect {
			t.Fatalf("incorrect match: key %q", tc.key)
		}
	}

	m.RemovePaths(paths)
	if len(m.Paths()) != 0 {
		t.Fatalf("removing all paths did not clear manager: paths %v", m.Paths())
	}
}
