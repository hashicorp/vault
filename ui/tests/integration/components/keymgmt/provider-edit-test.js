/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, settled, fillIn } from '@ember/test-helpers';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const ts = 'data-test-kms-provider';
const root = {
  path: 'vault.cluster.secrets.backend.list-root',
  model: 'keymgmt',
  label: 'keymgmt',
  text: 'keymgmt',
};

module('Integration | Component | keymgmt/provider-edit', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.store.push({
      data: {
        id: 'foo-bar',
        type: 'keymgmt/provider',
        attributes: {
          name: 'foo-bar',
          provider: 'azurekeyvault',
          keyCollection: 'keyvault-1',
          backend: 'keymgmt',
        },
      },
    });
    this.model = this.store.peekRecord('keymgmt/provider', 'foo-bar');
    this.root = root;
    this.owner.lookup('service:router').reopen({
      currentURL: '/ui/vault/secrets/keymgmt/show/foo-bar',
      currentRouteName: 'secrets.keymgmt.provider.show',
      urlFor() {
        return '';
      },
    });

    setRunOptions({
      rules: {
        // TODO: fix KMS provider-edit [data-test-kms-provider-delete] violates this rule
        // see https://dequeuniversity.com/rules/axe/4.8/scrollable-region-focusable
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });

  test('it should render show view', async function (assert) {
    assert.expect(10);

    // override capability getters
    Object.defineProperties(this.model, {
      canDelete: { value: true },
      canListKeys: { value: true },
    });

    this.server.post('/sys/capabilities-self', () => ({}));
    this.server.get('/keymgmt/kms/foo-bar/key', () => {
      return {
        data: {
          keys: ['testkey-1', 'testkey-2'],
        },
      };
    });

    const changeTab = async (tab) => {
      this.set('tab', tab);
      await settled();
    };

    await render(hbs`
      <Keymgmt::ProviderEdit
        @root={{this.root}}
        @model={{this.model}}
        @mode="show"
        @tab={{this.tab}}
      />`);

    assert.dom(`[${ts}-header]`).hasText('Provider foo-bar', 'Page header renders');
    assert.dom(`[${ts}-tab="details"]`).hasClass('active', 'Details tab is active');

    const infoRows = this.element.querySelectorAll('[data-test-component="info-table-row"]');
    assert.dom(infoRows[0]).hasText('Provider name foo-bar', 'Provider name field renders');
    assert.dom(infoRows[1]).hasText('Type Azure Key Vault', 'Type field renders');
    assert.dom('svg', infoRows[1]).hasAttribute('data-test-icon', 'azure-color', 'Icon renders for type');
    assert.dom(infoRows[2]).hasText('Key Vault instance name keyvault-1', 'Key collection field renders');
    assert.dom(infoRows[3]).hasText('Keys 2 keys', 'Keys field renders');

    await changeTab('keys');
    assert.dom(`[${ts}-details-actions]`).doesNotExist('Toolbar is hidden on keys tab');
    assert.dom('[data-test-secret-link]').exists({ count: 2 }, 'Keys list renders');

    await changeTab('details');
    await click(`[${ts}-delete]`);
    assert
      .dom('[data-test-confirm-action-message]')
      .hasText(
        'This provider cannot be deleted until all 2 key(s) distributed to it are revoked. This can be done from the Keys tab.',
        'Renders disabled message'
      );
    await click('[data-test-confirm-cancel-button]');
  });

  test('it should delete a provider', async function (assert) {
    assert.expect(5);

    // override capability getters
    Object.defineProperties(this.model, {
      canDelete: { value: true },
      canListKeys: { value: true },
    });

    this.server.post('/sys/capabilities-self', () => ({}));
    this.server.get('/keymgmt/kms/foo-bar/key', () => {
      return {
        data: {
          keys: [],
        },
      };
    });
    this.server.delete('/keymgmt/kms/foo-bar', () => {
      assert.ok(true, 'Request made to delete key');
      return {};
    });
    this.owner.lookup('service:router').reopen({
      transitionTo(path, model, { queryParams: { tab } }) {
        assert.strictEqual(path, root.path, 'Root path sent in transitionTo on delete');
        assert.strictEqual(model, root.model, 'Root model sent in transitionTo on delete');
        assert.deepEqual(tab, 'provider', 'Correct query params sent in transitionTo on delete');
      },
    });

    await render(hbs`
      <Keymgmt::ProviderEdit
        @root={{this.root}}
        @model={{this.model}}
        @mode="show"
        @tab={{this.tab}}
      />`);

    assert
      .dom('[data-test-value-div="Keys"]')
      .hasText('None', 'None is displayed when no keys exist for provider');

    await click(`[${ts}-delete]`);
    await click('[data-test-confirm-button]');
  });

  test('it should render create view', async function (assert) {
    assert.expect(14);

    this.server.put('/keymgmt/kms/foo', (schema, req) => {
      const params = {
        name: 'foo',
        backend: 'keymgmt',
        provider: 'gcpckms',
        key_collection: 'keyvault-1',
        credentials: {
          service_account_file: 'test',
        },
      };
      assert.deepEqual(JSON.parse(req.requestBody), params, 'PUT request made with correct data');
      return {};
    });
    this.owner.lookup('service:router').reopen({
      transitionTo(path, model, { queryParams: { itemType } }) {
        assert.strictEqual(
          path,
          'vault.cluster.secrets.backend.show',
          'Show route sent in transitionTo on save'
        );
        assert.strictEqual(model, 'foo', 'Model id sent in transitionTo on save');
        assert.deepEqual(itemType, 'provider', 'Correct query params sent in transitionTo on save');
      },
    });
    this.model = this.store.createRecord('keymgmt/provider', { backend: 'keymgmt' });

    await render(hbs`
      <Keymgmt::ProviderEdit
        @root={{this.root}}
        @model={{this.model}}
        @mode="create"
      />`);

    assert.dom(`[${ts}-header]`).hasText('Create Provider', 'Page header renders');
    assert.dom(`[${ts}-config-title]`).exists('Config header shown in create mode');
    assert.dom(`[${ts}-creds-title]`).doesNotExist('New credentials header hidden in create mode');

    await click(`[${ts}-submit]`);
    assert.dom('[data-test-inline-error-message]').exists('Validation error messages shown');

    await fillIn('[data-test-input="provider"]', 'azurekeyvault');
    ['client_id', 'client_secret', 'tenant_id'].forEach((prop) => {
      assert.dom(`[data-test-input="credentials.${prop}"]`).exists(`Azure - ${prop} field renders`);
    });

    await fillIn('[data-test-input="provider"]', 'awskms');
    ['access_key', 'secret_key'].forEach((prop) => {
      assert.dom(`[data-test-input="credentials.${prop}"]`).exists(`AWS - ${prop} field renders`);
    });

    await fillIn('[data-test-input="provider"]', 'gcpckms');
    assert.dom(`[data-test-input="credentials.service_account_file"]`).exists(`GCP - cred field renders`);

    await fillIn('[data-test-input="name"]', 'foo');
    await fillIn('[data-test-input="keyCollection"]', 'keyvault-1');
    await fillIn('[data-test-input="credentials.service_account_file"]', 'test');
    await click(`[${ts}-submit]`);
  });

  test('it should render edit view', async function (assert) {
    assert.expect(3);

    this.server.put('/keymgmt/kms/foo', (schema, req) => {
      const params = {
        name: 'foo-bar',
        backend: 'keymgmt',
        provider: 'azurekeyvault',
        key_collection: 'keyvault-1',
        credentials: {
          client_id: 'client_id test',
          client_secret: 'client_secret test',
          tenant_id: 'tenant_id test',
        },
      };
      assert.deepEqual(JSON.parse(req.requestBody), params, 'PUT request made with correct data');
      return {};
    });
    this.owner.lookup('service:router').reopen({
      transitionTo(path, model, { queryParams: { itemType } }) {
        assert.strictEqual(
          path,
          'vault.cluster.secrets.backend.show',
          'Show route sent in transitionTo on save'
        );
        assert.strictEqual(model, 'foo', 'Model id sent in transitionTo on save');
        assert.deepEqual(itemType, 'provider', 'Correct query params sent in transitionTo on save');
      },
    });
    await render(hbs`
      <Keymgmt::ProviderEdit
        @root={{this.root}}
        @model={{this.model}}
        @mode="edit"
      />`);

    assert.dom(`[${ts}-header]`).hasText('Update Credentials', 'Page header renders');
    assert.dom(`[${ts}-config-title]`).doesNotExist('Config header hidden in edit mode');
    assert.dom(`[${ts}-creds-title]`).exists('New credentials header shown in edit mode');

    for (const prop of ['client_id', 'client_secret', 'tenant_id']) {
      await fillIn(`[data-test-input="credentials.${prop}"]`, `${prop} test`);
    }
    await click(`[${ts}-submit]`);
  });
});
