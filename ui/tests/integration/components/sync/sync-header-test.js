/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import hbs from 'htmlbars-inline-precompile';
import { render } from '@ember/test-helpers';

module('Integration | Component | sync | SyncHeader', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.16.0+ent';
    this.title = 'Secrets sync';
    this.renderComponent = () => {
      return render(hbs`<SyncHeader @title={{this.title}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
    };
  });

  test('it should render default breadcrumb', async function (assert) {
    await this.renderComponent();
    assert.dom('[data-test-breadcrumbs]').exists({ count: 1 }, 'Correct number of breadcrumbs render');
    assert.dom('[data-test-crumb]').includesText('Secrets sync', 'renders default breadcrumb');
  });

  test('it should render breadcrumbs', async function (assert) {
    this.breadcrumbs = [{ label: 'Destinations', route: 'destinations' }];
    await this.renderComponent();
    assert.dom('[data-test-crumb]').includesText('Destinations', 'renders breadcrumb');
  });

  test('it should title for enterprise version', async function (assert) {
    await this.renderComponent();
    assert.dom('[data-test-page-title]').hasText('Secrets sync');
  });

  test('it should render title and promotional enterprise badge for community version', async function (assert) {
    this.version.version = '1.16.0';
    await this.renderComponent();
    assert.dom('[data-test-page-title]').hasText('Secrets sync Enterprise feature');
  });
});
