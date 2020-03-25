package dependency

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	// Ensure implements
	_ Dependency = (*FileQuery)(nil)

	// FileQuerySleepTime is the amount of time to sleep between queries, since
	// the fsnotify library is not compatible with solaris and other OSes yet.
	FileQuerySleepTime = 2 * time.Second
)

// FileQuery represents a local file dependency.
type FileQuery struct {
	stopCh chan struct{}

	path string
	stat os.FileInfo
}

// NewFileQuery creates a file dependency from the given path.
func NewFileQuery(s string) (*FileQuery, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("file: invalid format: %q", s)
	}

	return &FileQuery{
		stopCh: make(chan struct{}, 1),
		path:   s,
	}, nil
}

// Fetch retrieves this dependency and returns the result or any errors that
// occur in the process.
func (d *FileQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	log.Printf("[TRACE] %s: READ %s", d, d.path)

	select {
	case <-d.stopCh:
		log.Printf("[TRACE] %s: stopped", d)
		return "", nil, ErrStopped
	case r := <-d.watch(d.stat):
		if r.err != nil {
			return "", nil, errors.Wrap(r.err, d.String())
		}

		log.Printf("[TRACE] %s: reported change", d)

		data, err := ioutil.ReadFile(d.path)
		if err != nil {
			return "", nil, errors.Wrap(err, d.String())
		}

		d.stat = r.stat
		return respWithMetadata(string(data))
	}
}

// CanShare returns a boolean if this dependency is shareable.
func (d *FileQuery) CanShare() bool {
	return false
}

// Stop halts the dependency's fetch function.
func (d *FileQuery) Stop() {
	close(d.stopCh)
}

// String returns the human-friendly version of this dependency.
func (d *FileQuery) String() string {
	return fmt.Sprintf("file(%s)", d.path)
}

// Type returns the type of this dependency.
func (d *FileQuery) Type() Type {
	return TypeLocal
}

type watchResult struct {
	stat os.FileInfo
	err  error
}

// watch watchers the file for changes
func (d *FileQuery) watch(lastStat os.FileInfo) <-chan *watchResult {
	ch := make(chan *watchResult, 1)

	go func(lastStat os.FileInfo) {
		for {
			stat, err := os.Stat(d.path)
			if err != nil {
				select {
				case <-d.stopCh:
					return
				case ch <- &watchResult{err: err}:
					return
				}
			}

			changed := lastStat == nil ||
				lastStat.Size() != stat.Size() ||
				lastStat.ModTime() != stat.ModTime()

			if changed {
				select {
				case <-d.stopCh:
					return
				case ch <- &watchResult{stat: stat}:
					return
				}
			}

			time.Sleep(FileQuerySleepTime)
		}
	}(lastStat)

	return ch
}
