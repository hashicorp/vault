---
applyTo: "vault/ui/**/*.hbs"
description: "HashiCorp Ember.js UI Handlebars template coding standards"
---

# HashiCorp Ember.js Handlebars Template Guidelines

This document provides Handlebars template coding standards for HashiCorp Ember.js UI applications.

> **Note**: For general project context, framework information, and repository structure, see `ember_general.instructions.md`.

## Template Best Practices
- Check truthiness of arrays directly instead of using `.length` property
- Use string interpolation `"prefix/{{value}}"` instead of `{{concat}}` helper  
- Remove unnecessary quotes around dynamic component arguments
- Use `Hds::Link::Inline` for external documentation links instead of `<button>` elements
- Make `selected` attributes dynamic rather than static values - warn if static values are used
- Refactor conditionals to wrap content rather than entire elements when possible
- Avoid inline `style` attributes and `{{style ...}}` helpers - define CSS classes in `.scss` files instead
- Place `data-test-*` selectors as the last attribute on elements
- Remove quotes around dynamic data attributes: `data-test-id={{value}}` not `data-test-id="{{value}}"`

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

{{!-- Bad: style helper --}}
<Hds::Button @text="Save" style={{style margin-top="10px"}} data-test-save />

{{!-- Good: data-test selector at the end --}}
<Hds::Button @text="Save" @icon="loading" disabled={{this.isLoading}} data-test-save />

{{!-- Bad: data-test selector not at the end --}}
<Hds::Button data-test-save @text="Save" @icon="loading" disabled={{this.isLoading}} />

{{!-- Good: no quotes around dynamic values --}}
<div data-test-namespace-link={{option.label}}>

{{!-- Bad: unnecessary quotes --}}
<div data-test-namespace-link="{{option.label}}">
```

---

# Content and Terminology

## Title and Heading Case Rules
- **USE SENTENCE CASE**: All HTML headings (`<h1>`, `<h2>`, `<h3>`, etc.) should use sentence case (only first letter capitalized)
- **NO TITLE CASE**: Avoid title case where every major word is capitalized
- **Component arguments**: Use sentence case for `@title`, `@label`, and similar text properties
- End descriptive text with proper punctuation
- Follow proper grammar rules including ending sentences with periods
- Use consistent terminology for product-specific features and components

Examples:
```handlebars
{{!-- CORRECT: Sentence case in HTML headings --}}
<h2 class="title is-4">Quick actions</h2>
<h3 class="title is-marginless is-6">Configuration settings</h3>
<h1>Authentication methods</h1>

{{!-- INCORRECT: Title case in HTML headings --}}
<h2 class="title is-4">Quick Actions</h2>
<h3 class="title is-marginless is-6">Configuration Settings</h3>
<h1>Authentication Methods</h1>
```

```javascript
// CORRECT: Sentence case in component arguments
@title="Upload user's profile"
@label="Configuration path"
@placeholder="Enter mount path"

// INCORRECT: Title case or inconsistent casing
@title="Upload User's Profile"
@label="Configuration Path"
@placeholder="Enter Mount Path"
```

```handlebars
{{!-- CORRECT: Sentence case in data-test attributes when they contain readable text --}}
data-test-card-subtitle="configuration-settings"

{{!-- Component usage with proper casing --}}
<Hds::Button @text="Create configuration" />
<Hds::Alert @message="Operation completed successfully" />
```
