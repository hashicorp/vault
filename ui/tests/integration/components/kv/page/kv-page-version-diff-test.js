/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
    this.store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.server.get(`/kv-engine/data/my-secret?version=1`, () => {
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

    const metadata = this.server.create('kv-metadatum');
    metadata.id = kvMetadataPath(this.backend, this.secret);
    this.store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...metadata,
    });

    this.backend = 'kv-engine';
    this.secret = 'my-secret';
    this.metadata = this.store.peekRecord('kv/metadata', metadata.id);
    this.breadcrumbs = [{ label: 'version history', route: 'secret.metadata.versions' }, { label: 'diff' }];

    // compare version 4
    const dataId = kvDataPath(this.backend, this.secret, 4);
    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: dataId,
      secret_data: { foo: 'bar' },
      created_time: '2023-07-20T02:12:17.379762Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: 4,
    });
  });

  test('it renders compared data of the two versions and shows icons for deleted, destroyed and current', async function (assert) {
    assert.expect(12);

    await render(
      hbs`
       <Page::Secret::Metadata::VersionDiff
        @metadata={{this.metadata}} 
        @path={{this.secret}}
        @backend={{this.backend}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    const [left, right] = findAll(PAGE.detail.versionDropdown);
    assert.dom(PAGE.diff.visualDiff).hasText(
      `{ \"foo\": \"bar\" }`, // eslint-disable-line no-useless-escape
      'correctly pull in the data from version 4 and compared to version 1.'
    );
    assert.dom(PAGE.diff.deleted).hasText(`foo"bar"`);
    assert.dom(PAGE.diff.added).hasText(`foo3"bar3"`);
    assert.dom(left).hasText('Version 4', 'shows the current version for the left side default version.');
    assert
      .dom(right)
      .hasText(
        'Version 1',
        'shows the first version that is not deleted or destroyed for the right version on init.'
      );

    await click(right);

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
