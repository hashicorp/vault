package client

import (
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestClient(t *testing.T) {
	currentPatches, closeFunc := TestServer(t)
	defer closeFunc()

	client, err := New(hclog.Default())
	if err != nil {
		t.Fatal(err)
	}
	e := &env{
		client:         client,
		currentPatches: currentPatches,
	}
	e.TestGetPod(t)
	e.TestGetPodNotFound(t)
	e.TestUpdatePodTags(t)
	e.TestUpdatePodTagsNotFound(t)
}

type env struct {
	client         *Client
	currentPatches map[string]*Patch
}

func (e *env) TestGetPod(t *testing.T) {
	pod, err := e.client.GetPod(TestNamespace, TestPodname)
	if err != nil {
		t.Fatal(err)
	}
	if pod.Metadata.Name != "shell-demo" {
		t.Fatalf("expected %q but received %q", "shell-demo", pod.Metadata.Name)
	}
}

func (e *env) TestGetPodNotFound(t *testing.T) {
	_, err := e.client.GetPod(TestNamespace, "no-exist")
	if err == nil {
		t.Fatal("expected error because pod is unfound")
	}
	if err != ErrNotFound {
		t.Fatalf("expected %q but received %q", ErrNotFound, err)
	}
}

func (e *env) TestUpdatePodTags(t *testing.T) {
	if err := e.client.PatchPod(TestNamespace, TestPodname, &Patch{
		Operation: Add,
		Path:      "/metadata/labels/fizz",
		Value:     "buzz",
	}); err != nil {
		t.Fatal(err)
	}
	if len(e.currentPatches) != 1 {
		t.Fatalf("expected 1 label but received %+v", e.currentPatches)
	}
	if e.currentPatches["/metadata/labels/fizz"].Value != "buzz" {
		t.Fatalf("expected buzz but received %q", e.currentPatches["fizz"])
	}
}

func (e *env) TestUpdatePodTagsNotFound(t *testing.T) {
	err := e.client.PatchPod(TestNamespace, "no-exist", &Patch{
		Operation: Add,
		Path:      "/metadata/labels/fizz",
		Value:     "buzz",
	})
	if err == nil {
		t.Fatal("expected error because pod is unfound")
	}
	if err != ErrNotFound {
		t.Fatalf("expected %q but received %q", ErrNotFound, err)
	}
}
