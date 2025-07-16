/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { validate } from 'vault/utils/forms/validate';
import { module, test } from 'qunit';
import sinon from 'sinon';
import validators from 'vault/utils/forms/validators';

module('Unit | Utility | forms | validate', function (hooks) {
  hooks.beforeEach(function () {
    this.data = {
      foo: null,
      integer: null,
    };
    this.consoleStub = sinon.stub(console, 'error');
  });

  hooks.afterEach(function () {
    this.consoleStub.restore();
  });

  test('it should log error to console when validations are not passed as array', function (assert) {
    const validations = {
      foo: { type: 'presence', message: 'Foo is required' },
    };
    validate(this.data, validations);
    const message = 'Must provide validations as an array for property "foo".';
    assert.true(this.consoleStub.calledWith(message));
  });

  test('it should return valid when validations are not provided', async function (assert) {
    const expected = { isValid: true, state: {}, invalidFormMessage: '' };
    const validation = validate(this.data);
    assert.deepEqual(validation, expected, 'Data is considered valid when validations are not provided');
  });

  test('it should log error for incorrect validator type', function (assert) {
    const validations = {
      foo: [{ type: 'bar', message: 'Foo is bar' }],
    };
    validate(this.data, validations);
    const types = Object.keys(validators).join(', ');
    const message = `Validator type: "bar" not found. Available validators: ${types}`;
    assert.ok(this.consoleStub.calledWith(message));
  });

  test('it should validate', function (assert) {
    const message = 'This field is required';
    const validations = {
      foo: [{ type: 'presence', message }],
    };
    const v1 = validate(this.data, validations);
    assert.false(v1.isValid, 'isValid state is correct when errors exist');
    assert.deepEqual(
      v1.state,
      { foo: { isValid: false, errors: [message], warnings: [] } },
      'Correct state returned when property is invalid'
    );

    this.data.foo = true;
    const v2 = validate(this.data, validations);
    assert.true(v2.isValid, 'isValid state is correct when no errors exist');
    assert.deepEqual(
      v2.state,
      { foo: { isValid: true, errors: [], warnings: [] } },
      'Correct state returned when property is valid'
    );
  });

  test('invalid form message has correct error count', function (assert) {
    const message = 'This field is required';
    const messageII = 'This field must be a number';
    const validations = {
      foo: [{ type: 'presence', message }],
      integer: [{ type: 'number', messageII }],
    };
    const v1 = validate(this.data, validations);
    assert.strictEqual(
      v1.invalidFormMessage,
      'There are 2 errors with this form.',
      'error message says form as 2 errors'
    );

    this.data.integer = 9;
    const v2 = validate(this.data, validations);
    assert.strictEqual(
      v2.invalidFormMessage,
      'There is an error with this form.',
      'error message says form has an error'
    );

    this.data.foo = true;
    const v3 = validate(this.data, validations);
    assert.strictEqual(v3.invalidFormMessage, '', 'invalidFormMessage is empty when form is valid');
  });

  test('it should validate warnings', function (assert) {
    const message = 'Value contains whitespace.';
    const validations = {
      foo: [
        {
          type: 'containsWhiteSpace',
          message,
          level: 'warn',
        },
      ],
    };
    this.data.foo = 'foo bar';
    const { state, isValid } = validate(this.data, validations);
    assert.true(isValid, 'Data is considered valid when there are only warnings');
    assert.strictEqual(state.foo.warnings.join(' '), message, 'Warnings are returned');
  });

  test('it should accept a key to map validations to in state object', async function (assert) {
    this.data.foo = undefined;
    const validations = {
      foo: [{ type: 'presence', message: 'Foo is required' }],
    };
    const validation = validate(this.data, validations, 'data');
    const expected = {
      isValid: false,
      state: { 'data.foo': { isValid: false, errors: ['Foo is required'], warnings: [] } },
      invalidFormMessage: 'There is an error with this form.',
    };
    assert.deepEqual(validation, expected, 'Validation state is mapped to the correct key');
  });
});
