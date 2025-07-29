---
applyTo: "**/*.{js,ts,jsx,tsx}"
description: "JavaScript and TypeScript coding standards and best practices"
---

# JavaScript/TypeScript Best Practices

This document provides coding standards for JavaScript and TypeScript development, focusing on code quality, maintainability, and consistency.

## Project Context
This applies to JavaScript and TypeScript code across web applications, Node.js services, and utility scripts. Code should prioritize readability, type safety, and maintainability.

## Documentation Standards
- Include JSDoc comments for all public functions, classes, and modules
- Use `/** @module ModuleName */` for modules and `/** description */` for functions
- Document function parameters, return types, and any side effects
- Keep comments concise and focused on the "why" rather than the "what"

## Code Quality Standards
- Remove all unused imports, variables, and functions before committing
- Place comments directly above the code they describe, not inline or below
- Update comments when code changes to maintain accuracy

## TypeScript Guidelines
- Use explicit types instead of `any` - prefer `unknown` for truly dynamic content
- Define interfaces for object shapes and function signatures
- Use type guards and discriminated unions for runtime type checking
- Enable strict mode in TypeScript configuration

## Asynchronous Programming
- Use `async`/`await` with proper error handling in `try`/`catch` blocks
- Only use specialized async libraries when you need specific features like cancellation or state management
- Avoid `setTimeout` in favor of `requestAnimationFrame` for UI updates or proper async patterns for delays
- Handle promise rejections explicitly rather than relying on global handlers

## Date and Time Handling  
- Use `Date.UTC()` constructor instead of `new Date()` for consistent timezone handling
- Use UTC methods like `getUTCFullYear()`, `getUTCMonth()` for date manipulation
- Consider using a date library like `date-fns` for complex date operations

## Error Handling and Logging
- Avoid `console.error` in production code - use proper logging libraries or framework-specific methods
- Create meaningful error messages that include context about what failed
- Use structured logging with consistent log levels (debug, info, warn, error)

## Dependency Management
- Pin exact versions for critical dependencies or use tilde (`~`) for patch updates
- Avoid caret (`^`) operator which allows minor version updates that may introduce breaking changes
- Keep `package.json` changes minimal and focused on the specific feature or fix
- Always commit lock files (`yarn.lock`, `package-lock.json`) with dependency changes

Example dependency specification:
```json
{
  "dependencies": {
    "lodash": "4.17.21",        // exact version
    "express": "~4.18.0"        // patch updates only
  },
  "resolutions": {
    "minimist": "1.2.6"         // always exact for security fixes
  }
}
```
