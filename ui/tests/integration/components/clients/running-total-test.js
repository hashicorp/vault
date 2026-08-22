/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, findAll, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import clientsHandler, { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { getUnixTime } from 'date-fns';
import { formatNumber } from 'core/helpers/format-number';
import timestamp from 'core/utils/timestamp';
import { CLIENT_COUNT, CHARTS } from 'vault/tests/helpers/clients/client-count-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { parseAPITimestamp } from 'core/utils/date-formatters';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import {
  destructureClientCounts,
  formatByMonths,
  formatByNamespace,
} from 'core/utils/client-counts/serializers';

const START_TIME = getUnixTime(LICENSE_START);

module('Integration | Component | clients/running-total', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    setRunOptions({
      rules: {
        // Carbon Charts renders path.bar elements with role="graphics-symbol" without aria-label.
        // This is a known Carbon Charts library limitation; the rule is suppressed here.
        'svg-img-alt': { enabled: false },
      },
    });
    this.flags = this.owner.lookup('service:flags');
    this.version = this.owner.lookup('service:version');
    this.flags.activatedFlags = ['secrets-sync'];
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    clientsHandler(this.server);
    const activityResponse = await this.owner
      .lookup('service:api')
      .sys.internalClientActivityReportCounts(undefined, getUnixTime(timestamp.now()), undefined, START_TIME);
    this.activity = {
      ...activityResponse,
      by_namespace: formatByNamespace(activityResponse.by_namespace),
      by_month: formatByMonths(activityResponse.months),
      total: destructureClientCounts(activityResponse.total),
    };
    this.byMonthClients = this.activity.by_month.map((d) => d.new_clients);

    this.renderComponent = async () => {
      await render(hbs`
      <Clients::RunningTotal
        @byMonthClients={{this.byMonthClients}}
        @runningTotals={{this.activity.total}}
      />
    `);
    };
  });

  test('it text for community versions', async function (assert) {
    this.version.type = 'community';
    await this.renderComponent();
    assert
      .dom(`${CLIENT_COUNT.card('Client usage trends')} p`)
      .hasText(
        'Number of clients in the date range by client type, and a breakdown of new clients per month during the date range.'
      );
  });

  // abbreviating "ent" so the test is not filtered out in CE repo runs
  test('it renders text for ent versions', async function (assert) {
    this.version.type = 'enterprise';
    await this.renderComponent();
    assert
      .dom(`${CLIENT_COUNT.card('Client usage trends')} p`)
      .hasText(
        'Number of clients in the billing period by client type, and a breakdown of new clients per month during the billing period.'
      );
  });

  test('it renders text for HVD managed versions', async function (assert) {
    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    await this.renderComponent();
    assert
      .dom(`${CLIENT_COUNT.card('Client usage trends')} p`)
      .hasText(
        'Number of total unique clients in the data period by client type, and total number of unique clients per month. The monthly total is the relevant billing metric.'
      );
  });

  test('it renders with full monthly activity data', async function (assert) {
    await this.renderComponent();

    assert.dom(CLIENT_COUNT.card('Client usage trends')).exists('running total component renders');
    assert.dom(GENERAL.inputByAttr('toggle view')).exists('chart toggle renders');

    const donutLegendItems = findAll(CHARTS.carbonLegendLabel('Client count and type distribution'));
    const expectedDonutLabels = [
      'Entity clients',
      'Non-entity clients',
      'ACME clients',
      'Secret sync clients',
    ];
    assert.strictEqual(
      donutLegendItems.length,
      expectedDonutLabels.length,
      'donut chart legend has correct number of items'
    );
    donutLegendItems.forEach((el, i) => {
      assert
        .dom(el)
        .hasText(expectedDonutLabels[i], `donut legend item ${i + 1} is "${expectedDonutLabels[i]}"`);
    });

    assert.dom('[data-test-chart="Client usage by month (simple)"]').exists('simple chart container renders');
    assert.dom('[data-test-chart="Client usage by month (simple)"] svg').exists('Carbon chart renders SVG');
    assert
      .dom(CHARTS.carbonLegendLabel('Client usage by month (simple)'))
      .hasText('New clients', 'simple chart legend shows the data key label');

    const xTicks = findAll(CHARTS.carbonXAxisTick('Client usage by month (simple)'));
    assert.strictEqual(xTicks.length, this.byMonthClients.length, 'x-axis has one tick per month');
    assert
      .dom(xTicks[0])
      .hasText(
        parseAPITimestamp(this.byMonthClients[0].timestamp, 'M/yy'),
        'first x-axis tick shows the first month'
      );
    assert
      .dom(CHARTS.carbonBar('Client usage by month (simple)'))
      .exists({ count: this.byMonthClients.length }, 'renders one bar per month');
  });

  test('it toggles to split chart by client type', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.inputByAttr('toggle view'));

    assert.dom(CLIENT_COUNT.card('Client usage trends')).exists('running total component renders');
    assert
      .dom('[data-test-chart="Client usage by month (stacked)"]')
      .exists('stacked chart container renders');
    assert.dom('[data-test-chart="Client usage by month (stacked)"] svg').exists('Carbon chart renders SVG');

    const legendLabels = findAll(CHARTS.carbonLegendLabel('Client usage by month (stacked)'));
    const expectedLabels = ['Entity clients', 'Non-entity clients', 'ACME clients', 'Secret sync clients'];
    assert.strictEqual(
      legendLabels.length,
      expectedLabels.length,
      'stacked chart legend has correct number of items'
    );
    legendLabels.forEach((el, i) => {
      assert.dom(el).hasText(expectedLabels[i], `legend item ${i + 1} is "${expectedLabels[i]}"`);
    });

    const months = this.byMonthClients.length;
    const groupCount = expectedLabels.length;
    assert
      .dom(CHARTS.carbonBar('Client usage by month (stacked)'))
      .exists({ count: months * groupCount }, `renders ${groupCount} bars per month`);
  });

  test('it renders when no monthly breakdown is available', async function (assert) {
    this.byMonthClients = [];
    await this.renderComponent();
    const expectedStats = {
      Entity: formatNumber([this.activity.total.entity_clients]),
      'Non-entity': formatNumber([this.activity.total.non_entity_clients]),
      ACME: formatNumber([this.activity.total.acme_clients]),
      'Secret sync': formatNumber([this.activity.total.secret_syncs]),
    };
    for (const label in expectedStats) {
      assert
        .dom(CLIENT_COUNT.statTextValue(label))
        .hasText(
          `${expectedStats[label]}`,
          `stat label: ${label} renders single month new clients: ${expectedStats[label]}`
        );
    }
    assert.dom(CHARTS.chart('Client usage by month (simple)')).doesNotExist('bar chart does not render');
    assert.dom(CLIENT_COUNT.statTextValue()).exists({ count: 5 }, 'renders 5 stat text containers');
  });

  test('it hides secret sync totals when feature is not activated', async function (assert) {
    this.flags.activatedFlags = [];
    // reset secret sync clients to 0
    this.byMonthClients = this.byMonthClients.map((obj) => ({ ...obj, secret_syncs: 0 }));

    await this.renderComponent();

    assert.dom(CLIENT_COUNT.card('Client usage trends')).exists('running total component renders');
    assert.dom('[data-test-chart="Client usage by month (simple)"]').exists('simple chart container renders');
    const donutLegendItems = findAll(CHARTS.carbonLegendLabel('Client count and type distribution'));
    const expectedDonutLabels = ['Entity clients', 'Non-entity clients', 'ACME clients'];
    assert.strictEqual(
      donutLegendItems.length,
      expectedDonutLabels.length,
      'donut legend has 3 items — secret sync clients is not included'
    );
    donutLegendItems.forEach((el, i) => {
      assert
        .dom(el)
        .hasText(expectedDonutLabels[i], `donut legend item ${i + 1} is "${expectedDonutLabels[i]}"`);
    });

    // check toggle view
    await click(GENERAL.inputByAttr('toggle view'));
    assert
      .dom('[data-test-chart="Client usage by month (stacked)"]')
      .exists('stacked chart container renders');
    assert.dom('[data-test-chart="Client usage by month (stacked)"] svg').exists('Carbon chart renders SVG');

    const legendLabels = findAll(CHARTS.carbonLegendLabel('Client usage by month (stacked)'));
    const expectedLabels = ['Entity clients', 'Non-entity clients', 'ACME clients'];
    assert.strictEqual(
      legendLabels.length,
      expectedLabels.length,
      'stacked legend has 3 items — secret sync clients is not included'
    );
    legendLabels.forEach((el, i) => {
      assert.dom(el).hasText(expectedLabels[i], `legend item ${i + 1} is "${expectedLabels[i]}"`);
    });

    const months = this.byMonthClients.length;
    const groupCount = expectedLabels.length;
    assert
      .dom(CHARTS.carbonBar('Client usage by month (stacked)'))
      .exists({ count: months * groupCount }, `renders ${groupCount} bars per month without secret sync`);
  });
});
