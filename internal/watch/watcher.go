// Package watch provides file-system watching to trigger re-sync when
// the vaultpull config file changes.
package watch

import (
	"context"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Handler is called whenever the watched file changes.
type Handler func() error

// Watcher watches a single file and calls Handler on write events.
type Watcher struct {
	path    string
	debounce time.Duration
	handler Handler
}

// New creates a Watcher for the given file path.
func New(path string, debounce time.Duration, h Handler) *Watcher {
	return &Watcher{path: path, debounce: debounce, handler: h}
}

// Run starts the watch loop and blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer fw.Close()

	if err := fw.Add(w.path); err != nil {
		return err
	}

	var timer *time.Timer

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-fw.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(w.debounce, func() {
					if err := w.handler(); err != nil {
						log.Printf("watch: handler error: %v", err)
					}
				})
			}
		case err, ok := <-fw.Errors:
			if !ok {
				return nil
			}
			log.Printf("watch: fsnotify error: %v", err)
		}
	}
}
