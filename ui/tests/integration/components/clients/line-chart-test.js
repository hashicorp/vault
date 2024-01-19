/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { render, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { format, formatRFC3339, parseISO, subMonths } from 'date-fns';
import timestamp from 'core/utils/timestamp';

module('Integration | Component | clients/line-chart', function (hooks) {
  setupRenderingTest(hooks);
  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => new Date('2018-04-03T14:15:30'));
  });
  hooks.beforeEach(function () {
    this.set('xKey', 'foo');
    this.set('yKey', 'bar');
    this.set('dataset', [
      {
        foo: '2017-12-03T14:15:30',
        bar: 4,
      },
      {
        foo: '2018-01-03T14:15:30',
        bar: 8,
      },
      {
        foo: '2018-02-03T14:15:30',
        bar: 14,
      },
      {
        foo: '2018-03-03T14:15:30',
        bar: 10,
      },
    ]);
  });
  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it renders', async function (assert) {
    await render(hbs`
    <div class="chart-container-wide">
      <Clients::LineChart @dataset={{this.dataset}} @xKey={{this.xKey}} @yKey={{this.yKey}} />
    </div>
    `);

    assert.dom('[data-test-line-chart]').exists('Chart is rendered');
    assert
      .dom('[data-test-line-chart="plot-point"]')
      .exists({ count: this.dataset.length }, `renders ${this.dataset.length} plot points`);

    findAll('[data-test-x-axis] text').forEach((e, i) => {
      // For some reason the first x-axis label is missing
      const date = parseISO(this.dataset[i + 1][this.xKey]);
      const monthLabel = format(date, 'M/yy');
      assert.dom(e).hasText(monthLabel, `renders x-axis label: ${monthLabel}`);
    });
    assert.dom('[data-test-y-axis] text').hasText('0', `y-axis starts at 0`);
  });

  test('it renders upgrade data', async function (assert) {
    const now = timestamp.now();
    this.set('dataset', [
      {
        foo: formatRFC3339(subMonths(now, 4)),
        bar: 4,
      },
      {
        foo: formatRFC3339(subMonths(now, 3)),
        bar: 8,
      },
      {
        foo: formatRFC3339(subMonths(now, 2)),
        bar: 14,
      },
      {
        foo: formatRFC3339(subMonths(now, 1)),
        bar: 10,
      },
    ]);
    this.set('upgradeData', [
      {
        version: '1.10.1',
        previousVersion: '1.9.2',
        timestampInstalled: formatRFC3339(subMonths(now, 2)),
      },
    ]);
    await render(hbs`
    <div class="chart-container-wide">
      <Clients::LineChart
        @dataset={{this.dataset}}
        @upgradeData={{this.upgradeData}}
        @xKey={{this.xKey}}
        @yKey={{this.yKey}}
      />
    </div>
    `);
    assert.dom('[data-test-line-chart]').exists('Chart is rendered');
    assert
      .dom('[data-test-line-chart="plot-point"]')
      .exists({ count: this.dataset.length }, `renders ${this.dataset.length} plot points`);
    assert
      .dom(`[data-test-line-chart="upgrade-month-2-18"]`)
      .exists({ count: 1 }, `upgrade data point 2-18 has yellow highlight`);
  });

  test('it renders tooltip', async function (assert) {
    assert.expect(1);
    const now = timestamp.now();
    const tooltipData = [
      {
        timestamp: formatRFC3339(subMonths(now, 4)),
        clients: 4,
        new_clients: {
          clients: 0,
        },
      },
      {
        timestamp: formatRFC3339(subMonths(now, 3)),
        clients: 8,
        new_clients: {
          clients: 4,
        },
      },
      {
        timestamp: formatRFC3339(subMonths(now, 2)),
        clients: 14,
        new_clients: {
          clients: 6,
        },
      },
      {
        timestamp: formatRFC3339(subMonths(now, 1)),
        clients: 20,
        new_clients: {
          clients: 4,
        },
      },
    ];
    this.set('dataset', tooltipData);
    this.set('upgradeData', [
      {
        id: '1.10.1',
        previousVersion: '1.9.2',
        timestampInstalled: formatRFC3339(subMonths(now, 2)),
      },
    ]);
    await render(hbs`
    <div class="chart-container-wide">
    <Clients::LineChart
      @dataset={{this.dataset}}
      @upgradeData={{this.upgradeData}}
    />
    </div>
    `);

    assert
      .dom('[data-test-hover-circle]')
      .exists({ count: tooltipData.length }, 'all data circles are rendered');

    // FLAKY after adding a11y testing, skip for now
    // for (const [i, bar] of tooltipHoverCircles.entries()) {
    //   await triggerEvent(bar, 'mouseover');
    //   const tooltip = document.querySelector('.ember-modal-dialog');
    //   const { month, clients, new_clients } = tooltipData[i];
    //   assert
    //     .dom(tooltip)
    //     .includesText(
    //       `${formatChartDate(month)} ${clients} total clients ${new_clients.clients} new clients`,
    //       `tooltip text is correct for ${month}`
    //     );
    // }
  });

  test('it fails gracefully when upgradeData is an object', async function (assert) {
    this.set('upgradeData', { some: 'object' });
    await render(hbs`
    <div class="chart-container-wide">
    <Clients::LineChart
    @dataset={{this.dataset}}
    @upgradeData={{this.upgradeData}}
    @xKey={{this.xKey}}
    @yKey={{this.yKey}}
    />
    </div>
    `);

    assert
      .dom('[data-test-line-chart="plot-point"]')
      .exists({ count: this.dataset.length }, 'chart still renders when upgradeData is not an array');
  });

  test('it fails gracefully when upgradeData has incorrect key names', async function (assert) {
    this.set('upgradeData', [{ incorrect: 'key names' }]);
    await render(hbs`
    <div class="chart-container-wide">
    <Clients::LineChart
    @dataset={{this.dataset}}
    @upgradeData={{this.upgradeData}}
    @xKey={{this.xKey}}
    @yKey={{this.yKey}}
    />
    </div>
    `);

    assert
      .dom('[data-test-line-chart="plot-point"]')
      .exists({ count: this.dataset.length }, 'chart still renders when upgradeData has incorrect keys');
  });

  test('it renders empty state when no dataset', async function (assert) {
    await render(hbs`
    <div class="chart-container-wide">
    <Clients::LineChart @noDataMessage="this is a custom message to explain why you're not seeing a line chart"/>
    </div>
    `);

    assert.dom('[data-test-component="empty-state"]').exists('renders empty state when no data');
    assert
      .dom('[data-test-empty-state-subtext]')
      .hasText(
        `this is a custom message to explain why you're not seeing a line chart`,
        'custom message renders'
      );
  });
});
