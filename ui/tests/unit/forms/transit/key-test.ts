/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import TransitKeyForm from 'vault/forms/transit/key';
import { TransitCreateKeyRequestTypeEnum } from '@hashicorp/vault-client-typescript';

module('Unit | Form | transit key', function () {
  test('convergent encryption requires derivation, and disabling derivation disables convergence', function (assert) {
    const form = new TransitKeyForm({ name: 'key' }, { isNew: true });
    form.convergentEncryptionValue = true;
    assert.true(form.data.derived, 'enabling convergent encryption enables derivation');
    assert.true(form.data.convergent_encryption, 'convergent encryption is enabled');
    form.derivedValue = false;
    assert.false(form.data.derived, 'derivation is disabled');
    assert.false(form.data.convergent_encryption, 'convergent encryption cannot remain enabled');
  });

  test('creation fields depend on the key algorithm', function (assert) {
    for (const type of [
      TransitCreateKeyRequestTypeEnum.RSA_2048,
      TransitCreateKeyRequestTypeEnum.ECDSA_P256,
      TransitCreateKeyRequestTypeEnum.ED25519,
      TransitCreateKeyRequestTypeEnum.AES256_GCM96,
    ]) {
      const form = new TransitKeyForm({ name: 'key', type }, { isNew: true });
      const names = form.formFields.map((field) => field.name);
      assert.strictEqual(names.includes('derived'), ['ed25519', 'aes256-gcm96'].includes(type), type);
      assert.strictEqual(names.includes('convergent_encryption'), type === 'aes256-gcm96', type);
    }
  });

  test('edit version choices react to unsaved decryption minimum changes', function (assert) {
    const form = new TransitKeyForm({
      name: 'key',
      keys: { 1: 123, 2: 456, 3: 789 },
      min_decryption_version: 1,
    });
    assert.deepEqual(form.keyVersions, [1, 2, 3], 'all versions are available for decryption');
    assert.deepEqual(form.encryptionKeyVersions, [3, 2, 1], 'initial encryption choices');
    form.data.min_decryption_version = 2;
    assert.deepEqual(form.encryptionKeyVersions, [3, 2], 'unsaved changes constrain encryption choices');
    assert.false('encryptionKeyVersions' in form.toJSON().data, 'derived choices are not serialized');
  });

  test('serialization validates names and the minimum auto-rotation duration', function (assert) {
    const unnamed = new TransitKeyForm({}, { isNew: true });
    assert.false(unnamed.toJSON().isValid, 'a key name is required');
    for (const [period, expected] of [
      ['0s', true],
      ['3599s', false],
      ['1h', true],
      ['24h', true],
    ] as const) {
      const form = new TransitKeyForm({ name: 'key', auto_rotate_period: period }, { isNew: true });
      assert.strictEqual(form.toJSON().isValid, expected, `auto-rotation period ${period}`);
      assert.strictEqual(form.toJSON().data['auto_rotate_period'], period, 'duration is not mutated');
    }
  });
});
