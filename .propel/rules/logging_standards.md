# Logging Standards

Proper logging is essential for troubleshooting and monitoring. Follow these guidelines.

## Log Levels

- **ERROR**: Use for unrecoverable errors that require immediate attention
- **WARN**: Use for unexpected but recoverable situations
- **INFO**: Use for significant events in the normal operation flow
- **DEBUG**: Use for detailed troubleshooting information

## Context

- Always include request ID in log entries
- Add sufficient context (user ID, operation type, etc.)
- Never log sensitive information (passwords, tokens, PII)

## Structure

- Use structured logging (key-value pairs)
- Include timestamps in UTC
- Use consistent field names across services

## Example

```go
// POOR LOGGING: Insufficient context, unstructured
func poorLoggingExample() {
    log.Println("User operation failed")
}

// GOOD LOGGING: Structured with context
func goodLoggingExample(ctx context.Context, userID string, operation string) {
    logger.WithFields(log.Fields{
        "request_id": ctx.Value("request_id"),
        "user_id": userID,
        "operation": operation,
    }).Error("Operation failed due to insufficient permissions")
}
``` 