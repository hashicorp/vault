---
applyTo: "vault/ui/**/*.{js,ts,hbs,scss}"
description: "Vault UI coding standards and best practices - Ember.js, JavaScript, TypeScript, Handlebars, SCSS"
---

# Vault UI Best Practices

This document provides comprehensive coding standards for the HashiCorp Vault UI application, covering Ember.js, JavaScript, TypeScript, Handlebars templates, and SCSS.

## Project Context
The Vault UI is an Ember.js web application that provides a graphical interface for HashiCorp Vault, a secrets management and data protection platform. The application allows users to manage secrets, authentication methods, policies, and audit logs through a browser-based interface. This is an enterprise-grade application used by DevOps teams, security engineers, and system administrators to securely store and access sensitive data.

The application uses modern Ember Octane patterns, Handlebars templates, and SCSS for styling with the HashiCorp Design System components.

## Repository Structure
- `vault/ui/app/components/` - Reusable UI components and their templates
- `vault/ui/app/models/` - Ember Data models representing API entities (secrets, policies, auth methods)
- `vault/ui/app/routes/` - Route handlers for URL endpoints and data loading logic
- `vault/ui/app/templates/` - Page-level Handlebars templates 
- `vault/ui/app/services/` - Ember services for shared functionality and state management
- `vault/ui/app/helpers/` - Template helper functions for data formatting and logic
- `vault/ui/tests/` - Integration, unit, and acceptance tests
- `vault/ui/app/styles/` - SCSS stylesheets and component-specific styles
- `vault/ui/config/` - Ember CLI configuration and environment settings
- `vault/ui/mirage/` - Mock server configuration for development and testing

## Framework and Tools
- **Frontend Framework**: Ember.js 4.x with Ember Octane patterns and decorators
- **Data Layer**: Ember Data for API communication with Vault's REST API
- **Templating**: Handlebars templates with Ember's component architecture
- **Styling**: SCSS with HashiCorp Design System (HDS) components and Bulma CSS framework
- **Testing**: QUnit for unit/integration tests, Ember CLI Mirage for API mocking
- **Build System**: Ember CLI with Webpack for bundling and asset management
- **Development**: ESLint for code linting, Prettier for code formatting

## File Naming and TypeScript Adoption
- **New Files**: All new JavaScript files must use `.ts` extension for TypeScript
- **Legacy File Migration**: When editing existing `.js` files, assess the effort required to convert to `.ts` and prioritize conversion if time permits
- **Component Structure**: Components should have matching `.ts` and `.hbs` files in the same directory
- **Gradual Migration**: The codebase is transitioning from JavaScript to TypeScript - contribute to this effort when possible

## Build and Deployment Structure
- **Development Build**: `ember serve` creates development server with live reload
- **Production Build**: `ember build --environment=production` generates optimized assets in `dist/`
- **Asset Output**: Built files are served directly by the Vault binary's embedded web server
- **Integration**: UI is compiled and embedded into the main Vault binary during release builds
- **Testing**: `ember test` runs the full test suite with QUnit in headless Chrome

---

# JavaScript & TypeScript Guidelines

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
- Use Ember's `run.later` instead of `setTimeout` in tests for better runloop control
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

# Component Development

## Component Architecture
- Use `@tracked` only for internal component state that changes over time
- Never use `@tracked` on component arguments (properties passed from parent components)
- Component class names must match their file names exactly
- Place `.hbs` templates in the same directory as their `.ts` files (or `.js` for legacy files) within `vault/ui/app/components/`
- Organize new components in logical subdirectories by feature or domain
- Remove quotes around dynamic data attributes: `data-test-id={{value}}` not `data-test-id="{{value}}"`

Examples:
```javascript
// Good: tracked for internal state
@tracked isExpanded = false;

// Bad: tracked on argument
@tracked @secret; // @secret is an argument, not internal state
```

```handlebars
{{!-- Good: no quotes around dynamic values --}}
<div data-test-namespace-link={{option.label}}>

{{!-- Bad: unnecessary quotes --}}
<div data-test-namespace-link="{{option.label}}">
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

---

# Testing Standards

## Test Quality
- Use `assert.true()` or `assert.false()` instead of `assert.ok()` for boolean checks
- Provide descriptive assertion messages that explain what is being verified
- Use `data-test-*` selectors for DOM interactions to decouple from styling
- Place `this.server` setup at the top of test blocks when using ember-cli-mirage
- Interact with elements directly rather than using Page Object patterns

Example test assertions:
```javascript
// Good: specific assertion with descriptive message
assert.true(component.isVisible, 'Component should be visible after clicking toggle');

