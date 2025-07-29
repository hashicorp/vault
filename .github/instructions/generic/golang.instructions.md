---
applyTo: "**/*.go"
description: "Practical guidance for Golang programming"
---

# Go Programming Instructions

## Code Style and Formatting
- Use `gofmt` and `goimports` to format all Go code automatically
- Follow standard Go naming conventions:
  - camelCase for unexported identifiers
  - PascalCase for exported identifiers
  - Use MixedCaps or mixedCaps rather than underscores for multiword names
  - Words in names that are initialisms or acronyms (e.g. “URL” or “NATO”) have a consistent case. For example, “URL” should appear as “URL” or “url” (as in “urlPony”, or “URLPony”), never as “Url”. As an example: ServeHTTP not ServeHttp. For identifiers with multiple initialized “words”, use for example “xmlHTTPRequest” or “XMLHTTPRequest”.
- Keep line length reasonable (Go has no strict limit, but wrap long lines)
- Use meaningful variable and function names

## Package Organization
- Use lowercase, single-word package names, no hyphens or underscores
- Keep packages focused and cohesive
- Place main packages in `cmd/` directory for applications, unless the repository will only ever contain one executable, in which case a `main.go` at the root of the repository is ok
- Package name should be the base name of its source directory

## Naming Conventions
### Package Names
- Use short, concise, evocative names
- Prefer brevity since everyone using your package will type the name
- Examples: `bufio.Reader` not `bufio.BufReader`, `ring.New` not `ring.NewRing`

### Interface Names
- One-method interfaces use method name + "-er" suffix: `Reader`, `Writer`, `Formatter`
- Honor canonical method signatures: `Read`, `Write`, `Close`, `String`

### Getters and Setters
- Prefer exporting fields over adding getters and setters where possible
- When required, don't use "Get" prefix for getters. For example, if the field name is `owner`, the getter should be `Owner()`, and the setter `SetOwner()`

## Functions and Methods
- Use multiple return values to improve error handling
- Only use named result parameters if multiple parameters are of the same type
- Use defer for cleanup operations (closing files, unlocking mutexes)
  - Defer executes in LIFO order
  - Arguments to deferred functions are evaluated when defer executes

## Control Structures
- Use initialization statement in if/switch when appropriate:
  ```go
  if err := file.Chmod(0664); err != nil {
      return err
  }
  ```
- Use switch statements instead of repeated if statements for validation:
  ```go
  // Good:
  switch {
  case x < 0: return fmt.Errorf("negative value")
  case x > 10: return fmt.Errorf("value too large")
  }
  
  // Bad:
  if x < 0 { return fmt.Errorf("negative value") }
  if x > 10 { return fmt.Errorf("value too large") }
  ```
- Omit unnecessary else statements when if ends with return/break/continue
- Use type switches for interface type checking
- Use range for iteration over arrays, slices, strings, maps, channels
   ```go
   for i, v := range slice { /* ... */ }  // Good
   for i := 0; i < len(slice); i++ { /* ... */ }  // Bad
   ```

## Error Handling and Design
- Always handle errors explicitly
- Use descriptive error messages that include context
- Wrap errors with additional context when appropriate
- Return errors as the last return value
- Check error returns - they're provided for a reason
- Use the "comma ok" idiom for optional error checking:
  ```go
  if value, ok := someOperation(); !ok {
      // handle error case
  }
  ```
- Implement the error interface, `Is(error) bool` and `As(any) bool` for custom error types
- Provide detailed error context including operation and file paths
- Use panic only for truly exceptional cases or programming errors
- Use recover to handle panics gracefully in server applications
- Design errors to be useful when printed far from their origin

## Data Structures
- Prefer zero-value-useful types in your designs
- Prefer slices over arrays for most use cases
  - Understand slice sharing and capacity
  - Use `append()` built-in for growing slices
  - Use copy() for copying slice elements
- Maps
  - Use "comma ok" idiom to test for presence: `value, ok := m[key]`
  - Use `delete(m, key)` to remove entries safely

## Interfaces and Composition
- Design small, focused interfaces (often single-method)
- Accept interfaces, return concrete types
- Use empty interface `interface{}` or any sparingly
- Prefer composition over large interfaces
- Use type assertions and type switches for interface conversions
- Embed types to promote methods and satisfy interfaces
  - Use embedding for has-a relationships, not is-a
  - Understand method promotion and shadowing rules

## Concurrency
- Don't communicate by sharing memory; share memory by communicating
- Use channels for goroutine communication and synchronization
- Use goroutines for concurrent execution: `go functionCall()`
- Use buffered channels as semaphores for limiting concurrency
- Use select for non-blocking channel operations
- Prefer structured concurrency patterns over ad-hoc goroutine creation
- Unbuffered channels provide synchronization
- Buffered channels can improve performance and provide semaphore behavior
- Close channels to signal completion to range loops

## Documentation
- Write clear godoc comments for exported functions and types
- Start comments with the name of the item being documented
- Use complete sentences in documentation
- Provide examples in documentation when helpful
- Document any unusual behavior or requirements

## Best Practices
- Keep functions small and focused
- Use context for cancellation and timeouts
- Avoid global variables when possible
- Use the blank identifier `_` to discard unused values
- Make zero values useful in your type designs
- Use init() functions for package initialization
- Initialize complex data with composite literals
- Profile before optimizing
- Use build constraints for platform-specific code
- Understand escape analysis and stack vs heap allocation
- Use sync.Pool for frequently allocated objects

## Common Idioms
- Use `fmt.Stringer` interface for custom string representations
- Implement `io.Reader` and `io.Writer` interfaces when appropriate
- Use functional options pattern for configuration
- Use builder pattern for complex object construction
- Prefer early returns to reduce nesting
- Initialize complex data with composite literals

### Type Embedding
- Embed types to promote methods and satisfy interfaces
- Use embedding for has-a relationships, not is-a
- Understand method promotion and shadowing rules

### Performance
- Understand escape analysis and stack vs heap allocation
- Use sync.Pool for frequently allocated objects
- Profile before optimizing
- Use build constraints for platform-specific code

## Common Idioms
- Use `fmt.Stringer` interface for custom string representations
- Implement `io.Reader` and `io.Writer` interfaces when appropriate
- Use functional options pattern for configuration
- Use builder pattern for complex object construction
- Prefer early returns to reduce nesting
