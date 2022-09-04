package cert

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/go-hclog"
)

const timeoutDuration = 200 * time.Millisecond

type Watcher interface {
	Start(ctx context.Context)
	Stop() error
	Add(filename string) error
	Remove(filename string)
	Replace(oldFile, newFile string) error
	EventsCh() chan *FileWatcherEvent
}

type fileWatcher struct {
	watcher          *fsnotify.Watcher
	configFiles      map[string]*watchedFile
	configFilesLock  sync.RWMutex
	logger           hclog.Logger
	reconcileTimeout time.Duration
	cancel           context.CancelFunc
	done             chan interface{}
	stopOnce         sync.Once

	//eventsCh Channel where an event will be emitted when a file change is detected
	// a call to Start is needed before any event is emitted
	// after a Call to Stop succeed, the channel will be closed
	eventsCh chan *FileWatcherEvent
}

type watchedFile struct {
	modTime time.Time
}

type FileWatcherEvent struct {
	Filenames []string
}

//NewFileWatcher create a file watcher that will watch all the files/folders from configFiles
// if success a fileWatcher will be returned and a nil error
// otherwise an error and a nil fileWatcher are returned
func NewFileWatcher(configFiles []string, logger hclog.Logger) (Watcher, error) {
	ws, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &fileWatcher{
		watcher:          ws,
		logger:           logger.Named("file-watcher"),
		configFiles:      make(map[string]*watchedFile),
		eventsCh:         make(chan *FileWatcherEvent),
		reconcileTimeout: timeoutDuration,
		done:             make(chan interface{}),
	}
	for _, f := range configFiles {
		err = w.Add(f)
		if err != nil {
			return nil, fmt.Errorf("error adding file %q: %w", f, err)
		}
	}

	return w, nil
}

// Start start a file watcher, with a copy of the passed context.
// calling Start multiple times is a noop
func (w *fileWatcher) Start(ctx context.Context) {
	if w.cancel == nil {
		cancelCtx, cancel := context.WithCancel(ctx)
		w.cancel = cancel
		go w.watch(cancelCtx)
	}
}

// Stop the file watcher
// calling Stop multiple times is a noop, Stop must be called after a Start
func (w *fileWatcher) Stop() error {
	var err error
	w.stopOnce.Do(func() {
		w.cancel()
		<-w.done
		err = w.watcher.Close()
	})
	return err
}

// Add a file to the file watcher
// Add will lock the file watcher during the add
func (w *fileWatcher) Add(filename string) error {
	filename = filepath.Clean(filename)
	w.logger.Trace("adding file", "file", filename)
	if err := w.watcher.Add(filename); err != nil {
		return err
	}
	modTime, err := w.getFileModifiedTime(filename)
	if err != nil {
		return err
	}
	w.addFile(filename, modTime)
	return nil
}

// Remove a file from the file watcher
// Remove will lock the file watcher during the remove
func (w *fileWatcher) Remove(filename string) {
	w.removeFile(filename)
}

// Replace a file in the file watcher
// Replace will lock the file watcher during the replace
func (w *fileWatcher) Replace(oldFile, newFile string) error {
	if oldFile == newFile {
		return nil
	}
	newFile = filepath.Clean(newFile)
	w.logger.Trace("adding file", "file", newFile)
	if err := w.watcher.Add(newFile); err != nil {
		return err
	}
	modTime, err := w.getFileModifiedTime(newFile)
	if err != nil {
		return err
	}
	w.replaceFile(oldFile, newFile, modTime)
	return nil
}

func (w *fileWatcher) replaceFile(oldFile, newFile string, modTime time.Time) {
	w.configFilesLock.Lock()
	defer w.configFilesLock.Unlock()
	delete(w.configFiles, oldFile)
	w.configFiles[newFile] = &watchedFile{modTime: modTime}
}

func (w *fileWatcher) addFile(filename string, modTime time.Time) {
	w.configFilesLock.Lock()
	defer w.configFilesLock.Unlock()
	w.configFiles[filename] = &watchedFile{modTime: modTime}
}

func (w *fileWatcher) removeFile(filename string) {
	w.configFilesLock.Lock()
	defer w.configFilesLock.Unlock()
	delete(w.configFiles, filename)
}

func (w *fileWatcher) EventsCh() chan *FileWatcherEvent {
	return w.eventsCh
}

