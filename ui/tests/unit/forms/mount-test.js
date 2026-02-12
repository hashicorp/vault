/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import MountForm from 'vault/forms/mount';

module('Unit | Form | mount', function (hooks) {
  setupTest(hooks);

  module('Plugin Version Selection', function () {
    test('toJSON omits plugin_version when default is selected (empty value)', function (assert) {
      assert.expect(2);

      const form = new MountForm({});
      form.type = 'kv';
      form.data = {
        path: 'test-kv',
        description: 'Test KV engine',
        config: {
          plugin_version: '', // Empty string represents default selection
          max_lease_ttl: '8760h',
        },
      };

      const result = form.toJSON();

      assert.strictEqual(result.data.type, 'kv', 'type is set correctly');
      assert.notOk(
        Object.prototype.hasOwnProperty.call(result.data.config, 'plugin_version'),
        'plugin_version is omitted from config when empty'
      );
    });

    test('toJSON omits plugin_version when undefined', function (assert) {
      assert.expect(2);

      const form = new MountForm({});
      form.type = 'kv';
      form.data = {
        path: 'test-kv',
        description: 'Test KV engine',
        config: {
          max_lease_ttl: '8760h',
          // plugin_version not set
        },
      };

      const result = form.toJSON();

      assert.strictEqual(result.data.type, 'kv', 'type is set correctly');
      assert.notOk(
        Object.prototype.hasOwnProperty.call(result.data.config, 'plugin_version'),
        'plugin_version is omitted from config when undefined'
      );
    });

    test('toJSON includes plugin_version for builtin plugin selection', function (assert) {
      assert.expect(3);

      const form = new MountForm({});
      form.type = 'kv'; // Builtin type should remain as 'kv'
      form.data = {
        path: 'test-kv',
        description: 'Test KV engine',
        config: {
          plugin_version: 'v1.16.1+builtin',
          max_lease_ttl: '8760h',
        },
      };

      const result = form.toJSON();

      assert.strictEqual(result.data.type, 'kv', 'type remains builtin type for builtin plugins');
      assert.strictEqual(
        result.data.config.plugin_version,
        'v1.16.1+builtin',
        'plugin_version is included for builtin plugin'
      );
      assert.strictEqual(result.data.path, 'test-kv', 'other data is preserved');
    });

    test('toJSON includes plugin_version for external plugin selection', function (assert) {
      assert.expect(3);

      const form = new MountForm({});
      form.type = 'vault-plugin-secrets-kv'; // External plugin type
      form.data = {
        path: 'test-external-kv',
        description: 'Test external KV engine',
        config: {
          plugin_version: 'v0.25.0',
          max_lease_ttl: '8760h',
        },
      };

      const result = form.toJSON();

      assert.strictEqual(
        result.data.type,
        'vault-plugin-secrets-kv',
        'type is set to external plugin name for external plugins'
      );
      assert.strictEqual(
        result.data.config.plugin_version,
        'v0.25.0',
        'plugin_version is included for external plugin'
      );
      assert.strictEqual(result.data.path, 'test-external-kv', 'other data is preserved');
    });
  });

  module('setPluginVersionData', function () {
    test('sets config.plugin_version and preserves builtin type for builtin plugins', function (assert) {
      assert.expect(3);

      const form = new MountForm({});
      form.type = 'kv';
      form.data = { config: {} };

      const builtinVersionInfo = {
        version: 'v1.16.1+builtin',
        pluginName: 'vault-plugin-secrets-kv',
        isBuiltin: true,
        sha256: 'abc123',
      };

      form.setPluginVersionData(builtinVersionInfo);

      assert.strictEqual(form.data.config.plugin_version, 'v1.16.1+builtin', 'plugin_version is set');
      assert.strictEqual(form.type, 'kv', 'type remains builtin for builtin plugins');
      assert.ok(form.data.config, 'config object is preserved');
    });

    test('sets config.plugin_version and updates type for external plugins', function (assert) {
      assert.expect(3);

      const form = new MountForm({});
      form.type = 'kv'; // Initially set to builtin
      form.data = { config: {} };

      const externalVersionInfo = {
        version: 'v0.25.0',
        pluginName: 'vault-plugin-secrets-kv',
        isBuiltin: false,
        sha256: 'def456',
      };

      form.setPluginVersionData(externalVersionInfo);

      assert.strictEqual(form.data.config.plugin_version, 'v0.25.0', 'plugin_version is set');
      assert.strictEqual(
        form.type,
        'vault-plugin-secrets-kv',
        'type is updated to plugin name for external plugins'
      );
      assert.ok(form.data.config, 'config object is preserved');
    });
  });

  module('findVersionByLabel', function () {
    test('returns undefined for empty string (default selection)', function (assert) {
      assert.expect(1);

      const form = new MountForm({});
      form.data = { config: {} };

      const availableVersions = [
        {
          version: 'v1.16.1+builtin',
          pluginName: 'vault-plugin-secrets-kv',
          isBuiltin: true,
          sha256: 'abc123',
        },
        {
          version: 'v0.25.0',
          pluginName: 'vault-plugin-secrets-kv',
          isBuiltin: false,
          sha256: 'def456',
        },
      ];

      const result = form.findVersionByLabel('', availableVersions);

      assert.strictEqual(result, undefined, 'returns undefined for empty string (default)');
    });

    test('returns undefined for null/undefined selectedValue', function (assert) {
      assert.expect(2);

      const form = new MountForm({});
      form.data = { config: {} };

      const availableVersions = [
        {
          version: 'v1.16.1+builtin',
          pluginName: 'vault-plugin-secrets-kv',
          isBuiltin: true,
          sha256: 'abc123',
        },
      ];

      assert.strictEqual(
        form.findVersionByLabel(null, availableVersions),
        undefined,
        'returns undefined for null'
      );
      assert.strictEqual(
        form.findVersionByLabel(undefined, availableVersions),
        undefined,
        'returns undefined for undefined'
      );
    });

    test('finds matching version info by version string', function (assert) {
      assert.expect(2);

      const form = new MountForm({});
      form.data = { config: {} };

      const builtinVersion = {
        version: 'v1.16.1+builtin',
        pluginName: 'vault-plugin-secrets-kv',
        isBuiltin: true,
        sha256: 'abc123',
      };
      const externalVersion = {
        version: 'v0.25.0',
        pluginName: 'vault-plugin-secrets-kv',
        isBuiltin: false,
        sha256: 'def456',
      };
      const availableVersions = [builtinVersion, externalVersion];

      const builtinResult = form.findVersionByLabel('v1.16.1+builtin', availableVersions);
      const externalResult = form.findVersionByLabel('v0.25.0', availableVersions);

      assert.deepEqual(builtinResult, builtinVersion, 'finds builtin version correctly');
      assert.deepEqual(externalResult, externalVersion, 'finds external version correctly');
    });

    test('returns undefined for non-matching version', function (assert) {
      assert.expect(1);

      const form = new MountForm({});
      form.data = { config: {} };

      const availableVersions = [
        {
          version: 'v1.16.1+builtin',
          pluginName: 'vault-plugin-secrets-kv',
          isBuiltin: true,
          sha256: 'abc123',
        },
      ];

      const result = form.findVersionByLabel('v999.999.999', availableVersions);

      assert.strictEqual(result, undefined, 'returns undefined for non-matching version');
    });
  });

  module('setupPluginVersionField', function () {
    test('does nothing when no versions available', function (assert) {
      assert.expect(1);

      const form = new MountForm({});
      form.data = { config: {} };

      form.setupPluginVersionField(null);

      // Since the field is handled in the template now, just verify the method doesn't throw
      assert.ok(true, 'setupPluginVersionField handles null versions gracefully');
    });

    test('does nothing when only one version available', function (assert) {
      assert.expect(1);

      const form = new MountForm({});
      form.data = { config: {} };

      const singleVersion = [
        {
          version: 'v1.16.1+builtin',
          pluginName: 'vault-plugin-secrets-kv',
          isBuiltin: true,
          sha256: 'abc123',
        },
      ];

      form.setupPluginVersionField(singleVersion);

      // Since the field is handled in the template now, just verify the method doesn't throw
      assert.ok(true, 'setupPluginVersionField handles single version gracefully');
    });

    test('initializes plugin_version config when multiple versions available', function (assert) {
      assert.expect(1);

      const form = new MountForm({});
      form.data = { config: {} };

      const multipleVersions = [
        {
          version: 'v1.16.1+builtin',
          pluginName: 'vault-plugin-secrets-kv',
          isBuiltin: true,
          sha256: 'abc123',
        },
        {
          version: 'v0.25.0',
          pluginName: 'vault-plugin-secrets-kv',
          isBuiltin: false,
          sha256: 'def456',
        },
      ];

      form.setupPluginVersionField(multipleVersions);

      assert.strictEqual(form.data.config.plugin_version, '', 'plugin_version initialized as empty string');
    });
  });
});
