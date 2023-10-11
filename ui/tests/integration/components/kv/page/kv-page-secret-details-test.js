/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, find, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvDataPath, kvMetadataPath } from 'vault/utils/kv-path';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { FORM, PAGE, parseJsonEditor } from 'vault/tests/helpers/kv/kv-selectors';

module('Integration | Component | kv-v2 | Page::Secret::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.version = 2;
    this.dataId = kvDataPath(this.backend, this.path);
    this.metadataId = kvMetadataPath(this.backend, this.path);

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
    });

    const metadata = this.server.create('kv-metadatum');
    metadata.id = this.metadataId;
    this.store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...metadata,
    });

    this.metadata = this.store.peekRecord('kv/metadata', this.metadataId);
    this.secret = this.store.peekRecord('kv/data', this.dataId);

    // this is the route model, not an ember data model
    this.model = {
      backend: this.backend,
      path: this.path,
      secret: this.secret,
      metadata: this.metadata,
    };
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'list' },
      { label: this.model.path },
    ];
  });

  test('it renders secret details and toggles json view', async function (assert) {
    assert.expect(6);

    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(PAGE.title).includesText(this.model.path, 'renders secret path as page title');
    assert.dom(PAGE.infoRowValue('foo')).exists('renders row for secret data');
    assert.dom(PAGE.infoRowValue('foo')).hasText('***********');
    await click(FORM.toggleMasked);
    assert.dom(PAGE.infoRowValue('foo')).hasText('bar', 'renders secret value');
    await click(FORM.toggleJson);
    assert.propEqual(parseJsonEditor(find), this.secretData, 'json editor renders secret data');
    assert
      .dom(PAGE.detail.versionTimestamp)
      .includesText(`Version ${this.version} created`, 'renders version and time created');
  });

  test('it renders deleted empty state', async function (assert) {
    assert.expect(3);
    this.secret.deletionTime = '2023-07-23T02:12:17.379762Z';
    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
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
    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    assert.dom(PAGE.emptyStateTitle).hasText('Version 2 of this secret has been permanently destroyed');
    assert
      .dom(PAGE.emptyStateMessage)
      .hasText(
        'A version that has been permanently deleted cannot be restored. You can view other versions of this secret in the Version History tab above.'
      );
  });

  test('it renders secret version dropdown', async function (assert) {
    assert.expect(9);

    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

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
});
