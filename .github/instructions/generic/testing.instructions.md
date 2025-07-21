---
applyTo: "**/*_test.go,**/*_spec.rb"
description: "Practical guidance for writing tests across languages"
---

# Philosophy & Principles

## Testing: 
- Test coverage is not a metric to optimize for
- Test behavior and not implementation detail
- Write tests that are robust to refactoring the code under test

###  Value of a test
- Test is code that needs to be maintained
- If a test is not valuable, delete it

###  Tests are documentation
- Use the given-when-then pattern to communicate each test case
- Test names should express intent 

## Test Design & Quality 
*How to design reliable, maintainable tests*

### Test Isolation 
- A unit test should only fail for one reason
- Do not share state between tests
- Cleanup resources after the test
- If files are needed, use a temp directory

### Test Stability
- Do not use sleeping in a test

### Test Structure
- Consider parametrized tests where it makes sense

### Mocks
- Preferably do not use mocks
- If you use mocks, assert the mock expectations

## Supporting Code Quality
*Making the code under test more testable*

### Code under test
- If the code relies on magic numbers or constants: Extract them
- Consider introducing interfaces
