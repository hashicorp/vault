/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, render } from '@ember/test-helpers';
import { setupRenderingTest } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';
import { module, test } from 'qunit';
import sinon from 'sinon';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const DEFAULT_MOUNT_DATA = {
  accessor: 'test_accessor',
  config: {},
  description: '',
  external_entropy_access: false,
  local: false,
  plugin_version: '',
  running_plugin_version: '',
  running_sha256: '',
  seal_wrap: false,
  uuid: 'test-uuid',
};

const TEST_CASES = [
  {
    label: 'alicloud',
    type: 'alicloud',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  {
    label: 'azure',
    type: 'azure',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.plugin-settings',
  },
  {
    label: 'gcp',
    type: 'gcp',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.plugin-settings',
  },
  {
    label: 'gcpkms',
    type: 'gcpkms',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  {
    label: 'keymgmt',
    type: 'keymgmt',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  {
    label: 'kubernetes',
    type: 'kubernetes',
    expectedRoute: 'vault.cluster.secrets.backend.kubernetes.configuration',
  },
  {
    label: 'kvv1',
    type: 'kv',
    version: 1,
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  {
    label: 'kvv2',
    type: 'kv',
    version: 2,
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  {
    label: 'transform',
    type: 'transform',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  {
    label: 'transit',
    type: 'transit',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  { label: 'kmip', type: 'kmip', expectedRoute: 'vault.cluster.secrets.backend.kmip.configuration' },
  { label: 'ldap', type: 'ldap', expectedRoute: 'vault.cluster.secrets.backend.ldap.configuration' },
  { label: 'pki', type: 'pki', expectedRoute: 'vault.cluster.secrets.backend.pki.configuration' },
  {
    label: 'ssh',
    type: 'ssh',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.plugin-settings',
  },
  {
    label: 'totp',
    type: 'totp',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  {
    label: 'aws',
    type: 'aws',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.plugin-settings',
  },
  {
    label: 'consul',
    type: 'consul',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  {
    label: 'nomad',
    type: 'nomad',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  {
    label: 'rabbitmq',
    type: 'rabbitmq',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
  {
    label: 'database',
    type: 'database',
    expectedRoute: 'vault.cluster.secrets.backend.configuration.general-settings',
  },
];

module('Integration | Component | manage-dropdown | Configure link', function (hooks) {
  setupRenderingTest(hooks);

  const makeModel = ({ type, version, id }) => {
    const options = version ? { version } : undefined;
    return new SecretsEngineResource({
      ...DEFAULT_MOUNT_DATA,
      path: `${id}/`,
      type,
      options,
    });
  };

  hooks.beforeEach(function () {
    const router = this.owner.lookup('service:router');
    this.transitionStub = sinon.stub(router, 'transitionTo');
    this.refreshStub = sinon.stub(router, 'refresh');
    this.currentRouteStub = sinon.stub(router, 'currentRouteName');
    const api = this.owner.lookup('service:api');
    this.mountDisableApiStub = sinon.stub(api.sys, 'mountsDisableSecretsEngine');
  });

  hooks.afterEach(function () {
    this.transitionStub.restore();
    this.refreshStub.restore();
    this.mountDisableApiStub.restore();
  });

  test('it disables a mount', async function (assert) {
    this.model = makeModel({ type: 'ldap', id: 'ldap' });
    await render(
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    await click(GENERAL.confirmButton);
    const [id] = this.mountDisableApiStub.lastCall.args;
    assert.strictEqual(id, 'ldap', 'it calls disable with the secret engine id');
  });

  test('it calls refresh() when current route is secrets.backends', async function (assert) {
    this.currentRouteStub.value('vault.cluster.secrets.backends');
    this.model = makeModel({ type: 'ldap', id: 'ldap' });
    await render(
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    await click(GENERAL.confirmButton);
    assert.true(
      this.refreshStub.calledOnce,
      'refresh is called because the current route is vault.cluster.secrets.backends'
    );
    assert.true(this.transitionStub.notCalled, 'transitionTo is not called');
  });

  test('it calls transitionTo() when current route is NOT secrets.backends', async function (assert) {
    this.currentRouteStub.value('vault.cluster.secrets.backend.ldap.overview');
    this.model = makeModel({ type: 'ldap', id: 'ldap' });
    await render(
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    await click(GENERAL.confirmButton);
    assert.true(this.transitionStub.calledOnce, 'transitionTo() is called');
    assert.true(this.refreshStub.notCalled, 'refresh() is not called');
  });

  TEST_CASES.forEach(({ label, type, version, expectedRoute }) => {
    test(`Configure link routes correctly for ${label}`, async function (assert) {
      const routing = this.owner.lookup('service:-routing');
      const transitionSpy = sinon.stub(routing, 'transitionTo');
      const id = `${label}-integration-test`;
      this.model = makeModel({ type, version, id });

      await render(
        hbs`<ManageDropdown @model={{this.model}} @variant="icon" @configRoute={{this.model.backendConfigurationLink}} />`
      );

      await click(GENERAL.menuTrigger);
      await click(GENERAL.menuItem('Configure'));

      assert.true(transitionSpy.called, `Configure action for ${label} triggers a route transition`);
      assert.strictEqual(
        transitionSpy.firstCall.args[0],
        expectedRoute,
        `Configure action for ${label} transitions to ${expectedRoute}`
      );
      assert.true(
        JSON.stringify(transitionSpy.firstCall.args).includes(id),
        `Configure action for ${label} includes model id ${id}`
      );

      transitionSpy.restore();
    });
  });
});
