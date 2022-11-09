package logging

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogFile_openNew(t *testing.T) {
	logFile := LogFile{
		fileName: "vault-agent.log",
		logPath:  t.TempDir(),
	}
	err := logFile.openNew()
	require.NoError(t, err)

	msg := "[INFO] Something"
	_, err = logFile.Write([]byte(msg))
	require.NoError(t, err)

	content, err := os.ReadFile(logFile.fileInfo.Name())
	require.NoError(t, err)
	require.Contains(t, string(content), msg)
}
