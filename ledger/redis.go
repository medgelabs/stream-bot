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

func (l *RedisLedger) Absent(key string) bool {
	var res int
	err := l.pool.Do(radix.Cmd(&res, "EXISTS", key))
	if err != nil {
		log.Printf("ERROR: call EXISTS - %s", err)
		return false
	}

	return res == 0
}

func (l *RedisLedger) Add(key string) error {
	cmd := radix.FlatCmd(nil, "SET", key, 1, "EX", l.expireTimeInSeconds)

	if err := l.pool.Do(cmd); err != nil {
		log.Printf("ERROR: call SET - %s", err)
		return err
	}

	return nil
}
