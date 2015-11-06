package main

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

func get(c redis.Conn) {
	r, err := c.Do("GET", "ap-srv1")

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("get=%s", r)
}
func put(c redis.Conn) {
	r, err := c.Do("SET", "ap-srv1", time.Now())

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("put=%s", r)
}

func main() {
	c, err := redis.Dial("tcp", ":6379")

	if err != nil {
		log.Fatal(err)
	}

	get(c)

	defer c.Close()
}
