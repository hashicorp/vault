package benchhelpers

import (
	"testing"

	testing2 "github.com/mitchellh/go-testing-interface"
)

type tbWrapper struct {
	testing.TB
}

func (b tbWrapper) Parallel() {
	// no-op
}

func TBtoT(tb testing.TB) testing2.T {
	return tbWrapper{tb}
}
