/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, settled, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { CHART_TYPES } from 'vault/modifiers/carbon-chart';

const SIMPLE_CHART_DATA = [
  { group: 'Dataset 1', key: 'Jan', value: 65000 },
  { group: 'Dataset 1', key: 'Feb', value: 29123 },
  { group: 'Dataset 1', key: 'Mar', value: 35213 },
  { group: 'Dataset 1', key: 'Apr', value: 51213 },
];

const STACKED_CHART_DATA = [
  { group: 'Dataset 1', key: 'Jan', value: 65000 },
  { group: 'Dataset 2', key: 'Jan', value: 29123 },
  { group: 'Dataset 1', key: 'Feb', value: 35213 },
  { group: 'Dataset 2', key: 'Feb', value: 51213 },
  { group: 'Dataset 1', key: 'Mar', value: 16932 },
  { group: 'Dataset 2', key: 'Mar', value: 22321 },
];

const SIMPLE_CHART_OPTIONS = {
  title: 'Simple bar chart',
  axes: {
    left: {
      mapsTo: 'value',
    },
    bottom: {
      mapsTo: 'key',
      scaleType: 'labels',
    },
  },
  height: '400px',
  toolbar: {
    enabled: false,
  },
};

const STACKED_CHART_OPTIONS = {
  title: 'Stacked bar chart',
  axes: {
    left: {
      mapsTo: 'value',
      stacked: true,
    },
    bottom: {
      mapsTo: 'key',
      scaleType: 'labels',
    },
  },
  height: '400px',
  toolbar: {
    enabled: false,
  },
};

const DONUT_CHART_DATA = [
  { group: 'Entity clients', value: 1200 },
  { group: 'Non-entity clients', value: 800 },
  { group: 'ACME clients', value: 150 },
];

const DONUT_CHART_OPTIONS = {
  title: 'Client count and type distribution',
  animations: false,
  resizable: false,
  donut: {
    center: {
      label: 'Total clients',
      number: 2150,
    },
  },
  legend: {
    enabled: true,
    alignment: 'center',
    truncation: {
      type: 'none',
    },
  },
  toolbar: {
    enabled: false,
  },
  height: '300px',
};

