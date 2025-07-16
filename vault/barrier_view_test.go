// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestBarrierView_impl(t *testing.T) {
	var _ logical.Storage = new(BarrierView)
}

func TestBarrierView_spec(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")
	logical.TestStorage(t, view)
}

func TestBarrierView_BadKeysKeys(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")

	_, err := view.List(context.Background(), "../")
	if err == nil {
		t.Fatalf("expected error")
	}

	_, err = view.Get(context.Background(), "../")
	if err == nil {
		t.Fatalf("expected error")
	}

	err = view.Delete(context.Background(), "../foo")
	if err == nil {
		t.Fatalf("expected error")
	}

	le := &logical.StorageEntry{
		Key:   "../foo",
		Value: []byte("test"),
	}
	err = view.Put(context.Background(), le)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestBarrierView(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")

	// Write a key outside of foo/
	entry := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	if err := barrier.Put(context.Background(), entry); err != nil {
		t.Fatalf("bad: %v", err)
	}

	// List should have no visibility
	keys, err := view.List(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: %v", err)
	}

	// Get should have no visibility
	out, err := view.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Try to put the same entry via the view
	if err := view.Put(context.Background(), entry); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check it is nested
	entry, err = barrier.Get(context.Background(), "foo/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry == nil {
		t.Fatalf("missing nested foo/test")
	}

	// Delete nested
	if err := view.Delete(context.Background(), "test"); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the nested key
	entry, err = barrier.Get(context.Background(), "foo/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry != nil {
		t.Fatalf("nested foo/test should be gone")
	}

	// Check the non-nested key
	entry, err = barrier.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry == nil {
		t.Fatalf("root test missing")
	}
}

func TestBarrierView_SubView(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	root := NewBarrierView(barrier, "foo/")
	view := root.SubView("bar/")

	// List should have no visibility
	keys, err := view.List(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("bad: %v", err)
	}

	// Get should have no visibility
	out, err := view.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	// Try to put the same entry via the view
	entry := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	if err := view.Put(context.Background(), entry); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check it is nested
	bout, err := barrier.Get(context.Background(), "foo/bar/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if bout == nil {
		t.Fatalf("missing nested foo/bar/test")
	}

	// Check for visibility in root
	out, err = root.Get(context.Background(), "bar/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("missing nested bar/test")
	}

	// Delete nested
	if err := view.Delete(context.Background(), "test"); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the nested key
	bout, err = barrier.Get(context.Background(), "foo/bar/test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if bout != nil {
		t.Fatalf("nested foo/bar/test should be gone")
	}
}

func TestBarrierView_Scan(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "view/")

	expect := []string{}
	ent := []*logical.StorageEntry{
		{Key: "foo", Value: []byte("test")},
		{Key: "zip", Value: []byte("test")},
		{Key: "foo/bar", Value: []byte("test")},
		{Key: "foo/zap", Value: []byte("test")},
		{Key: "foo/bar/baz", Value: []byte("test")},
		{Key: "foo/bar/zoo", Value: []byte("test")},
	}

	for _, e := range ent {
		expect = append(expect, e.Key)
		if err := view.Put(context.Background(), e); err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	var out []string
	cb := func(path string) {
		out = append(out, path)
	}

	// Collect the keys
	if err := logical.ScanView(context.Background(), view, cb); err != nil {
		t.Fatalf("err: %v", err)
	}

	sort.Strings(out)
	sort.Strings(expect)
	if !reflect.DeepEqual(out, expect) {
		t.Fatalf("out: %v expect: %v", out, expect)
	}
}

func TestBarrierView_CollectKeys(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "view/")

	expect := []string{}
	ent := []*logical.StorageEntry{
		{Key: "foo", Value: []byte("test")},
		{Key: "zip", Value: []byte("test")},
		{Key: "foo/bar", Value: []byte("test")},
		{Key: "foo/zap", Value: []byte("test")},
		{Key: "foo/bar/baz", Value: []byte("test")},
		{Key: "foo/bar/zoo", Value: []byte("test")},
	}

	for _, e := range ent {
		expect = append(expect, e.Key)
		if err := view.Put(context.Background(), e); err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Collect the keys
	out, err := logical.CollectKeys(context.Background(), view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	sort.Strings(out)
	sort.Strings(expect)
	if !reflect.DeepEqual(out, expect) {
		t.Fatalf("out: %v expect: %v", out, expect)
	}
}

func TestBarrierView_ClearView(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "view/")

	expect := []string{}
	ent := []*logical.StorageEntry{
		{Key: "foo", Value: []byte("test")},
		{Key: "zip", Value: []byte("test")},
		{Key: "foo/bar", Value: []byte("test")},
		{Key: "foo/zap", Value: []byte("test")},
		{Key: "foo/bar/baz", Value: []byte("test")},
		{Key: "foo/bar/zoo", Value: []byte("test")},
	}

	for _, e := range ent {
		expect = append(expect, e.Key)
		if err := view.Put(context.Background(), e); err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Clear the keys
	if err := logical.ClearView(context.Background(), view); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Collect the keys
	out, err := logical.CollectKeys(context.Background(), view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("have keys: %#v", out)
	}
}

func TestBarrierView_Readonly(t *testing.T) {
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "foo/")

	// Add a key before enabling read-only
	entry := &logical.StorageEntry{Key: "test", Value: []byte("test")}
	if err := view.Put(context.Background(), entry); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Enable read only mode
	view.readOnlyErr = logical.ErrReadOnly

	// Put should fail in readonly mode
	if err := view.Put(context.Background(), entry); err != logical.ErrReadOnly {
		t.Fatalf("err: %v", err)
	}

	// Delete nested
	if err := view.Delete(context.Background(), "test"); err != logical.ErrReadOnly {
		t.Fatalf("err: %v", err)
	}

	// Check the non-nested key
	e, err := view.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if e == nil {
		t.Fatalf("key test missing")
	}
}