func (w *fileWatcher) watch(ctx context.Context) {
	ticker := time.NewTicker(w.reconcileTimeout)
	defer ticker.Stop()
	defer close(w.done)
	defer close(w.eventsCh)

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				w.logger.Error("watcher event channel is closed")
				return
			}
			w.logger.Trace("received watcher event", "event", event)
			if err := w.handleEvent(ctx, event); err != nil {
				w.logger.Error("error handling watcher event", "error", err, "event", event)
			}
		case _, ok := <-w.watcher.Errors:
			if !ok {
				w.logger.Error("watcher error channel is closed")
				return
			}
		case <-ticker.C:
			w.reconcile(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *fileWatcher) handleEvent(ctx context.Context, event fsnotify.Event) error {
	w.logger.Trace("event received ", "filename", event.Name, "OP", event.Op)
	// we only want Create and Remove events to avoid triggering a reload on file modification
	if !isCreateEvent(event) && !isRemoveEvent(event) && !isWriteEvent(event) && !isRenameEvent(event) {
		return nil
	}
	filename := filepath.Clean(event.Name)
	configFile, basename, ok := w.isWatched(filename)
	if !ok {
		return fmt.Errorf("file %s is not watched", event.Name)
	}

	// we only want to update mod time and re-add if the event is on the watched file itself
	if filename == basename {
		if isRemoveEvent(event) {
			// If the file was removed, try to reconcile and see if anything changed.
			w.logger.Trace("attempt a reconcile ", "filename", event.Name, "OP", event.Op)
			configFile.modTime = time.Time{}
			w.reconcile(ctx)
		}
	}
	if isCreateEvent(event) || isWriteEvent(event) || isRenameEvent(event) {
		w.logger.Trace("call the handler", "filename", event.Name, "OP", event.Op)
		select {
		case w.eventsCh <- &FileWatcherEvent{Filenames: []string{filename}}:
		case <-ctx.Done():
			return ctx.Err()
		}

	}
	return nil
}

func (w *fileWatcher) isWatched(filename string) (*watchedFile, string, bool) {
	path := filename
	w.configFilesLock.RLock()
	configFile, ok := w.configFiles[path]
	w.configFilesLock.RUnlock()
	if ok {
		return configFile, path, true
	}

	stat, err := os.Lstat(filename)

	// if the error is a not exist still try to find if the event for a configured file
	if os.IsNotExist(err) || (!stat.IsDir() && stat.Mode()&os.ModeSymlink == 0) {
		w.logger.Trace("not a dir and not a symlink to a dir")
		// try to see if the watched path is the parent dir
		newPath := filepath.Dir(path)
		w.logger.Trace("get dir", "dir", newPath)
		w.configFilesLock.RLock()
		configFile, ok = w.configFiles[newPath]
		w.configFilesLock.RUnlock()
	}
	return configFile, path, ok
}

func (w *fileWatcher) reconcile(ctx context.Context) {
	w.configFilesLock.Lock()
	defer w.configFilesLock.Unlock()
	for filename, configFile := range w.configFiles {
		newModTime, err := w.getFileModifiedTime(filename)
		if err != nil {
			w.logger.Error("failed to get file modTime", "file", filename, "err", err)
			continue
		}

		err = w.watcher.Add(filename)
		if err != nil {
			w.logger.Error("failed to add file to watcher", "file", filename, "err", err)
			continue
		}
		if !configFile.modTime.Equal(newModTime) {
			w.logger.Trace("call the handler", "filename", filename, "old modTime", configFile.modTime, "new modTime", newModTime)
			configFile.modTime = newModTime
			select {
			case w.eventsCh <- &FileWatcherEvent{Filenames: []string{filename}}:
			case <-ctx.Done():
				return
			}
		}
	}
}

func isCreateEvent(event fsnotify.Event) bool {
	return event.Op&fsnotify.Create == fsnotify.Create
}

func isRemoveEvent(event fsnotify.Event) bool {
	return event.Op&fsnotify.Remove == fsnotify.Remove
}

func isWriteEvent(event fsnotify.Event) bool {
	return event.Op&fsnotify.Write == fsnotify.Write
}

func isRenameEvent(event fsnotify.Event) bool {
	return event.Op&fsnotify.Rename == fsnotify.Rename
}

func (w *fileWatcher) getFileModifiedTime(filename string) (time.Time, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}

	return fileInfo.ModTime(), err
}
