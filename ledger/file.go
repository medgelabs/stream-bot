package ledger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// FileLedger is a in-memory cache backed by a FS file
// Note: the map is map[username]UNIX_EPOCH_SECONDS to account for expiring keys
type FileLedger struct {
	ledger         *os.File
	cache          map[string]int64
	lineSeparator  string
	fieldSeparator string
	expiration     int64
}

func NewFileLedger(ledger *os.File, keyExpirationTime int64) (FileLedger, error) {
	cache := make(map[string]int64)
	lineSeparator := "\n"
	fieldSeparator := ","

	bytes, err := ioutil.ReadAll(ledger)
	if err != nil {
		log.Printf("ERROR: read ledger - %v", err)
		return FileLedger{}, err
	}

	// Hydrate cache from Reader
	str := string(bytes)
	for _, line := range strings.Split(str, lineSeparator) {
		// skip empty lines
		if len(line) <= 1 {
			continue
		}

		tokens := strings.Split(line, fieldSeparator)

		// Don't add invalid lines
		if len(tokens) != 2 {
			log.Printf("ERROR: invalid ledger line - %s", line)
			continue
		}

		// And don't add keys that are expired
		key := tokens[0]
		ts, err := strconv.ParseInt(tokens[1], 10, 64)
		if err != nil {
			log.Printf("ERROR: invalid timestamp for key %s - %s", key, tokens[1])
			continue
		}

		expirationTime := time.Now().Unix() - keyExpirationTime
		if ts < expirationTime {
			continue
		}

		// Valid entry, add to the cache
		cache[key] = ts
	}

	// To ensure we remove stale data, we rewrite state to the ledger
	fileLedger := FileLedger{
		ledger:         ledger,
		cache:          cache,
		lineSeparator:  lineSeparator,
		fieldSeparator: fieldSeparator,
		expiration:     keyExpirationTime,
	}

	// Setup cache flushing
	fileLedger.flushCache()
	go func(ledger *FileLedger) {
		for {
			select {
			case <-time.After(10 * time.Second):
				fileLedger.flushCache()
			}
		}
	}(&fileLedger)

	return fileLedger, nil
}

func (l *FileLedger) flushCache() {
	// Truncate(0) clears the file contents
	l.ledger.Truncate(0)

	for username, ts := range l.cache {
		l.ledger.Write([]byte(l.line(username, ts)))
	}
}

func (l *FileLedger) Absent(key string) bool {
	entryTs, ok := l.cache[key]
	return !ok || l.expired(entryTs)
}

func (l *FileLedger) Add(key string) error {
	// Don't add a key if it already exists && is not expired
	entryTs, ok := l.cache[key]
	if ok || !l.expired(entryTs) {
		return nil
	}

	ts := time.Now().Unix()
	l.cache[key] = ts
	l.ledger.Write([]byte(l.line(key, ts)))

	return nil
}

// line returns a string representation of a formatted ledger line
func (l *FileLedger) line(username string, timestamp int64) string {
	return fmt.Sprintf("%s%s%d%s", username, l.fieldSeparator, timestamp, l.lineSeparator)
}

// expired checks if the given key's timestamp is beyond the expiration threshold
func (l *FileLedger) expired(entryTs int64) bool {
	expirationTime := time.Now().Unix() - l.expiration
	return entryTs < expirationTime
}
