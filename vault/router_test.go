package vault

import (
	"fmt"
	"strings"
	"testing"
)

type NoopBackend struct {
	Root  []string
	Paths []string
}

func (n *NoopBackend) HandleRequest(req *Request) (*Response, error) {
	n.Paths = append(n.Paths, req.Path)
	if req.View == nil {
		return nil, fmt.Errorf("missing view")
	}
	return nil, nil
}

func (n *NoopBackend) RootPaths() []string {
	return n.Root
}

func TestRouter_Mount(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	n := &NoopBackend{}
	err := r.Mount(n, "prod/aws/", view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Mount(n, "prod/aws/", view)
	if !strings.Contains(err.Error(), "cannot mount under existing mount") {
		t.Fatalf("err: %v", err)
	}

	req := &Request{
		Path: "prod/aws/foo",
	}
	resp, err := r.Route(req)
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

	n := &NoopBackend{}
	err := r.Mount(n, "prod/aws/", view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Unmount("prod/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := &Request{
		Path: "prod/aws/foo",
	}
	_, err = r.Route(req)
	if !strings.Contains(err.Error(), "no handler for route") {
		t.Fatalf("err: %v", err)
	}
}

func TestRouter_Remount(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	n := &NoopBackend{}
	err := r.Mount(n, "prod/aws/", view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Remount("prod/aws/", "stage/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Remount("prod/aws/", "stage/aws/")
	if !strings.Contains(err.Error(), "no mount at") {
		t.Fatalf("err: %v", err)
	}

	req := &Request{
		Path: "prod/aws/foo",
	}
	_, err = r.Route(req)
	if !strings.Contains(err.Error(), "no handler for route") {
		t.Fatalf("err: %v", err)
	}

	req = &Request{
		Path: "stage/aws/foo",
	}
	_, err = r.Route(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the path
	if len(n.Paths) != 1 || n.Paths[0] != "foo" {
		t.Fatalf("bad: %v", n.Paths)
	}
}

func TestRouter_RootPath(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	n := &NoopBackend{
		Root: []string{
			"root",
			"policy/*",
		},
	}
	err := r.Mount(n, "prod/aws/", view)
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
		out := r.RootPath(tc.path)
		if out != tc.expect {
			t.Fatalf("bad: path: %s expect: %v got %v", tc.path, tc.expect, out)
		}
	}
}
