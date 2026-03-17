/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | kv-v2 | Page::Secret::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.secretData = { foo: 'bar' };
    this.secretDataNested = {
      foo: {
        bar: 'baz',
      },
    };
    this.secret = {
      secretData: this.secretData,
      version: 1,
      destroyed: false,
      deletion_time: '',
      created_time: '2023-07-20T02:12:17.379762Z',
      custom_metadata: null,
    };
    this.metadata = {
      current_version: 1,
      updated_time: '2023-07-21T03:11:58.095971Z',
      versions: {
        1: {
          created_time: '2018-03-22T02:24:06.945319214Z',
          deletion_time: '',
          destroyed: false,
        },
      },
    };
    this.capabilities = { canReadData: true, canReadMetadata: true, canUpdateData: true };
    this.isPatchAllowed = true;
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: this.path },
    ];

    this.api = this.owner.lookup('service:api');
    this.queryParamsStub = sinon.stub(this.api, 'addQueryParams');
    this.syncStub = sinon
      .stub(this.api.sys, 'systemReadSyncAssociationsDestinations')
      .callsFake((initOverride) => {
        initOverride();
        return Promise.reject(getErrorResponse());
      });

    this.renderComponent = () =>
      render(
        hbs`
          <Page::Secret::Details
            @backend={{this.backend}}
            @breadcrumbs={{this.breadcrumbs}}
            @capabilities={{this.capabilities}}
            @isPatchAllowed={{this.isPatchAllowed}}
            @metadata={{this.metadata}}
            @path={{this.path}}
            @secret={{this.secret}}
          />
        `,
        { owner: this.engine }
      );
  });

  test('it renders secret details and toggles json view', async function (assert) {
    assert.expect(9);

    await this.renderComponent();
    assert.true(this.syncStub.calledOnce, 'sync status request made');
    assert.deepEqual(
      this.queryParamsStub.lastCall.args[1],
      { mount: this.backend, secret_name: this.path },
      'sync query params include mount and secret name'
    );
    assert
      .dom(PAGE.detail.syncAlert())
      .doesNotExist('sync page alert banner does not render when sync status errors');
    assert.dom(PAGE.title).includesText(this.path, 'renders secret path as page title');
    assert.dom(PAGE.infoRowValue('foo')).exists('renders row for secret data');
    assert.dom(PAGE.infoRowValue('foo')).hasText('***********');
    await click(GENERAL.button('toggle-masked'));
    assert.dom(PAGE.infoRowValue('foo')).hasText('bar', 'renders secret value');
    await click(GENERAL.toggleInput('json'));
    assert.dom(GENERAL.codeBlock('secret-data')).hasText(
      `Version data {
  "foo": "bar"
}`,
      'json editor renders secret data'
    );
    assert
      .dom(PAGE.detail.versionTimestamp)
      .includesText(`Version ${this.secret.version} created`, 'renders version and time created');
  });

  test('it renders hds codeblock view when secret is complex', async function (assert) {
    assert.expect(4);
    this.secret.secretData = this.secretDataNested;
    await this.renderComponent();
    assert.dom(PAGE.infoRowValue('foo')).doesNotExist('does not render rows of secret data');
    assert.dom(GENERAL.toggleInput('json')).isChecked();
    assert.dom(GENERAL.toggleInput('json')).isNotDisabled();
    assert.dom(GENERAL.codeBlock('secret-data')).exists('hds codeBlock exists');
  });

  test('it renders deleted empty state', async function (assert) {
    assert.expect(3);
    this.secret.deletion_time = '2023-07-23T02:12:17.379762Z';
    await this.renderComponent();

    assert.dom(PAGE.emptyStateTitle).hasText('Version 1 of this secret has been deleted');
    assert
      .dom(PAGE.emptyStateMessage)
      .hasText(
        'This version has been deleted but can be undeleted. View other versions of this secret by clicking the Version History tab above.'
      );
    assert
      .dom(PAGE.detail.versionTimestamp)
      .includesText(`Version ${this.secret.version} deleted`, 'renders version and time deleted');
  });

  test('it renders destroyed empty state', async function (assert) {
    assert.expect(2);
    this.secret.destroyed = true;
    await this.renderComponent();

    assert.dom(PAGE.emptyStateTitle).hasText('Version 1 of this secret has been permanently destroyed');
    assert
      .dom(PAGE.emptyStateMessage)
      .hasText(
        'A version that has been permanently deleted cannot be restored. You can view other versions of this secret in the Version History tab above.'
      );
  });

  test('it renders secret version dropdown', async function (assert) {
    assert.expect(6);
    this.metadata.versions[2] = {
      created_time: '2023-07-20T02:15:35.86465Z',
      deletion_time: '2023-07-25T00:36:19.950545Z',
      destroyed: false,
    };
    await this.renderComponent();

    assert.dom(PAGE.detail.versionTimestamp).includesText(this.secret.version, 'renders version');
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
      .dom(`${PAGE.detail.version(this.metadata.current_version)} [data-test-icon="check-circle"]`)
      .exists('renders current version icon');
  });

  test('it renders sync status page alert and refreshes', async function (assert) {
    assert.expect(3);

    this.syncStub.resolves({
      associated_destinations: {
        'aws-sm': {
          sync_status: 'SYNCED',
          name: 'my-destination',
          type: 'aws-sm',
          updated_at: '2023-09-01T12:00:00Z',
        },
      },
    });

    await this.renderComponent();
    assert
      .dom(PAGE.detail.syncAlert('my-destination'))
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
    assert.true(this.syncStub.calledTwice, 'sync status request made on refresh click');
  });

  test('it makes request to wrap a secret', async function (assert) {
    const wrapStub = sinon.stub(this.api.sys, 'wrap').resolves({ wrap_info: { token: 'hvs.token' } });

    await this.renderComponent();

    await click(PAGE.detail.copy);
    await click(GENERAL.button('wrap'));

    const { secretData: data, ...metadata } = this.secret;
    assert.true(
      wrapStub.calledWith({ data, metadata }, { headers: { 'X-Vault-Wrap-TTL': 1800 } }),
      'makes request to wrap secret with correct data'
    );
  });

  test('it renders sync status page alert for multiple destinations', async function (assert) {
    assert.expect(3);

    this.syncStub.resolves({
      associated_destinations: {
        'aws-sm': {
          sync_status: 'SYNCED',
          name: 'aws-dest',
          type: 'aws-sm',
          updated_at: '2023-09-01T12:00:00Z',
        },
        'gh-dest': {
          sync_status: 'SYNCING',
          name: 'gh-dest',
          type: 'gh',
          updated_at: '2023-09-01T12:00:00Z',
        },
      },
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
