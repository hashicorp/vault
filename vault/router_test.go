// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

func TestRouter_Mount(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	mountEntry := &MountEntry{
		Path:        "prod/aws/",
		UUID:        meUUID,
		Accessor:    "awsaccessor",
		NamespaceID: namespace.RootNamespaceID,
		namespace:   namespace.RootNamespace,
	}

	n := &NoopBackend{}
	err = r.Mount(n, "prod/aws/", mountEntry, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	meUUID, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	err = r.Mount(n, "prod/aws/", &MountEntry{UUID: meUUID, NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace}, view)
	if !strings.Contains(err.Error(), "cannot mount under existing mount") {
		t.Fatalf("err: %v", err)
	}

	meUUID, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	if path := r.MatchingMount(namespace.RootContext(nil), "prod/aws/foo"); path != "prod/aws/" {
		t.Fatalf("bad: %s", path)
	}

	if v := r.MatchingStorageByAPIPath(namespace.RootContext(nil), "prod/aws/foo"); v.(*BarrierView) != view {
		t.Fatalf("bad: %v", v)
	}

	if path := r.MatchingMount(namespace.RootContext(nil), "stage/aws/foo"); path != "" {
		t.Fatalf("bad: %s", path)
	}

	if v := r.MatchingStorageByAPIPath(namespace.RootContext(nil), "stage/aws/foo"); v != nil {
		t.Fatalf("bad: %v", v)
	}

	mountEntryFetched := r.MatchingMountByUUID(mountEntry.UUID)
	if mountEntryFetched == nil || !reflect.DeepEqual(mountEntry, mountEntryFetched) {
		t.Fatalf("failed to fetch mount entry using its ID; expected: %#v\n actual: %#v\n", mountEntry, mountEntryFetched)
	}

	_, mount, prefix, ok := r.MatchingAPIPrefixByStoragePath(namespace.RootContext(nil), "logical/foo")
	if !ok {
		t.Fatalf("missing storage prefix")
	}
	if mount != "prod/aws/" || prefix != "logical/" {
		t.Fatalf("Bad: %v - %v", mount, prefix)
	}

	req := &logical.Request{
		Path: "prod/aws/foo",
	}
	req.SetTokenEntry(&logical.TokenEntry{
		ID: "foo",
	})
	resp, err := r.Route(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
	if req.TokenEntry() == nil || req.TokenEntry().ID != "foo" {
		t.Fatalf("unexpected value for token entry: %v", req.TokenEntry())
	}

	// Verify the path
	if len(n.Paths) != 1 || n.Paths[0] != "foo" {
		t.Fatalf("bad: %v", n.Paths)
	}

	subMountEntry := &MountEntry{
		Path:        "prod/",
		UUID:        meUUID,
		Accessor:    "prodaccessor",
		NamespaceID: namespace.RootNamespaceID,
		namespace:   namespace.RootNamespace,
	}

	if r.MountConflict(namespace.RootContext(nil), "prod/aws/") == "" {
		t.Fatalf("bad: prod/aws/")
	}

	// No error is shown here because MountConflict is checked before Mount
	err = r.Mount(n, "prod/", subMountEntry, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if r.MountConflict(namespace.RootContext(nil), "prod/test") == "" {
		t.Fatalf("bad: prod/test/")
	}
}

func TestRouter_MountCredential(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, credentialBarrierPrefix)

	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	mountEntry := &MountEntry{
		Path:        "aws",
		UUID:        meUUID,
		Accessor:    "awsaccessor",
		NamespaceID: namespace.RootNamespaceID,
		namespace:   namespace.RootNamespace,
	}

	n := &NoopBackend{}
	err = r.Mount(n, "auth/aws/", mountEntry, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	meUUID, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	err = r.Mount(n, "auth/aws/", &MountEntry{UUID: meUUID, NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace}, view)
	if !strings.Contains(err.Error(), "cannot mount under existing mount") {
		t.Fatalf("err: %v", err)
	}

	if path := r.MatchingMount(namespace.RootContext(nil), "auth/aws/foo"); path != "auth/aws/" {
		t.Fatalf("bad: %s", path)
	}

	if v := r.MatchingStorageByAPIPath(namespace.RootContext(nil), "auth/aws/foo"); v.(*BarrierView) != view {
		t.Fatalf("bad: %v", v)
	}

	if path := r.MatchingMount(namespace.RootContext(nil), "auth/stage/aws/foo"); path != "" {
		t.Fatalf("bad: %s", path)
	}

	if v := r.MatchingStorageByAPIPath(namespace.RootContext(nil), "auth/stage/aws/foo"); v != nil {
		t.Fatalf("bad: %v", v)
	}

	mountEntryFetched := r.MatchingMountByUUID(mountEntry.UUID)
	if mountEntryFetched == nil || !reflect.DeepEqual(mountEntry, mountEntryFetched) {
		t.Fatalf("failed to fetch mount entry using its ID; expected: %#v\n actual: %#v\n", mountEntry, mountEntryFetched)
	}

	_, mount, prefix, ok := r.MatchingAPIPrefixByStoragePath(namespace.RootContext(nil), "auth/foo")
	if !ok {
		t.Fatalf("missing storage prefix")
	}
	if mount != "auth/aws" || prefix != credentialBarrierPrefix {
		t.Fatalf("Bad: %v - %v", mount, prefix)
	}

	req := &logical.Request{
		Path: "auth/aws/foo",
	}
	resp, err := r.Route(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	// Verify the path
	if len(n.Paths) != 1 || n.Paths[0] != "foo" {
		t.Fatalf("bad: %v", n.Paths)
	}
}

func TestRouter_Unmount(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	n := &NoopBackend{}
	err = r.Mount(n, "prod/aws/", &MountEntry{Path: "prod/aws/", UUID: meUUID, Accessor: "awsaccessor", NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Unmount(namespace.RootContext(nil), "prod/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := &logical.Request{
		Path: "prod/aws/foo",
	}
	_, err = r.Route(namespace.RootContext(nil), req)
	if !strings.Contains(err.Error(), "unsupported path") {
		t.Fatalf("err: %v", err)
	}

	if _, _, _, ok := r.MatchingAPIPrefixByStoragePath(namespace.RootContext(nil), "logical/foo"); ok {
		t.Fatalf("should not have matching storage prefix")
	}
}

func TestRouter_Remount(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	n := &NoopBackend{}
	me := &MountEntry{Path: "prod/aws/", UUID: meUUID, Accessor: "awsaccessor", NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace}
	err = r.Mount(n, "prod/aws/", me, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	me.Path = "stage/aws/"
	err = r.Remount(namespace.RootContext(nil), "prod/aws/", "stage/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Remount(namespace.RootContext(nil), "prod/aws/", "stage/aws/")
	if !strings.Contains(err.Error(), "no mount at") {
		t.Fatalf("err: %v", err)
	}

	req := &logical.Request{
		Path: "prod/aws/foo",
	}
	_, err = r.Route(namespace.RootContext(nil), req)
	if !strings.Contains(err.Error(), "unsupported path") {
		t.Fatalf("err: %v", err)
	}

	req = &logical.Request{
		Path: "stage/aws/foo",
	}
	_, err = r.Route(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the path
	if len(n.Paths) != 1 || n.Paths[0] != "foo" {
		t.Fatalf("bad: %v", n.Paths)
	}

	// Check the resolve from storage still works
	_, mount, prefix, _ := r.MatchingAPIPrefixByStoragePath(namespace.RootContext(nil), "logical/foobar")
	if mount != "stage/aws/" {
		t.Fatalf("bad mount: %s", mount)
	}
	if prefix != "logical/" {
		t.Fatalf("Bad prefix: %s", prefix)
	}
}

func TestRouter_RootPath(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	n := &NoopBackend{
		Root: []string{
			"root",
			"policy/*",
		},
	}
	err = r.Mount(n, "prod/aws/", &MountEntry{UUID: meUUID, Accessor: "awsaccessor", NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	type tcase struct {
		path   string
		expect bool
	}
	tcases := []tcase{
		{"random", false},
		{"prod/aws/foo", false},
		{"prod/aws/root", true},
		{"prod/aws/root-more", false},
		{"prod/aws/policy", false},
		{"prod/aws/policy/", true},
		{"prod/aws/policy/ops", true},
	}

	for _, tc := range tcases {
		out := r.RootPath(namespace.RootContext(nil), tc.path)
		if out != tc.expect {
			t.Fatalf("bad: path: %s expect: %v got %v", tc.path, tc.expect, out)
		}
	}
}

func TestRouter_LoginPath(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "auth/")

	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	n := &NoopBackend{
		Login: []string{
			"login",
			"oauth/*",
			"glob1*",
			"+/wildcard/glob2*",
			"end1/+",
			"end2/+/",
			"end3/+/*",
			"middle1/+/bar",
			"middle2/+/+/bar",
			"+/begin",
			"+/around/+/",
		},
	}
	err = r.Mount(n, "auth/foo/", &MountEntry{UUID: meUUID, Accessor: "authfooaccessor", NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	type tcase struct {
		path   string
		expect bool
	}
	tcases := []tcase{
		{"random", false},
		{"auth/foo/bar", false},
		{"auth/foo/login", true},
		{"auth/foo/login/", false},
		{"auth/invalid/login", false},
		{"auth/foo/oauth", false},
		{"auth/foo/oauth/", true},
		{"auth/foo/oauth/redirect", true},
		{"auth/foo/oauth/redirect/", true},
		{"auth/foo/oauth/redirect/bar", true},
		{"auth/foo/glob1", true},
		{"auth/foo/glob1/", true},
		{"auth/foo/glob1/redirect", true},

		// Wildcard cases

		// "+/wildcard/glob2*"
		{"auth/foo/bar/wildcard/glo", false},
		{"auth/foo/bar/wildcard/glob2", true},
		{"auth/foo/bar/wildcard/glob2222", true},
		{"auth/foo/bar/wildcard/glob2/", true},
		{"auth/foo/bar/wildcard/glob2/baz", true},

		// "end1/+"
		{"auth/foo/end1", false},
		{"auth/foo/end1/", true},
		{"auth/foo/end1/bar", true},
		{"auth/foo/end1/bar/", false},
		{"auth/foo/end1/bar/baz", false},
		// "end2/+/"
		{"auth/foo/end2", false},
		{"auth/foo/end2/", false},
		{"auth/foo/end2/bar", false},
		{"auth/foo/end2/bar/", true},
		{"auth/foo/end2/bar/baz", false},
		// "end3/+/*"
		{"auth/foo/end3", false},
		{"auth/foo/end3/", false},
		{"auth/foo/end3/bar", false},
		{"auth/foo/end3/bar/", true},
		{"auth/foo/end3/bar/baz", true},
		{"auth/foo/end3/bar/baz/", true},
		{"auth/foo/end3/bar/baz/qux", true},
		{"auth/foo/end3/bar/baz/qux/qoo", true},
		{"auth/foo/end3/bar/baz/qux/qoo/qaa", true},
		// "middle1/+/bar",
		{"auth/foo/middle1/bar", false},
		{"auth/foo/middle1/bar/", false},
		{"auth/foo/middle1/bar/qux", false},
		{"auth/foo/middle1/bar/bar", true},
		{"auth/foo/middle1/bar/bar/", false},
		// "middle2/+/+/bar",
		{"auth/foo/middle2/bar", false},
		{"auth/foo/middle2/bar/", false},
		{"auth/foo/middle2/bar/baz", false},
		{"auth/foo/middle2/bar/baz/", false},
		{"auth/foo/middle2/bar/baz/bar", true},
		{"auth/foo/middle2/bar/baz/bar/", false},
		// "+/begin"
		{"auth/foo/bar/begin", true},
		{"auth/foo/bar/begin/", false},
		{"auth/foo/begin", false},
		// "+/around/+/"
		{"auth/foo/bar/around", false},
		{"auth/foo/bar/around/", false},
		{"auth/foo/bar/around/baz", false},
		{"auth/foo/bar/around/baz/", true},
		{"auth/foo/bar/around/baz/qux", false},
	}

	for _, tc := range tcases {
		out := r.LoginPath(namespace.RootContext(nil), tc.path)
		if out != tc.expect {
			t.Fatalf("bad: path: %s expect: %v got %v", tc.path, tc.expect, out)
		}
	}
}

func TestRouter_Taint(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	n := &NoopBackend{}
	err = r.Mount(n, "prod/aws/", &MountEntry{UUID: meUUID, Accessor: "awsaccessor", NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Taint(namespace.RootContext(nil), "prod/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	_, err = r.Route(namespace.RootContext(nil), req)
	if err.Error() != "unsupported path" {
		t.Fatalf("err: %v", err)
	}

	// Rollback and Revoke should work
	req.Operation = logical.RollbackOperation
	_, err = r.Route(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req.Operation = logical.RevokeOperation
	_, err = r.Route(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestRouter_Untaint(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	n := &NoopBackend{}
	err = r.Mount(n, "prod/aws/", &MountEntry{UUID: meUUID, Accessor: "awsaccessor", NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Taint(namespace.RootContext(nil), "prod/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Untaint(namespace.RootContext(nil), "prod/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	_, err = r.Route(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestPathsToRadix(t *testing.T) {
	// Provide real paths
	paths := []string{
		"foo",
		"foo/*",
		"sub/bar*",
	}
	r := pathsToRadix(paths)

	raw, ok := r.Get("foo")
	if !ok || raw.(bool) != false {
		t.Fatalf("bad: %v (foo)", raw)
	}

	raw, ok = r.Get("foo/")
	if !ok || raw.(bool) != true {
		t.Fatalf("bad: %v (foo/)", raw)
	}

	raw, ok = r.Get("sub/bar")
	if !ok || raw.(bool) != true {
		t.Fatalf("bad: %v (sub/bar)", raw)
	}
}

func TestParseUnauthenticatedPaths(t *testing.T) {
	// inputs
	paths := []string{
		"foo",
		"foo/*",
		"sub/bar*",
	}
	wildcardPaths := []string{
		"end/+",
		"+/begin/*",
		"middle/+/bar*",
	}
	allPaths := append(paths, wildcardPaths...)

	p, err := parseUnauthenticatedPaths(allPaths)
	if err != nil {
		t.Fatal(err)
	}

	// outputs
	wildcardPathsEntry := []wildcardPath{
		{segments: []string{"end", "+"}, isPrefix: false},
		{segments: []string{"+", "begin", ""}, isPrefix: true},
		{segments: []string{"middle", "+", "bar"}, isPrefix: true},
	}
	expected := &specialPathsEntry{
		paths:         pathsToRadix(paths),
		wildcardPaths: wildcardPathsEntry,
	}

	if !reflect.DeepEqual(expected, p) {
		t.Fatalf("expected: %#v\n actual: %#v\n", expected, p)
	}
}

func TestParseUnauthenticatedPaths_Error(t *testing.T) {
	type tcase struct {
		paths []string
		err   string
	}
	tcases := []tcase{
		{
			[]string{"/foo/+*"},
			"path \"/foo/+*\": invalid use of wildcards ('+*' is forbidden)",
		},
		{
			[]string{"/foo/*/*"},
			"path \"/foo/*/*\": invalid use of wildcards (multiple '*' is forbidden)",
		},
		{
			[]string{"*/foo/*"},
			"path \"*/foo/*\": invalid use of wildcards (multiple '*' is forbidden)",
		},
		{
			[]string{"*/foo/"},
			"path \"*/foo/\": invalid use of wildcards ('*' is only allowed at the end of a path)",
		},
		{
			[]string{"/foo+"},
			"path \"/foo+\": invalid use of wildcards ('+' is not allowed next to a non-slash)",
		},
		{
			[]string{"/+foo"},
			"path \"/+foo\": invalid use of wildcards ('+' is not allowed next to a non-slash)",
		},
		{
			[]string{"/++"},
			"path \"/++\": invalid use of wildcards ('+' is not allowed next to a non-slash)",
		},
	}

	for _, tc := range tcases {
		_, err := parseUnauthenticatedPaths(tc.paths)
		if err == nil || err != nil && !strings.Contains(err.Error(), tc.err) {
			t.Fatalf("bad: path: %s expect: %v got %v", tc.paths, tc.err, err)
		}
	}
}

func TestWellKnownRedirectMatching(t *testing.T) {
	a := assert.New(t)
	// inputs
	redirs := map[string]string{
		"foo":     "v1/one-path",
		"bar/baz": "v1/two-paths",
		"baz/":    "v1/trailing-slash",
	}

	tests := map[string]struct {
		expected string
		mismatch bool
	}{
		"foo":           {"/v1/one-path", false},
		"foof":          {"", true},
		"foo/extra":     {"/v1/one-path/extra", false},
		"bar/baz":       {"/v1/two-paths", false},
		"bar/baz/extra": {"/v1/two-paths/extra", false},
		"baz":           {"/v1/trailing-slash", false},
		"baz/extra":     {"/v1/trailing-slash/extra", false},
	}
	apiRedir := NewWellKnownRedirects()
	for s, d := range redirs {
		if err := apiRedir.TryRegister(context.Background(), nil, "my-mount", s, d); err != nil {
			t.Fatal(err)
		}
	}

	for k, x := range tests {
		t.Run(k, func(t *testing.T) {
			v, s := apiRedir.Find(k)
			if x.mismatch && v != nil {
				t.Fail()
			} else if !x.mismatch && v == nil {
				t.Fail()
			} else if !x.mismatch {
				d, err := v.Destination(s)
				if err != nil {
					t.Fatal(err)
				}
				a.Equal(x.expected, d)
			}
		})
	}

	if found := apiRedir.DeregisterSource("my-mount", "bar/baz"); !found {
		t.Fail()
	}
}
