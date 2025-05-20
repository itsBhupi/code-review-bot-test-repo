# Documentation Standards

All code must be properly documented to ensure maintainability and knowledge sharing.

## Function Documentation

- Every exported function must have a comment
- Comment should start with the function name
- Describe what the function does, not how it does it
- Document parameters and return values
- Document any errors that might be returned

## Package Documentation

- Every package should have a package comment in one of its files
- Package comment should provide an overview of the package's purpose

## Example

```go
// ProcessUser transforms a user entity into a standardized format.
// It handles validation and normalization of user data.
//
// Parameters:
//   - user: The user entity to process
//
// Returns:
//   - The processed user entity
//   - error: If validation fails or processing cannot be completed
func ProcessUser(user *User) (*User, error) {
    // Implementation
}
```

## Type Documentation

- Document all exported types
- Explain the purpose of the type
- Document all fields if it's a struct
- Document any methods associated with the type

## Comments

- Comments should be complete sentences with proper punctuation
- Keep comments up-to-date when code changes
- Use `TODO:` or `FIXME:` prefixes for temporary comments