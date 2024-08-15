package repository

import (
	"github.com/gomodule/redigo/redis"
)

func NewRedisPool(host, password string) *redis.Pool {
	return &redis.Pool{
		MaxActive: 100,
		MaxIdle:   100,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host, redis.DialPassword(password))
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
	}
}
