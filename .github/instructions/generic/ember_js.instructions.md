---
applyTo: "vault/ui/**/*.{js,ts}"
description: "HashiCorp Ember.js UI JavaScript and TypeScript coding standards"
---

# HashiCorp Ember.js JavaScript & TypeScript Guidelines

This document provides JavaScript and TypeScript coding standards for HashiCorp Ember.js UI applications.

> **Note**: For general project context, framework information, and repository structure, see `ember_general.instructions.md`.

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
- **File Naming**: All new files should use `.ts` extension instead of `.js`
- **Migration Strategy**: When modifying existing `.js` files, evaluate the effort required to convert to `.ts` and prioritize conversion if time permits
- Use explicit types instead of `any` - prefer `unknown` for truly dynamic content
- Define interfaces for object shapes and function signatures
- Use type guards and discriminated unions for runtime type checking
- Enable strict mode in TypeScript configuration

## Modern Ember Patterns
- Replace `Route.extend`, `Model.extend`, or `Component.extend` with native JavaScript classes
- Use Ember Octane conventions: tracked properties, decorators, and native class syntax
- Create reusable components rather than one-off implementations
- Co-locate component templates (`.hbs`) with their TypeScript files (`.ts` preferred over `.js`)
- Prioritize reusability and maintainability when creating components - avoid overly complex or one-off implementations

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
- Use Ember's `assert` for runtime checks instead of `console.error`
- Avoid `console.error` in production code - use proper logging libraries or framework-specific methods
- Create meaningful error messages that include context about what failed
- Use structured logging with consistent log levels (debug, info, warn, error)

---

# Component Development (JavaScript/TypeScript)

## Component Architecture
- Use `@tracked` only for internal component state that changes over time
- Never use `@tracked` on component arguments (properties passed from parent components)
- Component class names must match their file names exactly
- Place `.hbs` templates in the same directory as their `.ts` files (or `.js` for legacy files) within `ui/app/components/`
- Organize new components in logical subdirectories by feature or domain

Examples:
```javascript
// Good: tracked for internal state
@tracked isExpanded = false;

// Bad: tracked on argument
@tracked @secret; // @secret is an argument, not internal state
```

---

# Model Definitions

## Ember Data Models
- Use single-line syntax for `@attr` declarations with simple options
- Avoid extra blank lines between consecutive single-line `@attr` declarations
- Group related attributes together logically

Example:
```javascript
// Good: compact single-line format
@attr('string', { label: 'Client ID' }) clientId;
@attr('string', { label: 'Client Secret' }) clientSecret;
@attr('boolean', { defaultValue: false }) isEnabled;

// Bad: unnecessary multi-line for simple attributes
@attr('string', {
  label: 'Client ID',
}) clientId;

@attr('string', { label: 'Client Secret' }) clientSecret;
```
