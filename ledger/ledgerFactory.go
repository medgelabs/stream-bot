package ledger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// ./bot -ledger=redis

// Acceptable types of Ledgers
const (
	REDIS = "redis"
	FILE  = "file"
	MEM   = "mem"
	DB    = "db"
)

func NewLedger(ledgerType string) (Ledger, error) {
	switch strings.ToLower(ledgerType) {
	case REDIS:
		redisHost := os.Getenv("REDIS_HOST")
		redisPort := os.Getenv("REDIS_PORT")
		redis, err := NewRedisLedger(redisHost, redisPort)
		return &redis, err
	case FILE:
		file, err := os.OpenFile("ledger.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		ledger, err := NewFileLedger(file)
		if err != nil {
			log.Fatalf("FATAL: read ledger file - %v", err)
		}
		return &ledger, err
	case MEM:
		mem := NewInMemoryLedger()
		return &mem, nil
	case DB:
		return nil, fmt.Errorf("Don't use this")
	default:
		return nil, fmt.Errorf("Invalid ledgerType " + ledgerType)
	}
}
