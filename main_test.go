package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

// TestMemCacheSetGet tests the Set and Get methods of MemCache
func TestMemCacheSetGet(t *testing.T) {
	cache := NewCache(WithTTL(5 * time.Second))

	// Test set operation
	err := cache.Set("bucket", "key", []byte("value"))
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Test get operation
	val, err := cache.Get("bucket", "key")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if !bytes.Equal(val, []byte("value")) {
		t.Errorf("Get returned incorrect value: got %v, want %v", val, []byte("value"))
	}
}

// TestMemCacheExpiration tests that entries are properly expired based on TTL
func TestMemCacheExpiration(t *testing.T) {
	cache := NewCache(WithTTL(1 * time.Second))

	_ = cache.Set("bucket", "key", []byte("value"))
	time.Sleep(2 * time.Second) // wait for the key to expire

	_, err := cache.Get("bucket", "key")
	if err == nil || err.Error() != "cache expired" {
		t.Errorf("Expected 'cache expired' error, got %v", err)
	}
}

// TestMemCacheEviction tests the eviction policy when the cache reaches its max size
func TestMemCacheEviction(t *testing.T) {
	cache := NewCache(WithEvictionPolicy("Oldest"))

	// Fill the cache to its max capacity
	for i := 0; i < MaxEntries; i++ {
		err := cache.Set("bucket", fmt.Sprintf("%d", i), []byte("value"))
		if err != nil {
			t.Fatalf("Set failed at iteration %d: %v", i, err)
		}
	}

	// Add one more item, triggering eviction
	err := cache.Set("bucket", "extra_key", []byte("value"))
	if err != nil {
		t.Fatalf("Set failed during eviction: %v", err)
	}

	// The first inserted item should be evicted
	_, err = cache.Get("bucket", "0")
	if err == nil || err.Error() != "not found" {
		t.Errorf("Expected 'not found' error, got %v", err)
	}
}

// TestMemCacheDelete tests the Delete method
func TestMemCacheDelete(t *testing.T) {
	cache := NewCache()

	_ = cache.Set("bucket", "key", []byte("value"))

	// Test delete operation
	err := cache.Delete("bucket", "key")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Try to get the deleted key
	_, err = cache.Get("bucket", "key")
	if err == nil || err.Error() != "not found" {
		t.Errorf("Expected 'not found' error after delete, got %v", err)
	}
}
