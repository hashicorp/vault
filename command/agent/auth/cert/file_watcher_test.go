package cert

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/stretchr/testify/require"
)

const defaultTimeout = 500 * time.Millisecond

func TestNewWatcher(t *testing.T) {
	w, err := NewFileWatcher([]string{}, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	require.NotNil(t, w)
}

func TestWatcherRenameEvent(t *testing.T) {

	fileTmp := createTempConfigFile(t, "temp_config3")
	filepaths := []string{createTempConfigFile(t, "temp_config1"), createTempConfigFile(t, "temp_config2")}
	wi, err := NewFileWatcher(filepaths, hclog.New(&hclog.LoggerOptions{}))
	w := wi.(*fileWatcher)

	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	require.NoError(t, err)
	err = os.Rename(fileTmp, filepaths[0])
	time.Sleep(w.reconcileTimeout + 50*time.Millisecond)
	require.NoError(t, err)
	require.NoError(t, assertEvent(filepaths[0], w.eventsCh, defaultTimeout))
	// make sure we consume all events
	_ = assertEvent(filepaths[0], w.eventsCh, defaultTimeout)
}

func TestWatcherAddRemove(t *testing.T) {
	var filepaths []string
	wi, err := NewFileWatcher(filepaths, hclog.New(&hclog.LoggerOptions{}))
	w := wi.(*fileWatcher)
	require.NoError(t, err)
	file1 := createTempConfigFile(t, "temp_config1")
	err = w.Add(file1)
	require.NoError(t, err)
	file2 := createTempConfigFile(t, "temp_config2")
	err = w.Add(file2)
	require.NoError(t, err)
	w.Remove(file2)
	_, ok := w.configFiles[file1]
	require.True(t, ok)
	_, ok = w.configFiles[file2]
	require.False(t, ok)

}

func TestWatcherReplace(t *testing.T) {
	var filepaths []string
	wi, err := NewFileWatcher(filepaths, hclog.New(&hclog.LoggerOptions{}))
	w := wi.(*fileWatcher)
	require.NoError(t, err)
	file1 := createTempConfigFile(t, "temp_config1")
	err = w.Add(file1)
	require.NoError(t, err)
	file2 := createTempConfigFile(t, "temp_config2")
	err = w.Replace(file1, file2)
	require.NoError(t, err)
	_, ok := w.configFiles[file1]
	require.False(t, ok)
	_, ok = w.configFiles[file2]
	require.True(t, ok)
}

func TestWatcherAddWhileRunning(t *testing.T) {
	var filepaths []string
	wi, err := NewFileWatcher(filepaths, hclog.New(&hclog.LoggerOptions{}))
	w := wi.(*fileWatcher)
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()
	file1 := createTempConfigFile(t, "temp_config1")
	err = w.Add(file1)
	require.NoError(t, err)
	file2 := createTempConfigFile(t, "temp_config2")
	err = w.Add(file2)
	require.NoError(t, err)
	w.Remove(file2)
	require.Len(t, w.configFiles, 1)
	_, ok := w.configFiles[file1]
	require.True(t, ok)
	_, ok = w.configFiles[file2]
	require.False(t, ok)
}

func TestWatcherRemoveNotFound(t *testing.T) {
	var filepaths []string
	w, err := NewFileWatcher(filepaths, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	file := createTempConfigFile(t, "temp_config2")
	w.Remove(file)
}

func TestWatcherAddNotExist(t *testing.T) {

	file := testutil.TempFile(t, "temp_config")
	filename := file.Name() + randomStr(16)
	w, err := NewFileWatcher([]string{filename}, hclog.New(&hclog.LoggerOptions{}))
	require.Error(t, err, "no such file or directory")
	require.Nil(t, w)
}

func TestEventWatcherWrite(t *testing.T) {

	file := testutil.TempFile(t, "temp_config")
	_, err := file.WriteString("test config")
	require.NoError(t, err)
	err = file.Sync()
	require.NoError(t, err)
	w, err := NewFileWatcher([]string{file.Name()}, hclog.New(&hclog.LoggerOptions{}))
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

func TestEventWatcherRead(t *testing.T) {

	filepath := createTempConfigFile(t, "temp_config1")
	w, err := NewFileWatcher([]string{filepath}, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	_, err = os.ReadFile(filepath)
	require.NoError(t, err)
	require.Error(t, assertEvent(filepath, w.EventsCh(), defaultTimeout), "timedout waiting for event")
}

func TestEventWatcherChmod(t *testing.T) {
	file := testutil.TempFile(t, "temp_config")
	defer func() {
		err := file.Close()
		require.NoError(t, err)
	}()
	_, err := file.WriteString("test config")
	require.NoError(t, err)
	err = file.Sync()
	require.NoError(t, err)

	w, err := NewFileWatcher([]string{file.Name()}, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	err = file.Chmod(0777)
	require.NoError(t, err)
	require.Error(t, assertEvent(file.Name(), w.EventsCh(), defaultTimeout), "timedout waiting for event")
}

func TestEventWatcherRemoveCreate(t *testing.T) {

	filepath := createTempConfigFile(t, "temp_config1")
	w, err := NewFileWatcher([]string{filepath}, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	require.NoError(t, err)
	err = os.Remove(filepath)
	require.NoError(t, err)
	recreated, err := os.Create(filepath)
	require.NoError(t, err)
	_, err = recreated.WriteString("config 2")
	require.NoError(t, err)
	err = recreated.Sync()
	require.NoError(t, err)
	// this an event coming from the reconcile loop
	require.NoError(t, assertEvent(filepath, w.EventsCh(), defaultTimeout))
}

func TestEventWatcherMove(t *testing.T) {

	filepath := createTempConfigFile(t, "temp_config1")

	w, err := NewFileWatcher([]string{filepath}, hclog.New(&hclog.LoggerOptions{}))
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
		require.NoError(t, assertEvent(filepath, w.EventsCh(), defaultTimeout))
	}
}

func TestEventReconcileMove(t *testing.T) {
	filepath := createTempConfigFile(t, "temp_config1")
	filepath2 := createTempConfigFile(t, "temp_config2")
	err := os.Chtimes(filepath, time.Now(), time.Now().Add(-1*time.Second))
	require.NoError(t, err)
	wi, err := NewFileWatcher([]string{filepath}, hclog.New(&hclog.LoggerOptions{}))
	w := wi.(*fileWatcher)
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	// remove the file from the internal watcher to only trigger the reconcile
	err = w.watcher.Remove(filepath)
	require.NoError(t, err)

	err = os.Rename(filepath2, filepath)
	time.Sleep(timeoutDuration + 50*time.Millisecond)
	require.NoError(t, err)
	require.NoError(t, assertEvent(filepath, w.EventsCh(), 2000*time.Millisecond))
}

func TestEventWatcherDirCreateRemove(t *testing.T) {
	filepath := testutil.TempDir(t, "temp_config1")
	w, err := NewFileWatcher([]string{filepath}, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()
	for i := 0; i < 1; i++ {
		name := filepath + "/" + randomStr(20)
		file, err := os.Create(name)
		require.NoError(t, err)
		err = file.Close()
		require.NoError(t, err)
		require.NoError(t, assertEvent(filepath, w.EventsCh(), defaultTimeout))

		err = os.Remove(name)
		require.NoError(t, err)
		require.NoError(t, assertEvent(filepath, w.EventsCh(), defaultTimeout))
	}
}

func TestEventWatcherDirMove(t *testing.T) {
	filepath := testutil.TempDir(t, "temp_config1")

	name := filepath + "/" + randomStr(20)
	file, err := os.Create(name)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)
	w, err := NewFileWatcher([]string{filepath}, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	for i := 0; i < 100; i++ {
		filepathTmp := createTempConfigFile(t, "temp_config2")
		err = os.Rename(filepathTmp, name)
		require.NoError(t, err)
		require.NoError(t, assertEvent(filepath, w.EventsCh(), defaultTimeout))
	}
}

func TestEventWatcherDirMoveTrim(t *testing.T) {
	filepath := testutil.TempDir(t, "temp_config1")

	name := filepath + "/" + randomStr(20)
	file, err := os.Create(name)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)
	w, err := NewFileWatcher([]string{filepath + "/"}, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	for i := 0; i < 100; i++ {
		filepathTmp := createTempConfigFile(t, "temp_config2")
		err = os.Rename(filepathTmp, name)
		require.NoError(t, err)
		require.NoError(t, assertEvent(filepath, w.EventsCh(), defaultTimeout))
	}
}

// Consul do not support configuration in sub-directories
func TestEventWatcherSubDirMove(t *testing.T) {
	filepath := testutil.TempDir(t, "temp_config1")
	err := os.Mkdir(filepath+"/temp", 0777)
	require.NoError(t, err)
	name := filepath + "/temp/" + randomStr(20)
	file, err := os.Create(name)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)
	w, err := NewFileWatcher([]string{filepath}, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	w.Start(context.Background())
	defer func() {
		_ = w.Stop()
	}()

	for i := 0; i < 2; i++ {
		filepathTmp := createTempConfigFile(t, "temp_config2")
		err = os.Rename(filepathTmp, name)
		require.NoError(t, err)
		require.Error(t, assertEvent(filepath, w.EventsCh(), defaultTimeout), "timedout waiting for event")
	}
}

func TestEventWatcherDirRead(t *testing.T) {
	filepath := testutil.TempDir(t, "temp_config1")

	name := filepath + "/" + randomStr(20)
	file, err := os.Create(name)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)
	w, err := NewFileWatcher([]string{filepath}, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	w.Start(context.Background())
	t.Cleanup(func() {
		_ = w.Stop()
	})

	_, err = os.ReadFile(name)
	require.NoError(t, err)
	require.Error(t, assertEvent(filepath, w.EventsCh(), defaultTimeout), "timedout waiting for event")
}

func TestEventWatcherMoveSoftLink(t *testing.T) {

	filepath := createTempConfigFile(t, "temp_config1")
	tempDir := testutil.TempDir(t, "temp_dir")
	name := tempDir + "/" + randomStr(20)
	err := os.Symlink(filepath, name)
	require.NoError(t, err)

	w, err := NewFileWatcher([]string{name}, hclog.New(&hclog.LoggerOptions{}))
	require.NoError(t, err)
	require.NotNil(t, w)

}

func assertEvent(name string, watcherCh chan *FileWatcherEvent, timeout time.Duration) error {
	select {
	case ev := <-watcherCh:
		if ev.Filenames[0] != name && !strings.Contains(ev.Filenames[0], name) {
			return fmt.Errorf("filename do not match %s %s", ev.Filenames[0], name)
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timedout waiting for event")
	}
}

func createTempConfigFile(t *testing.T, filename string) string {
	file := testutil.TempFile(t, filename)

	_, err1 := file.WriteString("test config")
	err2 := file.Close()

	require.NoError(t, err1)
	require.NoError(t, err2)

	return file.Name()
}

func randomStr(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
