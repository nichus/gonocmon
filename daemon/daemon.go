package main

import (
	"flag"
	"html/template"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool           *redis.Pool
	redisAddress   = flag.String("redisAddress", ":6379", "Address to the redis server")
	redisPassword  = flag.String("redisPassword", "", "Password for the redisAddress")
	maxConnections = flag.Int("max-connections", 10, "Max connections to redis")
	templates      *template.Template
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
			/*
			 * if _, err := c.Do("AUTH", password); err != nil {
			 * 	c.Close()
			 * 	return nil, err
			 * }
			 */
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func retrieveValue(key string) (reply interface{}, err error) {
	conn := pool.Get()
	defer conn.Close()

	count, err := conn.Do("EXISTS", key)
	if err == nil || count == 0 {
		return nil, err
	}
	value, err := conn.Do("")
	if err == nil {
		return nil, err
	}
	return value, err
}

func main() {
	var waitGroup WaitGroupWrapper

	flag.Parse()
	pool = newPool(*redisAddress, *redisPassword)
	defer pool.Close()

	waitGroup.Wrap(func() { httpServer() })

	waitGroup.Wait()
}
