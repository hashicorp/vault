---
applyTo: "vault/ui/**/*.{js,ts,hbs,scss}"
description: "Ember.js coding standards and best practices for Vault UI"
---

# Ember.js Best Practices for Vault UI

## 🚨 Important Review Points - Check First

### Title Case Violations (High Priority)
When reviewing ANY Ember template or component file, immediately check for:
- **HTML headings**: All `<h1>`, `<h2>`, `<h3>`, etc. must use sentence case (NOT title case)
- **Component arguments**: `@title`, `@text`, `@label` properties must use sentence case
- **Common violations**: "Quick Actions" → "Quick actions", "Secrets Engines" → "Secrets engines"

❌ **INCORRECT** (Title Case - Flag in Review):
```handlebars
<h2 class="title is-4">Quick Actions</h2>
<h3 class="title is-6">Secrets Engines</h3>
<Hds::Button @text="Create New Secret" />
```

✅ **CORRECT** (Sentence Case):
```handlebars
<h2 class="title is-4">Quick actions</h2>
<h3 class="title is-6">Secrets engines</h3>
<Hds::Button @text="Create new secret" />
```

---

This document provides Ember.js-specific coding standards for the HashiCorp Vault UI application.

## Project Context
The Vault UI is an Ember.js application that provides a web interface for HashiCorp Vault, a secrets management tool. The application uses Ember Octane patterns, Handlebars templates, and SCSS for styling.

## Repository Structure
- `vault/ui/app/components/` - Reusable UI components  
- `vault/ui/app/models/` - Ember Data models for API entities
- `vault/ui/app/routes/` - Route handlers and data loading
- `vault/ui/tests/` - Integration and unit tests
- `vault/ui/app/templates/` - Handlebars templates

## Framework and Tools
- Ember.js 4.x with Ember Octane patterns
- Ember Data for API communication  
- Handlebars for templating
- SCSS for styling with HashiCorp Design System components
- Ember CLI Mirage for test mocking
- QUnit for testing

## Modern Ember Patterns
- Use native JavaScript classes instead of `Route.extend`, `Model.extend`, or `Component.extend`
- Follow Ember Octane conventions and consult official Ember documentation
- Create reusable, maintainable components rather than one-off implementations

## Component Development
- Use `@tracked` only for component internal state that changes during the lifecycle
- Do not use `@tracked` on computed properties or arguments passed from parent components
- Match component class names to their file names for consistency
- Co-locate `.hbs` templates with their corresponding `.js` or `.ts` files within `vault/ui/app/components/`
- Place new components in logical subdirectories within `vault/ui/app/components/`
- Remove unnecessary quotes around dynamic data attributes in templates
- Prioritize reusability and maintainability when creating components - avoid overly complex or one-off implementations

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

## Model Definitions
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

## JavaScript/TypeScript in Ember Context
- Use Ember's `assert` for runtime checks instead of `console.error`
- Prefer `async`/`await` over `@task` from ember-concurrency unless you need specific features like cancellation or `task.isRunning` state
- Use Ember's `run.later` instead of `setTimeout` in tests for better runloop control
- Consider `requestAnimationFrame` or `@tracked` properties for reactive patterns instead of `setTimeout`
- Avoid `new Date()` - use `Date.UTC()` constructor with UTC methods like `getUTCFullYear()`, `getUTCMonth()` for consistent timezone handling
- Include documentation for new files using JSDoc (e.g., `/** @module ComponentName */`)
- Remove unused imports, variables, and functions before committing

## Testing Standards
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

## Handlebars Template Guidelines

### Title and Heading Requirements
- **ALL HTML headings** (`<h1>`, `<h2>`, `<h3>`, `<h4>`, `<h5>`, `<h6>`) MUST use sentence case
- **Component title arguments** (`@title`, `@text`, `@label`) MUST use sentence case
- **Watch for**: Common violations like "Quick Actions", "Secrets Engines", "Authentication Methods"

```handlebars
{{!-- CORRECT: Sentence case in headings --}}
<h2 class="title is-4">Quick actions</h2>
<h3 class="title is-marginless is-6" data-test-card-subtitle="secrets-engines">Secrets engines</h3>

{{!-- INCORRECT: Title case violations --}}
<h2 class="title is-4">Quick Actions</h2>
<h3 class="title is-marginless is-6" data-test-card-subtitle="secrets-engines">Secrets Engines</h3>

{{!-- CORRECT: Component arguments with sentence case --}}
<Hds::Button @text="Create new secret" />
<Hds::Alert @title="Operation completed" />

{{!-- INCORRECT: Component arguments with title case --}}
<Hds::Button @text="Create New Secret" />
<Hds::Alert @title="Operation Completed" />
```

### Other Template Best Practices
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

## SCSS Styling Guidelines
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

## Content and Terminology

### Title and Heading Case Rules
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

### Common Title Case Violations to Avoid
- "Quick Actions" → should be "Quick actions"
- "Secrets Engines" → should be "Secrets engines" 
- "Authentication Methods" → should be "Authentication methods"
- "Policy Rules" → should be "Policy rules"
- "Access Control" → should be "Access control"

```javascript
// Changelog entries
"ui (enterprise): Add advanced secret filtering"  // enterprise features
"ui: Fix secret list pagination"                  // community features

// KV engine naming
"KVv2" or "KV version 2"     // Correct
"KV v2"                      // Incorrect
```

## Changelog Guidelines
For files in the `changelog/` directory:
- **Enterprise features**: Use `ui (enterprise): descriptive text`  
- **Community features**: Use `ui: descriptive text`
- Always indicate enterprise-only features in the description

## Code Review Checklist
When reviewing Ember.js code changes, pay special attention to:

### 🚨 High Priority: Title and Heading Case
- **Important**: Check ALL HTML headings (`<h1>`, `<h2>`, `<h3>`, etc.) use sentence case
- **Important**: Verify component `@title`, `@label`, and `@text` arguments use sentence case
- **Look for**: Title case violations like "Quick Actions", "Secrets Engines", "Authentication Methods"
- **Flag**: Any heading or title that capitalizes words beyond the first word

### Template and Component Issues
- Missing `data-test-*` attributes for new UI elements
- Unnecessary quotes around dynamic handlebars expressions
- Use of `{{concat}}` instead of string interpolation
- Static `selected` attributes that should be dynamic
- Inline `style` attributes instead of CSS classes

### JavaScript/TypeScript Issues  
- Incorrect use of `@tracked` on component arguments
- Missing `@tracked` on internal component state
- Use of `console.error` instead of Ember's `assert`
- Poor test assertions using `assert.ok()` instead of `assert.true()`/`assert.false()`

### Model and Data Issues
- Multi-line `@attr` declarations for simple attributes
- Inconsistent attribute grouping
- Missing labels for form-related attributes


