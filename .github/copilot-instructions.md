# GitHub Copilot Instructions for Vault Enterprise

This repository contains coding guidelines and best practices for the Vault Enterprise project.

## Instruction Files

The `.github/instructions/` directory contains domain-specific coding guidelines:

### Go Development

- **[golang.instructions.md](instructions/generic/golang.instructions.md)**: General Go programming guidelines covering code style, package organization, error handling, and idiomatic patterns
- **[golang_tests.instructions.md](instructions/generic/golang_tests.instructions.md)**: Best practices for writing tests in Go, including table-driven tests, test structure, and integration testing

### Ember.js UI Development

- **[ember_general.instructions.md](instructions/generic/ember_general.instructions.md)**: Project context, repository structure, and cross-cutting concerns for Ember applications
- **[ember_js.instructions.md](instructions/generic/ember_js.instructions.md)**: JavaScript/TypeScript guidelines for Ember components and services
- **[ember_hbs.instructions.md](instructions/generic/ember_hbs.instructions.md)**: Handlebars template guidelines and best practices
- **[ember_styles.instructions.md](instructions/generic/ember_styles.instructions.md)**: SCSS styling guidelines and HashiCorp Design System usage
- **[ember_tests.instructions.md](instructions/generic/ember_tests.instructions.md)**: Testing guidelines for Ember applications

### General Development

- **[code_comments.instructions.md](instructions/generic/code_comments.instructions.md)**: Guidelines for writing effective code comments and self-documenting code
- **[testing.instructions.md](instructions/generic/testing.instructions.md)**: General testing best practices and guidelines

## How These Instructions Work

GitHub Copilot automatically applies these instructions based on the `applyTo` patterns defined in each file's frontmatter. For example:

- Go-specific rules apply to `**/*.go` files
- Go test rules apply to `**/*_test.go` files
- Ember UI rules apply to `vault/ui/**/*` files

When writing code, Copilot will reference the appropriate instruction files to provide context-aware suggestions that align with project standards.
