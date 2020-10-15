package ledger

import (
	"fmt"
	"log"

	"github.com/mediocregopher/radix/v3"
)

type Redis struct {
	pool                *radix.Pool
	expireTimeInSeconds int
}

func NewRedis(host, port string) (*Redis, error) {
	pool, err := radix.NewPool("tcp", fmt.Sprintf("%s:%s", host, port), 10)
	if err != nil {
		return nil, err
	}

	return &Redis{
		pool:                pool,
		expireTimeInSeconds: 60 * 60 * 12, // 12 hours
	}, nil
}

func (l *Redis) absent(key string) bool {
	var res int
	err := l.pool.Do(radix.Cmd(&res, "EXISTS", key))
	if err != nil {
		log.Printf("ERROR: call EXISTS - %s", err)
		return false
	}

	return res == 0
}

func (l *Redis) add(key string) error {
	cmd := radix.FlatCmd(nil, "SET", key, 1, "EX", l.expireTimeInSeconds)

	if err := l.pool.Do(cmd); err != nil {
		log.Printf("ERROR: call SET - %s", err)
		return err
	}

	return nil
}
