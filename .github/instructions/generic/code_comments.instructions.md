---
applyTo: '**/*.{js,ts,go}'
description: 'Guidelines for writing effective code comments'
---

# Code Comments Guidelines

## General Philosophy

Comments should explain **why** code exists, not **what** it does. The code itself should be self-documenting through clear naming and structure. Add comments only when they provide value beyond what the code itself communicates.

## When to Add Comments

### Do Comment When:

- **Explaining rationale**: Why a particular approach was chosen over alternatives

  ```javascript
  // Use binary search because the list is already sorted and can be very large (100k+ items)
  const index = binarySearch(sortedList, target);
  ```

- **Documenting non-obvious behavior**: Side effects, performance characteristics, or gotchas

  ```javascript
  /*
   * OpenAPI schemas can reference other schemas using JSON Schema $ref
   * (e.g., "$ref": "#/components/schemas/SomeSchema")
   */
  const resolveRef = (spec, ref) => { ... }
  ```

- **Explaining business logic**: Domain-specific rules or requirements

  ```go
  // Per PCI DSS compliance requirements, credit card data must be encrypted at rest
  encryptedData := encrypt(cardData)
  ```

- **Warning about edge cases**: Unexpected behavior in specific scenarios

  ```javascript
  /*
   * Note: This will return null if the operation is a plugin-based method
   * that hasn't been enabled yet
   */
  const operation = findOperation(spec, operationId);
  ```

- **Providing context for complex algorithms**: High-level explanation of what's happening
  ```go
  /*
   * Dijkstra's algorithm for shortest path - we use a priority queue
   * because the graph can have up to 10M nodes in production
   */
  pq := priorityQueue.New()
  ```

## When NOT to Add Comments

### Don't Comment When:

- **The code is self-explanatory**: Function and variable names clearly indicate purpose

  ```javascript
  // BAD: Parse CLI arguments
  const parseArgs = () => { ... }

  // GOOD: No comment needed - function name is clear
  const parseArgs = () => { ... }
  ```

- **Describing what a function does**: The function name should do this

  ```javascript
  // BAD: Convert camelCase to kebab-case
  const toKebabCase = (str) => { ... }

  // GOOD: No comment needed
  const toKebabCase = (str) => { ... }
  ```

- **Obvious operations**: Standard language constructs or library calls

  ```javascript
  // BAD: Loop through all items
  for (const item of items) { ... }

  // BAD: Create a new array
  const results = [];
  ```

- **Repeating variable names**: Comments that just restate identifiers
  ```go
  // BAD: userID is the user ID
  userID := getUserID()
  ```

## Improving Code Instead of Commenting

Before adding a comment, consider if you can make the code clearer instead:

### Extract to Named Functions

```javascript
// BAD:
// Check if user has admin permissions and is active
if (user.role === 'admin' && user.status === 'active' && !user.suspended) { ... }

// GOOD:
const canPerformAdminAction = (user) =>
  user.role === 'admin' && user.status === 'active' && !user.suspended;

if (canPerformAdminAction(user)) { ... }
```

### Use Descriptive Variable Names

```javascript
// BAD:
const d = 86400; // seconds in a day

// GOOD:
const SECONDS_PER_DAY = 86400;
```

### Break Down Complex Expressions

```javascript
// BAD:
// Calculate total with tax and discount
const total =
	price * quantity * (1 + taxRate) - price * quantity * discountRate;

// GOOD:
const subtotal = price * quantity;
const tax = subtotal * taxRate;
const discount = subtotal * discountRate;
const total = subtotal + tax - discount;
```

## Comment Style

### Keep Comments Concise

- Use brief, direct language
- Avoid overly formal or verbose explanations
- One or two sentences is usually enough

### Update Comments When Code Changes

- Outdated comments are worse than no comments
- Review and update comments during code changes
- Delete comments that are no longer relevant

### Use Proper Grammar

- Start with a capital letter
- Use complete sentences when helpful
- End with punctuation for multi-sentence comments

## Language-Specific Guidelines

### JavaScript/TypeScript

- Use JSDoc for public API documentation
- For implementation comments:
  - Use `/* */` for multi-line comments (easier to read, edit, and refactor)
  - Use `//` for single-line comments above the relevant code
  - Avoid inline comments - place comments on the line above instead

**Example - Block comments for context:**

```javascript
/*
 * Construct request type name to match the vault-client-typescript SDK conventions.
 * Example: SystemApi.mountsEnableSecretsEngine -> SystemApiMountsEnableSecretsEngineRequest
 */
const requestTypeName = constructTypeName(apiName, methodName);
```

**Example - Single-line above relevant code:**

```javascript
// Remove leading # from ref path
const parts = ref.split('/').slice(1);
```

**Avoid inline comments:**

```javascript
// BAD: Inline comment gets lost and makes line too long
const parts = ref.split('/').slice(1); // Remove leading #

// GOOD: Comment above is clearer and easier to read
// Remove leading # from ref path
const parts = ref.split('/').slice(1);
```

### Go

- Use complete sentences starting with the name of the element being described
- Package comments go before the package declaration
- Exported identifiers should have doc comments

```go
// Server handles HTTP requests for the application.
// It manages connection pooling and request routing.
type Server struct { ... }
```

## Examples

### Good Comment - Explains Why

```javascript
// Use pnpm instead of yarn because the project has a packageManager field
execSync(`pnpm prettier --write "${filePath}"`, { stdio: 'pipe' });
```

### Good Comment - Warns About Behavior

```javascript
// This regex will fail on Unicode emoji - use a proper parser for user-generated content
const pattern = /[a-zA-Z0-9]/g;
```

### Bad Comment - States the Obvious

```javascript
// Set the user name
user.name = newName;
```

### Bad Comment - Out of Date

```javascript
// TODO: Remove this after migration to v2 API (added 2019)
const legacyData = fetchFromV1API();
```

## Summary

- **Prefer** self-documenting code over comments
- **Explain** why, not what
- **Keep** comments concise and current
- **Remove** comments that add no value
- **Update** or delete comments when code changes
