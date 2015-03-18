package vault

import (
	"log"
	"os"
	"testing"
	"time"
)

// mockRollback returns a mock rollback manager
func mockRollback(t *testing.T) (*RollbackManager, *NoopBackend) {
	backend := new(NoopBackend)
	mounts := new(MountTable)
	router := NewRouter()

	mounts.Entries = []*MountEntry{
		&MountEntry{
			Path: "foo",
		},
	}
	if err := router.Mount(backend, "foo", nil); err != nil {
		t.Fatalf("err: %s", err)
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)
	return &RollbackManager{
		Logger: logger,
		Mounts: mounts,
		Router: router,
		Period: 10 * time.Millisecond,
	}, backend
}

func TestRollbackManager(t *testing.T) {
	m, backend := mockRollback(t)
	if len(backend.Paths) > 0 {
		t.Fatalf("bad: %#v", backend)
	}

	go m.Start()
	time.Sleep(100 * time.Millisecond)
	m.Stop()

	count := len(backend.Paths)
	if count == 0 {
		t.Fatalf("bad: %#v", backend)
	}
	if backend.Paths[0] != "" {
		t.Fatalf("bad: %#v", backend)
	}

	time.Sleep(100 * time.Millisecond)

	if count != len(backend.Paths) {
		t.Fatalf("should stop requests: %#v", backend)
	}
}
