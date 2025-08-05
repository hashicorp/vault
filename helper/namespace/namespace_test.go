// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package namespace

import (
	"context"
	"testing"
)

func TestSplitIDFromString(t *testing.T) {
	tcases := []struct {
		input  string
		id     string
		prefix string
	}{
		{
			"foo",
			"",
			"foo",
		},
		{
			"foo.id",
			"id",
			"foo",
		},
		{
			"foo.foo.id",
			"id",
			"foo.foo",
		},
		{
			"foo.foo/foo.id",
			"id",
			"foo.foo/foo",
		},
		{
			"foo.foo/.id",
			"id",
			"foo.foo/",
		},
		{
			"foo.foo/foo",
			"",
			"foo.foo/foo",
		},
		{
			"foo.foo/f",
			"",
			"foo.foo/f",
		},
		{
			"foo.foo/",
			"",
			"foo.foo/",
		},
		{
			"b.foo",
			"",
			"b.foo",
		},
		{
			"s.foo",
			"",
			"s.foo",
		},
		{
			"t.foo",
			"foo",
			"t",
		},
	}

	for _, c := range tcases {
		pre, id := SplitIDFromString(c.input)
		if pre != c.prefix || id != c.id {
			t.Fatalf("bad test case: %s != %s, %s != %s", pre, c.prefix, id, c.id)
		}
	}
}

func TestHasParent(t *testing.T) {
	// Create ns1
	ns1 := &Namespace{
		ID:   "id1",
		Path: "ns1/",
	}

	// Create ns1/ns2
	ns2 := &Namespace{
		ID:   "id2",
		Path: "ns1/ns2/",
	}

	// Create ns1/ns2/ns3
	ns3 := &Namespace{
		ID:   "id3",
		Path: "ns1/ns2/ns3/",
	}

	// Create ns4
	ns4 := &Namespace{
		ID:   "id4",
		Path: "ns4/",
	}

	// Create ns4/ns5
	ns5 := &Namespace{
		ID:   "id5",
		Path: "ns4/ns5/",
	}

	tests := []struct {
		name     string
		parent   *Namespace
		ns       *Namespace
		expected bool
	}{
		{
			"is root an ancestor of ns1",
			RootNamespace,
			ns1,
			true,
		},
		{
			"is ns1 an ancestor of ns2",
			ns1,
			ns2,
			true,
		},
		{
			"is ns2 an ancestor of ns3",
			ns2,
			ns3,
			true,
		},
		{
			"is ns1 an ancestor of ns3",
			ns1,
			ns3,
			true,
		},
		{
			"is root an ancestor of ns3",
			RootNamespace,
			ns3,
			true,
		},
		{
			"is ns4 an ancestor of ns3",
			ns4,
			ns3,
			false,
		},
		{
			"is ns5 an ancestor of ns3",
			ns5,
			ns3,
			false,
		},
		{
			"is ns1 an ancestor of ns5",
			ns1,
			ns5,
			false,
		},
	}

	for _, test := range tests {
		actual := test.ns.HasParent(test.parent)
		if actual != test.expected {
			t.Fatalf("bad ancestor calculation; name: %q, actual: %t, expected: %t", test.name, actual, test.expected)
		}
	}
}

func TestContextWithNamespace(t *testing.T) {
	ns1 := &Namespace{
		ID:   "id1",
		Path: "ns1/path/",
		CustomMetadata: map[string]string{
			"key1": "value1",
		},
	}

	ns2 := &Namespace{
		ID:   "id2",
		Path: "ns2/path/",
	}

	tests := map[string]struct {
		inputCtx          context.Context
		inputNamespace    *Namespace
		expectedNamespace *Namespace
		expectedErrorMsg  string
	}{
		"nil namespace": {
			inputCtx:         context.Background(),
			inputNamespace:   nil,
			expectedErrorMsg: ErrNoNamespace.Error(),
		},
		"valid context with custom namespace": {
			inputCtx:          context.Background(),
			inputNamespace:    ns1,
			expectedNamespace: ns1,
		},
		"valid context with root namespace": {
			inputCtx:          context.Background(),
			inputNamespace:    RootNamespace,
			expectedNamespace: RootNamespace,
		},
		"override existing namespace": {
			inputCtx:          ContextWithNamespace(context.Background(), RootNamespace),
			inputNamespace:    ns2,
			expectedNamespace: ns2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resultCtx := ContextWithNamespace(tc.inputCtx, tc.inputNamespace)
			if resultCtx == nil {
				t.Fatal("ContextWithNamespace should not return nil context")
			}

			ns, err := FromContext(resultCtx)

			if tc.expectedErrorMsg != "" {
				if err == nil {
					t.Fatalf("expected error but got nil")
				} else if err.Error() != tc.expectedErrorMsg {
					t.Fatalf("expected error message %q, got %q", tc.expectedErrorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if ns != tc.expectedNamespace {
				t.Fatalf("namespace does not match expected namespace: expected %v, got %v", tc.expectedNamespace, ns)
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	ns1 := &Namespace{
		ID:             "id1",
		Path:           "ns1/path/",
		CustomMetadata: map[string]string{"key1": "value1"},
	}

	tests := map[string]struct {
		inputCtx          context.Context
		expectedNamespace *Namespace
		expectedErrorMsg  string
	}{
		"nil context": {
			inputCtx:         nil,
			expectedErrorMsg: "context was nil",
		},
		"context without namespace": {
			inputCtx:         context.Background(),
			expectedErrorMsg: ErrNoNamespace.Error(),
		},
		"context with nil namespace value": {
			inputCtx:         context.WithValue(context.Background(), contextNamespace, (*Namespace)(nil)),
			expectedErrorMsg: ErrNoNamespace.Error(),
		},
		"context with custom namespace": {
			inputCtx:          ContextWithNamespace(context.Background(), ns1),
			expectedNamespace: ns1,
		},
		"context with root namespace": {
			inputCtx:          ContextWithNamespace(context.Background(), RootNamespace),
			expectedNamespace: RootNamespace,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ns, err := FromContext(tc.inputCtx)

			if tc.expectedErrorMsg != "" {
				if err == nil {
					t.Fatalf("expected error but got nil")
				} else if err.Error() != tc.expectedErrorMsg {
					t.Fatalf("expected error message %q, got %q", tc.expectedErrorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if ns != tc.expectedNamespace {
				t.Fatalf("namespace does not match expected namespace: expected %v, got %v", tc.expectedNamespace, ns)
			}
		})
	}
}
