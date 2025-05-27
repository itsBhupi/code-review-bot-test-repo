# Code Style Guidelines

Our codebase follows these style conventions to ensure consistency and readability.

## Naming Conventions

- **Variables**: Use camelCase for variables
- **Functions**: Use camelCase for functions
- **Types/Structs**: Use PascalCase for types and structs
- **Constants**: Use UPPER_SNAKE_CASE for constants

## Formatting

- Use tabs for indentation, not spaces
- Maximum line length is 100 characters
- Always use braces for control statements, even for single-line blocks

## Documentation

- All exported functions, types, and variables must have comments
- Comments should be complete sentences with proper punctuation
- Use godoc style for Go code documentation

## Example

```go
// BadExample shows an incorrectly formatted function
func badExample() {
  if (condition) return // No braces, incorrect indentation
}

// GoodExample demonstrates proper formatting and documentation.
// It follows all our style guidelines.
func GoodExample() {
	if condition {
		return
	}
}
```

<!-- @propel:reference=error-handling.md -->
