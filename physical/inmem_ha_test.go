package physical

import "testing"

func TestInmemHA(t *testing.T) {
	inm := NewInmemHA()
	testHABackend(t, inm, inm)
}
