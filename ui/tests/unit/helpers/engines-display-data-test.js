/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import engineDisplayData, { unknownEngineMetadata } from 'vault/helpers/engines-display-data';

module('Unit | Helper | engines-display-data', function () {
  test('it returns metadata for builtin engines', function (assert) {
    const keymgmtData = engineDisplayData('keymgmt');

    assert.strictEqual(keymgmtData.type, 'keymgmt', 'returns correct type for keymgmt');
    assert.strictEqual(keymgmtData.displayName, 'Key Management', 'returns correct displayName for keymgmt');
    assert.ok(keymgmtData.requiresEnterprise, 'keymgmt requires enterprise');
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
    const unknownMetadata = unknownEngineMetadata();

    assert.strictEqual(emptyData.type, unknownMetadata.type, 'returns unknown for empty string');
    assert.strictEqual(emptyData.displayName, 'Unknown plugin', 'uses default name for empty string');

    assert.strictEqual(nullData.type, unknownMetadata.type, 'returns unknown for null');
    assert.strictEqual(undefinedData.type, unknownMetadata.type, 'returns unknown for undefined');
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
