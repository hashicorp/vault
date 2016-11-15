package token

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
)

// TestCommand re-uses the existing Test function to ensure proper behavior of
// the internal token helper
func TestCommand(t *testing.T) {
	Test(t, &InternalTokenHelper{})
}

func TestInternalTokenHelper(t *testing.T) {
	testHashedStorage(t, "http://127.0.0.1:8200", "0769a29d")
	testHashedStorage(t, "https://127.0.0.1:8200", "0769a29d")
}

func testHashedStorage(t *testing.T, addr string, extension string) {
	// InternalTokenHelper evaluates the user's homedir, which changes for each test case
	homedir.DisableCache = true

	var tkn = "not a valid token"

	td, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	os.Setenv("HOME", td)
	homePath, err := homedir.Dir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	h := &InternalTokenHelper{}
	h.SetVaultAddress(addr)
	h.Store(tkn)

	_, err = os.Stat(homePath + "/.vault-token-" + extension)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual, err := h.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if actual != tkn {
		t.Fatalf("bad: expected %s, received %s", tkn, actual)
	}
}
