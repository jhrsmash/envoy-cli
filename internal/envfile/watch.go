package envfile

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchEvent describes a change detected in a watched .env file.
type WatchEvent struct {
	File      string
	ChangedAt time.Time
	OldHash   string
	NewHash   string
	Diff      []DiffEntry
}

// WatchOptions controls the behaviour of Watch.
type WatchOptions struct {
	// Interval between polls. Defaults to 2 seconds when zero.
	Interval time.Duration
	// MaxChecks limits the number of poll iterations (0 = unlimited).
	MaxChecks int
}

// Watch polls the given .env file for changes and sends a WatchEvent on the
// returned channel whenever the file content changes. The caller must close
// the done channel to stop watching. Any error encountered while reading the
// file is sent on errCh.
func Watch(path string, opts WatchOptions, done <-chan struct{}) (<-chan WatchEvent, <-chan error) {
	events := make(chan WatchEvent, 4)
	errors := make(chan error, 4)

	interval := opts.Interval
	if interval <= 0 {
		interval = 2 * time.Second
	}

	go func() {
		defer close(events)
		defer close(errors)

		currentHash, currentEnv, err := hashAndParse(path)
		if err != nil {
			errors <- err
			return
		}

		checks := 0
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				newHash, newEnv, err := hashAndParse(path)
				if err != nil {
					errors <- err
					continue
				}
				if newHash != currentHash {
					diff := Diff(currentEnv, newEnv)
					events <- WatchEvent{
						File:      path,
						ChangedAt: time.Now(),
						OldHash:   currentHash,
						NewHash:   newHash,
						Diff:      diff,
					}
					currentHash = newHash
					currentEnv = newEnv
				}
				checks++
				if opts.MaxChecks > 0 && checks >= opts.MaxChecks {
					return
				}
			}
		}
	}()

	return events, errors
}

func hashAndParse(path string) (string, map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", nil, fmt.Errorf("watch: open %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	data, err := io.ReadAll(io.TeeReader(f, h))
	if err != nil {
		return "", nil, fmt.Errorf("watch: read %s: %w", path, err)
	}

	env, err := Parse(string(data))
	if err != nil {
		return "", nil, fmt.Errorf("watch: parse %s: %w", path, err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), env, nil
}
