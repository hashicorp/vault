/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { setupTest } from 'ember-qunit';
import { module, test } from 'qunit';
import SecretsEngineResource from 'vault/resources/secrets/engine';

const makeResource = ({ type, version }) => {
  const options = version ? { version } : undefined;
  return new SecretsEngineResource({
    accessor: 'test_accessor',
    config: {},
    description: '',
    external_entropy_access: false,
    local: false,
    options,
    path: `${type}-test/`,
    plugin_version: '',
    running_plugin_version: '',
    running_sha256: '',
    seal_wrap: false,
    type,
    uuid: 'test-uuid',
  });
};

module('Unit | Component | manage-dropdown', function (hooks) {
  setupTest(hooks);

  test('backendConfigurationLink: addon engines with configRoute always use their config route', function (assert) {
    const kubernetes = makeResource({ type: 'kubernetes' });
    assert.strictEqual(
      kubernetes.backendConfigurationLink,
      'vault.cluster.secrets.backend.kubernetes.configuration',
      'kubernetes always routes to its configuration page'
    );

    const ldap = makeResource({ type: 'ldap' });
    assert.strictEqual(
      ldap.backendConfigurationLink,
      'vault.cluster.secrets.backend.ldap.configuration',
      'ldap always routes to its configuration page'
    );

    const pki = makeResource({ type: 'pki' });
    assert.strictEqual(
      pki.backendConfigurationLink,
      'vault.cluster.secrets.backend.pki.configuration',
      'pki always routes to its configuration page'
    );
  });

  test('backendConfigurationLink: configurable engines without configRoute route to plugin-settings', function (assert) {
    const ssh = makeResource({ type: 'ssh' });
    assert.strictEqual(
      ssh.backendConfigurationLink,
      'vault.cluster.secrets.backend.configuration.plugin-settings',
      'configurable engine routes to the plugin-settings view'
    );
  });

  test('backendConfigurationLink: non-configurable engines always route to general-settings', function (assert) {
    const alicloud = makeResource({ type: 'alicloud' });
    assert.strictEqual(
      alicloud.backendConfigurationLink,
      'vault.cluster.secrets.backend.configuration.general-settings'
    );
  });

  test('backendConfigurationLink: KV v1 routes to general-settings', function (assert) {
    const kvV1 = makeResource({ type: 'kv', version: 1 });
    assert.strictEqual(
      kvV1.backendConfigurationLink,
      'vault.cluster.secrets.backend.configuration.general-settings'
    );
  });

  test('backendConfigurationLink: KV v2 routes to general-settings (configRoute is display-only)', function (assert) {
    const kvV2 = makeResource({ type: 'kv', version: 2 });
    assert.strictEqual(
      kvV2.backendConfigurationLink,
      'vault.cluster.secrets.backend.configuration.general-settings',
      "kv's configRoute is skipped because it's for display only"
    );
  });
});
