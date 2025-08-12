/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  addVersionsToEngines,
  isValidPluginCatalogResponse,
} from 'vault/utils/plugin-catalog-helpers';

module('Unit | Utility | plugin-catalog-helpers', function () {
  test('addVersionsToEngines merges plugin data with static engines', function (assert) {
    const staticEngines = [
      { type: 'aws', displayName: 'AWS', mountCategory: ['secret'] },
      { type: 'kv', displayName: 'KV Version 2', mountCategory: ['secret'] },
      { type: 'unknown', displayName: 'Unknown Engine', mountCategory: ['secret'] },
    ];

    const pluginCatalogData = [
      { name: 'aws', type: 'secret', builtin: true, version: 'v1.12.0+builtin.vault' },
      { name: 'kv', type: 'secret', builtin: true, version: 'v0.24.1+builtin' },
    ];

    const result = addVersionsToEngines(staticEngines, pluginCatalogData);

    assert.strictEqual(result.length, 3, 'Should return same number of engines');

    // AWS engine should have version info
    const awsEngine = result.find((e) => e.type === 'aws');
    assert.strictEqual(awsEngine.version, 'v1.12.0', 'AWS should have cleaned version');
    assert.true(awsEngine.builtin, 'AWS should be marked as builtin');
    assert.true(awsEngine.isAvailable, 'AWS should be marked as available');

    // Test KV engine (has plugin data with +builtin suffix)
    const kvEngine = result.find((e) => e.type === 'kv');
    assert.strictEqual(kvEngine.version, 'v0.24.1', 'KV should have cleaned version');
    assert.true(kvEngine.builtin, 'KV should be marked as builtin');

    // Unknown engine should not have version info but should be marked unavailable
    const unknownEngine = result.find((e) => e.type === 'unknown');
    assert.strictEqual(unknownEngine.version, undefined, 'Unknown should not have version');
    assert.false(unknownEngine.isAvailable, 'Unknown should be marked as unavailable');
  });

  test('addVersionsToEngines handles empty or invalid plugin catalog data', function (assert) {
    const staticEngines = [{ type: 'aws', displayName: 'AWS', mountCategory: ['secret'] }];

    // Test with null
    let result = addVersionsToEngines(staticEngines, null);
    assert.deepEqual(result, staticEngines, 'Should return original engines with null data');

    // Test with undefined
    result = addVersionsToEngines(staticEngines, undefined);
    assert.deepEqual(result, staticEngines, 'Should return original engines with undefined data');

    // Test with empty array
    result = addVersionsToEngines(staticEngines, []);
    assert.strictEqual(result.length, 1, 'Should return engines with empty catalog');
    assert.false(result[0].isAvailable, 'Engine should be marked unavailable');
  });

  test('isValidPluginCatalogResponse validates response structure', function (assert) {
    // Valid response
    const validResponse = {
      data: {
        detailed: [{ name: 'aws', type: 'secret', builtin: true, version: 'v1.12.0' }],
        secret: ['aws'],
      },
    };
    assert.true(isValidPluginCatalogResponse(validResponse), 'Should validate correct response');

    // Invalid responses
    assert.false(isValidPluginCatalogResponse(null), 'Should reject null');
    assert.false(isValidPluginCatalogResponse(undefined), 'Should reject undefined');
    assert.false(isValidPluginCatalogResponse({}), 'Should reject empty object');
    assert.false(isValidPluginCatalogResponse({ data: {} }), 'Should reject missing detailed array');
    assert.false(
      isValidPluginCatalogResponse({ data: { detailed: 'not-array' } }),
      'Should reject non-array detailed'
    );
  });
});
