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

const { breadcrumb, title } = PAGE;

module('Integration | Component | sync | SyncHeader', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');

  hooks.beforeEach(function () {
    this.version = this.owner.lookup('service:version');
    this.version.type = 'enterprise';
    this.title = 'Secrets Sync';
    this.renderComponent = () => {
      return render(hbs`<SyncHeader @title={{this.title}} @breadcrumbs={{this.breadcrumbs}} />`, {
        owner: this.engine,
      });
    };
  });

  test('it should render default breadcrumb', async function (assert) {
    await this.renderComponent();
    assert.dom(breadcrumb).exists({ count: 1 }, 'Correct number of breadcrumbs render');
    assert.dom(breadcrumb).includesText('Secrets Sync', 'renders default breadcrumb');
  });

  test('it should render breadcrumbs', async function (assert) {
    this.breadcrumbs = [{ label: 'Destinations', route: 'destinations' }];
    await this.renderComponent();
    assert.dom(breadcrumb).includesText('Destinations', 'renders breadcrumb');
  });

  test('it should just render title for enterprise version', async function (assert) {
    await this.renderComponent();
    assert.dom(title).hasText('Secrets Sync Beta');
  });

  test('it should render title and promotional enterprise badge for community version', async function (assert) {
    this.version.type = 'community';
    await this.renderComponent();
    assert.dom(title).hasText('Secrets Sync Enterprise feature');
  });

  test('it should yield actions block', async function (assert) {
    await render(
      hbs`
      <SyncHeader @title={{this.title}} @breadcrumbs={{this.breadcrumbs}}>
        <:actions>
          <span data-test-action-block>Test</span>
        </:actions>
      </SyncHeader>
    `,
      { owner: this.engine }
    );

    assert.dom('[data-test-action-block]').exists('Component yields block for actions');
  });
});
