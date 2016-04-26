package physical

import (
	"log"
	"os"
	"testing"
)

func TestInmem(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	inm := NewInmem(logger)
	testBackend(t, inm)
	testBackend_ListPrefix(t, inm)
}
