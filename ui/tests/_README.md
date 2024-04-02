# Test Helpers Organization

Our test are constantly evolving, but here's a general overview of how the tests and selectors organization

## Folder organization

### acceptance

Acceptance tests should test the overall workflows and navigation within Vault. When possible, they should use the real API instead of mocked so that breaking changes from the backend can be caught. Reasons you may opt to use a mocked backend instead of the real one:

- Using the real backend would cause instability in concurrently-running tests (eg. seal/unseal flow)
- There isn't a way to set up a 3rd party dependency that the backend needs to run correctly (Database Secrets Engine, Sync Secrets)

### helpers

### integration

### pages

[DEPRECATED] This file should be removed in favor of selectors within the "helpers" folder

### unit

## Process

- Rename export from `general-selectors` to `GENERAL`
-
