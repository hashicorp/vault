/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';
import { CERTIFICATES } from 'vault/tests/helpers/pki/pki-helpers';

const { rootPem } = CERTIFICATES;

module('Unit | Serializer | pki/action', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    const store = this.owner.lookup('service:store');
    const serializer = store.serializerFor('pki/action');

    assert.ok(serializer);
  });

  module('actionType import', function (hooks) {
    hooks.beforeEach(function () {
      this.actionType = 'import';
      this.pemBundle = rootPem;
    });

    test('it serializes only valid params', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        pemBundle: this.pemBundle,
        issuerRef: 'do-not-send',
        keyType: 'do-not-send',
      });
      const expectedResult = {
        pem_bundle: this.pemBundle,
      };

      const serializedRecord = record.serialize(this.actionType);
      assert.deepEqual(
        serializedRecord,
        expectedResult,
        'Serializes only parameters valid for import action'
      );
    });
  });

  module('actionType generate-root', function (hooks) {
    hooks.beforeEach(function () {
      this.actionType = 'generate-root';
      this.allKeyFields = {
        keyName: 'key name',
        keyType: 'rsa',
        keyBits: '0',
        keyRef: 'key ref',
        managedKeyName: 'managed name',
        managedKeyId: 'managed id',
      };
      this.withDefaults = {
        exclude_cn_from_sans: false,
        format: 'pem',
        max_path_length: -1,
        not_before_duration: '30s',
        private_key_format: 'der',
      };
    });

    test('it serializes only params with values', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        excludeCnFromSans: false,
        format: 'pem',
        maxPathLength: -1,
        notBeforeDuration: '30s',
        privateKeyFormat: 'der',
        type: 'external', // only used for endpoint in adapter
        customTtl: '40m', // UI-only value
        issuerName: 'my issuer',
        commonName: undefined,
        foo: 'bar',
      });
      const expectedResult = {
        ...this.withDefaults,
        key_bits: '0',
        key_ref: 'default',
        key_type: 'rsa',
        issuer_name: 'my issuer',
      };

      // without passing `actionType` it will not compare against an allowlist
      const serializedRecord = record.serialize();
      assert.deepEqual(serializedRecord, expectedResult);
    });

    test('it serializes only valid params for type = external', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        ...this.allKeyFields,
        type: 'external',
        customTtl: '40m',
        issuerName: 'my issuer',
        commonName: 'my common name',
      });
      const expectedResult = {
        ...this.withDefaults,
        issuer_name: 'my issuer',
        common_name: 'my common name',
        key_name: 'key name',
        key_type: 'rsa',
        key_bits: '0',
      };

      const serializedRecord = record.serialize(this.actionType);
      assert.deepEqual(serializedRecord, expectedResult);
    });

    test('it serializes only valid params for type = internal', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        ...this.allKeyFields,
        type: 'internal',
        customTtl: '40m',
        issuerName: 'my issuer',
        commonName: 'my common name',
      });
      const expectedResult = {
        ...this.withDefaults,
        issuer_name: 'my issuer',
        common_name: 'my common name',
        key_name: 'key name',
        key_type: 'rsa',
        key_bits: '0',
      };

      const serializedRecord = record.serialize(this.actionType);
      assert.deepEqual(serializedRecord, expectedResult);
    });

    test('it serializes only valid params for type = existing', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        ...this.allKeyFields,
        type: 'existing',
        customTtl: '40m',
        issuerName: 'my issuer',
        commonName: 'my common name',
      });
      const expectedResult = {
        ...this.withDefaults,
        issuer_name: 'my issuer',
        common_name: 'my common name',
        key_ref: 'key ref',
      };

      const serializedRecord = record.serialize(this.actionType);
      assert.deepEqual(serializedRecord, expectedResult);
    });

    test('it serializes only valid params for type = kms', function (assert) {
      const store = this.owner.lookup('service:store');
      const record = store.createRecord('pki/action', {
        ...this.allKeyFields,
        type: 'kms',
        customTtl: '40m',
        issuerName: 'my issuer',
        commonName: 'my common name',
      });
      const expectedResult = {
        ...this.withDefaults,
        issuer_name: 'my issuer',
        common_name: 'my common name',
        key_name: 'key name',
        managed_key_name: 'managed name',
        managed_key_id: 'managed id',
      };

      const serializedRecord = record.serialize(this.actionType);
      assert.deepEqual(serializedRecord, expectedResult);
    });
  });
});