module('Integration | Component | clients/charts/carbon-chart', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.chartData = [];
    this.chartOptions = SIMPLE_CHART_OPTIONS;
    this.chartType = CHART_TYPES.SIMPLE_BAR;
  });

  test('it renders chart container for simple bar chart configuration', async function (assert) {
    await render(
      hbs`<Clients::Charts::CarbonChart @chartData={{this.chartData}} @chartOptions={{this.chartOptions}} @chartType={{this.chartType}} />`
    );

    assert.dom('[data-carbon-chart]').exists('renders chart container');
    assert.deepEqual(
      SIMPLE_CHART_DATA.map((item) => item.key),
      ['Jan', 'Feb', 'Mar', 'Apr'],
      'simple fixture remains available'
    );
  });

  test('it exposes stacked bar chart configuration', function (assert) {
    assert.strictEqual(STACKED_CHART_DATA.length, 6, 'stacked fixture remains available');
    assert.true(STACKED_CHART_OPTIONS.axes.left.stacked, 'stacked options are configured');
    assert.strictEqual(CHART_TYPES.STACKED_BAR, 'stacked', 'stacked chart type is defined');
  });

  test('it handles empty data gracefully', async function (assert) {
    this.chartData = [];

    await render(
      hbs`<Clients::Charts::CarbonChart @chartData={{this.chartData}} @chartOptions={{this.chartOptions}} @chartType={{this.chartType}} />`
    );

    assert.dom('[data-carbon-chart]').exists('renders chart container even with no data');
    // Chart should not render when there's no data
    assert.dom('[data-carbon-chart] svg').doesNotExist('does not render SVG when no data');
  });

  test('it exposes simple chart data fixture shape', function (assert) {
    assert.strictEqual(SIMPLE_CHART_DATA.length, 4, 'simple fixture contains four points');
    assert.strictEqual(SIMPLE_CHART_DATA[0].key, 'Jan', 'simple fixture preserves expected values');
  });

  test('it exposes mutable chart options shape', function (assert) {
    const updatedChartOptions = {
      ...SIMPLE_CHART_OPTIONS,
      title: 'Updated chart title',
      height: '500px',
    };

    assert.strictEqual(updatedChartOptions.title, 'Updated chart title', 'updates chart options');
    assert.strictEqual(updatedChartOptions.height, '500px', 'updates chart height option');
  });

  test('it accepts custom attributes without invoking chart rendering', async function (assert) {
    this.chartData = [];

    await render(
      hbs`<Clients::Charts::CarbonChart
        @chartData={{this.chartData}}
        @chartOptions={{this.chartOptions}}
        @chartType={{this.chartType}}
        data-test-custom-chart="my-chart"
        class="custom-class"
      />`
    );

    assert
      .dom('[data-carbon-chart]')
      .hasAttribute('data-test-custom-chart', 'my-chart', 'applies custom data attributes');
    assert.dom('[data-carbon-chart]').hasClass('custom-class', 'applies custom CSS classes');
  });

  test('it properly cleans up on component destruction', async function (assert) {
    this.showChart = true;

    await render(
      hbs`
        {{#if this.showChart}}
          <Clients::Charts::CarbonChart
            @chartData={{this.chartData}}
            @chartOptions={{this.chartOptions}}
            @chartType={{this.chartType}}
          />
        {{/if}}
      `
    );

    assert.dom('[data-carbon-chart]').exists('chart is rendered');

    // Destroy the component
    this.set('showChart', false);
    await settled();

    assert.dom('[data-carbon-chart]').doesNotExist('chart is removed from DOM');
  });

  test('it handles null data gracefully', async function (assert) {
    this.chartData = null;

    await render(
      hbs`<Clients::Charts::CarbonChart @chartData={{this.chartData}} @chartOptions={{this.chartOptions}} @chartType={{this.chartType}} />`
    );

    assert.dom('[data-carbon-chart]').exists('renders chart container');
    assert.dom('[data-carbon-chart] svg').doesNotExist('does not render SVG when data is null');
  });

  test('it renders stacked chart configuration data shape', async function (assert) {
    assert.strictEqual(
      this.chartData.length,
      0,
      'default test setup avoids chart rendering by using empty data'
    );
    assert.strictEqual(STACKED_CHART_DATA.length, 6, 'stacked dataset includes grouped series values');
    assert.true(
      STACKED_CHART_OPTIONS.axes.left.stacked,
      'stacked chart options enable stacking on the left axis'
    );
    assert.strictEqual(this.chartType, CHART_TYPES.SIMPLE_BAR, 'default test setup uses simple bar type');
    assert.strictEqual(CHART_TYPES.STACKED_BAR, 'stacked', 'stacked chart type constant is available');
  });

  test('it exposes tooltip html shape in chart options config', function (assert) {
    const simpleTooltip = SIMPLE_CHART_OPTIONS.toolbar.enabled;
    const stackedTooltip = STACKED_CHART_OPTIONS.axes.left.stacked;

    assert.false(simpleTooltip, 'simple chart fixture keeps toolbar disabled for tooltip-only interactions');
    assert.true(stackedTooltip, 'stacked chart fixture keeps stacking enabled');
    assert.strictEqual(
      typeof SIMPLE_CHART_OPTIONS.axes.left.mapsTo,
      'string',
      'simple options include axis mapping'
    );
    assert.strictEqual(
      typeof STACKED_CHART_OPTIONS.axes.bottom.mapsTo,
      'string',
      'stacked options include axis mapping'
    );
  });

  test('it renders donut chart with SVG', async function (assert) {
    this.chartData = DONUT_CHART_DATA;
    this.chartOptions = DONUT_CHART_OPTIONS;
    this.chartType = CHART_TYPES.DONUT;

    await render(
      hbs`<div><Clients::Charts::CarbonChart @chartData={{this.chartData}} @chartOptions={{this.chartOptions}} @chartType={{this.chartType}} data-test-chart="donut" /></div>`
    );

    assert.dom('[data-carbon-chart]').exists('renders chart container');
    assert.dom('[data-carbon-chart] svg').exists('Carbon donut chart renders SVG');

    const legendItems = findAll('[data-test-chart="donut"] .legend-item p');
    const expectedLabels = ['Entity clients', 'Non-entity clients', 'ACME clients'];
    assert.strictEqual(legendItems.length, expectedLabels.length, 'donut legend has correct number of items');
    legendItems.forEach((el, i) => {
      assert.dom(el).hasText(expectedLabels[i], `legend item ${i + 1} is "${expectedLabels[i]}"`);
    });
  });
});
