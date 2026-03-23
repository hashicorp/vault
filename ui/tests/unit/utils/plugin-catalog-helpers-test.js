/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  categorizeEnginesByStatus,
  enhanceEnginesWithCatalogData,
  getAllVersionsForEngineType,
  MOUNT_CATEGORIES,
  PLUGIN_CATEGORIES,
  PLUGIN_TYPES,
} from 'vault/utils/plugin-catalog-helpers';

module('Unit | Utility | plugin-catalog-helpers', function () {
  module('enhanceEnginesWithCatalogData', function () {
    test('it returns original engines when no catalog data provided', function (assert) {
      const staticEngines = [
        {
          type: 'kv',
          displayName: 'KV',
          pluginCategory: PLUGIN_CATEGORIES.GENERIC,
          mountCategory: [MOUNT_CATEGORIES.SECRET],
        },
      ];

      const result = enhanceEnginesWithCatalogData(staticEngines, []);
      assert.deepEqual(result, staticEngines, 'returns original engines when no catalog data');
    });

    test('it returns original engines when catalog data is null', function (assert) {
      const staticEngines = [
        {
          type: 'kv',
          displayName: 'KV',
          pluginCategory: PLUGIN_CATEGORIES.GENERIC,
          mountCategory: [MOUNT_CATEGORIES.SECRET],
        },
      ];

      const result = enhanceEnginesWithCatalogData(staticEngines, null);
      assert.deepEqual(result, staticEngines, 'handles null catalog data gracefully');
    });

    test('it enhances existing engines with catalog data', function (assert) {
      const staticEngines = [
        {
          type: 'kv',
          displayName: 'KV',
          pluginCategory: PLUGIN_CATEGORIES.GENERIC,
          mountCategory: [MOUNT_CATEGORIES.SECRET],
        },
        {
          type: 'pki',
          displayName: 'PKI',
          pluginCategory: PLUGIN_CATEGORIES.GENERIC,
          mountCategory: [MOUNT_CATEGORIES.SECRET],
        },
      ];

      const catalogData = [
        {
          name: 'kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
          version: '1.0.0',
          deprecation_status: 'supported',
        },
      ];

      const result = enhanceEnginesWithCatalogData(staticEngines, catalogData);

      assert.strictEqual(result.length, 2, 'maintains original engine count');
      assert.strictEqual(result[0].type, 'kv', 'preserves engine type');
      assert.true(result[0].builtin, 'adds builtin flag from catalog');
      assert.strictEqual(result[0].version, '1.0.0', 'adds version from catalog');
      assert.strictEqual(result[0].deprecationStatus, 'supported', 'adds deprecation status');
      assert.true(result[0].isAvailable, 'marks as available when in catalog');
      assert.ok(result[0].pluginData, 'includes plugin data');

      assert.false(result[1].isAvailable, 'marks as unavailable when not in catalog');
      assert.notOk(result[1].builtin, 'no builtin flag when not in catalog');
    });

    test('it handles database engine specially', function (assert) {
      const staticEngines = [
        {
          type: MOUNT_CATEGORIES.DATABASE,
          displayName: 'Database',
          pluginCategory: PLUGIN_CATEGORIES.GENERIC,
          mountCategory: [MOUNT_CATEGORIES.SECRET],
        },
      ];

      const databasePlugins = [
        {
          name: 'mysql-database-plugin',
          type: PLUGIN_TYPES.DATABASE,
          builtin: true,
          version: '1.0.0',
        },
      ];

      const result = enhanceEnginesWithCatalogData(staticEngines, [], databasePlugins);

      assert.true(result[0].isAvailable, 'database engine is available when database plugins exist');
      assert.true(result[0].builtin, 'uses representative database plugin data');
      assert.strictEqual(result[0].version, '1.0.0', 'uses representative plugin version');
    });

    test('it discovers external plugins not in static metadata', function (assert) {
      const staticEngines = [
        {
          type: 'kv',
          displayName: 'KV',
          pluginCategory: PLUGIN_CATEGORIES.GENERIC,
          mountCategory: [MOUNT_CATEGORIES.SECRET],
        },
      ];

      const catalogData = [
        {
          name: 'kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
        },
        {
          name: 'my-custom-plugin',
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
          version: '2.1.0',
        },
      ];

      const result = enhanceEnginesWithCatalogData(staticEngines, catalogData);

      assert.strictEqual(result.length, 2, 'adds external plugin');
      const externalPlugin = result.find((engine) => engine.type === 'my-custom-plugin');
      assert.ok(externalPlugin, 'external plugin is present');
      assert.strictEqual(externalPlugin.displayName, 'My Custom Plugin', 'converts name to Title Case');
      assert.strictEqual(
        externalPlugin.pluginCategory,
        PLUGIN_CATEGORIES.EXTERNAL,
        'marks as external category'
      );
      assert.true(externalPlugin.isAvailable, 'external plugin is available');
      assert.false(externalPlugin.builtin, 'external plugin is not builtin');
    });

    test('it excludes external plugins with builtin mappings from external category', function (assert) {
      const staticEngines = [
        {
          type: 'kv',
          displayName: 'KV',
          pluginCategory: PLUGIN_CATEGORIES.GENERIC,
          mountCategory: [MOUNT_CATEGORIES.SECRET],
        },
      ];

      const catalogData = [
        {
          name: 'kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
        },
        {
          name: 'vault-plugin-secrets-kv', // This has a builtin mapping
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
          version: '2.1.0',
        },
        {
          name: 'truly-external-plugin', // This does not have a builtin mapping
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
          version: '1.0.0',
        },
      ];

      const result = enhanceEnginesWithCatalogData(staticEngines, catalogData);

      // Should only add the truly external plugin, not the one with builtin mapping
      assert.strictEqual(result.length, 2, 'adds only truly external plugin');

      const kvEngine = result.find((engine) => engine.type === 'kv');
      const externalKv = result.find((engine) => engine.type === 'vault-plugin-secrets-kv');
      const trulyExternal = result.find((engine) => engine.type === 'truly-external-plugin');

      assert.ok(kvEngine, 'KV engine is present');
      assert.notOk(externalKv, 'external KV plugin is not added as separate engine');
      assert.ok(trulyExternal, 'truly external plugin is present');
      assert.strictEqual(
        trulyExternal.pluginCategory,
        PLUGIN_CATEGORIES.EXTERNAL,
        'truly external plugin is in external category'
      );
    });

    test('it matches external plugins with existing static engine glyphs', function (assert) {
      const staticEngines = [
        {
          type: 'aws',
          displayName: 'AWS',
          glyph: 'aws-color',
          pluginCategory: PLUGIN_CATEGORIES.CLOUD,
          mountCategory: [MOUNT_CATEGORIES.SECRET],
        },
      ];

      const catalogData = [
        {
          name: 'my-custom-aws-plugin',
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
        },
      ];

      const result = enhanceEnginesWithCatalogData(staticEngines, catalogData);

      const externalPlugin = result.find((engine) => engine.type === 'my-custom-aws-plugin');
      assert.strictEqual(externalPlugin.glyph, 'aws-color', 'uses matching static engine glyph');
    });

    test('it handles deprecation status correctly', function (assert) {
      const staticEngines = [
        {
          type: 'legacy-plugin',
          displayName: 'Legacy Plugin',
          pluginCategory: PLUGIN_CATEGORIES.GENERIC,
          mountCategory: [MOUNT_CATEGORIES.SECRET],
        },
      ];

      const catalogData = [
        {
          name: 'legacy-plugin',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
          deprecation_status: 'deprecated',
        },
      ];

      const result = enhanceEnginesWithCatalogData(staticEngines, catalogData);

      assert.strictEqual(
        result[0].deprecationStatus,
        'deprecated',
        'preserves deprecation status from catalog'
      );
    });
  });

  module('categorizeEnginesByStatus', function () {
    test('it separates enabled and disabled engines', function (assert) {
      const engines = [
        {
          type: 'kv',
          displayName: 'KV',
          isAvailable: true,
        },
        {
          type: 'pki',
          displayName: 'PKI',
          isAvailable: false,
        },
        {
          type: 'aws',
          displayName: 'AWS',
          // isAvailable not set (undefined)
        },
      ];

      const result = categorizeEnginesByStatus(engines);

      assert.strictEqual(result.enabled.length, 2, 'enabled includes available and undefined');
      assert.strictEqual(result.disabled.length, 1, 'disabled includes only false');
      assert.strictEqual(result.enabled[0].type, 'kv', 'includes available engine in enabled');
      assert.strictEqual(result.enabled[1].type, 'aws', 'includes undefined availability as enabled');
      assert.strictEqual(result.disabled[0].type, 'pki', 'includes unavailable engine in disabled');
    });

    test('it handles empty input', function (assert) {
      const result = categorizeEnginesByStatus([]);

      assert.strictEqual(result.enabled.length, 0, 'enabled is empty');
      assert.strictEqual(result.disabled.length, 0, 'disabled is empty');
    });

    test('it handles all enabled engines', function (assert) {
      const engines = [
        { type: 'kv', isAvailable: true },
        { type: 'pki' }, // undefined isAvailable
      ];

      const result = categorizeEnginesByStatus(engines);

      assert.strictEqual(result.enabled.length, 2, 'all engines are enabled');
      assert.strictEqual(result.disabled.length, 0, 'no engines are disabled');
    });

    test('it handles all disabled engines', function (assert) {
      const engines = [
        { type: 'kv', isAvailable: false },
        { type: 'pki', isAvailable: false },
      ];

      const result = categorizeEnginesByStatus(engines);

      assert.strictEqual(result.enabled.length, 0, 'no engines are enabled');
      assert.strictEqual(result.disabled.length, 2, 'all engines are disabled');
    });
  });

  module('constants', function () {
    test('MOUNT_CATEGORIES contains expected values', function (assert) {
      assert.strictEqual(MOUNT_CATEGORIES.SECRET, 'secret', 'SECRET category is correct');
      assert.strictEqual(MOUNT_CATEGORIES.AUTH, 'auth', 'AUTH category is correct');
      assert.strictEqual(MOUNT_CATEGORIES.DATABASE, 'database', 'DATABASE category is correct');
    });

    test('PLUGIN_TYPES contains expected values', function (assert) {
      assert.strictEqual(PLUGIN_TYPES.SECRET, 'secret', 'SECRET type is correct');
      assert.strictEqual(PLUGIN_TYPES.AUTH, 'auth', 'AUTH type is correct');
      assert.strictEqual(PLUGIN_TYPES.DATABASE, 'database', 'DATABASE type is correct');
    });

    test('PLUGIN_CATEGORIES contains expected values', function (assert) {
      assert.strictEqual(PLUGIN_CATEGORIES.GENERIC, 'generic', 'GENERIC category is correct');
      assert.strictEqual(PLUGIN_CATEGORIES.CLOUD, 'cloud', 'CLOUD category is correct');
      assert.strictEqual(PLUGIN_CATEGORIES.INFRA, 'infra', 'INFRA category is correct');
      assert.strictEqual(PLUGIN_CATEGORIES.EXTERNAL, 'external', 'EXTERNAL category is correct');
    });
  });

  module('getAllVersionsForEngineType', function () {
    test('it returns empty array when no catalog data provided', function (assert) {
      const result = getAllVersionsForEngineType(undefined, 'kv', 'secret');
      assert.deepEqual(
        result,
        { versions: [], hasUnversionedPlugins: false },
        'returns empty result for undefined catalog data'
      );

      const result2 = getAllVersionsForEngineType([], 'kv', 'secret');
      assert.deepEqual(
        result2,
        { versions: [], hasUnversionedPlugins: false },
        'returns empty result for empty catalog data'
      );
    });

    test('it returns versions for direct engine type matches', function (assert) {
      const catalogData = [
        {
          name: 'kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
          version: '1.0.0',
        },
        {
          name: 'kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
          version: '2.0.0',
        },
      ];

      const result = getAllVersionsForEngineType(catalogData, 'kv', 'secret');

      assert.strictEqual(result.versions.length, 2, 'returns both versions');
      assert.strictEqual(result.versions[0].version, '1.0.0', 'includes first version');
      assert.strictEqual(result.versions[1].version, '2.0.0', 'includes second version');
      assert.strictEqual(result.versions[0].pluginName, 'kv', 'includes plugin name');
      assert.true(result.versions[0].isBuiltin, 'marks builtin correctly');
      assert.false(result.hasUnversionedPlugins, 'no unversioned plugins detected');
    });

    test('it returns versions for external plugins that map to engine types', function (assert) {
      const catalogData = [
        {
          name: 'kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
          version: '1.0.0',
        },
        {
          name: 'vault-plugin-secrets-kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
          version: '2.1.0',
        },
      ];

      const result = getAllVersionsForEngineType(catalogData, 'kv', 'secret');

      assert.strictEqual(result.versions.length, 2, 'returns both builtin and external versions');

      const builtinVersion = result.versions.find((v) => v.isBuiltin);
      const externalVersion = result.versions.find((v) => !v.isBuiltin);

      assert.ok(builtinVersion, 'includes builtin version');
      assert.ok(externalVersion, 'includes external version');
      assert.strictEqual(builtinVersion.pluginName, 'kv', 'builtin uses engine name');
      assert.strictEqual(
        externalVersion.pluginName,
        'vault-plugin-secrets-kv',
        'external uses full plugin name'
      );
      assert.false(result.hasUnversionedPlugins, 'no unversioned plugins detected');
    });

    test('it excludes external plugins that do not map to the engine type', function (assert) {
      const catalogData = [
        {
          name: 'kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
          version: '1.0.0',
        },
        {
          name: 'vault-plugin-secrets-aws',
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
          version: '1.5.0',
        },
      ];

      const result = getAllVersionsForEngineType(catalogData, 'kv', 'secret');

      assert.strictEqual(result.versions.length, 1, 'only includes matching plugins');
      assert.strictEqual(result.versions[0].pluginName, 'kv', 'includes only KV engine');
      assert.false(result.hasUnversionedPlugins, 'no unversioned plugins detected');
    });

    test('it filters by plugin type correctly', function (assert) {
      const catalogData = [
        {
          name: 'gcp',
          type: 'auth',
          builtin: true,
          version: 'v0.22.0+builtin',
        },
        {
          name: 'gcp',
          type: 'secret',
          builtin: true,
          version: 'v0.23.0+builtin',
        },
        {
          name: 'vault-plugin-secrets-gcp',
          type: 'secret',
          builtin: false,
          version: 'v0.23.0',
        },
      ];

      // Test filtering for secret plugins only
      const secretResult = getAllVersionsForEngineType(catalogData, 'gcp', 'secret');
      assert.strictEqual(secretResult.versions.length, 2, 'returns only secret type plugins');
      assert.true(
        secretResult.versions.every(
          (plugin) => plugin.pluginName === 'gcp' || plugin.pluginName === 'vault-plugin-secrets-gcp'
        ),
        'includes correct secret plugins'
      );
      assert.false(secretResult.hasUnversionedPlugins, 'no unversioned plugins detected');

      // Test filtering for auth plugins only
      const authResult = getAllVersionsForEngineType(catalogData, 'gcp', 'auth');
      assert.strictEqual(authResult.versions.length, 1, 'returns only auth type plugins');
      assert.strictEqual(authResult.versions[0].pluginName, 'gcp', 'includes auth plugin');
      assert.false(authResult.hasUnversionedPlugins, 'no unversioned plugins detected');
    });

    test('it handles invalid catalog data gracefully', function (assert) {
      const invalidCatalogData = [
        null, // null entry
        { name: 'kv' }, // missing required fields
        { name: 'aws', version: '1.0.0' }, // missing builtin field
        {
          name: 'pki',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
          version: '1.0.0',
        }, // valid entry
      ];

      const result = getAllVersionsForEngineType(invalidCatalogData, 'pki');

      assert.strictEqual(result.versions.length, 1, 'filters out invalid entries');
      assert.strictEqual(result.versions[0].pluginName, 'pki', 'includes only valid entry');
      assert.false(result.hasUnversionedPlugins, 'no unversioned plugins detected');
    });

    test('it excludes unversioned plugins but detects them', function (assert) {
      const catalogData = [
        {
          name: 'vault-plugin-secrets-keymgmt',
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
          version: '', // Empty string when plugin registered without version
          sha256: '9433b2b37d30abf8f7cbf8c3e616dfc263034789681081ea4ba7918673d80086',
        },
        {
          name: 'vault-plugin-secrets-keymgmt',
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
          version: '1.5.0',
        },
      ];

      const result = getAllVersionsForEngineType(catalogData, 'keymgmt', 'secret');

      assert.strictEqual(result.versions.length, 1, 'excludes unversioned plugin from versions');
      assert.true(result.hasUnversionedPlugins, 'detects presence of unversioned plugins');

      const versionedPlugin = result.versions[0];
      assert.strictEqual(versionedPlugin.version, '1.5.0', 'includes only versioned plugin');
      assert.false(versionedPlugin.isBuiltin, 'versioned plugin is not builtin');
      assert.strictEqual(
        versionedPlugin.pluginName,
        'vault-plugin-secrets-keymgmt',
        'correct plugin name for versioned'
      );
    });

    test('it detects unversioned plugins for builtin engines', function (assert) {
      const catalogData = [
        {
          name: 'kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
          version: '1.0.0',
        },
        {
          name: 'kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
          version: '', // Unversioned external kv plugin
        },
      ];

      const result = getAllVersionsForEngineType(catalogData, 'kv', 'secret');

      assert.strictEqual(result.versions.length, 1, 'only includes versioned plugins');
      assert.true(result.hasUnversionedPlugins, 'detects unversioned plugin');
      assert.strictEqual(result.versions[0].version, '1.0.0', 'includes builtin version');
      assert.true(result.versions[0].isBuiltin, 'included plugin is builtin');
    });

    test('it handles multiple unversioned plugins for same engine type', function (assert) {
      const catalogData = [
        {
          name: 'vault-plugin-secrets-custom',
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
          version: '', // First unversioned plugin
        },
        {
          name: 'custom',
          type: PLUGIN_TYPES.SECRET,
          builtin: false,
          version: '', // Second unversioned plugin (direct match)
        },
        {
          name: 'custom',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
          version: '2.0.0', // Versioned plugin
        },
      ];

      const result = getAllVersionsForEngineType(catalogData, 'custom', 'secret');

      assert.strictEqual(result.versions.length, 1, 'excludes all unversioned plugins');
      assert.true(result.hasUnversionedPlugins, 'detects multiple unversioned plugins');
      assert.strictEqual(result.versions[0].version, '2.0.0', 'includes only versioned plugin');
    });

    test('it handles invalid engine type parameters', function (assert) {
      const catalogData = [
        {
          name: 'kv',
          type: PLUGIN_TYPES.SECRET,
          builtin: true,
          version: '1.0.0',
        },
      ];

      const result1 = getAllVersionsForEngineType(catalogData, null);
      assert.deepEqual(
        result1,
        { versions: [], hasUnversionedPlugins: false },
        'returns empty result for null engine type'
      );

      const result2 = getAllVersionsForEngineType(catalogData, '');
      assert.deepEqual(
        result2,
        { versions: [], hasUnversionedPlugins: false },
        'returns empty result for empty engine type'
      );

      const result3 = getAllVersionsForEngineType(catalogData, undefined);
      assert.deepEqual(
        result3,
        { versions: [], hasUnversionedPlugins: false },
        'returns empty result for undefined engine type'
      );
    });
  });
});
