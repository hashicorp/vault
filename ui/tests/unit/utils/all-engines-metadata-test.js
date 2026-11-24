/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { ALL_ENGINES, filterEnginesByMountCategory, isAddonEngine } from 'vault/utils/all-engines-metadata';

module('Unit | Utility | all-engines-metadata', function () {
  module('ALL_ENGINES', function () {
    test('it contains expected engine metadata', function (assert) {
      assert.true(Array.isArray(ALL_ENGINES), 'ALL_ENGINES is an array');
      assert.true(ALL_ENGINES.length > 0, 'ALL_ENGINES contains engines');

      // Check that at least some expected engines are present
      const engineTypes = ALL_ENGINES.map((engine) => engine.type);
      assert.true(engineTypes.includes('kv'), 'contains kv engine');
      assert.true(engineTypes.includes('pki'), 'contains pki engine');
      assert.true(engineTypes.includes('transit'), 'contains transit engine');
    });

    test('all engines have required properties', function (assert) {
      ALL_ENGINES.forEach((engine) => {
        assert.ok(engine.displayName, `${engine.type} has displayName`);
        assert.ok(engine.type, `${engine.type} has type`);
        assert.true(Array.isArray(engine.mountCategory), `${engine.type} has mountCategory array`);
        assert.true(engine.mountCategory.length > 0, `${engine.type} has at least one mount category`);
      });
    });
  });

  module('filterEnginesByMountCategory', function () {
    test('filters engines by secret mount category', function (assert) {
      const secretEngines = filterEnginesByMountCategory({
        mountCategory: 'secret',
        isEnterprise: false,
      });

      assert.true(Array.isArray(secretEngines), 'returns an array');
      assert.true(secretEngines.length > 0, 'returns some engines');

      // All returned engines should have 'secret' in mountCategory
      secretEngines.forEach((engine) => {
        assert.true(
          engine.mountCategory.includes('secret'),
          `${engine.type} should have 'secret' in mountCategory`
        );
      });
    });

    test('filters engines by auth mount category', function (assert) {
      const authEngines = filterEnginesByMountCategory({
        mountCategory: 'auth',
        isEnterprise: false,
      });

      assert.true(Array.isArray(authEngines), 'returns an array');
      assert.true(authEngines.length > 0, 'returns some engines');

      // All returned engines should have 'auth' in mountCategory
      authEngines.forEach((engine) => {
        assert.true(
          engine.mountCategory.includes('auth'),
          `${engine.type} should have 'auth' in mountCategory`
        );
      });
    });

    test('excludes enterprise engines when isEnterprise is false', function (assert) {
      const ossEngines = filterEnginesByMountCategory({
        mountCategory: 'secret',
        isEnterprise: false,
      });

      // Should not contain any engines that require enterprise
      ossEngines.forEach((engine) => {
        assert.notOk(
          engine.requiresEnterprise,
          `${engine.type} should not require enterprise when isEnterprise is false`
        );
      });
    });

    test('includes enterprise engines when isEnterprise is true', function (assert) {
      const allEngines = filterEnginesByMountCategory({
        mountCategory: 'secret',
        isEnterprise: true,
      });

      const ossEngines = filterEnginesByMountCategory({
        mountCategory: 'secret',
        isEnterprise: false,
      });

      // Enterprise should have same or more engines than OSS
      assert.true(
        allEngines.length >= ossEngines.length,
        'enterprise mode should include same or more engines'
      );
    });
  });

  module('isAddonEngine', function () {
    test('returns false for kv version 1', function (assert) {
      assert.false(isAddonEngine('kv', 1), 'kv version 1 is not an addon engine');
    });

    test('returns true for engines with engineRoute', function (assert) {
      assert.true(isAddonEngine('kv', 2), 'kv version 2 is an addon engine');
      assert.true(isAddonEngine('pki', 1), 'pki is an addon engine');
    });

    test('returns false for engines without engineRoute', function (assert) {
      assert.false(isAddonEngine('transit', 1), 'transit is not an addon engine');
      assert.false(isAddonEngine('cubbyhole', 1), 'cubbyhole is not an addon engine');
    });

    test('returns false for unknown engine types', function (assert) {
      assert.false(isAddonEngine('unknown-engine', 1), 'unknown engines are not addon engines');
    });
  });
});
