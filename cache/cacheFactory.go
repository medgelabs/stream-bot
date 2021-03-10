package cache

import (
	"fmt"
	"log"
	"medgebot/config"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Acceptable types of caches
const (
	FILE = "file"
	MEM  = "mem"
)

// New creates a new Cache instance from the given Config
func New(config config.Config) (Cache, error) {
	cacheType := config.Cache()
	if cacheType == "" {
		return nil, errors.New("config key - cache not found / empty")
	}

	keyExpirationTime := config.CacheExpirationTime()

	switch strings.ToLower(cacheType) {
	case FILE:
		file, err := cacheFile("cache.txt")
		if err != nil {
			log.Fatalf("FATAL: load cache file - %v", err)
		}

		cache, err := FilePersisted(file, keyExpirationTime)
		if err != nil {
			log.Fatalf("FATAL: read cache file - %v", err)
		}
		return &cache, err
	case MEM:
		mem, _ := InMemory(keyExpirationTime)
		return &mem, nil
	default:
		return nil, fmt.Errorf("Invalid cacheType - %s. Valid values are: %v", cacheType, []string{FILE, MEM})
	}
}

func cacheFile(filepath string) (*os.File, error) {
	// If the file is more than 12 hours ago, clear
	cacheFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	return cacheFile, err
}
