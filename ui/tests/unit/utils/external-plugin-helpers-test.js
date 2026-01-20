/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  EXTERNAL_PLUGIN_TO_BUILTIN_MAP,
  getBuiltinTypeFromExternalPlugin,
  isKnownExternalPlugin,
  getEffectiveEngineType,
} from 'vault/utils/external-plugin-helpers';

module('Unit | Utility | external-plugin-helpers', function () {
  module('EXTERNAL_PLUGIN_TO_BUILTIN_MAP', function () {
    test('it contains expected mappings', function (assert) {
      assert.strictEqual(
        EXTERNAL_PLUGIN_TO_BUILTIN_MAP['vault-plugin-secrets-keymgmt'],
        'keymgmt',
        'maps vault-plugin-secrets-keymgmt to keymgmt'
      );
    });

    test('it is a constant record', function (assert) {
      assert.strictEqual(typeof EXTERNAL_PLUGIN_TO_BUILTIN_MAP, 'object', 'is an object');
      assert.notStrictEqual(EXTERNAL_PLUGIN_TO_BUILTIN_MAP, null, 'is not null');
    });
  });

  module('getBuiltinTypeFromExternalPlugin', function () {
    test('it returns mapped builtin type for known external plugins', function (assert) {
      assert.strictEqual(
        getBuiltinTypeFromExternalPlugin('vault-plugin-secrets-keymgmt'),
        'keymgmt',
        'returns keymgmt for vault-plugin-secrets-keymgmt'
      );
    });

    test('it returns undefined for unknown external plugins', function (assert) {
      assert.strictEqual(
        getBuiltinTypeFromExternalPlugin('vault-plugin-secrets-unknown'),
        undefined,
        'returns undefined for unknown plugin'
      );
    });

    test('it returns undefined for builtin plugin names', function (assert) {
      assert.strictEqual(
        getBuiltinTypeFromExternalPlugin('keymgmt'),
        undefined,
        'returns undefined for builtin plugin name'
      );
    });

    test('it returns undefined for empty string', function (assert) {
      assert.strictEqual(
        getBuiltinTypeFromExternalPlugin(''),
        undefined,
        'returns undefined for empty string'
      );
    });
  });

  module('isKnownExternalPlugin', function () {
    test('it returns true for known external plugins', function (assert) {
      assert.true(
        isKnownExternalPlugin('vault-plugin-secrets-keymgmt'),
        'returns true for vault-plugin-secrets-keymgmt'
      );
    });

    test('it returns false for unknown external plugins', function (assert) {
      assert.false(isKnownExternalPlugin('vault-plugin-secrets-unknown'), 'returns false for unknown plugin');
    });

    test('it returns false for builtin plugin names', function (assert) {
      assert.false(isKnownExternalPlugin('keymgmt'), 'returns false for builtin plugin name');
    });

    test('it returns false for empty string', function (assert) {
      assert.false(isKnownExternalPlugin(''), 'returns false for empty string');
    });
  });

  module('getEffectiveEngineType', function () {
    test('it returns builtin type for known external plugins', function (assert) {
      assert.strictEqual(
        getEffectiveEngineType('vault-plugin-secrets-keymgmt'),
        'keymgmt',
        'returns keymgmt for vault-plugin-secrets-keymgmt'
      );
    });

    test('it returns original type for unknown external plugins', function (assert) {
      assert.strictEqual(
        getEffectiveEngineType('vault-plugin-secrets-unknown'),
        'vault-plugin-secrets-unknown',
        'returns original type for unknown plugin'
      );
    });

    test('it returns original type for builtin plugins', function (assert) {
      assert.strictEqual(
        getEffectiveEngineType('keymgmt'),
        'keymgmt',
        'returns original type for builtin plugin'
      );
    });

    test('it returns original type for standard engines', function (assert) {
      assert.strictEqual(getEffectiveEngineType('kv'), 'kv', 'returns original type for kv engine');
      assert.strictEqual(getEffectiveEngineType('pki'), 'pki', 'returns original type for pki engine');
      assert.strictEqual(getEffectiveEngineType('aws'), 'aws', 'returns original type for aws engine');
    });

    test('it handles empty string gracefully', function (assert) {
      assert.strictEqual(getEffectiveEngineType(''), '', 'returns empty string for empty input');
    });
  });

  module('future extensibility', function () {
    test('mapping can be easily extended', function (assert) {
      // Test that we can add more mappings (conceptually)
      const testMap = {
        ...EXTERNAL_PLUGIN_TO_BUILTIN_MAP,
        'vault-plugin-auth-example': 'example-auth',
      };

      assert.strictEqual(testMap['vault-plugin-secrets-keymgmt'], 'keymgmt', 'existing mapping is preserved');
      assert.strictEqual(testMap['vault-plugin-auth-example'], 'example-auth', 'new mapping can be added');
    });
  });
});
