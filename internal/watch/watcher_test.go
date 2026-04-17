package watch_test

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/watch"
)

func TestRun_CallsHandlerOnWrite(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, ".vaultpull.yaml")
	if err := os.WriteFile(file, []byte("initial"), 0o644); err != nil {
		t.Fatal(err)
	}

	var called atomic.Int32
	h := func() error {
		called.Add(1)
		return nil
	}

	w := watch.New(file, 50*time.Millisecond, h)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- w.Run(ctx) }()

	time.Sleep(80 * time.Millisecond)
	if err := os.WriteFile(file, []byte("updated"), 0o644); err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	cancel()
	<-done

	if called.Load() == 0 {
		t.Error("expected handler to be called at least once")
	}
}

func TestRun_CancelStopsWatcher(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, ".vaultpull.yaml")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	w := watch.New(file, 50*time.Millisecond, func() error { return nil })
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err == nil {
		t.Error("expected context cancellation error")
	}
}

func TestRun_MissingFileReturnsError(t *testing.T) {
	w := watch.New("/nonexistent/path.yaml", 50*time.Millisecond, func() error { return nil })
	err := w.Run(context.Background())
	if err == nil {
		t.Error("expected error for missing file")
	}
}
