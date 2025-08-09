/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import {
  addVersionsToEngines,
  categorizeEnginesByStatus,
  isValidPluginCatalogResponse,
} from 'vault/utils/plugin-catalog-helpers';

module('Unit | Utility | plugin-catalog-helpers', function () {
  test('addVersionsToEngines handles full plugin catalog response structure', function (assert) {
    const staticEngines = [
      { type: 'aws', displayName: 'AWS', mountCategory: ['secret'], glyph: 'aws-color' },
      { type: 'kv', displayName: 'KV', mountCategory: ['secret'], glyph: 'key-values' },
    ];

    // Full plugin catalog response structure
    const pluginCatalogResponse = {
      data: {
        secret: [
          'aws',
          'kv',
          'consul',
          'external-custom-plugin', // Plugin in secret list but not in detailed
          'my-custom-aws-variant', // External plugin with AWS-like name
        ],
        auth: ['approle', 'userpass'],
        detailed: [
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
            name: 'consul',
            type: 'secret',
            builtin: true,
            version: 'v1.21.0+builtin.vault',
            deprecation_status: 'supported',
          },
          {
            name: 'my-custom-aws-variant',
            type: 'secret',
            builtin: false,
            version: 'v2.0.0',
            deprecation_status: 'supported',
          },
          // Note: external-custom-plugin is in secret list but missing from detailed
        ],
      },
    };

    const result = addVersionsToEngines(
      staticEngines,
      pluginCatalogResponse.data.detailed.filter((plugin) => plugin.type === 'secret'),
      []
    );
    const realEngines = result.filter(
      (e) => !e.type.startsWith('demo-') && !e.type.startsWith('example-') && !e.type.startsWith('test-')
    );

    // Should have: 2 static + 3 external plugins (consul, external-custom-plugin, my-custom-aws-variant)
    assert.strictEqual(realEngines.length, 5, 'Should include all engines from secret list');

    // Check static engines are enhanced
    const awsEngine = realEngines.find((e) => e.type === 'aws');
    assert.ok(awsEngine, 'AWS engine should be included');
    assert.true(awsEngine.isAvailable, 'AWS should be available');
    assert.strictEqual(awsEngine.version, 'v1.12.0', 'AWS should have cleaned version');

    // Check external plugin with detailed info
    const customAwsVariant = realEngines.find((e) => e.type === 'my-custom-aws-variant');
    assert.ok(customAwsVariant, 'Custom AWS variant should be included');
    assert.strictEqual(customAwsVariant.pluginCategory, 'external', 'Should be marked as external');
    assert.strictEqual(customAwsVariant.glyph, 'aws-color', 'Should inherit AWS glyph');
    assert.strictEqual(customAwsVariant.version, 'v2.0.0', 'Should have version from detailed info');

    // Check external plugin in secret list but missing from detailed
    const externalCustomPlugin = realEngines.find((e) => e.type === 'external-custom-plugin');
    assert.ok(externalCustomPlugin, 'External custom plugin should be included');
    assert.strictEqual(externalCustomPlugin.pluginCategory, 'external', 'Should be marked as external');
    assert.strictEqual(externalCustomPlugin.glyph, 'file-text', 'Should use default glyph');
    assert.strictEqual(externalCustomPlugin.version, 'unknown', 'Should have unknown version');
    assert.false(externalCustomPlugin.builtin, 'Should not be marked as builtin');

    // Check external plugin that matches known type (consul)
    const consulEngine = realEngines.find((e) => e.type === 'consul');
    assert.ok(consulEngine, 'Consul engine should be included');
    assert.strictEqual(
      consulEngine.pluginCategory,
      'external',
      'Should be marked as external since not in static metadata'
    );
    assert.strictEqual(consulEngine.version, 'v1.21.0', 'Should have cleaned version');
  });

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

    const result = addVersionsToEngines(
      staticEngines,
      ['aws', 'kv'], // secret engines list
      pluginCatalogData // detailed data (already filtered for secret)
    );

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
    let result = addVersionsToEngines(staticEngines, null, []);
    assert.deepEqual(result, staticEngines, 'Should return original engines with null data');

    // Test with undefined
    result = addVersionsToEngines(staticEngines, undefined, []);
    assert.deepEqual(result, staticEngines, 'Should return original engines with undefined data');

    // Test with empty array
    result = addVersionsToEngines(staticEngines, [], []);
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

    const result = addVersionsToEngines(staticEngines, pluginCatalogData, []);

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

    const result = addVersionsToEngines(staticEngines, pluginCatalogData, []);

    const deprecatedPlugin = result.find((e) => e.type === 'deprecated-plugin');
    assert.strictEqual(
      deprecatedPlugin.deprecationStatus,
      'deprecated',
      'Should preserve deprecation status'
    );

    const pendingRemovalPlugin = result.find((e) => e.type === 'pending-removal');
    assert.strictEqual(
      pendingRemovalPlugin.deprecationStatus,
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

  test('addVersionsToEngines dynamically discovers plugins not in static metadata', function (assert) {
    const staticEngines = [
      { type: 'aws', displayName: 'AWS', mountCategory: ['secret'] },
      { type: 'kv', displayName: 'KV Version 2', mountCategory: ['secret'] },
    ];

    const pluginCatalogData = [
      { name: 'aws', type: 'secret', builtin: true, version: 'v1.12.0+builtin.vault' },
      { name: 'kv', type: 'secret', builtin: true, version: 'v0.24.1+builtin' },
      // Dynamic plugin not in static metadata
      { name: 'custom-plugin', type: 'secret', builtin: false, version: 'v2.0.0' },
      { name: 'another-external-plugin', type: 'secret', builtin: false, version: 'v1.5.2' },
      // Auth plugin should be ignored for now
      { name: 'custom-auth', type: 'auth', builtin: false, version: 'v1.0.0' },
    ];

    const result = addVersionsToEngines(staticEngines, pluginCatalogData);

    // Should have original 2 static engines + 2 dynamic secret engines (auth plugin ignored)
    const realEngines = result.filter(
      (e) => !e.type.startsWith('demo-') && !e.type.startsWith('example-') && !e.type.startsWith('test-')
    );
    assert.strictEqual(realEngines.length, 4, 'Should include static engines plus dynamic secret engines');

    // Check static engines are still enhanced
    const awsEngine = realEngines.find((e) => e.type === 'aws');
    assert.true(awsEngine.isAvailable, 'AWS should be available');

    // Check dynamic plugins
    const customPlugin = realEngines.find((e) => e.type === 'custom-plugin');
    assert.ok(customPlugin, 'Custom plugin should be included');
    assert.strictEqual(customPlugin.displayName, 'Custom Plugin', 'Should convert kebab-case to Title Case');
    assert.strictEqual(customPlugin.pluginCategory, 'external', 'Should be marked as external category');
    assert.strictEqual(customPlugin.glyph, 'file-text', 'Should default to file-text icon');
    assert.true(customPlugin.isAvailable, 'Dynamic plugin should be available');
    assert.false(customPlugin.builtin, 'Custom plugin should not be builtin');

    const anotherPlugin = realEngines.find((e) => e.type === 'another-external-plugin');
    assert.ok(anotherPlugin, 'Another external plugin should be included');
    assert.strictEqual(
      anotherPlugin.displayName,
      'Another External Plugin',
      'Should handle multiple hyphens'
    );

    // Auth plugin should not be included
    const authPlugin = realEngines.find((e) => e.type === 'custom-auth');
    assert.notOk(authPlugin, 'Auth plugin should not be included in secret engine discovery');
  });

  test('addVersionsToEngines uses glyph from matching static engine type', function (assert) {
    const staticEngines = [
      { type: 'aws', displayName: 'AWS', mountCategory: ['secret'], glyph: 'aws-color' },
      { type: 'kv', displayName: 'KV', mountCategory: ['secret'], glyph: 'key-values' },
    ];

    const pluginCatalogData = [
      // Static engines that exist in metadata
      { name: 'aws', type: 'secret', builtin: true, version: 'v1.12.0+builtin.vault' },

      // External plugins with names that contain known engine types
      { name: 'my-custom-aws-plugin', type: 'secret', builtin: false, version: 'v2.0.0' },
      { name: 'external-kv-store', type: 'secret', builtin: false, version: 'v1.5.0' },

      // External plugin with unknown type
      { name: 'completely-unknown-plugin', type: 'secret', builtin: false, version: 'v1.0.0' },
    ];

    const result = addVersionsToEngines(staticEngines, pluginCatalogData);
    const realEngines = result.filter(
      (e) => !e.type.startsWith('demo-') && !e.type.startsWith('example-') && !e.type.startsWith('test-')
    );

    // Should have 2 static + 3 external plugins
    assert.strictEqual(realEngines.length, 5, 'Should include static engines plus external plugins');

    // Check that external AWS plugin inherited the AWS glyph
    const externalAwsPlugin = realEngines.find((e) => e.type === 'my-custom-aws-plugin');
    assert.ok(externalAwsPlugin, 'External AWS plugin should be included');
    assert.strictEqual(externalAwsPlugin.glyph, 'aws-color', 'Should inherit glyph from AWS static engine');
    assert.strictEqual(externalAwsPlugin.pluginCategory, 'external', 'Should be marked as external');

    // Check that external KV plugin inherited the KV glyph
    const externalKvPlugin = realEngines.find((e) => e.type === 'external-kv-store');
    assert.ok(externalKvPlugin, 'External KV plugin should be included');
    assert.strictEqual(externalKvPlugin.glyph, 'key-values', 'Should inherit glyph from KV static engine');
    assert.strictEqual(externalKvPlugin.pluginCategory, 'external', 'Should be marked as external');

    // Check that unknown plugin uses default glyph
    const unknownPlugin = realEngines.find((e) => e.type === 'completely-unknown-plugin');
    assert.ok(unknownPlugin, 'Unknown plugin should be included');
    assert.strictEqual(unknownPlugin.glyph, 'file-text', 'Should use default glyph for unknown type');

    // Static engines should still work normally
    const staticAwsEngine = realEngines.find((e) => e.type === 'aws');
    assert.ok(staticAwsEngine, 'Static AWS engine should be included');
  });

  test('addVersionsToEngines handles complex plugin names correctly', function (assert) {
    const staticEngines = [];

    const pluginCatalogData = [
      { name: 'my-custom-vault-plugin', type: 'secret', builtin: false, version: 'v1.0.0' },
      { name: 'simple', type: 'secret', builtin: false, version: 'v2.0.0' },
      { name: 'multi-word-plugin-name', type: 'secret', builtin: false, version: 'v3.0.0' },
    ];

    const result = addVersionsToEngines(staticEngines, pluginCatalogData);
    const realEngines = result.filter(
      (e) => !e.type.startsWith('demo-') && !e.type.startsWith('example-') && !e.type.startsWith('test-')
    );

    assert.strictEqual(realEngines.length, 3, 'Should create engines for all secret plugins');

    const plugin1 = realEngines.find((e) => e.type === 'my-custom-vault-plugin');
    assert.strictEqual(plugin1.displayName, 'My Custom Vault Plugin', 'Should handle complex names');

    const plugin2 = realEngines.find((e) => e.type === 'simple');
    assert.strictEqual(plugin2.displayName, 'Simple', 'Should handle single word names');

    const plugin3 = realEngines.find((e) => e.type === 'multi-word-plugin-name');
    assert.strictEqual(plugin3.displayName, 'Multi Word Plugin Name', 'Should handle many words');
  });
});
