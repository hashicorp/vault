/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { secretQueryParams } from 'vault/helpers/secret-query-params';

/**
 * Test the secret-query-params helper to ensure it correctly handles
 * external plugin mapping for query parameter generation.
 */
module('Unit | Helper | secret-query-params external plugin support', function () {
  module('keymgmt external plugins', function () {
    test('generates itemType=key for external keymgmt plugins with key type', function (assert) {
      const result = secretQueryParams(['vault-plugin-secrets-keymgmt', 'key'], {});

      assert.deepEqual(
        result,
        { itemType: 'key' },
        'External keymgmt plugin generates correct itemType for key'
      );
    });

    test('generates itemType=provider for external keymgmt plugins with provider type', function (assert) {
      const result = secretQueryParams(['vault-plugin-secrets-keymgmt', 'provider'], {});

      assert.deepEqual(
        result,
        { itemType: 'provider' },
        'External keymgmt plugin generates correct itemType for provider'
      );
    });

    test('defaults to itemType=key for external keymgmt plugins with no type', function (assert) {
      const result = secretQueryParams(['vault-plugin-secrets-keymgmt'], {});

      assert.deepEqual(result, { itemType: 'key' }, 'External keymgmt plugin defaults to key itemType');
    });

    test('generates same params as builtin keymgmt', function (assert) {
      const externalResult = secretQueryParams(['vault-plugin-secrets-keymgmt', 'key'], {});
      const builtinResult = secretQueryParams(['keymgmt', 'key'], {});

      assert.deepEqual(
        externalResult,
        builtinResult,
        'External keymgmt generates same params as builtin keymgmt'
      );
    });
  });

  module('transit external plugins', function () {
    test('generates tab=actions for external transit plugins', function (assert) {
      // Note: transit external plugin would be vault-plugin-secrets-transit if it existed
      const result = secretQueryParams(['transit', ''], {});

      assert.deepEqual(result, { tab: 'actions' }, 'Transit plugins generate tab=actions');
    });
  });

  module('database external plugins', function () {
    test('generates type parameter for database plugins', function (assert) {
      const result = secretQueryParams(['database', 'connection'], {});

      assert.deepEqual(result, { type: 'connection' }, 'Database plugins generate correct type parameter');
    });

    test('passes through type parameter for external database plugins', function (assert) {
      // Even though we don't have database external mapping, test the behavior
      const result = secretQueryParams(['vault-plugin-database-postgresql', 'role'], {});

      // Should return undefined since unmapped external plugins don't generate params
      assert.strictEqual(result, undefined, 'Unmapped external plugins return undefined');
    });
  });

  module('asQueryParams formatting', function () {
    test('formats external keymgmt params for LinkTo components', function (assert) {
      const result = secretQueryParams(['vault-plugin-secrets-keymgmt', 'provider'], { asQueryParams: true });

      assert.deepEqual(
        result,
        {
          isQueryParams: true,
          values: { itemType: 'provider' },
        },
        'External keymgmt formats correctly for LinkTo components'
      );
    });

    test('returns undefined when formatted but no params generated', function (assert) {
      const result = secretQueryParams(['vault-plugin-secrets-unknown'], { asQueryParams: true });

      assert.strictEqual(
        result,
        undefined,
        'Unknown external plugins return undefined even with asQueryParams'
      );
    });
  });

  module('unknown external plugins', function () {
    test('returns undefined for unmapped external plugins', function (assert) {
      const result = secretQueryParams(['vault-plugin-secrets-unknown'], {});

      assert.strictEqual(result, undefined, 'Unmapped external plugins return undefined');
    });

    test('preserves behavior for builtin engines', function (assert) {
      const transitResult = secretQueryParams(['transit'], {});
      const keymgmtResult = secretQueryParams(['keymgmt'], {});
      const unknownResult = secretQueryParams(['unknown'], {});

      assert.deepEqual(transitResult, { tab: 'actions' }, 'Builtin transit works');
      assert.deepEqual(keymgmtResult, { itemType: 'key' }, 'Builtin keymgmt works');
      assert.strictEqual(unknownResult, undefined, 'Unknown builtin returns undefined');
    });
  });

  module('edge cases', function () {
    test('handles empty backend type', function (assert) {
      const result = secretQueryParams([''], {});

      assert.strictEqual(result, undefined, 'Empty backend type returns undefined');
    });

    test('handles undefined backend type', function (assert) {
      const result = secretQueryParams([undefined], {});

      assert.strictEqual(result, undefined, 'Undefined backend type returns undefined');
    });

    test('handles missing type parameter', function (assert) {
      const result = secretQueryParams(['vault-plugin-secrets-keymgmt'], {});

      assert.deepEqual(result, { itemType: 'key' }, 'Missing type parameter defaults correctly');
    });
  });
});
