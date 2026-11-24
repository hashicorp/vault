/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { getKeymgmtProviderIcon, isValidProvider } from 'vault/utils/keymgmt-provider-utils';
import { module, test } from 'qunit';

module('Unit | Util | keymgmt-provider-utils', function () {
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

  test('it returns provider-specific icons for known provider types', function (assert) {
    assert.strictEqual(getKeymgmtProviderIcon('azurekeyvault'), 'azure-color', 'Azure icon is returned');
    assert.strictEqual(getKeymgmtProviderIcon('awskms'), 'aws-color', 'AWS icon is returned');
    assert.strictEqual(getKeymgmtProviderIcon('gcpckms'), 'gcp-color', 'GCP icon is returned');
  });

  test('it returns default icon for unknown or empty provider types', function (assert) {
    assert.strictEqual(
      getKeymgmtProviderIcon('unknown-provider'),
      'key',
      'Unknown provider uses default icon'
    );
    assert.strictEqual(getKeymgmtProviderIcon(''), 'key', 'Empty string uses default icon');
    assert.strictEqual(getKeymgmtProviderIcon(), 'key', 'Undefined provider uses default icon');
  });
});
