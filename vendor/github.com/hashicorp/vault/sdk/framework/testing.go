package framework

import (
	"testing"
)

// TestBackendRoutes is a helper to test that all the given routes will
// route properly in the backend.
func TestBackendRoutes(t *testing.T, b *Backend, rs []string) {
	for _, r := range rs {
		if b.Route(r) == nil {
			t.Fatalf("bad route: %s", r)
		}
	}
}
