package store

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

// TXTFile is a store that reads and writes to a text file.
//
// Note: This is not a performant store, it is only intended for
// small amounts of data, and does not handle concurrent access
// beyond a single process.
type TXTFile struct {
	file string

	mx sync.Mutex
}

var _ Store = (*TXTFile)(nil)

func NewTXTFile(path string) *TXTFile {
	return &TXTFile{file: path}
}

func validate(key, value string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if strings.Contains(key, "->") {
		return fmt.Errorf("key cannot contain '->'")
	}
	if strings.Contains(key, "\n") || strings.Contains(value, "\n") {
		return fmt.Errorf("key and value cannot contain newlines")
	}

	return nil
}

func (t *TXTFile) load() ([]entry, error) {
	fd, err := os.Open(t.file)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var entries []entry
	s := bufio.NewScanner(fd)
	for s.Scan() {
		k, v, ok := strings.Cut(s.Text(), "->")
		if !ok {
			// skip invalid lines
			continue
		}

		entries = append(entries, entry{strings.TrimSpace(k), strings.TrimSpace(v)})
	}
	if s.Err() != nil {
		return nil, s.Err()
	}

	return entries, nil
}

func (t *TXTFile) save(e []entry) error {
	fd, err := os.Create(t.file)
	if err != nil {
		return err
	}
	defer fd.Close()

	var maxKeyLen int
	for _, e := range e {
		if len(e.k) > maxKeyLen {
			maxKeyLen = len(e.k)
		}
	}

	for _, e := range e {
		if len(e.k) < maxKeyLen {
			e.k = e.k + strings.Repeat(" ", maxKeyLen-len(e.k))
		}
		_, err = fmt.Fprintf(fd, "%s -> %s\n", e.k, e.v)
		if err != nil {
			return err
		}
	}

	return fd.Close()
}

func (t *TXTFile) Get(ctx context.Context, key string) (string, error) {
	t.mx.Lock()
	defer t.mx.Unlock()

	entries, err := t.load()
	if err != nil {
		return "", err
	}

	for _, e := range entries {
		if e.k != key {
			continue
		}

		return e.v, nil
	}

	return "", nil
}

func (t *TXTFile) Set(ctx context.Context, key, value string) error {
	err := validate(key, value)
	if err != nil {
		return err
	}

	t.mx.Lock()
	defer t.mx.Unlock()

	entries, err := t.load()
	if err != nil {
		return err
	}

	var found bool
	for i, e := range entries {
		if e.k != key {
			continue
		}
		found = true

		if value == "" {
			entries = append(entries[:i], entries[i+1:]...)
			break
		}

		entries[i].v = value
		break
	}

	if !found {
		entries = append(entries, entry{key, value})
	}

	return t.save(entries)
}

func (t *TXTFile) Scanner(ctx context.Context) (Scanner, error) {
	t.mx.Lock()
	defer t.mx.Unlock()

	entries, err := t.load()
	if err != nil {
		return nil, err
	}

	return &entryScanner{
		results: entries,
	}, nil
}
