package template

import "testing"

// TestNewServer is a simple test to make sure NewServer returns a Server and
// channel
func TestNewServer(t *testing.T) {
	ts, ch := NewServer(&ServerConfig{})
	if ts == nil {
		t.Fatal("nil server returned")
	}
	if ch == nil {
		t.Fatal("nil blocking channel returned")
	}
}
