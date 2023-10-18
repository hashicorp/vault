/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';

module('Integration | Component | sync header', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.16.0+ent';
  });

  test('it should render promotional enterprise badge for community version', async function (assert) {
    this.version.version = '1.16.0';
    await render(
      hbs`
     <SyncHeader @title="Secrets sync"/>
    `,
      { owner: this.engine }
    );

    assert.dom(PAGE.title).hasText('Secrets sync Enterprise feature');
    assert.dom(PAGE.breadcrumbs).doesNotExist('does not render breadcrumbs');
  });

  test('it should not render enterprise badge for enterprise versions', async function (assert) {
    this.version.version = '1.16.0+ent';
    await render(
      hbs`
        <SyncHeader @title="Secrets sync"/>
    `,
      { owner: this.engine }
    );

    assert.dom(PAGE.title).hasText('Secrets sync');
    assert.dom(PAGE.headerContainer).hasTextContaining('Secrets sync/', 'renders default breadcrumb');
  });

  test('it renders breadcrumbs', async function (assert) {
    this.breadcrumbs = [{ label: 'Destinations', route: 'destinations' }];

    await render(
      hbs`
         <SyncHeader @title="Secrets sync" @breadcrumbs={{this.breadcrumbs}}/>
    `,
      { owner: this.engine }
    );
    assert.dom(PAGE.headerContainer).hasTextContaining('Destinations', 'renders breadcrumbs');
  });
});
