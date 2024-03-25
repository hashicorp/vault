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
    this.isSecretsSyncActivated = false;
    this.counts = {};

    this.renderComponent = async () =>
      await render(
        hbs`<Clients::UsageStats @totalUsageCounts={{this.counts}} @isSecretsSyncActivated={{this.isSecretsSyncActivated}} />`
      );
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('[data-test-stat-text]').exists({ count: 3 }, 'Renders 3 Stat texts even with no data passed');
    assert.dom('[data-test-stat-text="total-clients"]').exists('Total clients exists');
    assert.dom('[data-test-stat-text="total-clients"] .stat-value').hasText('-', 'renders dash when no data');
    assert.dom('[data-test-stat-text="entity-clients"]').exists('Entity clients exists');
    assert
      .dom('[data-test-stat-text="entity-clients"] .stat-value')
      .hasText('-', 'renders dash when no data');
    assert.dom('[data-test-stat-text="non-entity-clients"]').exists('Non entity clients exists');
    assert
      .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
      .hasText('-', 'renders dash when no data');
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

    assert.dom('[data-test-stat-text]').exists({ count: 3 }, 'Renders 3 Stat texts');
    assert
      .dom('[data-test-stat-text="total-clients"] .stat-value')
      .hasText('17', 'Total clients shows passed value');
    assert
      .dom('[data-test-stat-text="entity-clients"] .stat-value')
      .hasText('7', 'entity clients shows passed value');
    assert
      .dom('[data-test-stat-text="non-entity-clients"] .stat-value')
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
      this.isSecretsSyncActivated = true;

      await this.renderComponent();

      assert.dom('[data-test-stat-text]').exists({ count: 4 }, 'Renders 4 Stat texts');
      assert
        .dom('[data-test-stat-text="secret-syncs"] .stat-value')
        .hasText('5', 'secrets sync clients shows passed value');
    });

    test('with secrets sync NOT activated', async function (assert) {
      this.isSecretsSyncActivated = false;

      await this.renderComponent();

      assert.dom('[data-test-stat-text]').exists({ count: 3 }, 'Renders 3 Stat texts');
      assert.dom('[data-test-stat-text="secret-syncs"] .stat-value').doesNotExist();
    });
  });
});
