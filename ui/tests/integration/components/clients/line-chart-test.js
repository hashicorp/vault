/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import { setupRenderingTest } from 'ember-qunit';
import { find, render, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { format, formatRFC3339, subMonths } from 'date-fns';
import timestamp from 'core/utils/timestamp';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | clients/line-chart', function (hooks) {
  setupRenderingTest(hooks);
  hooks.before(function () {
    this.timestampStub = sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2018-04-03T14:15:30')));
  });
  hooks.beforeEach(function () {
    this.set('xKey', 'foo');
    this.set('yKey', 'bar');
    this.set('dataset', [
      {
        foo: '2018-04-03T14:15:30',
        bar: 4,
        expectedLabel: '4/18',
      },
      {
        foo: '2018-05-03T14:15:30',
        bar: 8,
        expectedLabel: '5/18',
      },
      {
        foo: '2018-06-03T14:15:30',
        bar: 14,
        expectedLabel: '6/18',
      },
      {
        foo: '2018-07-03T14:15:30',
        bar: 10,
        expectedLabel: '7/18',
      },
    ]);
  });

  test('it renders', async function (assert) {
    await render(hbs`
    <div class="chart-container-wide">
      <Clients::Charts::Line @dataset={{this.dataset}} @xKey={{this.xKey}} @yKey={{this.yKey}} />
      </div>
    `);

    assert.dom('[data-test-chart="line chart"]').exists('Chart is rendered');
    assert
      .dom('[data-test-plot-point]')
      .exists({ count: this.dataset.length }, `renders ${this.dataset.length} plot points`);
    findAll('[data-test-x-axis] text').forEach((e, i) => {
      // For some reason the first axis label is not rendered
      assert
        .dom(e)
        .hasText(
          `${this.dataset[i].expectedLabel}`,
          `renders x-axis label: ${this.dataset[i].expectedLabel}`
        );
    });
    assert.dom('[data-test-y-axis] text').hasText('0', `y-axis starts at 0`);
  });

  test('it renders upgrade data', async function (assert) {
    const now = this.timestampStub();
    this.set('dataset', [
      {
        foo: formatRFC3339(subMonths(now, 4)),
        bar: 4,
        month: format(subMonths(now, 4), 'M/yy'),
      },
      {
        foo: formatRFC3339(subMonths(now, 3)),
        bar: 8,
        month: format(subMonths(now, 3), 'M/yy'),
      },
      {
        foo: formatRFC3339(subMonths(now, 2)),
        bar: 14,
        month: format(subMonths(now, 2), 'M/yy'),
      },
      {
        foo: formatRFC3339(subMonths(now, 1)),
        bar: 10,
        month: format(subMonths(now, 1), 'M/yy'),
      },
    ]);
    this.set('upgradeData', [
      {
        id: '1.10.1',
        previousVersion: '1.9.2',
        timestampInstalled: formatRFC3339(subMonths(now, 2)),
      },
    ]);
    await render(hbs`
    <div class="chart-container-wide">
      <Clients::Charts::Line
        @dataset={{this.dataset}}
        @upgradeData={{this.upgradeData}}
        @xKey={{this.xKey}}
        @yKey={{this.yKey}}
      />
    </div>
    `);
    assert.dom('[data-test-chart="line chart"]').exists('Chart is rendered');
    assert
      .dom('[data-test-plot-point]')
      .exists({ count: this.dataset.length }, `renders ${this.dataset.length} plot points`);
    assert
      .dom(find(`[data-test-line-chart="upgrade-${this.dataset[2].month}"]`))
      .hasStyle(
        { fill: 'rgb(253, 238, 186)' },
        `upgrade data point ${this.dataset[2].month} has yellow highlight`
      );
  });

  test('it renders tooltip', async function (assert) {
    assert.expect(1);
    const now = this.timestampStub();
    const tooltipData = [
      {
        month: format(subMonths(now, 4), 'M/yy'),
        timestamp: formatRFC3339(subMonths(now, 4)),
        clients: 4,
        new_clients: {
          clients: 0,
        },
      },
      {
        month: format(subMonths(now, 3), 'M/yy'),
        timestamp: formatRFC3339(subMonths(now, 3)),
        clients: 8,
        new_clients: {
          clients: 4,
        },
      },
      {
        month: format(subMonths(now, 2), 'M/yy'),
        timestamp: formatRFC3339(subMonths(now, 2)),
        clients: 14,
        new_clients: {
          clients: 6,
        },
      },
      {
        month: format(subMonths(now, 1), 'M/yy'),
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
      <Clients::Charts::Line
        @dataset={{this.dataset}}
        @upgradeData={{this.upgradeData}}
      />
    </div>
    `);

    const tooltipHoverCircles = findAll('[data-test-hover-circle]');
    assert.strictEqual(tooltipHoverCircles.length, tooltipData.length, 'all data circles are rendered');

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

  test('it fails gracefully when data is not formatted correctly', async function (assert) {
    this.set('dataset', [
      {
        foo: 1,
        bar: 4,
      },
      {
        foo: 2,
        bar: 8,
      },
      {
        foo: 3,
        bar: 14,
      },
      {
        foo: 4,
        bar: 10,
      },
    ]);
    await render(hbs`
    <div class="chart-container-wide">
      <Clients::Charts::Line
        @dataset={{this.dataset}}
        @xKey={{this.xKey}}
        @yKey={{this.yKey}}
      />
    </div>
    `);

    assert.dom('[data-test-line-chart]').doesNotExist('Chart is not rendered');
    assert
      .dom('[data-test-component="empty-state"]')
      .hasText('No data to display', 'Shows empty state when time date is not formatted correctly');
  });

  test('it fails gracefully when upgradeData is an object', async function (assert) {
    this.set('upgradeData', { some: 'object' });
    await render(hbs`
    <div class="chart-container-wide">
    <Clients::Charts::Line
    @dataset={{this.dataset}}
    @upgradeData={{this.upgradeData}}
    @xKey={{this.xKey}}
    @yKey={{this.yKey}}
    />
    </div>
    `);

    assert
      .dom('[data-test-plot-point]')
      .exists({ count: this.dataset.length }, 'chart still renders when upgradeData is not an array');
  });

  test('it fails gracefully when upgradeData has incorrect key names', async function (assert) {
    this.set('upgradeData', [{ incorrect: 'key names' }]);
    await render(hbs`
    <div class="chart-container-wide">
    <Clients::Charts::Line
    @dataset={{this.dataset}}
    @upgradeData={{this.upgradeData}}
    @xKey={{this.xKey}}
    @yKey={{this.yKey}}
    />
    </div>
    `);

    assert
      .dom('[data-test-plot-point]')
      .exists({ count: this.dataset.length }, 'chart still renders when upgradeData has incorrect keys');
  });

  test('it renders empty state when no dataset', async function (assert) {
    await render(hbs`
    <div class="chart-container-wide">
    <Clients::Charts::Line @noDataMessage="this is a custom message to explain why you're not seeing a line chart"/>
    </div>
    `);

    assert.dom('[data-test-component="empty-state"]').exists('renders empty state when no data');
    assert
      .dom(GENERAL.emptyStateSubtitle)
      .hasText(
        `this is a custom message to explain why you're not seeing a line chart`,
        'custom message renders'
      );
  });

  test('it updates axis when dataset updates', async function (assert) {
    const datasets = {
      small: [
        {
          foo: '2020-04-01',
          bar: 4,
          month: '4/20',
        },
        {
          foo: '2020-05-01',
          bar: 8,
          month: '5/20',
        },
        {
          foo: '2020-06-01',
          bar: 1,
        },
        {
          foo: '2020-07-01',
          bar: 10,
        },
      ],
      large: [
        {
          foo: '2020-08-01',
          bar: 4586,
          month: '8/20',
        },
        {
          foo: '2020-09-01',
          bar: 8928,
          month: '9/20',
        },
        {
          foo: '2020-10-01',
          bar: 11948,
          month: '10/20',
        },
        {
          foo: '2020-11-01',
          bar: 16943,
          month: '11/20',
        },
      ],
      broken: [
        {
          foo: '2020-01-01',
          bar: null,
          month: '1/20',
        },
        {
          foo: '2020-02-01',
          bar: 0,
          month: '2/20',
        },
        {
          foo: '2020-03-01',
          bar: 22,
          month: '3/20',
        },
        {
          foo: '2020-04-01',
          bar: null,
          month: '4/20',
        },
        {
          foo: '2020-05-01',
          bar: 70,
          month: '5/20',
        },
        {
          foo: '2020-06-01',
          bar: 50,
          month: '6/20',
        },
      ],
    };
    this.set('dataset', datasets.small);
    await render(hbs`
    <div class="chart-container-wide">
      <Clients::Charts::Line
        @dataset={{this.dataset}}
        @upgradeData={{this.upgradeData}}
        @xKey={{this.xKey}}
        @yKey={{this.yKey}}
      />
    </div>
    `);
    assert.dom('[data-test-y-axis]').hasText('0 2 4 6 8 10', 'y-axis renders correctly for small values');
    assert
      .dom('[data-test-x-axis]')
      .hasText('4/20 5/20 6/20 7/20', 'x-axis renders correctly for small values');

    // Update to large dataset
    this.set('dataset', datasets.large);
    assert.dom('[data-test-y-axis]').hasText('0 5k 10k 15k', 'y-axis renders correctly for new large values');
    assert
      .dom('[data-test-x-axis]')
      .hasText('8/20 9/20 10/20 11/20', 'x-axis renders correctly for small values');

    // Update to broken dataset
    this.set('dataset', datasets.broken);
    assert.dom('[data-test-y-axis]').hasText('0 20 40 60', 'y-axis renders correctly for new broken values');
    assert
      .dom('[data-test-x-axis]')
      .hasText('1/20 2/20 3/20 4/20 5/20 6/20', 'x-axis renders correctly for small values');
    assert.dom('[data-test-hover-circle]').exists({ count: 4 }, 'only render circles for non-null values');

    assert
      .dom('[data-test-hover-circle="1/20"]')
      .doesNotExist('first month dot does not exist because value is null');
    assert
      .dom('[data-test-hover-circle="4/20"]')
      .doesNotExist('other null count month dot also does not render');
    // Note: the line should also show a gap, but this is difficult to test for
  });
});
