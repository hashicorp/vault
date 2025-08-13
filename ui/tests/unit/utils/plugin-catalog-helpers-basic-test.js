/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { addVersionsToEngines } from 'vault/utils/plugin-catalog-helpers';

module('Unit | Utility | plugin-catalog-helpers | basic functionality', function () {
  test('addVersionsToEngines works with new signature', function (assert) {
    const staticEngines = [
      { type: 'aws', displayName: 'AWS', mountCategory: ['secret'], glyph: 'aws-color' },
      { type: 'kv', displayName: 'KV', mountCategory: ['secret'], glyph: 'key-values' },
    ];

    const secretEnginesList = ['aws', 'kv', 'external-plugin'];
    const secretEnginesDetailed = [
      {
        name: 'aws',
        type: 'secret',
        builtin: true,
        version: 'v1.12.0+builtin.vault',
        deprecation_status: 'supported',
      },
      {
        name: 'kv',
        type: 'secret',
        builtin: true,
        version: 'v0.24.1+builtin',
        deprecation_status: 'supported',
      },
      {
        name: 'external-plugin',
        type: 'secret',
        builtin: false,
        version: 'v2.0.0',
        deprecation_status: 'supported',
      },
    ];

    const result = addVersionsToEngines(staticEngines, secretEnginesList, secretEnginesDetailed);
    const realEngines = result.filter(
      (e) => !e.type.startsWith('demo-') && !e.type.startsWith('example-') && !e.type.startsWith('test-')
    );

    // Should have: 2 static + 1 external plugin
    assert.strictEqual(realEngines.length, 3, 'Should include all engines');

    // Check static engines are enhanced
    const awsEngine = realEngines.find((e) => e.type === 'aws');
    assert.ok(awsEngine, 'AWS engine should be included');
    assert.true(awsEngine.isAvailable, 'AWS should be available');
    assert.strictEqual(awsEngine.version, 'v1.12.0', 'AWS should have cleaned version');
    assert.false(awsEngine.isExternalPlugin, 'AWS should not be marked as external');

    const kvEngine = realEngines.find((e) => e.type === 'kv');
    assert.ok(kvEngine, 'KV engine should be included');
    assert.true(kvEngine.isAvailable, 'KV should be available');
    assert.strictEqual(kvEngine.version, 'v0.24.1', 'KV should have cleaned version');

    // Check external plugin
    const externalPlugin = realEngines.find((e) => e.type === 'external-plugin');
    assert.ok(externalPlugin, 'External plugin should be included');
    assert.strictEqual(externalPlugin.pluginCategory, 'external', 'Should be marked as external');
    assert.strictEqual(externalPlugin.version, 'v2.0.0', 'Should have version from detailed info');
    assert.true(externalPlugin.isExternalPlugin, 'Should be marked as external plugin');
  });

  test('addVersionsToEngines handles empty lists', function (assert) {
    const staticEngines = [{ type: 'aws', displayName: 'AWS', mountCategory: ['secret'] }];

    const result = addVersionsToEngines(staticEngines, [], []);
    
    // Should return the same number of engines with the same types
    assert.strictEqual(result.length, staticEngines.length, 'Should return same number of engines');
    
    const awsEngine = result.find((e) => e.type === 'aws');
    assert.ok(awsEngine, 'AWS engine should be present');
    assert.strictEqual(awsEngine.displayName, 'AWS', 'Should preserve display name');
    assert.false(awsEngine.isAvailable, 'Should be marked as unavailable when no catalog data');
  });
});
