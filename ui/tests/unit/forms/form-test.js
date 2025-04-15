/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import { module, test } from 'qunit';

module('Unit | forms | form', function (hooks) {
  hooks.beforeEach(function () {
    this.data = {
      foo: 'bar',
    };
    this.validations = {
      foo: [{ type: 'presence', message: 'Foo is required' }],
    };
  });

  test('it should set data and validations on class', async function (assert) {
    const form = new Form(this.data, this.validations);
    assert.deepEqual(form.data, this.data, 'Data is set correctly');
    assert.deepEqual(form.validations, this.validations, 'Validations are set correctly');
  });

  test('it should shim Ember object set method', async function (assert) {
    const form = new Form(this.data);
    form.set('data.foo', 'baz');
    assert.strictEqual(form.data.foo, 'baz', 'Property is updated using set method');
  });

  test('it should return data as valid without validations from toJSON method', async function (assert) {
    const expected = { isValid: true, state: {}, invalidFormMessage: '', data: this.data };
    const form = new Form(this.data);
    const json = form.toJSON();
    assert.deepEqual(json, expected, 'toJSON returns correct data and state');
  });

  test('it should validate and return data from toJSON method', async function (assert) {
    const state = { 'data.foo': { isValid: true, errors: [], warnings: [] } };
    const expected = { isValid: true, state, invalidFormMessage: '', data: this.data };

    const form = new Form(this.data, this.validations);
    const json = form.toJSON();

    assert.deepEqual(json, expected, 'toJSON returns correct data and validation state');
  });

  test('it should allow for custom data to be passed in toJSON method', async function (assert) {
    this.data.foo = 'string with whitespace';
    this.validations.foo = [{ type: 'containsWhiteSpace', message: 'Whitespace is not allowed' }];

    const serializedData = { foo: this.data.foo.replace(/\s/g, '_') };
    const state = { 'data.foo': { isValid: true, errors: [], warnings: [] } };
    const expected = { isValid: true, state, invalidFormMessage: '', data: serializedData };

    const form = new Form(this.data, this.validations);
    const json = form.toJSON(serializedData);

    assert.deepEqual(json, expected, 'toJSON allows data to be overridden');
  });
});
