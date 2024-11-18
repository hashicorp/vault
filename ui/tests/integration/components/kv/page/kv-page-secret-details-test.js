/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvDataPath } from 'vault/utils/kv-path';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { syncStatusResponse } from 'vault/mirage/handlers/sync';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { baseSetup } from 'vault/tests/helpers/kv/kv-run-commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | kv-v2 | Page::Secret::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    baseSetup(this);
    this.pathComplex = 'my-secret-object';
    this.version = 2;
    this.dataId = kvDataPath(this.backend, this.path);
    this.dataIdComplex = kvDataPath(this.backend, this.pathComplex);
    this.secretData = { foo: 'bar' };
    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.dataId,
      secret_data: this.secretData,
      created_time: '2023-07-20T02:12:17.379762Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: this.version,
      backend: this.backend,
      path: this.path,
    });
    // nested secret
    this.secretDataComplex = {
      foo: {
        bar: 'baz',
      },
    };
    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.dataIdComplex,
      secret_data: this.secretDataComplex,
      created_time: '2023-08-20T02:12:17.379762Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: this.version,
    });
    this.secret = this.store.peekRecord('kv/data', this.dataId);
    this.secretComplex = this.store.peekRecord('kv/data', this.dataIdComplex);
    // this is the route model, not an ember data model
    this.model = {
      backend: this.backend,
      // permissions are tested in navigation acceptance test, so just stub as all true here
      canReadData: true,
      canReadMetadata: true,
      canUpdateData: true,
      isPatchAllowed: true,
      metadata: this.metadata,
      path: this.path,
      secret: this.secret,
    };
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'list' },
      { label: this.model.path },
    ];
    this.modelComplex = {
      backend: this.backend,
      path: this.pathComplex,
      secret: this.secretComplex,
      metadata: this.metadata,
    };
    this.renderComponent = (model) => {
      this.model = model ? { ...this.model, ...model } : this.model;
      return render(
        hbs`
      <Page::Secret::Details
        @backend={{this.model.backend}}
        @breadcrumbs={{this.breadcrumbs}}
        @canReadData={{this.model.canReadData}}
        @canReadMetadata={{this.model.canReadMetadata}}
        @canUpdateData={{this.model.canUpdateData}}
        @isPatchAllowed={{this.model.isPatchAllowed}}
        @metadata={{this.model.metadata}}
        @path={{this.model.path}}
        @secret={{this.model.secret}}
      />
      `,
        { owner: this.engine }
      );
    };
  });

  test('it renders secret details and toggles json view', async function (assert) {
    assert.expect(9);
    this.server.get(`sys/sync/associations/destinations`, (schema, req) => {
      assert.ok(true, 'request made to fetch sync status');
      assert.propEqual(
        req.queryParams,
        {
          mount: this.backend,
          secret_name: this.path,
        },
        'query params include mount and secret name'
      );
      // no records so response returns 404
      return syncStatusResponse(schema, req);
    });
    await this.renderComponent();
    assert
      .dom(PAGE.detail.syncAlert())
      .doesNotExist('sync page alert banner does not render when sync status errors');
    assert.dom(PAGE.title).includesText(this.model.path, 'renders secret path as page title');
    assert.dom(PAGE.infoRowValue('foo')).exists('renders row for secret data');
    assert.dom(PAGE.infoRowValue('foo')).hasText('***********');
    await click(FORM.toggleMasked);
    assert.dom(PAGE.infoRowValue('foo')).hasText('bar', 'renders secret value');
    await click(FORM.toggleJson);
    assert.dom(GENERAL.codeBlock('secret-data')).hasText(
      `Version data {
  "foo": "bar"
}`,
      'json editor renders secret data'
    );
    assert
      .dom(PAGE.detail.versionTimestamp)
      .includesText(`Version ${this.version} created`, 'renders version and time created');
  });

  test('it renders hds codeblock view when secret is complex', async function (assert) {
    assert.expect(4);
    await this.renderComponent(this.modelComplex);
    assert.dom(PAGE.infoRowValue('foo')).doesNotExist('does not render rows of secret data');
    assert.dom(FORM.toggleJson).isChecked();
    assert.dom(FORM.toggleJson).isNotDisabled();
    assert.dom(GENERAL.codeBlock('secret-data')).exists('hds codeBlock exists');
  });

  test('it renders deleted empty state', async function (assert) {
    assert.expect(3);
    this.secret.deletionTime = '2023-07-23T02:12:17.379762Z';
    await this.renderComponent();

    assert.dom(PAGE.emptyStateTitle).hasText('Version 2 of this secret has been deleted');
    assert
      .dom(PAGE.emptyStateMessage)
      .hasText(
        'This version has been deleted but can be undeleted. View other versions of this secret by clicking the Version History tab above.'
      );
    assert
      .dom(PAGE.detail.versionTimestamp)
      .includesText(`Version ${this.version} deleted`, 'renders version and time deleted');
  });

  test('it renders destroyed empty state', async function (assert) {
    assert.expect(2);
    this.secret.destroyed = true;
    await this.renderComponent();

    assert.dom(PAGE.emptyStateTitle).hasText('Version 2 of this secret has been permanently destroyed');
    assert
      .dom(PAGE.emptyStateMessage)
      .hasText(
        'A version that has been permanently deleted cannot be restored. You can view other versions of this secret in the Version History tab above.'
      );
  });

  test('it renders secret version dropdown', async function (assert) {
    assert.expect(9);
    await this.renderComponent();

    assert.dom(PAGE.detail.versionTimestamp).includesText(this.version, 'renders version');
    assert.dom(PAGE.detail.versionDropdown).hasText(`Version ${this.secret.version}`);
    await click(PAGE.detail.versionDropdown);

    for (const version in this.metadata.versions) {
      const data = this.metadata.versions[version];
      assert.dom(PAGE.detail.version(version)).exists(`renders ${version} in dropdown menu`);

      if (data.destroyed || data.deletion_time) {
        assert
          .dom(`${PAGE.detail.version(version)} [data-test-icon="x-square-fill"]`)
          .hasClass(`${data.destroyed ? 'has-text-danger' : 'has-text-grey'}`);
      }
    }

    assert
      .dom(`${PAGE.detail.version(this.metadata.currentVersion)} [data-test-icon="check-circle"]`)
      .exists('renders current version icon');
  });

  test('it renders sync status page alert and refreshes', async function (assert) {
    assert.expect(6); // assert count important because confirms request made to fetch sync status twice
    const destinationName = 'my-destination';
    this.server.create('sync-association', {
      type: 'aws-sm',
      name: destinationName,
      mount: this.backend,
      secret_name: this.path,
    });
    this.server.get(`sys/sync/associations/destinations`, (schema, req) => {
      // these assertions should be hit twice, once on init and again when the 'Refresh' button is clicked
      assert.ok(true, 'request made to fetch sync status');
      assert.propEqual(
        req.queryParams,
        {
          mount: this.backend,
          secret_name: this.path,
        },
        'query params include mount and secret name'
      );
      return syncStatusResponse(schema, req);
    });

    await this.renderComponent();
    assert
      .dom(PAGE.detail.syncAlert(destinationName))
      .hasTextContaining(
        'Synced my-destination - last updated September',
        'renders sync status alert banner'
      );
    assert
      .dom(PAGE.detail.syncAlert())
      .hasTextContaining(
        'This secret has been synced from Vault to 1 destination. Updates to this secret will automatically sync to its destination.',
        'renders alert header referring to singular destination'
      );
    // sync status refresh button
    await click(`${PAGE.detail.syncAlert()} button`);
  });

  test('it makes request to wrap a secret', async function (assert) {
    assert.expect(2);
    const url = `${encodePath(this.backend)}/data/${encodePath(this.path)}`;

    this.server.get(url, (schema, { requestHeaders }) => {
      assert.true(true, `GET request made to url: ${url}`);
      assert.strictEqual(requestHeaders['X-Vault-Wrap-TTL'], '1800', 'request header includes wrap ttl');
      return {
        data: null,
        token: 'hvs.token',
        accessor: 'nTgqnw3S4GMz8NKHsOhTBhlk',
        ttl: 1800,
        creation_time: '2024-07-26T10:20:32.359107-07:00',
        creation_path: `${this.backend}/data/${this.path}}`,
      };
    });
    await this.renderComponent();

    await click(PAGE.detail.copy);
    await click(PAGE.detail.wrap);
  });

  test('it renders sync status page alert for multiple destinations', async function (assert) {
    assert.expect(3); // assert count important because confirms request made to fetch sync status twice
    this.server.create('sync-association', {
      type: 'aws-sm',
      name: 'aws-dest',
      mount: this.backend,
      secret_name: this.path,
    });
    this.server.create('sync-association', {
      type: 'gh',
      name: 'gh-dest',
      mount: this.backend,
      secret_name: this.path,
    });
    this.server.get(`sys/sync/associations/destinations`, (schema, req) => {
      return syncStatusResponse(schema, req);
    });

    await this.renderComponent();

    assert
      .dom(PAGE.detail.syncAlert('aws-dest'))
      .hasTextContaining('Synced aws-dest - last updated September', 'renders status for aws destination');
    assert
      .dom(PAGE.detail.syncAlert('gh-dest'))
      .hasTextContaining('Syncing gh-dest - last updated September', 'renders status for gh destination');
    assert
      .dom(PAGE.detail.syncAlert())
      .hasTextContaining(
        'This secret has been synced from Vault to 2 destinations. Updates to this secret will automatically sync to its destinations.',
        'renders alert title referring to plural destinations'
      );
  });
});
