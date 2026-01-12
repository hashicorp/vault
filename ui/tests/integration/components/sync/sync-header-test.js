/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const { breadcrumb } = PAGE;

module('Integration | Component | sync | SyncHeader', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.flags = this.owner.lookup('service:flags');
    this.title = 'Secrets Sync';
    this.renderComponent = () => {
      return render(hbs`<SyncHeader @title={{this.title}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
    };
  });

  test('it should render breadcrumbs', async function (assert) {
    this.breadcrumbs = [{ label: 'Destinations', route: 'destinations' }];
    await this.renderComponent();
    assert.dom(breadcrumb).includesText('Destinations', 'renders breadcrumb');
  });

  module('ent', function (hooks) {
    hooks.beforeEach(async function () {
      this.version.type = 'enterprise';
    });

    test('it should render title if license has secrets sync feature', async function (assert) {
      this.version.features = ['Secrets Sync'];
      await this.renderComponent();

      assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Secrets Sync');
    });
  });

  module('managed', function (hooks) {
    hooks.beforeEach(function () {
      this.version.type = 'enterprise';
      this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    });

    test('it should render title and plus badge', async function (assert) {
      await this.renderComponent();
      assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Secrets Sync');
      assert.dom(GENERAL.badge('Plus feature')).hasText('Plus feature', 'Plus feature badge renders');
    });
  });
});
