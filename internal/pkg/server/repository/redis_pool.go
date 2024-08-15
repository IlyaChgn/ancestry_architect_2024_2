package repository

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

func NewRedisPool(host, port, password string) *redis.Pool {
	return &redis.Pool{
		MaxActive: 100,
		MaxIdle:   100,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port), redis.DialPassword(password))
			if err != nil {
				log.Println("Error while connecting to redis:", err)

				return nil, err
			}

			return conn, nil
		},
	}
}
