/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, find, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import clientsHandler, { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { getUnixTime } from 'date-fns';
import { findAll } from '@ember/test-helpers';
import { formatNumber } from 'core/helpers/format-number';
import timestamp from 'core/utils/timestamp';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { CLIENT_COUNT, CHARTS } from 'vault/tests/helpers/clients/client-count-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const START_TIME = getUnixTime(LICENSE_START);

module('Integration | Component | clients/running-total', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    sinon.replace(timestamp, 'now', sinon.fake.returns(STATIC_NOW));
    clientsHandler(this.server);
    const store = this.owner.lookup('service:store');
    const activityQuery = {
      start_time: { timestamp: START_TIME },
      end_time: { timestamp: getUnixTime(timestamp.now()) },
    };
    const activity = await store.queryRecord('clients/activity', activityQuery);
    this.byMonthNewClients = activity.byMonth.map((d) => d.new_clients);

    this.totalUsageCounts = activity.total;
    this.isSecretsSyncActivated = true;
    this.isHistoricalMonth = false;

    this.renderComponent = async () => {
      await render(hbs`
      <Clients::RunningTotal
        @isSecretsSyncActivated={{this.isSecretsSyncActivated}}
        @byMonthNewClients={{this.byMonthNewClients}}
        @runningTotals={{this.totalUsageCounts}}
        @upgradeData={{this.upgradesDuringActivity}}
        @isHistoricalMonth={{this.isHistoricalMonth}}
      />
    `);
    };
    // Fails on #ember-testing-container
    setRunOptions({
      rules: {
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });

  test('it renders with full monthly activity data', async function (assert) {
    await this.renderComponent();

    assert
      .dom(CLIENT_COUNT.card('Client usage trends for selected billing period'))
      .exists('running total component renders');
    assert.dom(CHARTS.chart('Client usage by month')).exists('bar chart renders');
    assert.dom(CHARTS.legend).hasText('New clients');
    const expectedColor = 'rgb(28, 52, 95)';
    const color = getComputedStyle(find(CHARTS.legendDot(1))).backgroundColor;
    assert.strictEqual(color, expectedColor, `actual color: ${color}, expected color: ${expectedColor}`);

    const expectedValues = {
      'New client total and type distribution': formatNumber([this.totalUsageCounts.clients]),
      Entity: formatNumber([this.totalUsageCounts.entity_clients]),
      'Non-entity': formatNumber([this.totalUsageCounts.non_entity_clients]),
      ACME: formatNumber([this.totalUsageCounts.acme_clients]),
      'Secret sync': formatNumber([this.totalUsageCounts.secret_syncs]),
    };
    for (const label in expectedValues) {
      assert
        .dom(CLIENT_COUNT.statTextValue(label))
        .hasText(
          `${expectedValues[label]}`,
          `stat label: ${label} renders correct total: ${expectedValues[label]}`
        );
    }

    // assert bar chart is correct
    findAll(CHARTS.xAxisLabel).forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${this.byMonthNewClients[i].month}`,
          `renders x-axis labels for bar chart: ${this.byMonthNewClients[i].month}`
        );
    });
    assert
      .dom(CHARTS.verticalBar)
      .exists({ count: this.byMonthNewClients.length }, 'renders correct number of bars ');
  });

  test('it toggles to split chart by client type', async function (assert) {
    await this.renderComponent();
    await click(GENERAL.inputByAttr('toggle view'));

    assert
      .dom(CLIENT_COUNT.card('Client usage trends for selected billing period'))
      .exists('running total component renders');
    assert.dom(CHARTS.chart('Client usage by month')).exists('bar chart renders');
    assert.dom(CHARTS.legend).hasText('Entity clients Non-entity clients Secret sync clients Acme clients');

    // assert each legend item is correct
    const expectedLegend = [
      { label: 'Entity clients', color: 'rgb(28, 52, 95)' },
      { label: 'Non-entity clients', color: 'rgb(6, 208, 146)' },
      { label: 'Secret sync clients', color: 'rgb(145, 28, 237)' },
      { label: 'Acme clients', color: 'rgb(2, 168, 239)' },
    ];

    findAll('.legend-item').forEach((e, i) => {
      const { label, color } = expectedLegend[i];
      assert.dom(e).hasText(label, `legend renders label: ${label}`);
      const dotColor = getComputedStyle(find(CHARTS.legendDot(i + 1))).backgroundColor;
      assert.strictEqual(dotColor, color, `actual color: ${dotColor}, expected color: ${color}`);
    });

    // assert bar chart is correct
    findAll(CHARTS.xAxisLabel).forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${this.byMonthNewClients[i].month}`,
          `renders x-axis labels for bar chart: ${this.byMonthNewClients[i].month}`
        );
    });

    const months = this.byMonthNewClients.length;
    const barsPerMonth = expectedLegend.length;
    assert
      .dom(CHARTS.verticalBar)
      .exists({ count: months * barsPerMonth }, `renders ${barsPerMonth} bars per month`);
  });

  test('it renders with single historical month data', async function (assert) {
    const singleMonthNew = this.byMonthNewClients[this.byMonthNewClients.length - 1];
    this.byMonthNewClients = [singleMonthNew];
    this.isHistoricalMonth = true;
    await this.renderComponent();
    const expectedStats = {
      'New clients': formatNumber([singleMonthNew.clients]),
      Entity: formatNumber([singleMonthNew.entity_clients]),
      'Non-entity': formatNumber([singleMonthNew.non_entity_clients]),
      ACME: formatNumber([singleMonthNew.acme_clients]),
      'Secret sync': formatNumber([singleMonthNew.secret_syncs]),
    };
    for (const label in expectedStats) {
      assert
        .dom(`[data-test-new] ${CLIENT_COUNT.statTextValue(label)}`)
        .hasText(
          `${expectedStats[label]}`,
          `stat label: ${label} renders single month new clients: ${expectedStats[label]}`
        );
    }
    assert.dom(CHARTS.chart('Client usage by month')).doesNotExist('bar chart does not render');
    assert.dom(CLIENT_COUNT.statTextValue()).exists({ count: 5 }, 'renders 5 stat text containers');
  });

  test('it hides secret sync totals when feature is not activated', async function (assert) {
    this.isSecretsSyncActivated = false;
    // reset secret sync clients to 0
    this.byMonthNewClients = this.byMonthNewClients.map((obj) => ({ ...obj, secret_syncs: 0 }));

    await this.renderComponent();

    assert
      .dom(CLIENT_COUNT.card('Client usage trends for selected billing period'))
      .exists('running total component renders');
    assert.dom(CHARTS.chart('Client usage by month')).exists('bar chart renders');
    assert.dom(CLIENT_COUNT.statTextValue('Entity')).exists();
    assert.dom(CLIENT_COUNT.statTextValue('Non-entity')).exists();
    assert.dom(CLIENT_COUNT.statTextValue('Secret sync')).doesNotExist('does not render secret syncs');

    // check toggle view
    await click(GENERAL.inputByAttr('toggle view'));
    assert
      .dom(CHARTS.legend)
      .hasText('Entity clients Non-entity clients Acme clients', 'legend does not include sync clients');

    // assert each legend item is correct
    const expectedLegend = [
      { label: 'Entity clients', color: 'rgb(28, 52, 95)' },
      { label: 'Non-entity clients', color: 'rgb(6, 208, 146)' },
      { label: 'Acme clients', color: 'rgb(2, 168, 239)' },
    ];

    findAll('.legend-item').forEach((e, i) => {
      const { label, color } = expectedLegend[i];
      assert.dom(e).hasText(label, `legend renders label: ${label}`);
      const dotColor = getComputedStyle(find(CHARTS.legendDot(i + 1))).backgroundColor;
      assert.strictEqual(dotColor, color, `actual color: ${dotColor}, expected color: ${color}`);
    });

    const months = this.byMonthNewClients.length;
    const barsPerMonth = expectedLegend.length;
    assert
      .dom(CHARTS.verticalBar)
      .exists({ count: months * barsPerMonth }, `renders ${barsPerMonth} bars per month`);
  });
});
