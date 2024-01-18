/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const EXAMPLE = [
  {
    month: '7/22',
    timestamp: '2022-07-01T00:00:00-07:00',
    clients: null,
    entity_clients: null,
    non_entity_clients: null,
    secret_syncs: null,
    namespaces: [],
    new_clients: {},
    namespaces_by_key: {},
  },
  {
    month: '8/22',
    timestamp: '2022-08-01T00:00:00-07:00',
    clients: 6440,
    entity_clients: 1471,
    non_entity_clients: 4389,
    secret_syncs: 4207,
    namespaces: [],
    new_clients: {},
    namespaces_by_key: {},
  },
  {
    month: '9/22',
    timestamp: '2022-09-01T00:00:00-07:00',
    clients: 9583,
    entity_clients: 149,
    non_entity_clients: 20,
    secret_syncs: 5802,
    namespaces: [],
    new_clients: {},
    namespaces_by_key: {},
  },
];
module('Integration | Component | clients/sync-bar-chart', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.data = EXAMPLE;
  });

  test('it renders when some months have no data', async function (assert) {
    await render(hbs`<Clients::SyncBarChart @data={{this.data}} @dataKey="secret_syncs" />`);
    assert.dom('[data-test-sync-bar-chart]').exists('renders chart container');
    assert.dom('[data-test-vertical-bar]').exists({ count: 3 }, 'renders 3 vertical bars');
    // Tooltips
    await triggerEvent('[data-test-interactive-area="9/22"]', 'mouseover');
    assert.dom('[data-test-tooltip]').exists({ count: 1 }, 'renders tooltip on mouseover');
    assert.dom('[data-test-tooltip-count]').hasText('5,802 secret syncs', 'tooltip has exact count');
    assert.dom('[data-test-tooltip-month]').hasText('Sep 2022', 'tooltip has humanized month and year');
    await triggerEvent('[data-test-interactive-area="9/22"]', 'mouseout');
    assert.dom('[data-test-tooltip]').doesNotExist('removes tooltip on mouseout');
    await triggerEvent('[data-test-interactive-area="7/22"]', 'mouseover');
    assert
      .dom('[data-test-tooltip-count]')
      .hasText('No data', 'renders tooltip with no data message when no data is available');
    // Axis
    assert.dom('[data-test-x-axis]').hasText('7/22 8/22 9/22', 'renders x-axis labels');
    assert.dom('[data-test-y-axis]').hasText('0 2k 4k 6k', 'renders y-axis labels');
    // Table
    assert.dom('[data-test-underlying-data]').doesNotExist('does not render underlying data by default');
  });

  test('it renders underlying data', async function (assert) {
    await render(
      hbs`<Clients::SyncBarChart @data={{this.data}} @dataKey="secret_syncs" @showTable={{true}} />`
    );
    assert.dom('[data-test-sync-bar-chart]').exists('renders chart container');
    assert.dom('[data-test-underlying-data]').exists('renders underlying data when showTable=true');
    assert
      .dom('[data-test-underlying-data] thead')
      .hasText('Month Count of secret syncs', 'renders correct table headers');
  });
});
