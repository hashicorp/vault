// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package namespace

import (
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
