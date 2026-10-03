/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

/**
 * Test the secret-engine model to ensure external plugin mapping
 * works correctly for key getters that affect routing and UI behavior.
 */
module('Unit | Model | secret-engine external plugin support', function (hooks) {
  setupTest(hooks);

  module('isV2KV getter', function () {
    test('returns true for external KV v2 plugins', function (assert) {
      const store = this.owner.lookup('service:store');

      const externalKvV2 = store.createRecord('secret-engine', {
        type: 'vault-plugin-secrets-kv',
        version: 2,
      });

      assert.true(externalKvV2.isV2KV, 'External KV v2 plugin is recognized as V2 KV');
    });

    test('returns false for external KV v1 plugins', function (assert) {
      const store = this.owner.lookup('service:store');

      const externalKvV1 = store.createRecord('secret-engine', {
        type: 'vault-plugin-secrets-kv',
        version: 1,
      });

      assert.false(externalKvV1.isV2KV, 'External KV v1 plugin is not V2 KV');
    });

    test('returns true for builtin KV v2 engines', function (assert) {
      const store = this.owner.lookup('service:store');

      const builtinKvV2 = store.createRecord('secret-engine', {
        type: 'kv',
        version: 2,
      });

      assert.true(builtinKvV2.isV2KV, 'Builtin KV v2 engine is recognized as V2 KV');
    });

    test('returns true for generic v2 engines', function (assert) {
      const store = this.owner.lookup('service:store');

      const genericV2 = store.createRecord('secret-engine', {
        type: 'generic',
        version: 2,
      });

      assert.true(genericV2.isV2KV, 'Generic v2 engine is recognized as V2 KV');
    });

    test('returns false for non-KV external plugins', function (assert) {
      const store = this.owner.lookup('service:store');

      const externalKeymgmt = store.createRecord('secret-engine', {
        type: 'vault-plugin-secrets-keymgmt',
        version: 1,
      });

      assert.false(externalKeymgmt.isV2KV, 'External keymgmt plugin is not V2 KV');
    });
  });

  module('backendLink getter', function () {
    test('returns KV engine route for external KV v2 plugins', function (assert) {
      const store = this.owner.lookup('service:store');

      const externalKvV2 = store.createRecord('secret-engine', {
        type: 'vault-plugin-secrets-kv',
        version: 2,
      });

      const backendLink = externalKvV2.backendLink;
      assert.true(backendLink.includes('kv.list'), `External KV v2 uses KV engine route: ${backendLink}`);
    });

    test('returns correct route for external database plugins', function (assert) {
      const store = this.owner.lookup('service:store');

      // Mock external database plugin (though not in our current mapping)
      const externalDb = store.createRecord('secret-engine', {
        type: 'vault-plugin-database-postgresql',
      });

      const backendLink = externalDb.backendLink;
      // Should fall back to list-root for unmapped plugins
      assert.strictEqual(
        backendLink,
        'vault.cluster.secrets.backend.list-root',
        'Unmapped external plugin uses generic route'
      );
    });

    test('handles external keymgmt plugins correctly', function (assert) {
      const store = this.owner.lookup('service:store');

      const externalKeymgmt = store.createRecord('secret-engine', {
        type: 'vault-plugin-secrets-keymgmt',
      });

      const backendLink = externalKeymgmt.backendLink;
      // External keymgmt should route to generic since keymgmt doesn't have engineRoute
      assert.strictEqual(
        backendLink,
        'vault.cluster.secrets.backend.list-root',
        'External keymgmt uses list-root route'
      );
    });
  });

  module('backendConfigurationLink getter', function () {
    test('returns effective type configuration route for external plugins', function (assert) {
      const store = this.owner.lookup('service:store');

      const externalAzure = store.createRecord('secret-engine', {
        type: 'vault-plugin-secrets-azure',
      });

      const configLink = externalAzure.backendConfigurationLink;
      // Note: The old secret-engine model uses isAddonEngine logic, so Azure (not an addon)
      // falls back to general-settings rather than plugin-settings
      assert.strictEqual(
        configLink,
        'vault.cluster.secrets.backend.configuration.general-settings',
        `External Azure uses general settings route in old model: ${configLink}`
      );
    });
    test('fallback to generic configuration for unmapped plugins', function (assert) {
      const store = this.owner.lookup('service:store');

      const unknownExternal = store.createRecord('secret-engine', {
        type: 'vault-plugin-secrets-unknown',
      });

      const configLink = unknownExternal.backendConfigurationLink;
      assert.strictEqual(
        configLink,
        'vault.cluster.secrets.backend.configuration.general-settings',
        'Unknown external plugin uses generic configuration route'
      );
    });
  });
});
