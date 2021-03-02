package ledger

import (
	"fmt"
	"log"
	"medgebot/config"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Acceptable types of Ledgers
const (
	FILE = "file"
	MEM  = "mem"
)

func NewCache(config config.Config) (Cache, error) {
	cacheType := config.Cache()
	if cacheType == "" {
		return nil, errors.New("config key - cache not found / empty")
	}

	keyExpirationTime := config.CacheExpirationTime()

	switch strings.ToLower(cacheType) {
	case FILE:
		file, err := cacheFile("ledger.txt")
		if err != nil {
			log.Fatalf("FATAL: load ledger file - %v", err)
		}

		cache, err := NewFileCache(file, keyExpirationTime)
		if err != nil {
			log.Fatalf("FATAL: read ledger file - %v", err)
		}
		return &cache, err
	case MEM:
		mem, _ := NewInMemoryCache(keyExpirationTime)
		return &mem, nil
	default:
		return nil, fmt.Errorf("Invalid ledgerType - %s. Valid values are: %v", cacheType, []string{FILE, MEM})
	}
}

func cacheFile(filepath string) (*os.File, error) {
	// If the file is more than 12 hours ago, clear
	ledgerFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	return ledgerFile, err
}
