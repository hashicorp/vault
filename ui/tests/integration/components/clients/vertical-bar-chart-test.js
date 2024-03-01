/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, findAll, find } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | clients/vertical-bar-chart', function (hooks) {
  setupRenderingTest(hooks);
  hooks.beforeEach(function () {
    this.set('chartLegend', [
      { label: 'entity clients', key: 'entity_clients' },
      { label: 'non-entity clients', key: 'non_entity_clients' },
    ]);
  });

  test('it renders chart and tooltip for total clients', async function (assert) {
    const barChartData = [
      { month: 'january', clients: 141, entity_clients: 91, non_entity_clients: 50, new_clients: 5 },
      { month: 'february', clients: 251, entity_clients: 101, non_entity_clients: 150, new_clients: 5 },
    ];
    this.set('barChartData', barChartData);

    await render(hbs`
    <div class="chart-container-wide">
      <Clients::VerticalBarChart
        @dataset={{this.barChartData}}
        @chartLegend={{this.chartLegend}}
      />
    </div>
    `);
    assert.dom('[data-test-vertical-bar-chart]').exists('renders chart');
    assert
      .dom('[data-test-vertical-chart="data-bar"]')
      .exists({ count: barChartData.length * 2 }, 'renders correct number of bars'); // multiply length by 2 because bars are stacked

    assert.dom(find('[data-test-vertical-chart="y-axis-labels"] text')).hasText('0', `y-axis starts at 0`);
    findAll('[data-test-vertical-chart="x-axis-labels"] text').forEach((e, i) => {
      assert.dom(e).hasText(`${barChartData[i].month}`, `renders x-axis label: ${barChartData[i].month}`);
    });

    // FLAKY after adding a11y testing, skip for now
    // const tooltipHoverBars = findAll('[data-test-vertical-bar-chart] rect.tooltip-rect');
    // for (const [i, bar] of tooltipHoverBars.entries()) {
    //   await triggerEvent(bar, 'mouseover');
    //   const tooltip = document.querySelector('.ember-modal-dialog');
    //   assert
    //     .dom(tooltip)
    //     .includesText(
    //       `${barChartData[i].clients} total clients ${barChartData[i].entity_clients} entity clients ${barChartData[i].non_entity_clients} non-entity clients`,
    //       'tooltip text is correct'
    //     );
    // }
  });

  test('it renders chart and tooltip for new clients', async function (assert) {
    const barChartData = [
      { month: 'january', entity_clients: 91, non_entity_clients: 50, clients: 0 },
      { month: 'february', entity_clients: 101, non_entity_clients: 150, clients: 110 },
    ];
    this.set('barChartData', barChartData);

    await render(hbs`
    <div class="chart-container-wide">
      <Clients::VerticalBarChart
        @dataset={{this.barChartData}}
        @chartLegend={{this.chartLegend}}
      />
    </div>
    `);

    assert.dom('[data-test-vertical-bar-chart]').exists('renders chart');
    assert
      .dom('[data-test-vertical-chart="data-bar"]')
      .exists({ count: barChartData.length * 2 }, 'renders correct number of bars'); // multiply length by 2 because bars are stacked

    assert.dom(find('[data-test-vertical-chart="y-axis-labels"] text')).hasText('0', `y-axis starts at 0`);
    findAll('[data-test-vertical-chart="x-axis-labels"] text').forEach((e, i) => {
      assert.dom(e).hasText(`${barChartData[i].month}`, `renders x-axis label: ${barChartData[i].month}`);
    });

    // FLAKY after adding a11y testing, skip for now
    // const tooltipHoverBars = findAll('[data-test-vertical-bar-chart] rect.tooltip-rect');
    // for (const [i, bar] of tooltipHoverBars.entries()) {
    //   await triggerEvent(bar, 'mouseover');
    //   const tooltip = document.querySelector('.ember-modal-dialog');
    //   assert
    //     .dom(tooltip)
    //     .includesText(
    //       `${barChartData[i].clients} new clients ${barChartData[i].entity_clients} entity clients ${barChartData[i].non_entity_clients} non-entity clients`,
    //       'tooltip text is correct'
    //     );
    // }
  });

  test('it renders empty state when no dataset', async function (assert) {
    await render(hbs`
    <div class="chart-container-wide">
    <Clients::VerticalBarChart @noDataMessage="this is a custom message to explain why you're not seeing a vertical bar chart"/>
    </div>
    `);

    assert.dom('[data-test-component="empty-state"]').exists('renders empty state when no data');
    assert
      .dom('[data-test-empty-state-subtext]')
      .hasText(
        `this is a custom message to explain why you're not seeing a vertical bar chart`,
        'custom message renders'
      );
  });
});
