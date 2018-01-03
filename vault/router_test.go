package vault

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/logical"
	log "github.com/mgutz/logxi/v1"
)

type NoopBackend struct {
	sync.Mutex

	Root          []string
	Login         []string
	Paths         []string
	Requests      []*logical.Request
	Response      *logical.Response
	Invalidations []string
}

func (n *NoopBackend) HandleRequest(req *logical.Request) (*logical.Response, error) {
	n.Lock()
	defer n.Unlock()

	requestCopy := *req
	n.Paths = append(n.Paths, req.Path)
	n.Requests = append(n.Requests, &requestCopy)
	if req.Storage == nil {
		return nil, fmt.Errorf("missing view")
	}

	return n.Response, nil
}

func (n *NoopBackend) HandleExistenceCheck(req *logical.Request) (bool, bool, error) {
	return false, false, nil
}

func (n *NoopBackend) SpecialPaths() *logical.Paths {
	return &logical.Paths{
		Root:            n.Root,
		Unauthenticated: n.Login,
	}
}

func (n *NoopBackend) System() logical.SystemView {
	return logical.StaticSystemView{
		DefaultLeaseTTLVal: time.Hour * 24,
		MaxLeaseTTLVal:     time.Hour * 24 * 32,
	}
}

func (n *NoopBackend) Cleanup() {
	// noop
}

func (n *NoopBackend) InvalidateKey(k string) {
	n.Invalidations = append(n.Invalidations, k)
}

func (n *NoopBackend) Setup(config *logical.BackendConfig) error {
	return nil
}

func (n *NoopBackend) Logger() log.Logger {
	return logformat.NewVaultLoggerWithWriter(ioutil.Discard, log.LevelOff)
}

func (n *NoopBackend) Initialize() error {
	return nil
}

func (n *NoopBackend) Type() logical.BackendType {
	return logical.TypeLogical
}

func (n *NoopBackend) RegisterLicense(license interface{}) error {
	return nil
}

