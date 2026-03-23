/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { getModelTypeForEngine } from 'vault/utils/model-helpers/secret-engine-helpers';

module('Unit | Utility | transform-engine-logic', function () {
  module('Secret prefix-based model type detection', function () {
    test('it returns "transform" when secret is empty/null/undefined', function (assert) {
      // Check that empty/null/undefined secrets return default transform type
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: null }),
        'transform',
        'returns transform for null secret'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: undefined }),
        'transform',
        'returns transform for undefined secret'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: '' }),
        'transform',
        'returns transform for empty string secret'
      );
    });

    test('it returns "transform/role" when secret starts with "role/"', function (assert) {
      // Check that role/ prefix returns transform/role
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'role/my-role' }),
        'transform/role',
        'returns transform/role for role/ prefixed secret'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'role/' }),
        'transform/role',
        'returns transform/role for just role/ prefix'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'role/test-role-name' }),
        'transform/role',
        'returns transform/role for complex role name'
      );
    });

    test('it returns "transform/template" when secret starts with "template/"', function (assert) {
      // Check that template/ prefix returns transform/template
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'template/my-template' }),
        'transform/template',
        'returns transform/template for template/ prefixed secret'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'template/' }),
        'transform/template',
        'returns transform/template for just template/ prefix'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'template/test-template-name' }),
        'transform/template',
        'returns transform/template for complex template name'
      );
    });

    test('it returns "transform/alphabet" when secret starts with "alphabet/"', function (assert) {
      // Check that alphabet/ prefix returns transform/alphabet
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'alphabet/my-alphabet' }),
        'transform/alphabet',
        'returns transform/alphabet for alphabet/ prefixed secret'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'alphabet/' }),
        'transform/alphabet',
        'returns transform/alphabet for just alphabet/ prefix'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'alphabet/test-alphabet-name' }),
        'transform/alphabet',
        'returns transform/alphabet for complex alphabet name'
      );
    });

    test('it returns "transform" as default for other secret names', function (assert) {
      // Check that non-recognized prefixes return default transform type
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'some-other-secret' }),
        'transform',
        'returns transform for non-prefixed secret'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'transformation/test' }),
        'transform',
        'returns transform for transformation/ prefix (TODO case)'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'random-name' }),
        'transform',
        'returns transform for random secret name'
      );
    });

    test('it handles edge cases correctly', function (assert) {
      // Test cases that might cause issues
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'role' }),
        'transform',
        'returns transform for just "role" (no slash)'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'template' }),
        'transform',
        'returns transform for just "template" (no slash)'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'alphabet' }),
        'transform',
        'returns transform for just "alphabet" (no slash)'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { secret: 'role/template/alphabet' }),
        'transform/role',
        'returns transform/role for complex path starting with role/'
      );
    });
  });

  module('Tab-based model type selection', function () {
    test('it returns correct model types based on tab parameter', function (assert) {
      // Check tab-based model type selection (switch statement logic)
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'role' }),
        'transform/role',
        'returns transform/role for role tab'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'template' }),
        'transform/template',
        'returns transform/template for template tab'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'alphabet' }),
        'transform/alphabet',
        'returns transform/alphabet for alphabet tab'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'other' }),
        'transform',
        'returns transform for unknown tab (default case)'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: null }),
        'transform',
        'returns transform for null tab'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: undefined }),
        'transform',
        'returns transform for undefined tab'
      );
    });
  });

  module('Context-based transform type resolution', function () {
    test('it handles tab context parameter', function (assert) {
      // Simplified to use only tab parameter
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'role' }),
        'transform/role',
        'uses tab when available'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'template' }),
        'transform/template',
        'uses tab when available'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'alphabet' }),
        'transform/alphabet',
        'uses tab when available'
      );
    });

    test('it validates tab is in allowed list', function (assert) {
      // Check that only allowed tab values return specific types
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'invalid' }),
        'transform',
        'returns default for invalid tab'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: 'transformation' }),
        'transform',
        'returns default for "transformation" (not in allowed list)'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { tab: '' }),
        'transform',
        'returns default for empty tab'
      );
    });

    test('it returns default when no valid transform type is found', function (assert) {
      // Check default behavior when no valid context is provided
      assert.strictEqual(
        getModelTypeForEngine('transform', {}),
        'transform',
        'returns default for empty context'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', { someOtherParam: 'value' }),
        'transform',
        'returns default when no transform-related parameters'
      );
    });
  });

  module('Combined logic scenarios', function () {
    test('secret parameter takes precedence over tab', function (assert) {
      // When secret is provided with a prefix, it should override tab
      assert.strictEqual(
        getModelTypeForEngine('transform', {
          secret: 'role/my-role',
          tab: 'template',
        }),
        'transform/role',
        'secret prefix overrides tab'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', {
          secret: 'template/my-template',
          tab: 'alphabet',
        }),
        'transform/template',
        'secret prefix overrides tab'
      );
    });

    test('tab used when secret has no recognized prefix', function (assert) {
      // When secret doesn't have a recognized prefix, fall back to tab
      assert.strictEqual(
        getModelTypeForEngine('transform', {
          secret: 'some-other-secret',
          tab: 'role',
        }),
        'transform/role',
        'uses tab when secret has no prefix'
      );
      assert.strictEqual(
        getModelTypeForEngine('transform', {
          secret: 'random-name',
          tab: 'template',
        }),
        'transform/template',
        'uses tab when secret has no prefix'
      );
    });

    test('handles all original edge cases', function (assert) {
      // Comprehensive test covering various combinations
      const testCases = [
        // Secret-based detection
        { input: { secret: 'role/test' }, expected: 'transform/role' },
        { input: { secret: 'template/test' }, expected: 'transform/template' },
        { input: { secret: 'alphabet/test' }, expected: 'transform/alphabet' },
        { input: { secret: 'other/test' }, expected: 'transform' },
        { input: { secret: '' }, expected: 'transform' },
        { input: { secret: null }, expected: 'transform' },

        // Tab-based selection
        { input: { tab: 'role' }, expected: 'transform/role' },
        { input: { tab: 'template' }, expected: 'transform/template' },
        { input: { tab: 'alphabet' }, expected: 'transform/alphabet' },

        // Default cases
        { input: {}, expected: 'transform' },
        { input: { tab: 'invalid' }, expected: 'transform' },

        // Precedence cases
        { input: { secret: 'role/test', tab: 'template' }, expected: 'transform/role' },
        { input: { secret: 'other', tab: 'role' }, expected: 'transform/role' },
      ];

      testCases.forEach(({ input, expected }) => {
        const result = getModelTypeForEngine('transform', input);
        assert.strictEqual(
          result,
          expected,
          `getModelTypeForEngine('transform', ${JSON.stringify(input)}) should return '${expected}'`
        );
      });
    });
  });
});
