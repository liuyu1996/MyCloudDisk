package redis

import (
	"MyCloudDisk/config"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisCli struct {
	Pool *redis.Pool
}

var RedisClient *RedisCli

func Default()  {
	rediscli := new(RedisCli)
	rediscli.Pool = &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			conn, e = redis.Dial(
				"tcp",
				config.Configs.RedisHost,
				redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
			)
			if e != nil {
				return nil, e
			}
			if _, e =conn.Do("AUTH", config.Configs.RedisPwd); e != nil{
				_ = conn.Close()
				return nil, e
			}
			return conn, nil
		},
		MaxIdle:         256,
		MaxActive:       30,
		IdleTimeout:     time.Duration(120),
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	RedisClient = rediscli
}
