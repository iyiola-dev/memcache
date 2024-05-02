package main

import (
	"container/list"
	"errors"
	"log"
	"sync"
	"time"
)

const MaxEntries = 256

type Options struct {
	TTL            time.Duration
	EvictionPolicy string
}

type Option func(*Options) error

// Cache structure
type Cache struct {
	cache   map[string]*list.Element
	list    *list.List
	lock    sync.Mutex
	options Options
}

// entry is the data structure stored in the linked list
type entry struct {
	bucket string
	key    string
	value  []byte
	time   time.Time
}

// NewCache creates a new Cache with default options
func NewCache(opts ...Option) *Cache {
	cacheOptions := Options{
		EvictionPolicy: "Oldest", // Default policy
	}
	for _, opt := range opts {
		opt(&cacheOptions)
	}

	cache := &Cache{
		cache:   make(map[string]*list.Element),
		list:    list.New(),
		options: cacheOptions,
	}

	//automatic delete if the cache has expired
	go func() {
		for e := cache.list.Front(); e != nil; e = e.Next() {
			go automaticDelete(cache, e)
		}
	}()

	return cache

}

// Set inserts a value into the cache.
func (c *Cache) Set(bucket, key string, value []byte, opts ...Option) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Apply options
	for _, opt := range opts {
		if err := opt(&c.options); err != nil {
			return err
		}
	}

	// check if the cache is full
	if c.list.Len() == MaxEntries {
		//evict the oldest item from the cache
		c.evict(c.options.EvictionPolicy)
	}

	fullKey := bucket + ":" + key
	el, ok := c.cache[fullKey]

	if ok {
		c.list.MoveToFront(el)

		//update the value
		el.Value.(*entry).value = value

		//update the time
		el.Value.(*entry).time = time.Now()

		return nil

	}

	ent := &entry{bucket: bucket, key: fullKey, value: value, time: time.Now()}
	element := c.list.PushFront(ent)
	c.cache[fullKey] = element
	return nil
}

// Get retrieves a value from the cache.
func (c *Cache) Get(bucket, key string, opts ...Option) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Apply options
	for _, opt := range opts {
		if err := opt(&c.options); err != nil {
			return nil, err
		}
	}

	fullKey := bucket + ":" + key

	if ele, hit := c.cache[fullKey]; hit {
		c.list.MoveToFront(ele)
		if c.options.TTL > 0 && time.Since(ele.Value.(*entry).time) > c.options.TTL {
			c.list.Remove(ele)
			delete(c.cache, fullKey)
			return nil, errors.New("cache expired")
		}
		return ele.Value.(*entry).value, nil
	}

	return nil, errors.New("not found")

}

// Delete removes an entry from the cache.
func (c *Cache) Delete(bucket, key string, opts ...Option) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Apply options
	for _, opt := range opts {
		if err := opt(&c.options); err != nil {
			return err
		}
	}

	fullKey := bucket + ":" + key

	if ele, hit := c.cache[fullKey]; hit {
		c.removeElement(ele)
		return nil
	}
	return errors.New("not found")
}

// evict removes the oldest item from the cache based on the eviction policy.
func (c *Cache) evict(policy string) {
	switch policy {
	case "Oldest":
		c.removeElement(c.list.Back())
	}

}

// removeElement handles the removal of an element from the cache.
func (c *Cache) removeElement(e *list.Element) {
	c.list.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}

// automaticDelete is a helper function that automatically deletes an entry if it has expired.
func automaticDelete(c *Cache, e *list.Element) {
	//run every minute
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			if time.Since(e.Value.(*entry).time) > c.options.TTL {
				c.removeElement(e)
			}
		}
	}
}

// ApplyOptions applies given option functions to the cache options.
func (c *Cache) ApplyOptions(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(&c.options); err != nil {
			return err
		}
	}
	return nil
}

// WithTTL sets the time-to-live for cache entries.
func WithTTL(ttl time.Duration) Option {
	return func(o *Options) error {
		o.TTL = ttl
		return nil
	}
}

// WithEvictionPolicy sets the eviction policy for the cache.
func WithEvictionPolicy(policy string) Option {
	return func(o *Options) error {
		o.EvictionPolicy = policy
		return nil
	}
}

func main() {
	cache := NewCache(WithTTL(10*time.Second), WithEvictionPolicy("Oldest"))
	err := cache.Set("bucket", "key", []byte("value"))
	if err != nil {
		log.Println(err)
		panic(err)
	}

	value, err := cache.Get("bucket", "key")

	if err != nil {
		log.Println(err)
		panic(err)
	}

	log.Println(string(value))

	//sleep for 11 seconds
	time.Sleep(11 * time.Second)

	_, err = cache.Get("bucket", "key")
	if err != nil {
		log.Println(err)
	}

	err = cache.Delete("bucket", "key")

	if err != nil {
		log.Println(err)

	}

}
