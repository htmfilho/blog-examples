package main

import (
	"testing"
)

func TestRedisCache_Clean(t *testing.T) {
	cache := GetCachingMechanism()
	cache.Put("single", "Single Record")

	cache.Clean("single")
	if "Single Record" == cache.Get("single") {
		t.Fail()
	}
}

func TestRedisCache_CleanAll(t *testing.T) {
	cache := GetCachingMechanism()

	keys := []string{"multiple1", "multiple2", "multiple3"}
	entries := make(map[string]interface{})
	entries[keys[0]] = "Multiple 1"
	entries[keys[1]] = "Multiple 2"
	entries[keys[2]] = "Multiple 3"
	cache.PutAll(entries)

	cache.CleanAll()

	for k, _ := range entries {
		if entries[k] == cache.Get(k) {
			t.Fail()
		}
	}
}

func TestRedisCache_Get(t *testing.T) {
	cache := GetCachingMechanism()
	cache.Put("single", "Single Record")
	if "Single Record" != cache.Get("single") {
		t.Fail()
	}
}

func TestRedisCache_GetAll(t *testing.T) {
	cache := GetCachingMechanism()

	keys := []string{"multiple1", "multiple2", "multiple3"}
	entries := make(map[string]interface{})
	entries[keys[0]] = "Multiple 1"
	entries[keys[1]] = "Multiple 2"
	entries[keys[2]] = "Multiple 3"
	cache.PutAll(entries)

	results := cache.GetAll(keys)

	for k, v := range results {
		if entries[k] != v {
			t.Fail()
		}
	}
}

func TestRedisCache_Put(t *testing.T) {
	cache := GetCachingMechanism()
	cache.Put("single", "Single Record")
	if "Single Record" != cache.Get("single") {
		t.Fail()
	}
}

func TestRedisCache_PutAll(t *testing.T) {
	cache := GetCachingMechanism()

	keys := []string{"multiple1", "multiple2", "multiple3"}
	entries := make(map[string]interface{})
	entries[keys[0]] = "Multiple 1"
	entries[keys[1]] = "Multiple 2"
	entries[keys[2]] = "Multiple 3"
	cache.PutAll(entries)

	results := cache.GetAll(keys)

	for k, v := range results {
		if entries[k] != v {
			t.Fail()
		}
	}
}
