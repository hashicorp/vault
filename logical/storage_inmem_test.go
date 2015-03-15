package logical

import (
	"testing"
)

func TestInmemStorage(t *testing.T) {
	TestStorage(t, new(InmemStorage))
}
