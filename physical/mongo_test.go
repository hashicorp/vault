package physical

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func TestMongoBackend(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		t.SkipNow()
	}

	var mongoConfBackend map[string]string = map[string]string{
		"url": mongoUri,
	}

	b, err := NewBackend("mongo", logger, mongoConfBackend)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

}

func TestMongoHABackend(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	mongoUri := os.Getenv("MONGODB_URI")
	if mongoUri == "" {
		t.SkipNow()
	}

	var mongoConfHABackend map[string]string = map[string]string{
		"ha_enabled": "1",
		"collection": "vault_ha",
		"url":        mongoUri,
	}

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
