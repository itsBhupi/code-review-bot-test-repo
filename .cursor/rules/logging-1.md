# Logging Standards

Proper logging is essential for debugging and monitoring applications in production.

## Log Levels

- **ERROR**: Use for errors that affect functionality
- **WARN**: Use for unexpected conditions that don't break functionality
- **INFO**: Use for important state changes and major operations
- **DEBUG**: Use for detailed diagnostic information

## Guidelines

- Use structured logging with fields instead of string concatenation
- Include relevant context in logs (request ID, user ID, etc.)
- Don't log sensitive information (passwords, tokens, etc.)
- Use consistent format for similar log messages

## Example

```go
// Bad logging
log.Println("User " + username + " logged in at " + time.Now().String())

// Good logging
log.WithFields(log.Fields{
    "username": username,
    "timestamp": time.Now(),
    "request_id": requestID,
}).Info("User logged in")
```

## Error Logs

All error logs should:
1. Include the full error message
2. Provide context about what operation was being performed
3. Include relevant variables that would help debugging