/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import LocalStorage from 'vault/lib/local-storage';

module('Unit | lib | local-storage', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    window.localStorage.clear();
  });

  test('it does not error if nothing is in local storage', async function (assert) {
    assert.expect(1);
    assert.strictEqual(
      LocalStorage.cleanupStorage('something', 'something-key'),
      undefined,
      'returns undefined and does not throw an error when method is called and nothing exist in localStorage.'
    );
  });

  test('it does not remove anything in localStorage that does not start with the string or we have specified to keep.', async function (assert) {
    assert.expect(3);
    LocalStorage.setItem('string-key-remove', 'string-key-remove-value');
    LocalStorage.setItem('beep-boop-bop-key', 'beep-boop-bop-value');
    LocalStorage.setItem('string-key', 'string-key-value');
    const storageLengthBefore = window.localStorage.length;
    LocalStorage.cleanupStorage('string', 'string-key');
    const storageLengthAfter = window.localStorage.length;
    assert.strictEqual(
      storageLengthBefore - storageLengthAfter,
      1,
      'the method should only remove one key from localStorage.'
    );
    assert.strictEqual(
      LocalStorage.getItem('string-key'),
      'string-key-value',
      'the key we asked to keep still exists in localStorage.'
    );
    assert.strictEqual(
      LocalStorage.getItem('string-key-remove'),
      null,
      'the key we did not specify to keep was removed from localStorage.'
    );
    // clear storage
    window.localStorage.clear();
  });
});
