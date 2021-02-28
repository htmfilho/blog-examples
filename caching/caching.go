package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
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
	conn *redis.Client
	ctx  context.Context
}

// GetCachingMechanism initializes and returns a caching mechanism.
func GetCachingMechanism() Cache {
	cch := &RedisCache{
		conn: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}

	cch.ctx = context.Background()

	return cch
}

// Put adds an entry in the cache.
func (rc *RedisCache) Put(key string, value interface{}) {
	if err := rc.conn.Set(rc.ctx, key, value, 0); err != nil {
		fmt.Println(err)
	}
}

// PutAll adds the entries of a map in the cache.
func (rc *RedisCache) PutAll(entries map[string]interface{}) {
	for k, v := range entries {
		rc.Put(k, v)
	}
}

// Get gets an entry from the cache.
func (rc *RedisCache) Get(key string) interface{} {
	value, err := rc.conn.Get(rc.ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return value
}

// GetAll gets all the entries of a map from the cache.
func (rc *RedisCache) GetAll(keys []string) map[string]interface{} {
	entries := make(map[string]interface{})
	for _, k := range keys {
		entries[k] = rc.Get(k)
	}

	return entries
}

// Clean cleans a entry from the cache.
func (rc *RedisCache) Clean(key string) {
	if err := rc.conn.Del(rc.ctx, key); err != nil {
		fmt.Println(err)
	}
}

// CleanAll cleans the entire cache.
func (rc *RedisCache) CleanAll() {
	rc.conn.FlushDB(rc.ctx)
}
