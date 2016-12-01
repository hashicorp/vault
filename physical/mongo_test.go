package physical

import (
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

var mongoConfBackend map[string]string = map[string]string{
// nothing needed
}

var mongoConfHABackend map[string]string = map[string]string{
	"ha_enabled": "1",
	"collection": "vault_ha",
}

func TestMongoBackend(t *testing.T) {

	logger := logformat.NewVaultLogger(log.LevelTrace)
	b, err := NewBackend("mongo", logger, mongoConfBackend)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

}

func TestMongoHABackend(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)
	b, err := NewBackend("mongo", logger, mongoConfHABackend)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	ha, ok := b.(HABackend)
	if !ok {
		t.Fatalf("mongo does not implement HABackend")
	}
	testHABackend(t, ha, ha)

}
