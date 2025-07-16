/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

const CHACHA = {
  request_id: 'a5695685-584c-6b25-fada-35304d3d583d',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    allow_plaintext_backup: false,
    deletion_allowed: false,
    derived: false,
    exportable: false,
    keys: {
      1: 1559598610,
    },
    latest_version: 1,
    min_available_version: 0,
    min_decryption_version: 1,
    min_encryption_version: 0,
    name: 'anewone',
    supports_decryption: true,
    supports_derivation: true,
    supports_encryption: true,
    supports_signing: false,
    type: 'chacha20-poly1305',
  },
  wrap_info: null,
  warnings: null,
  auth: null,
  backend: 'its-a-transit',
};

const AES = {
  request_id: '90c327a8-9a68-6fab-13a1-f51b68cb24d7',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    allow_plaintext_backup: false,
    deletion_allowed: false,
    derived: false,
    exportable: false,
    keys: {
      1: 1559577523,
    },
    latest_version: 1,
    min_available_version: 0,
    min_decryption_version: 1,
    min_encryption_version: 0,
    name: 'new',
    supports_decryption: true,
    supports_derivation: true,
    supports_encryption: true,
    supports_signing: false,
    type: 'aes256-gcm96',
  },
  wrap_info: null,
  warnings: null,
  auth: null,
  backend: 'its-a-transit',
};

module('Unit | Serializer | transit-key', function (hooks) {
  setupTest(hooks);
  test('it expands the timestamp for aes and chacha-poly keys', function (assert) {
    const serializer = this.owner.lookup('serializer:transit-key');
    const aesExpected = AES.data.keys[1] * 1000;
    const chachaExpected = CHACHA.data.keys[1] * 1000;
    const aesData = serializer.normalizeSecrets({ ...AES });
    assert.strictEqual(aesData[0].keys[1], aesExpected, 'converts seconds to millis for aes keys');

    const chachaData = serializer.normalizeSecrets({ ...CHACHA });
    assert.strictEqual(chachaData[0].keys[1], chachaExpected, 'converts seconds to millis for chacha keys');
  });

  test('it includes backend from the payload on the normalized data', function (assert) {
    const serializer = this.owner.lookup('serializer:transit-key');
    const data = serializer.normalizeSecrets({ ...AES });
    assert.strictEqual(data[0].backend, 'its-a-transit', 'pulls backend from the payload onto the data');
  });
});
