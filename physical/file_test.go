package physical

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileBackend(t *testing.T) {
	dir, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	b, err := NewBackend("file", map[string]string{
		"path": dir,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
}
