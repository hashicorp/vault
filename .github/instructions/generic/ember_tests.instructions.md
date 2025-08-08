---
applyTo: "vault/ui/tests/**/*.{js,ts}"
description: "HashiCorp Ember.js testing standards and best practices"
---

# HashiCorp Ember.js Testing Standards

This document provides testing standards and best practices for HashiCorp Ember.js UI applications.

> **Note**: For general project context, framework information, and repository structure, see `ember_general.instructions.md`.

## Testing Framework and Tools
- **Testing Framework**: QUnit for unit, integration, and acceptance tests
- **Mock Server**: Ember CLI Mirage for API mocking and test data
- **Test Runners**: Ember CLI test runner with headless Chrome

## Test Directory Structure
- `ui/tests/integration/` - Component integration tests
- `ui/tests/unit/` - Service, helper, and model unit tests  
- `ui/tests/acceptance/` - End-to-end user workflow tests
- `ui/mirage/` - Mock server configuration for development and testing
- `ui/tests/helpers/` - Custom test helper functions and utilities

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

## Asynchronous Testing
- Use Ember's `run.later` instead of `setTimeout` in tests for better runloop control
- Handle async operations with proper waiting patterns in tests
- Ensure test isolation by resetting state between tests

## Mirage Server Setup
- Place `this.server` setup at the top of test blocks when using ember-cli-mirage
- Configure mock data that reflects realistic API responses
- Use mirage factories for generating test data consistently
- Reset server state between tests to ensure test isolation

Example mirage setup:
```javascript
module('Integration | Component | my-component', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.server = new Server();
    this.server.create('model', { id: 1, name: 'Test Item' });
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });
});
```

## DOM Testing Best Practices
- Use `data-test-*` selectors for DOM interactions to decouple from styling
- Test user interactions through realistic user flows
- Verify state changes after user actions
- Test error states and edge cases
- Ensure accessibility features work correctly

Example DOM testing:
```javascript
// Good: using data-test selectors
await click('[data-test-submit-button]');
assert.true(find('[data-test-success-message]'), 'Success message should appear after submission');

// Bad: using CSS classes for testing
await click('.btn-submit');
assert.ok(find('.alert-success'));
```

## Test Organization
- Group related tests in logical modules
- Use descriptive test names that explain the scenario being tested
- Include both positive and negative test cases
- Test error handling and edge cases
- Keep tests focused on a single behavior or outcome

Example test structure:
```javascript
module('Integration | Component | secret-form', function (hooks) {
  setupRenderingTest(hooks);

  test('it displays validation error when secret name is empty', async function (assert) {
    await render(hbs`<SecretForm @onSubmit={{this.handleSubmit}} />`);
    
    await click('[data-test-submit-button]');
    
    assert.true(
      find('[data-test-name-error]').textContent.includes('Secret name is required'),
      'Should display validation error for empty secret name'
    );
  });

  test('it calls onSubmit with form data when valid', async function (assert) {
    let submittedData;
    this.handleSubmit = (data) => { submittedData = data; };

    await render(hbs`<SecretForm @onSubmit={{this.handleSubmit}} />`);
    
    await fillIn('[data-test-secret-name]', 'my-secret');
    await fillIn('[data-test-secret-value]', 'secret-value');
    await click('[data-test-submit-button]');
    
    assert.deepEqual(submittedData, {
      name: 'my-secret',
      value: 'secret-value'
    }, 'Should submit form data with correct values');
  });
});
```

---

# Debugging Tests

## Debugging Tips
When debugging Ember templates in tests:
- Use `{{debugger}}` inside templates to inspect values in the browser console
- In the console during debugging:
  - Use `get('property.name')` to inspect nested properties
  - Use `context` to explore the current template context
- Add console.log statements in test code to track execution flow
- Use browser developer tools to inspect DOM state during test execution

Example debugging usage:
```handlebars
{{#each this.items as |item|}}
  {{debugger}}
{{/each}}
```

```javascript
test('debugging example', async function (assert) {
  await render(hbs`<MyComponent @data={{this.testData}} />`);
  
  // Pause execution to inspect DOM
  debugger;
  
  console.log('Component state:', this.element.querySelector('[data-test-component]'));
  
  await click('[data-test-button]');
  
  // Inspect state after interaction
  debugger;
});
```

## Test Isolation
- Reset component state between tests
- Clear any global state or services
- Ensure tests don't depend on execution order
- Use hooks for setup and teardown consistently
