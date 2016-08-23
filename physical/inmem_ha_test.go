package physical

import (
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func TestInmemHA(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	inm := NewInmemHA(logger)
	testHABackend(t, inm, inm)
}
