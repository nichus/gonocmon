package main

import (
	"log"
	"time"


	"github.com/garyburd/redigo/redis"
)

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if redisPassword != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func redisDo(args ...interface{}) {
	conn := pool.Get()
	defer conn.Close()

	result, err := conn.Do(args...)
	if err != nil {
		log.Printf("Unable to execute redis command '%s' received error '%s'", string.Join(args, " "), err)
	}
	return result, err
}

/*
func loadSettings() {
	count, err := redisDo("EXISTS", "settings")
	if err != nil && count == 0 {
		redisDo()
	}
}
*/
