/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import SecretForm from 'vault/forms/secret';

module('Unit | Form | secret', function () {
  test('create requires a path and serializes only secret data', function (assert) {
    const form = new SecretForm(
      { path: '', secretData: { username: 'alice', optional: '', enabled: false } },
      { isNew: true }
    );
    const { isValid, state, data } = form.toJSON();

    assert.false(isValid, 'a create form without a path is invalid');
    assert.deepEqual(state['path']?.errors, ["Path can't be blank."], 'the path validation is returned');
    assert.deepEqual(
      data,
      { username: 'alice', optional: '', enabled: false },
      'the write payload preserves values and omits the path'
    );
  });

  test('edit does not require a path and serializes nested secret data', function (assert) {
    const form = new SecretForm(
      { id: 'existing', secretData: { metadata: { enabled: true } } },
      { isNew: false }
    );
    const { isValid, data } = form.toJSON();

    assert.true(isValid, 'an edit form is not invalidated by create-only path validation');
    assert.deepEqual(
      data,
      { metadata: { enabled: true } },
      'nested secret values remain in the write payload'
    );
  });
});
