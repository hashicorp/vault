/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { getModelTypeForEngine } from 'vault/utils/model-helpers/secret-engine-helpers';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';

module('Unit | Utility | model-helpers/secret-engine-helpers', function () {
  module('getModelTypeForEngine', function () {
    test('returns correct model types for basic engines', function (assert) {
      assert.strictEqual(getModelTypeForEngine('transit'), 'transit-key');
      assert.strictEqual(getModelTypeForEngine('ssh'), 'role-ssh');
      assert.strictEqual(getModelTypeForEngine('aws'), 'role-aws');
      assert.strictEqual(getModelTypeForEngine('cubbyhole'), 'secret');
      assert.strictEqual(getModelTypeForEngine('kv'), 'secret');
      assert.strictEqual(getModelTypeForEngine('generic'), 'secret');
      assert.strictEqual(getModelTypeForEngine('totp'), 'totp-key');
    });

    test('returns correct model types for database engine with context', function (assert) {
      assert.strictEqual(
        getModelTypeForEngine('database', { isRole: true }),
        'database/role',
        'returns database/role when isRole is true'
      );
      assert.strictEqual(
        getModelTypeForEngine('database', { tab: 'role' }),
        'database/role',
        'returns database/role when tab is role'
      );
      assert.strictEqual(
        getModelTypeForEngine('database', { secret: 'role/my-role' }),
        'database/role',
        'returns database/role when secret starts with role/'
      );
      assert.strictEqual(
        getModelTypeForEngine('database', {}),
        'database/connection',
        'returns database/connection for empty context'
      );
      assert.strictEqual(
        getModelTypeForEngine('database'),
        'database/connection',
        'returns database/connection with no context'
      );
    });

    test('returns correct model types for transform engine', function (assert) {
      // Test secret name prefix logic (takes priority)
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'role/my-role' }),
        'transform/role',
        'returns transform/role for secret starting with role/'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'template/my-template' }),
        'transform/template',
        'returns transform/template for secret starting with template/'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'alphabet/my-alphabet' }),
        'transform/alphabet',
        'returns transform/alphabet for secret starting with alphabet/'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'other/my-other' }),
        'transform',
        'returns transform for secret with unknown prefix'
      );

      // Test query parameter logic (fallback)
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'role' }),
        'transform/role',
        'returns transform/role when tab is role'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { transformType: 'template' }),
        'transform/template',
        'returns transform/template when transformType is template'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'alphabet' }),
        'transform/alphabet',
        'returns transform/alphabet when tab is alphabet'
      );

      // Test precedence: secret name should override query params
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'role/my-role', tab: 'template' }),
        'transform/role',
        'secret name prefix takes precedence over tab parameter'
      );

      // Test fallback cases
      assert.strictEqual(
        getModelTypeForEngine('transform', { transformType: 'unknown' }),
        'transform',
        'returns transform for unknown transformType'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', {}),
        'transform',
        'returns transform for empty context'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform'),
        'transform',
        'returns transform with no context'
      );
    });

    test('returns correct model types for keymgmt engine', function (assert) {
      assert.strictEqual(
        getModelTypeForEngine('keymgmt', { itemType: 'key' }),
        'keymgmt/key',
        'returns keymgmt/key when itemType is key'
      );
      assert.strictEqual(
        getModelTypeForEngine('keymgmt', { tab: 'provider' }),
        'keymgmt/provider',
        'returns keymgmt/provider when tab is provider'
      );
      assert.strictEqual(
        getModelTypeForEngine('keymgmt', { itemType: 'provider' }),
        'keymgmt/provider',
        'returns keymgmt/provider when itemType is provider'
      );
      assert.strictEqual(
        getModelTypeForEngine('keymgmt', {}),
        'keymgmt/key',
        'returns keymgmt/key for empty context (default)'
      );
      assert.strictEqual(
        getModelTypeForEngine('keymgmt'),
        'keymgmt/key',
        'returns keymgmt/key with no context (default)'
      );
    });

    test('returns default "secret" for unknown engines', function (assert) {
      assert.strictEqual(
        getModelTypeForEngine('unknown-engine'),
        'secret',
        'returns secret for unknown engine'
      );
      assert.strictEqual(
        getModelTypeForEngine('custom-plugin'),
        'secret',
        'returns secret for custom plugin'
      );
      assert.strictEqual(getModelTypeForEngine(''), 'secret', 'returns secret for empty string');
    });

    test('works with external plugin mapping', function (assert) {
      // Test that external plugins get correct model types via effective type mapping
      const externalKeymgmtType = getEffectiveEngineType('vault-plugin-secrets-keymgmt');
      assert.strictEqual(externalKeymgmtType, 'keymgmt', 'external plugin maps to builtin');

      const modelType = getModelTypeForEngine(externalKeymgmtType, { itemType: 'provider' });
      assert.strictEqual(modelType, 'keymgmt/provider', 'external plugin gets correct model type');
    });
  });
});
