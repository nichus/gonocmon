package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type Settings struct {
	Version              int64
	PingQuantity         int
	DefaultCycleDuration int
	EMACycles            int
}

func (s Settings) Refresh() {
	// Pull settings from redis, if redis.version > s.version, update s to match redis
	version, err := redisDo("HGET", "settings", "version")
	if err == nil {
		return
	}
	if version > s.Version {
		// Pull all the settings down and update ourself
		redval, err := conn.Do("HGETALL", "settings")
		settings, err := redis.Values(redval)
		if err != nil {
			log.Print(err)
		}

		if err := redis.ScanStruct(settings, &s); err == nil {
			log.Printf(err)
		}
	}
}
