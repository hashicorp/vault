// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tokenhelper

import (
	"testing"
)

// test is a public function that can be used in other tests to
// test that a helper is functioning properly.
func test(t *testing.T, h TokenHelper) {
	if err := h.Store("foo"); err != nil {
		t.Fatalf("err: %s", err)
	}

	v, err := h.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if v != "foo" {
		t.Fatalf("bad: %#v", v)
	}

	if err := h.Erase(); err != nil {
		t.Fatalf("err: %s", err)
	}

	v, err = h.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if v != "" {
		t.Fatalf("bad: %#v", v)
	}
}
