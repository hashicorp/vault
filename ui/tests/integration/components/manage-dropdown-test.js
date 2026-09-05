/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, render } from '@ember/test-helpers';
import { setupRenderingTest } from 'vault/tests/helpers';
import hbs from 'htmlbars-inline-precompile';
import { module, test } from 'qunit';
import sinon from 'sinon';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const CONFIRM_MODAL = '[data-test-confirm-modal]';

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

module('Integration | Component | manage-dropdown | Delete modal', function (hooks) {
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
    this.kvV1ListStub = sinon.stub(api.secrets, 'kvV1List').resolves({ keys: ['foo', 'bar'] });
    this.kvV2ListStub = sinon.stub(api.secrets, 'kvV2List').resolves({ keys: ['a', 'b', 'c'] });
  });

  hooks.afterEach(function () {
    this.transitionStub.restore();
    this.refreshStub.restore();
    this.mountDisableApiStub.restore();
    this.kvV1ListStub.restore();
    this.kvV2ListStub.restore();
  });

  test('it renders the delete modal with title and body on delete click', async function (assert) {
    this.model = makeModel({ type: 'ldap', id: 'my-ldap' });
    await render(
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @showDelete={{true}} @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    assert.dom(CONFIRM_MODAL).exists('delete modal is shown');
    assert
      .dom('[data-test-confirm-action-title]')
      .hasText('Delete secrets engine my-ldap?', 'modal title includes engine id');
    assert.dom('[data-test-confirm-action-message]').containsText('will be permanently deleted');
    assert.dom('[data-test-confirm-action-message]').containsText('my-ldap secrets engine');
    assert
      .dom('[data-test-confirm-action-message]')
      .containsText('Engine configuration and third-party integrations');
  });

  test('it shows secret count for KV v1 engines', async function (assert) {
    this.model = makeModel({ type: 'kv', version: 1, id: 'kv-v1' });
    await render(
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @showDelete={{true}} @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    assert.dom('[data-test-confirm-action-message]').containsText('2 secrets', 'shows KV v1 secret count');
  });

  test('it shows secret count for KV v2 engines', async function (assert) {
    this.model = makeModel({ type: 'kv', version: 2, id: 'kv-v2' });
    await render(
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @showDelete={{true}} @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    assert.dom('[data-test-confirm-action-message]').containsText('3 secrets', 'shows KV v2 secret count');
  });

  test('it does not fire onConfirm when wrong text is entered', async function (assert) {
    this.model = makeModel({ type: 'ldap', id: 'ldap' });
    await render(
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @showDelete={{true}} @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    await fillIn(GENERAL.confirmTextInput, 'wrong-text');
    await click(GENERAL.confirmButton);
    assert.dom(CONFIRM_MODAL).exists('modal stays open');
    assert.false(this.mountDisableApiStub.called, 'disable API is not called with wrong input');
    assert.dom(GENERAL.confirmWarning).exists('warning is shown after failed confirm attempt');
  });

  test('modal closes and resets on cancel', async function (assert) {
    this.model = makeModel({ type: 'ldap', id: 'ldap' });
    await render(
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @showDelete={{true}} @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    assert.dom(CONFIRM_MODAL).exists('modal is open');
    await click(GENERAL.cancelButton);
    assert.dom(CONFIRM_MODAL).doesNotExist('modal is closed after cancel');
  });
});

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
    this.kvV1ListStub = sinon.stub(api.secrets, 'kvV1List').resolves({ keys: [] });
    this.kvV2ListStub = sinon.stub(api.secrets, 'kvV2List').resolves({ keys: [] });
  });

  hooks.afterEach(function () {
    this.transitionStub.restore();
    this.refreshStub.restore();
    this.mountDisableApiStub.restore();
    this.kvV1ListStub.restore();
    this.kvV2ListStub.restore();
  });

  test('it disables a mount', async function (assert) {
    this.model = makeModel({ type: 'ldap', id: 'ldap' });
    await render(
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @showDelete={{true}} @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    assert.dom(CONFIRM_MODAL).exists('confirm modal is shown');
    await fillIn(GENERAL.confirmTextInput, 'delete-engine');
    await click(GENERAL.confirmButton);
    const [id] = this.mountDisableApiStub.lastCall.args;
    assert.strictEqual(id, 'ldap', 'it calls disable with the secret engine id');
  });

  test('it calls refresh() when current route is secrets.backends', async function (assert) {
    this.currentRouteStub.value('vault.cluster.secrets.backends');
    this.model = makeModel({ type: 'ldap', id: 'ldap' });
    await render(
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @showDelete={{true}} @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    assert.dom(CONFIRM_MODAL).exists('confirm modal is shown');
    await fillIn(GENERAL.confirmTextInput, 'delete-engine');
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
      hbs`<ManageDropdown @model={{this.model}} @variant="icon" @showDelete={{true}} @configRoute={{this.model.backendConfigurationLink}} />`
    );
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('Delete'));
    assert.dom(CONFIRM_MODAL).exists('confirm modal is shown');
    await fillIn(GENERAL.confirmTextInput, 'delete-engine');
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
        hbs`<ManageDropdown @model={{this.model}} @variant="icon" @showDelete={{true}} @configRoute={{this.model.backendConfigurationLink}} />`
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
