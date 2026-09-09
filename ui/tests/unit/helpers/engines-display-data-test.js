/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import engineDisplayData, { unknownEngineMetadata } from 'core/helpers/engines-display-data';
import { module, test } from 'qunit';
import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';

module('Unit | Helper | engines-display-data', function () {
  test('it returns correct display data for known engine types', function (assert) {
    // Test keymgmt engine
    const keymgmtData = engineDisplayData('keymgmt');
    assert.strictEqual(keymgmtData.type, 'keymgmt', 'returns correct type for keymgmt');
    assert.strictEqual(keymgmtData.displayName, 'Key Management', 'returns correct displayName for keymgmt');
    assert.ok(keymgmtData.requiresEnterprise, 'keymgmt requires enterprise');

    // Test aws engine with ALL_ENGINES comparison
    const awsData = engineDisplayData('aws');
    const expectedAws = ALL_ENGINES.find((e) => e.type === 'aws');
    assert.propEqual(awsData, expectedAws, 'Returns correct display data for aws');

    // Test enterprise-only engine
    const kmipData = engineDisplayData('kmip');
    assert.true(kmipData.requiresEnterprise, 'KMIP requires enterprise');
    assert.strictEqual(kmipData.displayName, 'KMIP', 'KMIP displayName is correct');
  });

  test('it returns metadata for external plugins that map to builtins', function (assert) {
    const externalKeymgmtData = engineDisplayData('vault-plugin-secrets-keymgmt');

    // Should return keymgmt metadata but with the external plugin type preserved
    assert.strictEqual(
      externalKeymgmtData.type,
      'vault-plugin-secrets-keymgmt',
      'preserves external plugin type'
    );
    assert.strictEqual(externalKeymgmtData.displayName, 'Key Management', 'returns builtin displayName');
    assert.ok(externalKeymgmtData.requiresEnterprise, 'inherits enterprise requirement from builtin');
    assert.strictEqual(externalKeymgmtData.glyph, 'key', 'inherits glyph from builtin');
  });

  test('it returns unknown plugin metadata for unmapped external plugins', function (assert) {
    const unknownData = engineDisplayData('vault-plugin-secrets-unknown');
    const unknownMetadata = unknownEngineMetadata('vault-plugin-secrets-unknown');

    assert.strictEqual(unknownData.type, unknownMetadata.type, 'returns unknown type');
    assert.strictEqual(
      unknownData.displayName,
      'vault-plugin-secrets-unknown',
      'uses plugin name as displayName'
    );
    assert.strictEqual(unknownData.glyph, unknownMetadata.glyph, 'uses default lock glyph');
    assert.deepEqual(
      unknownData.mountCategory,
      unknownMetadata.mountCategory,
      'has correct mount categories'
    );
  });

  test('it returns unknown plugin metadata for empty/null inputs', function (assert) {
    const emptyData = engineDisplayData('');
    const nullData = engineDisplayData(null);
    const undefinedData = engineDisplayData(undefined);

    // Test empty string
    assert.strictEqual(emptyData.type, 'unknown', 'returns unknown type for empty string');
    assert.strictEqual(emptyData.displayName, 'Unknown plugin', 'uses default name for empty string');
    assert.propEqual(
      emptyData.mountCategory,
      ['secret', 'auth'],
      'mountCategory is correct for empty string'
    );
    assert.strictEqual(emptyData.glyph, 'lock', 'default glyph is a lock for empty string');

    // Test null
    assert.strictEqual(nullData.type, 'unknown', 'returns unknown type for null');
    assert.strictEqual(nullData.displayName, 'Unknown plugin', 'uses default name for null');
    assert.propEqual(nullData.mountCategory, ['secret', 'auth'], 'mountCategory is correct for null');
    assert.strictEqual(nullData.glyph, 'lock', 'default glyph is a lock for null');

    // Test undefined
    assert.strictEqual(undefinedData.type, 'unknown', 'returns unknown type for undefined');
    assert.strictEqual(undefinedData.displayName, 'Unknown plugin', 'uses default name for undefined');
    assert.propEqual(
      undefinedData.mountCategory,
      ['secret', 'auth'],
      'mountCategory is correct for undefined'
    );
    assert.strictEqual(undefinedData.glyph, 'lock', 'default glyph is a lock for undefined');
  });

  test('it returns fallback display data for unknown engine types', function (assert) {
    const unknownData = engineDisplayData('not-an-engine');
    assert.strictEqual(unknownData.displayName, 'not-an-engine', 'uses passed type as fallback displayName');
    assert.strictEqual(unknownData.type, 'not-an-engine', 'returns methodType type');
    assert.propEqual(unknownData.mountCategory, ['secret', 'auth'], 'mountCategory is correct');
    assert.strictEqual(unknownData.glyph, 'lock', 'default glyph is a lock');
  });

  test('it handles case sensitivity correctly', function (assert) {
    // Should not match due to case sensitivity
    const upperCaseData = engineDisplayData('KEYMGMT');
    const upperCaseUnknownMetadata = unknownEngineMetadata('KEYMGMT');

    const mixedCaseData = engineDisplayData('KeyMgmt');
    const mixedCaseUnknownMetadata = unknownEngineMetadata('KeyMgmt');

    assert.strictEqual(
      upperCaseData.type,
      upperCaseUnknownMetadata.type,
      'case sensitive - KEYMGMT not recognized'
    );
    assert.strictEqual(
      mixedCaseData.type,
      mixedCaseUnknownMetadata.type,
      'case sensitive - KeyMgmt not recognized'
    );
  });
});
