package main

import (
	"flag"
	"fmt"
	"html/template"
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

func DashboardHandler(response http.ResponseWriter, request *http.Request) {
	conn := pool.Get()
	defer conn.Close()
	timestamp, err := conn.Do("GET", "ap-srv1")
	if err != nil {
		fmt.Fprintf(response, "Error retrieving value from redis: %s", err)
	}
	details := struct {
		Time string
	}{
		fmt.Sprintf("%s", timestamp),
	}
	t, _ := template.ParseFiles("templates/dashboard.html")
	t.Execute(response, details)
}

func CreateSystemHandler(response http.ResponseWriter, request *http.Request) {
	conn := pool.Get()
	defer conn.Close()
	timestamp, err := conn.Do("GET", "systemList")
	if err != nil {
		fmt.Fprintf(response, "Error retrieving value from redis: %s", err)
	}
	details := struct {
		Time string
	}{
		fmt.Sprintf("%s", timestamp),
	}
	t, _ := template.ParseFiles("templates/dashboard.html")
	t.Execute(response, details)
}

func main() {
	flag.Parse()
	pool = newPool(*redisAddress, *redisPassword)
	defer pool.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", DashboardHandler)
	router.HandleFunc("/systems/create", CreateSystemHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
