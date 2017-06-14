package plugin

import (
	"testing"

	log "github.com/mgutz/logxi/v1"
)

func TestLogger_impl(t *testing.T) {
	var _ log.Logger = new(LoggerClient)
}
