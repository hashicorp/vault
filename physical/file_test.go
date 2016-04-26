package physical

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestFileBackend(t *testing.T) {
	dir, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	logger := log.New(os.Stderr, "", log.LstdFlags)
	b, err := NewBackend("file", logger, map[string]string{
		"path": dir,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
}
