/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvMetadataPath } from 'vault/utils/kv-path';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

module('Integration | Component | kv | Page::Secret::Metadata::Version-History', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    const store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    const metadata = this.server.create('kv-metadatum');
    // we want to test a scenario where the current version is also destroyed so there are two icons.
    // we override the mirage factory to account for this use case.
    metadata.data.versions[4] = {
      created_time: '2023-07-21T03:11:58.095971Z',
      deletion_time: '',
      destroyed: true,
    };
    metadata.id = kvMetadataPath('kv-engine', 'my-secret');
    store.pushPayload('kv/metadatum', {
      modelName: 'kv/metadata',
      ...metadata,
    });
    this.metadata = store.peekRecord('kv/metadata', metadata.id);
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.metadata.backend, route: 'list' },
      { label: this.metadata.path, route: 'secret.details', model: this.metadata.path },
      { label: 'version history' },
    ];
  });

  test('it renders version history and shows icons for deleted, destroyed and current', async function (assert) {
    assert.expect(7); // 4 linked blocks, 2 destroyed, 1 deleted.

    await render(
      hbs`
       <Page::Secret::Metadata::VersionHistory
        @path={{this.metadata.path}}
        @metadata={{this.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    for (const version in this.metadata.versions) {
      const data = this.metadata.versions[version];
      assert.dom(PAGE.versions.linkedBlock(version)).exists(`renders the linked blocks for each version`);

      if (data.destroyed) {
        assert
          .dom(`${PAGE.versions.icon(version)} [data-test-icon="x-square-fill"]`)
          .hasStyle({ color: 'rgb(199, 52, 69)' });
      }
      if (data.deletion_time) {
        assert
          .dom(`${PAGE.versions.icon(version)} [data-test-icon="x-square-fill"]`)
          .hasStyle({ color: 'rgb(111, 118, 130)' });
      }
    }
  });

  test('it gives the option to create a new version from a secret from the popup menu', async function (assert) {
    assert.expect(1);
    await render(
      hbs`
       <Page::Secret::Metadata::VersionHistory
        @path={{this.metadata.path}}
        @metadata={{this.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    // because the popup menu is nested in a linked block we must combine the two selectors
    const popupSelector = `${PAGE.versions.linkedBlock(2)} ${PAGE.popup}`;
    await click(popupSelector);
    assert
      .dom('[data-test-create-new-version-from="2"]')
      .exists('Shows the option to create a new version from that secret.');
  });
});
