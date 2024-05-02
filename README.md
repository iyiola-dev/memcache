# MemCache

MemCache is a simple in-memory cache implementation written in Go. It provides basic caching functionalities such as storing, retrieving, and deleting key-value pairs with optional time-to-live (TTL) support and eviction policies.

## Features

- **In-Memory Storage**: Data is stored in memory, providing fast read and write access.
- **Key-Value Storage**: Data is stored as key-value pairs, allowing efficient retrieval based on keys.
- **Time-to-Live (TTL)**: Supports setting a TTL for cache entries, after which they are automatically removed.
- **Eviction Policies**: Supports eviction policies to manage cache size when reaching capacity.
- **Concurrency-Safe**: Designed to support concurrent read and write operations safely.

## Creating a Cache Instance

To create a new MemCache instance, use the `NewCache` function along with optional configuration options:

```go
cache := NewCache(WithTTL(10*time.Second), WithEvictionPolicy("Oldest"))
```

## Setting Values

Insert key-value pairs into the cache using the `Set` method:

```go
err := cache.Set("bucket", "key", []byte("value"))
if err != nil {
    // Handle error
}
```

## Retrieving Values

To retrieve a value from the cache, use the `Get` method:

```go
value, err := cache.Get("bucket", "key")
if err != nil {
    // Handle error
}
```

## Deleting Values

Delete an entry from the cache using the `Delete` method:

```go
err := cache.Delete("bucket", "key")
if err != nil {
    // Handle error
}
```


## Configuration Options

MemCache supports configuration options such as TTL and eviction policy. You can apply these options during cache creation or update them later:

```go
// Apply options during cache creation
cache := NewCache(WithTTL(10*time.Second), WithEvictionPolicy("Oldest"))

// Update options later
cache.ApplyOptions(WithTTL(20*time.Second))
```

##### Available Eviction Policy: Determines how Memcached handles cache capacity limitations. The `WithEvictionPolicy` option accepts a string argument specifying the desired policy:

- "Oldest" (default): Evicts the oldest item.




