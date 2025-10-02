/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import isDeleted from 'kv/helpers/is-deleted';

module('Integration | Component | kv-v2 | Page::Secret::Metadata::Version-History', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    // we want to test a scenario where the current version is also destroyed so there are two icons.
    this.metadata = {
      current_version: 4,
      updated_time: '2023-07-21T03:11:58.095971Z',
      versions: {
        1: {
          created_time: '2018-03-22T02:24:06.945319214Z',
          deletion_time: '',
          destroyed: false,
        },
        2: {
          created_time: '2023-07-20T02:15:35.86465Z',
          deletion_time: '2023-07-25T00:36:19.950545Z',
          destroyed: false,
        },
        3: {
          created_time: '2023-07-20T02:15:40.164549Z',
          deletion_time: '',
          destroyed: true,
        },
        4: {
          created_time: '2023-07-21T03:11:58.095971Z',
          deletion_time: '',
          destroyed: true,
        },
      },
    };
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: this.path, route: 'secret.details', model: this.path },
      { label: 'Version History' },
    ];
    this.capabilities = { canReadMetadata: true, canCreateVersionData: true };

    this.renderComponent = () =>
      render(
        hbs`
          <Page::Secret::Metadata::VersionHistory
            @metadata={{this.metadata}}
            @path={{this.path}}
            @backend={{this.backend}}
            @breadcrumbs={{this.breadcrumbs}}
            @capabilities={{this.capabilities}}
          />
      `,
        { owner: this.engine }
      );
  });

  test('it renders version history and shows icons for deleted, destroyed and current', async function (assert) {
    assert.expect(7); // 4 linked blocks, 2 destroyed, 1 deleted.

    await this.renderComponent();

    for (const version in this.metadata.versions) {
      const data = this.metadata.versions[version];
      assert.dom(PAGE.versions.linkedBlock(version)).exists(`renders the linked blocks for each version`);

      if (data.destroyed) {
        assert
          .dom(`${PAGE.versions.icon(version)} [data-test-icon="x-square-fill"]`)
          .hasStyle({ color: 'rgb(229, 34, 40)' });
      }
      if (isDeleted(data.deletion_time)) {
        assert
          .dom(`${PAGE.versions.icon(version)} [data-test-icon="x-square-fill"]`)
          .hasStyle({ color: 'rgb(101, 106, 118)' });
      }
    }
  });

  test('it gives the option to create a new version from a secret from the popup menu', async function (assert) {
    assert.expect(1);

    await this.renderComponent();
    // because the popup menu is nested in a linked block we must combine the two selectors
    const popupSelector = `${PAGE.versions.linkedBlock(1)} ${PAGE.popup}`;
    await click(popupSelector);
    assert
      .dom('[data-test-create-new-version-from="1"]')
      .exists('Shows the option to create a new version from that secret.');
  });
});
