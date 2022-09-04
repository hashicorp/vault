package cert

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
)

type rateLimitedFileWatcher struct {
	watcher          Watcher
	eventCh          chan *FileWatcherEvent
	coalesceInterval time.Duration
}

func (r *rateLimitedFileWatcher) Start(ctx context.Context) {
	r.watcher.Start(ctx)
	r.coalesceTimer(ctx, r.watcher.EventsCh(), r.coalesceInterval)
}

func (r rateLimitedFileWatcher) Stop() error {
	return r.watcher.Stop()
}

func (r rateLimitedFileWatcher) Add(filename string) error {
	return r.watcher.Add(filename)
}

func (r rateLimitedFileWatcher) Remove(filename string) {
	r.watcher.Remove(filename)
}

func (r rateLimitedFileWatcher) Replace(oldFile, newFile string) error {
	return r.watcher.Replace(oldFile, newFile)
}

func (r rateLimitedFileWatcher) EventsCh() chan *FileWatcherEvent {
	return r.eventCh
}

func NewRateLimitedFileWatcher(configFiles []string, logger hclog.Logger, coalesceInterval time.Duration) (Watcher, error) {

	watcher, err := NewFileWatcher(configFiles, logger)
	if err != nil {
		return nil, err
	}
	return &rateLimitedFileWatcher{
		watcher:          watcher,
		coalesceInterval: coalesceInterval,
		eventCh:          make(chan *FileWatcherEvent),
	}, nil
}

func (r rateLimitedFileWatcher) coalesceTimer(ctx context.Context, inputCh chan *FileWatcherEvent, coalesceDuration time.Duration) {
	var (
		coalesceTimer     *time.Timer
		sendCh            = make(chan struct{})
		fileWatcherEvents []string
	)

	go func() {
		for {
			select {
			case event, ok := <-inputCh:
				if !ok {
					if len(fileWatcherEvents) > 0 {
						r.eventCh <- &FileWatcherEvent{Filenames: fileWatcherEvents}
					}
					close(r.eventCh)
					return
				}
				fileWatcherEvents = append(fileWatcherEvents, event.Filenames...)
				if coalesceTimer == nil {
					coalesceTimer = time.AfterFunc(coalesceDuration, func() {
						// This runs in another goroutine so we can't just do the send
						// directly here as access to fileWatcherEvents is racy. Instead,
						// signal the main loop above.
						sendCh <- struct{}{}
					})
				}
			case <-sendCh:
				coalesceTimer = nil
				r.eventCh <- &FileWatcherEvent{Filenames: fileWatcherEvents}
				fileWatcherEvents = make([]string, 0)
			case <-ctx.Done():
				return
			}
		}
	}()
}
