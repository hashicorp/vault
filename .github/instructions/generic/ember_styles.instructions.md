---
applyTo: "vault/ui/**/*.scss"
description: "HashiCorp Ember.js UI SCSS styling coding standards"
---

# HashiCorp Ember.js SCSS Styling Guidelines

This document provides SCSS styling coding standards for HashiCorp Ember.js UI applications.

> **Note**: For general project context, framework information, and repository structure, see `ember_general.instructions.md`.

## CSS Best Practices
- Avoid `z-index` - manage stacking order through DOM structure instead
- Avoid `!important` - use more specific selectors for overrides

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
