---
applyTo: "vault/ui/**/*"
description: "HashiCorp Ember.js UI general guidelines - project context, build structure, and cross-cutting concerns"
---

# HashiCorp Ember.js UI General Guidelines

This document provides general coding standards and project context for HashiCorp Ember.js UI applications. This serves as the central reference for project structure, framework information, and cross-cutting concerns.

> **Related Documents**: See domain-specific guidelines in `ember_javascript.instructions.md`, `ember_hbs.instructions.md`, `ember_styles.instructions.md`, and `ember_tests.instructions.md`.

## Project Context
HashiCorp Ember.js UI applications provide web-based interfaces for HashiCorp's infrastructure tools and cloud platforms. These applications serve enterprise-grade interfaces used by DevOps teams, platform engineers, security professionals, and system administrators to manage infrastructure, security policies, and operational workflows.

Applications use modern Ember Octane patterns, Handlebars templates, and SCSS for styling with the HashiCorp Design System components.

## Repository Structure
- `ui/app/components/` - Reusable UI components and their templates
- `ui/app/models/` - Ember Data models representing API entities
- `ui/app/routes/` - Route handlers for URL endpoints and data loading logic
- `ui/app/templates/` - Page-level Handlebars templates 
- `ui/app/services/` - Ember services for shared functionality and state management
- `ui/app/helpers/` - Template helper functions for data formatting and logic
- `ui/tests/` - Integration, unit, and acceptance tests
- `ui/app/styles/` - SCSS stylesheets and component-specific styles
- `ui/config/` - Ember CLI configuration and environment settings
- `ui/mirage/` - Mock server configuration for development and testing

## Framework and Tools
- **Frontend Framework**: Ember.js 4.x with Ember Octane patterns and decorators
- **Data Layer**: Ember Data for API communication with backend services
- **Templating**: Handlebars templates with Ember's component architecture
- **Styling**: SCSS with HashiCorp Design System (HDS) components and Bulma CSS framework
- **Build System**: Ember CLI with Webpack for bundling and asset management
- **Development**: ESLint for code linting, Prettier for code formatting

## File Naming Conventions
- **Component Structure**: Components should have matching JavaScript/TypeScript and `.hbs` files in the same directory
- **Directory Organization**: Organize new components in logical subdirectories by feature or domain

## Build and Deployment Structure
- **Development Build**: `ember serve` creates development server with live reload
- **Production Build**: `ember build --environment=production` generates optimized assets in `dist/`
- **Asset Output**: Built files are typically served by the application's backend server or embedded in the main binary
- **Integration**: UI is often compiled and integrated into the main application during release builds
- **Testing**: `ember test` runs the full test suite with QUnit in headless Chrome

---

# Changelog Guidelines

For files in the `changelog/` directory:
- **Enterprise features**: Use `ui (enterprise): descriptive text`  
- **Community features**: Use `ui: descriptive text`
- Always indicate enterprise-only features in the description

```javascript
// Changelog entries
"ui (enterprise): Add advanced policy filtering"  // enterprise features
"ui: Fix configuration list pagination"           // community features
```

---

# Dependency Management

## package.json Guidelines
- Pin exact versions for critical dependencies or use tilde (`~`) for patch updates only
- **WARNING**: Avoid caret (`^`) operator which allows minor version updates that may introduce breaking changes
- Use tilde (`~`) for regular dependencies, exact versions for security-critical packages
- Dependencies in `resolutions` block MUST be pinned (no `~` or `^`)
- Keep `package.json` changes minimal and focused on the specific feature or fix
- Always commit lock files (`yarn.lock`, `package-lock.json`) with dependency changes
- Ensure package.json changes are independent of other code changes (except lock files)
- **Reminder**: Consider coordinating dependency changes with backend teams when the UI is embedded in application binaries

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
