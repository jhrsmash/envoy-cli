package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnvFile: %v", err)
	}
	return p
}

func TestWatch_DetectsChange(t *testing.T) {
	path := writeTempEnvFile(t, "FOO=bar\n")

	done := make(chan struct{})
	defer close(done)

	events, errs := Watch(path, WatchOptions{Interval: 50 * time.Millisecond, MaxChecks: 10}, done)

	// Give the watcher one tick to settle, then update the file.
	time.Sleep(80 * time.Millisecond)
	if err := os.WriteFile(path, []byte("FOO=bar\nBAZ=qux\n"), 0o644); err != nil {
		t.Fatalf("update file: %v", err)
	}

	select {
	case ev := <-events:
		if ev.File != path {
			t.Errorf("expected file %s, got %s", path, ev.File)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
		if len(ev.Diff) == 0 {
			t.Error("expected non-empty diff")
		}
	case err := <-errs:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for watch event")
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	path := writeTempEnvFile(t, "FOO=bar\n")

	done := make(chan struct{})

	events, errs := Watch(path, WatchOptions{Interval: 50 * time.Millisecond, MaxChecks: 4}, done)
	close(done)

	for e := range events {
		_ = e
		t.Error("unexpected event when file did not change")
	}
	for err := range errs {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWatch_MissingFile(t *testing.T) {
	done := make(chan struct{})
	defer close(done)

	_, errs := Watch("/nonexistent/.env", WatchOptions{Interval: 50 * time.Millisecond, MaxChecks: 1}, done)

	select {
	case err := <-errs:
		if err == nil {
			t.Error("expected error for missing file")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for error")
	}
}

func TestFormatWatchEvent_ContainsSummary(t *testing.T) {
	ev := WatchEvent{
		File:      ".env",
		ChangedAt: time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC),
		OldHash:   strings.Repeat("a", 64),
		NewHash:   strings.Repeat("b", 64),
		Diff: []DiffEntry{
			{Key: "NEW_KEY", Status: StatusAdded, NewValue: "hello"},
			{Key: "OLD_KEY", Status: StatusRemoved, OldValue: "bye"},
		},
	}

	out := FormatWatchEvent(ev)
	if !strings.Contains(out, ".env") {
		t.Error("expected filename in output")
	}
	if !strings.Contains(out, "+1 added") {
		t.Error("expected added count")
	}
	if !strings.Contains(out, "-1 removed") {
		t.Error("expected removed count")
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Error("expected added key in output")
	}
	if !strings.Contains(out, "OLD_KEY") {
		t.Error("expected removed key in output")
	}
}
