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
  },
  {
    month: '8/22',
    timestamp: '2022-08-01T00:00:00-07:00',
    clients: 6440,
    entity_clients: 1471,
    non_entity_clients: 4389,
    secret_syncs: 4207,
  },
  {
    month: '9/22',
    timestamp: '2022-09-01T00:00:00-07:00',
    clients: 9583,
    entity_clients: 149,
    non_entity_clients: 20,
    secret_syncs: 5802,
  },
];

module('Integration | Component | clients/charts/vertical-bar-basic', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.data = EXAMPLE;
  });

  test('it renders when some months have no data', async function (assert) {
    assert.expect(10);
    await render(
      hbs`<Clients::Charts::VerticalBarBasic @data={{this.data}} @dataKey="secret_syncs" @chartTitle="My chart"/>`
    );
    assert.dom('[data-test-chart="My chart"]').exists('renders chart container');
    assert.dom('[data-test-vertical-bar]').exists({ count: 3 }, 'renders 3 vertical bars');

    // Tooltips
    await triggerEvent('[data-test-interactive-area="9/22"]', 'mouseover');
    assert.dom('[data-test-tooltip]').exists({ count: 1 }, 'renders tooltip on mouseover');
    assert.dom('[data-test-tooltip-count]').hasText('5,802 secret syncs', 'tooltip has exact count');
    assert.dom('[data-test-tooltip-month]').hasText('September 2022', 'tooltip has humanized month and year');
    await triggerEvent('[data-test-interactive-area="9/22"]', 'mouseout');
    assert.dom('[data-test-tooltip]').doesNotExist('removes tooltip on mouseout');
    await triggerEvent('[data-test-interactive-area="7/22"]', 'mouseover');
    assert
      .dom('[data-test-tooltip-count]')
      .hasText('No data', 'renders tooltip with no data message when no data is available');
    // Axis
    assert.dom('[data-test-x-axis]').hasText('7/22 8/22 9/22', 'renders x-axis labels');
    assert.dom('[data-test-y-axis]').hasText('0 2k 4k', 'renders y-axis labels');
    // Table
    assert.dom('[data-test-underlying-data]').doesNotExist('does not render underlying data by default');
  });

  // 0 is different than null (no data)
  test('it renders when all months have 0 clients', async function (assert) {
    assert.expect(9);

    this.data = [
      {
        month: '6/22',
        timestamp: '2022-06-01T00:00:00-07:00',
        clients: 0,
        entity_clients: 0,
        non_entity_clients: 0,
        secret_syncs: 0,
      },
      {
        month: '7/22',
        timestamp: '2022-07-01T00:00:00-07:00',
        clients: 0,
        entity_clients: 0,
        non_entity_clients: 0,
        secret_syncs: 0,
      },
    ];
    await render(
      hbs`<Clients::Charts::VerticalBarBasic @data={{this.data}} @dataKey="secret_syncs" @chartTitle="My chart"/>`
    );

    assert.dom('[data-test-chart="My chart"]').exists('renders chart container');
    assert.dom('[data-test-vertical-bar]').exists({ count: 2 }, 'renders 2 vertical bars');
    assert.dom('[data-test-vertical-bar]').hasAttribute('height', '0', 'rectangles have 0 height');
    // Tooltips
    await triggerEvent('[data-test-interactive-area="6/22"]', 'mouseover');
    assert.dom('[data-test-tooltip]').exists({ count: 1 }, 'renders tooltip on mouseover');
    assert.dom('[data-test-tooltip-count]').hasText('0 secret syncs', 'tooltip has exact count');
    assert.dom('[data-test-tooltip-month]').hasText('June 2022', 'tooltip has humanized month and year');
    await triggerEvent('[data-test-interactive-area="6/22"]', 'mouseout');
    assert.dom('[data-test-tooltip]').doesNotExist('removes tooltip on mouseout');
    // Axis
    assert.dom('[data-test-x-axis]').hasText('6/22 7/22', 'renders x-axis labels');
    assert.dom('[data-test-y-axis]').hasText('0 1 2 3 4', 'renders y-axis labels');
  });

  test('it renders underlying data', async function (assert) {
    assert.expect(3);
    await render(
      hbs`<Clients::Charts::VerticalBarBasic @data={{this.data}} @dataKey="secret_syncs" @showTable={{true}} @chartTitle="My chart"/>`
    );
    assert.dom('[data-test-chart="My chart"]').exists('renders chart container');
    assert.dom('[data-test-underlying-data]').exists('renders underlying data when showTable=true');
    assert
      .dom('[data-test-underlying-data] thead')
      .hasText('Month Secret syncs Count', 'renders correct table headers');
  });
});
