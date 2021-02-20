package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Cache is the adapter contract between the application and the cache library.
type Cache interface {
	Put(key string, value interface{})
	PutAll(map[string]interface{})
	Get(key string) interface{}
	GetAll(keys []string) map[string]interface{}
	Clean(key string)
	CleanAll()
}

// RedisCache holds a Redis connection pool.
type RedisCache struct {
	conn *redis.Pool
}

// Put adds an entry in the cache.
func (rc *RedisCache) Put(key string, value interface{}) {
	_, err := rc.conn.Get().Do("SET", key, value)
	if err != nil {
		fmt.Println(err)
	}
}

// PutAll adds the entries of a map in the cache.
func (rc *RedisCache) PutAll(map[string]interface{}) {

}

// Get gets an entry from the cache.
func (rc *RedisCache) Get(key string) interface{} {
	value, err := redis.String(rc.conn.Get().Do("GET", key))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return value
}

// GetAll gets all the entries of a map from the cache.
func (rc *RedisCache) GetAll(keys []string) map[string]interface{} {
	return nil
}

// Clean cleans a entry from the cache.
func (rc *RedisCache) Clean(key string) {

}

// CleanAll cleans the entire cache.
func (rc *RedisCache) CleanAll() {

}

// GetCachingMechanism initializes and returns a caching mechanism.
func GetCachingMechanism() Cache {
	cache := &RedisCache{
		conn: &redis.Pool{
			MaxIdle:     7,
			MaxActive:   30,
			IdleTimeout: 60 * time.Second,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", "localhost:6379")
				if err != nil {
					fmt.Println(err)
					return nil, err
				}

				if _, err := conn.Do("SELECT", 0); err != nil {
					conn.Close()
					fmt.Println(err)
					return nil, err
				}

				return conn, nil
			},
			TestOnBorrow: func(conn redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := conn.Do("PING")
				fmt.Println(err)
				return err
			},
		},
	}
	return cache
}

func main() {
	cache := GetCachingMechanism()

	cache.Put("test", "data")

	fmt.Println(cache.Get("test"))
}
