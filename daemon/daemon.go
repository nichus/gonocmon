package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

var (
	pool           *redis.Pool
	redisAddress   = flag.String("redisAddress", ":6379", "Address to the redis server")
	redisPassword  = flag.String("redisPassword", "", "Password for the redisAddress")
	maxConnections = flag.Int("max-connections", 10, "Max connections to redis")
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

func DashboardHandler(response http.ResponseWriter, request *http.Request) {
	conn := pool.Get()
	defer conn.Close()
	fmt.Fprintf(response, "Welcome to the Dashboard on: %s.", request.URL.Path[1:])
	fmt.Fprintf(response, "\n\n<br /><br />\n\n")

	timestamp, err := conn.Do("GET", "ap-srv1")
	if err != nil {
		fmt.Fprintf(response, "Error retrieving value from redis: %s", err)
	}
	fmt.Fprintf(response, "Last Timestamp was: %s.", timestamp)
}

func main() {
	flag.Parse()
	pool = newPool(*redisAddress, *redisPassword)
	defer pool.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", DashboardHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
