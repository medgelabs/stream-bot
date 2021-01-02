package ledger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	ENTRY_LINE_SIZE = 3 // key | value | createdTs

	// Indices
	KEY   = 0
	VALUE = 1
	TS    = 2
)

// PersistableLedger is a in-memory cache backed by a FS file, if persistent
type PersistableLedger struct {
	cache          map[string]Entry
	persistent     bool
	ledger         *os.File
	lineSeparator  string
	fieldSeparator string
	expiration     int64
}

type Entry struct {
	key       string
	value     string
	timestamp int64
}

// NewInMemoryLedger is a Ledger that does not persist its in memory cache
func NewInMemoryLedger(keyExpirationSeconds int64) (PersistableLedger, error) {
	return NewFileLedger(nil, keyExpirationSeconds)
}

// FileLedger that persists its cache to the given os.File
// if os.File is nil, it is assumed to be a non-persisting Ledger
func NewFileLedger(ledger *os.File, keyExpirationSeconds int64) (PersistableLedger, error) {
	cache := make(map[string]Entry)
	lineSeparator := "\n"
	fieldSeparator := "|"

	// To ensure we remove stale data, we rewrite state to the ledger
	fileLedger := PersistableLedger{
		ledger:         ledger,
		persistent:     (ledger != nil),
		cache:          cache,
		lineSeparator:  lineSeparator,
		fieldSeparator: fieldSeparator,
		expiration:     keyExpirationSeconds,
	}

	if fileLedger.persistent {
		fileLedger.rehydrate()

		// Setup cache flushing
		fileLedger.flushCache()
		go func(ledger *PersistableLedger) {
			for {
				select {
				case <-time.After(10 * time.Second):
					fileLedger.flushCache()
				}
			}
		}(&fileLedger)
	}

	return fileLedger, nil
}

// Get the value at the given key. Returns error if key not found
func (l *PersistableLedger) Get(key string) (string, error) {
	if l.Absent(key) {
		return "", errors.Errorf("Key not found: %s", key)
	}

	entry, _ := l.cache[key]
	return entry.value, nil
}

// Put the given key/value. If already present, the timestamp will
// be updated
func (l *PersistableLedger) Put(key, value string) error {
	ts := time.Now().Unix()
	entry := Entry{
		key:       key,
		value:     value,
		timestamp: ts,
	}

	// Add to cache and write to ledger immediately
	l.cache[key] = entry

	if l.persistent {
		_, err := l.ledger.Write([]byte(l.line(key, entry)))
		if err != nil {
			return errors.Wrap(err, "ERROR: failed to add "+key+" to ledger")
		}
	}

	return nil
}

// Absent is true if the key is either not present or expired
func (l *PersistableLedger) Absent(key string) bool {
	entry, ok := l.cache[key]
	return !ok || l.expired(entry.timestamp)
}

// line returns a string representation of a formatted ledger line
func (l *PersistableLedger) line(key string, entry Entry) string {
	var buf strings.Builder
	buf.WriteString(key + l.fieldSeparator)
	buf.WriteString(entry.value + l.fieldSeparator)
	buf.WriteString(fmt.Sprintf("%d", entry.timestamp))
	return buf.String()
}

// expired checks if the given key's timestamp is beyond the expiration threshold
func (l *PersistableLedger) expired(entryTs int64) bool {
	expirationTime := time.Now().Unix() - l.expiration
	return entryTs < expirationTime
}

// rehydrate reads the given file and returns a hydrated FileLedger
func (l *PersistableLedger) rehydrate() error {
	bytes, err := ioutil.ReadAll(l.ledger)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("ERROR: read ledger - %v", err))
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
			log.Printf("ERROR: invalid ledger line - %s", line)
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

		expirationTime := time.Now().Unix() - l.expiration
		if ts < expirationTime {
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
func (l *PersistableLedger) flushCache() {
	if l.ledger == nil {
		log.Println("WARN: flushCache called for a nil ledger")
		return
	}

	// Truncate(0) clears the file contents
	l.ledger.Truncate(0)

	for username, ts := range l.cache {
		l.ledger.Write([]byte(l.line(username, ts)))
	}
}
