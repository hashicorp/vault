package physical

import (
	"log"
	"os"
	"testing"
)

func TestInmemHA(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	inm := NewInmemHA(logger)
	testHABackend(t, inm, inm)
}
