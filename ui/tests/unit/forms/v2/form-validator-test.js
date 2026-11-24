/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { validateAllFields, validateField } from 'vault/forms/v2/form-validator';

module('Unit | forms/v2 | form-validator', function () {
  test('validateField returns empty array for field without validations', function (assert) {
    const field = {
      name: 'username',
      type: 'TextInput',
      label: 'Username',
    };

    const errors = validateField(field, 'test-value', {});

    assert.strictEqual(errors.length, 0, 'no errors for field without validations');
  });

  test('validateField returns errors for invalid required field', function (assert) {
    const field = {
      name: 'username',
      type: 'TextInput',
      label: 'Username',
      validations: [{ type: 'required', message: 'Username is required' }],
    };

    const errors = validateField(field, '', {});

    assert.strictEqual(errors.length, 1, 'has one error');
    assert.strictEqual(errors[0], 'Username is required', 'error message matches');
  });

  test('validateField returns empty array for valid required field', function (assert) {
    const field = {
      name: 'username',
      type: 'TextInput',
      label: 'Username',
      validations: [{ type: 'required', message: 'Username is required' }],
    };

    const errors = validateField(field, 'john_doe', {});

    assert.strictEqual(errors.length, 0, 'no errors for valid field');
  });

  test('validateField validates multiple rules', function (assert) {
    const field = {
      name: 'username',
      type: 'TextInput',
      label: 'Username',
      validations: [
        { type: 'required', message: 'Username is required' },
        { type: 'minLength', options: { minLength: 3 }, message: 'Must be at least 3 characters' },
        { type: 'maxLength', options: { maxLength: 20 }, message: 'Must be at most 20 characters' },
      ],
    };

    // Test too short
    let errors = validateField(field, 'ab', {});
    assert.strictEqual(errors.length, 1, 'has one error for too short');
    assert.ok(errors[0].includes('3 characters'), 'error mentions minimum length');

    // Test too long
    errors = validateField(field, 'a'.repeat(25), {});
    assert.strictEqual(errors.length, 1, 'has one error for too long');
    assert.ok(errors[0].includes('20 characters'), 'error mentions maximum length');

    // Test valid
    errors = validateField(field, 'john_doe', {});
    assert.strictEqual(errors.length, 0, 'no errors for valid value');
  });

  test('validateField supports custom validator functions', function (assert) {
    const field = {
      name: 'password',
      type: 'TextInput',
      label: 'Password',
      validations: [
        {
          validator: (formData) => {
            return formData.password === formData.confirmPassword;
          },
          message: 'Passwords must match',
        },
      ],
    };

    const payload = {
      password: 'secret123',
      confirmPassword: 'different',
    };

    const errors = validateField(field, 'secret123', payload);

    assert.strictEqual(errors.length, 1, 'has one error');
    assert.strictEqual(errors[0], 'Passwords must match', 'custom validation error');
  });

  test('validateField supports dynamic error messages', function (assert) {
    const field = {
      name: 'username',
      type: 'TextInput',
      label: 'Username',
      validations: [
        {
          type: 'required',
          message: (formData) => `${formData.fieldLabel || 'This field'} is required`,
        },
      ],
    };

    const payload = { fieldLabel: 'User Name' };
    const errors = validateField(field, '', payload);

    assert.strictEqual(errors.length, 1, 'has one error');
    assert.strictEqual(errors[0], 'User Name is required', 'dynamic message generated');
  });

  test('validateField uses default message when none provided', function (assert) {
    const field = {
      name: 'email',
      type: 'TextInput',
      label: 'Email',
      validations: [
        { type: 'email' }, // No message provided
      ],
    };

    const errors = validateField(field, 'invalid-email', {});

    assert.strictEqual(errors.length, 1, 'has one error');
    assert.ok(errors[0].includes('email'), 'uses default email error message');
  });

  test('validateField returns all failing validations', function (assert) {
    const field = {
      name: 'username',
      type: 'TextInput',
      label: 'Username',
      validations: [
        { type: 'required', message: 'Required' },
        { type: 'minLength', options: { minLength: 3 }, message: 'Too short' },
        { type: 'maxLength', options: { maxLength: 20 }, message: 'Too long' },
      ],
    };

    // Empty value should fail required and minLength
    const errors = validateField(field, '', {});

    assert.deepEqual(errors, ['Required', 'Too short'], 'returns messages for all failing rules');
  });

  test('validateAllFields validates all fields in form', function (assert) {
    const fields = [
      {
        name: 'username',
        type: 'TextInput',
        label: 'Username',
        validations: [{ type: 'required', message: 'Username is required' }],
      },
      {
        name: 'email',
        type: 'TextInput',
        label: 'Email',
        validations: [
          { type: 'required', message: 'Email is required' },
          { type: 'email', message: 'Must be valid email' },
        ],
      },
      {
        name: 'age',
        type: 'TextInput',
        label: 'Age',
        validations: [{ type: 'min', options: { min: 18 }, message: 'Must be 18 or older' }],
      },
    ];

    const payload = {
      username: '',
      email: 'invalid',
      age: 15,
    };

    const errorMap = validateAllFields(fields, payload);

    assert.strictEqual(errorMap.size, 3, 'has errors for 3 fields');
    assert.ok(errorMap.has('username'), 'has username error');
    assert.ok(errorMap.has('email'), 'has email error');
    assert.ok(errorMap.has('age'), 'has age error');
  });

  test('validateAllFields returns empty map for valid form', function (assert) {
    const fields = [
      {
        name: 'username',
        type: 'TextInput',
        label: 'Username',
        validations: [{ type: 'required', message: 'Username is required' }],
      },
      {
        name: 'email',
        type: 'TextInput',
        label: 'Email',
        validations: [
          { type: 'required', message: 'Email is required' },
          { type: 'email', message: 'Must be valid email' },
        ],
      },
    ];

    const payload = {
      username: 'john_doe',
      email: 'john@example.com',
    };

    const errorMap = validateAllFields(fields, payload);

    assert.strictEqual(errorMap.size, 0, 'no errors for valid form');
  });

  test('validateAllFields handles nested field paths', function (assert) {
    const fields = [
      {
        name: 'user.profile.name',
        type: 'TextInput',
        label: 'Name',
        validations: [{ type: 'required', message: 'Name is required' }],
      },
    ];

    const payload = {
      user: {
        profile: {
          name: '',
        },
      },
    };

    const errorMap = validateAllFields(fields, payload);

    assert.strictEqual(errorMap.size, 1, 'has one error');
    assert.ok(errorMap.has('user.profile.name'), 'error keyed by full path');
  });

  test('validateAllFields only includes fields with errors', function (assert) {
    const fields = [
      {
        name: 'username',
        type: 'TextInput',
        label: 'Username',
        validations: [{ type: 'required', message: 'Username is required' }],
      },
      {
        name: 'email',
        type: 'TextInput',
        label: 'Email',
        validations: [{ type: 'required', message: 'Email is required' }],
      },
    ];

    const payload = {
      username: 'john_doe', // Valid
      email: '', // Invalid
    };

    const errorMap = validateAllFields(fields, payload);

    assert.strictEqual(errorMap.size, 1, 'only one field has errors');
    assert.notOk(errorMap.has('username'), 'valid field not in error map');
    assert.ok(errorMap.has('email'), 'invalid field in error map');
  });

  test('validateField handles validator options correctly', function (assert) {
    const field = {
      name: 'code',
      type: 'TextInput',
      label: 'Code',
      validations: [
        {
          type: 'pattern',
          options: { pattern: '^[A-Z]{3}$', flags: 'i' },
          message: 'Must be 3 letters',
        },
      ],
    };

    // Test valid (case insensitive due to 'i' flag)
    let errors = validateField(field, 'abc', {});
    assert.strictEqual(errors.length, 0, 'no errors for valid pattern');

    // Test invalid
    errors = validateField(field, '123', {});
    assert.strictEqual(errors.length, 1, 'has error for invalid pattern');
  });
});
