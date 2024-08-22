/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, triggerEvent } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CHARTS } from 'vault/tests/helpers/clients/client-count-selectors';
import { assertBarChart } from 'vault/tests/helpers/clients/client-count-helpers';

module('Integration | Component | clients/charts/vertical-bar-grouped', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.legend = [
      { key: 'clients', label: 'Total clients' },
      { key: 'foo', label: 'Foo' },
    ];
    this.data = [
      {
        timestamp: '2018-04-03T14:15:30',
        clients: 14,
        foo: 4,
        month: '4/18',
      },
      {
        timestamp: '2018-05-03T14:15:30',
        clients: 18,
        foo: 8,
        month: '5/18',
      },
      {
        timestamp: '2018-06-03T14:15:30',
        clients: 114,
        foo: 14,
        month: '6/18',
      },
      {
        timestamp: '2018-07-03T14:15:30',
        clients: 110,
        foo: 10,
        month: '7/18',
      },
    ];
    this.renderComponent = async () => {
      await render(
        hbs`<div class="has-top-padding-xxl">
        <Clients::Charts::VerticalBarGrouped @data={{this.data}} @legend={{this.legend}} @upgradeData={{this.upgradeData}} />
        </div>`
      );
    };
  });

  test('it renders empty state when no data', async function (assert) {
    this.data = [];
    await this.renderComponent();
    assert.dom(CHARTS.chart('grouped vertical bar chart')).doesNotExist();
    assert.dom(GENERAL.emptyStateSubtitle).hasText('No data to display');
  });

  test('it renders chart with data as grouped bars', async function (assert) {
    await this.renderComponent();
    assert.dom(CHARTS.chart('grouped vertical bar chart')).exists();
    const barCount = this.data.length * this.legend.length;
    // bars are what we expect
    assert.dom(CHARTS.verticalBar).exists({ count: barCount });
    assert.dom(`.custom-bar-clients`).exists({ count: 4 }, 'clients bars have correct class');
    assert.dom(`.custom-bar-foo`).exists({ count: 4 }, 'foo bars have correct class');
    assertBarChart(assert, 'grouped vertical bar chart', this.data, true);
  });

  test('it renders chart with tooltips when some is missing', async function (assert) {
    assert.expect(13);
    this.data = [
      {
        timestamp: '2018-04-03T14:15:30',
        month: '4/18',
        expectedTooltip: 'April 2018 No data',
      },
      {
        timestamp: '2018-05-03T14:15:30',
        month: '5/18',
        clients: 0,
        foo: 0,
      },
      {
        timestamp: '2018-06-03T14:15:30',
        month: '6/18',
        clients: 14,
        foo: 4,
        expectedTooltip: 'June 2018 14 Total clients 4 Foo',
      },
    ];
    await this.renderComponent();
    assert.dom(CHARTS.chart('grouped vertical bar chart')).exists();
    const barCount = this.data.length * this.legend.length;
    assert.dom(CHARTS.verticalBar).exists({ count: barCount });
    assertBarChart(assert, 'grouped vertical bar chart', this.data, true);

    // TOOLTIPS - NO DATA
    await triggerEvent(CHARTS.hover(this.data[0].timestamp), 'mouseover');
    assert.dom(CHARTS.tooltip).isVisible(`renders tooltip on mouseover`);
    assert
      .dom(CHARTS.tooltip)
      .hasText(this.data[0].expectedTooltip, 'renders formatted timestamp with no data message');
    await triggerEvent(CHARTS.hover(this.data[2].timestamp), 'mouseout');
    assert.dom(CHARTS.tooltip).doesNotExist('removes tooltip on mouseout');

    // TOOLTIPS - WITH DATA
    await triggerEvent(CHARTS.hover(this.data[2].timestamp), 'mouseover');
    assert.dom(CHARTS.tooltip).isVisible(`renders tooltip on mouseover`);
    assert.dom(CHARTS.tooltip).hasText(this.data[2].expectedTooltip, 'renders formatted timestamp with data');
    await triggerEvent(CHARTS.hover(this.data[2].timestamp), 'mouseout');
    assert.dom(CHARTS.tooltip).doesNotExist('removes tooltip on mouseout');
  });

  test('it renders upgrade data', async function (assert) {
    this.upgradeData = [
      {
        version: '1.10.1',
        previousVersion: '1.9.2',
        timestampInstalled: '2018-05-03T14:15:30',
      },
    ];
    await this.renderComponent();
    assert.dom(CHARTS.chart('grouped vertical bar chart')).exists();
    const barCount = this.data.length * this.legend.length;
    // bars are what we expect
    assert.dom(CHARTS.verticalBar).exists({ count: barCount });
    assert.dom(`.custom-bar-clients`).exists({ count: 4 }, 'clients bars have correct class');
    assert.dom(`.custom-bar-foo`).exists({ count: 4 }, 'foo bars have correct class');
    assertBarChart(assert, 'grouped vertical bar chart', this.data, true);

    // TOOLTIP
    await triggerEvent(CHARTS.hover('2018-05-03T14:15:30'), 'mouseover');
    assert.dom(CHARTS.tooltip).isVisible(`renders tooltip on mouseover`);
    assert
      .dom(CHARTS.tooltip)
      .hasText(
        'May 2018 18 Total clients 8 Foo Vault was upgraded from 1.9.2 to 1.10.1',
        'renders formatted timestamp with data'
      );
    await triggerEvent(CHARTS.hover('2018-05-03T14:15:30'), 'mouseout');
    assert.dom(CHARTS.tooltip).doesNotExist('removes tooltip on mouseout');
  });

  test('it updates axis when dataset updates', async function (assert) {
    const datasets = {
      small: [
        {
          timestamp: '2020-04-01',
          bar: 4,
          month: '4/20',
        },
        {
          timestamp: '2020-05-01',
          bar: 8,
          month: '5/20',
        },
        {
          timestamp: '2020-06-01',
          bar: 1,
        },
        {
          timestamp: '2020-07-01',
          bar: 10,
        },
      ],
      large: [
        {
          timestamp: '2020-08-01',
          bar: 4586,
          month: '8/20',
        },
        {
          timestamp: '2020-09-01',
          bar: 8928,
          month: '9/20',
        },
        {
          timestamp: '2020-10-01',
          bar: 11948,
          month: '10/20',
        },
        {
          timestamp: '2020-11-01',
          bar: 16943,
          month: '11/20',
        },
      ],
      broken: [
        {
          timestamp: '2020-01-01',
          bar: null,
          month: '1/20',
        },
        {
          timestamp: '2020-02-01',
          bar: 0,
          month: '2/20',
        },
        {
          timestamp: '2020-03-01',
          bar: 22,
          month: '3/20',
        },
        {
          timestamp: '2020-04-01',
          bar: null,
          month: '4/20',
        },
        {
          timestamp: '2020-05-01',
          bar: 70,
          month: '5/20',
        },
        {
          timestamp: '2020-06-01',
          bar: 50,
          month: '6/20',
        },
      ],
    };
    this.legend = [{ key: 'bar', label: 'Some thing' }];
    this.set('data', datasets.small);
    await this.renderComponent();
    assert.dom('[data-test-y-axis]').hasText('0 2 4 6 8 10', 'y-axis renders correctly for small values');
    assert
      .dom('[data-test-x-axis]')
      .hasText('4/20 5/20 6/20 7/20', 'x-axis renders correctly for small values');

    // Update to large dataset
    this.set('data', datasets.large);
    assert.dom('[data-test-y-axis]').hasText('0 5k 10k 15k', 'y-axis renders correctly for new large values');
    assert
      .dom('[data-test-x-axis]')
      .hasText('8/20 9/20 10/20 11/20', 'x-axis renders correctly for small values');

    // Update to broken dataset
    this.set('data', datasets.broken);
    assert.dom('[data-test-y-axis]').hasText('0 20 40 60', 'y-axis renders correctly for new broken values');
    assert
      .dom('[data-test-x-axis]')
      .hasText('1/20 2/20 3/20 4/20 5/20 6/20', 'x-axis renders correctly for small values');
  });
});
