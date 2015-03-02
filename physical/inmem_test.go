package physical

import "testing"

func TestInmem(t *testing.T) {
	inm := NewInmem()
	testBackend(t, inm)
	testBackend_ListPrefix(t, inm)
}
