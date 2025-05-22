# Security Guidelines

These guidelines help prevent common security vulnerabilities in our codebase.

## Input Validation

- All user input must be validated before use
- Use whitelisting (allow-list) approach rather than blacklisting
- Never trust client-side validation alone

## Authentication & Authorization

- Use our central auth library for all authentication flows
- Implement proper role-based access control for all endpoints
- Never store passwords in plain text, always use our HashPassword utility

## Data Protection

- Use prepared statements for all SQL queries
- Apply proper escaping for all HTML output
- Encrypt all sensitive data at rest using our encryption library

## Example

```go
// VULNERABLE: Direct use of user input in SQL query
func badQueryExample(userInput string) {
    db.Query("SELECT * FROM users WHERE name = '" + userInput + "'")
}

// SECURE: Using prepared statements
func secureQueryExample(userInput string) {
    db.Query("SELECT * FROM users WHERE name = ?", userInput)
}
``` 