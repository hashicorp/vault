/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import SecretsEngineResource from 'vault/resources/secrets/engine';

import type { Mount } from 'vault/mount';

const createEngine = (overrides: Partial<Mount> = {}) =>
  new SecretsEngineResource({
    accessor: 'kv_12345678',
    config: {},
    description: '',
    external_entropy_access: false,
    local: false,
    path: 'example/',
    plugin_version: '',
    running_plugin_version: '',
    running_sha256: '',
    seal_wrap: false,
    type: 'kv',
    uuid: '12345678-1234-1234-1234-123456789012',
    ...overrides,
  });

module('Unit | Resource | secrets engine', function () {
  test('derives engine type and icon from mount data', function (assert) {
    const namespacedEngine = createEngine({ type: 'ns_identity' });
    const keyManagementEngine = createEngine({ type: 'keymgmt' });
    const unknownEngine = createEngine({ type: 'unknown-engine' });

    assert.strictEqual(namespacedEngine.engineType, 'identity', 'namespace prefix is removed');
    assert.strictEqual(keyManagementEngine.icon, 'key', 'key management uses the key icon');
    assert.strictEqual(unknownEngine.icon, 'lock', 'unknown engines use the fallback icon');
  });

  test('identifies KV v2 mounts and list visibility', function (assert) {
    const externalKvV2 = createEngine({
      type: 'vault-plugin-secrets-kv',
      options: { version: 2 },
    });

    assert.true(externalKvV2.isV2KV, 'external KV v2 plugins use the KV v2 experience');
    assert.false(createEngine({ type: 'system' }).shouldIncludeInList, 'internal mounts are excluded');
    assert.true(createEngine().shouldIncludeInList, 'visible mounts remain in the list');
  });

  test('derives supported status and backend routes', function (assert) {
    const kvV2 = createEngine({ options: { version: 2 } });
    const database = createEngine({ type: 'database' });
    const pki = createEngine({ type: 'pki' });
    const unknownEngine = createEngine({ type: 'unknown-engine' });

    assert.true(kvV2.isSupportedBackend, 'KV is a supported backend');
    assert.false(unknownEngine.isSupportedBackend, 'unknown engines are not supported backends');
    assert.strictEqual(kvV2.backendLink, 'vault.cluster.secrets.backend.kv.list');
    assert.strictEqual(database.backendLink, 'vault.cluster.secrets.backend.overview');
    assert.strictEqual(pki.backendLink, 'vault.cluster.secrets.backend.pki.overview');
    assert.strictEqual(unknownEngine.backendLink, 'vault.cluster.secrets.backend.list-root');
  });

  test('derives configuration route and local display', function (assert) {
    const azure = createEngine({ type: 'vault-plugin-secrets-azure' });
    const localEngine = createEngine({ local: true });

    assert.strictEqual(
      azure.backendConfigurationLink,
      'vault.cluster.secrets.backend.configuration.plugin-settings',
      'configurable external plugins use plugin settings'
    );
    assert.strictEqual(localEngine.localDisplay, 'local');
    assert.strictEqual(createEngine().localDisplay, 'replicated');
  });
});
