# Code Formatting Standards

All Go code must follow consistent formatting to ensure readability and maintainability.

## Indentation and Spacing

- Use **tabs** for indentation, not spaces
- Maintain consistent indentation levels (each new block adds one level)
- Leave one blank line between functions
- No trailing whitespace
- Maximum line length should be 100 characters

## Braces and Parentheses

- Opening braces should be on the same line as the statement
- Closing braces should align with the start of the opening statement
- No spaces inside parentheses

## Example

```go
// Correct formatting
func processItem(item string) error {
	if item == "" {
		return fmt.Errorf("empty item")
	}
	
	for i := 0; i < len(item); i++ {
		if item[i] == ' ' {
			continue
		}
		// Process character
	}
	
	return nil
}
```

## Tools

Always run `gofmt` or `goimports` on your code before committing.