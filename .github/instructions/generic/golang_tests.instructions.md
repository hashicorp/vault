---
applyTo: "**/*_test.go"
description: "Practical guidance for writing tests in Go"
---

# Golang Testing Guidelines
*Practical guidance for writing tests in Go*

## General Guidelines
- Write unit tests for all public functions
- Use table-driven tests when appropriate
- Follow naming convention: `Test_FunctionName`
- Use `t.Parallel()` for faster test execution and to encourage writing tests that don't share state

## Tools
- Use testify/assert and testify/require for consistent assertions

## Test Structure
- Use t.Helper() to mark helper functions
- Use subtests for organizing related test cases with t.Run()

## Integration Tests 
- Use testcontainers to spin up dependencies like postgres or redis

## Black box testing 
- The test should be in feature_test package to be the first client of the feature package

## End-to-End Testing (Future Enhancement)
<!-- TODO: Add guidance for end-to-end testing with enos and cloud-first testing in Vault Cloud -->
