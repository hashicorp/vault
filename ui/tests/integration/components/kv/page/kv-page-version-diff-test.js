/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvMetadataPath, kvDataPath } from 'vault/utils/kv-path';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

module('Integration | Component | kv | Page::Secret::Metadata::VersionDiff', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.breadcrumbs = [{ label: 'version history', route: 'secret.metadata.versions' }, { label: 'diff' }];

    this.store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());

    const metadata = this.server.create('kv-metadatum');
    metadata.id = kvMetadataPath(this.backend, this.path);
    this.store.pushPayload('kv/metadata', { modelName: 'kv/metadata', ...metadata });
    this.metadata = this.store.peekRecord('kv/metadata', metadata.id);
    // push current secret version record into the store to assert only one request is made
    const dataId = kvDataPath(this.backend, this.path, 4);
    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: dataId,
      secret_data: { foo: 'bar' },
      version: this.metadata.currentVersion,
    });
  });

  test('it renders empty states when current version is deleted or destroyed', async function (assert) {
    assert.expect(4);
    this.server.get(`/${this.backend}/data/${this.path}`, () => {});
    const { currentVersion } = this.metadata;

    // destroyed
    this.metadata.versions[currentVersion].destroyed = true;
    await render(
      hbs`
       <Page::Secret::Metadata::VersionDiff
        @metadata={{this.metadata}} 
        @path={{this.path}}
        @backend={{this.backend}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    assert.dom(PAGE.emptyStateTitle).hasText(`Version ${currentVersion} has been destroyed`);
    assert
      .dom(PAGE.emptyStateMessage)
      .hasText('The current version of this secret has been destroyed. Select another version to compare.');

    // deleted
    this.metadata.versions[currentVersion].destroyed = false;
    this.metadata.versions[currentVersion].deletion_time = '2023-07-25T00:36:19.950545Z';
    await render(
      hbs`
       <Page::Secret::Metadata::VersionDiff
        @metadata={{this.metadata}} 
        @path={{this.path}}
        @backend={{this.backend}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(PAGE.emptyStateTitle).hasText(`Version ${currentVersion} has been deleted`);
    assert
      .dom(PAGE.emptyStateMessage)
      .hasText('The current version of this secret has been deleted. Select another version to compare.');
  });

  test('it renders compared data of the two versions and shows icons for deleted, destroyed and current', async function (assert) {
    assert.expect(14);
    this.server.get(`/${this.backend}/data/${this.path}`, (schema, req) => {
      assert.ok('request made to the fetch version 1 data.');
      // request should not be made for version 4 (current version) because that record already exists in the store
      assert.strictEqual(req.queryParams.version, '1', 'request includes version param');
      return {
        request_id: 'foobar',
        data: {
          data: { hello: 'world' },
          metadata: {
            created_time: '2023-06-20T21:26:47.592306Z',
            custom_metadata: null,
            deletion_time: '',
            destroyed: false,
            version: 1,
          },
        },
      };
    });

    await render(
      hbs`
       <Page::Secret::Metadata::VersionDiff
        @metadata={{this.metadata}} 
        @path={{this.path}}
        @backend={{this.backend}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    const [left, right] = findAll(PAGE.detail.versionDropdown);
    assert.dom(PAGE.diff.visualDiff).hasText(
      `foo\"bar\"hello\"world\"`, // eslint-disable-line no-useless-escape
      'correctly pull in the data from version 4 and compared to version 1.'
    );
    assert.dom(PAGE.diff.deleted).hasText(`hello"world"`);
    assert.dom(PAGE.diff.added).hasText(`foo"bar"`);
    assert.dom(right).hasText('Version 4', 'shows the current version for the left side default version.');
    assert.dom(left).hasText('Version 1', 'shows the latest active version on init.');

    await click(left);

    for (const num in this.metadata.versions) {
      const data = this.metadata.versions[num];
      assert.dom(PAGE.detail.version(num)).exists('renders the button for each version.');

      if (data.destroyed || data.deletion_time) {
        assert
          .dom(`${PAGE.detail.version(num)} [data-test-icon="x-square-fill"]`)
          .hasClass(`${data.destroyed ? 'has-text-danger' : 'has-text-grey'}`);
      }
    }
    assert
      .dom(`${PAGE.detail.version('1')} button`)
      .hasClass('is-active', 'correctly shows the selected version 1 as active.');
  });
});
