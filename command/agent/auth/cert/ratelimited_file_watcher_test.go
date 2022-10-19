package cert

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/stretchr/testify/require"
)

func TestNewRateLimitedWatcher(t *testing.T) {
	w, err := NewRateLimitedFileWatcher([]string{}, hclog.New(&hclog.LoggerOptions{}), 1*time.Nanosecond)
	require.NoError(t, err)
	require.NotNil(t, w)
}

func TestRateLimitedWatcherRenameEvent(t *testing.T) {

	fileTmp := createTempConfigFile(t, "temp_config3")
	filepaths := []string{createTempConfigFile(t, "temp_config1"), createTempConfigFile(t, "temp_config2")}
	w, err := NewRateLimitedFileWatcher(filepaths, hclog.New(&hclog.LoggerOptions{}), 1*time.Nanosecond)

	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	require.NoError(t, err)
	err = os.Rename(fileTmp, filepaths[0])
	time.Sleep(timeoutDuration + 50*time.Millisecond)
	require.NoError(t, err)
	require.NoError(t, assertEvent(filepaths[0], w.EventsCh(), defaultTimeout))
	// make sure we consume all events
	_ = assertEvent(filepaths[0], w.EventsCh(), defaultTimeout)
}

func TestRateLimitedWatcherAddNotExist(t *testing.T) {

	file := testutil.TempFile(t, "temp_config")
	filename := file.Name() + randomStr(16)
	w, err := NewRateLimitedFileWatcher([]string{filename}, hclog.New(&hclog.LoggerOptions{}), 1*time.Nanosecond)
	require.Error(t, err, "no such file or directory")
	require.Nil(t, w)
}

func TestEventRateLimitedWatcherWrite(t *testing.T) {

	file := testutil.TempFile(t, "temp_config")
	_, err := file.WriteString("test config")
	require.NoError(t, err)
	err = file.Sync()
	require.NoError(t, err)
	w, err := NewRateLimitedFileWatcher([]string{file.Name()}, hclog.New(&hclog.LoggerOptions{}), 1*time.Nanosecond)
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	_, err = file.WriteString("test config 2")
	require.NoError(t, err)
	err = file.Sync()
	require.NoError(t, err)
	require.NoError(t, assertEvent(file.Name(), w.EventsCh(), defaultTimeout))
}

func TestEventRateLimitedWatcherMove(t *testing.T) {

	filepath := createTempConfigFile(t, "temp_config1")

	w, err := NewRateLimitedFileWatcher([]string{filepath}, hclog.New(&hclog.LoggerOptions{}), 1*time.Second)
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	for i := 0; i < 10; i++ {
		filepath2 := createTempConfigFile(t, "temp_config2")
		err = os.Rename(filepath2, filepath)
		time.Sleep(timeoutDuration + 50*time.Millisecond)
		require.NoError(t, err)
	}
	require.NoError(t, assertEvent(filepath, w.EventsCh(), defaultTimeout))
	require.Error(t, assertEvent(filepath, w.EventsCh(), defaultTimeout), "expected timeout error")
}
