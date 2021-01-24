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
	REDIS = "redis"
	FILE  = "file"
	MEM   = "mem"
)

func NewLedger(config config.Config) (Ledger, error) {
	ledgerType := config.Ledger()
	if ledgerType == "" {
		return nil, errors.New("config key - ledger not found / empty")
	}

	keyExpirationTime := config.LedgerExpirationTime()

	switch strings.ToLower(ledgerType) {
	case REDIS:
		redisHost := config.RedisHost()
		redisPort := config.RedisPort()
		redis, err := NewRedisLedger(redisHost, redisPort, keyExpirationTime)
		return &redis, err
	case FILE:
		file, err := ledgerFile("ledger.txt")
		if err != nil {
			log.Fatalf("FATAL: load ledger file - %v", err)
		}

		ledger, err := NewFileLedger(file, keyExpirationTime)
		if err != nil {
			log.Fatalf("FATAL: read ledger file - %v", err)
		}
		return &ledger, err
	case MEM:
		mem, _ := NewInMemoryLedger(keyExpirationTime)
		return &mem, nil
	default:
		return nil, fmt.Errorf("Invalid ledgerType - %s. Valid values are: %s, %s, %s", ledgerType, REDIS, FILE, MEM)
	}
}

func ledgerFile(filepath string) (*os.File, error) {
	// If the file is more than 12 hours ago, clear
	ledgerFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	return ledgerFile, err
}