// Bad: vague assertion without context
assert.ok(component.isVisible);
```

---

# Handlebars Template Guidelines

## Title and Heading Requirements
- **ALL HTML headings** (`<h1>`, `<h2>`, `<h3>`, `<h4>`, `<h5>`, `<h6>`) MUST use sentence case
- **Component title arguments** (`@title`, `@text`, `@label`) MUST use sentence case

```handlebars
{{!-- CORRECT: Sentence case in headings --}}
<h2 class="title is-4">Quick actions</h2>

{{!-- INCORRECT: Title case violations --}}
<h2 class="title is-4">Quick Actions</h2>

{{!-- CORRECT: Component arguments with sentence case --}}
<Hds::Button @text="Create new secret" />
<Hds::Alert @title="Operation completed" />

{{!-- INCORRECT: Component arguments with title case --}}
<Hds::Button @text="Create New Secret" />
<Hds::Alert @title="Operation Completed" />
```

## Template Best Practices
- Check truthiness of arrays directly instead of using `.length` property
- Use string interpolation `"prefix/{{value}}"` instead of `{{concat}}` helper  
- Remove unnecessary quotes around dynamic component arguments
- Use `Hds::Link::Inline` for `vault/docs` links instead of `<button>` elements
- Make `selected` attributes dynamic rather than static values - warn if static values are used
- Refactor conditionals to wrap content rather than entire elements when possible
- Avoid inline `style` attributes - define CSS classes in `.scss` files instead
- Place `data-test-*` selectors as the last attribute on elements

Examples:
```handlebars
{{!-- Good: direct array check --}}
{{#if this.model.allowed_roles}}

{{!-- Bad: unnecessary .length check --}}
{{#if (gt this.model.allowed_roles.length 0)}}

{{!-- Good: string interpolation --}}
@secret="role/{{@model.id}}"

{{!-- Bad: concat helper --}}
@secret={{concat "role/" @model.id}}

{{!-- Good: conditional content, not element --}}
<PH.Title>{{if this.version.isEnterprise "Enterprise" "Community"}} features</PH.Title>

{{!-- Bad: conditional wrapping entire element --}}
{{#if this.version.isEnterprise}}
  <PH.Title>Enterprise features</PH.Title>
{{else}}
  <PH.Title>Community features</PH.Title>
{{/if}}

{{!-- Good: CSS classes instead of inline styles --}}
<Hds::Button @text="Save" class="custom-button" data-test-save />

{{!-- Bad: inline style attribute --}}
<Hds::Button @text="Save" style="margin-top: 10px;" data-test-save />

{{!-- Good: data-test selector at the end --}}
<Hds::Button @text="Save" @icon="loading" disabled={{this.isLoading}} data-test-save />

{{!-- Bad: data-test selector not at the end --}}
<Hds::Button data-test-save @text="Save" @icon="loading" disabled={{this.isLoading}} />
```

---

# SCSS Styling Guidelines

## CSS Best Practices
- Avoid `z-index` - manage stacking order through DOM structure instead
- Avoid `!important` - use more specific selectors for overrides
- Define CSS classes in `.scss` files rather than using inline `style` attributes

Examples:
```scss
// Good: specific selector for overrides
.namespace-picker .button {
  color: red;
}

// Bad: using !important
.button {
  color: red !important;
}

// Bad: using z-index for stacking
.modal {
  position: absolute;
  z-index: 10;
}

// Good: manage stacking through DOM order in template
// Place modal element after background overlay in template
```

---

# Content and Terminology

## Title and Heading Case Rules
- **USE SENTENCE CASE**: All HTML headings (`<h1>`, `<h2>`, `<h3>`, etc.) should use sentence case (only first letter capitalized)
- **NO TITLE CASE**: Avoid title case where every major word is capitalized
- **Component arguments**: Use sentence case for `@title`, `@label`, and similar text properties
- End descriptive text with proper punctuation
- Use "KVv1" or "KVv2" for KV secret engines, or "KV version 1/2" when space allows (not "KV v2")
- Indicate "enterprise" in changelog entries for enterprise-only features
- Follow proper grammar rules including ending sentences with periods

Examples:
```handlebars
{{!-- CORRECT: Sentence case in HTML headings --}}
<h2 class="title is-4">Quick actions</h2>
<h3 class="title is-marginless is-6">Secrets engines</h3>
<h1>Authentication methods</h1>

{{!-- INCORRECT: Title case in HTML headings --}}
<h2 class="title is-4">Quick Actions</h2>
<h3 class="title is-marginless is-6">Secrets Engines</h3>
<h1>Authentication Methods</h1>
```

```javascript
// CORRECT: Sentence case in component arguments
@title="Upload user's profile"
@label="Secret path"
@placeholder="Enter mount path"

// INCORRECT: Title case or inconsistent casing
@title="Upload User's Profile"
@label="Secret Path"
@placeholder="Enter Mount Path"
```

```handlebars
{{!-- CORRECT: Sentence case in data-test attributes when they contain readable text --}}
data-test-card-subtitle="secrets-engines"

{{!-- Component usage with proper casing --}}
<Hds::Button @text="Create secret" />
<Hds::Alert @message="Operation completed successfully" />
```

```javascript
// KV engine naming
"KVv2" or "KV version 2"     // Correct
"KV v2"                      // Incorrect
```

---

# Changelog Guidelines

For files in the `changelog/` directory:
- **Enterprise features**: Use `ui (enterprise): descriptive text`  
- **Community features**: Use `ui: descriptive text`
- Always indicate enterprise-only features in the description

```javascript
// Changelog entries
"ui (enterprise): Add advanced secret filtering"  // enterprise features
"ui: Fix secret list pagination"                  // community features
```

---

# Dependency Management

## Package.json Guidelines
- Pin exact versions for critical dependencies or use tilde (`~`) for patch updates only
- **WARNING**: Avoid caret (`^`) operator which allows minor version updates that may introduce breaking changes
- Use tilde (`~`) for regular dependencies, exact versions for security-critical packages
- Dependencies in `resolutions` block MUST be pinned (no `~` or `^`)
- Keep `package.json` changes minimal and focused on the specific feature or fix
- Always commit lock files (`yarn.lock`, `package-lock.json`) with dependency changes
- Ensure package.json changes are independent of other code changes (except lock files)
- **Reminder**: Consider backporting package.json changes to LTS versions of the Vault binary when possible

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

---

# Code Review Checklist

## Title and Heading Case
- Check ALL HTML headings (`<h1>`, `<h2>`, `<h3>`, etc.) use sentence case
- Verify component `@title`, `@label`, and `@text` arguments use sentence case
- **Look for**: Title case violations like "Quick Actions", "Secrets Engines", "Authentication Methods"
- **Flag**: Any heading or title that capitalizes words beyond the first word

**Quick Reference Examples:**
```handlebars
{{!-- CORRECT: Sentence case --}}
<h2 class="title is-4">Quick actions</h2>
<Hds::Button @text="Create new secret" />

{{!-- INCORRECT: Title case - Flag in Review --}}
<h2 class="title is-4">Quick Actions</h2>
<Hds::Button @text="Create New Secret" />
```

## Template and Component Issues
- Missing `data-test-*` attributes for new UI elements
- Unnecessary quotes around dynamic handlebars expressions
- Use of `{{concat}}` instead of string interpolation
- Static `selected` attributes that should be dynamic
- Inline `style` attributes instead of CSS classes

## JavaScript/TypeScript Issues  
- Incorrect use of `@tracked` on component arguments
- Missing `@tracked` on internal component state
- Use of `console.error` instead of Ember's `assert`
- Poor test assertions using `assert.ok()` instead of `assert.true()`/`assert.false()`
- Missing documentation for new files
- Use of `any` type in TypeScript
- Incorrect async patterns (`setTimeout` instead of proper alternatives)
- Use of `new Date()` instead of UTC methods

## Model and Data Issues
- Multi-line `@attr` declarations for simple attributes
- Inconsistent attribute grouping
- Missing labels for form-related attributes

## Debugging Tips
When debugging Ember templates:
- Use `{{debugger}}` inside templates to inspect values in the browser console
- In the console during debugging:
  - Use `get('property.name')` to inspect nested properties
  - Use `context` to explore the current template context
- Example usage:
  ```handlebars
  {{#each this.items as |item|}}
    {{debugger}}
  {{/each}}
  ```
