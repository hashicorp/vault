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
- **Required**: Add documentation for all new files (components, helpers, services)

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
- Only use `@task` from ember-concurrency when you need specific features like cancellation or `task.isRunning` state management
- Avoid `setTimeout` in favor of `requestAnimationFrame` for UI updates or proper async patterns for delays
- **Warning**: `setTimeout` is prone to testing issues and event loop management problems
- Handle promise rejections explicitly rather than relying on global handlers

## Date and Time Handling  
- **WARNING**: Avoid `new Date()` as it uses the browser's timezone
- Use `Date.UTC()` constructor instead of `new Date()` for consistent timezone handling
- Use UTC methods like `getUTCFullYear()`, `getUTCMonth()` for date manipulation to ensure dates are calculated consistently
- Consider using a date library like `date-fns` for complex date operations

## Error Handling and Logging
- Avoid `console.error` in production code - use proper logging libraries or framework-specific methods
- Create meaningful error messages that include context about what failed
- Use structured logging with consistent log levels (debug, info, warn, error)

## Dependency Management
- Pin exact versions for critical dependencies or use tilde (`~`) for patch updates only
- **WARNING**: Avoid caret (`^`) operator which allows minor version updates that may introduce breaking changes
- Use tilde (`~`) for regular dependencies, exact versions for security-critical packages
- Dependencies in `resolutions` block MUST be pinned (no `~` or `^`)
- Keep `package.json` changes minimal and focused on the specific feature or fix
- Always commit lock files (`yarn.lock`, `package-lock.json`) with dependency changes
- Ensure package.json changes are independent of other code changes (except lock files)

Example dependency specification:
```json
{
  "dependencies": {
    "lodash": "4.17.21",        // exact version for critical packages
    "express": "~4.18.0"        // tilde for patch updates only
  },
  "resolutions": {
    "minimist": "1.2.6"         // always exact for security fixes
  }
}
```
