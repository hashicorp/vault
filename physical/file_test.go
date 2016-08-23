package physical

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func TestFileBackend(t *testing.T) {
	dir, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("file", logger, map[string]string{
		"path": dir,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
}
