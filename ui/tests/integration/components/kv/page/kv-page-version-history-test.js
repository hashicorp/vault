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
import { PAGE } from 'vault/tests/helpers/kv/kv-page-selectors';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

module('Integration | Component | kv | Page::Secret::Version-History', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    const store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    const metadata = this.server.create('kv-metadatum');
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
    assert.expect(8);

    await render(
      hbs`
       <Page::Secret::VersionHistory
        @path={{this.metadata.path}}
        @metadata={{this.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    for (const version in this.metadata.versions) {
      const data = this.metadata.versions[version];
      assert.dom(PAGE.list.linkedBlock(version)).exists(`renders the linked blocks for each version`);

      if (data.destroyed) {
        assert
          .dom(`${PAGE.list.icon(version)} [data-test-icon="x-square-fill"]`)
          .hasStyle({ color: 'rgb(199, 52, 69)' });
      }
      if (data.deletion_time) {
        assert
          .dom(`${PAGE.list.icon(version)} [data-test-icon="x-square-fill"]`)
          .hasStyle({ color: 'rgb(111, 118, 130)' });
      }
    }

    assert
      .dom(`${PAGE.list.icon(this.metadata.currentVersion)} [data-test-icon="check-circle-fill"]`)
      .exists('renders the current version');
  });

  test('it gives the option to create a new version from a secret from the popup menu', async function (assert) {
    assert.expect(1);
    await render(
      hbs`
       <Page::Secret::VersionHistory
        @path={{this.metadata.path}}
        @metadata={{this.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    // because the popup menu is nested in a linked block we must combine the two selectors
    const popupSelector = `${PAGE.list.linkedBlock(2)} ${PAGE.list.popup}`;
    await click(popupSelector);
    assert
      .dom('[data-test-create-new-version-from="2"]')
      .exists('Shows the option to create a new version from that secret.');
  });
});
