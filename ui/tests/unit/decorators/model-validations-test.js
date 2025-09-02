/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { withModelValidations } from 'vault/decorators/model-validations';
import sinon from 'sinon';
import Model from '@ember-data/model';
import * as validateUtil from 'vault/utils/forms/validate';

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
    this.validateSpy = sinon.spy(validateUtil, 'validate');
  });
  hooks.afterEach(function () {
    sinon.restore();
  });

  test('it should throw error when validations object is not provided', function (assert) {
    assert.expect(1);

    try {
      createClass();
    } catch (e) {
      assert.strictEqual(e.message, 'Validations object must be provided to constructor for setup');
    }
  });

  test('it should validate', function (assert) {
    const message = 'This field is required';
    const validations = {
      foo: [{ type: 'presence', message }],
    };
    const fooClass = createClass(validations);
    fooClass.validate();
    assert.true(
      this.validateSpy.calledWith(fooClass, validations),
      'validate util called with correct arguments'
    );
  });
});
