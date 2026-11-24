/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import NamespaceForm from 'vault/forms/namespace';

module('Unit | Form | namespace', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.validate = (path) => new NamespaceForm({ path }, { isNew: true }).toJSON();
  });

  test('it rejects a blank path', function (assert) {
    const { isValid, state, invalidFormMessage } = this.validate('');

    assert.false(isValid, 'a blank path is invalid');
    assert.deepEqual(state.path.errors, [`Path can't be blank.`], 'presence message is returned');
    assert.strictEqual(invalidFormMessage, 'There is an error with this form.', 'error count message');
  });

  test('it rejects a path ending in a forward slash', function (assert) {
    const { isValid, state } = this.validate('invalid-namespace/');

    assert.false(isValid, 'a trailing slash is invalid');
    assert.deepEqual(
      state.path.errors,
      [`Path can't end in forward slash '/'.`],
      'endsInSlash message is returned'
    );
  });

  test('it rejects a path containing whitespace', function (assert) {
    const { isValid, state } = this.validate('invalid namespace');

    assert.false(isValid, 'whitespace is invalid');
    assert.deepEqual(
      state.path.errors,
      [`Path can't contain whitespace.`],
      'containsWhiteSpace message is returned'
    );
  });

  test('it accepts a valid path', function (assert) {
    const { isValid, state, invalidFormMessage } = this.validate('valid-namespace');

    assert.true(isValid, 'a simple path is valid');
    assert.deepEqual(state.path.errors, [], 'no errors are returned');
    assert.strictEqual(invalidFormMessage, '', 'no error count message');
  });

  test('it reports every failing rule for a single path', function (assert) {
    const { isValid, state, invalidFormMessage } = this.validate('invalid namespace/');

    assert.false(isValid, 'a path breaking two rules is invalid');
    assert.deepEqual(
      state.path.errors,
      [`Path can't end in forward slash '/'.`, `Path can't contain whitespace.`],
      'both messages are returned in validation order'
    );
    assert.strictEqual(invalidFormMessage, 'There are 2 errors with this form.', 'error count message');
  });

  test('it exposes the path field for the create form', function (assert) {
    const form = new NamespaceForm({}, { isNew: true });

    assert.deepEqual(
      form.formFields.map((field) => field.name),
      ['path'],
      'the create form renders only the path field'
    );
  });

  test('it returns the entered path for the create request', function (assert) {
    const { data } = this.validate('valid-namespace');

    assert.deepEqual(data, { path: 'valid-namespace' }, 'the path is carried through to the API call');
  });
});
