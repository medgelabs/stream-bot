package ledger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Acceptable types of Ledgers
const (
	REDIS = "redis"
	FILE  = "file"
	MEM   = "mem"
)

func NewLedger(ledgerType string) (Ledger, error) {
	switch strings.ToLower(ledgerType) {
	case REDIS:
		redisHost := os.Getenv("REDIS_HOST")
		redisPort := os.Getenv("REDIS_PORT")
		redis, err := NewRedisLedger(redisHost, redisPort)
		return &redis, err
	case FILE:
		file, err := ledgerFile("ledger.txt")
		if err != nil {
			log.Fatalf("FATAL: load ledger file - %v", err)
		}

		ledger, err := NewFileLedger(file)
		if err != nil {
			log.Fatalf("FATAL: read ledger file - %v", err)
		}
		return &ledger, err
	case MEM:
		mem := NewInMemoryLedger()
		return &mem, nil
	default:
		return nil, fmt.Errorf("Invalid ledgerType " + ledgerType)
	}
}

func ledgerFile(filepath string) (*os.File, error) {
	stats, err := os.Stat(filepath)
	if err != nil {
		log.Printf("ERROR: failed to stat %s - %v", filepath, err)
		return &os.File{}, err
	}
	lastModTime := stats.ModTime()

	// Due to Sub()'s weird rejection of time.Duration, we add negative
	// to get 12 hours in the past
	twelveHoursAgo := time.Now().Add(-12 * time.Hour)

	// If the file is more than 12 hours ago, clear
	if lastModTime.Before(twelveHoursAgo) {
		log.Println("Ledger file older than 12 hours. Clearing..")

		err := os.Remove(filepath)
		if err != nil {
			log.Printf("ERROR: failed to remove ledger file %s - %v", filepath, err)
			return &os.File{}, err
		}
	}

	ledgerFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	return ledgerFile, err
}
