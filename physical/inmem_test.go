package physical

import "testing"

func TestInmem(t *testing.T) {
	inm := newInmem()
	testBackend(t, inm)
	testBackend_ListPrefix(t, inm)
}
