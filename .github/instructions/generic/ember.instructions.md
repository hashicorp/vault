---
applyTo: "vault/ui/**/*.{js,ts,hbs,scss}"
description: "Ember.js coding standards and best practices for Vault UI"
---

# Ember.js Best Practices for Vault UI

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
- Co-locate `.hbs` templates with their corresponding `.js` or `.ts` files
- Place new components in logical subdirectories within `vault/ui/app/components/`
- Remove unnecessary quotes around dynamic data attributes in templates

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
- Prefer `async`/`await` over `@task` from ember-concurrency unless you need specific features like cancellation
- Use Ember's `run.later` instead of `setTimeout` in tests for better runloop control
- Consider `requestAnimationFrame` or `@tracked` properties for reactive patterns instead of `setTimeout`

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
- Check truthiness of arrays directly instead of using `.length` property
- Use string interpolation `"prefix/{{value}}"` instead of `{{concat}}` helper  
- Remove unnecessary quotes around dynamic component arguments
- Use `Hds::Link::Inline` for vault/docs links instead of `<button>` elements
- Make `selected` attributes dynamic rather than static values
- Refactor conditionals to wrap content rather than entire elements when possible
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
```

## Content and Terminology
- Use sentence case for titles and component arguments
- End descriptive text with proper punctuation
- Use "KVv1" or "KVv2" for KV secret engines, or "KV version 1/2" when space allows
- Indicate "enterprise" in changelog entries for enterprise-only features

Examples:
```javascript
// Good: sentence case
@title="Upload user's profile"

// Bad: incorrect casing  
@title="Upload Users'sProfile"

// Changelog entries
"ui (enterprise): Add advanced secret filtering"  // enterprise features
"ui: Fix secret list pagination"                  // community features
```


