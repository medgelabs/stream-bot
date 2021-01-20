package ledger

import (
	"fmt"
	"log"

	"github.com/mediocregopher/radix/v3"
)

type RedisLedger struct {
	pool                *radix.Pool
	expireTimeInSeconds int64
}

func NewRedisLedger(host, port string, keyExpirationSeconds int64) (RedisLedger, error) {
	pool, err := radix.NewPool("tcp", fmt.Sprintf("%s:%s", host, port), 10)
	if err != nil {
		return RedisLedger{}, err
	}

	return RedisLedger{
		pool:                pool,
		expireTimeInSeconds: keyExpirationSeconds,
	}, nil
}

// Get the value at the specified key, if it exists
func (l *RedisLedger) Get(key string) (string, error) {
	var result string
	cmd := radix.Cmd(result, "GET", key)

	if err := l.pool.Do(cmd); err != nil {
		log.Printf("ERROR: call GET - %s", err)
		return "", err
	}

	return result, nil
}

// Put a value to the Redis store
func (l *RedisLedger) Put(key, value string) error {
	cmd := radix.FlatCmd(nil, "SET", key, value, "EX", l.expireTimeInSeconds)

	if err := l.pool.Do(cmd); err != nil {
		log.Printf("ERROR: call SET - %s", err)
		return err
	}

	return nil
}

// Check if the given key exists in the Redis store
func (l *RedisLedger) Absent(key string) bool {
	var res int
	err := l.pool.Do(radix.Cmd(&res, "EXISTS", key))
	if err != nil {
		log.Printf("ERROR: call EXISTS - %s", err)
		return false
	}

	return res == 0
}
