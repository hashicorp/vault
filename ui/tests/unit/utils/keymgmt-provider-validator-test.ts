/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isValidProvider } from 'vault/utils/keymgmt-provider-validator';
import { module, test } from 'qunit';

module('Unit | Util | keymgmt-provider-validator', function () {
  test('it returns true for valid provider strings', function (assert) {
    assert.true(isValidProvider('azure-provider'), 'Valid provider name');
    assert.true(isValidProvider('test-provider'), 'Another valid provider');
    assert.true(isValidProvider('a'), 'Single character string');
    assert.true(isValidProvider('  valid-provider  '), 'Provider with leading/trailing spaces');
  });

  test('it returns false for objects', function (assert) {
    assert.false(isValidProvider({ permissionsError: true }), 'Object with permissionsError');
    assert.false(isValidProvider({}), 'Empty object');
    assert.false(isValidProvider({ name: 'provider' }), 'Object with properties');
  });

  test('it returns false for non-string primitives', function (assert) {
    assert.false(isValidProvider(123), 'Number returns false');
    assert.false(isValidProvider(true), 'Boolean true returns false');
    assert.false(isValidProvider(false), 'Boolean false returns false');
  });

  test('it returns false for arrays', function (assert) {
    assert.false(isValidProvider([]), 'Empty array');
    assert.false(isValidProvider(['provider']), 'Array with string');
  });
});
