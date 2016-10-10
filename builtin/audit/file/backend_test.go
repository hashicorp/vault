package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/salt"
)

func TestAuditFile_fileModeNew(t *testing.T) {
	salter, _ := salt.NewSalt(nil, nil)

	modeStr := "0777"
	mode, err := strconv.ParseUint(modeStr, 8, 32)

	path, err := ioutil.TempDir("", "test")
	defer os.RemoveAll(path)

	file := filepath.Join(path, "auditTest.txt")

	config := map[string]string{
		"path": file,
		"mode": modeStr,
	}

	_, err = Factory(&audit.BackendConfig{
		Salt:   salter,
		Config: config,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file)

	info, _ := os.Stat(file)
	createdMode := info.Mode()
	if createdMode != os.FileMode(mode) {
		t.Fatalf("File mode does not match.")
	}
}

func TestAuditFile_fileModeExisting(t *testing.T) {
	salter, _ := salt.NewSalt(nil, nil)

	f, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("Failure to create test file.")
	}
	defer os.Remove(f.Name())
	err = f.Close()
	if err != nil {
		t.Fatalf("Failure to close the file.")
	}

	config := map[string]string{
		"path": f.Name(),
	}

	_, err = Factory(&audit.BackendConfig{
		Salt:   salter,
		Config: config,
	})
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(f.Name())
	if err != nil {
		t.Fatalf("cannot retrieve file mode from `Stat`")
	}
	createdMode := info.Mode()
	if createdMode != os.FileMode(0600) {
		t.Fatalf("File mode does not match.")
	}
}
