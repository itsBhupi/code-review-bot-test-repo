# Error Handling Guidelines

Proper error handling is crucial for building robust and maintainable Go applications. Follow these guidelines to ensure consistent and effective error management.

## Core Principles

- **Always handle errors explicitly** - Never ignore errors with `_`
- **Fail fast** - Return errors immediately when they occur
- **Provide context** - Add meaningful information to errors
- **Use structured logging** - Log errors with relevant context fields
- **Don't panic in libraries** - Reserve panic for truly exceptional cases

## Error Handling Patterns

### 1. Basic Error Handling

```go
// Good: Always check and handle errors
result, err := someFunction()
if err != nil {
    log.WithFields(log.Fields{
        "operation": "someFunction",
        "error": err.Error(),
    }).Error("Operation failed")
    return fmt.Errorf("failed to process: %w", err)
}

// Bad: Ignoring errors
result, _ := someFunction() // Never do this
```

### 2. Error Wrapping

```go
// Good: Wrap errors with context
func ProcessUser(userID string) error {
    user, err := fetchUser(userID)
    if err != nil {
        return fmt.Errorf("failed to fetch user %s: %w", userID, err)
    }
    
    if err := validateUser(user); err != nil {
        return fmt.Errorf("user validation failed for %s: %w", userID, err)
    }
    
    return nil
}
```

### 3. Custom Error Types

```go
// Define custom error types for specific error conditions
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// Usage
func ValidateEmail(email string) error {
    if !isValidEmail(email) {
        return &ValidationError{
            Field:   "email",
            Value:   email,
            Message: "invalid email format",
        }
    }
    return nil
}
```

## HTTP Error Handling

### 1. Structured Error Responses

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code"`
    Details string `json:"details,omitempty"`
}

func handleError(w http.ResponseWriter, err error, statusCode int) {
    log.WithFields(log.Fields{
        "error": err.Error(),
        "status_code": statusCode,
    }).Error("HTTP request failed")
    
    response := ErrorResponse{
        Error: http.StatusText(statusCode),
        Code:  fmt.Sprintf("ERR_%d", statusCode),
        Details: err.Error(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(response)
}
```

### 2. Error Middleware

```go
func ErrorMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.WithFields(log.Fields{
                    "panic": err,
                    "path": r.URL.Path,
                    "method": r.Method,
                }).Error("Panic recovered")
                
                handleError(w, fmt.Errorf("internal server error"), http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}
```

## Database Error Handling

### 1. Connection Errors

```go
func ConnectDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database connection: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.WithFields(log.Fields{
        "database": "postgres",
        "status": "connected",
    }).Info("Database connection established")
    
    return db, nil
}
```

### 2. Query Errors

```go
func GetUser(db *sql.DB, userID int) (*User, error) {
    user := &User{}
    
    query := "SELECT id, name, email FROM users WHERE id = $1"
    err := db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found: %d", userID)
        }
        return nil, fmt.Errorf("failed to query user %d: %w", userID, err)
    }
    
    return user, nil
}
```

## Error Logging Standards

### 1. Structured Error Logging

```go
// Good: Use structured logging with context
log.WithFields(log.Fields{
    "operation": "user_creation",
    "user_id": userID,
    "error": err.Error(),
    "timestamp": time.Now(),
}).Error("Failed to create user")

// Bad: String concatenation
log.Println("Error creating user " + userID + ": " + err.Error())
```

### 2. Error Levels

- **ERROR**: Use for errors that affect functionality
- **WARN**: Use for recoverable errors or unexpected conditions
- **INFO**: Use for important operational information
- **DEBUG**: Use for detailed diagnostic information

## What NOT to Do

### 1. Don't Ignore Errors

```go
// Bad: Ignoring errors
result, _ := riskyOperation()

// Bad: Empty error handling
if err != nil {
    // TODO: handle error
}
```

### 2. Don't Use log.Fatal in Libraries

```go
// Bad: Using log.Fatal in library code
func ProcessData(data []byte) {
    if len(data) == 0 {
        log.Fatal("data is empty") // Don't do this
    }
}

// Good: Return error instead
func ProcessData(data []byte) error {
    if len(data) == 0 {
        return fmt.Errorf("data cannot be empty")
    }
    return nil
}
```

### 3. Don't Log and Return Errors

```go
// Bad: Both logging and returning the same error
func ProcessUser(userID string) error {
    user, err := fetchUser(userID)
    if err != nil {
        log.Error("Failed to fetch user: ", err) // Don't log here
        return err // And also return
    }
    return nil
}

// Good: Either log OR return, not both
func ProcessUser(userID string) error {
    user, err := fetchUser(userID)
    if err != nil {
        return fmt.Errorf("failed to fetch user %s: %w", userID, err)
    }
    return nil
}
```

## Error Recovery

### 1. Graceful Degradation

```go
func GetUserWithFallback(userID string) (*User, error) {
    // Try primary source
    user, err := getUserFromPrimary(userID)
    if err == nil {
        return user, nil
    }
    
    log.WithFields(log.Fields{
        "user_id": userID,
        "error": err.Error(),
    }).Warn("Primary user source failed, trying fallback")
    
    // Try fallback source
    user, err = getUserFromFallback(userID)
    if err != nil {
        return nil, fmt.Errorf("both primary and fallback failed for user %s: %w", userID, err)
    }
    
    return user, nil
}
```

### 2. Retry Logic

```go
func RetryOperation(operation func() error, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        if err := operation(); err != nil {
            lastErr = err
            
            log.WithFields(log.Fields{
                "attempt": i + 1,
                "max_retries": maxRetries,
                "error": err.Error(),
            }).Warn("Operation failed, retrying")
            
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        return nil
    }
    
    return fmt.Errorf("operation failed after %d retries: %w", maxRetries, lastErr)
}
```

## Testing Error Conditions

Always test error conditions in your unit tests:

```go
func TestProcessUser_InvalidInput(t *testing.T) {
    err := ProcessUser("")
    if err == nil {
        t.Error("Expected error for empty user ID")
    }
    
    if !strings.Contains(err.Error(), "user ID cannot be empty") {
        t.Errorf("Expected specific error message, got: %v", err)
    }
}
```
