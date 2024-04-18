/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | clients/usage-stats', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.showSecretSyncs = false;
    this.counts = {};

    this.renderComponent = async () =>
      await render(
        hbs`<Clients::UsageStats @totalUsageCounts={{this.counts}} @showSecretSyncs={{this.showSecretSyncs}} />`
      );
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('[data-test-stat-text]').exists({ count: 4 }, 'Renders 4 Stat texts even with no data passed');
    assert.dom('[data-test-stat-text="Total clients"]').exists('Total clients exists');
    assert.dom('[data-test-stat-text="Total clients"] .stat-value').hasText('-', 'renders dash when no data');
    assert.dom('[data-test-stat-text="Entity"]').exists('Entity exists');
    assert.dom('[data-test-stat-text="Entity"] .stat-value').hasText('-', 'renders dash when no data');
    assert.dom('[data-test-stat-text="Non-entity"]').exists('Non entity clients exists');
    assert.dom('[data-test-stat-text="Non-entity"] .stat-value').hasText('-', 'renders dash when no data');
    assert
      .dom('a')
      .hasAttribute('href', 'https://developer.hashicorp.com/vault/tutorials/monitoring/usage-metrics');
  });

  test('it renders with token data', async function (assert) {
    this.counts = {
      clients: 17,
      entity_clients: 7,
      non_entity_clients: 10,
    };

    await this.renderComponent();

    assert.dom('[data-test-stat-text]').exists({ count: 4 }, 'Renders 4 Stat texts');
    assert
      .dom('[data-test-stat-text="Total clients"] .stat-value')
      .hasText('17', 'Total clients shows passed value');
    assert
      .dom('[data-test-stat-text="Entity"] .stat-value')
      .hasText('7', 'entity clients shows passed value');
    assert
      .dom('[data-test-stat-text="Non-entity"] .stat-value')
      .hasText('10', 'non entity clients shows passed value');
  });

  module('it renders with full totals data', function (hooks) {
    hooks.beforeEach(function () {
      this.counts = {
        clients: 22,
        entity_clients: 7,
        non_entity_clients: 10,
        secret_syncs: 5,
      };
    });

    test('with secrets sync activated', async function (assert) {
      this.showSecretSyncs = true;

      await this.renderComponent();

      assert.dom('[data-test-stat-text]').exists({ count: 5 }, 'Renders 5 Stat texts');
      assert
        .dom('[data-test-stat-text="Secret sync"] .stat-value')
        .hasText('5', 'secrets sync clients shows passed value');
    });

    test('with secrets sync NOT activated', async function (assert) {
      this.showSecretSyncs = false;

      await this.renderComponent();

      assert.dom('[data-test-stat-text]').exists({ count: 4 }, 'Renders 4 Stat texts');
      assert.dom('[data-test-stat-text="Secret sync"] .stat-value').doesNotExist();
    });
  });
});
