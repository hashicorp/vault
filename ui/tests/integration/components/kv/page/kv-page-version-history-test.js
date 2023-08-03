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
    assert.expect(5);

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
    const rows = document.querySelectorAll(PAGE.list.rows);
    rows.forEach((element, index) => {
      if (index === 0) {
        // using querySelector to search for nested classes instead of hasClass
        assert
          .dom(element.querySelector('.has-text-danger'))
          .exists('version 4 has the destroyed flight icon');
        assert
          .dom(element.querySelector('.has-text-success'))
          .exists('version 4 has the current version icon');
      }
      if (index === 1) {
        assert
          .dom(element.querySelector('.has-text-danger'))
          .exists('version 3 has the destroyed flight icon');
      }
      if (index === 2) {
        assert.dom(element.querySelector('.has-text-grey')).exists('version 2 has the deleted flight icon');
      }
      if (index === 3) {
        assert
          .dom(element.querySelector('.flight-icon'))
          .exists(
            { count: 1 },
            'only shows the version icon if secret is not deleted or destroyed or current'
          );
      }
    });
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
