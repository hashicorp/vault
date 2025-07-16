// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
			t.Fatalf("path should not exist in filtered paths %q", path)
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
			t.Fatalf("path should exist in filtered paths %q", path)
		}
	}

	// remove the paths
	m.RemovePaths(paths)

	for _, path := range paths {
		if m.HasPath(path) {
			t.Fatalf("path should not exist in filtered paths %q", path)
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
			t.Fatalf("path should not exist in filtered paths %q", path)
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
			t.Fatalf("path should exist in filtered paths %q", path)
		}
	}

	// remove the paths
	m.RemovePathPrefix("path")

	if m.Len() != 0 {
		t.Fatalf("bad: path length expect 0, got %d", len(m.Paths()))
	}

	for _, path := range paths {
		if m.HasPath(path) {
			t.Fatalf("path should not exist in filtered paths %q", path)
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
		{"path1/key1", true},
		{"path2/key1", true},
		{"path3/key1", true},
		{"path1/key1/subkey1", true},
		{"path1/key1/subkey99", false},
		{"path2/key1/subkey1", true},
		{"path1/key1/subkey1/subkey1", false},
		{"nonexistentpath/key1", false},
		{"path4/key1", false},
		{"path5/key1/subkey1", false},
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

func TestPathManager_HasPath(t *testing.T) {
	m := New()

	m.AddPaths([]string{"a/b/c/"})
	if m.HasPath("a/") {
		t.Fatal("should not have path 'a/'")
	}
	if m.HasPath("a/b/") {
		t.Fatal("should not have path 'a/b/'")
	}
	if !m.HasPath("a/b/c/") {
		t.Fatal("should have path 'a/b/c'")
	}

	m.AddPaths([]string{"a/"})
	if !m.HasPath("a/") {
		t.Fatal("should have path 'a/'")
	}
	if !m.HasPath("a/b/") {
		t.Fatal("should have path 'a/b/'")
	}
	if !m.HasPath("a/b/c/") {
		t.Fatal("should have path 'a/b/c'")
	}

	m.RemovePaths([]string{"a/"})
	if m.HasPath("a/") {
		t.Fatal("should not have path 'a/'")
	}
	if m.HasPath("a/b/") {
		t.Fatal("should not have path 'a/b/'")
	}
	if !m.HasPath("a/b/c/") {
		t.Fatal("should have path 'a/b/c'")
	}
}
