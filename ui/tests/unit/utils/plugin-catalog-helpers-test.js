/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  addVersionsToEngines,
  isValidPluginCatalogResponse,
  categorizeEnginesByStatus,
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

  test('addVersionsToEngines handles external plugins correctly', function (assert) {
    const staticEngines = [
      { type: 'custom-plugin', displayName: 'Custom Plugin', mountCategory: ['secret'] },
      { type: 'aws', displayName: 'AWS', mountCategory: ['secret'] },
    ];

    const pluginCatalogData = [
      { name: 'custom-plugin', type: 'secret', builtin: false, version: 'v2.1.0' },
      { name: 'aws', type: 'secret', builtin: true, version: 'v1.12.0+builtin.vault' },
    ];

    const result = addVersionsToEngines(staticEngines, pluginCatalogData);

    // Custom external plugin
    const customPlugin = result.find((e) => e.type === 'custom-plugin');
    assert.strictEqual(customPlugin.version, 'v2.1.0', 'External plugin should have exact version');
    assert.false(customPlugin.builtin, 'Custom plugin should be marked as external');
    assert.true(customPlugin.isAvailable, 'External plugin should be marked as available');

    // Builtin plugin
    const awsEngine = result.find((e) => e.type === 'aws');
    assert.strictEqual(awsEngine.version, 'v1.12.0', 'Builtin version should be cleaned');
    assert.true(awsEngine.builtin, 'AWS should be marked as builtin');
  });

  test('addVersionsToEngines handles deprecation status', function (assert) {
    const staticEngines = [
      { type: 'deprecated-plugin', displayName: 'Deprecated Plugin', mountCategory: ['secret'] },
      { type: 'pending-removal', displayName: 'Pending Removal', mountCategory: ['secret'] },
    ];

    const pluginCatalogData = [
      {
        name: 'deprecated-plugin',
        type: 'secret',
        builtin: true,
        version: 'v1.0.0',
        deprecation_status: 'deprecated',
      },
      {
        name: 'pending-removal',
        type: 'secret',
        builtin: true,
        version: 'v0.9.0',
        deprecation_status: 'pending-removal',
      },
    ];

    const result = addVersionsToEngines(staticEngines, pluginCatalogData);

    const deprecatedPlugin = result.find((e) => e.type === 'deprecated-plugin');
    assert.strictEqual(
      deprecatedPlugin.deprecation_status,
      'deprecated',
      'Should preserve deprecation status'
    );

    const pendingRemovalPlugin = result.find((e) => e.type === 'pending-removal');
    assert.strictEqual(
      pendingRemovalPlugin.deprecation_status,
      'pending-removal',
      'Should preserve pending removal status'
    );
  });

  test('addVersionsToEngines handles plugins with OCI image and runtime', function (assert) {
    const staticEngines = [{ type: 'oci-plugin', displayName: 'OCI Plugin', mountCategory: ['secret'] }];

    const pluginCatalogData = [
      {
        name: 'oci-plugin',
        type: 'secret',
        builtin: false,
        version: 'v1.5.0',
        oci_image: 'hashicorp/vault-plugin-secrets-custom:v1.5.0',
        runtime: 'container-runtime',
      },
    ];

    const result = addVersionsToEngines(staticEngines, pluginCatalogData);

    const ociPlugin = result.find((e) => e.type === 'oci-plugin');
    assert.strictEqual(ociPlugin.version, 'v1.5.0', 'Should preserve OCI plugin version');
    assert.false(ociPlugin.builtin, 'OCI plugin should be marked as external');
    assert.strictEqual(
      ociPlugin.pluginData.oci_image,
      'hashicorp/vault-plugin-secrets-custom:v1.5.0',
      'Should preserve OCI image information'
    );
    assert.strictEqual(
      ociPlugin.pluginData.runtime,
      'container-runtime',
      'Should preserve runtime information'
    );
  });

  test('addVersionsToEngines preserves static engine metadata', function (assert) {
    const staticEngines = [
      {
        type: 'aws',
        displayName: 'AWS Secrets Engine',
        mountCategory: ['secret'],
        pluginCategory: 'cloud',
        glyph: 'aws',
        searchTags: ['cloud', 'dynamic'],
      },
    ];

    const pluginCatalogData = [
      { name: 'aws', type: 'secret', builtin: true, version: 'v1.12.0+builtin.vault' },
    ];

    const result = addVersionsToEngines(staticEngines, pluginCatalogData);

    const awsEngine = result.find((e) => e.type === 'aws');
    assert.strictEqual(awsEngine.displayName, 'AWS Secrets Engine', 'Should preserve display name');
    assert.strictEqual(awsEngine.pluginCategory, 'cloud', 'Should preserve plugin category');
    assert.strictEqual(awsEngine.glyph, 'aws', 'Should preserve glyph');
    assert.deepEqual(awsEngine.searchTags, ['cloud', 'dynamic'], 'Should preserve search tags');
  });

  test('categorizeEnginesByStatus separates enabled and disabled engines', function (assert) {
    const engines = [
      {
        type: 'aws',
        displayName: 'AWS',
        isAvailable: true,
        pluginCategory: 'cloud',
        mountCategory: ['secret'],
      },
      {
        type: 'kv',
        displayName: 'KV Version 2',
        isAvailable: true,
        pluginCategory: 'generic',
        mountCategory: ['secret'],
      },
      {
        type: 'disabled-plugin',
        displayName: 'Disabled Plugin',
        isAvailable: false,
        pluginCategory: 'generic',
        mountCategory: ['secret'],
      },
      {
        type: 'unavailable-plugin',
        displayName: 'Unavailable Plugin',
        pluginCategory: 'infra',
        mountCategory: ['secret'],
      },
    ];

    const result = categorizeEnginesByStatus(engines);

    assert.strictEqual(result.enabled.length, 2, 'Should have 2 enabled engines');
    assert.strictEqual(result.disabled.length, 2, 'Should have 2 disabled engines');

    const enabledTypes = result.enabled.map((e) => e.type);
    assert.deepEqual(enabledTypes, ['aws', 'kv'], 'Should have correct enabled engines');

    const disabledTypes = result.disabled.map((e) => e.type);
    assert.deepEqual(
      disabledTypes,
      ['disabled-plugin', 'unavailable-plugin'],
      'Should have correct disabled engines'
    );
  });

  test('categorizeEnginesByStatus handles empty arrays', function (assert) {
    const result = categorizeEnginesByStatus([]);

    assert.strictEqual(result.enabled.length, 0, 'Should have no enabled engines');
    assert.strictEqual(result.disabled.length, 0, 'Should have no disabled engines');
    assert.ok(Array.isArray(result.enabled), 'Enabled should be an array');
    assert.ok(Array.isArray(result.disabled), 'Disabled should be an array');
  });

  test('categorizeEnginesByStatus treats undefined isAvailable as enabled', function (assert) {
    const engines = [
      {
        type: 'aws',
        displayName: 'AWS',
        pluginCategory: 'cloud',
        mountCategory: ['secret'],
        // isAvailable is undefined
      },
      {
        type: 'kv',
        displayName: 'KV Version 2',
        isAvailable: true,
        pluginCategory: 'generic',
        mountCategory: ['secret'],
      },
    ];

    const result = categorizeEnginesByStatus(engines);

    assert.strictEqual(result.enabled.length, 2, 'Should treat undefined isAvailable as enabled');
    assert.strictEqual(result.disabled.length, 0, 'Should have no disabled engines');
  });
});
