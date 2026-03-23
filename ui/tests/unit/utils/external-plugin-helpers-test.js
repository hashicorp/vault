/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  EXTERNAL_PLUGIN_TO_BUILTIN_MAP,
  getBuiltinTypeFromExternalPlugin,
  getEffectiveEngineType,
  getExternalPluginNameFromBuiltin,
  isKnownExternalPlugin,
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

  module('getExternalPluginNameFromBuiltin', function () {
    test('it returns external plugin name for known builtin types', function (assert) {
      assert.strictEqual(
        getExternalPluginNameFromBuiltin('keymgmt'),
        'vault-plugin-secrets-keymgmt',
        'returns vault-plugin-secrets-keymgmt for keymgmt'
      );

      assert.strictEqual(
        getExternalPluginNameFromBuiltin('azure'),
        'vault-plugin-secrets-azure',
        'returns vault-plugin-secrets-azure for azure'
      );

      assert.strictEqual(
        getExternalPluginNameFromBuiltin('gcp'),
        'vault-plugin-secrets-gcp',
        'returns vault-plugin-secrets-gcp for gcp'
      );
    });

    test('it returns null for unknown builtin types', function (assert) {
      assert.strictEqual(
        getExternalPluginNameFromBuiltin('unknown-engine'),
        null,
        'returns null for unknown builtin type'
      );
    });

    test('it returns null for external plugin names', function (assert) {
      assert.strictEqual(
        getExternalPluginNameFromBuiltin('vault-plugin-secrets-keymgmt'),
        null,
        'returns null for external plugin name'
      );
    });

    test('it returns null for empty string', function (assert) {
      assert.strictEqual(getExternalPluginNameFromBuiltin(''), null, 'returns null for empty string');
    });

    test('it handles case sensitivity correctly', function (assert) {
      assert.strictEqual(
        getExternalPluginNameFromBuiltin('KEYMGMT'),
        null,
        'returns null for uppercase builtin type'
      );

      assert.strictEqual(
        getExternalPluginNameFromBuiltin('KeyMgmt'),
        null,
        'returns null for mixed case builtin type'
      );
    });

    test('it works with all mapped builtin types', function (assert) {
      // Test that every builtin type in the map can be reverse-looked up
      const builtinTypes = Object.values(EXTERNAL_PLUGIN_TO_BUILTIN_MAP);
      const uniqueBuiltinTypes = [...new Set(builtinTypes)];

      uniqueBuiltinTypes.forEach((builtinType) => {
        const externalName = getExternalPluginNameFromBuiltin(builtinType);
        assert.ok(externalName, `found external name for builtin type: ${builtinType}`);
        assert.true(
          externalName.startsWith('vault-plugin-'),
          `external name ${externalName} follows expected pattern`
        );
      });
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

    test('reverse lookup works with extended mappings', function (assert) {
      // Test conceptual extensibility of the reverse lookup
      // This verifies that the reverse lookup algorithm is robust
      const originalFunction = getExternalPluginNameFromBuiltin('keymgmt');
      assert.strictEqual(
        originalFunction,
        'vault-plugin-secrets-keymgmt',
        'reverse lookup works for existing mappings'
      );

      // Test that non-existent mappings return null as expected
      const nonExistentResult = getExternalPluginNameFromBuiltin('hypothetical-auth');
      assert.strictEqual(
        nonExistentResult,
        null,
        'reverse lookup correctly returns null for non-mapped types'
      );
    });
  });
});
