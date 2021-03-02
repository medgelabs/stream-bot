package ledger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	// ENTRY_LINE_SIZE declares how entries are written to the File
	ENTRY_LINE_SIZE = 3 // key | value | createdTs

	// Indices
	KEY   = 0
	VALUE = 1
	TS    = 2
)

// PersistableCache is a in-memory cache backed by a FS file, if persistent.
// Key expiration is enabled by setting an expiration > 0. Disabled if <= 0
type PersistableCache struct {
	cache          map[string]Entry
	persistent     bool
	persistTarget  *os.File
	lineSeparator  string
	fieldSeparator string
	expiration     int64
}

// Entry represents entries in the cache
type Entry struct {
	key       string
	value     string
	timestamp int64
}

// NewInMemoryCache is a Cache that does not persist its in memory cache
func NewInMemoryCache(keyExpirationSeconds int64) (PersistableCache, error) {
	return NewFileCache(nil, keyExpirationSeconds)
}

// NewFileCache that persists its cache to the given os.File
// if os.File is nil, it is assumed to be a non-persisting Cache
func NewFileCache(file *os.File, keyExpirationSeconds int64) (PersistableCache, error) {
	cache := make(map[string]Entry)
	lineSeparator := "\n"
	fieldSeparator := "|"

	// To ensure we remove stale data, we rewrite state to the cache persistence target
	pl := PersistableCache{
		persistTarget:  file,
		persistent:     (file != nil),
		cache:          cache,
		lineSeparator:  lineSeparator,
		fieldSeparator: fieldSeparator,
		expiration:     keyExpirationSeconds,
	}

	if pl.persistent {
		pl.rehydrate()

		// Setup cache flushing
		pl.flushCache()
		go func(cache *PersistableCache) {
			for {
				select {
				case <-time.After(10 * time.Second):
					cache.flushCache()
				}
			}
		}(&pl)
	}

	return pl, nil
}

// Get the value at the given key. Returns error if key not found
func (l *PersistableCache) Get(key string) (string, error) {
	if l.Absent(key) {
		return "", errors.Errorf("Key not found: %s", key)
	}

	entry, _ := l.cache[key]
	return entry.value, nil
}

// Put the given key/value. If already present, the timestamp will
// be updated
func (l *PersistableCache) Put(key, value string) error {
	ts := time.Now().Unix()
	entry := Entry{
		key:       key,
		value:     value,
		timestamp: ts,
	}

	// Add to cache and write to persistence immediately
	l.cache[key] = entry

	if l.persistent {
		_, err := l.persistTarget.Write([]byte(l.line(key, entry)))
		if err != nil {
			return errors.Wrap(err, "ERROR: failed to persist "+key+"")
		}
	}

	return nil
}

// Absent is true if the key is either not present or expired
func (l *PersistableCache) Absent(key string) bool {
	entry, ok := l.cache[key]
	return !ok || l.expired(entry.timestamp)
}

// line returns a string representation of a formatted file line
func (l *PersistableCache) line(key string, entry Entry) string {
	var buf strings.Builder
	buf.WriteString(key + l.fieldSeparator)
	buf.WriteString(entry.value + l.fieldSeparator)
	buf.WriteString(fmt.Sprintf("%d", entry.timestamp))
	buf.WriteString(l.lineSeparator)
	return buf.String()
}

// expired checks if the given key's timestamp is beyond the expiration threshold.
// Always returns false if an expiration is not set
func (l *PersistableCache) expired(entryTs int64) bool {
	// If expiration is not set, all keys are considered active
	if l.expiration <= 0 {
		return false
	}

	expirationTime := time.Now().Unix() - l.expiration
	return entryTs < expirationTime
}

// rehydrate reads the given file and returns a hydrated FileCache
func (l *PersistableCache) rehydrate() error {
	bytes, err := io.ReadAll(l.persistTarget)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("ERROR: read cache persistence - %v", err))
	}

	// Hydrate cache from Reader
	str := string(bytes)
	for _, line := range strings.Split(str, l.lineSeparator) {
		// skip empty lines
		if len(line) <= 1 {
			continue
		}

		tokens := strings.Split(line, l.fieldSeparator)

		// Don't add invalid lines
		if len(tokens) != ENTRY_LINE_SIZE {
			log.Printf("ERROR: invalid cache line - %s", line)
			continue
		}

		// And don't add keys that are expired
		key := tokens[KEY]
		value := tokens[VALUE]
		ts, err := strconv.ParseInt(tokens[TS], 10, 64)
		if err != nil {
			log.Printf("ERROR: invalid timestamp for key %s - %s", key, tokens[TS])
			continue
		}

		if l.expired(ts) {
			continue
		}

		// Valid entry, add to the cache
		l.cache[key] = Entry{
			key:       key,
			value:     value,
			timestamp: ts,
		}
	}

	return nil
}

// flushCache writes the state of l.cache to the given os.File
func (l *PersistableCache) flushCache() {
	if l.persistTarget == nil {
		log.Println("WARN: flushCache called for a nil cache")
		return
	}

	// Truncate(0) clears the file contents
	l.persistTarget.Truncate(0)

	for key, val := range l.cache {
		line := l.line(key, val)
		l.persistTarget.Write([]byte(line))
	}
}
