/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import clientsHandler, { LICENSE_START, STATIC_NOW } from 'vault/mirage/handlers/clients';
import sinon from 'sinon';
import { formatRFC3339, getUnixTime } from 'date-fns';
import { findAll } from '@ember/test-helpers';
import { formatNumber } from 'core/helpers/format-number';
import timestamp from 'core/utils/timestamp';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { CLIENT_COUNT } from 'vault/tests/helpers/clients/client-count-selectors';

const START_TIME = getUnixTime(LICENSE_START);

module('Integration | Component | clients/running-total', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => STATIC_NOW);
  });

  hooks.beforeEach(async function () {
    clientsHandler(this.server);
    const store = this.owner.lookup('service:store');
    const activityQuery = {
      start_time: { timestamp: START_TIME },
      end_time: { timestamp: getUnixTime(timestamp.now()) },
    };
    const activity = await store.queryRecord('clients/activity', activityQuery);
    this.byMonthActivity = activity.byMonth;
    this.newActivity = this.byMonthActivity.map((d) => d.new_clients);
    this.totalUsageCounts = activity.total;
    this.set('timestamp', formatRFC3339(timestamp.now()));
    this.set('chartLegend', [
      { label: 'entity clients', key: 'entity_clients' },
      { label: 'non-entity clients', key: 'non_entity_clients' },
    ]);
    this.isSecretsSyncActivated = true;
    this.isHistoricalMonth = false;

    this.renderComponent = async () => {
      await render(hbs`
      <Clients::RunningTotal
        @isSecretsSyncActivated={{this.isSecretsSyncActivated}}
        @byMonthActivityData={{this.byMonthActivity}}
        @runningTotals={{this.totalUsageCounts}}
        @upgradeData={{this.upgradesDuringActivity}}
        @responseTimestamp={{this.timestamp}}
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

  hooks.after(function () {
    timestamp.now.restore();
  });

  test('it renders with full monthly activity data', async function (assert) {
    await this.renderComponent();

    assert.dom(CLIENT_COUNT.chartContainer('Vault client counts')).exists('running total component renders');
    assert.dom(CLIENT_COUNT.charts.lineChart).exists('line chart renders');

    const expectedValues = {
      'Running client total': formatNumber([this.totalUsageCounts.clients]),
      'Entity clients': formatNumber([this.totalUsageCounts.entity_clients]),
      'Non-entity clients': formatNumber([this.totalUsageCounts.non_entity_clients]),
      'ACME clients': formatNumber([this.totalUsageCounts.acme_clients]),
      'Secrets sync clients': formatNumber([this.totalUsageCounts.secret_syncs]),
    };
    for (const label in expectedValues) {
      assert
        .dom(CLIENT_COUNT.charts.statTextValue(label))
        .hasText(
          `${expectedValues[label]}`,
          `stat label: ${label} renders correct total: ${expectedValues[label]}`
        );
    }

    // assert line chart is correct
    findAll(CLIENT_COUNT.charts.line.xAxisLabel).forEach((e, i) => {
      assert
        .dom(e)
        .hasText(
          `${this.byMonthActivity[i].month}`,
          `renders x-axis labels for line chart: ${this.byMonthActivity[i].month}`
        );
    });
    assert
      .dom(CLIENT_COUNT.charts.line.plotPoint)
      .exists(
        { count: this.byMonthActivity.filter((m) => m.clients).length },
        'renders correct number of plot points'
      );
  });

  test('it renders with no new monthly data', async function (assert) {
    this.byMonthActivity = this.byMonthActivity.map((d) => ({
      ...d,
      new_clients: { month: d.month },
    }));

    await this.renderComponent();

    assert.dom(CLIENT_COUNT.chartContainer('Vault client counts')).exists('running total component renders');
    assert.dom(CLIENT_COUNT.charts.lineChart).exists('line chart renders');

    const expectedValues = {
      'Entity clients': formatNumber([this.totalUsageCounts.entity_clients]),
      'Non-entity clients': formatNumber([this.totalUsageCounts.non_entity_clients]),
      'ACME clients': formatNumber([this.totalUsageCounts.acme_clients]),
      'Secrets sync clients': formatNumber([this.totalUsageCounts.secret_syncs]),
    };
    for (const label in expectedValues) {
      assert
        .dom(CLIENT_COUNT.charts.statTextValue(label))
        .hasText(
          `${expectedValues[label]}`,
          `stat label: ${label} renders correct total: ${expectedValues[label]}`
        );
    }
  });

  test('it renders with single historical month data', async function (assert) {
    const singleMonth = this.byMonthActivity[this.byMonthActivity.length - 1];
    const singleMonthNew = this.newActivity[this.newActivity.length - 1];
    this.byMonthActivity = [singleMonth];
    this.isHistoricalMonth = true;

    await this.renderComponent();

    let expectedStats = {
      'Total monthly clients': formatNumber([singleMonth.clients]),
      'Entity clients': formatNumber([singleMonth.entity_clients]),
      'Non-entity clients': formatNumber([singleMonth.non_entity_clients]),
      'ACME clients': formatNumber([singleMonth.acme_clients]),
      'Secrets sync clients': formatNumber([singleMonth.secret_syncs]),
    };
    for (const label in expectedStats) {
      assert
        .dom(`[data-test-total] ${CLIENT_COUNT.charts.statTextValue(label)}`)
        .hasText(
          `${expectedStats[label]}`,
          `stat label: ${label} renders single month total: ${expectedStats[label]}`
        );
    }

    expectedStats = {
      'New clients': formatNumber([singleMonthNew.clients]),
      'Entity clients': formatNumber([singleMonthNew.entity_clients]),
      'Non-entity clients': formatNumber([singleMonthNew.non_entity_clients]),
      'ACME clients': formatNumber([singleMonthNew.acme_clients]),
      'Secrets sync clients': formatNumber([singleMonthNew.secret_syncs]),
    };
    for (const label in expectedStats) {
      assert
        .dom(`[data-test-new] ${CLIENT_COUNT.charts.statTextValue(label)}`)
        .hasText(
          `${expectedStats[label]}`,
          `stat label: ${label} renders single month new clients: ${expectedStats[label]}`
        );
    }
    assert.dom(CLIENT_COUNT.charts.lineChart).doesNotExist('line chart does not render');
    assert.dom(CLIENT_COUNT.charts.statTextValue()).exists({ count: 10 }, 'renders 10 stat text containers');
  });

  test('it hides secret sync totals when feature is not activated', async function (assert) {
    this.isSecretsSyncActivated = false;

    await this.renderComponent();

    assert.dom(CLIENT_COUNT.chartContainer('Vault client counts')).exists('running total component renders');
    assert.dom(CLIENT_COUNT.charts.lineChart).exists('line chart renders');
    assert.dom(CLIENT_COUNT.charts.statTextValue('Entity clients')).exists();
    assert.dom(CLIENT_COUNT.charts.statTextValue('Non-entity clients')).exists();
    assert
      .dom(CLIENT_COUNT.charts.statTextValue('Secrets sync clients'))
      .doesNotExist('does not render secret syncs');
  });
});
