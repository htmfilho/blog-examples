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
func (rc *RedisCache) PutAll(entries map[string]interface{}) {
	c := rc.conn.Get()
	for k, v := range entries {
		err := c.Send("SET", k, v)
		if err != nil {
			fmt.Println(err)
		}
	}

	err := c.Flush()
	if err != nil {
		fmt.Println(err)
	}

	_, err = c.Receive()
	if err != nil {
		fmt.Println(err)
	}
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
	// Converts []string to []interface{} since Go doesn't do it explicitly
	// because it doesn't want the syntax to hide a O(n) operation.
	intKeys := make([]interface{}, len(keys))
	for i, _ := range keys {
		intKeys[i] = keys[i]
	}

	c := rc.conn.Get()

	values, err := redis.Strings(c.Do("MGET", intKeys...))

	entries := make(map[string]interface{})
	for i, k := range keys {
		entries[k] = values[i]
		if err != nil {
			fmt.Println(err)
		}
	}

	return entries
}

// Clean cleans a entry from the cache.
func (rc *RedisCache) Clean(key string) {
	_, err := rc.conn.Get().Do("DEL", key)
	if err != nil {
		fmt.Println(err)
	}
}

// CleanAll cleans the entire cache.
func (rc *RedisCache) CleanAll() {
	_, err := rc.conn.Get().Do("FLUSHDB")
	if err != nil {
		fmt.Println(err)
	}
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

	cache.Put("single", "Single Record")

	fmt.Println(cache.Get("single"))

	keys := []string{"multiple1", "multiple2", "multiple3"}
	entries := make(map[string]interface{})
	entries[keys[0]] = "Multiple 1"
	entries[keys[1]] = "Multiple 2"
	entries[keys[2]] = "Multiple 3"
	cache.PutAll(entries)

	entries = cache.GetAll(keys)

	for k, v := range entries {
		fmt.Print(k)
		fmt.Print(" = ")
		fmt.Println(v)
	}

	fmt.Println(cache.Get("single"))
	cache.Clean("single")
	fmt.Println(cache.Get("single"))

	cache.CleanAll()
	entries = cache.GetAll(keys)
	for k, v := range entries {
		fmt.Print(k)
		fmt.Print(" = ")
		fmt.Println(v)
	}
}