func TestRouter_Mount(t *testing.T) {
	r := NewRouter()
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	mountEntry := &MountEntry{
		Path:     "prod/aws/",
		UUID:     meUUID,
		Accessor: "awsaccessor",
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

	err = r.Mount(n, "prod/aws/", &MountEntry{UUID: meUUID}, view)
	if !strings.Contains(err.Error(), "cannot mount under existing mount") {
		t.Fatalf("err: %v", err)
	}

	meUUID, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	if path := r.MatchingMount("prod/aws/foo"); path != "prod/aws/" {
		t.Fatalf("bad: %s", path)
	}

	if v := r.MatchingStorageByAPIPath("prod/aws/foo"); v.(*BarrierView) != view {
		t.Fatalf("bad: %v", v)
	}

	if path := r.MatchingMount("stage/aws/foo"); path != "" {
		t.Fatalf("bad: %s", path)
	}

	if v := r.MatchingStorageByAPIPath("stage/aws/foo"); v != nil {
		t.Fatalf("bad: %v", v)
	}

	mountEntryFetched := r.MatchingMountByUUID(mountEntry.UUID)
	if mountEntryFetched == nil || !reflect.DeepEqual(mountEntry, mountEntryFetched) {
		t.Fatalf("failed to fetch mount entry using its ID; expected: %#v\n actual: %#v\n", mountEntry, mountEntryFetched)
	}

	mount, prefix, ok := r.MatchingStoragePrefixByStoragePath("logical/foo")
	if !ok {
		t.Fatalf("missing storage prefix")
	}
	if mount != "prod/aws/" || prefix != "logical/" {
		t.Fatalf("Bad: %v - %v", mount, prefix)
	}

	req := &logical.Request{
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

	subMountEntry := &MountEntry{
		Path:     "prod/",
		UUID:     meUUID,
		Accessor: "prodaccessor",
	}

	if r.MountConflict("prod/aws/") == "" {
		t.Fatalf("bad: prod/aws/")
	}

	// No error is shown here because MountConflict is checked before Mount
	err = r.Mount(n, "prod/", subMountEntry, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if r.MountConflict("prod/test") == "" {
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
		Path:     "aws",
		UUID:     meUUID,
		Accessor: "awsaccessor",
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

	err = r.Mount(n, "auth/aws/", &MountEntry{UUID: meUUID}, view)
	if !strings.Contains(err.Error(), "cannot mount under existing mount") {
		t.Fatalf("err: %v", err)
	}

	if path := r.MatchingMount("auth/aws/foo"); path != "auth/aws/" {
		t.Fatalf("bad: %s", path)
	}

	if v := r.MatchingStorageByAPIPath("auth/aws/foo"); v.(*BarrierView) != view {
		t.Fatalf("bad: %v", v)
	}

	if path := r.MatchingMount("auth/stage/aws/foo"); path != "" {
		t.Fatalf("bad: %s", path)
	}

	if v := r.MatchingStorageByAPIPath("auth/stage/aws/foo"); v != nil {
		t.Fatalf("bad: %v", v)
	}

	mountEntryFetched := r.MatchingMountByUUID(mountEntry.UUID)
	if mountEntryFetched == nil || !reflect.DeepEqual(mountEntry, mountEntryFetched) {
		t.Fatalf("failed to fetch mount entry using its ID; expected: %#v\n actual: %#v\n", mountEntry, mountEntryFetched)
	}

	mount, prefix, ok := r.MatchingStoragePrefixByStoragePath("auth/foo")
	if !ok {
		t.Fatalf("missing storage prefix")
	}
	if mount != "auth/aws" || prefix != credentialBarrierPrefix {
		t.Fatalf("Bad: %v - %v", mount, prefix)
	}

	req := &logical.Request{
		Path: "auth/aws/foo",
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

	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	n := &NoopBackend{}
	err = r.Mount(n, "prod/aws/", &MountEntry{Path: "prod/aws/", UUID: meUUID, Accessor: "awsaccessor"}, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Unmount("prod/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := &logical.Request{
		Path: "prod/aws/foo",
	}
	_, err = r.Route(req)
	if !strings.Contains(err.Error(), "unsupported path") {
		t.Fatalf("err: %v", err)
	}

	if _, _, ok := r.MatchingStoragePrefixByStoragePath("logical/foo"); ok {
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
	me := &MountEntry{Path: "prod/aws/", UUID: meUUID, Accessor: "awsaccessor"}
	err = r.Mount(n, "prod/aws/", me, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	me.Path = "stage/aws/"
	err = r.Remount("prod/aws/", "stage/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Remount("prod/aws/", "stage/aws/")
	if !strings.Contains(err.Error(), "no mount at") {
		t.Fatalf("err: %v", err)
	}

	req := &logical.Request{
		Path: "prod/aws/foo",
	}
	_, err = r.Route(req)
	if !strings.Contains(err.Error(), "unsupported path") {
		t.Fatalf("err: %v", err)
	}

	req = &logical.Request{
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

	// Check the resolve from storage still works
	mount, prefix, _ := r.MatchingStoragePrefixByStoragePath("logical/foobar")
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
	err = r.Mount(n, "prod/aws/", &MountEntry{UUID: meUUID, Accessor: "awsaccessor"}, view)
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
		},
	}
	err = r.Mount(n, "auth/foo/", &MountEntry{UUID: meUUID, Accessor: "authfooaccessor"}, view)
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
		{"auth/foo/oauth", false},
		{"auth/foo/oauth/redirect", true},
	}

	for _, tc := range tcases {
		out := r.LoginPath(tc.path)
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
	err = r.Mount(n, "prod/aws/", &MountEntry{UUID: meUUID, Accessor: "awsaccessor"}, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Taint("prod/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	_, err = r.Route(req)
	if err.Error() != "unsupported path" {
		t.Fatalf("err: %v", err)
	}

	// Rollback and Revoke should work
	req.Operation = logical.RollbackOperation
	_, err = r.Route(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req.Operation = logical.RevokeOperation
	_, err = r.Route(req)
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
	err = r.Mount(n, "prod/aws/", &MountEntry{UUID: meUUID, Accessor: "awsaccessor"}, view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Taint("prod/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = r.Untaint("prod/aws/")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	_, err = r.Route(req)
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
