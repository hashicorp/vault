/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { withModelValidations } from 'vault/decorators/model-validations';
import validators from 'vault/utils/model-helpers/validators';
import sinon from 'sinon';
import Model from '@ember-data/model';

// create class using decorator
const createClass = (validations) => {
  @withModelValidations(validations)
  class Foo extends Model {}
  const foo = Foo.extend({
    modelName: 'bar',
    foo: null,
    integer: null,
  });
  return new foo();
};

module('Unit | Decorators | ModelValidations', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.spy = sinon.spy(console, 'error');
  });
  hooks.afterEach(function () {
    this.spy.restore();
  });

  test('it should throw error when validations object is not provided', function (assert) {
    assert.expect(1);

    try {
      createClass();
    } catch (e) {
      assert.strictEqual(e.message, 'Validations object must be provided to constructor for setup');
    }
  });

  test('it should log error to console when validations are not passed as array', function (assert) {
    const validations = {
      foo: { type: 'presence', message: 'Foo is required' },
    };
    const fooClass = createClass(validations);
    fooClass.validate();
    const message = 'Must provide validations as an array for property "foo" on bar model';
    assert.ok(this.spy.calledWith(message));
  });

  test('it should log error for incorrect validator type', function (assert) {
    const validations = {
      foo: [{ type: 'bar', message: 'Foo is bar' }],
    };
    const fooClass = createClass(validations);
    fooClass.validate();
    const message = `Validator type: "bar" not found. Available validators: ${Object.keys(validators).join(
      ', '
    )}`;
    assert.ok(this.spy.calledWith(message));
  });

  test('it should validate', function (assert) {
    const message = 'This field is required';
    const validations = {
      foo: [{ type: 'presence', message }],
    };
    const fooClass = createClass(validations);
    const v1 = fooClass.validate();
    assert.false(v1.isValid, 'isValid state is correct when errors exist');
    assert.deepEqual(
      v1.state,
      { foo: { isValid: false, errors: [message], warnings: [] } },
      'Correct state returned when property is invalid'
    );

    fooClass.foo = true;
    const v2 = fooClass.validate();
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
    const fooClass = createClass(validations);
    const v1 = fooClass.validate();
    assert.strictEqual(
      v1.invalidFormMessage,
      'There are 2 errors with this form.',
      'error message says form as 2 errors'
    );

    fooClass.integer = 9;
    const v2 = fooClass.validate();
    assert.strictEqual(
      v2.invalidFormMessage,
      'There is an error with this form.',
      'error message says form has an error'
    );

    fooClass.foo = true;
    const v3 = fooClass.validate();
    assert.strictEqual(v3.invalidFormMessage, null, 'invalidFormMessage is null when form is valid');
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
    const fooClass = createClass(validations);
    fooClass.foo = 'foo bar';
    const { state, isValid } = fooClass.validate();
    assert.true(isValid, 'Model is considered valid when there are only warnings');
    assert.strictEqual(state.foo.warnings.join(' '), message, 'Warnings are returned');
  });
});
