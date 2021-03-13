package cache

import (
	"fmt"
	"io"
	log "medgebot/logger"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
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
	// key is duplicated in the entry for rehydration purposes
	key       string
	value     string
	timestamp int64
}

// InMemory is a Cache that does not persist its in memory cache
func InMemory(keyExpirationSeconds int64) (PersistableCache, error) {
	return createCache(nil, keyExpirationSeconds)
}

// FilePersisted is a Cache that persists its cache to the given os.File
func FilePersisted(file *os.File, keyExpirationSeconds int64) (PersistableCache, error) {
	return createCache(file, keyExpirationSeconds)
}

// createCache constructs a cache instance.
// if os.File is nil, it is assumed to be a non-persisting Cache
// if keyExpirationSeconds <= 0, it is assumed to never expire keys
func createCache(file *os.File, keyExpirationSeconds int64) (PersistableCache, error) {
	cache := make(map[string]Entry)
	lineSeparator := "\n"
	fieldSeparator := "|"

	// To ensure we remove stale data, we rewrite state to the cache persistence target
	pc := PersistableCache{
		persistTarget:  file,
		persistent:     (file != nil),
		cache:          cache,
		lineSeparator:  lineSeparator,
		fieldSeparator: fieldSeparator,
		expiration:     keyExpirationSeconds,
	}

	if pc.persistent {
		pc.rehydrate()

		// Immediately flush to ensure any expired keys are removed
		pc.flushCache()

		// Setup cache flushing
		go func(cache *PersistableCache) {
			for {
				select {
				case <-time.After(10 * time.Second):
					cache.flushCache()
				}
			}
		}(&pc)
	}

	return pc, nil
}

// Get the value at the given key. Returns error if key not found
func (cache *PersistableCache) Get(key string) (string, error) {
	if cache.Absent(key) {
		return "", errors.Errorf("Key not found: %s", key)
	}

	entry, _ := cache.cache[key]
	return entry.value, nil
}

// Put the given key/value. If already present, the timestamp will
// be updated
func (cache *PersistableCache) Put(key, value string) error {
	ts := time.Now().Unix()
	entry := Entry{
		value:     value,
		timestamp: ts,
	}

	// Add to cache and write to persistence immediately
	cache.cache[key] = entry

	if cache.persistent {
		line := cache.line(key, entry)
		_, err := cache.persistTarget.Write([]byte(line))
		if err != nil {
			return errors.Wrap(err, "ERROR: failed to persist "+key+"")
		}
	}

	return nil
}

// Absent is true if the key is either not present or expired
func (cache *PersistableCache) Absent(key string) bool {
	entry, ok := cache.cache[key]
	return !ok || cache.expired(entry.timestamp)
}

// expired checks if the given key's timestamp is beyond the expiration threshold.
// Always returns false if an expiration is not set
func (cache *PersistableCache) expired(entryTs int64) bool {
	// If expiration is not set, all keys are considered active
	if cache.expiration <= 0 {
		return false
	}

	expirationTime := time.Now().Unix() - cache.expiration
	return entryTs < expirationTime
}

// rehydrate reads the given file and returns a hydrated FileCache
func (cache *PersistableCache) rehydrate() error {
	bytes, err := io.ReadAll(cache.persistTarget)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("ERROR: read cache persistence - %v", err))
	}

	// Hydrate cache from Reader
	str := string(bytes)
	for _, line := range strings.Split(str, cache.lineSeparator) {
		// skip empty lines
		if len(line) <= 1 {
			continue
		}

		entry, err := cache.fromLine(line)
		if err != nil {
			log.Error(err, "ERROR: invalid cache line - %s", line)
			continue
		}

		// Don't cache expired entries
		if cache.expired(entry.timestamp) {
			continue
		}

		// Valid entry, add to the cache
		cache.cache[entry.key] = entry
	}

	return nil
}

// line returns a string representation of a formatted file line
func (cache *PersistableCache) line(key string, entry Entry) string {
	var buf strings.Builder
	buf.WriteString(key + cache.fieldSeparator)
	buf.WriteString(entry.value + cache.fieldSeparator)
	buf.WriteString(fmt.Sprintf("%d", entry.timestamp))
	buf.WriteString(cache.lineSeparator)
	return buf.String()
}

// fromLine extracts an Entry from a cache persistence line
func (cache *PersistableCache) fromLine(line string) (Entry, error) {
	// Constants regarding line composition
	entryLineSize := 3
	keyIdx := 0
	valueIdx := 1
	tsIdx := 2

	tokens := strings.Split(line, cache.fieldSeparator)

	if len(tokens) != entryLineSize {
		return Entry{}, errors.Errorf("Invalid line length: %d", len(tokens))
	}

	key := tokens[keyIdx]
	value := tokens[valueIdx]
	ts, err := strconv.ParseInt(tokens[tsIdx], 10, 64)
	if err != nil {
		return Entry{}, errors.Errorf("invalid timestamp for key %s - %s", key, tokens[tsIdx])
	}

	return Entry{
		key:       key,
		value:     value,
		timestamp: ts,
	}, nil
}

// flushCache writes the state of l.cache to the given os.File
func (cache *PersistableCache) flushCache() {
	if cache.persistTarget == nil {
		log.Warn("flushCache called for a nil cache")
		return
	}

	// Truncate(0) clears the file contents
	cache.persistTarget.Truncate(0)

	for key, val := range cache.cache {
		line := cache.line(key, val)
		cache.persistTarget.Write([]byte(line))
	}
}
