/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import TransitKeyResource from 'vault/resources/transit/key';

module('Unit | Resource | transit key', function () {
  test('actions and export types follow API capabilities in display order', function (assert) {
    const key = new TransitKeyResource({
      backend: 'transit',
      id: 'key',
      supports_encryption: true,
      supports_decryption: true,
    });
    assert.deepEqual(
      key.supportedActions.map(({ name }) => name),
      ['encrypt', 'decrypt', 'datakey', 'rewrap', 'hmac', 'verify'],
      'symmetric keys expose encryption but not signing or export'
    );
    key.supports_signing = true;
    key.exportable = true;
    assert.deepEqual(
      key.supportedActions.map(({ name }) => name),
      ['encrypt', 'decrypt', 'datakey', 'rewrap', 'sign', 'hmac', 'verify', 'export'],
      'signing and export are added when supported'
    );
    assert.deepEqual(key.exportKeyTypes, ['encryption', 'signing', 'hmac'], 'export order is preserved');
    key.supports_encryption = false;
    key.supports_decryption = false;
    assert.deepEqual(key.exportKeyTypes, ['signing', 'hmac'], 'signing is the default for signing keys');
    assert.deepEqual(
      key.supportedActions.map(({ name }) => name),
      ['sign', 'hmac', 'verify', 'export'],
      'signing-only keys omit encryption actions'
    );
  });

  test('version choices preserve archived gaps and update when limits change', function (assert) {
    const key = new TransitKeyResource({
      backend: 'transit',
      id: 'key',
      keys: { 3: 123, 4: 456, 5: 789 },
      latest_version: 5,
      min_decryption_version: 3,
      min_encryption_version: 0,
    });
    assert.deepEqual(key.validKeyVersions, ['3', '4', '5'], 'export uses only retained versions');
    assert.deepEqual(key.keyVersions, [1, 2, 3, 4, 5], 'configuration includes historical versions');
    assert.deepEqual(key.encryptionKeyVersions, [5, 4, 3], 'decryption minimum limits choices');
    assert.deepEqual(key.keysForEncryption, [5, 4, 3, 2, 1], 'zero minimum permits all versions');
    key.min_decryption_version = 4;
    key.min_encryption_version = 4;
    assert.deepEqual(key.encryptionKeyVersions, [5, 4], 'decryption choices update');
    assert.deepEqual(key.keysForEncryption, [5, 4], 'encryption choices update');
    key.min_encryption_version = 10;
    assert.deepEqual(key.keysForEncryption, [], 'an out-of-range minimum is clamped');
  });

  test('empty version data produces empty choices', function (assert) {
    const key = new TransitKeyResource({ backend: 'transit', id: 'key' });
    assert.deepEqual(key.keyVersions, [], 'no configuration versions are invented');
    assert.deepEqual(key.validKeyVersions, [], 'no export versions are invented');
    assert.deepEqual(key.keysForEncryption, [], 'no encryption versions are invented');
  });
});
