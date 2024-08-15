package repository

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func NewRedisPool(host, port, password string) *redis.Pool {
	return &redis.Pool{
		MaxActive: 100,
		MaxIdle:   100,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port), redis.DialPassword(password))
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
	}
}
