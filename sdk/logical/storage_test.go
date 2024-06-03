// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"
	"testing"

	"github.com/go-test/deep"
)

var keyList = []string{
	"a",
	"b",
	"d",
	"foo",
	"foo42",
	"foo/a/b/c",
	"c/d/e/f/g",
}

func TestScanView(t *testing.T) {
	s := prepKeyStorage(t)

	keys := make([]string, 0)
	err := ScanView(context.Background(), s, func(path string) {
		keys = append(keys, path)
	})
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(keys, keyList); diff != nil {
		t.Fatal(diff)
	}
}

func TestScanView_CancelContext(t *testing.T) {
	s := prepKeyStorage(t)

	ctx, cancelCtx := context.WithCancel(context.Background())
	var i int
	err := ScanView(ctx, s, func(path string) {
		cancelCtx()
		i++
	})

	if err == nil {
		t.Error("Want context cancel err, got none")
	}
	if i != 1 {
		t.Errorf("Want i==1, got %d", i)
	}
}

func TestCollectKeys(t *testing.T) {
	s := prepKeyStorage(t)

	keys, err := CollectKeys(context.Background(), s)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(keys, keyList); diff != nil {
		t.Fatal(diff)
	}
}

func TestCollectKeysPrefix(t *testing.T) {
	s := prepKeyStorage(t)

	keys, err := CollectKeysWithPrefix(context.Background(), s, "foo")
	if err != nil {
		t.Fatal(err)
	}

	exp := []string{
		"foo",
		"foo42",
		"foo/a/b/c",
	}

	if diff := deep.Equal(keys, exp); diff != nil {
		t.Fatal(diff)
	}
}

func prepKeyStorage(t *testing.T) Storage {
	t.Helper()
	s := &InmemStorage{}

	for _, key := range keyList {
		if err := s.Put(context.Background(), &StorageEntry{
			Key:      key,
			Value:    nil,
			SealWrap: false,
		}); err != nil {
			t.Fatal(err)
		}
	}

	return s
}
