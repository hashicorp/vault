/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  enhanceEnginesWithCatalogData,
  categorizeEnginesByStatus,
  MOUNT_CATEGORIES,
  PLUGIN_TYPES,
  PLUGIN_CATEGORIES,
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
});
