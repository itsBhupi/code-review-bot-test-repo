# Performance Best Practices

Follow these guidelines to ensure our application remains performant at scale.

## Database Interactions

- Keep transactions as short as possible
- Use indexes for frequently queried columns
- Paginate results for large datasets, default limit 50 items

## Memory Management

- Avoid unnecessary allocations in hot paths
- Use object pooling for frequently created/destroyed objects
- Consider using sync.Pool for temporary objects

## Concurrency

- Use goroutines judiciously, monitor creation rate
- Implement backpressure mechanisms for high-throughput services
- Always use timeouts for external service calls, maximum 30 seconds

## Example

```go
// NOT OPTIMIZED: Creates a new buffer for each request
func unoptimizedHandler(w http.ResponseWriter, r *http.Request) {
    buffer := make([]byte, 1024*1024) // 1MB buffer for each request
    // Use buffer...
}

// OPTIMIZED: Reuses buffers from a pool
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024*1024)
    },
}

func optimizedHandler(w http.ResponseWriter, r *http.Request) {
    buffer := bufferPool.Get().([]byte)
    defer bufferPool.Put(buffer)
    // Use buffer...
}
``` 