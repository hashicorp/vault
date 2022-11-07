package logging

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogFile_openNew(t *testing.T) {
	logFile := LogFile{
		fileName: "vault-agent.log",
		logPath:  TempTestDir(t, ""),
	}
	err := logFile.openNew()
	require.NoError(t, err)

	msg := "[INFO] Something"
	_, err = logFile.Write([]byte(msg))
	require.NoError(t, err)

	content, err := os.ReadFile(logFile.FileInfo.Name())
	require.NoError(t, err)
	require.Contains(t, string(content), msg)
}

// TempDir creates a temporary directory within tmpdir with the name 'testname-name'.
// If the directory cannot be created t.Fatal is called.
// The directory will be removed when the test ends.
func TempTestDir(t testing.TB, name string) string {
	if t == nil {
		panic("argument t must be non-nil")
	}
	name = t.Name() + "-" + name
	name = strings.Replace(name, "/", "_", -1)
	d, err := os.MkdirTemp("", name)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(d)
	})
	return d
}
