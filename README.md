# MemCache

MemCache is a simple in-memory cache implementation written in Go. It provides basic caching functionalities such as storing, retrieving, and deleting key-value pairs with optional time-to-live (TTL) support and eviction policies.

## Features

- **In-Memory Storage**: Data is stored in memory, providing fast read and write access.
- **Key-Value Storage**: Data is stored as key-value pairs, allowing efficient retrieval based on keys.
- **Time-to-Live (TTL)**: Supports setting a TTL for cache entries, after which they are automatically removed.
- **Eviction Policies**: Supports eviction policies to manage cache size when reaching capacity.
- **Concurrency-Safe**: Designed to support concurrent read and write operations safely.

## Installation

To use MemCache in your Go project, simply import it as a package:

```go
import "your/module/path/memcache"
