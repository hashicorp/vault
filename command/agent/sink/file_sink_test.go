package sink

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/go-uuid"
)

const (
	fileServerTestDir = "vsi-file-test"
)

func testFileServer(t *testing.T, format string) (Server, string) {
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("%s.", fileServerTestDir))
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf("%s/token", tmpDir)

	config := map[string]string{
		"name": "testserver",
		"path": path,
	}

	if format != "" {
		config["format"] = format
	}

	fs, err := NewFileServer(config)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("[INFO] FileServer created with path %s", path)
	return fs, tmpDir
}

func TestFileServer(t *testing.T) {
	t.Log("[INFO] Starting TestFileServer")
	defer t.Log("[INFO] Finished TestFileServer")

	_ = TestCore(t)

	fs, tmpDir := testFileServer(t, "")
	defer os.RemoveAll(tmpDir)

	path := fmt.Sprintf("%s/token", tmpDir)

	uuidStr, _ := uuid.GenerateUUID()
	if err := fs.WriteToken(uuidStr); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode() != os.FileMode(0640) {
		t.Fatalf("wrong file mode was detected at %s", path)
	}
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(fileBytes) != uuidStr {
		t.Fatalf("expected %s, got %s", uuidStr, string(fileBytes))
	}
}

func TestFileServer_formatToken(t *testing.T) {
	t.Log("[INFO] Starting TestFileServer_formatToken")
	defer t.Log("[INFO] Finished TestFileServer_formatToken")

	_ = TestCore(t)

	fs, tmpDir := testFileServer(t, "token")
	defer os.RemoveAll(tmpDir)

	path := fmt.Sprintf("%s/token", tmpDir)

	uuidStr, _ := uuid.GenerateUUID()
	if err := fs.WriteToken(uuidStr); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode() != os.FileMode(0640) {
		t.Fatalf("wrong file mode was detected at %s", path)
	}
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(fileBytes) != uuidStr {
		t.Fatalf("expected %s, got %s", uuidStr, string(fileBytes))
	}
}

func TestFileServer_formatEnvironment(t *testing.T) {
	t.Log("[INFO] Starting TestFileServer_formatEnvironment")
	defer t.Log("[INFO] Finished TestFileServer_formatEnvironment")

	_ = TestCore(t)

	fs, tmpDir := testFileServer(t, "")
	defer os.RemoveAll(tmpDir)

	path := fmt.Sprintf("%s/token", tmpDir)

	uuidStr, _ := uuid.GenerateUUID()
	if err := fs.WriteToken(uuidStr); err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode() != os.FileMode(0640) {
		t.Fatalf("wrong file mode was detected at %s", path)
	}
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(fileBytes) != fmt.Sprintf("VAULT_TOKEN=%s", uuidStr) {
		t.Fatalf("expected %s, got %s", uuidStr, string(fileBytes))
	}

}
