package ledger

import (
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// TODO expiration consideration
// Do we use fs timestamp (per shito86)? or do we need structured files like JSON with a ts?

type FileLedger struct {
	ledger    io.ReadWriter
	cache     map[string]int
	separator string
}

func NewFileLedger(ledger io.ReadWriter) (FileLedger, error) {
	cache := make(map[string]int)
	separator := "\n"

	bytes, err := ioutil.ReadAll(ledger)
	if err != nil {
		log.Printf("ERROR: read ledger - %v", err)
		return FileLedger{}, err
	}

	// Hydrate cache from Reader
	str := string(bytes)
	for _, user := range strings.Split(str, separator) {
		cache[user] = 1
	}

	return FileLedger{
		ledger:    ledger,
		cache:     cache,
		separator: separator,
	}, nil
}

func (l *FileLedger) Absent(key string) bool {
	_, ok := l.cache[key]
	return !ok
}

func (l *FileLedger) Add(key string) error {
	l.cache[key] = 1
	l.ledger.Write([]byte(key + l.separator))

	return nil
}
