package main

import (
	"flag"
	"html/template"

	"github.com/garyburd/redigo/redis"
)

var (
	pool                        *redis.Pool
	templates                   *template.Template
	templateDirectory           string
	redisPassword, redisAddress string
	redisMaxConn                int
)

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

func init() {
	flag.StringVar(&redisPassword, "redisPassword", "", "Password used to connect to redisAddress")
	flag.StringVar(&redisAddress, "redisAddress", ":6379", "Address to the redis server")
	flag.IntVar(&redisMaxConn, "max-connections", 10, "Max connections to redis")
	flag.StringVar(&templateDirectory, "templateDirectory", "/home/ovandenb/gnm-templates", "Path to directory containing the html templates")

	flag.Parse()
	return
}

func main() {
	var waitGroup WaitGroupWrapper

	flag.Parse()
	pool = newPool(redisAddress, redisPassword)
	defer pool.Close()

	waitGroup.Wrap(func() { httpServer() })

	waitGroup.Wait()
}
