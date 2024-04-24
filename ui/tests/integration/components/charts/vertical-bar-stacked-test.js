/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { findAll, render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { CHARTS } from 'vault/tests/helpers/clients/client-count-selectors';

const EXAMPLE = [
  {
    timestamp: '2022-09-01T00:00:00',
    total: null,
    fuji_apples: null,
    gala_apples: null,
    red_delicious: null,
  },
  {
    timestamp: '2022-10-01T00:00:00',
    total: 6440,
    fuji_apples: 1471,
    gala_apples: 4389,
    red_delicious: 4207,
  },
  {
    timestamp: '2022-11-01T00:00:00',
    total: 9583,
    fuji_apples: 149,
    gala_apples: 20,
    red_delicious: 5802,
  },
];

module('Integration | Component | clients/charts/vertical-bar-stacked', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.data = EXAMPLE;
    this.legend = [
      { key: 'fuji_apples', label: 'Fuji counts' },
      { key: 'gala_apples', label: 'Gala counts' },
    ];
  });

  test('it renders when some months have no data', async function (assert) {
    assert.expect(10);
    await render(
      hbs`<Clients::Charts::VerticalBarStacked @data={{this.data}} @chartLegend={{this.legend}} @chartTitle="My chart"/>`
    );

    assert.dom(CHARTS.chart('My chart')).exists('renders chart container');

    const visibleBars = findAll(CHARTS.verticalBar).filter((e) => e.getAttribute('height') !== '0');
    const count = this.data.filter((d) => d.total !== null).length * 2;
    assert.strictEqual(visibleBars.length, count, `renders ${count} vertical bars`);

    // Tooltips
    await triggerEvent(CHARTS.hover('2022-09-01T00:00:00'), 'mouseover');
    assert.dom(CHARTS.tooltip).isVisible('renders tooltip on mouseover');
    assert
      .dom(CHARTS.tooltip)
      .hasText('September 2022 No data', 'renders formatted timestamp with no data message');
    await triggerEvent(CHARTS.hover('2022-09-01T00:00:00'), 'mouseout');
    assert.dom(CHARTS.tooltip).doesNotExist('removes tooltip on mouseout');

    await triggerEvent(CHARTS.hover('2022-10-01T00:00:00'), 'mouseover');
    assert
      .dom(CHARTS.tooltip)
      .hasText('October 2022 1,471 Fuji counts 4,389 Gala counts', 'October tooltip has exact count');
    await triggerEvent(CHARTS.hover('2022-10-01T00:00:00'), 'mouseout');

    await triggerEvent(CHARTS.hover('2022-11-01T00:00:00'), 'mouseover');
    assert
      .dom(CHARTS.tooltip)
      .hasText('November 2022 149 Fuji counts 20 Gala counts', 'November tooltip has exact count');
    await triggerEvent(CHARTS.hover('2022-11-01T00:00:00'), 'mouseout');

    // Axis
    assert.dom(CHARTS.xAxis).hasText('9/22 10/22 11/22', 'renders x-axis labels');
    assert.dom(CHARTS.yAxis).hasText('0 2k 4k', 'renders y-axis labels');
    // Table
    assert.dom(CHARTS.table).doesNotExist('does not render underlying data by default');
  });

  // 0 is different than null (no data)
  test('it renders when all months have 0 clients', async function (assert) {
    assert.expect(14);

    this.data = [
      {
        month: '10/22',
        timestamp: '2022-10-01T00:00:00',
        total: 40,
        fuji_apples: 0,
        gala_apples: 0,
        red_delicious: 40,
      },
      {
        month: '11/22',
        timestamp: '2022-11-01T00:00:00',
        total: 180,
        fuji_apples: 0,
        gala_apples: 0,
        red_delicious: 180,
      },
    ];
    await render(
      hbs`<Clients::Charts::VerticalBarStacked @data={{this.data}} @chartLegend={{this.legend}} @chartTitle="My chart"/>`
    );

    assert.dom(CHARTS.chart('My chart')).exists('renders chart container');
    findAll(CHARTS.verticalBar).forEach((b, idx) =>
      assert.dom(b).isNotVisible(`bar: ${idx} does not render`)
    );
    findAll(CHARTS.verticalBar).forEach((b, idx) =>
      assert.dom(b).hasAttribute('height', '0', `rectangle: ${idx} have 0 height`)
    );

    // Tooltips
    await triggerEvent(CHARTS.hover('2022-10-01T00:00:00'), 'mouseover');
    assert.dom(CHARTS.tooltip).isVisible('renders tooltip on mouseover');
    assert.dom(CHARTS.tooltip).hasText('October 2022 0 Fuji counts 0 Gala counts', 'tooltip has 0 counts');
    await triggerEvent(CHARTS.hover('2022-10-01T00:00:00'), 'mouseout');
    assert.dom(CHARTS.tooltip).isNotVisible('removes tooltip on mouseout');

    // Axis
    assert.dom(CHARTS.xAxis).hasText('10/22 11/22', 'renders x-axis labels');
    assert.dom(CHARTS.yAxis).hasText('0 1 2 3 4', 'renders y-axis labels');
  });

  test('it renders underlying data', async function (assert) {
    assert.expect(3);
    await render(
      hbs`<Clients::Charts::VerticalBarStacked @data={{this.data}} @chartLegend={{this.legend}} @showTable={{true}} @chartTitle="My chart"/>`
    );
    assert.dom(CHARTS.chart('My chart')).exists('renders chart container');
    assert.dom(CHARTS.table).exists('renders underlying data when showTable=true');
    assert
      .dom(`${CHARTS.table} thead`)
      .hasText('Timestamp Fuji apples Gala apples', 'renders correct table headers');
  });
});
